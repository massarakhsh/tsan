package one

import (
	"github.com/massarakhsh/lik"
	"github.com/jinzhu/gorm"
)

const KeyBonuses = "bonuses"

//	Объект "Контакт"
type Bonuses struct {
	One						//	Общий объект
}

//	Инициализация таблицы контактов
func InitializeBonuses() {
	if !ODB.HasTable(KeyBonuses) {
		DBBonuses().CreateTable(&Bonuses{})
	} else {
		DBBonuses().AutoMigrate(&Bonuses{})
	}
}

//	Позиционирование интерфейса
func DBBonuses() *gorm.DB {
	return ODB.Table(KeyBonuses)
}

//	Получить объект
func GetBonuses(id lik.IDB) (Bonuses,bool) {
	it := Bonuses{}
	ok := it.read(KeyBonuses, id, &it)
	return it,ok
}

//	Новый объект
func NewBonuses(datas... interface{}) (Bonuses,bool) {
	it := Bonuses{}
	ok := it.create(KeyBonuses, it)
	if ok && len(datas) > 0 {
		ok = it.Update(datas...)
	}
	return it,ok
}

//	Выбрать объекты
func SelectBonuses(query interface{}, args... interface{}) []Bonuses {
	var bonusess []Bonuses
	if query != nil {
		DBBonuses().Where(query, args...).Find(&bonusess)
	} else {
		DBBonuses().Find(&bonusess)
	}
	return bonusess
}

//	Получить таблицу
func (it *Bonuses) Table() string {
	return KeyBonuses
}

//	Сохранить
func (it *Bonuses) Save() bool {
	return it.save(KeyBonuses, it)
}

//	Изменить
func (it *Bonuses) Update(datas... interface{}) bool {
	return it.update(it, datas)
}

//	Удалить
func (it *Bonuses) Delete() {
	it.delete(KeyBonuses, it)
}

