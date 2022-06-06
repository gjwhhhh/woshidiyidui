package controller

import (
	"douyin/src/dao"
	"douyin/src/service"
	"douyin/src/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(c *gin.Context) {
	// 参数校验
	videoId, err1 := strconv.ParseInt(c.Query("video_Id"), 10, 64)
	actionType, err2 := strconv.ParseInt(c.Query("action_type"), 10, 32)
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  "Illegal params, parse err"})
		return
	}

	// 校验token
	token := c.Query("token")
	claims, err := util.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  fmt.Sprintf("Parse token err:%s", err.Error())})
		return
	}

	// 判断用户是否存在
	userId, exist := dao.IsExist(claims.Username, claims.Password)
	if !exist {
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  "No userInfo corresponding to token"})
		return
	}

	err = service.FavoriteAction(userId, videoId, int32(actionType))
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  fmt.Sprintf("Failed to like, %s", err.Error())})
		return
	}
	c.JSON(http.StatusOK, Response{StatusCode: 0})
}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	// 获取参数
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  fmt.Sprintf("Illegal params, userId parse err:%s", err.Error()),
		})
		return
	}

	// 校验token
	token := c.Query("token")
	claims, err := util.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  fmt.Sprintf("Parse token err:%s", err.Error())})
		return
	}

	// 判断用户是否存在
	_, exist := dao.IsExist(claims.Username, claims.Password)
	if !exist {
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  "No userInfo corresponding to token"})
		return
	}

	videos, err := service.FavoriteList(userId)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  fmt.Sprintf("Get favorite video list err:%s", err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: videos,
	})
}
