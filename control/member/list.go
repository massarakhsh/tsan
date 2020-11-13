//	Модуль списка сотрудников
package member

import (
	"github.com/massarakhsh/tsan/control"
	"github.com/massarakhsh/tsan/fancy"
	"github.com/massarakhsh/tsan/repo"
	"github.com/massarakhsh/lik"
	"github.com/massarakhsh/lik/likbase"
)

//	Дескриптор списка сотрудников
type MemberControl struct {
	control.ListControl
}

//	Интерфейс вывода таблицы
type dealMemberGridFill struct {
	It	*MemberControl
}
func (it *dealMemberGridFill) Run(rule *repo.DataRule) {
	it.It.MemberGridFill(rule)
}

//	Интерфейс проверки строки
type dealMemberListRowsProbe struct {
	It	*MemberControl
}
func (it *dealMemberListRowsProbe) Run(rule *repo.DataRule, elm *likbase.ItElm) bool {
	return it.It.MemberListRowsProbe(rule, elm)
}

//	Интерфейс вывода формы
type dealMemberElmForm struct {
	It	*MemberControl
}
func (it *dealMemberElmForm) Run(rule *repo.DataRule, elm *likbase.ItElm) {
	it.It.MemberElmForm(rule, elm)
}

//	Интерфейс вывода строки
type dealMemberRowFill struct {
	It	*MemberControl
}
func (it *dealMemberRowFill) Run(rule *repo.DataRule, elm *likbase.ItElm, row lik.Seter) {
	it.It.MemberRowFill(rule, elm, row)
}

//	Интерфейс входа в троку
type dealMemberEnterRow struct {
	It	*MemberControl
}
func (it *dealMemberEnterRow) Run(rule *repo.DataRule) {
	it.It.MemberEnterRow(rule)
}

//	Интерфейс исполнения команд
type dealMemberExecute struct {
	It	*MemberControl
}
func (it *dealMemberExecute) Run(rule *repo.DataRule, cmd string, data lik.Seter) {
	it.It.MemberExecute(rule, cmd, data)
}

//	Конструктор лескриптора
func BuildMemberList(rule *repo.DataRule, id lik.IDB) *MemberControl {
	it := &MemberControl{}
	it.ControlInitialize("member", id)
	it.ListInitialize(rule, "member", "member")
	it.ItGridFill = &dealMemberGridFill{it}
	it.ItListMakeProbe = &dealMemberListRowsProbe{it}
	it.ItRowFill = &dealMemberRowFill{it}
	//it.ItEnterRow = &dealMemberEnterRow{it}
	it.ItFormElm = &dealMemberElmForm{it}
	it.ItExecute = &dealMemberExecute{it}
	return it
}

//	Исполнение команд
func (it *MemberControl) MemberExecute(rule *repo.DataRule, cmd string, data lik.Seter) {
	if cmd == "all" || cmd == it.Main {
		it.MemberExecute(rule, rule.Shift(), data)
	} else if cmd == "toshow" || cmd == "toenter" {
		if sel := rule.Top(); sel != "" {
			it.Sel = sel
		}
		it.GoWindowMode(rule, "member/membercard", it.Sel)
	} else if cmd == "newmember" {
		path := "/member/membercard0"
		rule.SetResponse(path, "_function_lik_window_part")
	} else {
		it.ListExecute(rule, cmd, data)
	}
}

//	Вывод таблицы
func (it *MemberControl) MemberGridFill(rule *repo.DataRule) {
	it.ListGridFill(rule)
	if rule.IAmAdmin() {
		it.AddCommandItem(rule, fancy.OrdUnit+9, lik.BuildSet(
			"type=button", "tip=Кабинет", "imageCls=imgmember", "disabled=true", "handler=function_member_enter_cabinet",
		))
		it.AddCommandItem(rule, fancy.OrdUnit+10, lik.BuildSet(
			"type=button", "tip=Создать запись сотрудника", "imageCls=imgadd", "handler=function_member_create_cabinet",
		))
		//it.AddCommandImg(rule, fancy.OrdUnit+10, "Создать запись сотрудника", "add")
	}
}

//	Вывод формы
func (it *MemberControl) MemberElmForm(rule *repo.DataRule, elm *likbase.ItElm) {
	it.ListFormElm(rule, elm)
}

//	Проверка строки
func (it *MemberControl) MemberListRowsProbe(rule *repo.DataRule, elm *likbase.ItElm) bool {
	accept := (elm != nil)
	if accept && !it.ListMakeProbe(rule, elm) {
		accept = false
	}
	return accept
}

//	Вывод строки
func (it *MemberControl) MemberRowFill(rule *repo.DataRule, elm *likbase.ItElm, row lik.Seter) bool {
	it.ListRowFill(rule, elm, row)
	if row == nil {
		return false
	}
	bonus := repo.CalculeBonus(elm.Id)
	row.SetItem(bonus, "bonus")
	//row.SetItem(fmt.Sprintf("/%s%d/membercard%d?_tp=1", it.Main, int(elm.Id), int(elm.Id)), "pathopen")
	return true
}

//	Вход в строку
func (it *MemberControl) MemberEnterRow(rule *repo.DataRule) {
	//rule.SetPagePush(fmt.Sprintf("membercard%d", int(it.IdMain)))
}

