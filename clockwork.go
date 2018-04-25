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
	j._schedule_next_run()
	j.scheduler.jobs = append(j.scheduler.jobs, *j)
	return j.identifier
}

func (j *Job) _schedule_next_run() {
	/*	examples from python/schedule:
		schedule.every(10).minutes.do(job)
		schedule.every().hour.do(job)
		schedule.every().day.at("10:30").do(job)
		schedule.every().monday.do(job)
		schedule.every().wednesday.at("13:15").do(job)
	*/
	if j.frequency == 1 {
		// If Every(frequency) == 1, unit can be anything .
		// at() can be used only with day and WEEKDAY
		if j.use_at == true && (j.unit == minute || j.unit == hour || j.unit == week) {
			panic("Cannot schedule Every(1) with at() when unit is not day or WEEKDAY") // TODO: Turn this into err
		} else {
			fmt.Println("OK")
		}
	} else {
		// If Every(frequency) > 1, unit has to be either second, minute, hour, day, week - not a WEEKDAY
		// at() can be used only with day

		if j.unit == second || j.unit == minute || j.unit == hour || j.unit == day || j.unit == week {
			if j.use_at == true && (j.unit == day) {
				fmt.Println("OK")
			} else {
				panic("Cannot schedule Every(>1) with at() when unit is not day")
			}
		} else {
			panic("Cannot schedule Every(>1) when unit is WEEKDAY") // TODO: Turn this into err
		}
	}
	return
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

func (j *Job) Minute() *Job {
	j.unit = minute
	return j
}

func (j *Job) Minutes() *Job {
	j.unit = minute
	return j
}

func (j *Job) Hours() *Job {
	j.unit = hour
	return j
}

func (j *Job) Hour() *Job {
	j.unit = hour
	return j
}

func (j *Job) Days() *Job {
	j.unit = day
	return j
}

func (j *Job) Day() *Job {
	j.unit = day
	return j
}

func (j *Job) Weeks() *Job {
	j.unit = week
	return j
}

func (j *Job) Week() *Job {
	j.unit = week
	return j
}

func (j *Job) Monday() *Job {
	j.unit = monday
	return j
}

func (j *Job) Tuesday() *Job {
	j.unit = tuesday
	return j
}

func (j *Job) Wednesday() *Job {
	j.unit = wednesday
	return j
}

func (j *Job) Thursday() *Job {
	j.unit = thursday
	return j
}

func (j *Job) Friday() *Job {
	j.unit = friday
	return j
}

func (j *Job) Saturday() *Job {
	j.unit = saturday
	return j
}

func (j *Job) Sunday() *Job {
	j.unit = sunday
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
		frequency:  1,
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
	sched.Schedule().Every(2).Days().at("12:23").Do(printaj)
	sched.Run()
}

func printaj() {
	//fmt.Println("bok")

}
