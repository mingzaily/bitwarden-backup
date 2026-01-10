package scheduler

import (
	"sync"

	"github.com/mingzaily/bitwarden-backup/internal/logger"
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
	logger.Module(logger.ModuleScheduler).Info("Scheduler started")
	s.cron.Start()
}

func (s *Scheduler) Stop() {
	logger.Module(logger.ModuleScheduler).Info("Scheduler stopped")
	s.cron.Stop()
}
