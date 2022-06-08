package dao

import (
	"douyin/src/global"
	"douyin/src/pkg/errcode"
	"douyin/src/pojo/vo"
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
		var video videoValid
		var user userValid
		if err = rows.Scan(&video.Id, &video.PlayUrl, &video.CoverUrl, &video.FavoriteCount, &video.CommentCount, &user.Id, &user.Name, &user.FollowCount, &user.FollowerCount); err != nil {
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
		video.Author = user
		videos = append(videos, *voVideo)
	}

	return videos, nil
}
