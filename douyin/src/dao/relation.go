package dao

import (
	"douyin/src/global"
	"github.com/jinzhu/gorm"
)

func FindFollowerIdsByFollowing(following int64) ([]int64, error) {
	db := global.DBEngine
	var followerIds []int64
	if err := db.Table("dy_relation").Where("follower_id = ?", following).Find(&followerIds).Error; err == nil || err == gorm.ErrRecordNotFound {
		return followerIds, nil
	} else {
		return followerIds, err
	}
}
