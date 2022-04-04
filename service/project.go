package service

import (
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"math/rand"
	"mime/multipart"
	"os"
	"path"
	"strings"
	"uroborus/common"
	"uroborus/model"
	"uroborus/store"
)

type ProjectService struct {
	projectStore         *store.ProjectStore
	imageService         *BaseImageService
	deployHistoryService *DeployHistoryService
	gitService           *GitService
	groupService         *GroupService
	rdsCli               *redis.Client
}

func NewProjectService(
	projectStore *store.ProjectStore,
	imageService *BaseImageService,
	gitService *GitService,
	deployHistoryService *DeployHistoryService,
	groupService *GroupService,
	rdsCli *redis.Client,
) *ProjectService {
	return &ProjectService{
		projectStore:         projectStore,
		imageService:         imageService,
		gitService:           gitService,
		deployHistoryService: deployHistoryService,
		groupService:         groupService,
		rdsCli:               rdsCli,
	}
}

func (s ProjectService) Save(req *model.RegisterProjectReq) error {
	if req.Type == "tool" || req.Name == "" {
		req.Name = fmt.Sprintf("%s-%s", strings.Split(req.Image, ":")[0], req.UserName)
	}
	if req.Env != nil {
		req.Project.Env = strings.Join(req.Env, ",")
	}
	if req.Cmd != nil {
		req.Project.Command = strings.Join(req.Cmd, ",")
	}

	imageName := strings.Split(req.Image, ":")[0]
	if port, err := s.getImagePort(imageName); err != nil {
		return err
	} else if port != "" {
		req.Port = port
	}

	has, err := s.Get(&model.Project{
		Name: req.Name,
	})
	if err != nil {
		return err
	}

	if !has {
		if req.Type != "tool" {
			if err := s.initProjectPath(&req.Project); err != nil {
				return err
			}
			if err := s.clone(&req.Project); err != nil {
				return err
			}
		}
		if err := s.generatePort(&req.Project); err != nil {
			return err
		}
	}
	req.Project.AccessUrl = fmt.Sprintf("http://%s-%s.%s:%d", req.Project.Branch, req.Project.UserName, viper.GetString("baseUrl"), req.Project.BindPort)
	//req.Project.AccessUrl = fmt.Sprintf("http://121.196.214.245:%d", req.Project.BindPort)
	return s.projectStore.Save(&req.Project)
}

func (s ProjectService) Delete(req *model.Project) error {
	return s.projectStore.Delete(req)
}
func (s ProjectService) Update(req *model.Project) error {
	return s.projectStore.Update(req)
}

func (s ProjectService) generatePort(project *model.Project) error {
	for {
		project.BindPort = rand.Intn(65535)
		err := s.projectStore.Get(&model.Project{
			BindPort: project.BindPort,
		})
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				break
			}
			return err
		}
	}
	return nil
}

func (s ProjectService) getImagePort(name string) (string, error) {
	resp, err := s.imageService.Get(model.BaseImage{
		Name: name,
	})
	if err != nil {
		return "", err
	}
	if len(resp) == 0 {
		return "", errors.New("image not found")
	}
	return resp[0].Port, nil
}

func (s ProjectService) getBasePath(project *model.Project) string {
	root := viper.GetString("root")
	return fmt.Sprintf("%s/%s/%s/%s", root, project.UserName, project.Branch, project.Name)
}

func (s ProjectService) initProjectPath(project *model.Project) error {
	basePath := s.getBasePath(project)
	project.LocalRepo = basePath + model.RepoBasePath
	dockerfilePath := basePath + model.DockerfileBasePath
	logPath := basePath + model.LogfileBasePath
	imgPath := basePath + model.ImgfileBasePath
	for _, path := range []string{project.LocalRepo, dockerfilePath, logPath, imgPath} {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}

func (s ProjectService) Find(req model.Project) ([]model.GetProjectResp, error) {
	ans := make([]model.GetProjectResp, 0)
	projects, err := s.projectStore.Find(req)
	if err != nil {
		return nil, err
	}
	for _, project := range projects {
		deploy := model.DeployHistory{OriginId: project.ID}
		err := s.deployHistoryService.Get(&deploy)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		if req.Status != 0 && deploy.Status != req.Status {
			continue
		}
		resp := model.GetProjectResp{Project: project, DeployHistory: deploy}
		if project.GroupId != 0 {
			groups, err := s.groupService.Find(&model.Group{Model: model.Model{ID: project.GroupId}})
			if err != nil {
				return nil, err
			}
			resp.Group = groups[0]
		}
		ans = append(ans, resp)
	}
	return ans, nil
}

func (s ProjectService) clone(project *model.Project) error {
	if err := s.gitService.Clone(project.LocalRepo, false, &git.CloneOptions{
		URL:               project.RemoteRepo,
		ReferenceName:     plumbing.NewBranchReferenceName(project.Branch),
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	}); err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func (s ProjectService) Get(project *model.Project) (bool, error) {
	err := s.projectStore.Get(project)
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s ProjectService) GetGroupProjects(req *model.Group) ([]model.GetProjectResp, error) {
	projReq := model.Project{GroupId: req.ID, IsShow: model.ShowProjStatus, Status: model.DEPLOY_STATUS_SUCCESS, NeedGroup: true}
	resp, err := s.Find(projReq)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s ProjectService) RegisterGroupProject(req *model.RegisterGroupProjectReq, user string) error {
	groups, err := s.groupService.Find(&model.Group{Code: req.Code})
	if err != nil {
		return err
	}
	proj := model.Project{GroupId: groups[0].ID, IsShow: 1}
	proj.Name = req.Name
	if err := s.Update(&proj); err != nil {
		return err
	}
	key := fmt.Sprintf("%s:%d", model.RedisKeyPrefix, groups[0].ID)
	err = s.rdsCli.SAdd(key+model.KeyProjCountSuffix, req.Name).Err()
	proj.UserName = user
	projs, err := s.projectStore.Find(proj)
	if err != nil {
		return err
	}
	if len(projs) == 0 {
		s.rdsCli.SAdd(key+model.KeyUserCountSuffix, user)
	}
	return nil
}

func (s ProjectService) SaveProjectInfo(req *model.Project, file *multipart.FileHeader) (string, error) {
	if file == nil {
		req.SurfacePath = viper.GetString("default.surface")
	} else {
		fileExt := strings.ToLower(path.Ext(file.Filename))
		if fileExt != ".png" && fileExt != ".jpg" && fileExt != ".gif" && fileExt != ".jpeg" {
			return "", errors.New("上传失败!只允许png,jpg,gif,jpeg文件")
		}
		fileName := common.Md5(file)
		basePath := s.getBasePath(req)
		imgPath := basePath + model.ImgfileBasePath
		req.SurfacePath = fmt.Sprintf("%s%s%s", imgPath, fileName, fileExt)
	}
	if err := s.projectStore.Update(req); err != nil {
		return "", err
	}
	return req.SurfacePath, nil
}
