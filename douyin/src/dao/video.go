package dao

import (
	"database/sql"
	"douyin/src/global"
	"douyin/src/pkg/errcode"
	"douyin/src/pojo/entity"
	"douyin/src/pojo/vo"
	"douyin/src/util"
	"github.com/jinzhu/gorm"
	"time"
)

// BatchVideoByTimeAndUId 通过时间和用户id获取视频流，需要指定用户是否被关注
func BatchVideoByTimeAndUId(latestTime, userId int64, pageSize int) ([]vo.Video, error) {
	// 定义通道
	followerIdsChan := make(chan map[int64]struct{})
	errorChan := make(chan error)
	var followerIdMap map[int64]struct{}

	// 定义查询超时时间
	timer := time.NewTimer(time.Second)

	// 启用协程查询用户的关注列表
	go findFollowerIdsByFollowing(followerIdsChan, errorChan, userId)

	var db = global.DBEngine
	sqlTimeFormat := util.ParseTimeUnixToDbTime(latestTime)
	//查询视频list
	rows, err := db.Table("dy_video").Select([]string{"id", "user_id", "play_url", "cover_url", "favorite_count", "comment_count", "title"}).
		Where("create_date < ? and is_deleted = ?", sqlTimeFormat, 0).
		Order("create_date DESC").
		Limit(pageSize).
		Rows()
	videoVoList := make([]vo.Video, 0)
	if err != nil {
		return make([]vo.Video, 0), err
	}
	defer rows.Close()

	// 等待协程结果，超时直接返回Timeout
loop:
	for {
		select {
		case err = <-errorChan:
			return videoVoList, err
		case followerIdMap = <-followerIdsChan:
			break loop
		default:
			select {
			case <-timer.C:
				return videoVoList, errcode.TimeOutFail
			default:
				continue
			}
		}
	}

	// 读取数据并判断是否关注
	for rows.Next() {
		var video entity.DyVideo
		if err = rows.Scan(&video.Id, &video.UserId, &video.PlayUrl, &video.CoverUrl, &video.FavoriteCount, &video.CommentCount, &video.Title); err != nil {
			return videoVoList, err
		}

		voVideo := video.NewVoVideo()
		voUser := GetUserInfo(video.UserId.Int64)
		_, voUser.IsFollow = followerIdMap[video.UserId.Int64]
		voVideo.Author = *voUser

		videoVoList = append(videoVoList, *voVideo)
	}
	return videoVoList, nil
}

// BatchVideoByTime 通过时间获取视频流，不需要指定用户是否被关注
func BatchVideoByTime(latestTime int64, pageSize int) ([]vo.Video, error) {
	var db = global.DBEngine
	sqlTimeFormat := util.ParseTimeUnixToDbTime(latestTime)
	//查询视频list
	rows, err := db.Table("dy_video").Select([]string{"id", "user_id", "play_url", "cover_url", "favorite_count", "comment_count", "title"}).
		Where("create_date < ? and is_deleted = ?", sqlTimeFormat, 0).
		Order("create_date DESC").
		Limit(pageSize).
		Rows()
	videoVoList := make([]vo.Video, 0)
	if err != nil {
		return make([]vo.Video, 0), err
	}
	defer rows.Close()

	// 读取数据并判断是否关注
	for rows.Next() {
		var video entity.DyVideo
		if err = rows.Scan(&video.Id, &video.UserId, &video.PlayUrl, &video.CoverUrl, &video.FavoriteCount, &video.CommentCount, &video.Title); err != nil {
			return videoVoList, err
		}

		voVideo := video.NewVoVideo()
		voUser := GetUserInfo(video.UserId.Int64)
		voVideo.Author = *voUser

		videoVoList = append(videoVoList, *voVideo)
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
		UserId:   sql.NullInt64{Int64: userId, Valid: true},
		PlayUrl:  sql.NullString{String: videoUrl, Valid: true},
		CoverUrl: sql.NullString{String: coverUrl, Valid: true},
		Title:    sql.NullString{String: title, Valid: true},
	}
	var db = global.DBEngine
	err := db.Create(&dyVideo).Error
	if err != nil {
		return err
	}
	return nil
}

// BatchPublishVideoByUId 批量查询视频信息，略缩图简略信息
func BatchPublishVideoByUId(userId int64) ([]vo.Video, error) {
	var db = global.DBEngine
	//查询视频list
	rows, err := db.Table("dy_video").Select([]string{"id", "cover_url", "favorite_count"}).
		Where("user_id = ? and is_deleted = ?", userId, 0).
		Order("create_date DESC").
		Rows()
	videoVoList := make([]vo.Video, 0)
	if err != nil {
		return make([]vo.Video, 0), err
	}
	defer rows.Close()

	// 读取数据并判断是否关注
	for rows.Next() {
		var video entity.DyVideo
		if err = rows.Scan(&video.Id, &video.CoverUrl, &video.FavoriteCount); err != nil {
			return videoVoList, err
		}

		voVideo := &vo.Video{
			Id:            video.Id.Int64,
			CoverUrl:      video.CoverUrl.String,
			FavoriteCount: video.FavoriteCount.Int64,
		}

		videoVoList = append(videoVoList, *voVideo)
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
			Id:    video.Id.Int64,
			Title: video.Title.String,
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
