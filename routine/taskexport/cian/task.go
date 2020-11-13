//	Процесс взаимодействия с площадкой ЦИАН
package cian

import (
	"github.com/massarakhsh/tsan/routine/taskexport"
)

//	Дескриптор задачи
type TaskCian struct {
	taskexport.TaskExport
}

//	Статический указатель на задачу
var ItTaskCian *TaskCian

//	Запуск задачи
func GoIt() {
	ItTaskCian = &TaskCian{}
	ItTaskCian.Initialize("cian", ItTaskCian)
}

