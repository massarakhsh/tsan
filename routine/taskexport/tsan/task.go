//	Процесс взаимодействия с площадкой ЦАН
package tsan

import (
	"bitbucket.org/961961/tsan/routine/taskexport"
)

type TaskTsan struct {
	taskexport.TaskExport
	ListKeys	[]string		//	Список необходимых полей
	//	Фотографии добавляются после этого списка
}

//	Статический указатель на задачу
var ItTaskTsan *TaskTsan

//	Запуск задачи
func GoIt() {
	ItTaskTsan = &TaskTsan{}
	ItTaskTsan.Initialize("tsan", ItTaskTsan)
}

