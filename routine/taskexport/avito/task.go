//	Процесс взаимодействия с площадкой Авито
package avito

import (
	"github.com/massarakhsh/tsan/routine/taskexport"
)

//	Дескриптор задачи
type TaskAvito struct {
	taskexport.TaskExport
}

//	Статический указатель на дескриптор
var ItTaskAvito *TaskAvito

//	Запуск задачи
func GoIt() {
	ItTaskAvito = &TaskAvito{}
	ItTaskAvito.Initialize("avito", ItTaskAvito)
}

