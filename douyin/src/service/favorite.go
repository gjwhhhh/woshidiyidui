package service

import (
	"douyin/src/dao"
	"errors"
	"fmt"
)

const LikeOpt = 1   // 点赞
const UnLikeOpt = 2 // 取消点赞

func FavoriteAction(userId, videoId int64, actionType int32) error {
	if actionType == LikeOpt {
		return dao.AddFavorite(userId, videoId)
	} else if actionType == UnLikeOpt {
		return dao.DeleteFavorite(userId, videoId)
	} else {
		return errors.New(fmt.Sprintf("unsupported operation, action_type = %d", actionType))
	}
}
