package one

import (
	"github.com/massarakhsh/lik"
	"github.com/jinzhu/gorm"
)

const KeyMember = "member"

//	Объект "Сотрудник"
type Member struct {
	One
	Family      string		//	Фамилия
	Namely      string		//	Имя
	Paterly     string		//	Отчество
	Phone       string		//	Телефон
	DepartId	int			//	ID подразделения
	Role        string		//	Роль
	Email       string		//	Почта
	ProPhone    string		//	Телефон рекламы
	Pin         string		//	Внутренний телефон
	Photo       string		//	Фотография
	Notes       string		//	Примечание
}

//	Инициализация таблицы
func InitializeMember() {
	if !ODB.HasTable(KeyMember) {
		DBMember().CreateTable(&Member{})
	} else {
		DBMember().AutoMigrate(&Member{})
	}
}

//	Позиционирование интерфейса
func DBMember() *gorm.DB {
	return ODB.Table(KeyMember)
}

//	Получить объект
func GetMember(id lik.IDB) (Member,bool) {
	it := Member{}
	ok := it.read(KeyMember, id, &it)
	return it,ok
}

//	Новый объект
func NewMember(datas... interface{}) (Member,bool) {
	it := Member{}
	ok := it.create(KeyMember, it)
	if ok && len(datas) > 0 {
		ok = it.Update(datas...)
	}
	return it,ok
}

//	Выбрать объекты
func SelectMember(query interface{}, args... interface{}) []Member {
	var members []Member
	if query != nil {
		DBMember().Where(query, args...).Find(&members)
	} else {
		DBMember().Find(&members)
	}
	return members
}

//	Получить таблицу
func (it *Member) Table() string {
	return KeyMember
}

//	Сохранить
func (it *Member) Save() bool {
	return it.save(KeyMember, it)
}

//	Изменить
func (it *Member) Update(datas... interface{}) bool {
	return it.update(it, datas)
}

//	Удалить
func (it *Member) Delete() {
	it.delete(KeyMember, it)
}

