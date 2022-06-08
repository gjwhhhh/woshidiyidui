package entity

import (
	"database/sql"
	"douyin/src/pojo/vo"
	"time"
)

type DyUser struct {
	Id            sql.NullInt64  `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	Username      sql.NullString `gorm:"column:username" json:"username"`
	Password      sql.NullString `gorm:"column:password" json:"password"`
	FollowCount   sql.NullInt64  `gorm:"column:follow_count;default:0;NOT NULL" json:"follow_count"`
	FollowerCount sql.NullInt64  `gorm:"column:follower_count;default:0" json:"follower_count"`
	CreateDate    time.Time      `gorm:"column:create_date;default:CURRENT_TIMESTAMP;NOT NULL" json:"create_date"` // 创建时间
	UpdateDate    time.Time      `gorm:"column:update_date;default:CURRENT_TIMESTAMP;NOT NULL" json:"update_date"` // 更新时间
	IsDeleted     sql.NullInt32  `gorm:"column:is_deleted;default:0" json:"is_deleted"`
}

func (u DyUser) NewVoUser() *vo.User {
	if !u.Id.Valid || !u.Username.Valid || !u.FollowCount.Valid || !u.FollowerCount.Valid {
		return nil
	}
	return &vo.User{
		Id:            u.Id.Int64,
		Name:          u.Username.String,
		FollowCount:   u.FollowCount.Int64,
		FollowerCount: u.FollowerCount.Int64,
	}
}

func (m *DyUser) TableName() string {
	return "dy_user"
}
