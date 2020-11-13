package repo

import (
	"bitbucket.org/shaman/lik"
)

//	Дескриптор ядра контроллера
type DataControl struct {
	Frame		string	//	Ключ фрейма
	Mode      	string	//	Ключ режима
	Mark		string	//	Ключ марки
	IdMain    	lik.IDB	//	Основной объект
	Stable		bool
}

//	Интерфейс ядра контроллера
type Controller interface {
	GetFrame() string		//	Фрейм
	GetMark() string		//	Метка
	GetMode() string		//	Режим
	SetMark(mark string)	//	Установить метку
	GetId() lik.IDB			//	Получить идентификатор
	SetId(id lik.IDB)		//	Установить идентификатор
	IsStable() bool
}

//	Инициализировать ядро
func (it *DataControl) RepoInitialize(frame string, id lik.IDB, mode string) {
	it.Frame = frame
	it.IdMain = id
	it.Mode = mode
}

//	Получить фрейм
func (it *DataControl) GetFrame() string {
	return it.Frame
}

//	Получить метку
func (it *DataControl) GetMark() string {
	return it.Mark
}

//	Получить режим
func (it *DataControl) GetMode() string {
	return it.Mode
}

//	Установить метку
func (it *DataControl) SetMark(mark string) {
	it.Mark = mark
}

//	Получить идентификатор
func (it *DataControl) GetId() lik.IDB {
	return it.IdMain
}

//	Установить идентификатор
func (it *DataControl) SetId(id lik.IDB) {
	it.IdMain = id
}

//	Определить стабильность
func (it *DataControl) IsStable() bool {
	return it.Stable
}

