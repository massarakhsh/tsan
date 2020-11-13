// Инструментальный интерфейс базы данных.
package jone

import (
	"bitbucket.org/shaman/lik"
	"bitbucket.org/shaman/lik/likbase"
)

//	Дескриптор таблицы базы данных
type FlatTable struct {
	Part	string
	Fields	[]likbase.DBField
}

//	Дескриптор базы данных
var DB likbase.JsonBaser

//	Декрипторы резидентных таблиц
var TableSystem, TableParam, TableDepart, TableMember, TableComplex,
	TableObject, TableBell, TableClient, TableOffer, TableDeal  *likbase.ItTable
var ListTables   []*likbase.ItTable

//	Дескриптор учетной записи системы
var SysElm *likbase.ItElm

//	Флаги применения полей
const (
	TagGrid    = 0x01	//	Поле таблицы
	TagHide    = 0x02	//	Скрытое поле
	TagForm    = 0x04	//	Поле формы
	TagShow    = 0x08	//	Поле презентации
	TagEdit    = 0x10	//	Редактируемое поле
	TagSale    = 0x20	//	Поле продаж
	TagBuy     = 0x40	//	Поле покупок
	TagMust    = 0x80	//	Обязательное поле
	TagFlat    = 0x100	//	Поле квартир и комнат
	TagHouse   = 0x200	//	Поле домов
	TagLand    = 0x400	//	Поле участков
	TagTune    = 0x8000	//	Поле настроек
	Tag_Target = TagSale | TagBuy
	Tag_Realty = TagFlat | TagHouse | TagLand
)

//	Инициализация базы данных
func GoIt(serv string, base string, user string, pass string) {
	DB = likbase.OpenJsonBase(serv, base, user, pass)
	InitTables()
	StartBase()
}

//	Инициализация таблиц
func InitTables() {
	TableSystem = DB.BuildTable("system", "Система")
	TableParam = DB.BuildTable("param", "Параметры")
	TableDepart = DB.BuildTable("depart", "Отделы")
	TableMember = DB.BuildTable("member", "Сотрудники")
	TableComplex = DB.BuildTable("complex", "Комплексы")
	TableObject = DB.BuildTable("object", "Объекты")
	TableBell = DB.BuildTable("bell", "Обращения")
	TableClient = DB.BuildTable("client", "Клиенты")
	TableOffer = DB.BuildTable("offer", "Заявки")
	TableDeal = DB.BuildTable("deal", "Сделки")
	ListTables = []*likbase.ItTable{
		TableSystem, TableParam,
		TableDepart, TableMember, TableClient,
		TableComplex, TableObject,
		TableBell, TableOffer, TableDeal,
	}
}

//	Запуск базы данных
func StartBase() {
	LoadListTables()
	UpgradeListTables()
}

//	Останов базы данных
func StopBase() {
	DB.StopDB()
}

//	ПОлучить таблицу по имени
func GetTable(part string) *likbase.ItTable {
	if part == "receptor" {
		part = "member"
	} else if part == "sale" {
		part = "offer"
	} else if part == "buy" {
		part = "offer"
	}
	for _, table := range ListTables {
		if table.Part == part {
			return table
		}
	}
	return nil
}

//	Получить элемент
func GetElm(part string, id lik.IDB) *likbase.ItElm {
	var elm *likbase.ItElm
	if table := GetTable(part); table != nil {
		elm = table.GetElm(id)
	}
	return elm
}

//	Удалить элемент
func DeleteElm(part string, id lik.IDB) {
	if table := GetTable(part); table != nil {
		table.DeleteElm(id)
	}
}

//	Загрузить список таблиц
func LoadListTables() {
	SysElm = nil
	for _, table := range ListTables {
		table.LoadElms()
	}
	for id,elm := range TableSystem.Elms {
		if SysElm == nil || id > SysElm.Id {
			SysElm = elm
		}
	}
}

//	Обновить список таблиц
func UpgradeListTables() {
	UpgradeStart()
	for _, table := range ListTables {
		for _, elm := range table.Elms {
			UpgradeElm(elm)
		}
		UpgradeTable(table)
	}
	UpgradeStop()
}

//	Очистить список таблиц
func PurgeListTables() {
	for _, table := range ListTables {
		table.Purge()
	}
}

