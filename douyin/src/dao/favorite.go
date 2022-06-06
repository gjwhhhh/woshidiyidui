package dao

import (
	"database/sql"
	"douyin/src/global"
	"douyin/src/pkg/errcode"
	"douyin/src/pojo/entity"
	"github.com/jinzhu/gorm"
)

type FavoriteActionDO struct {
	UserId  int64 `json:"user_id"`
	VideoId int64 `json:"video_id"`
}

const (
	ACTION_TYPE_LKIE = iota + 1
	ACTION_TYPE_UNLIKE
)

// FavoriteAction 更新点赞
func FavoriteAction(favoriteActionDO FavoriteActionDO, actionType int32) error {
	var db = global.DBEngine
	tx := db.Begin()
	defer tx.Callback()
	dyFavorite := &entity.DyFavorite{
		UserId:  sql.NullInt64{Int64: favoriteActionDO.UserId, Valid: true},
		VideoId: sql.NullInt64{Int64: favoriteActionDO.VideoId, Valid: true},
	}
	// 点赞操作
	if actionType == ACTION_TYPE_LKIE {
		// 保存中间表
		if rowsAffected := tx.Save(dyFavorite).RowsAffected; tx.Error != nil || rowsAffected != 1 {
			return errcode.LikeFail
		}
		tx.Table("dy_video").Where("id = ? and isdeleted = ?", favoriteActionDO.VideoId, 0).UpdateColumn("favorite_count", gorm.Expr("favorite_count + ?", 1))
		if tx.Error != nil {
			return errcode.LikeFail
		}
		tx.Commit()
		return nil
	}
	// 取消点赞
	if rowsAffected := tx.Delete(dyFavorite).RowsAffected; tx.Error != nil || rowsAffected != 1 {
		return errcode.LikeFail
	}
	tx.Table("dy_video").Where("id = ? and isdeleted = ?", favoriteActionDO.VideoId, 0).UpdateColumn("favorite_count", gorm.Expr("favorite_count + ?", 1))
	if tx.Error != nil {
		return errcode.LikeFail
	}
	tx.Commit()
	return nil
}
