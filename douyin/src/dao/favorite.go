package dao

import (
	"douyin/src/global"
	"douyin/src/pojo/vo"
)

const FindFavoriteVideoListByUIdSQL = `SELECT
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
	LEFT JOIN dy_user ON dy_favorite.user_id = dy_user.id
	LEFT JOIN dy_video ON dy_favorite.video_id = dy_video.id 
WHERE
	dy_user.id = ?`

func FindFavoriteVideoListByUId(uid int64) ([]vo.Video, error) {
	db := global.DBEngine
	rows, err := db.DB().Query(FindFavoriteVideoListByUIdSQL, uid)
	videos := make([]vo.Video, 0)
	if err != nil {
		return videos, err
	}
	for rows.Next() {
		var video vo.Video
		var user vo.User
		if err = rows.Scan(&video.Id, &video.PlayUrl, &video.CoverUrl, &video.FavoriteCount, &video.CommentCount, &user.Id, &user.Name, &user.FollowCount, &user.FollowerCount); err != nil {
			return videos, err
		}
		user.IsFollow = true
		video.IsFavorite = true
		video.Author = user
		videos = append(videos, video)
	}
	return videos, nil
}
