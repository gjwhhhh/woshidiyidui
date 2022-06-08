package service

import (
	"douyin/src/dao"
	"douyin/src/pojo/vo"
	"fmt"
	"github.com/pkg/errors"
)

const FollowOpt = 1   // 关注
const UnFollowOpt = 2 // 取消关注

// FollowerList 查询某个用户粉丝列表
func FollowerList(curUserId, userId int64) ([]vo.User, error) {
	return dao.FindFollowerList(curUserId, userId)
}

// FollowList 查询某个用户关注的人的列表
func FollowList(curUserId, userId int64) ([]vo.User, error) {
	return dao.FindFollowList(curUserId, userId)
}

// RelationAction 关系操作
func RelationAction(userId, toUserId int64, actionType int32) error {
	if userId == toUserId {
		return errors.New("can't follow on oneself")
	}
	if actionType == FollowOpt {
		return dao.Follow(userId, toUserId)
	} else if actionType == UnFollowOpt {
		return dao.UnFollow(userId, toUserId)
	}
	return errors.New(fmt.Sprintf("unsupported operation, action_type = %d", actionType))
}
