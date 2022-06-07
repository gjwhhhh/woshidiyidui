package dao

import (
	"database/sql"
	"douyin/src/global"
	"douyin/src/pkg/errcode"
	"douyin/src/pojo/entity"
	"douyin/src/pojo/vo"
	"github.com/jinzhu/gorm"
)

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
	tx.Table("dy_video").Where("id = ?", videoId).UpdateColumn("favorite_count", gorm.Expr("favorite_count + ?", 1))
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
	//dyFavorite := &entity.DyFavorite{
	//	UserId:  sql.NullInt64{Int64: userId, Valid: true},
	//	VideoId: sql.NullInt64{Int64: videoId, Valid: true},
	//}
	// 保存中间表
	if rowsAffected := tx.Where("user_id = ? and video_id = ?", userId, videoId).Delete(&entity.DyFavorite{}).RowsAffected; tx.Error != nil || rowsAffected != 1 {
		return errcode.UnLikeFail
	}
	tx.Table("dy_video").Where("id = ?", videoId).UpdateColumn("favorite_count", gorm.Expr("favorite_count - ?", 1))
	if tx.Error != nil {
		return errcode.UnLikeFail
	}
	tx.Commit()
	return nil
}

// AddComment 新增评论
// 开启事务
// 1.新增一条评论记录
// 2.在video表中评论数+1（是否需要CAS）
// 提交事务
func AddComment(videoId, userId int64, commentStr string) (*vo.Comment, error) {
	var db = global.DBEngine
	tx := db.Begin()
	defer tx.Callback()
	dyComment := &entity.DyComment{
		UserId:  sql.NullInt64{Int64: userId, Valid: true},
		VideoId: sql.NullInt64{Int64: videoId, Valid: true},
		Content: sql.NullString{String: commentStr, Valid: true},
	}
	// 保存中间表
	if rowsAffected := tx.Save(dyComment).RowsAffected; tx.Error != nil || rowsAffected != 1 {
		return nil, errcode.AddCommentFail
	}
	tx.Table("dy_video").Where("id = ? and isdeleted = ?", videoId, 0).UpdateColumn("comment_count", gorm.Expr("comment_count + ?", 1))
	if tx.Error != nil {
		return nil, errcode.AddCommentFail
	}
	tx.Commit()
	user := GetUserInfo(userId)
	comment := &vo.Comment{
		Id:         dyComment.Id,
		Content:    dyComment.Content.String,
		CreateDate: dyComment.CreateDate.String(),
		User:       *user,
	}
	return comment, nil
}

// DeleteComment 删除评论
// 开启事务
// 1.删除对应评论记录
// 2.在video表中评论数-1（是否需要CAS）
// 提交事务
func DeleteComment(commentId, videoId int64) error {
	var db = global.DBEngine
	tx := db.Begin()
	defer tx.Callback()
	dyComment := &entity.DyComment{
		VideoId: sql.NullInt64{Int64: videoId, Valid: true},
		Id:      commentId,
	}
	// 保存中间表
	if rowsAffected := tx.Delete(dyComment).RowsAffected; tx.Error != nil || rowsAffected != 1 {
		return errcode.DeleteCommentFail
	}
	tx.Table("dy_video").Where("id = ? and isdeleted = ?", videoId, 0).UpdateColumn("comment_count", gorm.Expr("comment_count - ?", 1))
	if tx.Error != nil {
		return errcode.DeleteCommentFail
	}
	tx.Commit()
	return nil
}
