package controller

import (
	"douyin/src/pkg/errcode"
	"douyin/src/service"
	"douyin/src/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User User `json:"User"`
}

// Register 用户注册
func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	// 调用业务注册用户，返回用户id和token
	userId, token, err := service.Register(username, password)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  fmt.Sprintf("%s, %s", errcode.OptionFail.Msg(), err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, UserLoginResponse{
		Response: Response{StatusCode: 0},
		UserId:   userId,
		Token:    token,
	})
}

// Login 用户登录
func Login(c *gin.Context) {
	// 获取参数
	username := c.Query("username")
	password := c.Query("password")

	// 获取token
	userId, token, err := service.GetToken(username, password)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  fmt.Sprintf("%s, %s", errcode.UnauthorizedTokenGenerate.Msg(), err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, UserLoginResponse{
		Response: Response{StatusCode: 0},
		UserId:   userId,
		Token:    token,
	})
}

// UserInfo 获取用户信息（可能是当前用户，也可能是其它用户，统一处理）
func UserInfo(c *gin.Context) {
	// 校验token
	token := c.Query("token")
	claims, err := util.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  fmt.Sprintf("%s, %s", errcode.UnauthorizedTokenError.Msg(), err.Error()),
			},
		})
		return
	}

	// 参数校验
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  fmt.Sprintf("%s, %s", errcode.InvalidParams.Msg(), err.Error()),
		})
		return
	}

	// 调用业务获取用户信息vo
	userInfo, err := service.GetUserInfo(claims.Username, claims.Password, userId)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  fmt.Sprintf("%s, %s", errcode.RequestFail.Msg(), err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, UserResponse{
		Response: Response{StatusCode: 0},
		User:     *userInfo,
	})
}
