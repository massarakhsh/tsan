package one

import (
	"github.com/massarakhsh/lik"
	"github.com/jinzhu/gorm"
)

const KeyClient = "client"

//	Объект "Сделка"
type Client struct {
	One						//	Общий объект
}

//	Инициализация таблицы сделок
func InitializeClient() {
	if !ODB.HasTable(KeyClient) {
		DBClient().CreateTable(&Client{})
	} else {
		DBClient().AutoMigrate(&Client{})
	}
}

//	Позиционирование интерфейса
func DBClient() *gorm.DB {
	return ODB.Table(KeyClient)
}

//	Получить объект
func GetClient(id lik.IDB) (Client,bool) {
	it := Client{}
	ok := it.read(KeyClient, id, &it)
	return it,ok
}

//	Новый объект
func NewClient(datas... interface{}) (Client,bool) {
	it := Client{}
	ok := it.create(KeyClient, it)
	if ok && len(datas) > 0 {
		ok = it.Update(datas...)
	}
	return it,ok
}

//	Выбрать объекты
func SelectClient(query interface{}, args... interface{}) []Client {
	var clients []Client
	if query != nil {
		DBClient().Where(query, args...).Find(&clients)
	} else {
		DBClient().Find(&clients)
	}
	return clients
}

//	Получить таблицу
func (it *Client) Table() string {
	return KeyClient
}

//	Сохранить
func (it *Client) Save() bool {
	return it.save(KeyClient, it)
}

//	Изменить
func (it *Client) Update(datas... interface{}) bool {
	return it.update(it, datas)
}

//	Удалить
func (it *Client) Delete() {
	it.delete(KeyClient, it)
}

