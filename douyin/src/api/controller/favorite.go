package controller

import (
	"douyin/src/api/controller/common"
	"douyin/src/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type (
	FavoriteActionReq struct {
		UserId     int64  `json:"user_id"`
		Token      string `json:"token"`
		VideoId    int64  `json:"video_id"`
		ActionType int32  `json:"action_type"`
	}
	FavoriteActionRep struct {
		Response
	}
)

type (
	FavoriteListReq struct {
	}
)

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(c *gin.Context) {
	//token := c.Query("token")
	favoriteActionReq := &FavoriteActionReq{}
	if err := c.BindJSON(favoriteActionReq); err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: common.REQUEST_PARAMETER_ERROR})
	}
	//if _, exist := usersLoginInfo[favoriteActionReq.Token]; !exist {
	//	c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: common.TOKEN_ERROR})
	//}
	favoriteActionDTO := service.FavoriteActionDTO{
		UserId:     favoriteActionReq.UserId,
		VideoId:    favoriteActionReq.VideoId,
		ActionType: favoriteActionReq.ActionType,
	}
	service.FavoriteAction(favoriteActionDTO)

}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: DemoVideos,
	})
}
