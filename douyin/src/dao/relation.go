package dao

import (
	"douyin/src/global"
	"douyin/src/pojo/vo"
	"github.com/jinzhu/gorm"
)

// FindFollowerIdsByFollowing 查询 uid 关注人id集合
func FindFollowerIdsByFollowing(uid int64) ([]int64, error) {
	db := global.DBEngine
	var followerIds []int64
	if err := db.Table("dy_relation").Where("follower_id = ? and isdeleted = ?", uid, 0).Find(&followerIds).Error; err == nil || err == gorm.ErrRecordNotFound {
		return followerIds, nil
	} else {
		return followerIds, err
	}
}

// 查询 uid 关注人id集合
func findFollowerIdsByFollowing(followerIdsChan chan<- map[int64]struct{}, errorChan chan<- error, uid int64) {
	followerIds, err := FindFollowerIdsByFollowing(uid)
	if err != nil {
		errorChan <- err
		return
	}
	followerIdMap := make(map[int64]struct{})
	for _, followerId := range followerIds {
		followerIdMap[followerId] = struct{}{}
	}
	followerIdsChan <- followerIdMap
}

// TODO dao完成

// FindFollowerList 查询某个用户粉丝列表
func FindFollowerList(userid int64) ([]vo.User, error) {
	return []vo.User{}, nil
}

// FindFollowList 查询某个用户关注的人的列表
func FindFollowList(userid int64) ([]vo.User, error) {
	return []vo.User{}, nil
}
