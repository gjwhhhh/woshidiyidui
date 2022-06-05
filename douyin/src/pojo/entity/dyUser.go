package entity

import (
	"time"
)

type DyUser struct {
	Id            int64     `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	Username      string    `gorm:"column:username" json:"username"`
	Password      string    `gorm:"column:password" json:"password"`
	FollowCount   int       `gorm:"column:follow_count;default:0;NOT NULL" json:"follow_count"`
	FollowerCount int       `gorm:"column:follower_count;default:0" json:"follower_count"`
	CreateDate    time.Time `gorm:"column:create_date;default:CURRENT_TIMESTAMP;NOT NULL" json:"create_date"` // 创建时间
	UpdateDate    time.Time `gorm:"column:update_date;default:CURRENT_TIMESTAMP;NOT NULL" json:"update_date"` // 更新时间
	Isdeleted     int       `gorm:"column:isdeleted;default:0" json:"isdeleted"`
}

func (m *DyUser) TableName() string {
	return "dy_user"
}
