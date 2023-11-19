package control

import (
	tskPkg "github.com/anatoliy9697/c2vocab/internal/model/task"
	res "github.com/anatoliy9697/c2vocab/internal/resources"
	// "github.com/anatoliy9697/c2vocab/internal/usecases"
)

type TaskHandler struct {
	Code string
	Res  res.Resources
}

func (th TaskHandler) Run(done chan string, tasks []tskPkg.Task) {
	defer func() { done <- th.Code }()

	th.Res.Logger = th.Res.Logger.With("taskHandlerCode", th.Code)

	// for _ = range tasks {
	// 1. Берем блокировку по задаче
	// 2. Обрабатываем
	// 3. Снимаем блокировку
	// }

	th.Res.Logger.Info("Task handler execution completed")

}
