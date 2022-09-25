# clockwork

<p align="center"><img height=250 src ="https://github.com/sbrki/clockwork/raw/master/assets/logo/large.png" /></p>

[![Awesome](https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg)](https://github.com/sindresorhus/awesome) 
[![GoDoc](https://godoc.org/github.com/sbrki/clockwork?status.svg)](https://godoc.org/github.com/sbrki/clockwork)
[![Go Report Card](https://goreportcard.com/badge/github.com/sbrki/clockwork)](https://goreportcard.com/report/github.com/sbrki/clockwork)
![Coverage](http://gocover.io/_badge/github.com/sbrki/clockwork)


A simple and intuitive scheduling library in Go.

Inspired by [python's schedule](https://github.com/dbader/schedule) and [ruby's clockwork](https://github.com/adamwiggins/clockwork) libraries.


## Example use

```go
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

The package used to use [go dep](https://golang.github.io/dep/) for dependency management.
It has switched to go modules as of commit `5f1b50934f209adb9930ef98fe654f814156a858`, which
became available under `v1.0.0`
