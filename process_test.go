package gopool

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"time"
)

func TestProcess_Start_Stop(t *testing.T) {
	a := assert.New(t)

	p := NewProcess("process-1", func(process *Process, commands <-chan ProcessCommand) error {
		for {
			select {
			case cmd := <-commands:
				if cmd == StopProcessCommand {
					time.Sleep(time.Millisecond * 40)
					return nil
				}
			}
		}
		return nil
	})

	a.Equal(ProcessStopped, p.Status())

	p.Start()
	time.Sleep(time.Millisecond * 10) // sleep for just enough time for the go routine to start up
	a.Equal(ProcessRunning, p.Status())

	err := p.Start()
	a.EqualError(err, "process is not stopped: running")

	p.Stop()
	time.Sleep(time.Millisecond * 50)
	a.Equal(ProcessStopped, p.Status())
}

func TestProcess_Start_Finish(t *testing.T) {
	a := assert.New(t)

	p := NewProcess("process-1", func(process *Process, commands <-chan ProcessCommand) error {
		time.Sleep(time.Millisecond * 100)
		return nil
	})
	finChan := p.FinishedChan()

	a.Equal(ProcessStopped, p.Status())

	p.Start()
	time.Sleep(time.Millisecond * 10) // sleep for just enough time for the go routine to start up
	a.Equal(ProcessRunning, p.Status())

	<-finChan
	time.Sleep(time.Millisecond * 100)
	a.Equal(ProcessFinished, p.Status())
}

func TestProcess_Stop(t *testing.T) {
	a := assert.New(t)

	output := make(chan string, 100)

	p := NewProcess("process-1", func(process *Process, commands <-chan ProcessCommand) error {
		for {
			select {
			case cmd := <-commands:
				output <- "stop-" + process.ID()
				if cmd == StopProcessCommand {
					close(output)
					return nil
				}
			default:
				break
			}

			output <- process.ID()
			time.Sleep(time.Millisecond * 25)
		}
		return nil
	})

	p.Start()

	time.Sleep(time.Millisecond * 100)
	p.Stop()
	time.Sleep(time.Millisecond * 100)

	outputArr := make([]string, 100)
	indexToUse := 0

outerLoop:
	for {
		select {
		case o, open := <-output:
			if ! open {
				break outerLoop
			}

			outputArr[indexToUse] = o
			indexToUse++
		}
	}

	for k, v := range outputArr {
		switch k {
		case 0:
			a.Equal("process-1", v)
		case 1:
			a.Equal("process-1", v)
		case 2:
			a.Equal("process-1", v)
		case 3:
			a.Equal("process-1", v)
		case 4:
			a.Equal("stop-process-1", v)
		default:
			a.Equal("", v)
		}
	}
}
