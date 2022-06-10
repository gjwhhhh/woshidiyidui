package middleware

import (
	"douyin/src/api/controller"
	"douyin/src/dao"
	"douyin/src/pkg/errcode"
	"douyin/src/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func JWT() gin.HandlerFunc {
	// 判断url 中的token是否正确
	return func(c *gin.Context) {
		// 获取token
		token := c.Query("token")
		// 判断token
		claims, err := util.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusOK, controller.Response{
				StatusCode: 1,
				StatusMsg:  fmt.Sprintf("%s, %s", errcode.UnauthorizedTokenError.Msg(), err.Error()),
			})
			c.Abort()
			return
		}
		// 判断用户是否存在
		_, exist := dao.IsExist(claims.Username, claims.Password)
		if !exist {
			c.JSON(http.StatusOK, controller.Response{
				StatusCode: 1,
				StatusMsg:  fmt.Sprintf("%s, no user corresponding to token", errcode.UnauthorizedTokenError.Msg()),
			})
			c.Abort()
			return
		}
		// 没有问题继续下一步操作
		c.Next()
	}
}
