package scheduler

import (
	"log"

	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	cron *cron.Cron
}

func New() *Scheduler {
	return &Scheduler{
		cron: cron.New(cron.WithSeconds()),
	}
}

func (s *Scheduler) Start() {
	log.Println("Scheduler started")
	s.cron.Start()
}

func (s *Scheduler) Stop() {
	log.Println("Scheduler stopped")
	s.cron.Stop()
}
