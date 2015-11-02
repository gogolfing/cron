package sched

import "time"

type Schedule interface {
	NextTime(from time.Time) (time.Time, bool)
}
