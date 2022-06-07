package service

import (
	"douyin/src/dao"
	"douyin/src/pkg/errcode"
	"douyin/src/pojo/vo"
)

// PublishVideo 发布视频
func PublishVideo(userId int64, videoUrl, coverUrl, title string) error {
	return dao.AddVideo(userId, videoUrl, coverUrl, title)
}

// PublishList 获取发布列表
func PublishList(curUserId, userId int64) (videos []vo.Video, err error) {
	var userInfo *vo.User
	if curUserId == 0 || userId == curUserId { // 未登录获取，或者获取自己的发布列表
		videos, err = dao.BatchVideoByUId(userId)
		userInfo = dao.GetUserInfo(curUserId)
	} else { //获取别人的发布列表
		// 临时保存参数
		var exist bool
		// 获取别人的用户信息
		userInfo, exist = dao.GetOtherUserInfo(curUserId, userId)
		if !exist {
			return nil, errcode.UserNotExistFail
		}
		// 获取别人的发布视频列表信息
		videos, err = dao.BatchVideoByUIdAndOtherUId(curUserId, userId)
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
