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

type FeedResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
	NextTime  int64   `json:"next_time,omitempty"`
}

// Feed 返回视频流数据
func Feed(c *gin.Context) {
	// 获取参数
	token := c.Query("token")
	latestTime, err := strconv.ParseInt(c.Query("latest_time"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  fmt.Sprintf("Illegal params, latestTime parse err:%s", err.Error()),
			},
		})
		return
	}

	// 校验token
	var userId int64
	if token != "" {
		claims, err := util.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusOK, UserResponse{
				Response: Response{
					StatusCode: 1,
					StatusMsg:  fmt.Sprintf("Parse token err:%s", err.Error()),
				},
			})
			return
		}
		// 获取用户id
		var exist bool
		userId, exist = dao.IsExist(claims.Username, claims.Password)
		if !exist {
			c.JSON(http.StatusOK, UserResponse{
				Response: Response{
					StatusCode: 1,
					StatusMsg:  "No userInfo corresponding to token",
				},
			})
			return
		}
	}

	// 调用service获取视频流
	feed, nextTime, err := service.Feed(latestTime, userId)

	if err != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  fmt.Sprintf("Getting video feed error:%s", err.Error()),
			},
		})
		return
	}

	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: feed,
		NextTime:  nextTime,
	})
}
