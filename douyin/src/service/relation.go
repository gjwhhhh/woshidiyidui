package service

import (
	"douyin/src/dao"
	"douyin/src/pojo/vo"
)

//查询某个用户关注的人的列表
func FollowList(userid int64) ([]vo.User, error) {
	return dao.FindFollowList(userid)
}
