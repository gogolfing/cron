package cron

import "time"

type schedule struct {
}

func (s *schedule) nextTime(from time.Time) (time.Time, bool) {
	return from, true
}
