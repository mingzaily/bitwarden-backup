package scheduler

import (
	"sync"
	"time"

	"github.com/mingzaily/bitwarden-backup/internal/logger"
	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	cron        *cron.Cron
	taskEntries map[uint]cron.EntryID // 任务ID -> cron entry ID 映射
	mu          sync.RWMutex          // 保护 taskEntries 的并发访问

	taskQueue   chan uint
	queuedTasks map[uint]bool
	queueMu     sync.Mutex
	stopChan    chan struct{}
	workerDone  chan struct{}
}

func New() *Scheduler {
	return &Scheduler{
		cron:        cron.New(cron.WithSeconds()),
		taskEntries: make(map[uint]cron.EntryID),
		taskQueue:   make(chan uint, 100),
		queuedTasks: make(map[uint]bool),
		stopChan:    make(chan struct{}),
		workerDone:  make(chan struct{}),
	}
}

func (s *Scheduler) Start() {
	logger.Module(logger.ModuleScheduler).Info("Starting scheduler")
	s.startWorker()
	s.cron.Start()
	logger.Module(logger.ModuleScheduler).Info("Scheduler started")
}

func (s *Scheduler) Stop() {
	logger.Module(logger.ModuleScheduler).Info("Stopping scheduler")
	s.cron.Stop()
	close(s.stopChan)

	select {
	case <-s.workerDone:
		logger.Module(logger.ModuleScheduler).Info("Worker stopped gracefully")
	case <-time.After(30 * time.Second):
		logger.Module(logger.ModuleScheduler).Error("Worker stop timeout")
	}

	logger.Module(logger.ModuleScheduler).Info("Scheduler stopped")
}

func (s *Scheduler) startWorker() {
	go func() {
		defer close(s.workerDone)
		logger.Module(logger.ModuleScheduler).Info("Task queue worker started")

		for {
			select {
			case taskID := <-s.taskQueue:
				s.processTask(taskID)
			case <-s.stopChan:
				logger.Module(logger.ModuleScheduler).Info("Worker received stop signal, draining queue")
				for {
					select {
					case taskID := <-s.taskQueue:
						s.processTask(taskID)
					default:
						logger.Module(logger.ModuleScheduler).Info("Worker stopped")
						return
					}
				}
			}
		}
	}()
}
