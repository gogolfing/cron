package cron

import "sync"

type Job struct {
	Data     interface{}
	schedule *schedule
}

type Cron struct {
	lock *sync.Mutex
	jobs map[*Job]bool
}

func NewCron() *Cron {
	return &Cron{
		lock: &sync.Mutex{},
	}
}
