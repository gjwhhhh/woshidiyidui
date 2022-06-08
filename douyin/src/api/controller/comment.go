package controller

import (
	"douyin/src/dao"
	"douyin/src/pkg/errcode"
	"douyin/src/service"
	"douyin/src/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	Response
	Comment *Comment `json:"comment,omitempty"`
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {
	// 校验token
	token := c.Query("token")
	claims, err := util.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  fmt.Sprintf("%s, %s", errcode.UnauthorizedTokenError.Msg(), err.Error()),
		})
		return
	}

	// 判断用户是否存在
	userId, exist := dao.IsExist(claims.Username, claims.Password)
	if !exist {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  fmt.Sprintf("%s, no user corresponding to token", errcode.UnauthorizedTokenError.Msg()),
		})
		return
	}

	// 参数校验
	videoId, err1 := strconv.ParseInt(c.Query("video_id"), 10, 64)
	actionType, err2 := strconv.ParseInt(c.Query("action_type"), 10, 32)
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  errcode.InvalidParams.Error(),
		})
		return
	}

	// 取消评论comment返回空结构体，新增评论返回comment
	comment, err := service.CommentAction(videoId, userId, int32(actionType), c)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  fmt.Sprintf("%s, %s", errcode.OptionFail.Msg(), err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, CommentActionResponse{
		Response: Response{StatusCode: 0},
		Comment:  comment},
	)
}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
	// 校验token
	token := c.Query("token")
	var curUserId int64
	if token != "" {
		claims, err := util.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  fmt.Sprintf("%s, %s", errcode.UnauthorizedTokenError.Msg(), err.Error()),
			})
			return
		}

		// 判断用户是否存在
		var exist bool
		curUserId, exist = dao.IsExist(claims.Username, claims.Password)
		if !exist {
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  fmt.Sprintf("%s, 没有关于此token的用户信息", errcode.UnauthorizedTokenError.Msg()),
			})
			return
		}
	}

	// 参数校验
	videoId, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  fmt.Sprintf("%s, %s", errcode.InvalidParams.Msg(), err.Error()),
		})
		return
	}

	// 调用业务
	comments, err := service.CommentList(videoId, curUserId)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  fmt.Sprintf("%s, %s", errcode.RequestFail.Msg(), err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0},
		CommentList: comments,
	})
}
