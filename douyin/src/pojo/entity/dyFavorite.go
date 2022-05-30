package entity

import "database/sql"

type DyFavorite struct {
	Id      int64         `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	UserId  sql.NullInt64 `gorm:"column:user_id" json:"user_id"`
	VideoId sql.NullInt64 `gorm:"column:video_id;default:0" json:"video_id"`
}

func (m *DyFavorite) TableName() string {
	return "`dy_favorite`"
}
