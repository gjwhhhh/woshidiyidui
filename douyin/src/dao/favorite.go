package dao

import (
	"douyin/src/global"
	"douyin/src/pkg/errcode"
	"douyin/src/pojo/entity"
	"douyin/src/pojo/vo"
	"github.com/jinzhu/gorm"
	"time"
)

const FindFavoriteVideoListByUidSql = `SELECT
	dy_video.Id,
	dy_video.play_url,
	dy_video.cover_url,
	dy_video.favorite_count,
	dy_video.comment_count,
	dy_user.id,
	dy_user.username,
	dy_user.follow_count,
	dy_user.follower_count 
FROM
	dy_favorite
	LEFT JOIN dy_video ON dy_favorite.video_id = dy_video.id
	LEFT JOIN dy_user ON dy_video.user_id = dy_user.id
WHERE
	dy_favorite.user_id = ? AND dy_favorite.is_deleted = ?`

// FindFavoriteVideoListByUId 根据用户id获取点赞的视频
func FindFavoriteVideoListByUId(uid int64) ([]vo.Video, error) {
	followerIdsChan := make(chan map[int64]struct{})
	errorChan := make(chan error)
	var followerIdMap map[int64]struct{}

	// 定义查询超时时间
	timer := time.NewTimer(time.Second)

	// 启用协程查询用户的关注列表
	go findFollowerIdsByFollowing(followerIdsChan, errorChan, uid)

	// 根据用户id获取点赞的视频
	db := global.DBEngine
	rows, err := db.DB().Query(FindFavoriteVideoListByUidSql, uid, 0)
	videos := make([]vo.Video, 0)
	if err != nil {
		return videos, err
	}
	defer rows.Close()

	// 等待协程结果，超时直接返回Timeout
loop:
	for {
		select {
		case err = <-errorChan:
			return videos, err
		case followerIdMap = <-followerIdsChan:
			break loop
		default:
			select {
			case <-timer.C:
				return videos, errcode.TimeOutFail
			default:
				continue
			}
		}
	}

	// 读取数据并判断是否关注
	for rows.Next() {
		var video entity.DyVideo
		var user entity.DyUser
		if err = rows.Scan(&video.Id, &video.PlayUrl, &video.CoverUrl, &video.FavoriteCount, &video.CommentCount, &user.Id, &user.Username, &user.FollowCount, &user.FollowerCount); err != nil {
			return videos, err
		}

		voUser := user.NewVoUser()
		if voUser == nil {
			continue
		}

		voVideo := video.NewVoVideo()
		if voVideo == nil {
			continue
		}

		voVideo.IsFavorite = true
		_, voUser.IsFollow = followerIdMap[voUser.Id]
		voVideo.Author = *voUser

		videos = append(videos, *voVideo)
	}

	return videos, nil
}

// 查询 uid 点赞视频id集合
func findFavoriteVideoIdsByUId(videoIdsChan chan<- map[int64]struct{}, errorChan chan<- error, uid int64) {
	videoIds, err := FindFavoriteVideoIdsByUId(uid)
	if err != nil {
		errorChan <- err
		return
	}
	videoIdMap := make(map[int64]struct{})
	for _, video := range videoIds {
		videoIdMap[video] = struct{}{}
	}
	videoIdsChan <- videoIdMap
}

// FindFavoriteVideoIdsByUId 查询 uid 点赞视频id集合
func FindFavoriteVideoIdsByUId(uId int64) ([]int64, error) {
	db := global.DBEngine
	var res []int64
	if err := db.Table("dy_favorite").Where("user_id = ? and is_deleted = ?", uId, 0).Pluck("video_id", &res).Error; err == nil || err == gorm.ErrRecordNotFound {
		return res, nil
	} else {
		return make([]int64, 0), err
	}
}
