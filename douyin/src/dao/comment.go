package dao

import (
	"douyin/src/global"
	"douyin/src/pkg/errcode"
	"douyin/src/pojo/entity"
	"douyin/src/pojo/vo"
	"time"
)

const FindCommentListByVideoIdAndUIdSQL = `SELECT
	dy_comment.id,
	dy_comment.content,
	dy_comment.create_date,
	dy_user.id,
	dy_user.username,
	dy_user.follow_count,
	dy_user.follower_count 
FROM
	dy_comment
	LEFT JOIN dy_user ON dy_comment.user_id = dy_user.id 
WHERE
	dy_comment.video_id = ?`

// FindCommentListByVideoIdAndUId 评论列表
func FindCommentListByVideoIdAndUId(videoId, userId int64) ([]vo.Comment, error) {
	followerIdsChan := make(chan map[int64]struct{})
	errorChan := make(chan error)
	var followerIdMap map[int64]struct{}

	timer := time.NewTimer(time.Second)
	// 首先查询当前用户的关注列表
	go findFollowerIdsByFollowing(followerIdsChan, errorChan, userId)

	// 从数据库查询评论列表
	db := global.DBEngine
	rows, err := db.DB().Query(FindCommentListByVideoIdAndUIdSQL, videoId)
	comments := make([]vo.Comment, 0)
	if err != nil {
		return comments, err
	}
	defer rows.Close()

	// 等待上方查询关注列表的协程返回结果，超时则直接返回错误码
loop:
	for {
		select {
		case err = <-errorChan:
			return comments, err
		case followerIdMap = <-followerIdsChan:
			break loop
		default:
			select {
			case <-timer.C:
				return comments, errcode.TimeOutFail
			default:
				continue
			}
		}
	}

	// 拼装数据为vo并加入列表
	for rows.Next() {
		var comment entity.DyComment
		var commentator entity.DyUser
		if err = rows.Scan(&comment.Id, &comment.Content, &comment.CreateDate, &commentator.Id,
			&commentator.Username, &commentator.FollowCount, &commentator.FollowerCount); err != nil {
			return comments, err
		}
		voUser := commentator.NewVoUser()
		if voUser == nil {
			continue
		}
		voComment := comment.NewVoComment()
		if voComment == nil {
			continue
		}

		// 如当前用户的关注列表中有此评论者，则评论者的IsFollow为true
		_, voUser.IsFollow = followerIdMap[voUser.Id]
		voComment.User = *voUser
		comments = append(comments, *voComment)
	}
	return comments, nil
}

// FindCommentListByVideoWithoutLogin 未登录查看评论列表
func FindCommentListByVideoWithoutLogin(videoId int64) ([]vo.Comment, error) {

	// 从数据库查询评论列表
	db := global.DBEngine
	rows, err := db.DB().Query(FindCommentListByVideoIdAndUIdSQL, videoId)
	comments := make([]vo.Comment, 0)
	if err != nil {
		return comments, err
	}
	defer rows.Close()

	// 拼装数据为vo并加入列表
	for rows.Next() {
		var comment entity.DyComment
		var commentator entity.DyUser
		if err = rows.Scan(&comment.Id, &comment.Content, &comment.CreateDate, &commentator.Id,
			&commentator.Username, &commentator.FollowCount, &commentator.FollowerCount); err != nil {
			return comments, err
		}

		voUser := commentator.NewVoUser()
		if voUser == nil {
			continue
		}

		voComment := comment.NewVoComment()
		if voComment == nil {
			continue
		}

		voComment.User = *voUser
		comments = append(comments, *voComment)
	}
	return comments, nil
}
