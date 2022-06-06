package service

import "douyin/src/dao"

const (
	ACTION_TYPE_LKIE = iota + 1
	ACTION_TYPE_UNLIKE
)

type FavoriteActionDTO struct {
	UserId     int64 `json:"user_id"`
	VideoId    int64 `json:"video_id"`
	ActionType int32 `json:"action_type"`
}

func FavoriteAction(favoriteActionDTO FavoriteActionDTO) {
	favoriteActionDO := dao.FavoriteActionDO{VideoId: favoriteActionDTO.VideoId, UserId: favoriteActionDTO.UserId}
	dao.FavoriteAction(favoriteActionDO, favoriteActionDTO.ActionType)
}
