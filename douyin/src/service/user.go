package service

import (
	"douyin/src/dao"
	"douyin/src/pkg/errcode"
	"douyin/src/pojo/vo"
	"douyin/src/util"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
)

type auth struct {
	Username string `valid:"Required;MaxSize(32)"`
	Password string `valid:"Required;MaxSize(32)"`
}

var validation = validator.New()

func Register(username, password string) (int64, string, error) {
	// 参数校验
	err := validation.Struct(auth{
		Username: username,
		Password: password,
	})
	if err != nil {
		return 0, "", errors.New(fmt.Sprintf("%s, %s", errcode.InvalidParams.Error(), err.Error()))
	}

	// 判断用户是否存在
	_, exist := dao.IsExistByUName(username)
	if exist {
		return 0, "", errors.New("user already exist")
	}

	// 新增用户
	userId, err := dao.AddUser(username, password)
	if err != nil {
		return 0, "", errors.New(fmt.Sprintf("%s, %s", errcode.OptionFail.Error(), err.Error()))
	}

	// 生成token
	token, err := util.GenerateToken(username, password)
	if err != nil {
		return 0, "", errors.New(fmt.Sprintf("%s, %s", errcode.OptionFail.Error(), err.Error()))
	}
	return userId, token, nil
}

// GetToken 获取token
func GetToken(username, password string) (int64, string, error) {
	// 参数校验
	err := validation.Struct(auth{
		Username: username,
		Password: password,
	})
	if err != nil {
		return 0, "", err
	}

	// 判断用户是否存在
	userId, exist := dao.IsExist(username, password)
	if !exist {
		return 0, "", errcode.UserNotExistFail
	}

	// 生成token
	token, err := util.GenerateToken(username, password)
	if err != nil {
		return 0, "", err
	}
	return userId, token, nil
}

// GetUserInfo 获取用户信息
func GetUserInfo(username, password string, userId int64) (*vo.User, error) {
	// 判断用户是否存在
	curUserId, exist := dao.IsExist(username, password)
	if !exist {
		return nil, errcode.UserNotExistFail
	}

	// 获取用户信息
	if curUserId == userId { // 获取自己的用户信息
		return dao.GetUserInfo(curUserId), nil
	} else { // 获取他人的用户信息
		// 判断用户是否存在
		userInfo, err := dao.GetOtherUserInfo(curUserId, userId)
		if err != nil {
			return nil, err
		}
		return userInfo, nil
	}
}
