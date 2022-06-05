package dao

import "douyin/src/pojo/vo"

// TODO 编写DAO

// BatchVideoByTimeAndUId 通过时间和用户id获取视频流，需要指定用户是否点赞
func BatchVideoByTimeAndUId(latestTime, userId int64, pageSize int) ([]vo.Video, error) {
	return make([]vo.Video, 0), nil
}

// BatchVideoByTime 通过时间获取视频流，不需要指定用户是否点赞
func BatchVideoByTime(latestTime int64, pageSize int) ([]vo.Video, error) {
	return make([]vo.Video, 0), nil
}

// GetVideoTimeById 通过视频id获取视频发布时间
func GetVideoTimeById(userId int64) (int64, error) {
	return 0, nil
}

// AddVideo 用户新增视频
func AddVideo(userId int64, videoUrl, coverUrl, title string) error {
	return nil
}
