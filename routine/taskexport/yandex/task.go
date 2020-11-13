//	Процесс взаимодействия с площадкой Яндекс
package yandex

import (
	"github.com/massarakhsh/tsan/routine/taskexport"
)

//	Дескриптор задачи
type TaskYandex struct {
	taskexport.TaskExport
	AnswerAt	int
}

//	Статический уазатель на задачу
var ItTaskYandex *TaskYandex

//	Запуск задачи
func GoIt() {
	ItTaskYandex = &TaskYandex{}
	ItTaskYandex.Initialize("yandex", ItTaskYandex)
}

