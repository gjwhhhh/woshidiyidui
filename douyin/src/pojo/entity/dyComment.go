package entity

import (
	"database/sql"
	"douyin/src/pojo/vo"
	"time"
)

type DyComment struct {
	Id         sql.NullInt64  `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	VideoId    sql.NullInt64  `gorm:"column:video_id" json:"video_id"`
	UserId     sql.NullInt64  `gorm:"column:user_id" json:"user_id"`
	Content    sql.NullString `gorm:"column:content" json:"content"`
	CreateDate time.Time      `gorm:"column:create_date;default:CURRENT_TIMESTAMP;NOT NULL" json:"create_date"` // 创建时间
	UpdateDate time.Time      `gorm:"column:update_date;default:CURRENT_TIMESTAMP;NOT NULL" json:"update_date"` // 更新时间
	IsDeleted  sql.NullInt32  `gorm:"column:is_deleted;default:0" json:"is_deleted"`
}

func (c DyComment) NewVoComment() *vo.Comment {
	if !c.Id.Valid || !c.Content.Valid {
		return nil
	}
	return &vo.Comment{
		Id:         c.Id.Int64,
		Content:    c.Content.String,
		CreateDate: c.CreateDate.String(),
	}
}

func (m *DyComment) TableName() string {
	return "dy_comment"
}
