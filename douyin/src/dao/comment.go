package dao

import "douyin/src/pojo/vo"

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
	// TODO 完成 dao
	return nil, nil
}
