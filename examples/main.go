package main

import (
	"fmt"

	"github.com/sbrki/clockwork"
)

func main() {
	sched := clockwork.NewScheduler()

	sched.Schedule().Every(10).Seconds().Do(something)
	sched.Schedule().Every(3).Minutes().Do(something)
	sched.Schedule().Every(4).Hours().Do(something)
	sched.Schedule().Every().Day().At("12:00")
	sched.Schedule().Every(2).Days().At("12:32").Do(something)
	sched.Schedule().Every(12).Weeks().Do(something)

	sched.Schedule().Every().Second().Do(something) // Every() is "shorthand" for Every(1)
	sched.Schedule().Every().Monday().Do(something)
	sched.Schedule().Every().Saturday().At("8:00").Do(something)

	// Polling interval defaults to 333 milliseconds, but you can set it manually
	sched.SetPollingInterval(500) // set polling interval to 500 milliseconds

	sched.Run()
}

func something() {
	fmt.Println("foo")

}
