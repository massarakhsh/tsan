package one

import (
	"bitbucket.org/shaman/lik"
	"github.com/jinzhu/gorm"
)

const KeyDeal = "deal"

//	Объект "Сделка"
type Deal struct {
	One						//	Общий объект
}

//	Инициализация таблицы сделок
func InitializeDeal() {
	if !ODB.HasTable(KeyDeal) {
		DBDeal().CreateTable(&Deal{})
	} else {
		DBDeal().AutoMigrate(&Deal{})
	}
}

//	Позиционирование интерфейса
func DBDeal() *gorm.DB {
	return ODB.Table(KeyDeal)
}

//	Получить объект
func GetDeal(id lik.IDB) (Deal,bool) {
	it := Deal{}
	ok := it.read(KeyDeal, id, &it)
	return it,ok
}

//	Новый объект
func NewDeal(datas... interface{}) (Deal,bool) {
	it := Deal{}
	ok := it.create(KeyDeal, it)
	if ok && len(datas) > 0 {
		ok = it.Update(datas...)
	}
	return it,ok
}

//	Выбрать объекты
func SelectDeal(query interface{}, args... interface{}) []Deal {
	var deals []Deal
	if query != nil {
		DBDeal().Where(query, args...).Find(&deals)
	} else {
		DBDeal().Find(&deals)
	}
	return deals
}

//	Получить таблицу
func (it *Deal) Table() string {
	return KeyDeal
}

//	Сохранить
func (it *Deal) Save() bool {
	return it.save(KeyDeal, it)
}

//	Изменить
func (it *Deal) Update(datas... interface{}) bool {
	return it.update(it, datas)
}

//	Удалить
func (it *Deal) Delete() {
	it.delete(KeyDeal, it)
}

