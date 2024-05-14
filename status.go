package orar

import (
	"time"

	"github.com/robfig/cron/v3"
)

type StatusData struct {
	Id        cron.EntryID
	JobRunner *Job
	Next      time.Time
	Prev      time.Time
}


func Entries() []cron.Entry {
	return MainCron.Entries()
}

func StatusPage() []StatusData {

	ents := MainCron.Entries()

	Statuses := make([]StatusData, len(ents))
	for k, v := range ents {
		Statuses[k].Id = v.ID
		Statuses[k].JobRunner = AddJob(v.Job)
		Statuses[k].Next = v.Next
		Statuses[k].Prev = v.Prev

	}

	return Statuses
}

func StatusJson() map[string]interface{} {

	return map[string]interface{}{
		"jobrunner": StatusPage(),
	}

}

func AddJob(job cron.Job) *Job {
	return job.(*Job)
}
