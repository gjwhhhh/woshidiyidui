package controller

import (
	"douyin/src/dao"
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
	token := c.PostForm("token")

	// 校验token
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

	// 判断用户是否存在
	userId, exist := dao.IsExist(claims.Username, claims.Password)
	if !exist {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "No userInfo corresponding to token",
			},
		})
		return
	}

	file, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  fmt.Sprintf("File parsing failed，%s", err.Error()),
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
			StatusMsg:  fmt.Sprintf("File stored on server failed，%s", err.Error()),
		})
		return
	}

	// 将文件上传的oss，返回视频url和封面url
	videoUrl, coverUrl, err := oss.UploadVideo(filename)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  fmt.Sprintf("File upload to cloud failed，%s", err.Error()),
		})
		return
	}

	// 将视频存储z
	title := c.PostForm("title")
	err = service.PublishVideo(userId, videoUrl, coverUrl, title)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  fmt.Sprintf("File storge in db failed，%s", err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  title + " uploaded successfully",
	})
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	// 获取c
	token := c.PostForm("token")
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  fmt.Sprintf("Illegal params, userId parse err:%s", err.Error()),
			},
		})
		return
	}

	// 校验token
	claims, err := util.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  fmt.Sprintf("Parse token err:%s", err.Error()),
		})
		return
	}

	videos, err := service.PublishList(claims.Username, claims.Password, userId)
	if err != nil {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  fmt.Sprintf("Get publish video list err:%s", err.Error()),
			},
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
