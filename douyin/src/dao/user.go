package dao

import (
	"database/sql"
	"douyin/src/cache/user_id"
	"douyin/src/cache/user_uname_pwd"
	"douyin/src/global"
	"douyin/src/pojo/entity"
	"douyin/src/pojo/vo"
	"github.com/jinzhu/gorm"
)

var UserCacheById = user_id.UserCacheConstructor(50)
var UserCacheByUnameAndPwd = user_uname_pwd.UserCacheConstructor(50)

// IsExist 判断用户是否存在
func IsExist(username, password string) (int64, bool) {
	var db = global.DBEngine
	// 查缓存,是否包含当前用户信息
	kw := username + password
	cacheUser := UserCacheByUnameAndPwd.Get(kw)
	// 包含，直接返回
	if cacheUser != nil {
		return cacheUser.Id.Int64, true
	}

	var dyUser entity.DyUser
	err := db.Where("username=? and password=? and is_deleted = ?", username, password, 0).First(&dyUser).Error
	if err == gorm.ErrRecordNotFound {
		return 0, false
	}

	// 添加到缓存中
	UserCacheByUnameAndPwd.Put(kw, &dyUser)
	if err != nil {
		return 0, false
	}
	return dyUser.Id.Int64, true
}

// IsExistByUName 判断用户是否存在
func IsExistByUName(username string) (int64, bool) {
	var db = global.DBEngine
	var dy entity.DyUser
	err := db.Where("username = ? and is_deleted = ?", username, 0).First(&dy).Error
	if err == gorm.ErrRecordNotFound {
		return 0, false
	}
	if err != nil {
		return 0, false
	}
	return dy.Id.Int64, true
}

// GetUserInfo 获取当前用户的信息
func GetUserInfo(curUserId int64) *vo.User {
	var db = global.DBEngine
	// 查缓存,是否包含当前用户信息
	cacheUser := UserCacheById.Get(curUserId)
	// 包含，直接返回
	if cacheUser != nil {
		return cacheUser.NewVoUser()
	}

	var dyUser entity.DyUser
	err := db.Where("id = ? and is_deleted = ?", curUserId, 0).First(&dyUser).Error
	if err == gorm.ErrRecordNotFound {
		return nil
	}

	// 添加到缓存中
	UserCacheById.Put(dyUser.Id.Int64, &dyUser)
	return dyUser.NewVoUser()
}

// GetOtherUserInfo 获取其它用户的信息
// bool 返回用户是否存在
func GetOtherUserInfo(curUserId, otherUserId int64) (*vo.User, error) {
	var db = global.DBEngine

	//使用缓存获取用户信息
	voUser := GetUserInfo(otherUserId)

	//获取是否关注信息
	var count int
	err := db.Table("dy_relation").Where("follower_id = ? and following_id = ? and is_deleted = ?", curUserId, otherUserId, 0).Count(&count).Error
	if err != nil {
		return nil, err
	}
	voUser.IsFollow = count != 0
	return voUser, nil
}

// AddUser 新增用户
func AddUser(username, password string) (int64, error) {
	var db = global.DBEngine
	user := entity.DyUser{
		Username: sql.NullString{String: username, Valid: true},
		Password: sql.NullString{String: password, Valid: true},
	}
	err := db.Create(&user).Error
	if err != nil {
		return 0, err
	}
	id, _ := IsExist(username, password)
	return id, nil
}
