package dao

import (
	"douyin/src/global"
	"douyin/src/pojo/entity"
	"douyin/src/pojo/vo"
	"github.com/jinzhu/gorm"
)

// IsExist 判断用户是否存在
func IsExist(username, password string) (int64, bool) {
	var db = global.DBEngine
	var dy entity.DyUser
	err := db.Where("username=? AND password=?", username, password).Find(&dy).Error
	if err == gorm.ErrRecordNotFound {
		return 0, false
	}
	if err != nil {
		return 0, false
	}
	return dy.Id, true

}

// GetUserInfo 获取当前用户的信息
func GetUserInfo(curUserId int64) vo.User {
	var db = global.DBEngine
	var dyUser entity.DyUser
	err := db.Where("Id=?", curUserId).Find(&dyUser).Error
	if err == gorm.ErrRecordNotFound {
		return vo.User{}
	}

	var voUser = vo.User{
		Id:            dyUser.Id,
		Name:          dyUser.Username,
		FollowerCount: int64(dyUser.FollowCount),
		FollowCount:   int64(dyUser.FollowCount),
	}
	return voUser
}

// GetOtherUserInfo 获取其它用户的信息
// bool 返回用户是否存在
func GetOtherUserInfo(curUserId, otherUserId int64) (vo.User, bool) {
	var db = global.DBEngine
	//获取用户信息
	var dyUser entity.DyUser
	err := db.Where("Id=?", otherUserId).Find(&dyUser).Error
	if err == gorm.ErrRecordNotFound {
		return vo.User{}, false
	}
	var isFol = false
	var dyRela entity.DyRelation
	//获取是否关注信息
	err = db.Where("follower_id=? AND following_id=?", curUserId, otherUserId).Find(dyRela).Error
	if err == gorm.ErrRecordNotFound {
		isFol = false
	} else {
		isFol = true
	}
	var voUser = vo.User{
		Id:            dyUser.Id,
		Name:          dyUser.Username,
		FollowerCount: int64(dyUser.FollowCount),
		FollowCount:   int64(dyUser.FollowCount),
		IsFollow:      isFol,
	}

	return voUser, true
}

// AddUser 新增用户
func AddUser(username, password string) (int64, error) {
	var db = global.DBEngine
	user := entity.DyUser{
		Username: username,
		Password: password,
	}
	err := db.Create(&user).Error
	if err != nil {
		return 0, err
	}
	id, _ := IsExist(username, password)
	return id, nil
}
