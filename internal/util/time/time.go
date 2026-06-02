package time

import "time"

func NowUtc0() int64 {
	return time.Now().UnixMilli()
}
