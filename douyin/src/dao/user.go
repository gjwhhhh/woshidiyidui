package dao

import "douyin/src/pojo/vo"

// TODO 编写 DAO

// IsExist 判断用户是否存在
func IsExist(username, password string) (int64, bool) {
	return 0, false
}

// GetUserInfo 获取当前用户的信息
func GetUserInfo(curUserId int64) vo.User {
	return vo.User{}
}

// GetOtherUserInfo 获取其它用户的信息
func GetOtherUserInfo(curUserId, otherUserId int64) (vo.User, bool) {
	return vo.User{}, false
}

// AddUser 新增用户
func AddUser(username, password string) (int64, error) {
	return 0, nil
}
