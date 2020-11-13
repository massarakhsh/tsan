package one

import (
	"bitbucket.org/shaman/lik"
	"github.com/jinzhu/gorm"
)

const KeyBell = "bell"

//	Объект "Контакт"
type Bell struct {
	One						//	Общий объект
}

//	Инициализация таблицы контактов
func InitializeBell() {
	if !ODB.HasTable(KeyBell) {
		DBBell().CreateTable(&Bell{})
	} else {
		DBBell().AutoMigrate(&Bell{})
	}
}

//	Позиционирование интерфейса
func DBBell() *gorm.DB {
	return ODB.Table(KeyBell)
}

//	Получить объект
func GetBell(id lik.IDB) (Bell,bool) {
	it := Bell{}
	ok := it.read(KeyBell, id, &it)
	return it,ok
}

//	Новый объект
func NewBell(datas... interface{}) (Bell,bool) {
	it := Bell{}
	ok := it.create(KeyBell, it)
	if ok && len(datas) > 0 {
		ok = it.Update(datas...)
	}
	return it,ok
}

//	Выбрать объекты
func SelectBell(query interface{}, args... interface{}) []Bell {
	var bells []Bell
	if query != nil {
		DBBell().Where(query, args...).Find(&bells)
	} else {
		DBBell().Find(&bells)
	}
	return bells
}

//	Получить таблицу
func (it *Bell) Table() string {
	return KeyBell
}

//	Сохранить
func (it *Bell) Save() bool {
	return it.save(KeyBell, it)
}

//	Изменить
func (it *Bell) Update(datas... interface{}) bool {
	return it.update(it, datas)
}

//	Удалить
func (it *Bell) Delete() {
	it.delete(KeyBell, it)
}

