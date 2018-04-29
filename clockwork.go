// Package clockwork enables simple and intuitive scheduling in Go.
//
// Examples:
//		sched.Schedule().Every(10).Seconds().Do(something)
//		sched.Schedule().Every(3).Minutes().Do(something)
//		sched.Schedule().Every(4).Hours().Do(something)
//		sched.Schedule().Every(2).Days().At("12:32").Do(something)
//		sched.Schedule().Every(12).Weeks().Do(something)
//		sched.Schedule().Every(1).Monday().Do(something)
//		sched.Schedule().Every(1).Saturday().At("8:00").Do(something)
package clockwork

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type TimeUnit int

const (
	none = iota
	second
	minute
	hour
	day
	week
	monday
	tuesday
	wednesday
	thursday
	friday
	saturday
	sunday
)

var timeNow = func() time.Time {
	return time.Now()
}

type Job struct {
	identifier string
	scheduler  *Scheduler
	unit       TimeUnit
	frequency  int
	useAt      bool
	atHour     int
	atMinute   int
	workFunc   func()

	nextScheduledRun time.Time
}

func (j *Job) Every(frequencies ...int) *Job {
	l := len(frequencies)

	switch l {
	case 0:
		j.frequency = 1
	case 1:
		if frequencies[0] <= 0 {
			panic("Every expects frequency to be greater than of equal to 1")
		}
		j.frequency = frequencies[0]
	default:
		panic("Every expects 0 or 1 arguments")
	}

	return j
}

// Deprecating sooner than later
func (j *Job) EverySingle() *Job {
	return j.Every()
}

func (j *Job) At(t string) *Job {
	j.useAt = true
	j.atHour, _ = strconv.Atoi(strings.Split(t, ":")[0])
	j.atMinute, _ = strconv.Atoi(strings.Split(t, ":")[1])
	return j
}

func (j *Job) Do(function func()) string {
	j.workFunc = function
	j.scheduleNextRun()
	j.scheduler.jobs = append(j.scheduler.jobs, *j)
	return j.identifier
}

func (j *Job) due() bool {
	now := timeNow()
	if now.After(j.nextScheduledRun) {
		return true
	} else {
		return false
	}
}

func (j *Job) scheduleNextRun() {
	// If Every(frequency) == 1, unit can be anything .
	// At() can be used only with day and WEEKDAY
	if j.frequency == 1 {
		// Panic if usage of "At()" is incorrect
		if j.useAt == true && (j.unit == minute || j.unit == hour || j.unit == week) {
			panic("Cannot schedule Every(1) with At() when unit is not day or WEEKDAY") // TODO: Turn this into err
		}

		// Handle everything except day and WEEKDAY -- these guys don't use At()
		if j.unit == second || j.unit == minute || j.unit == hour || j.unit == week {
			if j.nextScheduledRun == (time.Time{}) {
				j.nextScheduledRun = timeNow()
			}

			switch j.unit {
			case second:
				j.nextScheduledRun = j.nextScheduledRun.Add(1 * time.Second)
			case minute:
				j.nextScheduledRun = j.nextScheduledRun.Add(1 * time.Minute)
			case hour:
				j.nextScheduledRun = j.nextScheduledRun.Add(1 * time.Hour)
			case week:
				j.nextScheduledRun = j.nextScheduledRun.Add(168 * time.Hour) // 168 hours in a week

			}
		} else {
			// Handle Day and WEEKDAY  --  these guys use At()
			switch j.unit {
			case day:
				if j.nextScheduledRun == (time.Time{}) {
					now := timeNow()
					last_midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
					if j.useAt == true {
						j.nextScheduledRun = last_midnight.Add(
							time.Duration(j.atHour)*time.Hour +
								time.Duration(j.atMinute)*time.Minute,
						)
					} else {
						j.nextScheduledRun = last_midnight
					}
				}
				j.nextScheduledRun = j.nextScheduledRun.Add(24 * time.Hour)

			case monday:
				j.scheduleWeekday(time.Monday)
			case tuesday:
				j.scheduleWeekday(time.Tuesday)
			case wednesday:
				j.scheduleWeekday(time.Wednesday)
			case thursday:
				j.scheduleWeekday(time.Thursday)
			case friday:
				j.scheduleWeekday(time.Friday)
			case saturday:
				j.scheduleWeekday(time.Saturday)
			case sunday:
				j.scheduleWeekday(time.Sunday)
			}

		}

		fmt.Println("Scheduled for ", j.nextScheduledRun)

	} else {
		// If Every(frequency) > 1, unit has to be either second, minute, hour, day, week - not a WEEKDAY
		// At() can be used only with day

		if j.unit == second || j.unit == minute || j.unit == hour || j.unit == day || j.unit == week {
			if j.useAt == true && (j.unit != day) {
				panic("Cannot schedule Every(>1) with At() when unit is not day") // TODO: Turn this into err
			}
			// Handle everything except  day  -- these guys don't use At()
			if j.unit == second || j.unit == minute || j.unit == hour || j.unit == week {
				if j.nextScheduledRun == (time.Time{}) {
					j.nextScheduledRun = timeNow()
				}

				switch j.unit {
				case second:
					j.nextScheduledRun = j.nextScheduledRun.Add(
						time.Duration(j.frequency) * time.Second,
					)
				case minute:
					j.nextScheduledRun = j.nextScheduledRun.Add(
						time.Duration(j.frequency) * time.Minute,
					)
				case hour:
					j.nextScheduledRun = j.nextScheduledRun.Add(
						time.Duration(j.frequency) * time.Hour,
					)
				case week:
					j.nextScheduledRun = j.nextScheduledRun.Add(
						time.Duration(j.frequency*168) * time.Hour,
					) // 168 hours in a week

				}
			} else {
				// Handle Day  --  these guy uses At()
				switch j.unit { // switch is here not really neccesarry since day is
				case day: // the only option left.
					if j.nextScheduledRun == (time.Time{}) {
						now := timeNow()
						last_midnight := time.Date(
							now.Year(),
							now.Month(),
							now.Day(),
							0, 0, 0, 0,
							time.Local,
						)
						if j.useAt == true {
							j.nextScheduledRun = last_midnight.Add(
								time.Duration(j.atHour)*time.Hour +
									time.Duration(j.atMinute)*time.Minute,
							)
						} else {
							j.nextScheduledRun = last_midnight
						}
					}
					j.nextScheduledRun = j.nextScheduledRun.Add(time.Duration(j.frequency*24) * time.Hour)
				}
			}

		} else {
			panic("Cannot schedule Every(>1) when unit is WEEKDAY") // TODO: Turn this into err
		}
		fmt.Println("Scheduled for ", j.nextScheduledRun)

	}
	return
}

func (j *Job) scheduleWeekday(day_of_week time.Weekday) {
	if j.nextScheduledRun == (time.Time{}) {
		now := timeNow()
		lastWeekdayMidnight := time.Date(
			now.Year(),
			now.Month(),
			now.Day()-int(now.Weekday()-day_of_week),
			0, 0, 0, 0,
			time.Local)
		if j.useAt == true {
			j.nextScheduledRun = lastWeekdayMidnight.Add(
				time.Duration(j.atHour)*time.Hour +
					time.Duration(j.atMinute)*time.Minute,
			)
		} else {
			j.nextScheduledRun = lastWeekdayMidnight
		}
	}
	j.nextScheduledRun = j.nextScheduledRun.Add(7 * 24 * time.Hour)
}

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

func (j *Job) Hour() *Job {
	j.unit = hour
	return j
}

func (j *Job) Hours() *Job {
	j.unit = hour
	return j
}

func (j *Job) Day() *Job {
	j.unit = day
	return j
}

func (j *Job) Days() *Job {
	j.unit = day
	return j
}

func (j *Job) Week() *Job {
	j.unit = week
	return j
}

func (j *Job) Weeks() *Job {
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

type Scheduler struct {
	identifier string
	jobs       []Job
}

func NewScheduler() Scheduler {
	return Scheduler{
		identifier: uuid.New().String(),
		jobs:       make([]Job, 0),
	}
}

func (s *Scheduler) isTest() {
	timeNow = func() time.Time {
		return time.Date(1, 1, 1, 1, 1, 0, 0, time.Local)
	}
}

func (s *Scheduler) Run() {
	for {
		for jobIdx := range s.jobs {
			job := &s.jobs[jobIdx]
			if job.due() {
				job.scheduleNextRun()
				go job.workFunc()
			}
		}

	}
}

func (s *Scheduler) Schedule() *Job {
	new_job := Job{
		identifier:       uuid.New().String(),
		scheduler:        s,
		unit:             none,
		frequency:        1,
		useAt:            false,
		atHour:           0,
		atMinute:         0,
		workFunc:         nil,
		nextScheduledRun: time.Time{}, // zero value
	}
	return &new_job
}
