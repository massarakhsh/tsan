package member

import (
	"github.com/massarakhsh/tsan/control"
	"github.com/massarakhsh/tsan/fancy"
	"github.com/massarakhsh/tsan/jone"
	"github.com/massarakhsh/tsan/one"
	"github.com/massarakhsh/tsan/repo"
	"github.com/massarakhsh/tsan/show"
	"github.com/massarakhsh/lik"
	"github.com/massarakhsh/lik/likbase"
	"github.com/massarakhsh/lik/likdom"
)

//	Дескриптор окна редактора
type MemberEditor struct {
	control.DataControl
	fancy.DataFancy
}

//	Интерфейс команд
type dealEditorExecute struct {
	It	*MemberEditor
}
func (it *dealEditorExecute) Run(rule *repo.DataRule, cmd string, data lik.Seter) {
	it.It.EditorExecute(rule, cmd, data)
}

//	Конструктор дескриптора
func BuildEditor(rule *repo.DataRule, main string, id lik.IDB) *MemberEditor {
	it := &MemberEditor{ }
	it.ControlInitializeZone(main, id, "editor")
	it.Sel = likbase.IDBToStr(it.IdMain)
	if id > 0 {
		it.Fun = fancy.FunShow
	} else {
		it.Fun = fancy.FunAdd
	}
	it.Form.Tab = 0
	it.FancyInitialize(main, "member", "editor")
	it.ItExecute = &dealEditorExecute{it}
	return it
}

//	Выполнение команд редактора
func (it *MemberEditor) EditorExecute(rule *repo.DataRule, cmd string, data lik.Seter) {
	if cmd == "showform" {
		it.cmdShowForm(rule)
	} else if cmd == "edit" {
		it.Fun = fancy.FunMod
		it.FormFixTab(rule.Shift())
		rule.OnChangeData()
	} else if cmd == "cancel" {
		it.Fun = fancy.FunShow
		it.FormFixTab(rule.Shift())
		rule.OnChangeData()
	} else if cmd == "write" {
		it.FormFixTab(rule.Shift())
		it.cmdWrite(rule, data)
		it.Fun = fancy.FunShow
	}
}

//	Отображение окна редактора
func (it *MemberEditor) BuildShow(rule *repo.DataRule, sx int, sy int) likdom.Domer {
	it.SetSize(sx, sy)
	return show.BuildFancyForm(it.Main,"editor")
}

//	Отображение формы
func (it *MemberEditor) cmdShowForm(rule *repo.DataRule) {
	elm := jone.GetElm(it.Part, it.IdMain)
	it.FormClear()
	it.RunFieldsFill(rule)
	it.Form.Title = "Карточка сотрудника"
	it.Form.Tabs.AddItems("Карточка")
	it.Form.Items.AddItemSet("type=tab", "items", it.editorFieldsCard(rule, elm))
	it.Form.Tabs.AddItems("Функции")
	it.Form.Items.AddItemSet("type=tab", "items", it.editorFieldsFunction(rule, elm))
	if it.IsEdit() {
		it.Form.Tools.AddItemSet("text=Записать", "handler=function_fancy_edit_write")
		it.Form.Tools.AddItemSet("text=Отменить", "handler=function_fancy_edit_cancel")
	} else if rule.IAmAdmin() {
		it.Form.Tools.AddItemSet("text=Редактировать", "handler=function_fancy_edit_start")
	}
	it.Form.SetParameter(it.Form.Tab, "activeTab")
	it.Form.SetSize(it.Sx,0)
	it.ShowForm(rule)
}

//	Заполнение полей объекта
func (it *MemberEditor) editorFieldsCard(rule *repo.DataRule, elm *likbase.ItElm) lik.Lister {
	items := it.FormElmCollect(rule, elm,"","member_all")
	return items
}

//	Заполнение полей полномочий
func (it *MemberEditor) editorFieldsFunction(rule *repo.DataRule, elm *likbase.ItElm) lik.Lister {
	items := lik.BuildList()
	for _,name := range []string {jone.DoSecond, jone.DoRent, jone.DoNew, jone.DoVilla, jone.DoArea} {
		item := lik.BuildSet("type=checkbox", "name", "h_do_" + name, "editable")
		value := elm != nil && elm.GetBool("do_" + name)
		if name == jone.DoSecond {
			item.SetItem("Вторичка", "label")
		} else if name == jone.DoRent {
			item.SetItem("Аренда", "label")
		} else if name == jone.DoNew {
			item.SetItem("Новостройки", "label")
		} else if name == jone.DoVilla {
			item.SetItem("Загород", "label")
		} else if name == jone.DoArea {
			item.SetItem("Участки", "label")
		}
		item.SetItem(value, "value")
		if !it.IsEdit() || !rule.IAmAdmin() || name == jone.DoNew {
			item.SetItem(false, "editable")
			item.SetItem("readonly", "cls")
		}
		items.AddItems(item)
	}
	return items
}

//	Запись изменений
func (it *MemberEditor) cmdWrite(rule *repo.DataRule, data lik.Seter) {
	elm := jone.GetElm(it.Part, it.IdMain)
	if elm == nil {
		elm = jone.GetTable(it.Part).CreateElm()
	}
	if elm != nil {
		it.IdMain = elm.Id
		it.Sel = likbase.IDBToStr(elm.Id)
		it.UpdateElmData(rule, elm, data)
		if it,ok := one.GetMember(elm.Id); ok {
			repo.SynchronizeMemberOne(&it, elm)
		}
	}
	rule.OnChangeData()
}

