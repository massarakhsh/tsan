//	Модуль внутренней задачи.
package routine

import "sync"

//	Дескриптор задачи
type Task struct {
	Name	string		//	Наименование задачи
	isStoping	bool	//	Признак необходима остановка
	isStoped	bool	//	Признак задача остановлена
}

//	Интерфейс задачи
type Tasker interface {
	GetName() string		//	Имя задачи
	OnStoping()				//	Начать остановку задачи
	IsStoping() bool		//	Проверка. что задача останавливается
	OnStoped()				//	Задача остановлена
	IsStoped() bool			//	Проверка, что задача остановлена
}

var SyncList sync.Mutex		//	Семафор списка задач
var TaskList []Tasker		//	Список интерфейсов задач

//	Получить имя задачи
func (it *Task) GetName() string {
	return it.Name
}

//	Начать остановку задачи
func (it *Task) OnStoping() {
	it.isStoping = true
}

//	Проверка, что задача останавливается
func (it *Task) IsStoping() bool {
	return it.isStoping
}

//	Определить, что задача остановлена
func (it *Task) OnStoped() {
	it.isStoped = true
}

//	Проверка, что задача остановлена
func (it *Task) IsStoped() bool {
	return it.isStoped
}

//	Регистрация интерфейса задачи
func RegisterTask(task Tasker) {
	SyncList.Lock()
	TaskList = append(TaskList, task)
	SyncList.Unlock()
}

