package service

import (
	"douyin/src/dao"
	"douyin/src/pojo/vo"
	"errors"
)

// PublishVideo 发布视频
func PublishVideo(userId int64, videoUrl, coverUrl, title string) error {
	return dao.AddVideo(userId, videoUrl, coverUrl, title)
}

// PublishList 获取发布列表
func PublishList(username, password string, userId int64) (videos []vo.Video, err error) {
	userInfo, err := GetUserInfo(username, password, userId)
	if err != nil {
		return nil, err
	}

	if userId == userInfo.Id { // 获取自己的发布列表
		videos, err = dao.BatchVideoByUId(userId)
	} else { //获取别人的发布列表
		// 临时保存参数
		curUserId := userInfo.Id
		var exist bool
		// 获取别人的用户信息
		userInfo, exist = dao.GetOtherUserInfo(userInfo.Id, userId)
		if !exist {
			return nil, errors.New("user don't exist")
		}
		videos, err = dao.BatchVideoByUIdAndOtherUId(curUserId, userId)
	}
	if err != nil {
		return nil, err
	}

	// 封装作者
	for _, video := range videos {
		video.Author = userInfo
	}
	return videos, nil
}
