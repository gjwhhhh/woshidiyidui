package service

import (
	"douyin/src/dao"
	"douyin/src/pojo/vo"
	"douyin/src/util"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

const AddCommentOpt = 1    // 新增评论
const DeleteCommentOpt = 2 // 删除评论

// CommentAction 评论操作
func CommentAction(videoId, userId int64, actionType int32, c *gin.Context) (*vo.Comment, error) {
	if actionType == AddCommentOpt {
		commentStr := c.Query("comment_text")
		// comment校验
		if commentStr == "" {
			return nil, errors.New("comment cannot be empty")
		}
		comment, err := dao.AddComment(videoId, userId, commentStr)
		if err != nil {
			return nil, err
		}
		time, err := util.ParseDbTimeToVoTime(comment.CreateDate)
		if err != nil {
			return nil, err
		}
		comment.CreateDate = time
		return comment, nil
	} else if actionType == DeleteCommentOpt {
		commentId, err := strconv.ParseInt(c.Query("comment_id"), 10, 64)

		if err != nil {
			return nil, errors.New("illegal param, parsing comment_id err")
		}
		return nil, dao.DeleteComment(commentId, videoId)

	}
	return nil, errors.New(fmt.Sprintf("unsupported operation, action_type = %d", actionType))
}

// CommentList 评论列表
func CommentList(videoId, userId int64) ([]vo.Comment, error) {
	comments, err := dao.FindCommentListByVideoIdAndUId(videoId, userId)
	if err != nil {
		return nil, err
	}
	for i, comment := range comments {
		time, _ := util.ParseRFC3339TimeToVoTime(comment.CreateDate)
		if err != nil {
			return nil, err
		}
		comments[i].CreateDate = time
	}
	return comments, nil
}
