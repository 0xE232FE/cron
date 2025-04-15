Fork of https://godoc.org/github.com/robfig/cron  

Changes:  
Jobs has a context  
Jobs recover panic by default  
Jobs can have a description (label)  
Can have a job that only run once using cron.Once  
Can change the *time.Location while running  
Do not sort all entries at evey iteration  
Tests runs in under a second (instead of over a minute)  

```go
package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/alaingilbert/cron"
	"log"
	"os"
	"time"
)

type SomeJob struct{}

// Run implements the Job interface for SomeJob
func (s SomeJob) Run(context.Context) error {
	fmt.Println("Some job")
	return nil
}

func main() {
	appCtx := context.Background()

	l := log.New(os.Stdout, "", log.LstdFlags)

	c := cron.New(
		cron.WithContext(appCtx), // Can provide a custom context
		cron.WithLogger(l),       // Can provide a custom logger
	)

	// A simple function can be used as a job
	_, _ = c.AddJob("* * * * * *", cron.WithTimeout(2*time.Second, func() {
		fmt.Println("hello world")
	}))

	// If your function returns an error, it will be logged by the cron logger
	_, _ = c.AddJob("* * * * * *", func() error {
		return errors.New("this error will be logged by the cron logger")
	})

	// You can give a description to your jobs using the cron.Label helper
	// The label will be part of the log if the job returns an error
	_, _ = c.AddJob("* * * * * *", func() error {
		return errors.New("this error will be logged by the cron logger")
	}, cron.Label("some description"))

	// A job can also receive a context if you want to handle graceful shutdown
	_, _ = c.AddJob("* * * * * *", func(ctx context.Context) {
		// ...
	})

	// This is also valid
	_, _ = c.AddJob("* * * * * *", func(ctx context.Context) error {
		return nil
	})

	// WithTimeout will ensure that the context get cancelled after the given duration
	_, _ = c.AddJob("*/3 * * * * *", cron.WithTimeout(2*time.Second, func(ctx context.Context) {
		select {
		case <-time.After(time.Minute):
		case <-ctx.Done():
			return
		}
	}))

	// This is also valid
	_, _ = c.AddJob("*/3 * * * * *", cron.WithTimeout(2*time.Second, func(ctx context.Context) error {
		select {
		case <-time.After(time.Minute):
		case <-ctx.Done():
			return ctx.Err()
		}
		return nil
	}))

	// Anything that implements the Job interface can be used for job
	_, _ = c.AddJob("* * * * * *", SomeJob{})

	// When using cron.Once, the job will remove itself from the cron entries
	// after being executed once
	_, _ = c.AddJob("* * * * * *", cron.Once(func() {
		fmt.Println("Will only be executed once")
	}))

	// cron.SkipIfStillRunning will ensure that the job is skipped if a previous
	// invocation is still running
	_, _ = c.AddJob("* * * * * *", cron.SkipIfStillRunning(func() {
		fmt.Println("Slow job")
		time.Sleep(3 * time.Second)
	}))

	// You can use cron.Chain to chain wrappers to your job
	_, _ = c.AddJob("* * * * * *", cron.Chain(func() {
		fmt.Println("job with chained wrappers")
	}, cron.TimeoutWrapper(time.Second), cron.Once))

	// Jitter adds some random delay before running the job
	_, _ = c.AddJob("0 */5 * * * *", cron.WithJitter(time.Minute, func() {
		fmt.Println("Runs every 5min, but wait 0-1 minute before starting")
	}))
	
	c.Run()
}
```

```go
go test -v ./...
go clean -testcache && go test -race ./...
```