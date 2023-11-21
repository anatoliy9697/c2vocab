package control

import (
	tskPkg "github.com/anatoliy9697/c2vocab/internal/model/task"
	res "github.com/anatoliy9697/c2vocab/internal/resources"
	"github.com/anatoliy9697/c2vocab/internal/usecases"
)

type TaskHandler struct {
	Code string
	Res  res.Resources
}

func (th TaskHandler) Run(done chan string, tasks []tskPkg.Task) {
	defer func() { done <- th.Code }()

	th.Res.Logger = th.Res.Logger.With("taskHandlerCode", th.Code)

	var err error
	for _, task := range tasks {

		th.Res.Logger.Debug("Handling task", "task", task)
		if err = usecases.HandleTask(th.Res, task); err != nil {
			th.Res.Logger.Error(err.Error())
		}

		th.Res.Logger.Debug("Unlocking task", "task", task)
		if err = th.Res.TskRepo.UnlockTaskByUserId(task.UserId); err != nil {
			th.Res.Logger.Error(err.Error())
		}

	}

	th.Res.Logger.Info("Task handler execution completed")

}
