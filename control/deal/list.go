//	Модуль списка сделок
package listdeal

import (
	"github.com/massarakhsh/tsan/control"
	"github.com/massarakhsh/tsan/repo"
	"github.com/massarakhsh/lik"
	"github.com/massarakhsh/lik/likbase"
)

//	Дескриптор списка сделок
type DealControl struct {
	control.ListControl
}

//	Обработчик заполнения таблицы
type dealDealGridFill struct {
	It	*DealControl
}
func (it *dealDealGridFill) Run(rule *repo.DataRule) {
	it.It.DealGridFill(rule)
}

//	Обработчик проверки строк
type dealDealListRowsProbe struct {
	It	*DealControl
}
func (it *dealDealListRowsProbe) Run(rule *repo.DataRule, elm *likbase.ItElm) bool {
	return it.It.DealListRowsProbe(rule, elm)
}

//	Обработчик заполнения формы
type dealDealElmForm struct {
	It	*DealControl
}
func (it *dealDealElmForm) Run(rule *repo.DataRule, elm *likbase.ItElm) {
	it.It.DealElmForm(rule, elm)
}

//	Обработчик заполнения строки
type dealDealRowFill struct {
	It	*DealControl
}
func (it *dealDealRowFill) Run(rule *repo.DataRule, elm *likbase.ItElm, row lik.Seter) {
	it.It.DealRowFill(rule, elm, row)
}

//	Интерфейс обработки команд
type dealDealExecute struct {
	It	*DealControl
}
func (it *dealDealExecute) Run(rule *repo.DataRule, cmd string, data lik.Seter) {
	it.It.DealExecute(rule, cmd, data)
}

//	Конструктор дескриптора списка сделок
func BuildListDeal(rule *repo.DataRule, id lik.IDB) *DealControl {
	it := &DealControl{}
	it.ControlInitialize("deal", id)
	it.ListInitialize(rule, "deal", "deal")
	it.ItGridFill = &dealDealGridFill{it}
	it.ItListMakeProbe = &dealDealListRowsProbe{it}
	it.ItRowFill = &dealDealRowFill{it}
	it.ItFormElm = &dealDealElmForm{it}
	it.ItExecute = &dealDealExecute{it}
	return it
}

//	Обработка команд
func (it *DealControl) DealExecute(rule *repo.DataRule, cmd string, data lik.Seter) {
	if cmd == "all" || cmd == it.Main {
		it.DealExecute(rule, rule.Shift(), data)
	} else {
		it.ListExecute(rule, cmd, data)
	}
}

//	Заполнение таблицы
func (it *DealControl) DealGridFill(rule *repo.DataRule) {
	it.ListGridFill(rule)
}

//	Заполнение формы
func (it *DealControl) DealElmForm(rule *repo.DataRule, elm *likbase.ItElm) {
	it.ListFormElm(rule, elm)
}

//	Проверка строк
func (it *DealControl) DealListRowsProbe(rule *repo.DataRule, elm *likbase.ItElm) bool {
	accept := (elm != nil)
	if accept && !it.ListMakeProbe(rule, elm) {
		accept = false
	}
	return accept
}

//	Заполнение строки
func (it *DealControl) DealRowFill(rule *repo.DataRule, elm *likbase.ItElm, row lik.Seter) bool {
	it.ListRowFill(rule, elm, row)
	if row == nil {
		return false
	}
	//row.SetItem(fmt.Sprintf("/%s%d/show%d?_tp=1", it.Main, int(elm.Id), int(elm.Id)), "pathopen")
	return true
}

