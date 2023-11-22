package control

import (
	"context"
	"time"

	tskPkg "github.com/anatoliy9697/c2vocab/internal/model/task"
	res "github.com/anatoliy9697/c2vocab/internal/resources"
	"github.com/google/uuid"
)

type Scheduler struct {
	MaxTaskHandlers int
	TaskWaitingTime int // ms
	TaskBatchSize   int
	Res             res.Resources
}

func (s Scheduler) Run(ctx context.Context, done chan struct{}) {
	defer func() { done <- struct{}{} }()

	s.Res.Logger.Info("Scheduler is running")

	handlers := make(map[string]struct{}, s.MaxTaskHandlers)
	handlerDone := make(chan string, 10)
	handlerCode := ""

	noTaskTicker := time.NewTicker(time.Millisecond * time.Duration(s.TaskWaitingTime))
	haveTaskTicker := time.NewTicker(time.Millisecond)
	ticker := haveTaskTicker

	var tasks []tskPkg.Task
	var err error

loop:
	for {
		select {

		// Scheduler shutdown
		case <-ctx.Done():
			break loop

		// Handler had finished
		case handlerCode = <-handlerDone:
			delete(handlers, handlerCode)

		// Next scheduler iteration
		case <-ticker.C:
			if len(handlers) >= s.MaxTaskHandlers {
				s.Res.Logger.Info("No free task handlers. Waiting for handler")
				handlerCode = <-handlerDone
				delete(handlers, handlerCode)
			}

			handlerCode = uuid.NewString()[:7]

			if tasks, err = s.Res.TskRepo.TasksWithLocking("taskHandler-"+handlerCode, s.TaskBatchSize, s.Res.LockConf.TimeForReassign); err != nil {
				s.Res.Logger.Error("Scheduler fatal error: " + err.Error())
				panic(err) // TODO: Надо бы сделать отлов паники
			}

			if len(tasks) > 0 {
				handlers[handlerCode] = struct{}{}
				s.Res.Logger.Info("Running task handler "+handlerCode, "tasks", tasks)
				go TaskHandler{
					Code: handlerCode,
					Res:  s.Res,
				}.Run(handlerDone, tasks)
				ticker = haveTaskTicker
			} else {
				s.Res.Logger.Debug("No scheduler tasks")
				ticker = noTaskTicker
			}

		}

	}

	// Whaiting for rest handlers finishing
	for len(handlers) > 0 {
		handlerCode := <-handlerDone
		delete(handlers, handlerCode)
	}

	s.Res.Logger.Info("Scheduler execution completed")
}
