package util

import "time"

const DbTimeFormat = "2006-01-02 15:04:05 -0700 MST"
const VoTimeFormat = "1-2"

// ParseDbTimeToVoTime 将db时间格式转换为vo时间格式
func ParseDbTimeToVoTime(timeStr string) (string, error) {
	parse, err := time.Parse(DbTimeFormat, timeStr)
	if err != nil {
		return "", err
	}
	return parse.Format(VoTimeFormat), nil
}
