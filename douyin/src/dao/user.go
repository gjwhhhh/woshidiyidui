package dao

import (
	"douyin/src/cache"
	"douyin/src/global"
	"douyin/src/pojo/entity"
	"douyin/src/pojo/vo"
	"github.com/jinzhu/gorm"
)

var UserCache = cache.UserCacheConstructor(50)

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

// IsExistByUName 判断用户是否存在
func IsExistByUName(username string) (int64, bool) {
	var db = global.DBEngine
	var dy entity.DyUser
	err := db.Where("username=?", username).Find(&dy).Error
	if err == gorm.ErrRecordNotFound {
		return 0, false
	}
	if err != nil {
		return 0, false
	}
	return dy.Id, true
}

// GetUserInfo 获取当前用户的信息
func GetUserInfo(curUserId int64) *vo.User {
	var db = global.DBEngine
	// 查缓存,是否包含当前用户信息
	cacheUser := UserCache.Get(curUserId)
	// 包含，直接返回
	if cacheUser != nil {
		return &vo.User{
			Id:            cacheUser.Id,
			Name:          cacheUser.Username,
			FollowerCount: int64(cacheUser.FollowCount),
			FollowCount:   int64(cacheUser.FollowCount),
		}
	}

	var dyUser entity.DyUser
	err := db.Where("id=?", curUserId).Find(&dyUser).Error
	if err == gorm.ErrRecordNotFound {
		return nil
	}

	// 添加到缓存中
	UserCache.Put(dyUser.Id, &dyUser)
	return &vo.User{
		Id:            dyUser.Id,
		Name:          dyUser.Username,
		FollowerCount: int64(dyUser.FollowCount),
		FollowCount:   int64(dyUser.FollowCount),
	}
}

// GetOtherUserInfo 获取其它用户的信息
// bool 返回用户是否存在
func GetOtherUserInfo(curUserId, otherUserId int64) (*vo.User, bool) {
	var db = global.DBEngine
	//获取用户信息
	var dyUser entity.DyUser
	err := db.Where("id=?", otherUserId).Find(&dyUser).Error
	if err == gorm.ErrRecordNotFound {
		return nil, false
	}
	var isFol = false
	//获取是否关注信息
	err = db.Where("follower_id=? AND following_id=?", curUserId, otherUserId).Find(&entity.DyRelation{}).Error
	if err == gorm.ErrRecordNotFound {
		isFol = false
	} else {
		isFol = true
	}
	var voUser = &vo.User{
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
