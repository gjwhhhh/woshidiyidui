package service

import "douyin/src/dao"

// PublishVideo 发布视频
func PublishVideo(userId int64, videoUrl, coverUrl, title string) error {
	return dao.AddVideo(userId, videoUrl, coverUrl, title)
}
