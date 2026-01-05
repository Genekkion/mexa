package mexadomain

import (
	"fmt"
	"mexa/internal/utils/set"
	"time"
)

type ExLogId = int

type ExLog struct {
	Id         ExLogId    `json:"-"`
	ExerciseId ExerciseId `json:"-"`
	UserId     UserId     `json:"-"`
	CreatedAt  time.Time  `json:"-"`
	Type       ExLogType  `json:"-"`
}

type ExLogType string

const (
	LogTypeExStart ExLogType = "ex_start"
	LogTypeExEnd   ExLogType = "ex_end"
)

var (
	validTypes = set.New(set.WithSlice([]ExLogType{
		LogTypeExStart,
		LogTypeExEnd,
	}))
)

func ParseExLogType(s string) (t *ExLogType, err error) {
	v := ExLogType(s)
	if validTypes.Contains(v) {
		return &v, nil
	}
	return nil, fmt.Errorf("invalid ExLogType, s: %s", s)
}

type LogExStartEndBase struct {
	Id         ExLogId
	ExerciseId ExerciseId
	UserId     UserId
	CreatedAt  time.Time
	Type       ExLogType
}

type LogExStart struct {
	LogExStartEndBase
}

func NewLogExStart() LogExStart {
	return LogExStart{}
}

type LogExEnd struct {
	LogExStartEndBase
}

func NewLogExEnd() LogExEnd {
	return LogExEnd{}
}
