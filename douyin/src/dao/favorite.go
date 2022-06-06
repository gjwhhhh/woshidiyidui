package dao

import (
	"douyin/src/global"
	"douyin/src/pkg/errcode"
	"douyin/src/pojo/vo"
	"time"
)

const FIND_FAVORITE_VIDEO_LIST_BY_UID_SQL = `SELECT
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
	dy_favorite.user_id = ?`

func FindFavoriteVideoListByUId(uid int64) ([]vo.Video, error) {

	followerIdsChan := make(chan map[int64]struct{})
	errorChan := make(chan error)
	timer := time.NewTimer(time.Second)
	go findFollowerIdsByFollowing(followerIdsChan, errorChan, uid)
	var followerIdMap map[int64]struct{}
	db := global.DBEngine
	rows, err := db.DB().Query(FIND_FAVORITE_VIDEO_LIST_BY_UID_SQL, uid)
	videos := make([]vo.Video, 0)
	if err != nil {
		return videos, err
	}
	defer rows.Close()

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

	for rows.Next() {
		var video vo.Video
		var user vo.User
		if err = rows.Scan(&video.Id, &video.PlayUrl, &video.CoverUrl, &video.FavoriteCount, &video.CommentCount, &user.Id, &user.Name, &user.FollowCount, &user.FollowerCount); err != nil {
			return videos, err
		}
		video.IsFavorite = true
		_, user.IsFollow = followerIdMap[user.Id]
		video.Author = user
		videos = append(videos, video)
	}
	return videos, nil
}

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
