package service

import (
	"douyin/src/dao"
	"douyin/src/pojo/vo"
	"errors"
	"fmt"
)

const LikeOpt = 1   // 点赞
const UnLikeOpt = 2 // 取消点赞

// FavoriteAction 爱心操作
func FavoriteAction(userId, videoId int64, actionType int32) error {
	if actionType == LikeOpt {
		return dao.Like(userId, videoId)
	} else if actionType == UnLikeOpt {
		return dao.UnLike(userId, videoId)
	}
	return errors.New(fmt.Sprintf("未知操作, action_type = %d", actionType))
}

// FavoriteList 喜欢列表
func FavoriteList(userId int64) ([]vo.Video, error) {
	return dao.FindFavoriteVideoListByUId(userId)
}
