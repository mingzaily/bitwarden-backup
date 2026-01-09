package scheduler

import (
	"log"
	"sync"

	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	cron        *cron.Cron
	taskEntries map[uint]cron.EntryID // 任务ID -> cron entry ID 映射
	mu          sync.RWMutex          // 保护 taskEntries 的并发访问
}

func New() *Scheduler {
	return &Scheduler{
		cron:        cron.New(cron.WithSeconds()),
		taskEntries: make(map[uint]cron.EntryID),
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
