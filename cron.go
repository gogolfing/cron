package cron

import (
	"sync"
	"time"

	"github.com/gogolfing/cron/sched"
	"github.com/gogolfing/timequeue"
)

type Job struct {
	sched.Schedule
	Data interface{}
}

type Event struct {
	*Job
	Time time.Time
}

type Cron struct {
	lock     *sync.Mutex
	jobs     map[*Job]bool
	events   chan *Event
	location *time.Location
}

func NewCron(loc *time.Location) *Cron {
	return &Cron{
		lock:   &sync.Mutex{},
		jobs:   map[*Job]bool{},
		events: make(chan *Event),
	}
}

func (c *Cron) Add(schedStr string, data interface{}) (*Job, error) {
	return nil, nil
}

func (c *Cron) AddSchedule(sched sched.Schedule, data interface{}) *Job {
	return nil
}

func (c *Cron) AddJob(job *Job) {
}

func (c *Cron) Remove(job *Job, emit bool) bool {
	return false
}

func (c *Cron) SetJobSchedule(job *Job, sched sched.Schedule) bool {
	return false
}

func (c *Cron) SetJobParseSchedule(job *Job, schedStr string) (bool, error) {
	return false, nil
}

func (c *Cron) Start() {
}

func (c *Cron) Stop() {
}

func (c *Cron) IsRunning() bool {
	return false
}

func (c *Cron) Events() <-chan *Event {
	return c.events
}

func (c *Cron) run() {
}

func newJob(s sched.Schedule, data interface{}) *Job {
	return nil
}

func createEventFromMessage(message *timequeue.Message, job *Job) *Event {
	return nil
}

func newEvent(job *Job, time time.Time) *Event {
	return nil
}
