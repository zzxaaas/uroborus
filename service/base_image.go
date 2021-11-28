package service

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
	"strings"
	"uroborus/model"
	"uroborus/store"
)

type BaseImageService struct {
	baseImageStore *store.BaseImageStore
	restyClient    *resty.Client
}

func NewBaseImageService(baseImageStore *store.BaseImageStore) *BaseImageService {
	return &BaseImageService{
		baseImageStore: baseImageStore,
		restyClient:    resty.New(),
	}
}

func (s BaseImageService) Save(req model.BaseImage, header string) error {
	enterPoint := fmt.Sprintf("/v1/repositories/%s/tags", req.Name)
	tags := model.RegistryImageInfoResp{}
	resp, err := s.restyClient.NewRequest().
		SetHeader("User-Agent", header).
		SetResult(&tags.Tags).
		Get(viper.GetString("registry.url") + enterPoint)
	if err != nil {
		return err
	}
	if resp.IsError() {
		return fmt.Errorf("satus:%s, message:%s", resp.Status(), resp.String())
	}

	for _, tag := range tags.Tags {
		req.Tags += tag.Name + ","
	}
	req.Tags = strings.TrimRight(req.Tags, ",")
	return s.baseImageStore.Save(req)
}

func (s BaseImageService) Get(req model.BaseImage) ([]model.GetBaseImageResp, error) {
	resp := make([]model.GetBaseImageResp, 0)
	images, err := s.baseImageStore.Find(req)
	if err != nil {
		return nil, err
	}
	for _, image := range images {
		resp = append(resp, model.GetBaseImageResp{
			BaseImage: image,
			Tags:      strings.Split(image.Tags, ","),
		})
	}
	return resp, nil
}
