package one

import (
	"github.com/massarakhsh/lik"
	"github.com/jinzhu/gorm"
)

const KeyDepart = "depart"

//	Объект "Подразделение"
type Depart struct {
	One						//	Общий объект
	UpDepartId	int			//	ID вышестоящего подразделения
	Name        string		//	Наименование подразделения
	Notes       string		//	Примечание
}

//	Инициализация таблицы
func InitializeDepart() {
	if !ODB.HasTable(KeyDepart) {
		DBDepart().CreateTable(&Depart{})
	} else {
		DBDepart().AutoMigrate(&Depart{})
	}
}

//	Позиционирование интерфейса
func DBDepart() *gorm.DB {
	return ODB.Table(KeyDepart)
}

//	Получить объект
func GetDepart(id lik.IDB) (Depart,bool) {
	it := Depart{}
	ok := it.read(KeyDepart, id, &it)
	return it,ok
}

//	Новый объект
func NewDepart(datas... interface{}) (Depart,bool) {
	it := Depart{}
	ok := it.create(KeyDepart, it)
	if ok && len(datas) > 0 {
		ok = it.Update(datas...)
	}
	return it,ok
}

//	Выбрать объекты
func SelectDepart(query interface{}, args... interface{}) []Depart {
	var departs []Depart
	if query != nil {
		DBDepart().Where(query, args...).Find(&departs)
	} else {
		DBDepart().Find(&departs)
	}
	return departs
}

//	Получить таблицу
func (it *Depart) Table() string {
	return KeyDepart
}

//	Сохранить
func (it *Depart) Save() bool {
	return it.save(KeyDepart, it)
}

//	Изменить
func (it *Depart) Update(datas... interface{}) bool {
	return it.update(it, datas)
}

//	Удалить
func (it *Depart) Delete() {
	it.delete(KeyDepart, it)
}

