package service

import (
	"fmt"
	"github.com/go-redis/redis"
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
	return s.groupStore.Save(req)
}

func (s GroupService) Find(req *model.Group) ([]model.Group, error) {
	groups, err := s.groupStore.Find(req)
	if err != nil {
		return nil, err
	}
	for i := range groups {
		key := fmt.Sprintf("%s:%s", model.RedisKeyPrefix, strconv.Itoa(int(req.ID)))
		groups[i].UserCount, _ = s.rdsCli.Get(key + model.KeyUserCountSuffix).Int()
		groups[i].ProjectCount, _ = s.rdsCli.Get(key + model.KeyProjCountSuffix).Int()
	}
	return groups, nil
}

func (s GroupService) Delete(req *model.Group) error {
	return s.groupStore.Delete(req)
}
