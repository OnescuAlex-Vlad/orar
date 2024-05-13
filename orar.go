package orar

import (
	"bytes"
	"log"
	"reflect"
	"runtime/debug"
	"sync"
	"sync/atomic"

	"github.com/robfig/cron/v3"
)


type Job struct {
	Name    string
	inner   cron.Job
	status  uint32
	running sync.Mutex
}


var (
	// Singleton instance of the underlying job scheduler.
	MainCron *cron.Cron

	// This limits the number of jobs allowed to run concurrently.
	workPermits chan struct{}

	// Is a single job allowed to run concurrently with itself?
	selfConcurrent bool
)

const UNNAMED = "(unnamed)"

func New(job cron.Job) *Job {
	name := reflect.TypeOf(job).Name()
	if name == "Func" {
		name = UNNAMED
	}
	return &Job{
		Name:  name,
		inner: job,
	}
}

func (j *Job) Status() string {
	if atomic.LoadUint32(&j.status) > 0 {
		return "RUNNING"
	}
	return "IDLE"
}

func (j *Job) Run() {
	// If the job panics, just print a stack trace.
	// Don't let the whole process die.
	defer func() {
		if err := recover(); err != nil {
			//if revelError := revel.NewErrorFromPanic(err); revelError != nil {
			// revel.ERROR.Print(err, "\n", revelError.Stack)
			// } else {
			// revel.ERROR.Print(err, "\n", string(debug.Stack()))

			var buf bytes.Buffer
			logger := log.New(&buf, "JobRunner Log: ", log.Lshortfile)
			logger.Panic(err, "\n", string(debug.Stack()))
			// }
		}
	}()

	if !selfConcurrent {
		j.running.Lock()
		defer j.running.Unlock()
	}

	if workPermits != nil {
		workPermits <- struct{}{}
		defer func() { <-workPermits }()
	}

	atomic.StoreUint32(&j.status, 1)
	defer atomic.StoreUint32(&j.status, 0)

	j.inner.Run()
}