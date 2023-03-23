package utils

import "time"

const (
	TIME_YMDHMT = "20060102150405"
)

func GetNowTimeString() string {
	return time.Now().Format(TIME_YMDHMT)
}
