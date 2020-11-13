package controls

import (
	"bitbucket.org/961961/tsan/control"
	"bitbucket.org/961961/tsan/fancy"
	"bitbucket.org/961961/tsan/jone"
	"bitbucket.org/961961/tsan/repo"
	"bitbucket.org/961961/tsan/show"
	"bitbucket.org/shaman/lik"
	"bitbucket.org/shaman/lik/likbase"
	"bitbucket.org/shaman/lik/likdom"
)

//	Дескриптор выбора режима
type CommandControl struct {
	control.DataControl
	fancy.DataFancy
}

//	Обработчик событий
type dealRoleExecute struct {
	It	*CommandControl
}
func (it *dealRoleExecute) Run(rule *repo.DataRule, cmd string, data lik.Seter) {
	it.It.RoleExecute(rule, cmd, data)
}

//	Конструктор дескриптора
func BuildCommand(rule *repo.DataRule, id lik.IDB) *CommandControl {
	it := &CommandControl{}
	it.ControlInitialize("command", id)
	it.FancyInitialize("command", "member","all")
	it.Sel = likbase.IDBToStr(it.IdMain)
	it.Fun = fancy.FunShow
	it.ItExecute = &dealRoleExecute{it}
	return it
}

//	Обработка событий
func (it *CommandControl) RoleExecute(rule *repo.DataRule, cmd string, data lik.Seter) {
	if cmd == "all" {
		it.RoleExecute(rule, rule.Shift(), data)
	} else if cmd == "role" || cmd == "segment" {
		it.Zone = cmd
		it.RoleExecute(rule, rule.Shift(), data)
	} else if cmd == "showform" {
		it.cmdShowForm(rule)
	} else if cmd == "logoff" {
		it.cmdLogoff(rule)
	} else if cmd == jone.ItDispatch && rule.ICanRole(jone.ItDispatch) {
		it.cmdToRole(rule, cmd)
	} else if cmd == jone.ItRealtor && rule.ICanRole(jone.ItRealtor) {
		it.cmdToRole(rule, cmd)
	} else if cmd == jone.ItManager && rule.ICanRole(jone.ItManager) {
		it.cmdToRole(rule, cmd)
	} else if cmd == jone.ItAdvert && rule.ICanRole(jone.ItAdvert) {
		it.cmdToRole(rule, cmd)
	} else if cmd == jone.ItAdmin && rule.ICanRole(jone.ItAdmin) {
		it.cmdToRole(rule, cmd)
	} else if it.Zone == "segment" {
		it.cmdToSegment(rule, cmd)
	}
}

//	Отображение окна
func (it *CommandControl) BuildShow(rule *repo.DataRule, sx int, sy int) likdom.Domer {
	it.SetSize(sx, sy)
	return show.BuildFancyForm(it.Main,"editor")
}

//	Отображение формы
func (it *CommandControl) cmdShowForm(rule *repo.DataRule) {
	it.FormClear()
	it.RunFieldsFill(rule)
	it.Form.Tools.AddItemSet("text=Отмена", "handler=function_fancy_form_cancel")
	if it.Zone == "role" {
		it.Form.Title = "Выбор режима"
		if rule.ICanRole(jone.ItDispatch) {
			it.Form.Items.AddItemSet("type=html", "value",
				it.showCmdWhat(rule, jone.ItDispatch, "Как диспетчер"))
		}
		if rule.ICanRole(jone.ItRealtor) {
			it.Form.Items.AddItemSet("type=html", "value",
				it.showCmdWhat(rule, jone.ItRealtor, "Как риэлтор"))
		}
		if rule.ICanRole(jone.ItManager) {
			it.Form.Items.AddItemSet("type=html", "value",
				it.showCmdWhat(rule, jone.ItManager, "Как менеджер"))
		}
		if rule.ICanRole(jone.ItAdvert) {
			it.Form.Items.AddItemSet("type=html", "value",
				it.showCmdWhat(rule, jone.ItAdvert, "Как рекламный менеджер"))
		}
		if rule.ICanRole(jone.ItAdmin) {
			it.Form.Items.AddItemSet("type=html", "value",
				it.showCmdWhat(rule, jone.ItAdmin, "Как администратор"))
		}
		it.Form.SetParameter(200, "likLeft")
		it.Form.SetParameter(50, "likTop")
		it.Form.SetSize(350,250)
	} else if it.Zone == "segment" {
		it.Form.Title = "Выбор сегмента"
		if true {
			it.Form.Items.AddItemSet("type=html", "value",
				it.showCmdWhat(rule, jone.DoCall, "Колл-центр"))
		}
		if rule.ICanDo(jone.DoSecond) || rule.IAmAdmin() {
			it.Form.Items.AddItemSet("type=html", "value",
				it.showCmdWhat(rule, jone.DoSecond, "Вторичка"))
		}
		if rule.ICanDo(jone.DoNew) || rule.IAmAdmin() {
			it.Form.Items.AddItemSet("type=html", "value",
				"<b>Новостройки</b>") //it.showCmdWhat(rule, jone.DoNew, "Новостройки"))
		}
		if rule.ICanDo(jone.DoRent) || rule.IAmAdmin() {
			it.Form.Items.AddItemSet("type=html", "value",
				it.showCmdWhat(rule, jone.DoRent, "Аренда"))
		}
		if rule.ICanDo(jone.DoVilla) || rule.IAmAdmin() {
			it.Form.Items.AddItemSet("type=html", "value",
				it.showCmdWhat(rule, jone.DoVilla, "Загород"))
		}
		if rule.ICanDo(jone.DoArea) || rule.IAmAdmin() {
			it.Form.Items.AddItemSet("type=html", "value",
				it.showCmdWhat(rule, jone.DoArea, "Участки"))
		}
		if rule.IAmAdmin() {
			it.Form.Items.AddItemSet("type=html", "value", "<hr>")
			it.Form.Items.AddItemSet("type=html", "value",
				it.showCmdWhat(rule, jone.DoTune, "Настройки"))
		}
		it.Form.SetParameter(100, "likLeft")
		it.Form.SetParameter(50, "likTop")
		it.Form.SetSize(250,250)
	}
	it.Form.Items.AddItemSet("type=html", "value", "<hr>")
	it.Form.Items.AddItemSet("type=html", "value",
		it.showCmdWhat(rule, "logoff", "Выйти из системы"))
	it.Form.Items.AddItemSet("type=html", "value", "<br>")
	it.ShowForm(rule)
}

//	Обработка команды
func (it *CommandControl) showCmdWhat(rule *repo.DataRule, code string, text string) string {
	return show.LinkTextProc("cmd", text, "choose_command('" + code + "')").ToString()
}

//	Выход
func (it *CommandControl) cmdLogoff(rule *repo.DataRule) {
	rule.SessionLogin(0, 0, true)
	rule.ItSession.ExitLogin = true
	rule.OnChangeData()
}

//	Выбор роли
func (it *CommandControl) cmdToRole(rule *repo.DataRule, role string) {
	if rule.ICanRole(role) {
		jone.SetElmValue(rule.GetMember(), role, "rolit")
		rule.ItSession.RoleMember = role
		rule.OnChangePage()
	}
}

//	Выбор сегмента
func (it *CommandControl) cmdToSegment(rule *repo.DataRule, segment string) {
	if rule.ICanDo(segment) || rule.IAmAdmin() {
		rule.ItPage.Params.SetItem(segment, "segment")
		rule.OnChangePage()
	}
}

