package service

import (
	"errors"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"os"
	"uroborus/common/auth"
	"uroborus/model"
	"uroborus/store"
)

type UserService struct {
	userStore *store.UserStore
}

func NewUserService(userStore *store.UserStore) *UserService {
	return &UserService{
		userStore: userStore,
	}
}

func (s UserService) Register(user *model.User) error {
	if has, err := s.Get(&model.User{
		Email: user.Email,
	}); err != nil {
		return err
	} else if has {
		return errors.New("该邮箱已注册")
	}
	if has, err := s.Get(&model.User{
		UserName: user.UserName,
	}); err != nil {
		return err
	} else if has {
		return errors.New("用户名重复")
	}
	if err := s.generatePassword(user); err != nil {
		return err
	}
	userRootPath := viper.GetString("user.rootPath") + user.UserName
	if err := os.MkdirAll(userRootPath, os.ModePerm); err != nil {
		return err
	}
	return s.userStore.Save(user)
}

func (s UserService) Login(condition *model.User) (string, error) {
	user := model.User{
		UserName: condition.UserName,
	}
	if has, err := s.Get(&user); err != nil {
		return "", err
	} else if !has {
		return "", errors.New("未注册")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(condition.Password)); err != nil {
		return "", errors.New("密码错误")
	}
	token, err := auth.SetToken(user.UserName)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s UserService) Get(user *model.User) (bool, error) {
	err := s.userStore.Get(user)
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s UserService) generatePassword(user *model.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hash)
	return nil
}
