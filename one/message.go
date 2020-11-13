package one

import (
	"bitbucket.org/shaman/lik"
	"github.com/jinzhu/gorm"
)

const KeyMessage = "message"

//	Объект "Сообщение"
type Message struct {
	One
	Proto		string		//	Тип сообщения
	Scope		string		//	Область сообщения
	Body		string		`gorm:"size:65535"` //	Тело сообщения
	OfferId		lik.IDB		//	Ключ заявки
	MemberId	lik.IDB		//	Ключ оператора
	TimeAt		int			//	Время сообщения
	ReadAt		int			//	Время доставки
}

//	Инициализация таблицы
func InitializeMessage() {
	if !ODB.HasTable(KeyMessage) {
		DBMessage().CreateTable(&Message{})
	} else {
		DBMessage().AutoMigrate(&Message{})
	}
}

//	Позиционирование интерфейса
func DBMessage() *gorm.DB {
	return ODB.Table(KeyMessage)
}

//	Получить объект
func GetMessage(id lik.IDB) (Message,bool) {
	it := Message{}
	ok := it.read(KeyMessage, id, &it)
	return it,ok
}

//	Новый объект
func NewMessage(datas... interface{}) (Message,bool) {
	it := Message{}
	ok := it.create(KeyMessage, it)
	if ok && len(datas) > 0 {
		ok = it.Update(datas...)
	}
	return it,ok
}

//	Выбрать объекты
func SelectMessage(query interface{}, args... interface{}) []Message {
	var messages []Message
	if query != nil {
		DBMessage().Where(query, args...).Find(&messages)
	} else {
		DBMessage().Find(&messages)
	}
	return messages
}

//	Получить таблицу
func (it *Message) Table() string {
	return KeyMessage
}

//	Сохранить
func (it *Message) Save() bool {
	return it.save(KeyMessage, it)
}

//	Изменить
func (it *Message) Update(datas... interface{}) bool {
	return it.update(it, datas)
}

//	Удалить
func (it *Message) Delete() {
	it.delete(KeyMessage, it)
}

func (it *Message) GetSource() string {
	source := ""
	if it.Proto == "public" {
		if it.Scope == "avito" {
			source = "Авито"
		} else if it.Scope == "yandex" {
			source = "Яндекс"
		} else if it.Scope == "cian" {
			source = "ЦИАН"
		} else if it.Scope == "tsan" {
			source = "ЦАН"
		} else {
			source = it.Scope
		}
	} else {
		source = it.Proto
	}
	return source
}