package dao

// TODO 事务sql操作

// Like 新增点赞记录
/**
开启事务
1.新增一条点赞记录
2.在video表中点赞数+1（是否需要CAS）
提交事务
*/
func Like(userId, videoId int64) error {
	return nil
}

// UnLike 删除记录
/**
开启事务
1.删除对应点赞记录
2.在video表中点赞数-1（是否需要CAS）
提交事务
*/
func UnLike(userId, videoId int64) error {
	return nil
}
