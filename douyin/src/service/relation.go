package service

import (
	"douyin/src/dao"
	"douyin/src/pojo/vo"
	"fmt"
	"github.com/pkg/errors"
	"sync"
)

const FollowOpt = 1   // 关注
const UnFollowOpt = 2 // 取消关注

var lock sync.Mutex

// FollowerList 查询某个用户粉丝列表
func FollowerList(curUserId, userId int64) ([]vo.User, error) {
	if curUserId == 0 {
		dao.FindFollowerListWithoutLogin(userId)
	}
	return dao.FindFollowerList(curUserId, userId)
}

// FollowList 查询某个用户关注的人的列表
func FollowList(curUserId, userId int64) ([]vo.User, error) {
	if curUserId == 0 {
		dao.FindFollowListWithoutLogin(userId)
	}
	return dao.FindFollowList(curUserId, userId)
}

// RelationAction 关系操作
func RelationAction(userId, toUserId int64, actionType int32) error {
	if userId == toUserId {
		return errors.New("不能关注自己")
	}
	if actionType == FollowOpt {
		err := dao.Follow(userId, toUserId)
		if err != nil {
			return err
		}
	} else if actionType == UnFollowOpt {
		err := dao.UnFollow(userId, toUserId)
		if err != nil {
			return err
		}
	} else {
		return errors.New(fmt.Sprintf("未知操作, action_type = %d", actionType))
	}
	// 删除缓存，保证并发安全
	lock.Lock()
	defer lock.Unlock()
	dao.UserCacheById.Delete(userId)
	dao.UserCacheById.Delete(toUserId)
	return nil
}
