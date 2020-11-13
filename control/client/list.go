//	Модуль списка клиеннтов
package client

import (
	"bitbucket.org/961961/tsan/control"
	"bitbucket.org/961961/tsan/fancy"
	"bitbucket.org/961961/tsan/repo"
	"bitbucket.org/shaman/lik"
	"bitbucket.org/shaman/lik/likbase"
)

//	Дескриптор списка клиентов
type ClientControl struct {
	control.ListControl
}

//	Обработчик заполнения таблицы
type dealClientGridFill struct {
	It	*ClientControl
}
func (it *dealClientGridFill) Run(rule *repo.DataRule) {
	it.It.ClientGridFill(rule)
}

//	Обработчик проверки строки
type dealClientListRowsProbe struct {
	It	*ClientControl
}
func (it *dealClientListRowsProbe) Run(rule *repo.DataRule, elm *likbase.ItElm) bool {
	return it.It.ClientListRowsProbe(rule, elm)
}

//	Обработчик заполнения формы
type dealClientElmForm struct {
	It	*ClientControl
}
func (it *dealClientElmForm) Run(rule *repo.DataRule, elm *likbase.ItElm) {
	it.It.ClientElmForm(rule, elm)
}

//	Обработчик заполнения строки
type dealClientRowFill struct {
	It	*ClientControl
}
func (it *dealClientRowFill) Run(rule *repo.DataRule, elm *likbase.ItElm, row lik.Seter) {
	it.It.ClientRowFill(rule, elm, row)
}

//	Обработчик обработки событий
type dealClientExecute struct {
	It	*ClientControl
}
func (it *dealClientExecute) Run(rule *repo.DataRule, cmd string, data lik.Seter) {
	it.It.ClientExecute(rule, cmd, data)
}

//	Конструктор списка клиентов
func BuildClientList(rule *repo.DataRule, id lik.IDB) *ClientControl {
	it := &ClientControl{}
	it.ControlInitialize("client", id)
	it.ListInitialize(rule, "client", "client")
	it.ItGridFill = &dealClientGridFill{it}
	it.ItListMakeProbe = &dealClientListRowsProbe{it}
	it.ItRowFill = &dealClientRowFill{it}
	it.ItFormElm = &dealClientElmForm{it}
	it.ItExecute = &dealClientExecute{it}
	return it
}

//	Обработка событий
func (it *ClientControl) ClientExecute(rule *repo.DataRule, cmd string, data lik.Seter) {
	if cmd == "all" || cmd == it.Main {
		it.ClientExecute(rule, rule.Shift(), data)
	} else if cmd == "toshow" || cmd == "toenter" {
		if sel := rule.Top(); sel != "" {
			it.Sel = sel
		}
		it.GoWindowMode(rule, "clientcard", it.Sel)
	} else if cmd == "newclient" {
		path := "/clientcard0"
		rule.SetResponse(path, "_function_lik_window_part")
	} else {
		it.ListExecute(rule, cmd, data)
	}
}

//	Заполнение таблицы
func (it *ClientControl) ClientGridFill(rule *repo.DataRule) {
	it.ListGridFill(rule)
	if rule.IAmAdmin() {
		it.AddCommandImg(rule, fancy.OrdUnit+1, "Создать запись клиента", "toadd", "add")
	}
}

//	Заполнение формы
func (it *ClientControl) ClientElmForm(rule *repo.DataRule, elm *likbase.ItElm) {
	it.ListFormElm(rule, elm)
}

//	Проверка строки
func (it *ClientControl) ClientListRowsProbe(rule *repo.DataRule, elm *likbase.ItElm) bool {
	accept := (elm != nil)
	if accept && !it.ListMakeProbe(rule, elm) {
		accept = false
	}
	return accept
}

//	Заполнение строки
func (it *ClientControl) ClientRowFill(rule *repo.DataRule, elm *likbase.ItElm, row lik.Seter) bool {
	it.ListRowFill(rule, elm, row)
	if row == nil {
		return false
	}
	//row.SetItem(fmt.Sprintf("/%s%d/show%d?_tp=1", it.Main, int(elm.Id), int(elm.Id)), "pathopen")
	return true
}

