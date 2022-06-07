package dao

import (
	"database/sql"
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
	dy_favorite.user_id = ?`

type videoValid struct {
	Id            sql.NullInt64  `json:"id,omitempty"`
	Author        userValid      `json:"author"`
	PlayUrl       sql.NullString `json:"play_url" json:"play_url,omitempty"`
	CoverUrl      sql.NullString `json:"cover_url,omitempty"`
	FavoriteCount sql.NullInt64  `json:"favorite_count,omitempty"`
	CommentCount  sql.NullInt64  `json:"comment_count,omitempty"`
}

type userValid struct {
	Id            sql.NullInt64  `json:"id,omitempty"`
	Name          sql.NullString `json:"name,omitempty"`
	FollowCount   sql.NullInt64  `json:"follow_count,omitempty"`
	FollowerCount sql.NullInt64  `json:"follower_count,omitempty"`
}

// NewVoVideo 根据video获取vo.videoValid
func (v videoValid) NewVoVideo() *vo.Video {
	if !v.Id.Valid || !v.PlayUrl.Valid || !v.CoverUrl.Valid || !v.FavoriteCount.Valid || !v.CommentCount.Valid {
		return nil
	}
	return &vo.Video{
		Id:            v.Id.Int64,
		PlayUrl:       v.PlayUrl.String,
		CoverUrl:      v.CoverUrl.String,
		FavoriteCount: v.FavoriteCount.Int64,
		CommentCount:  v.CommentCount.Int64,
	}
}

// NewVoUser 根据User获取vo.userValid
func (u userValid) NewVoUser() *vo.User {
	if !u.Id.Valid || !u.Name.Valid || !u.FollowCount.Valid || !u.FollowerCount.Valid {
		return nil
	}
	return &vo.User{
		Id:            u.Id.Int64,
		Name:          u.Name.String,
		FollowCount:   u.FollowCount.Int64,
		FollowerCount: u.FollowerCount.Int64,
	}
}

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
	rows, err := db.DB().Query(FindFavoriteVideoListByUidSql, uid)
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
