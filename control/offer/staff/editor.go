package staff

import (
	"bitbucket.org/961961/tsan/control"
	"bitbucket.org/961961/tsan/fancy"
	"bitbucket.org/961961/tsan/jone"
	"bitbucket.org/961961/tsan/repo"
	"bitbucket.org/961961/tsan/show"
	"bitbucket.org/shaman/lik"
	"bitbucket.org/shaman/lik/likbase"
	"bitbucket.org/shaman/lik/likdom"
	"fmt"
)

//	Дескриптор окна редактора
type ShowEditor struct {
	control.DataControl
	fancy.DataFancy
}

//	Интерфейс команд
type dealEditorExecute struct {
	It	*ShowEditor
}
func (it *dealEditorExecute) Run(rule *repo.DataRule, cmd string, data lik.Seter) {
	it.It.EditorExecute(rule, cmd, data)
}

//	Интерфейс проверки поля
type dealEditorFieldProbe struct {
	It	*ShowEditor
}
func (it *dealEditorFieldProbe) Run(rule *repo.DataRule, field lik.Seter) bool {
	return it.It.EditorFieldProbe(rule, field)
}

//	Конструктор дескриптора
func BuildEditor(rule *repo.DataRule, main string, id lik.IDB) *ShowEditor {
	it := &ShowEditor{ }
	it.ControlInitializeZone(main, id, "editor")
	it.Sel = likbase.IDBToStr(it.IdMain)
	it.Fun = fancy.FunShow
	it.Form.Tab = 0
	it.FancyInitialize(main, "offer", "editor")
	it.ItExecute = &dealEditorExecute{it}
	it.ItFieldProbe = &dealEditorFieldProbe{it}
	return it
}

//	Выполнение команд редактора
func (it *ShowEditor) EditorExecute(rule *repo.DataRule, cmd string, data lik.Seter) {
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
func (it *ShowEditor) BuildShow(rule *repo.DataRule, sx int, sy int) likdom.Domer {
	it.SetSize(sx, sy)
	return show.BuildFancyForm(it.Main,"editor")
}

//	Отображение формы
func (it *ShowEditor) cmdShowForm(rule *repo.DataRule) {
	elm := jone.GetElm(it.Part, it.IdMain)
	it.FormClear()
	it.RunFieldsFill(rule)
	it.Form.Title = "Объект недвижимости"
	if it.Part == "offer" {
		it.Form.Title = "Заявка"
		segment := jone.CalculateElmString(elm,"segment")
		target := jone.CalculateElmString(elm,"target")
		if segment == jone.DoRent && target == "sale" {
			it.Form.Title += " \"сдать\""
		} else if segment == jone.DoRent && target == "buy" {
			it.Form.Title += " \"снять\""
		} else if target == "sale" {
			it.Form.Title += " на продажу"
		} else if target == "buy" {
			it.Form.Title += " на покупку"
		}
		it.Form.Tabs.AddItems("Заявка")
		it.Form.Items.AddItemSet("type=tab", "items", it.editorFieldsOffer(rule, elm))
		if target == "sale" {
			it.Form.Tabs.AddItems("Объект")
			it.Form.Items.AddItemSet("type=tab", "items", it.editorFieldsObject(rule, elm))
			it.Form.Tabs.AddItems("Адрес")
			it.Form.Items.AddItemSet("type=tab", "items", it.editorFieldsAddress(rule, elm))
		} else if target == "buy" {
			it.Form.Tabs.AddItems("Требования")
			it.Form.Items.AddItemSet("type=tab", "items", it.editorFieldsRequire(rule, elm))
		}
		it.Form.Tabs.AddItems("Клиент")
		it.Form.Items.AddItemSet("type=tab", "items", it.editorFieldsClient(rule, elm))
		it.Form.Tabs.AddItems("Договор")
		it.Form.Items.AddItemSet("type=tab", "items", it.editorFieldsContract(rule, elm))
	}
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

//	Заполнение полей заявки
func (it *ShowEditor) editorFieldsOffer(rule *repo.DataRule, elm *likbase.ItElm) lik.Lister {
	items := it.FormElmCollect(rule, elm,"","object_offer")
	idu := elm.GetIDB("idu")
	if idu == 0 { idu = elm.Id }
	did := likdom.BuildDiv("width='100%'")
	did.BuildString("Номер заявки: ")
	did.BuildItem("B").BuildString(lik.IDBToStr(idu))
	if it.IsEdit() {
		did.AppendItem(show.LinkTextProc("cmd", "Изменить", fmt.Sprintf("change_id_new(%d)", int(idu))))
	}
	items.InsertItem(lik.BuildSet("type=html", "value", did.ToString()), 0)
	target := elm.GetString("target")
	post := lik.IfString(target == "sale", "objectid", "require")
	items.AddItemSet("type=textarea", "label=Описание", "emptyText=Введите описание объекта",
		"name=s_" + post + "__definition", "cls=definition",
		"value", jone.CalculateElmString(elm, post + "/definition"), "editable", it.IsEdit())
	return items
}

//	Заполнение полей объекта
func (it *ShowEditor) editorFieldsObject(rule *repo.DataRule, elm *likbase.ItElm) lik.Lister {
	items := it.FormElmCollect(rule, elm,"objectid/define","define")
	return items
}

//	Заполнение полей адреса
func (it *ShowEditor) editorFieldsAddress(rule *repo.DataRule, elm *likbase.ItElm) lik.Lister {
	items := it.FormElmCollect(rule, elm,"objectid/address", "address")
	return items
}

//	Заполнение полец требований
func (it *ShowEditor) editorFieldsRequire(rule *repo.DataRule, elm *likbase.ItElm) lik.Lister {
	items := it.FormElmCollect(rule, elm,"require", "require")
	return items
}

//	Заполнение полей клиента
func (it *ShowEditor) editorFieldsClient(rule *repo.DataRule, elm *likbase.ItElm) lik.Lister {
	items := it.FormElmCollect(rule, elm,"clientid","offer_client")
	return items
}

//	Заполнение полей контракта
func (it *ShowEditor) editorFieldsContract(rule *repo.DataRule, elm *likbase.ItElm) lik.Lister {
	items := it.FormElmCollect(rule, elm,"contract","contract")
	present := false
	if docs := jone.CalculateElmList(elm,"objectid/picture"); docs != nil {
		for nd := 0; nd < docs.Count(); nd++ {
			if doc := docs.GetSet(nd); doc != nil && doc.GetString("media") == "doc" {
				if !present {
					items.AddItemSet("type=html", "value", "<b>Список документов:</b>")
					present = true
				}
				url := doc.GetString("url")
				text := doc.GetString("comment")
				if text == "" { text = "(без названия)" }
				a := fmt.Sprintf("<a target=_blank href='%s'>%s</a>", url, text)
				items.AddItemSet("type=html", "value", a)
			}
		}
	}
	return items
}

//	Проверка поля
func (it *ShowEditor) EditorFieldProbe(rule *repo.DataRule, field lik.Seter) bool {
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
func (it *ShowEditor) cmdWrite(rule *repo.DataRule, data lik.Seter) {
	elm := jone.GetElm(it.Part, it.IdMain)
	if elm.Table.Part == "offer" && jone.CalculateElmIDB(elm, "objectid") == 0 {
		obj := jone.TableObject.CreateElm()
		jone.SetElmValue(elm, rule.ItSession.IdMember, "memberid")
		jone.SetElmValue(elm, obj.Id, "objectid")
	}
	it.UpdateElmData(rule, elm, data)
	rule.OnChangeData()
}

//	Загрузка файла
func (it *ShowEditor) editorLoadFile(rule *repo.DataRule) {
	bufs := rule.GetBuffers()
	_ = bufs
	rule.IsJson = false
	//rule.ItPage.NeedUrl = true
}

//	Переделка заявки под новый номер
func (it *ShowEditor) editorReDoId(rule *repo.DataRule) {
	if elm := jone.GetElm(it.Part, it.IdMain); elm != nil {
		offer := jone.TableOffer.CreateElm()
		offer.Wait()
		elm.SetValue(offer.Id, "idu")
		offer.Delete()
		rule.OnChangePage()
	}
}

