package one

import (
	"bitbucket.org/shaman/lik"
	"github.com/jinzhu/gorm"
)

const KeyOffer = "offer"

//	Объект "Заявка"
type Offer struct {
	One
	MemberId	int			//	ID ответственного риэлтора
	Cost		int			//	Текущая цена
	Notes       string		//	Примечание
}

//	Инициализация таблицы
func InitializeOffer() {
	if !ODB.HasTable(KeyOffer) {
		DBOffer().CreateTable(&Offer{})
	} else {
		DBOffer().AutoMigrate(&Offer{})
	}
}

//	Позиционирование интерфейса
func DBOffer() *gorm.DB {
	return ODB.Table(KeyOffer)
}

//	Получить объект
func GetOffer(id lik.IDB) (Offer,bool) {
	it := Offer{}
	ok := it.read(KeyOffer, id, &it)
	return it,ok
}

//	Новый объект
func NewOffer(datas... interface{}) (Offer,bool) {
	it := Offer{}
	ok := it.create(KeyOffer, it)
	if ok && len(datas) > 0 {
		ok = it.Update(datas...)
	}
	return it,ok
}

//	Выбрать объекты
func SelectOffer(query interface{}, args... interface{}) []Offer {
	var offers []Offer
	if query != nil {
		DBOffer().Where(query, args...).Find(&offers)
	} else {
		DBOffer().Find(&offers)
	}
	return offers
}

//	Получить таблицу
func (it *Offer) Table() string {
	return KeyOffer
}

//	Сохранить
func (it *Offer) Save() bool {
	return it.save(KeyOffer, it)
}

//	Изменить
func (it *Offer) Update(datas... interface{}) bool {
	return it.update(it, datas)
}

//	Удалить
func (it *Offer) Delete() {
	it.delete(KeyOffer, it)
}

