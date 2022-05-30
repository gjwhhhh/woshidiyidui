package entity

import (
	"database/sql"
	"time"
)

type DyRelation struct {
	Id          int64         `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	FollowerId  sql.NullInt64 `gorm:"column:follower_id;default:0" json:"follower_id"`                          // 用户id
	FollowingId sql.NullInt64 `gorm:"column:following_id;default:0" json:"following_id"`                        // 对方的id
	CreateDate  time.Time     `gorm:"column:create_date;default:CURRENT_TIMESTAMP;NOT NULL" json:"create_date"` // 创建时间
	UpdateDate  time.Time     `gorm:"column:update_date;default:CURRENT_TIMESTAMP;NOT NULL" json:"update_date"` // 更新时间
	Isdeleted   sql.NullInt32 `gorm:"column:isdeleted;default:0" json:"isdeleted"`                              // 数据存在则用户关注了这个人
}

func (m *DyRelation) TableName() string {
	return "`dy_relation`"
}
