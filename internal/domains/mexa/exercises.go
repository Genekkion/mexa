package mexadomain

import "time"

type ExerciseId = int

type Exercise struct {
	Id        ExerciseId `json:"-"`
	Code      string     `json:"code"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
}
