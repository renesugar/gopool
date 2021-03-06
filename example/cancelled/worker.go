package main

import (
	"github.com/tomwright/gopool"
	"context"
	"fmt"
	"time"
)

func main() {
	var nameChan = make(chan string, 10)

	// define the work we want to complete
	var work = func(ctx context.Context) error {
		for {
			select {
			case <-ctx.Done():
				// if we get a message here, the context has been cancelled
				// and we should return the context err
				return ctx.Err()
			case name, ok := <-nameChan:
				// if we get a message here, either we have a name to do some work with
				// or our job channel has been closed
				if ! ok {
					return nil
				}
				fmt.Printf("Hello %s\n", name)
			}
		}
	}

	// create a worker, giving it a unique id, the piece of work, and some context
	w := gopool.NewWorker("worker-2", work, context.TODO())

	// start the worker
	cancel := w.Start()

	// the worker is now running in the background...
	// we can pass jobs to it by writing to the nameChan created above
	nameChan <- "Tom"
	nameChan <- "Jess"
	nameChan <- "Frank"
	nameChan <- "Joe"

	// let's sleep for a few microseconds to allow some processing
	time.Sleep(time.Microsecond * 50)
	// then assume some "problem" has occurred
	// and we need to stop the above from processing names immediately
	cancel()

	// we can now wait until the worker has finished (due to the cancel)
	<-w.Done()

	// check to see if an error was returned from the worker and log it if found
	err := w.Err()
	if err != nil {
		fmt.Printf("worker %s was stopped due to an error: %s\n", w.ID(), err)
	}
}
