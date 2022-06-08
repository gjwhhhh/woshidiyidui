package controller

import (
	"douyin/src/dao"
	"douyin/src/pkg/errcode"
	"douyin/src/service"
	"douyin/src/util"
	"douyin/src/util/oss"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"strconv"
)

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}

// Publish 用户上传文件
func Publish(c *gin.Context) {
	// 校验token
	token := c.PostForm("token")
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
	title := c.PostForm("title")
	if title == "" {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  fmt.Sprintf("%s, titile can not be empty", errcode.InvalidParams.Error()),
		})
		return
	}

	// 获取数据
	file, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  fmt.Sprintf("%s, %s", errcode.FileReadFail.Msg(), err.Error()),
		})
		return
	}

	// 格式化文件名
	filename := fmt.Sprintf("%d_%s", userId, filepath.Base(file.Filename))
	// 指定存储本地文件的全路径
	localFilePath := filepath.Join(oss.LocalFilePathPrefix, filename)
	// 将文件存储到本地
	if err := c.SaveUploadedFile(file, localFilePath); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  fmt.Sprintf("%s, %s", errcode.FileUploadFail.Msg(), err.Error()),
		})
		return
	}

	// 将文件上传的oss，返回视频url和封面url
	videoUrl, coverUrl, err := oss.UploadVideo(filename)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  fmt.Sprintf("%s, %s", errcode.FileUploadFail.Msg(), err.Error()),
		})
		return
	}

	err = service.PublishVideo(userId, videoUrl, coverUrl, title)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  fmt.Sprintf("%s, %s", errcode.FileUploadFail.Msg(), err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  fmt.Sprintf("%s uploaded successfully", title),
	})
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	// 校验token
	var curUserId int64
	token := c.Query("token")
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
				StatusMsg:  fmt.Sprintf("%s, no user corresponding to token", errcode.UnauthorizedTokenError.Msg()),
			})
			return
		}
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

	videos, err := service.PublishList(curUserId, userId)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  fmt.Sprintf("%s, %s", errcode.RequestFail.Msg(), err.Error()),
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
