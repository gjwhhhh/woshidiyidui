package entity

import (
	"time"
)

type DyVideo struct {
	Id            int64     `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	UserId        int64     `gorm:"column:user_id" json:"user_id"`
	PlayUrl       string    `gorm:"column:play_url" json:"play_url"`
	CoverUrl      string    `gorm:"column:cover_url" json:"cover_url"`
	FavoriteCount int       `gorm:"column:favorite_count;default:0" json:"favorite_count"`
	CommentCount  int       `gorm:"column:comment_count;default:0" json:"comment_count"`
	Title         string    `gorm:"column:title" json:"title"`
	CreateDate    time.Time `gorm:"column:create_date;default:CURRENT_TIMESTAMP;NOT NULL" json:"create_date"` // 创建时间
	UpdateDate    time.Time `gorm:"column:update_date;default:CURRENT_TIMESTAMP;NOT NULL" json:"update_date"` // 更新时间
	Isdeleted     int32     `gorm:"column:isdeleted;default:0" json:"isdeleted"`
}

func (m *DyVideo) TableName() string {
	return "dy_video"
}
