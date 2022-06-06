package service

import (
	"douyin/src/dao"
	"douyin/src/pojo/vo"
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
		commentStr := c.Query("comment")
		// comment校验
		if commentStr == "" {
			return nil, errors.New("comment cannot be empty")
		}
		return dao.AddComment(videoId, userId, commentStr)

	} else if actionType == DeleteCommentOpt {
		commentId, err := strconv.ParseInt(c.Query("comment_id"), 10, 64)
		if err != nil {
			return nil, errors.New("illegal param, parsing comment_id err")
		}
		return nil, dao.DeleteComment(commentId, videoId)

	}
	return nil, errors.New(fmt.Sprintf("unsupported operation, action_type = %d", actionType))
}
