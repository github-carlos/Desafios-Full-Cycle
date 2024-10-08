package jobexecutor

import (
	"fmt"
	"trevas-bot/pkg/jobs"

	cron "github.com/robfig/cron/v3"
	"go.mau.fi/whatsmeow"
)

type JobExecutor interface {
	Execute(client *whatsmeow.Client)
	CronConfig() string
}

type JobsExecutor struct {
	client whatsmeow.Client
	jobsExecutors   []JobExecutor
	Cron   *cron.Cron
}

func NewJobsExecutor(client *whatsmeow.Client) *JobsExecutor {
  c := cron.New()

  var jobsExecutors []JobExecutor

  jobsExecutors = append(jobsExecutors, jobs.NewMemeSenderJob())

  for _, job := range jobsExecutors {
    fmt.Println("Adding job")
    c.AddFunc(job.CronConfig(), func() {
      fmt.Println("Running Job")
      job.Execute(client)
    })
  }

  c.Start()

	return &JobsExecutor{
		jobsExecutors: jobsExecutors,
    Cron: c,
	}
}
