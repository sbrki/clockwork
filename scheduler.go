package main

import (
	"fmt"
	"github.com/google/uuid"
	"strconv"
	"strings"
	"time"
)

type TimeUnit int

const (
	None = iota
	second
	//seconds
	minute
	//minutes
	hour
	//hours
	day
	//days
	week
	//weeks
	monday
	tuesday
	wednesday
	thursday
	friday
	saturday
	sunday
)

/*****************************************************************************/

type Job struct {
	identifier string
	scheduler  *Scheduler
	unit       TimeUnit
	frequency  int
	use_at     bool
	at_hour    int
	at_minute  int
	jobfunc    func()

	last_run time.Time
	next_run time.Time
}

func (j *Job) Every(frequency int) *Job {
	j.frequency = frequency
	return j
}

func (j *Job) at(t string) *Job {
	j.use_at = true
	j.at_hour, _ = strconv.Atoi(strings.Split(t, ":")[0])
	j.at_minute, _ = strconv.Atoi(strings.Split(t, ":")[1])
	return j
}

func (j *Job) Do(function func()) string {
	j.jobfunc = function
	j.scheduler.jobs = append(j.scheduler.jobs, *j)
	return j.identifier
}

/**********************/
func (j *Job) Second() *Job {
	j.unit = second
	return j
}
func (j *Job) Seconds() *Job {
	j.unit = second
	return j
}

/*****************************************************************************/

type Scheduler struct {
	identifier string
	jobs       []Job
}

func New() Scheduler {
	return Scheduler{
		identifier: uuid.New().String(),
		jobs:       make([]Job, 0),
	}
}

func (s *Scheduler) Run() {
	for {
		for _, job := range s.jobs {
			if job.jobfunc == nil {
				fmt.Println("nil")
			} else {
				job.jobfunc()
			}
		}
	}
}

func (s *Scheduler) Schedule() *Job {
	new_job := Job{
		identifier: uuid.New().String(),
		scheduler:  s,
		unit:       None,
		frequency:  0,
		use_at:     false,
		at_hour:    0,
		at_minute:  0,
		jobfunc:    nil,
		last_run:   time.Time{}, // zero value
		next_run:   time.Time{}, // zero value
	}
	return &new_job
}

/*****************************************************************************/

func main() {
	sched := New()
	sched.Schedule().Every(10).Seconds().Do(printaj)
	sched.Run()
}

func printaj() {
	fmt.Println("bok")
}
