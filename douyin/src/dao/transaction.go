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
	tx.Table("dy_video").Where("id = ? and is_deleted = ?", videoId, 0).UpdateColumn("favorite_count", gorm.Expr("favorite_count + ?", 1))
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
	// 保存中间表
	if rowsAffected := tx.Table("dy_favorite").Where("user_id = ? and video_id = ? and is_deleted = ?", userId, videoId, 0).UpdateColumn("is_deleted", 1).RowsAffected; tx.Error != nil || rowsAffected != 1 {
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
	commentAdd := &sql.NullString{String: commentStr, Valid: true}
	dyComment := &entity.DyComment{
		UserId:  sql.NullInt64{Int64: userId, Valid: true},
		VideoId: sql.NullInt64{Int64: videoId, Valid: true},
		Content: *commentAdd,
	}
	// 保存中间表
	if rowsAffected := tx.Save(dyComment).RowsAffected; tx.Error != nil || rowsAffected != 1 {
		return nil, errcode.AddCommentFail
	}
	tx.Table("dy_video").Where("id = ? and is_deleted = ?", videoId, 0).UpdateColumn("comment_count", gorm.Expr("comment_count + ?", 1))
	if tx.Error != nil {
		return nil, errcode.AddCommentFail
	}
	tx.Commit()
	user := GetUserInfo(userId)
	comment := &vo.Comment{
		Id:         dyComment.Id.Int64,
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
	// 保存中间表
	if rowsAffected := tx.Table("dy_comment").Where("id = ? and is_deleted = ?", commentId, 0).UpdateColumn("is_deleted", 1).RowsAffected; tx.Error != nil || rowsAffected != 1 {
		return errcode.DeleteCommentFail
	}
	tx.Table("dy_video").Where("id = ? and is_deleted = ?", videoId, 0).UpdateColumn("comment_count", gorm.Expr("comment_count - ?", 1))
	if tx.Error != nil {
		return errcode.DeleteCommentFail
	}
	tx.Commit()
	return nil
}

// Follow 新增关注操作
// 开启事务
// 1.新增一条关注记录
// 2.在user表中关注数+1（是否需要CAS）
// 提交事务
func Follow(userId, toUserId int64) error {
	var db = global.DBEngine
	tx := db.Begin()
	defer tx.Callback()
	dyRelation := &entity.DyRelation{
		FollowerId:  sql.NullInt64{Int64: userId, Valid: true},
		FollowingId: sql.NullInt64{Int64: toUserId, Valid: true},
	}
	// 保存中间表
	if rowsAffected := tx.Save(dyRelation).RowsAffected; tx.Error != nil || rowsAffected != 1 {
		return errcode.FollowFail
	}
	tx.Table("dy_user").Where("id = ? and is_deleted = ?", userId).UpdateColumn("follow_count", gorm.Expr("follow_count + ?", 1))
	if tx.Error != nil {
		return errcode.FollowFail
	}
	tx.Table("dy_user").Where("id = ? and is_deleted = ?", toUserId).UpdateColumn("follower_count", gorm.Expr("follower_count + ?", 1))
	if tx.Error != nil {
		return errcode.FollowFail
	}
	tx.Commit()
	return nil
}

// UnFollow 取消关注操作
// 开启事务
// 1.删除对应关注记录
// 2.在video表中点赞数-1（是否需要CAS）
// 提交事务
func UnFollow(userId, toUserId int64) error {
	var db = global.DBEngine
	tx := db.Begin()
	defer tx.Callback()
	// 保存中间表
	if rowsAffected := tx.Table("dy_relation").Where("follower_id = ? and following_id = ? and is_deleted = ?", userId, toUserId, 0).UpdateColumn("is_deleted", 1).RowsAffected; tx.Error != nil || rowsAffected != 1 {
		return errcode.UnFollowFail
	}
	tx.Table("dy_user").Where("id = ? and is_deleted = ?", userId, 0).UpdateColumn("follow_count", gorm.Expr("follow_count - ?", 1))
	if tx.Error != nil {
		return errcode.UnLikeFail
	}
	tx.Table("dy_user").Where("id = ? and is_deleted = ?", toUserId, 0).UpdateColumn("follower_count", gorm.Expr("follow_count - ?", 1))
	if tx.Error != nil {
		return errcode.UnLikeFail
	}
	tx.Commit()
	return nil
}
