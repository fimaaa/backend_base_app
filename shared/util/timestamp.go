package util

import "time"

func MakeTimestamp() int64 {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)
	return now.UnixNano() / int64(time.Millisecond)
}
