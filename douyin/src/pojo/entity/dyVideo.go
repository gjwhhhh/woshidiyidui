package entity

import (
	"database/sql"
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

func (m *DyVideo) TableName() string {
	return "dy_video"
}
