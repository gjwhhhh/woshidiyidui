package dao

import (
	"database/sql"
	"douyin/src/global"
	"douyin/src/pkg/errcode"
	"douyin/src/pojo/entity"
	"douyin/src/pojo/vo"
	"github.com/jinzhu/gorm"
	"time"
)

// FindFollowerIdsByUId 查询 uid 关注人id集合
func FindFollowerIdsByUId(uId int64) ([]int64, error) {
	db := global.DBEngine
	type Res struct {
		FollowingId sql.NullInt64 `gorm:"column:following_id;default:0" json:"following_id"`
	}
	var temp []Res
	if err := db.Table("dy_relation").Where("follower_id = ? and is_deleted = ?", uId, 0).Select("following_id").Scan(&temp).Error; err == nil || err == gorm.ErrRecordNotFound {
		res := make([]int64, len(temp))
		for i, r := range temp {
			res[i] = r.FollowingId.Int64
		}
		return res, nil
	} else {
		return make([]int64, 0), err
	}
}

// 查询 uid 关注人id集合
func findFollowerIdsByFollowing(followerIdsChan chan<- map[int64]struct{}, errorChan chan<- error, uid int64) {
	followerIds, err := FindFollowerIdsByUId(uid)
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

const FindFollowerListByUidSql = `SELECT
	dy_user.id,
	dy_user.username,
	dy_user.follow_count,
	dy_user.follower_count 
FROM
	dy_relation
	LEFT JOIN dy_user ON  dy_relation.follower_id = dy_user.id
WHERE
	dy_relation.following_id = ? AND dy_relation.is_deleted = ?`

// FindFollowerList 查询某个用户粉丝列表
func FindFollowerList(curUserId, userId int64) ([]vo.User, error) {
	// 定义通道
	followerIdsChan := make(chan map[int64]struct{})
	errorChan := make(chan error)
	var followerIdMap map[int64]struct{}

	// 定义查询超时时间
	timer := time.NewTimer(time.Second)

	// 启用协程查询用户的关注列表
	go findFollowerIdsByFollowing(followerIdsChan, errorChan, curUserId)

	var db = global.DBEngine
	rows, err := db.DB().Query(FindFollowerListByUidSql, userId, 0)
	followerList := make([]vo.User, 0)
	if err != nil {
		return followerList, err
	}
	defer rows.Close()

	// 等待协程结果，超时直接返回Timeout
loop:
	for {
		select {
		case err = <-errorChan:
			return followerList, err
		case followerIdMap = <-followerIdsChan:
			break loop
		default:
			select {
			case <-timer.C:
				return followerList, errcode.TimeOutFail
			default:
				continue
			}
		}
	}

	// 读取数据并判断是否关注
	for rows.Next() {
		var user entity.DyUser
		if err = rows.Scan(&user.Id, &user.Username, &user.FollowerCount, &user.FollowCount); err != nil {
			return followerList, err
		}

		voUser := user.NewVoUser()
		// TODO 对自己是否关注
		if user.Id.Int64 == curUserId {
			voUser.IsFollow = true
		} else {
			_, voUser.IsFollow = followerIdMap[user.Id.Int64]
		}
		followerList = append(followerList, *voUser)
	}
	return followerList, nil
}

// FindFollowerListWithoutLogin 未登录查询某个用户粉丝列表
func FindFollowerListWithoutLogin(userId int64) ([]vo.User, error) {
	var db = global.DBEngine
	rows, err := db.DB().Query(FindFollowerListByUidSql, userId, 0)
	followerList := make([]vo.User, 0)
	if err != nil {
		return followerList, err
	}
	defer rows.Close()

	for rows.Next() {
		var user entity.DyUser
		if err = rows.Scan(&user.Id, &user.Username, &user.FollowerCount, &user.FollowCount); err != nil {
			return followerList, err
		}

		voUser := user.NewVoUser()
		followerList = append(followerList, *voUser)
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
	dy_relation.follower_id = ? AND dy_relation.is_deleted = ?`

// FindFollowList 查询某个用户关注的人的列表
func FindFollowList(curUserId, userId int64) ([]vo.User, error) {
	// 定义通道
	followerIdsChan := make(chan map[int64]struct{})
	errorChan := make(chan error)
	var followerIdMap map[int64]struct{}

	// 定义查询超时时间
	timer := time.NewTimer(time.Second)

	// 启用协程查询用户的关注列表
	go findFollowerIdsByFollowing(followerIdsChan, errorChan, curUserId)

	var db = global.DBEngine
	rows, err := db.DB().Query(FindFollowListByUidSql, userId, 0)
	followList := make([]vo.User, 0)
	if err != nil {
		return followList, err
	}
	defer rows.Close()

	// 等待协程结果，超时直接返回Timeout
loop:
	for {
		select {
		case err = <-errorChan:
			return followList, err
		case followerIdMap = <-followerIdsChan:
			break loop
		default:
			select {
			case <-timer.C:
				return followList, errcode.TimeOutFail
			default:
				continue
			}
		}
	}

	for rows.Next() {
		var user entity.DyUser
		if err = rows.Scan(&user.Id, &user.Username, &user.FollowerCount, &user.FollowCount); err != nil {
			return followList, err
		}

		voUser := user.NewVoUser()
		// TODO 对自己是否关注
		if user.Id.Int64 == curUserId {
			voUser.IsFollow = true
		} else {
			_, voUser.IsFollow = followerIdMap[user.Id.Int64]
		}
		followList = append(followList, *voUser)
	}
	return followList, nil
}

// FindFollowListWithoutLogin 未登录查询某个用户关注的人的列表
func FindFollowListWithoutLogin(userId int64) ([]vo.User, error) {
	var db = global.DBEngine
	rows, err := db.DB().Query(FindFollowListByUidSql, userId, 0)
	followList := make([]vo.User, 0)
	if err != nil {
		return followList, err
	}
	defer rows.Close()

	for rows.Next() {
		var user entity.DyUser
		if err = rows.Scan(&user.Id, &user.Username, &user.FollowerCount, &user.FollowCount); err != nil {
			return followList, err
		}

		voUser := user.NewVoUser()
		followList = append(followList, *voUser)
	}
	return followList, nil
}
