//	Модуль списка сообщений.
package message

import (
	"bitbucket.org/961961/tsan/control/window"
	"bitbucket.org/961961/tsan/one"
	"bitbucket.org/961961/tsan/repo"
	"bitbucket.org/shaman/lik"
	"time"
)

//	Дескриптор списка
type ClientList struct {
	window.ClientBox
}

//	Интерфейс команд
type dealListExecute struct {
	It	*ClientList
}
func (it *dealListExecute) Run(rule *repo.DataRule, cmd string, data lik.Seter) {
	it.It.ListExecute(rule, cmd, data)
}

//	Интерфейс отображения таблицы
type dealShowGrid struct {
	It	*ClientList
}
func (it *dealShowGrid) Run(rule *repo.DataRule) {
	it.It.ListShowGrid(rule)
}

//	Интерфейс отображения страницы
type dealShowPage struct {
	It	*ClientList
}
func (it *dealShowPage) Run(rule *repo.DataRule) lik.Lister {
	return it.It.ListShowPage(rule)
}

//	Конструктор дескриптора
func BuildList(rule *repo.DataRule, frame string, id lik.IDB) *ClientList {
	it := &ClientList{}
	it.CollectInitialize(rule, frame, id,"message")
	it.ItExecute = &dealListExecute{it}
	it.ItShowGrid = &dealShowGrid{it}
	it.ItShowPage = &dealShowPage{it}
	it.IsLockRemote = true
	return it
}

//	Обработка команд
func (it *ClientList) ListExecute(rule *repo.DataRule, cmd string, data lik.Seter) {
	if cmd == "delete" {
		it.listDelete(rule)
	} else if cmd == "append" {
		it.listAppend(rule)
	} else {
		it.BoxExecute(rule, cmd, data)
	}
}

//	Отображение таблицы
func (it *ClientList) ListShowGrid(rule *repo.DataRule) {
	//it.Parameters.SetItem(64, "cellHeight")
	it.AddCmdTitle(rule, 100, "Сообщения")
	if rule.IAmAdmin() || rule.IAmAdvert() ||
		repo.ProbeItMy(rule, "member", it.IdMain) ||
		repo.ProbeItDep(rule, "member", it.IdMain) {
		it.AddCmdImg(rule, 200, "Добавить", "add", "/new_message_append")
	}
	if rule.IAmAdmin() || rule.IAmAdvert() {
		it.AddCmdImg(rule, 210, "Удалить", "del", "/ask_message_delete")
	}
	it.AddColumnItem(rule, "index=who", "title=Источник", "width=125", "autoHeight=true")
	it.AddColumnItem(rule, "index=body", "title=Сообщение", "width=500", "autoHeight=true")
	it.CollectShowGrid(rule)
}

//	Отображение страницы
func (it *ClientList) ListShowPage(rule *repo.DataRule) lik.Lister {
	rows := it.CollectShowPage(rule)
	var mess []one.Message
	one.DBMessage().Where("proto='public' AND offer_id=?", int(it.IdMain)).
		Order("updated_at desc").Find(&mess)
	for nm := 0; nm < len(mess); nm++ {
		msg := mess[nm]
		row := lik.BuildSet("id", msg.ID)
		src := msg.GetSource()
		src += "<br>" + time.Unix(int64(msg.TimeAt), 0).Format("2006/01/02 15:04")
		row.SetItem(src, "who")
		row.SetItem(msg.Body, "body")
		rows.AddItems(row)
	}
	return rows
}

//	Удаление элемента
func (it *ClientList) listDelete(rule *repo.DataRule) {
	if id := lik.StrToIDB(it.IdSelected); id > 0 {
		if msg,ok := one.GetMessage(id); ok {
			msg.Delete()
			rule.OnChangeData()
		}
	}
}

//	Добавление элемента
func (it *ClientList) listAppend(rule *repo.DataRule) {
	if text := lik.StringFromXS(rule.GetContext("text")); text != "" {
		message := &one.Message{Proto: "public", Scope: "member"}
		message.OfferId = it.IdMain
		message.Body = text
		message.TimeAt = int(time.Now().Unix())
		message.Save()
		it.IdSelected = lik.IntToStr(int(message.ID))
		rule.OnChangeData()
	}
}

