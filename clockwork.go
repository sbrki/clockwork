package clockwork

import (
	//	"fmt"
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

	next_scheduled_run time.Time
}

func (j *Job) Every(frequency int) *Job {
	if frequency <= 0 {
		panic("Every(frequency) has to be >= 1")
	}
	j.frequency = frequency
	return j
}

func (j *Job) EverySingle() *Job {
	j.frequency = 1
	return j
}

func (j *Job) At(t string) *Job {
	j.use_at = true
	j.at_hour, _ = strconv.Atoi(strings.Split(t, ":")[0])
	j.at_minute, _ = strconv.Atoi(strings.Split(t, ":")[1])
	return j
}

func (j *Job) Do(function func()) string {
	j.jobfunc = function
	j.schedule_next_run()
	j.scheduler.jobs = append(j.scheduler.jobs, *j)
	return j.identifier
}

func (j *Job) due() bool {
	now := time.Now()
	if now.After(j.next_scheduled_run) {
		return true
	} else {
		return false
	}
}

func (j *Job) schedule_next_run() {
	/*	examples from python/schedule:
		schedule.every(10).minutes.do(job)
		schedule.every().hour.do(job)
		schedule.every().day.at("10:30").do(job)
		schedule.every().monday.do(job)
		schedule.every().wednesday.at("13:15").do(job)
	*/

	// If Every(frequency) == 1, unit can be anything .
	// At() can be used only with day and WEEKDAY
	if j.frequency == 1 {
		// Panic if usage of "At()" is incorrect
		if j.use_at == true && (j.unit == minute || j.unit == hour || j.unit == week) {
			panic("Cannot schedule Every(1) with At() when unit is not day or WEEKDAY") // TODO: Turn this into err
		}

		// Handle everything except day and WEEKDAY -- these guys don't use At()
		if j.unit == second || j.unit == minute || j.unit == hour || j.unit == week {
			if j.next_scheduled_run == (time.Time{}) {
				j.next_scheduled_run = time.Now()
			}

			switch j.unit {
			case second:
				j.next_scheduled_run = j.next_scheduled_run.Add(1 * time.Second)
			case minute:
				j.next_scheduled_run = j.next_scheduled_run.Add(1 * time.Minute)
			case hour:
				j.next_scheduled_run = j.next_scheduled_run.Add(1 * time.Hour)
			case week:
				j.next_scheduled_run = j.next_scheduled_run.Add(168 * time.Hour) // 168 hours in a week

			}
		} else {
			// Handle Day and WEEKDAY  --  these guys use At()
			switch j.unit {
			case day:
				if j.next_scheduled_run == (time.Time{}) {
					now := time.Now()
					last_midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
					if j.use_at == true {
						j.next_scheduled_run = last_midnight.Add(
							time.Duration(j.at_hour)*time.Hour +
								time.Duration(j.at_minute)*time.Minute,
						)
					} else {
						j.next_scheduled_run = last_midnight
					}
				}
				j.next_scheduled_run = j.next_scheduled_run.Add(24 * time.Hour)

			case monday:
				if j.next_scheduled_run == (time.Time{}) {
					now := time.Now()
					last_monday_midnight := time.Date(
						now.Year(),
						now.Month(),
						now.Day()-int(now.Weekday()-time.Monday),
						0, 0, 0, 0,
						time.Local)
					if j.use_at == true {
						j.next_scheduled_run = last_monday_midnight.Add(
							time.Duration(j.at_hour)*time.Hour +
								time.Duration(j.at_minute)*time.Minute,
						)
					} else {
						j.next_scheduled_run = last_monday_midnight
					}
				}
				j.next_scheduled_run = j.next_scheduled_run.Add(7 * 24 * time.Hour)

			case tuesday:
				if j.next_scheduled_run == (time.Time{}) {
					now := time.Now()
					last_tuesday_midnight := time.Date(
						now.Year(),
						now.Month(),
						now.Day()-int(now.Weekday()-time.Tuesday),
						0, 0, 0, 0,
						time.Local)
					if j.use_at == true {
						j.next_scheduled_run = last_tuesday_midnight.Add(
							time.Duration(j.at_hour)*time.Hour +
								time.Duration(j.at_minute)*time.Minute,
						)
					} else {
						j.next_scheduled_run = last_tuesday_midnight
					}
				}
				j.next_scheduled_run = j.next_scheduled_run.Add(7 * 24 * time.Hour)

			case wednesday:
				if j.next_scheduled_run == (time.Time{}) {
					now := time.Now()
					last_wednesday_midnight := time.Date(
						now.Year(),
						now.Month(),
						now.Day()-int(now.Weekday()-time.Wednesday),
						0, 0, 0, 0,
						time.Local)
					if j.use_at == true {
						j.next_scheduled_run = last_wednesday_midnight.Add(
							time.Duration(j.at_hour)*time.Hour +
								time.Duration(j.at_minute)*time.Minute,
						)
					} else {
						j.next_scheduled_run = last_wednesday_midnight
					}
				}
				j.next_scheduled_run = j.next_scheduled_run.Add(7 * 24 * time.Hour)

			case thursday:
				if j.next_scheduled_run == (time.Time{}) {
					now := time.Now()
					last_thursday_midnight := time.Date(
						now.Year(),
						now.Month(),
						now.Day()-int(now.Weekday()-time.Thursday),
						0, 0, 0, 0,
						time.Local)
					if j.use_at == true {
						j.next_scheduled_run = last_thursday_midnight.Add(
							time.Duration(j.at_hour)*time.Hour +
								time.Duration(j.at_minute)*time.Minute,
						)
					} else {
						j.next_scheduled_run = last_thursday_midnight
					}
				}
				j.next_scheduled_run = j.next_scheduled_run.Add(7 * 24 * time.Hour)

			case saturday:
				if j.next_scheduled_run == (time.Time{}) {
					now := time.Now()
					last_saturday_midnight := time.Date(
						now.Year(),
						now.Month(),
						now.Day()-int(now.Weekday()-time.Saturday),
						0, 0, 0, 0,
						time.Local)
					if j.use_at == true {
						j.next_scheduled_run = last_saturday_midnight.Add(
							time.Duration(j.at_hour)*time.Hour +
								time.Duration(j.at_minute)*time.Minute,
						)
					} else {
						j.next_scheduled_run = last_saturday_midnight
					}
				}
				j.next_scheduled_run = j.next_scheduled_run.Add(7 * 24 * time.Hour)

			case sunday:
				if j.next_scheduled_run == (time.Time{}) {
					now := time.Now()
					last_sunday_midnight := time.Date(
						now.Year(),
						now.Month(),
						now.Day()-int(now.Weekday()-time.Sunday),
						0, 0, 0, 0,
						time.Local)
					if j.use_at == true {
						j.next_scheduled_run = last_sunday_midnight.Add(
							time.Duration(j.at_hour)*time.Hour +
								time.Duration(j.at_minute)*time.Minute,
						)
					} else {
						j.next_scheduled_run = last_sunday_midnight
					}
				}
				j.next_scheduled_run = j.next_scheduled_run.Add(7 * 24 * time.Hour)

			}

		}

		//fmt.Println("Scheduled for ", j.next_scheduled_run)

	} else {
		// If Every(frequency) > 1, unit has to be either second, minute, hour, day, week - not a WEEKDAY
		// At() can be used only with day

		if j.unit == second || j.unit == minute || j.unit == hour || j.unit == day || j.unit == week {
			if j.use_at == true && (j.unit != day) {
				panic("Cannot schedule Every(>1) with At() when unit is not day") // TODO: Turn this into err
			}
			// Handle everything except  day  -- these guys don't use At()
			if j.unit == second || j.unit == minute || j.unit == hour || j.unit == week {
				if j.next_scheduled_run == (time.Time{}) {
					j.next_scheduled_run = time.Now()
				}

				switch j.unit {
				case second:
					j.next_scheduled_run = j.next_scheduled_run.Add(time.Duration(j.frequency) * time.Second)
				case minute:
					j.next_scheduled_run = j.next_scheduled_run.Add(time.Duration(j.frequency) * time.Minute)
				case hour:
					j.next_scheduled_run = j.next_scheduled_run.Add(time.Duration(j.frequency) * time.Hour)
				case week:
					j.next_scheduled_run = j.next_scheduled_run.Add(time.Duration(j.frequency*168) * time.Hour) // 168 hours in a week

				}
			} else {
				// Handle Day  --  these guy uses At()
				switch j.unit { // switch is here not really neccesarry since day is
				case day: // the only option left.
					if j.next_scheduled_run == (time.Time{}) {
						now := time.Now()
						last_midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
						if j.use_at == true {
							j.next_scheduled_run = last_midnight.Add(
								time.Duration(j.at_hour)*time.Hour +
									time.Duration(j.at_minute)*time.Minute,
							)
						} else {
							j.next_scheduled_run = last_midnight
						}
					}
					j.next_scheduled_run = j.next_scheduled_run.Add(time.Duration(j.frequency*24) * time.Hour)
				}
			}

		} else {
			panic("Cannot schedule Every(>1) when unit is WEEKDAY") // TODO: Turn this into err
		}
		//fmt.Println("Scheduled for ", j.next_scheduled_run)

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

/*****************************************************************************/

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

func (s *Scheduler) Run() {
	for {
		for jobIdx := range s.jobs {
			job := &s.jobs[jobIdx]
			if job.due() {
				job.schedule_next_run()
				go job.jobfunc()
			}
		}

	}
}

func (s *Scheduler) Schedule() *Job {
	new_job := Job{
		identifier:         uuid.New().String(),
		scheduler:          s,
		unit:               None,
		frequency:          1,
		use_at:             false,
		at_hour:            0,
		at_minute:          0,
		jobfunc:            nil,
		next_scheduled_run: time.Time{}, // zero value
	}
	return &new_job
}

/*****************************************************************************/
