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
	if err := db.Table("dy_relation").Where("follower_id = ? and is_deleted = ?", uid, 0).Find(&followerIds).Error; err == nil || err == gorm.ErrRecordNotFound {
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
const FindFollowerListByUidSql = `SELECT
	dy_user.id,
	dy_user.username,
	dy_user.follow_count,
	dy_user.follower_count 
FROM
	dy_relation
	LEFT JOIN dy_user ON  dy_relation.follower_id = dy_user.id
WHERE
	dy_relation.following_id = ?`

// FindFollowerList 查询某个用户粉丝列表
func FindFollowerList(userid int64) ([]vo.User, error) {
	var db = global.DBEngine
	rows, err := db.DB().Query(FindFollowerListByUidSql, userid)
	followerList := make([]vo.User, 0)
	if err != nil {
		return followerList, err
	}
	defer rows.Close()
	for rows.Next() {
		var user vo.User
		if err = rows.Scan(&user.Id, &user.Name, &user.FollowerCount, &user.FollowCount); err != nil {
			return followerList, err
		}
		followerList = append(followerList, user)
	}
	return followerList, nil

}

const FindFollowListByUidSql = `SELECT
	dy_user.id,
	dy_user.username,
	dy_user.follow_count,
	dy_user.follower_count 
FROM
	dy_relation
	LEFT JOIN dy_user ON  dy_relation.following_id = dy_user.id
WHERE
	dy_relation.follower_id = ?`

// FindFollowList 查询某个用户关注的人的列表
func FindFollowList(userid int64) ([]vo.User, error) {
	var db = global.DBEngine
	rows, err := db.DB().Query(FindFollowListByUidSql, userid)
	followerList := make([]vo.User, 0)
	if err != nil {
		return followerList, err
	}
	defer rows.Close()
	for rows.Next() {
		var user vo.User
		if err = rows.Scan(&user.Id, &user.Name, &user.FollowerCount, &user.FollowCount); err != nil {
			return followerList, err
		}
		followerList = append(followerList, user)
	}
	return followerList, nil

}
