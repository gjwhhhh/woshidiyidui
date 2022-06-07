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

type UserListResponse struct {
	Response
	UserList []User `json:"user_list"`
}

// RelationAction no practical effect, just check if token is valid
func RelationAction(c *gin.Context) {
	// 获取参数
	toUserId, err1 := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	actionType, err2 := strconv.ParseInt(c.Query("action_type"), 10, 32)
	if err1 != nil || err2 != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "Illegal params",
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
	userId, exist := dao.IsExist(claims.Username, claims.Password)
	if !exist {
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  "No userInfo corresponding to token"})
		return
	}

	err = service.RelationAction(userId, toUserId, int32(actionType))
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  fmt.Sprintf("Relationship operation err, %s", err.Error())})
		return
	}

	c.JSON(http.StatusOK, Response{StatusCode: 0})
}

// FollowList all users have same follow list
func FollowList(c *gin.Context) {
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

	follows, err := service.FollowList(userId)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  fmt.Sprintf("Get following list err:%s", err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: follows,
	})
}

// FollowerList all users have same follower list
func FollowerList(c *gin.Context) {
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

	followers, err := service.FollowerList(userId)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  fmt.Sprintf("Get follower list err:%s", err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: followers,
	})
}
