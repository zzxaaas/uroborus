package service

import (
	"fmt"
	"github.com/go-redis/redis"
	"math/rand"
	"strconv"
	"uroborus/model"
	"uroborus/store"
)

type GroupService struct {
	groupStore *store.GroupStore
	rdsCli     *redis.Client
}

func NewGroupService(groupStore *store.GroupStore, rdsCli *redis.Client) *GroupService {
	return &GroupService{
		groupStore: groupStore,
		rdsCli:     rdsCli,
	}
}

func (s GroupService) Register(req *model.Group) error {
	req.Code = strconv.Itoa(rand.Intn(100000000))
	return s.groupStore.Save(req)
}

func (s GroupService) Find(req *model.Group) ([]model.Group, error) {
	groups, err := s.groupStore.FindCreateGroup(req)
	if err != nil {
		return nil, err
	}
	joinGroup, err := s.groupStore.FindJoinGroup(req)
	groups = append(groups, joinGroup...)
	for i := range groups {
		key := fmt.Sprintf("%s:%d", model.RedisKeyPrefix, groups[i].ID)
		groups[i].UserCount = s.rdsCli.SCard(key + model.KeyUserCountSuffix).Val()
		groups[i].ProjectCount = s.rdsCli.SCard(key + model.KeyProjCountSuffix).Val()
		if err != nil {
			fmt.Println(err)
		}
	}
	return groups, nil
}

func (s GroupService) Delete(req *model.Group) error {
	return s.groupStore.Delete(req)
}
