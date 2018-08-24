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
	"io"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

type (
	LogLevel = uint8

	// Logger defines the logging interface.
	Logger interface {
		Output() io.Writer
		SetOutput(w io.Writer)
		Prefix() string
		SetPrefix(p string)
		Level() LogLevel
		SetLevel(v LogLevel)
		Print(i ...interface{})
		Printf(format string, args ...interface{})
		Debug(i ...interface{})
		Debugf(format string, args ...interface{})
		Info(i ...interface{})
		Infof(format string, args ...interface{})
		Warn(i ...interface{})
		Warnf(format string, args ...interface{})
		Error(i ...interface{})
		Errorf(format string, args ...interface{})
		Fatal(i ...interface{})
		Fatalf(format string, args ...interface{})
		Panic(i ...interface{})
		Panicf(format string, args ...interface{})
	}

	// TimeUnit is an numeration used for handling
	// time units internally.
	TimeUnit int

	// Job struct handles all the data required to
	// schedule and run jobs.
	Job struct {
		stopped bool

		identifier string
		desc       string
		scheduler  *Scheduler
		unit       TimeUnit
		frequency  int
		useAt      bool
		atHour     int
		atMinute   int
		workFunc   func()

		nextScheduledRun time.Time
	}

	// Scheduler type is used to store a group of jobs (Job structs)
	Scheduler struct {
		mtx    sync.RWMutex
		stopCh chan struct{}

		identifier string
		jobs       []Job
		logger     Logger
	}

	JobInfo = map[string]string
)

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

var (
	timeUnitName = []string{"None", "Second", "Minute", "Hour", "Day", "Week", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}

	timeNow = func() time.Time {
		return time.Now()
	}

	time1970 = time.Unix(0, 0)
)

func (j *Job) Info() JobInfo {
	ret := make(JobInfo)
	ret["id"] = j.identifier
	ret["desc"] = j.desc
	ret["unit"] = timeUnitName[j.unit]
	ret["frequency"] = strconv.Itoa(j.frequency)
	ret["useAt"] = "false"
	if j.useAt {
		ret["useAt"] = "true"
		ret["atHour"] = strconv.Itoa(j.atHour)
		ret["atMinute"] = strconv.Itoa(j.atMinute)
	}
	return ret
}

func (j *Job) Stopped() bool {
	return j.stopped
}

func (j *Job) Stop() {
	if !j.stopped {
		j.nextScheduledRun = time1970
		j.stopped = true
	}

	j.scheduler.mtx.Lock()
	defer j.scheduler.mtx.Unlock()

	n := -1
	for i, job := range j.scheduler.jobs {
		if job.ID() == j.ID() {
			n = i
			break
		}
	}
	if n == -1 {
		return
	}

	if n == 0 {
		j.scheduler.jobs = j.scheduler.jobs[1:]
	} else if n == len(j.scheduler.jobs)-1 {
		j.scheduler.jobs = j.scheduler.jobs[:n]
	} else {
		newjobs := make([]Job, len(j.scheduler.jobs)-1)
		copy(newjobs[0:n], j.scheduler.jobs[0:n])
		copy(newjobs[n:], j.scheduler.jobs[n+1:])
		j.scheduler.jobs = newjobs
	}
}

func (j *Job) ID() string {
	return j.identifier
}

func (j *Job) AddDesc(d string) *Job {
	j.desc = d
	return j
}

func (j *Job) Desc() string {
	return j.desc
}

// Every is a method that fills the given Job struct with the given frequency
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

// EverySingle is deprecated predecessor to Every()
func (j *Job) EverySingle() *Job {
	return j.Every()
}

// At method fills the given Job struct atHout and atMinute fields
//with the provided information
func (j *Job) At(t string) *Job {
	j.useAt = true
	j.atHour, _ = strconv.Atoi(strings.Split(t, ":")[0])
	j.atMinute, _ = strconv.Atoi(strings.Split(t, ":")[1])
	return j
}

// Do method fills the given job struct with the function pointer
// to the job (user provided task) itself.
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
	}
	return false
}

// Generally, At() can only be used then unit is day or WEEKDAY
func (j *Job) isAsUsedIncorrectly() bool {
	if j.useAt == true &&
		(j.unit == second || j.unit == minute ||
			j.unit == hour || j.unit == week) {
		return true
	}
	return false
}

// Returns false when job unit is Day or any of the weekdays, vice versa.
// Used for scheduling when job frequency is 1, because day and WEEKDAY
// can be used with At() function which requires different scheduling approach.
func (j *Job) unitNotDayOrWEEKDAY() bool {
	if j.unit == second || j.unit == minute ||
		j.unit == hour || j.unit == week {
		return true
	}
	return false

}

// Returns false when job unit is or any of the weekdays, vice versa.
// Used for scheduling when job frequency is > 1, because we need to
// manually check for unit since we can't schedule WEEKDAYS with
// frequency > 1 .
func (j *Job) unitNotWEEKDAY() bool {
	if j.unit == second || j.unit == minute ||
		j.unit == hour || j.unit == day ||
		j.unit == week {
		return true
	}
	return false
}

func (j *Job) scheduleNextRun() {
	if j.stopped {
		return
	}

	// If Every(frequency) == 1, unit can be anything .
	// At() can be used only with day and WEEKDAY
	if j.frequency == 1 {

		// Panic if usage of "At()" is incorrect
		if j.isAsUsedIncorrectly() {
			panic(
				`Cannot schedule Every(1) with At()
				 when unit is not day or WEEKDAY`,
			) // TODO: Turn this into err
		}

		// Handle everything except day and WEEKDAY -- these guys don't use At()
		if j.unitNotDayOrWEEKDAY() {
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
				// 168 hours in a week
				j.nextScheduledRun = j.nextScheduledRun.Add(168 * time.Hour)
			}
		} else {
			// Handle Day and WEEKDAY  --  these guys use At()
			switch j.unit {
			case day:
				if j.nextScheduledRun == (time.Time{}) {
					now := timeNow()
					lastMidnight := time.Date(
						now.Year(),
						now.Month(),
						now.Day(),
						0, 0, 0, 0,
						time.Local,
					)
					if j.useAt == true {
						j.nextScheduledRun = lastMidnight.Add(
							time.Duration(j.atHour)*time.Hour +
								time.Duration(j.atMinute)*time.Minute,
						)
					} else {
						j.nextScheduledRun = lastMidnight
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

		j.scheduler.logger.Debug("Scheduled for ", j.nextScheduledRun)

	} else {
		// If Every(frequency) > 1, unit has to be either
		// second, minute, hour, day, week - not a WEEKDAY .
		// At() can be used only with day

		// Panic if usage of "At()" is incorrect
		if j.isAsUsedIncorrectly() {
			panic("Cannot schedule Every(>1) with At() when unit is not day")
			// TODO: Turn this into err
		}

		// Unlike when frequency = 1, here unit can't be anyhing.
		// We have to check that it isn't a WEEKDAY
		if j.unitNotWEEKDAY() {

			// Handle everything except day -- these guys don't use At()
			if j.unit != day {
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
				if j.nextScheduledRun == (time.Time{}) {
					now := timeNow()
					lastMidnight := time.Date(
						now.Year(),
						now.Month(),
						now.Day(),
						0, 0, 0, 0,
						time.Local,
					)
					if j.useAt == true {
						j.nextScheduledRun = lastMidnight.Add(
							time.Duration(j.atHour)*time.Hour +
								time.Duration(j.atMinute)*time.Minute,
						)
					} else {
						j.nextScheduledRun = lastMidnight
					}
				}
				j.nextScheduledRun = j.nextScheduledRun.Add(
					time.Duration(j.frequency*24) * time.Hour,
				)
			}

		} else {
			panic("Cannot schedule Every(>1) when unit is WEEKDAY")
			// TODO: Turn this into err
		}

		j.scheduler.logger.Debug("Scheduled for ", j.nextScheduledRun)
	}
	return
}

func (j *Job) scheduleWeekday(dayOfWeek time.Weekday) {
	if j.nextScheduledRun == (time.Time{}) {
		now := timeNow()
		lastWeekdayMidnight := time.Date(
			now.Year(),
			now.Month(),
			now.Day()-int(now.Weekday()-dayOfWeek),
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

// Second method fills the given job struct with seconds
func (j *Job) Second() *Job {
	j.unit = second
	return j
}

// Seconds method fills the given job struct with seconds
func (j *Job) Seconds() *Job {
	j.unit = second
	return j
}

// Minute method fills the given job struct with minutes
func (j *Job) Minute() *Job {
	j.unit = minute
	return j
}

// Minutes method fills the given job struct with minutes
func (j *Job) Minutes() *Job {
	j.unit = minute
	return j
}

// Hour method fills the given job struct with hours
func (j *Job) Hour() *Job {
	j.unit = hour
	return j
}

// Hours method fills the given job struct with hours
func (j *Job) Hours() *Job {
	j.unit = hour
	return j
}

// Day method fills the given job struct with days
func (j *Job) Day() *Job {
	j.unit = day
	return j
}

// Days method fills the given job struct with days
func (j *Job) Days() *Job {
	j.unit = day
	return j
}

// Week method fills the given job struct with weeks
func (j *Job) Week() *Job {
	j.unit = week
	return j
}

// Weeks method fills the given job struct with weeks
func (j *Job) Weeks() *Job {
	j.unit = week
	return j
}

// Monday method fills the given job struct with monday
func (j *Job) Monday() *Job {
	j.unit = monday
	return j
}

// Tuesday method fills the given job struct with tuesday
func (j *Job) Tuesday() *Job {
	j.unit = tuesday
	return j
}

// Wednesday method fills the given job struct with wednesday
func (j *Job) Wednesday() *Job {
	j.unit = wednesday
	return j
}

// Thursday method fills the given job struct with thursday
func (j *Job) Thursday() *Job {
	j.unit = thursday
	return j
}

// Friday method fills the given job struct with friday
func (j *Job) Friday() *Job {
	j.unit = friday
	return j
}

// Saturday method fills the given job struct with saturday
func (j *Job) Saturday() *Job {
	j.unit = saturday
	return j
}

// Sunday method fills the given job struct with sunday
func (j *Job) Sunday() *Job {
	j.unit = sunday
	return j
}

// NewScheduler creates and returns a new Scheduler
func NewScheduler(logger Logger) Scheduler {
	return Scheduler{
		identifier: uuid.New().String(),
		jobs:       make([]Job, 0),
		logger:     logger,
	}
}

// activateTestMode method sets the timeNow func for testing,
// by setting the current time to a fixed value
func (s *Scheduler) activateTestMode() {
	timeNow = func() time.Time {
		return time.Date(1, 1, 1, 1, 1, 0, 0, time.Local)
	}
}

// Run method on the Scheduler type runs the scheduler.
// This is a blocking method, and should be run as a goroutine.
func (s *Scheduler) Run() {
	s.stopCh = make(chan struct{}, 1)
	s.logger.Debugf("scheduler: %s starts at %s\n", s.identifier, timeNow())
	go func() {
	LOOP:
		for {
			select {
			case <-s.stopCh:
				s.logger.Debugf("scheduler: %s stops at %s\n", s.identifier, timeNow())
				break LOOP
			default:
				s.mtx.RLock()
				for jobIdx := range s.jobs {
					job := &s.jobs[jobIdx]
					if job.due() {
						job.scheduleNextRun()
						go job.workFunc()
					}
				}
				s.mtx.RUnlock()

				time.Sleep(1 * time.Second)
			}
		}
	}()
}

// Schedule method on the Scheduler creates a new Job
// and prepares is for "filling"
func (s *Scheduler) Schedule() *Job {
	newJob := Job{
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
	return &newJob
}

func (s *Scheduler) GetJob(id string) *Job {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	for _, job := range s.jobs {
		if job.ID() == id {
			return &job
		}
	}

	return nil
}

func (s *Scheduler) Size() int {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return len(s.jobs)
}

func (s *Scheduler) Stop() {
	s.stopCh <- struct{}{}
}

func (s *Scheduler) AllJobs() []Job {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	cpJobs := make([]Job, len(s.jobs))
	copy(cpJobs, s.jobs)
	return cpJobs
}
