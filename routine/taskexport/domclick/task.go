//	Процесс взаимодействия с площадкой Авито
package domclick

import (
	"bitbucket.org/961961/tsan/routine/taskexport"
)

//	Дескриптор задачи
type TaskDomclick struct {
	taskexport.TaskExport
}

//	Статический указатель на дескриптор
var ItTaskDomclick *TaskDomclick

//	Запуск задачи
func GoIt() {
	ItTaskDomclick = &TaskDomclick{}
	ItTaskDomclick.Initialize("domclick", ItTaskDomclick)
}
