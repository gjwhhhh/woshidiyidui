package dao

import (
	"database/sql"
	"douyin/src/global"
	"douyin/src/pojo/entity"
	"douyin/src/pojo/vo"
	"douyin/src/util"
	"github.com/jinzhu/gorm"
)

// BatchVideoByTimeAndUId 通过时间和用户id获取视频流，需要指定用户是否点赞
func BatchVideoByTimeAndUId(latestTime, userId int64, pageSize int) ([]vo.Video, error) {
	var dyVideoList []entity.DyVideo
	var db = global.DBEngine
	sqlTimeFormat := util.ParseTimeUnixToDbTime(latestTime)
	//查询视频list
	err := db.Where("create_date < ? and is_deleted = ?", sqlTimeFormat, 0).Order("create_date DESC").Limit(pageSize).Find(&dyVideoList).Error
	if err != nil {
		return make([]vo.Video, 0), err
	}

	//转换voList
	var videoVoList = make([]vo.Video, 0)
	for i := 0; i < len(dyVideoList); i++ {
		// 通过缓存查询
		userTmp := GetUserInfo(dyVideoList[i].UserId.Int64)
		//db.Where("id=?", dyVideoList[i].UserId).Find(&userTmp)
		var dyFavorite entity.DyFavorite
		err := db.Where("user_id = ? and video_id = ? and is_deleted = ?", userId, dyVideoList[i].Id, 0).Find(&dyFavorite).Error
		var isFavorite = false
		if err != gorm.ErrRecordNotFound {
			isFavorite = true
		}
		videoTmp := vo.Video{
			Id:    dyVideoList[i].Id.Int64,
			Title: dyVideoList[i].Title.String,
			Author: vo.User{
				Id:            userTmp.Id,
				Name:          userTmp.Name,
				FollowCount:   userTmp.FollowerCount,
				FollowerCount: userTmp.FollowerCount,
				IsFollow:      false,
			},
			PlayUrl:       dyVideoList[i].PlayUrl.String,
			CoverUrl:      dyVideoList[i].CoverUrl.String,
			FavoriteCount: dyVideoList[i].FavoriteCount.Int64,
			CommentCount:  dyVideoList[i].CommentCount.Int64,
			IsFavorite:    isFavorite,
		}
		videoVoList = append(videoVoList, videoTmp)
	}
	return videoVoList, nil
}

// BatchVideoByTime 通过时间获取视频流，不需要指定用户是否点赞
func BatchVideoByTime(latestTime int64, pageSize int) ([]vo.Video, error) {
	var dyVideoList []entity.DyVideo
	var db = global.DBEngine
	sqlTimeFormat := util.ParseTimeUnixToDbTime(latestTime)
	//查询视频list
	err := db.Where("create_date < ? and is_deleted = ?", sqlTimeFormat, 0).Order("create_date DESC").Limit(pageSize).Find(&dyVideoList).Error
	if err != nil {
		return make([]vo.Video, 0), err
	}

	//转换voList
	var videoVoList = make([]vo.Video, 0)
	for i := 0; i < len(dyVideoList); i++ {
		// 通过缓存查询
		userTmp := GetUserInfo(dyVideoList[i].UserId.Int64)
		//db.Where("id=?", dyVideoList[i].UserId).Find(&userTmp)
		videoTmp := vo.Video{
			Id:    dyVideoList[i].Id.Int64,
			Title: dyVideoList[i].Title.String,
			Author: vo.User{
				Id:            userTmp.Id,
				Name:          userTmp.Name,
				FollowCount:   userTmp.FollowerCount,
				FollowerCount: userTmp.FollowerCount,
				IsFollow:      false,
			},
			PlayUrl:       dyVideoList[i].PlayUrl.String,
			CoverUrl:      dyVideoList[i].CoverUrl.String,
			FavoriteCount: dyVideoList[i].FavoriteCount.Int64,
			CommentCount:  dyVideoList[i].CommentCount.Int64,
			IsFavorite:    false,
		}
		videoVoList = append(videoVoList, videoTmp)
	}
	return videoVoList, nil
}

// GetVideoTimeById 通过视频id获取视频发布时间
func GetVideoTimeById(userId int64) (int64, error) {
	var dyVideo entity.DyVideo
	var db = global.DBEngine
	err := db.Where("id = ? and is_deleted = ?", userId, 0).Find(&dyVideo).Error
	if err != nil {
		return 0, err
	}
	// 返回毫秒
	return dyVideo.CreateDate.UnixMilli(), nil
}

// AddVideo 用户新增视频
func AddVideo(userId int64, videoUrl, coverUrl, title string) error {
	dyVideo := entity.DyVideo{
		UserId:   sql.NullInt64{Int64: userId},
		PlayUrl:  sql.NullString{String: videoUrl},
		CoverUrl: sql.NullString{String: coverUrl},
		Title:    sql.NullString{String: title},
	}
	var db = global.DBEngine
	err := db.Create(&dyVideo).Error
	if err != nil {
		return err
	}
	return nil
}

// BatchVideoByUId 批量查询视频信息，不需要指定是否点赞
func BatchVideoByUId(userId int64) ([]vo.Video, error) {
	var dyVideoList []entity.DyVideo
	var db = global.DBEngine
	err := db.Where("user_id = ? and is_deleted = ?", userId, 0).Find(&dyVideoList).Error
	if err != nil {
		return make([]vo.Video, 0), err
	}
	var videoVoList []vo.Video
	//遍历视频，然后循环中查询用户的信息，指定是否点赞
	for _, video := range dyVideoList {
		// 通过缓存查询
		userTmp := GetUserInfo(video.UserId.Int64)
		//db.Where("id=?", video.UserId).Find(&userTmp)

		videoVoTmp := vo.Video{
			Id: video.Id.Int64,
			Author: vo.User{
				Id:            userTmp.Id,
				Name:          userTmp.Name,
				FollowCount:   userTmp.FollowCount,
				FollowerCount: userTmp.FollowerCount,
				IsFollow:      false,
			},
			PlayUrl:       video.PlayUrl.String,
			CoverUrl:      video.CoverUrl.String,
			FavoriteCount: video.FavoriteCount.Int64,
			CommentCount:  video.CommentCount.Int64,
			IsFavorite:    false,
		}
		videoVoList = append(videoVoList, videoVoTmp)
	}
	return videoVoList, nil
}

// BatchVideoByUIdAndOtherUId 批量查询otherUId对应视频信息，需要通过curUserId指定是否点赞
func BatchVideoByUIdAndOtherUId(curUserId, otherUId int64) ([]vo.Video, error) {
	var dyVideoList []entity.DyVideo
	var db = global.DBEngine
	err := db.Where("user_id = ? and is_deleted = ?", otherUId, 0).Find(&dyVideoList).Error
	if err != nil {
		return make([]vo.Video, 0), err
	}
	var videoVoList []vo.Video
	for _, video := range dyVideoList {
		// 通过缓存查询
		userTmp := GetUserInfo(video.UserId.Int64)
		//db.Where("id=?", video.UserId).Find(&userTmp)
		var dyRelation entity.DyRelation
		err := db.Where("follower_id = ? and following_id = ? and is_deleted = ?", curUserId, otherUId, 0).Find(&dyRelation).Error
		var isFollow = false
		if err != gorm.ErrRecordNotFound {
			isFollow = true
		}
		var dyFavorite entity.DyFavorite
		err2 := db.Where("user_id = ? and video_id = ? and is_deleted = ?", curUserId, video.Id, 0).Find(&dyFavorite).Error
		var isFavorite = false
		if err2 != gorm.ErrRecordNotFound {
			isFavorite = true
		}
		videoVoTmp := vo.Video{
			Id: video.Id.Int64,
			Author: vo.User{
				Id:            userTmp.Id,
				Name:          userTmp.Name,
				FollowCount:   userTmp.FollowCount,
				FollowerCount: userTmp.FollowerCount,
				IsFollow:      isFollow,
			},
			PlayUrl:       video.PlayUrl.String,
			CoverUrl:      video.CoverUrl.String,
			FavoriteCount: video.FavoriteCount.Int64,
			CommentCount:  video.CommentCount.Int64,
			IsFavorite:    isFavorite,
		}
		videoVoList = append(videoVoList, videoVoTmp)
	}
	return videoVoList, nil
}
