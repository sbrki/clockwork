# clockwork

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

	sched.Schedule().EverySingle().Second().Do(something) // EverySingle is "shorthand" for Every(1)
	sched.Schedule().EverySingle().Monday().At("9:10").Do(something)
	sched.Schedule().EverySingle().Saturday().At("8:00").Do(something)

	sched.Run()
}

func something() {
	fmt.Println("foo")

}
```
