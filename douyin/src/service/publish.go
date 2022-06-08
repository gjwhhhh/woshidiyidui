package service

import (
	"douyin/src/dao"
	"douyin/src/pojo/vo"
)

// PublishVideo 发布视频
func PublishVideo(userId int64, videoUrl, coverUrl, title string) error {
	return dao.AddVideo(userId, videoUrl, coverUrl, title)
}

// PublishList 获取发布列表
func PublishList(curUserId, userId int64) (videos []vo.Video, err error) {
	var userInfo *vo.User
	if curUserId == 0 { // 未登录获取别人信息
		videos, err = dao.BatchPublishVideoByUId(userId)
		userInfo = dao.GetUserInfo(userId)
	} else if userId == curUserId { // 获取自己的发布列表
		videos, err = dao.BatchPublishVideoByUId(curUserId)
		userInfo = dao.GetUserInfo(curUserId)
	} else { //获取别人的发布列表
		// 临时保存参数
		// 获取别人的用户信息
		userInfo, err = dao.GetOtherUserInfo(curUserId, userId)
		if err != nil {
			return nil, err
		}
		// 获取别人的发布视频列表信息
		videos, err = dao.BatchPublishVideoByUId(userId)
	}
	if err != nil {
		return nil, err
	}

	// 封装作者
	for i, _ := range videos {
		videos[i].Author = *userInfo
	}
	return videos, nil
}
