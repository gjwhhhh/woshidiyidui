package dao

import (
	"database/sql"
	"douyin/src/global"
	"douyin/src/pkg/errcode"
	"douyin/src/pojo/entity"
	"github.com/jinzhu/gorm"
)

// TODO 事务sql操作

// Like 新增点赞记录
// 开启事务
// 1.新增一条点赞记录
// 2.在video表中点赞数+1（是否需要CAS）
// 提交事务
func Like(userId, videoId int64) error {
	var db = global.DBEngine
	tx := db.Begin()
	defer tx.Callback()
	dyFavorite := &entity.DyFavorite{
		UserId:  sql.NullInt64{Int64: userId, Valid: true},
		VideoId: sql.NullInt64{Int64: videoId, Valid: true},
	}
	// 保存中间表
	if rowsAffected := tx.Save(dyFavorite).RowsAffected; tx.Error != nil || rowsAffected != 1 {
		return errcode.LikeFail
	}
	tx.Table("dy_video").Where("id = ? and isdeleted = ?", videoId, 0).UpdateColumn("favorite_count", gorm.Expr("favorite_count + ?", 1))
	if tx.Error != nil {
		return errcode.LikeFail
	}
	tx.Commit()
	return nil
}

// UnLike 删除记录
// 开启事务
// 1.删除对应点赞记录
// 2.在video表中点赞数-1（是否需要CAS）
// 提交事务
func UnLike(userId, videoId int64) error {
	var db = global.DBEngine
	tx := db.Begin()
	defer tx.Callback()
	dyFavorite := &entity.DyFavorite{
		UserId:  sql.NullInt64{Int64: userId, Valid: true},
		VideoId: sql.NullInt64{Int64: videoId, Valid: true},
	}
	// 保存中间表
	if rowsAffected := tx.Delete(dyFavorite).RowsAffected; tx.Error != nil || rowsAffected != 1 {
		return errcode.LikeFail
	}
	tx.Table("dy_video").Where("id = ? and isdeleted = ?", videoId, 0).UpdateColumn("favorite_count", gorm.Expr("favorite_count + ?", 1))
	if tx.Error != nil {
		return errcode.LikeFail
	}
	tx.Commit()
	return nil
}
