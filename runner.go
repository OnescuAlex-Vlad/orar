package orar

import (
	"time"

	"github.com/robfig/cron/v3"
)

type Func func()

func Schedule(spec string, job cron.Job) error {
	sched, err := cron.ParseStandard(spec)
	if err != nil {
		return err
	}
	MainCron.Schedule(sched, New(job))
	return nil
}


func Every(duration time.Duration, job cron.Job) {

	MainCron.Schedule(cron.Every(duration), New(job))
}

// Run the given job right now.
func Now(job cron.Job) {
	go New(job).Run()
}

func In(duration time.Duration, job cron.Job) {
	go func() {
		time.Sleep(duration)
		New(job).Run()
	}()
}
