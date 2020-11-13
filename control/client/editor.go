package client

import (
	"github.com/massarakhsh/tsan/control"
	"github.com/massarakhsh/tsan/fancy"
	"github.com/massarakhsh/tsan/jone"
	"github.com/massarakhsh/tsan/repo"
	"github.com/massarakhsh/tsan/show"
	"github.com/massarakhsh/lik"
	"github.com/massarakhsh/lik/likbase"
	"github.com/massarakhsh/lik/likdom"
)

//	Дескриптор окна редактора
type ClientEditor struct {
	control.DataControl
	fancy.DataFancy
}

//	Интерфейс команд
type dealEditorExecute struct {
	It	*ClientEditor
}
func (it *dealEditorExecute) Run(rule *repo.DataRule, cmd string, data lik.Seter) {
	it.It.EditorExecute(rule, cmd, data)
}

//	Интерфейс проверки поля
type dealEditorFieldProbe struct {
	It	*ClientEditor
}
func (it *dealEditorFieldProbe) Run(rule *repo.DataRule, field lik.Seter) bool {
	return it.It.EditorFieldProbe(rule, field)
}

//	Конструктор дескриптора
func BuildEditor(rule *repo.DataRule, main string, id lik.IDB) *ClientEditor {
	it := &ClientEditor{ }
	it.ControlInitializeZone(main, id, "editor")
	it.Sel = likbase.IDBToStr(it.IdMain)
	if id > 0 {
		it.Fun = fancy.FunShow
	} else {
		it.Fun = fancy.FunAdd
	}
	it.Form.Tab = 0
	it.FancyInitialize(main, "client", "editor")
	it.ItExecute = &dealEditorExecute{it}
	it.ItFieldProbe = &dealEditorFieldProbe{it}
	return it
}

//	Выполнение команд редактора
func (it *ClientEditor) EditorExecute(rule *repo.DataRule, cmd string, data lik.Seter) {
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
	} else if cmd == "loadfile" {
		it.editorLoadFile(rule)
	} else if cmd == "write" {
		it.FormFixTab(rule.Shift())
		it.cmdWrite(rule, data)
		it.Fun = fancy.FunShow
	} else if cmd == "cancel" {
		it.FormFixTab(rule.Shift())
		it.Fun = fancy.FunShow
	} else if cmd == "loadfile" {
		it.editorLoadFile(rule)
	} else if cmd == "redoid" {
		it.editorReDoId(rule)
	}
}

//	Отображение окна редактора
func (it *ClientEditor) BuildShow(rule *repo.DataRule, sx int, sy int) likdom.Domer {
	it.SetSize(sx, sy)
	return show.BuildFancyForm(it.Main,"editor")
}

//	Отображение формы
func (it *ClientEditor) cmdShowForm(rule *repo.DataRule) {
	elm := jone.GetElm(it.Part, it.IdMain)
	it.FormClear()
	it.RunFieldsFill(rule)
	it.Form.Title = "Карточка клиента"
	it.Form.Tabs.AddItems("Карточка")
	it.Form.Items.AddItemSet("type=tab", "items", it.editorFieldsCard(rule, elm))
	it.Form.Tabs.AddItems("Функции")
	it.Form.Items.AddItemSet("type=tab", "items", it.editorFieldsFunction(rule, elm))
	if !it.IsEdit() {
		it.Form.Tools.AddItemSet("text=Редактировать", "handler=function_fancy_edit_start")
	} else {
		it.Form.Tools.AddItemSet("text=Записать", "handler=function_fancy_edit_write")
		it.Form.Tools.AddItemSet("text=Отменить", "handler=function_fancy_edit_cancel")
	}
	it.Form.SetParameter(it.Form.Tab, "activeTab")
	it.Form.SetSize(it.Sx,0)
	it.ShowForm(rule)
}

//	Заполнение полей объекта
func (it *ClientEditor) editorFieldsCard(rule *repo.DataRule, elm *likbase.ItElm) lik.Lister {
	items := it.FormElmCollect(rule, elm,"","client_all")
	return items
}

//	Заполнение полей адреса
func (it *ClientEditor) editorFieldsFunction(rule *repo.DataRule, elm *likbase.ItElm) lik.Lister {
	items := it.FormElmCollect(rule, elm,"", "address")
	return items
}

//	Проверка поля
func (it *ClientEditor) EditorFieldProbe(rule *repo.DataRule, field lik.Seter) bool {
	if !it.FancyFieldProbe(rule, field) { return false }
	tags := field.GetInt("tags")
	target := ""
	realty := ""
	if elm := jone.GetElm(it.Part, it.IdMain); elm != nil {
		if target = jone.CalculateElmString(elm,"target"); target == "sale" {
			realty = jone.CalculateElmString(elm, "objectid/realty")
		} else if target == "buy" {
			realty = jone.CalculateElmString(elm, "require/realty")
		}
	}
	return it.FancyProbeTags(target, realty, tags)
}

//	Запись изменений
func (it *ClientEditor) cmdWrite(rule *repo.DataRule, data lik.Seter) {
	elm := jone.GetElm(it.Part, it.IdMain)
	/*if elm.Table.Part == "offer" && jone.CalculateElmIDB(elm, "objectid") == 0 {
		obj := jone.TableObject.CreateElm()
		jone.SetElmValue(elm, rule.ItSession.IdClient, "clientid")
		jone.SetElmValue(elm, obj.Id, "objectid")
	}*/
	it.UpdateElmData(rule, elm, data)
	rule.OnChangeData()
}

//	Загрузка файла
func (it *ClientEditor) editorLoadFile(rule *repo.DataRule) {
	bufs := rule.GetBuffers()
	_ = bufs
	rule.IsJson = false
	//rule.ItPage.NeedUrl = true
}

//	Переделка заявки под новый номер
func (it *ClientEditor) editorReDoId(rule *repo.DataRule) {
	if elm := jone.GetElm(it.Part, it.IdMain); elm != nil {
		offer := jone.TableOffer.CreateElm()
		offer.Wait()
		elm.SetValue(offer.Id, "idu")
		offer.Delete()
		rule.OnChangePage()
	}
}

