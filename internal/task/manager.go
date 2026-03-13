package task

import (
	"xiaomiao-home-system/internal/conf"
	"xiaomiao-home-system/internal/data"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/robfig/cron/v3"
)

type TaskManager struct {
	log    *log.Helper
	config *conf.Config
	data   *data.Data
	task   *cron.Cron
	quit   chan struct{}
}

func NewTaskManager(config *conf.Config, data *data.Data, logger log.Logger) *TaskManager {
	return &TaskManager{
		config: config,
		data:   data,
		log:    log.NewHelper(log.With(logger, "task", "TaskManager")),
		task:   cron.New(),
		quit:   make(chan struct{}),
	}
}

func (u *TaskManager) Start() error {
	u.task.Start()

	// go func() {
	// 	for {
	// 		time.Sleep(5 * time.Second)

	// 		select {
	// 		case <-u.quit:
	// 			return
	// 		default:
	// 			// 执行DNS任务
	// 			u.RunDnsTask()
	// 		}
	// 	}
	// }()

	return nil
}

func (u *TaskManager) Stop() error {
	close(u.quit)
	u.task.Stop()

	return nil
}

// RunDnsTask 执行dns任务
func (u *TaskManager) RunDnsTask() {

	u.log.Infof("DNS任务执行成功")
}
