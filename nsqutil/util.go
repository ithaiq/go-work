package nsqutil

import (
	"time"
)

func FormatUnixTimestamp() string {
	timeLayout := "2006-01-02 15:04:05"
	t := time.Unix(time.Now().Unix(), 0)
	return t.Format(timeLayout)
}
