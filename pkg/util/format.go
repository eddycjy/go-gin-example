package util

import (
	"errors"
	"strings"
	"time"
)

func FormatDatetimeToTimestamp(datetime string) (int64, error) {
	if strings.Contains(datetime, "T") && strings.Contains(datetime, "+") {
		arr1 := strings.Split(datetime, "+")
		arr2 := strings.Split(arr1[0], "T")
		timeLayout := "2006-01-02 15:04:05"
		times, err := time.Parse(timeLayout, arr2[0]+" "+arr2[1])
		if err != nil {
			return 0, err
		}
		return times.Unix(), nil
	} else {
		return 0, errors.New("datetime格式错误,无法转化")
	}
}
