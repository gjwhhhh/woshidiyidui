package dao

import "douyin/src/pojo/vo"

// TODO 完成dao操作

// BatchFavoriteVideoByUId 批量查询用户点赞的视频信息，需要封装对应视频的作者信息
func BatchFavoriteVideoByUId(userId int64) ([]vo.Video, error) {
	return make([]vo.Video, 0), nil
}
