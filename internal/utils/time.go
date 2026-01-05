package utils

import "time"

var (
	TNow = func() func() time.Time {
		loc, err := time.LoadLocation("Asia/Singapore")
		if err != nil {
			return func() time.Time {
				return time.Now()
			}
		}

		return func() time.Time {
			return time.Now().In(loc)
		}
	}()
)
