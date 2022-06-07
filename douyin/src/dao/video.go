package dao

import (
	"douyin/src/global"
	"douyin/src/pojo/entity"
	"douyin/src/pojo/vo"
	fmt "fmt"
	"github.com/jinzhu/gorm"
	"time"
)

// BatchVideoByTimeAndUId 通过时间和用户id获取视频流，需要指定用户是否点赞
func BatchVideoByTimeAndUId(latestTime, userId int64, pageSize int) ([]vo.Video, error) {
	var dyVideoList []entity.DyVideo
	var db = global.DBEngine
	sqlTimeFormat := timeunixTotimeString(latestTime)
	//查询视频list
	err := db.Where("create_date<?", sqlTimeFormat).Order("create_date DESC").Limit(pageSize).Find(&dyVideoList).Error
	if err != nil {
		return make([]vo.Video, 0), err
	}

	//转换voList
	var videoVoList = make([]vo.Video, 0)
	for i := 0; i < len(dyVideoList); i++ {
		var userTmp entity.DyUser
		db.Where("id=?", dyVideoList[i].UserId).Find(&userTmp)
		var dyFavorite entity.DyFavorite
		err := db.Where("user_id=? AND video_id=?", userId, dyVideoList[i].Id).Find(&dyFavorite).Error
		var isFavorite = false
		if err != gorm.ErrRecordNotFound {
			isFavorite = true
		}
		videoTmp := vo.Video{
			Id: dyVideoList[i].Id,
			Author: vo.User{
				Id:            userTmp.Id,
				Name:          userTmp.Username,
				FollowCount:   int64(userTmp.FollowerCount),
				FollowerCount: int64(userTmp.FollowerCount),
				IsFollow:      false,
			},
			PlayUrl:       dyVideoList[i].PlayUrl,
			CoverUrl:      dyVideoList[i].CoverUrl,
			FavoriteCount: int64(dyVideoList[i].FavoriteCount),
			CommentCount:  int64(dyVideoList[i].CommentCount),
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
	sqlTimeFormat := timeunixTotimeString(latestTime)
	//查询视频list
	err := db.Where("create_date<?", sqlTimeFormat).Order("create_date DESC").Limit(pageSize).Find(&dyVideoList).Error
	if err != nil {
		return make([]vo.Video, 0), err
	}

	//转换voList
	var videoVoList = make([]vo.Video, 0)
	for i := 0; i < len(dyVideoList); i++ {
		var userTmp entity.DyUser
		db.Where("id=?", dyVideoList[i].UserId).Find(&userTmp)
		videoTmp := vo.Video{
			Id: dyVideoList[i].Id,
			Author: vo.User{
				Id:            userTmp.Id,
				Name:          userTmp.Username,
				FollowCount:   int64(userTmp.FollowerCount),
				FollowerCount: int64(userTmp.FollowerCount),
				IsFollow:      false,
			},
			PlayUrl:       dyVideoList[i].PlayUrl,
			CoverUrl:      dyVideoList[i].CoverUrl,
			FavoriteCount: int64(dyVideoList[i].FavoriteCount),
			CommentCount:  int64(dyVideoList[i].CommentCount),
			IsFavorite:    false,
		}
		videoVoList = append(videoVoList, videoTmp)
	}
	return videoVoList, nil
}

//将时间戳转换为数据库时间格式
func timeunixTotimeString(lastestTime int64) string {
	timeObj := time.UnixMilli(lastestTime)

	year := timeObj.Year()     //年
	month := timeObj.Month()   //月
	day := timeObj.Day()       //日
	hour := timeObj.Hour()     //小时
	minute := timeObj.Minute() //分钟
	second := timeObj.Second() //秒
	timeString := fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d\n", year, month, day, hour, minute, second)
	return timeString
}

// GetVideoTimeById 通过视频id获取视频发布时间
func GetVideoTimeById(userId int64) (int64, error) {
	var dyVideo entity.DyVideo
	var db = global.DBEngine
	err := db.Where("id=?", userId).Find(&dyVideo).Error
	if err != nil {
		return 0, err
	}
	return dyVideo.CreateDate.Unix(), nil
}

// AddVideo 用户新增视频
func AddVideo(userId int64, videoUrl, coverUrl, title string) error {
	dyVideo := entity.DyVideo{
		UserId:   userId,
		PlayUrl:  videoUrl,
		CoverUrl: coverUrl,
		Title:    title,
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
	err := db.Where("user_id=?", userId).Find(&dyVideoList).Error
	if err != nil {
		return make([]vo.Video, 0), err
	}
	var videoVoList []vo.Video
	//遍历视频，然后循环中查询用户的信息，指定是否点赞
	for _, video := range dyVideoList {
		var userTmp entity.DyUser
		db.Where("id=?", video.UserId).Find(&userTmp)

		videoVoTmp := vo.Video{
			Id: video.Id,
			Author: vo.User{
				Id:            userTmp.Id,
				Name:          userTmp.Username,
				FollowCount:   int64(userTmp.FollowCount),
				FollowerCount: int64(userTmp.FollowerCount),
				IsFollow:      false,
			},
			PlayUrl:       video.PlayUrl,
			CoverUrl:      video.CoverUrl,
			FavoriteCount: int64(video.FavoriteCount),
			CommentCount:  int64(video.CommentCount),
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
	err := db.Where("user_id=?", otherUId).Find(&dyVideoList).Error
	if err != nil {
		return make([]vo.Video, 0), err
	}
	var videoVoList []vo.Video
	for _, video := range dyVideoList {
		var userTmp entity.DyUser
		db.Where("id=?", video.UserId).Find(&userTmp)
		var dyRelation entity.DyRelation
		err := db.Where("follower_id=? AND following_id=?", curUserId, otherUId).Find(&dyRelation).Error
		var isFollow = false
		if err != gorm.ErrRecordNotFound {
			isFollow = true
		}
		var dyFavorite entity.DyFavorite
		err2 := db.Where("user_id=? AND video_id=?", curUserId, video.Id).Find(&dyFavorite).Error
		var isFavorite = false
		if err2 != gorm.ErrRecordNotFound {
			isFavorite = true
		}
		videoVoTmp := vo.Video{
			Id: video.Id,
			Author: vo.User{
				Id:            userTmp.Id,
				Name:          userTmp.Username,
				FollowCount:   int64(userTmp.FollowCount),
				FollowerCount: int64(userTmp.FollowerCount),
				IsFollow:      isFollow,
			},
			PlayUrl:       video.PlayUrl,
			CoverUrl:      video.CoverUrl,
			FavoriteCount: int64(video.FavoriteCount),
			CommentCount:  int64(video.CommentCount),
			IsFavorite:    isFavorite,
		}
		videoVoList = append(videoVoList, videoVoTmp)
	}
	return videoVoList, nil
}
