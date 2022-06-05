package entity

import (
	"database/sql"
	"time"
)

type DyUser struct {
	Id            int64          `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	Username      sql.NullString `gorm:"column:username" json:"username"`
	Password      sql.NullString `gorm:"column:password" json:"password"`
	FollowCount   int            `gorm:"column:follow_count;default:0;NOT NULL" json:"follow_count"`
	FollowerCount sql.NullInt32  `gorm:"column:follower_count;default:0" json:"follower_count"`
	CreateDate    time.Time      `gorm:"column:create_date;default:CURRENT_TIMESTAMP;NOT NULL" json:"create_date"` // 创建时间
	UpdateDate    time.Time      `gorm:"column:update_date;default:CURRENT_TIMESTAMP;NOT NULL" json:"update_date"` // 更新时间
	Isdeleted     sql.NullInt32  `gorm:"column:isdeleted;default:0" json:"isdeleted"`
}

func (m *DyUser) TableName() string {
	return "`dy_user`"
}