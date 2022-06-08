package util

import "time"

const DbTimeFormat1 = "2006-01-02 15:04:05 -0700 MST"
const VoTimeFormat = "2006年1月2日 15时4分"
const DbTimeFormat2 = "2006-01-02 15:04:05"

// ParseDbTimeToVoTime 将db时间格式转换为vo时间格式
func ParseDbTimeToVoTime(timeStr string) (string, error) {
	parse, err := time.Parse(DbTimeFormat1, timeStr)
	if err != nil {
		return "", err
	}
	return parse.Format(VoTimeFormat), nil
}

// ParseRFC3339TimeToVoTime 将RFC3339时间格式转换为vo时间格式
func ParseRFC3339TimeToVoTime(timeStr string) (string, error) {
	parse, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return "", err
	}
	return parse.Format(VoTimeFormat), nil
}

// ParseTimeUnixToDbTime 将时间戳(精确到毫秒)转换为数据库时间格式
func ParseTimeUnixToDbTime(latestTime int64) string {
	time := time.UnixMilli(latestTime)
	return time.Format(DbTimeFormat2)
}

// GetTimeUnixNow 获取当前时间,毫秒
func GetTimeUnixNow() int64 {
	return time.Now().Unix()
}
