package entity

import (
	"database/sql"
	"douyin/src/pojo/vo"
	"time"
)

type DyVideo struct {
	Id            sql.NullInt64  `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	UserId        sql.NullInt64  `gorm:"column:user_id" json:"user_id"`
	PlayUrl       sql.NullString `gorm:"column:play_url" json:"play_url"`
	CoverUrl      sql.NullString `gorm:"column:cover_url" json:"cover_url"`
	FavoriteCount sql.NullInt64  `gorm:"column:favorite_count;default:0" json:"favorite_count"`
	CommentCount  sql.NullInt64  `gorm:"column:comment_count;default:0" json:"comment_count"`
	Title         sql.NullString `gorm:"column:title" json:"title"`
	CreateDate    time.Time      `gorm:"column:create_date;default:CURRENT_TIMESTAMP;NOT NULL" json:"create_date"` // 创建时间
	UpdateDate    time.Time      `gorm:"column:update_date;default:CURRENT_TIMESTAMP;NOT NULL" json:"update_date"` // 更新时间
	IsDeleted     sql.NullInt32  `gorm:"column:is_deleted;default:0" json:"is_deleted"`
}

// NewVoVideo 根据video获取vo.videoValid
func (v DyVideo) NewVoVideo() *vo.Video {
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

func (m *DyVideo) TableName() string {
	return "dy_video"
}
