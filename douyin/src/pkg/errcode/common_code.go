package errcode

//定义错误码
var (
	Success                   = NewError(0, "成功")
	UnKnowError               = NewError(1, "失败")
	ServerError               = NewError(10000000, "服务内部错误")
	InvalidParams             = NewError(10000001, "入参错误")
	NotFound                  = NewError(10000002, "找不到")
	UnauthorizedAuthNotExist  = NewError(10000003, "鉴权失败")
	UnauthorizedTokenError    = NewError(10000004, "鉴权失败，Token错误")
	UnauthorizedTokenTimeout  = NewError(10000005, "鉴权失败，Token超时")
	UnauthorizedTokenGenerate = NewError(10000006, "鉴权失败，Token生产失败")
	TooManyRequest            = NewError(10000007, "请求过多")
	LikeFail                  = NewError(10000008, "点赞失败")
)
