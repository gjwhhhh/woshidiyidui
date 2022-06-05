package service

import (
	"douyin/src/dao"
	"douyin/src/global"
	"douyin/src/pojo/vo"
)

// Feed 获取视频流
func Feed(latestTime, userId int64) (feed []vo.Video, nextTime int64, err error) {
	// 获取视频feed流
	if userId == 0 {
		feed, err = feedByTime(latestTime)
	} else {
		feed, err = feedByTimeAndUId(latestTime, userId)
	}
	if err != nil {
		return nil, 0, err
	}

	// 获取nextTime
	len := len(feed)
	if len == 0 {
		return feed, latestTime, nil
	}

	nextTime, err = dao.GetVideoTimeById(feed[len-1].Id)
	if err != nil {
		return nil, 0, err
	}
	return feed, nextTime, nil
}

func feedByTimeAndUId(latestTime, userId int64) ([]vo.Video, error) {
	return dao.BatchVideoByTimeAndUId(latestTime, userId, global.DatabaseSetting.FeedPageSize)
}

func feedByTime(latestTime int64) ([]vo.Video, error) {
	return dao.BatchVideoByTime(latestTime, global.DatabaseSetting.FeedPageSize)
}
