package dao

import (
	"database/sql"
	"douyin/src/pojo/vo"
)

type commentValid struct {
	Id         sql.NullInt64  `json:"id,omitempty"`
	User       userValid      `json:"user"`
	Content    sql.NullString `json:"content,omitempty"`
	CreateDate sql.NullString `json:"create_date,omitempty"`
}

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

func (c commentValid) NewVoComment() *vo.Comment {
	if !c.Id.Valid || !c.Content.Valid || !c.CreateDate.Valid {
		return nil
	}
	return &vo.Comment{
		Id:         c.Id.Int64,
		Content:    c.Content.String,
		CreateDate: c.CreateDate.String,
	}
}
