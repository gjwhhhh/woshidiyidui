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

type FeedResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
	NextTime  int64   `json:"next_time,omitempty"`
}

// Feed 返回视频流数据
func Feed(c *gin.Context) {
	// 校验token
	token := c.Query("token")
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
					StatusMsg:  fmt.Sprintf("%s, no user corresponding to token", errcode.UnauthorizedTokenError.Msg()),
				},
			})
			return
		}
	}

	// 校验参数
	latestTime, err := strconv.ParseInt(c.Query("latest_time"), 10, 64)
	if latestTime == 0 {
		latestTime = util.GetTimeUnixNow()
	} else {
		if err != nil {
			c.JSON(http.StatusOK, UserResponse{
				Response: Response{
					StatusCode: 1,
					StatusMsg:  fmt.Sprintf("%s, %s", errcode.InvalidParams.Msg(), err.Error()),
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
				StatusMsg:  fmt.Sprintf("%s, %s", errcode.RequestFail.Msg(), err.Error()),
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
