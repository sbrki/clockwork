# clockwork

<p align="center"><img src ="https://github.com/whiteShtef/clockwork/raw/master/assets/logo/small.png" /></p>

A simple and intuitive scheduling library in Go.

Inspired by [python's schedule](https://github.com/dbader/schedule) and [ruby's clockwork](https://github.com/adamwiggins/clockwork) libraries.


## Example use

```go
package main

import (
	"fmt"
	"github.com/whiteshtef/clockwork"
)

func main() {
	sched := clockwork.NewScheduler()

	sched.Schedule().Every(10).Seconds().Do(something)
	sched.Schedule().Every(3).Minutes().Do(something)
	sched.Schedule().Every(4).Hours().Do(something)
	sched.Schedule().Every(2).Days().At("12:32").Do(something)
	sched.Schedule().Every(12).Weeks().Do(something)

	sched.Schedule().Every().Second().Do(something) // Every() is "shorthand" for Every(1)
	sched.Schedule().Every().Monday().Do(something)
	sched.Schedule().Every().Saturday().At("8:00").Do(something)

	sched.Run()
}

func something() {
	fmt.Println("foo")

}
```
