//	Контроллер списка контактов.
//
//	Модуль реализует интерфейс контроллера списка контактов (звонков)
//
//	Одна из задач модуля - мониторинг таблицы звонков АТС и автоматическая инициализация контактов
package bell

import (
	"bitbucket.org/961961/tsan/control"
	"bitbucket.org/961961/tsan/fancy"
	"bitbucket.org/961961/tsan/jone"
	"bitbucket.org/961961/tsan/repo"
	"bitbucket.org/961961/tsan/routine/taskcall"
	"bitbucket.org/961961/tsan/show"
	"bitbucket.org/shaman/lik"
	"bitbucket.org/shaman/lik/likbase"
	"bitbucket.org/shaman/lik/likdom"
	"fmt"
	"strings"
	"time"
)

//	Декскриптор контроллера списка контактов
type BellControl struct {
	control.ListControl         //	Основан на контроллере списка
	cCall               int     //	Число звонков в очереди
	cOpen               int     //	Число активных звонков
	cMe                 lik.IDB //	Активный контакт текущего контроллера
	onCalls             map[lik.IDB]bool
	IsBellChange        bool
	lastBellForm        time.Time
}

//	Обработчик событий команд
type dealBellExecute struct {
	It	*BellControl
}
func (it *dealBellExecute) Run(rule *repo.DataRule, cmd string, data lik.Seter) {
	it.It.BellExecute(rule, cmd, data)
}

//	Обработчик дежурного опроса
type dealBellMarshal struct {
	It	*BellControl
}
func (it *dealBellMarshal) Run(rule *repo.DataRule) {
	it.It.BellMarshal(rule)
}

//	Обработчик инициализации таблицы
type dealBellGridFill struct {
	It	*BellControl
}
func (it *dealBellGridFill) Run(rule *repo.DataRule) {
	it.It.BellGridFill(rule)
}

//	Обработчик заполнения строки
type dealBellRowFill struct {
	It	*BellControl
}
func (it *dealBellRowFill) Run(rule *repo.DataRule, elm *likbase.ItElm, row lik.Seter) {
	it.It.BellRowFill(rule, elm, row)
}

//	Обработчик проверки строки
type dealBellMakeProbe struct {
	It	*BellControl
}
func (it *dealBellMakeProbe) Run(rule *repo.DataRule, elm *likbase.ItElm) bool {
	return it.It.BellMakeProbe(rule, elm)
}

//	Обработчик заполнения формы
type dealBellElmForm struct {
	It	*BellControl
}
func (it *dealBellElmForm) Run(rule *repo.DataRule, elm *likbase.ItElm) {
	it.It.BellElmForm(rule, elm)
}

//	Обрабюотчик импорта элементов
type dealBellImport struct {
	It	*BellControl
}
func (it *dealBellImport) Run(rule *repo.DataRule, part string, sid string) {
	it.It.BellImportContact(rule, part, sid)
}

//	Конструктор дескриптора списка контактов
func BuildListBell(rule *repo.DataRule, id lik.IDB) *BellControl {
	it := &BellControl{}
	it.onCalls = make(map[lik.IDB]bool)
	it.ControlInitialize("bell", id)
	it.ListInitialize(rule, "bell", "bell")
	it.ItExecute = &dealBellExecute{it}
	it.ItMarshal = &dealBellMarshal{it}
	it.ItGridFill = &dealBellGridFill{it}
	it.ItRowFill = &dealBellRowFill{it}
	it.ItListMakeProbe = &dealBellMakeProbe{it}
	it.ItFormElm = &dealBellElmForm{it}
	it.ItImport = &dealBellImport{it}
	return it
}

//	Обработка события
func (it *BellControl) BellExecute(rule *repo.DataRule, cmd string, data lik.Seter) {

	if cmd == "all" || cmd == "bell" {
		it.BellExecute(rule, rule.Shift(), data)
	} else if cmd == "gooffer" {
		it.bellGoOffer(rule, likbase.StrToIDB(rule.Shift()))
	} else if cmd == "offercreate" {
		it.cmdOfferCreate(rule, data)
	} else if cmd == "bellaccept" {
		it.cmdBellCancel(rule, "", data)
	} else if cmd == "bellcancel" {
		it.cmdBellCancel(rule, rule.Shift(), data)
	} else if cmd == "phonesearch" {
		it.cmdPhoneSearch(rule, rule.Shift())
	} else if cmd == "offersearch" {
		it.cmdOfferSearch(rule, rule.Shift())
	} else if cmd == "cancel" {
		it.FormFixTab(rule.Shift())
		rule.SetResponse("bell_all_show", "_function_fancy_trio")
	} else if cmd == "write" {
		it.FormFixTab(rule.Shift())
		it.BellCmdWrite(rule, data)
	} else if cmd == "setpin" {
		pin := lik.StrToInt(lik.StringFromXS(rule.Shift()))
		rule.SetPinMember(pin)
		rule.OnChangeData()
	} else if cmd == "beta" {
		rule.OnChangeData()
	} else if cmd == "bellform" {
		it.onCalls[lik.StrToIDB(rule.Shift())] = true
	} else if cmd == "toenter" {
		it.BellEnterRow(rule)
	} else {
		it.ListExecute(rule, cmd, data)
	}
}

//	Обработка маршализации
func (it *BellControl) BellMarshal(rule *repo.DataRule) {
	it.bellMonitorCall(rule)
	it.ListMarshal(rule)
}

//	Заполнение таблицы
func (it *BellControl) BellGridFill(rule *repo.DataRule) {
	it.ListGridFill(rule)
	it.AddCommandImg(rule, fancy.OrdUnit+ 1, "Создать контакт", "toadd", "add")
	it.bellGridFillPhone(rule)
}

//	Заполнение индикатора звонков
func (it *BellControl) bellGridFillPhone(rule *repo.DataRule) {
	it.AddCommandItem(rule, fancy.OrdPhone-1, lik.BuildSet())
	it.AddCommandItem(rule, fancy.OrdPhone, lik.BuildSet("type=text", "text=Телефон:", "cls=small"))
	it.bellProbeCall(rule)
	it.AddCommandItem(rule, fancy.OrdPhone+1, lik.BuildSet("type=text", "text", it.bellBuildCall(rule), "cls=topcmd"))
}

//	Мониторинг звонков
func (it *BellControl) bellMonitorCall(rule *repo.DataRule) {
	it.bellProbeCall(rule)
	if idcall := it.cMe; idcall != 0 && time.Now().Sub(it.lastBellForm) > time.Second * 5 {
		if _,use := it.onCalls[idcall]; !use {
			it.lastBellForm = time.Now()
			//it.onCalls[idcall] = true
			if elm := jone.TableBell.GetElm(idcall); elm != nil {
				if rid := elm.GetIDB("receptorid"); rid == 0 {
					elm.SetValue(rule.ItSession.IdMember, "receptorid")
					elm.SetValue(int(time.Now().Unix()), "date")
				}
				it.CmdFormDo = fmt.Sprintf("edit_%d", int(idcall))
			}
		}
	}
	if it.IsBellChange {
		it.IsBellChange = false
		serial := it.bellBuildCall(rule)
		rule.SetResponse(serial, "listcall")
		it.CmdGridDo = "refresh"
	}
}

//	Проверка звонков
func (it *BellControl) bellProbeCall(rule *repo.DataRule) {
	if member := rule.GetMember(); member != nil {
		pin := member.GetInt("pin")
		inc, oth, idb := taskcall.RequestCall(pin)
		if inc != it.cCall || oth != it.cOpen || idb != it.cMe {
			it.cCall, it.cOpen, it.cMe = inc, oth, idb
			it.IsBellChange = true
		}
	}
}

//	Вывод индикатора звонков
func (it *BellControl) bellBuildCall(rule *repo.DataRule) string {
	div := likdom.BuildDivClass("top", "id=listcall")
	pint := "нет"
	cls := "toc tocin"
	pin := rule.GetMember().GetInt("pin")
	if pin > 0 {
		pint = lik.IntToStr(pin)
		if it.cMe > 0 {
			cls += " tocact"
		}
	}
	cmd := likdom.BuildItem("a", "href=#", "onclick", fmt.Sprintf("change_pin(%d)",pin))
	cmd.BuildString(pint)
	div.BuildDivClass(cls).AppendItem(cmd)
	if it.cOpen == 0 && it.cCall == 0 {
		div.BuildDivClass("toc").BuildUnpairItem("img", "src=/images/pho_no.png", "title=Нет активных звонков")
	}
	if it.cOpen > 0 {
		div.BuildDivClass("toc").BuildString("&nbsp;")
		if it.cOpen == 1 {
			div.BuildDivClass("toc").BuildUnpairItem("img", "src=/images/pho_oth.png", "title=Разговор")
		} else {
			div.BuildDivClass("toc").BuildUnpairItem("img", "src=/images/pho_oth.png", "title=Разговоры")
			div.BuildDivClass("toc").BuildString(fmt.Sprintf("<small>(%d)</small>", it.cOpen))
		}
	}
	if it.cCall > 0 {
		div.BuildDivClass("toc").BuildString(fmt.Sprintf("+"))
		div.BuildDivClass("toc").BuildUnpairItem("img", "src=/images/pho_inc.png", "title=Звонки в очереди")
		if it.cCall > 1 {
			div.BuildDivClass("toc").BuildString(fmt.Sprintf("<small>(%d)</small>", it.cCall))
		}
	}
	return div.ToString()
}

//	Проверка строки для списка
func (it *BellControl) BellMakeProbe(rule *repo.DataRule, elm *likbase.ItElm) bool {
	if rule.IAmRealtor() {
		if it.Selector.ItLocate != jone.ItMy {
			it.Selector.ItLocate = jone.ItMy
		}
	} else if rule.IAmManager() {
		if it.Selector.ItLocate == "" || it.Selector.ItLocate == jone.ItAll {
			it.Selector.ItLocate = jone.ItDep
		}
	}
	accept := false
	if rid := jone.CalculateElmIDB(elm, "receptorid"); repo.ProbeItMy(rule, "member", rid) {
		accept = true
	} else if mid := jone.CalculateElmIDB(elm, "memberid"); repo.ProbeItMy(rule, "member", mid) {
		accept = true
	} else if it.Selector.ItLocate == jone.ItMy {
		accept = false
	} else if it.Selector.ItLocate == jone.ItDep {
		if rid != 0 && repo.ProbeItDep(rule, "member", rid) {
			accept = true
		} else if mid != 0 && repo.ProbeItDep(rule, "member", mid) {
			accept = true
		} else {
			accept = false
		}
	} else {
		accept = true
	}
	if accept && it.Selector.ItStatus == "active" {
		if offer := jone.TableOffer.GetElm(jone.CalculateElmIDB(elm,"offerid")); offer != nil {
			accept = false
		} else if jone.CalculateElmString(elm,"cancel") != "" {
			accept = false
		}
	}
	if accept {
		accept = it.ListMakeProbe(rule, elm)
	}
	return accept
}

//	Заполнение строки
func (it *BellControl) BellRowFill(rule *repo.DataRule, elm *likbase.ItElm, row lik.Seter) {
	it.ListRowFill(rule, elm, row)
	var pic likdom.Domer
	status := jone.CalculateElmString(elm,"status")
	offerid := jone.CalculateElmIDB(elm,"offerid")
	if offer := jone.TableOffer.GetElm(offerid); offer != nil {
		text := jone.CalculateElmSid(offer)
		if status != jone.ItReady {
			status = jone.ItReady
			elm.SetValue(jone.ItReady, "status")
		}
		//mode := repo.CalculateElmString(elm, "target")
		path := fmt.Sprintf("/offershow%d?_tp=1", int(offerid))
		//rule.SetResponse(path,"_function_lik_window_part")
		status = show.LinkTextProc("cmd", text, fmt.Sprintf("lik_window_part('%s')", path)).ToString()
	} else if cancel := jone.CalculateElmString(elm,"cancel"); cancel != "" {
		offerid = 0
		pic = likdom.BuildUnpairItem("img", "src","/images/cancel.png", "width=16px", "height=16px")
		text := jone.SystemStringTranslate("bell_cancel", cancel)
		pic.SetAttr("title", text)
		if status != jone.ItCancel {
			status = jone.ItCancel
			elm.SetValue(jone.ItCancel, "status")
		}
		status = "отмена"
	} else {
		offerid = 0
		atstart := jone.CalculateElmInt(elm,"date")
		atnow := int(time.Now().Unix())
		txdel := it.bellShowDelta(atstart, atnow)
		if clm := it.SeekPartField("delay"); clm != nil {
			row.SetItem(lik.BuildItem(txdel), clm.GetString("index"))
		}
		if len(it.bellTestDiagnosis(rule, elm.Info)) == 0 {
			pic = likdom.BuildUnpairItem("img", "src","/images/watch.png", "width=16px", "height=16px")
			pic.SetAttr("title", "Ожидает обработки")
			if status != jone.ItReady {
				status = jone.ItReady
				elm.SetValue(jone.ItReady, "status")
			}
			status = "ожидает"
		} else {
			pic = likdom.BuildUnpairItem("img", "src","/images/control.png", "width=16px", "height=16px")
			pic.SetAttr("title", "Не заполнена")
			if status != jone.ItError {
				status = jone.ItError
				elm.SetValue(jone.ItError, "status")
			}
			status = "ошибки"
		}
	}
	if elm.Id == it.cMe {
		pic = likdom.BuildUnpairItem("img", "src","/images/pho_act.png", "width=16px", "height=16px")
		pic.SetAttr("title", "Текущий звонок")
		if offerid == 0 {
			status = "..."
		}
	} else if stcall := taskcall.CallerProbe(elm.Id); stcall == "connected" {
		pic = likdom.BuildUnpairItem("img", "src","/images/pho_oth.png", "width=16px", "height=16px")
		pic.SetAttr("title", "Идет разговор")
		if offerid == 0 {
			status = "..."
		}
	} else if stcall == "calling" {
		pic = likdom.BuildUnpairItem("img", "src","/images/pho_inc.png", "width=16px", "height=16px")
		pic.SetAttr("title", "Звонок в очереди")
		if offerid == 0 {
			status = "..."
		}
	}
	if pic != nil {
		row.SetItem(pic.ToString(), "pic")
	}
	if status != "" {
		row.SetItem(status, "status")
	}
}

//	Визуализация длительности периода
func (it *BellControl) bellShowDelta(from int, to int) string {
	text := ""
	if to == 0 || from == 0 {
		text = "?"
	} else if to - from < 60 {
		text = fmt.Sprintf("%d сек", to - from)
	} else if to - from < 60*60 {
		text = fmt.Sprintf("%d мин",(to - from) / 60)
	} else if to - from < 60*60*24 {
		text = fmt.Sprintf("%d час", (to - from)/ 3600)
	} else {
		text = fmt.Sprintf("%d дн.", to / (3600*24) - from / (3600*24))
	}
	return text
}

//	Поиск телефона
func (it *BellControl) cmdPhoneSearch(rule *repo.DataRule, phone string) {
	ok := 0
	if phone != "" {
		if elm := repo.SearchClient(phone); elm != nil {
			ok = 1
			rule.SetResponse(elm.Id, "id")
			rule.SetResponse(jone.CalculateElmString(elm,"namely"), "namely")
			rule.SetResponse(jone.CalculateElmString(elm,"paterly"), "paterly")
			rule.SetResponse(jone.CalculateElmString(elm,"family"), "family")
		} else if elm := repo.SearchBell(phone); elm != nil {
			ok = 1
			rule.SetResponse(jone.CalculateElmString(elm,"clientid/namely"), "namely")
			rule.SetResponse(jone.CalculateElmString(elm,"clientid/paterly"), "paterly")
			rule.SetResponse(jone.CalculateElmString(elm,"clientid/family"), "family")
		}
	}
	rule.SetResponse(ok, "ok")
}

//	Поиск заявки
func (it *BellControl) cmdOfferSearch(rule *repo.DataRule, sid string) {
	if elm := jone.TableOffer.GetElm(likbase.StrToIDB(sid)); elm != nil {
		rule.SetResponse(elm.Id, "id")
		comp := ""
		if target := jone.CalculateElmString(elm,"target"); target == "sale" {
			comp = "buy"
		} else if target == "buy" {
			comp = "sale"
		}
		rule.SetResponse(comp, "target")
		rule.SetResponse(jone.CalculateElmString(elm,"segment"), "segment")
		rule.SetResponse(jone.CalculateElmString(elm,"objectid/realty"), "realty")
		rule.SetResponse(jone.CalculateElmString(elm,"objectid/address/subcity"), "subcity")
		rule.SetResponse(jone.CalculateElmString(elm,"objectid/define/rooms"), "rooms")
		rule.SetResponse(jone.CalculateElmString(elm,"memberid"), "memberid")
	} else {
		rule.SetResponse(0, "id")
	}
}

//	Отмена редактирования контакта
func (it *BellControl) cmdBellCancel(rule *repo.DataRule, parm string, data lik.Seter) {
	if elm := jone.GetElm(it.Part, likbase.StrToIDB(it.Sel)); elm != nil {
		if data != nil {
			it.BellCmdWrite(rule, data)
			it.CmdFormDo = ""
		}
		if parm != "" {
			jone.SetElmValue(elm, parm, "cancel")
			jone.SetElmValue(elm, int(time.Now().Unix()), "datedone")
		} else {
			jone.SetElmValue(elm,nil, "cancel")
			jone.SetElmValue(elm, nil, "datedone")
		}
		it.Form.Tab = 3
		rule.OnChangeData()
	}
}

//	Отображение окна редактирования
func (it *BellControl) BellElmForm(rule *repo.DataRule, elm *likbase.ItElm) {
	var info lik.Seter
	if it.Fun == fancy.FunAdd {
		info = lik.BuildSet()
	} else {
		info = elm.Info
	}
	it.Form.SetSize(480, 0)
	if it.Fun == fancy.FunAdd {
		it.bellFillAdd(rule, info)
	} else if it.Fun == fancy.FunShow {
		it.bellFillShow(rule, info)
	} else if it.Fun == fancy.FunMod || it.Fun == fancy.FunEdit {
		it.bellFillMod(rule, info)
	} else if it.Fun == fancy.FunDel {
		it.bellFillDel(rule, info)
	}
}

//	Окно в режиме отображения
func (it *BellControl) bellFillShow(rule *repo.DataRule, info lik.Seter) {
	it.SetTitle(rule, it.Fun, "Контакт с клиентом")
	idmem := jone.CalculateIDB(info, "memberid")
	idrec := jone.CalculateIDB(info, "receptorid")
	if rule.IAmAdmin() ||
		repo.ProbeItMy(rule, "member", idmem) ||
		repo.ProbeItMy(rule, "member", idrec) ||
		rule.IAmManager() && repo.ProbeItDep(rule, "member", idmem)  ||
		rule.IAmManager() && repo.ProbeItDep(rule, "member", idrec) {
		it.AddTitleToolText(rule, "Изменить", "function_fancy_form_toedit")
	}
	if rule.IAmAdmin() {
		it.AddTitleToolText(rule, "Удалить", "function_fancy_form_todelete")
	}
	it.AddTitleToolText(rule, "Закрыть", "function_fancy_form_cancel")
	it.bellFillTabs(rule, info)
}

//	Окно в режиме редактирования
func (it *BellControl) bellFillMod(rule *repo.DataRule, info lik.Seter) {
	if target := info.GetString("target"); target == "sale" {
		it.Form.SingleCho = true
	}
	it.SetTitle(rule, it.Fun, fmt.Sprintf("Редактирование контакта <b id=belledit>%s</b>", it.Sel))
	it.AddTitleToolText(rule, "Сохранить", "function_fancy_form_write")
	if it.Fun == fancy.FunMod {
		it.AddTitleToolText(rule, "Отменить", "function_fancy_form_toshow")
	} else {
		it.AddTitleToolText(rule, "Отменить", "function_fancy_form_cancel")
	}
	it.Form.AddEventAction("set", "function_bell_set_change")
	it.bellFillTabs(rule, info)
}

//	Окно в режиме добавления
func (it *BellControl) bellFillAdd(rule *repo.DataRule, info lik.Seter) {
	it.SetTitle(rule, it.Fun, "Создание контакта")
	it.AddTitleToolText(rule, "Сохранить", "function_fancy_form_write")
	it.AddTitleToolText(rule, "Отменить", "function_fancy_form_cancel")
	jone.SetInfoValue(info, int(time.Now().Unix()), "date")
	jone.SetInfoValue(info, rule.ItSession.IdMember, "receptorid")
	jone.SetInfoValue(info, rule.ItSession.IdMember, "memberid")
	it.Form.AddEventAction("set", "function_bell_set_change")
	it.Form.Tab = 0
	it.bellFillTabs(rule, info)
}

//	Окно в режиме удаления
func (it *BellControl) bellFillDel(rule *repo.DataRule, info lik.Seter) {
	it.SetTitle(rule, it.Fun,fmt.Sprintf("Удаление контакта"))
	it.AddTitleToolText(rule, "Действительно удалить?", "function_fancy_real_delete")
	it.AddTitleToolText(rule, "Отменить", "function_fancy_form_toshow")
	it.bellFillTabs(rule, info)
}

//	Заполнение закладок
func (it *BellControl) bellFillTabs(rule *repo.DataRule, info lik.Seter) {
	it.Form.Tabs.AddItems("Контакт")
	it.Form.Items.AddItemSet("type=tab", "items", it.bellFieldsMain(rule, info))
	it.Form.Tabs.AddItems("Объект")
	it.Form.Items.AddItemSet("type=tab", "items", it.bellFieldsObject(rule, info))
	it.Form.Tabs.AddItems("Заявка")
	it.Form.Items.AddItemSet("type=tab", "items", it.collectOfferFields(rule, info))
}

//	Закладка основных сведений
func (it *BellControl) bellFieldsMain(rule *repo.DataRule, info lik.Seter) lik.Lister {
	items := it.FormInfoCollect(rule, info, "", "bell_main")
	cid := jone.CalculateIDB(info, "clientid")
	items.InsertItem(lik.BuildSet("type=hidden", "name=u_clientid", "value", cid), 0)
	rid := jone.CalculateIDB(info, "receptorid")
	oper := jone.CalculatePartIdText("member", rid)
	if rid == 0 && (it.Fun == fancy.FunAdd || it.Fun == fancy.FunMod || it.Fun == fancy.FunEdit) {
		oper = jone.CalculateElmText(rule.GetMember()) + "?"
	}
	items.InsertItem(lik.BuildSet("type=string", "label=Кто принял", "name=u__receptor",
		"editable=false", "cls=readonly", "value", oper), 2)
	if pin := jone.CalculateString(info, "pin"); pin != "" {
		items.InsertItem(lik.BuildSet("type=string", "label=Пин", "name=u_pin",
			"editable=false", "cls=readonly", "value", pin), 3)
	}
	if _,item := it.FindItem(items, "clientid__phone1"); item != nil {
		item.SetItem("function_fancy_phone_control", "format/inputFn")
	}
	if _,item := it.FindItem(items, "phone_search"); item != nil {
		item.SetItem("html","type")
		div := likdom.BuildDivClassId("phone_search","phone_search","align=center")
		div.BuildString("...")
		item.SetItem(div.ToString(),"value")
	}
	return items
}

//	Закладка характеристик объекта
func (it *BellControl) bellFieldsObject(rule *repo.DataRule, info lik.Seter) lik.Lister {
	items := it.FormInfoCollect(rule, info,"", "bell_object")
	targetid := jone.CalculateIDB(info, "targetid")
	items.InsertItem(lik.BuildSet("type=hidden", "name=u_targetid", "value", targetid), 0)
	if _,item := it.FindItem(items, "offer_search"); item != nil {
		item.SetItem("html","type")
		div := likdom.BuildDivClassId("offer_search","offer_search","align=center")
		if targetid == 0 {
			div.BuildString("...")
		} else if offer := jone.TableOffer.GetElm(targetid); offer != nil {
			div.BuildString(fmt.Sprintf("Объект по заявке №%03d", int(targetid)))
		} else {
			div.BuildString("...")
		}
		item.SetItem(div.ToString(),"value")
	}
	def := jone.CalculateString(info,"definition")
	items.AddItemSet("type=textarea", "label=Комментарии",
		"name=s_definition", "cls=definition",
		"value", def, "editable", it.IsEdit())
	return items
}

//	Закладка заявки
func (it *BellControl) collectOfferFields(rule *repo.DataRule, info lik.Seter) lik.Lister {
	items := lik.BuildList()
	cancel := jone.CalculateString(info,"cancel")
	offerid := jone.CalculateIDB(info,"offerid")
	offer := jone.TableOffer.GetElm(offerid)
	if offer != nil {
		text := fmt.Sprintf("По контакту создана заявка №%03d", int(offerid))
		items.AddItemSet("type=html", "value", text)
		text = jone.CalculateElmText(offer)
		items.AddItemSet("type=html", "value", text)
	} else if cancel != "" {
		text := jone.CalculatePartIdText("bell", likbase.StrToIDB(it.Sel))
		items.AddItemSet("type=html", "value", text)
	} else {
		if data := jone.CalculateTranslate(info, "target"); data != "" {
			items.AddItemSet("type=html", "value", fmt.Sprintf("Цель обращения: <b>%s</b>", data))
		}
		if data := jone.CalculateTranslate(info, "segment"); data != "" {
			items.AddItemSet("type=html", "value", fmt.Sprintf("Сегмент рынка: <b>%s</b>", data))
		}
		if data := jone.CalculateTranslate(info, "realty"); data != "" {
			items.AddItemSet("type=html", "value", fmt.Sprintf("Тип недвижимости: <b>%s</b>", data))
		}
	}
	fun := lik.IfString(it.Fun == fancy.FunAdd || it.Fun == fancy.FunMod || it.Fun == fancy.FunEdit, "(1)", "(0)")
	if it.Fun != fancy.FunDel {
		if offer != nil {
			path := fmt.Sprintf("/offershow%d?_tp=1", int(offerid))
			code := show.LinkTextProc("cmd", "Открыть заявку", fmt.Sprintf("lik_window_part('%s')", path))
			items.AddItemSet("type=html", "value", code.ToString())
		} else if cancel != "" {
			code := show.LinkTextProc("cmd", "Отменить дисквалификацию", "bell_accept" + fun)
			items.AddItemSet("type=html", "value", code.ToString())
		} else {
			cmd := show.LinkTextProc("cmd","Создать заявку","bell_newoffer" + fun)
			items.AddItemSet("type=html", "value", cmd.ToString())
			items.AddItemSet("type=html", "value='&nbsp;'")
			cmd = show.LinkTextProc("cmd must","Дисквалифицировать заявку","bell_cancel" + fun)
			items.AddItemSet("type=html", "value", cmd.ToString())
			cmd = it.bellCancelSelect(rule)
			items.AddItemSet("type=html", "value", cmd.ToString())
		}
		if diagns := it.bellTestDiagnosis(rule, info); len(diagns) > 0 {
			items.AddItemSet("type=html", "value", "<hr>")
			for _, dia := range diagns {
				items.AddItemSet("type=html", "cls=must", "value", dia)
			}
		}
	}
	return items
}

//	Дисквалификация заявки
func (it *BellControl) bellCancelSelect(rule *repo.DataRule) likdom.Domer {
	code := likdom.BuildSpace()
	code.BuildString("Причина:")
	sel := code.BuildItem("select", "id=why_cancel")
	if ent := repo.GenDiction.FindEnt("bell_cancel"); ent != nil {
		content := ent.GetContent()
		for _,canc := range content {
			part := canc.GetString("part")
			text := canc.GetString("name")
			sel.BuildItem("option value=" + part).BuildString(text)
		}
	}
	return code
}

//	Инпорт контакта
func (it *BellControl) BellImportContact(rule *repo.DataRule, part string, sid string) {
	rule.SetResponse("ok", "answer")
}

//	Создание заявки
func (it *BellControl) cmdOfferCreate(rule *repo.DataRule, data lik.Seter) {
	if data != nil {
		it.BellCmdWrite(rule, data)
		it.CmdFormDo = ""
	}
	elmbell := jone.GetElm(it.Part, likbase.StrToIDB(it.Sel))
	if elmbell == nil {
		return
	}
	elmoffer := jone.TableOffer.CreateElm()
	target := jone.CalculateElmString(elmbell,"target")
	segment := jone.CalculateElmString(elmbell,"segment")
	realty := jone.CalculateElmString(elmbell,"realty")
	source := jone.CalculateElmString(elmbell,"source")
	if target != "" { jone.SetElmValue(elmoffer, target,"target") }
	if segment != "" { jone.SetElmValue(elmoffer, segment,"segment") }
	jone.SetElmValue(elmoffer, source,"source")
	jone.SetElmValue(elmoffer,"pass","status")
	jone.SetElmValue(elmbell, elmoffer.Id, "offerid")

	if target == "sale" {
		elmobj := jone.TableObject.CreateElm()
		jone.SetElmValue(elmoffer, elmobj.Id, "objectid")
		if realty != "" {
			jone.SetElmValue(elmoffer, realty, "objectid/realty")
		}
		if datas := jone.CalculateElmSet(elmbell,"address"); datas != nil {
			jone.SetElmValue(elmoffer, datas.Clone(),"objectid/address")
		}
		if datas := jone.CalculateElmSet(elmbell,"define"); datas != nil {
			jone.SetElmValue(elmoffer, datas.Clone(),"objectid/define")
		}
	} else {
		if realty != "" {
			jone.SetElmValue(elmoffer, realty, "require/realty")
		}
		if datas := jone.CalculateElmSet(elmbell,"address"); datas != nil {
			for _,set := range datas.Values() {
				jone.SetElmValue(elmoffer, set.Val, "require/" + set.Key)
			}
		}
		if dat := jone.CalculateElmString(elmbell,"define/rooms"); dat != "" {
			jone.SetElmValue(elmoffer, dat,"require/rooms")
		}
	}

	phone := jone.CalculateElmString(elmbell,"clientid/phone1")
	isitrealtor := strings.HasPrefix(jone.CalculateElmString(elmbell,"isitrealtor"), "y")
	if elmcli := repo.SearchClient(phone); elmcli != nil {
		jone.SetElmValue(elmoffer, elmcli.Id, "clientid")
	} else if !isitrealtor {
		elmcli := jone.TableClient.CreateElm()
		if datas := jone.CalculateElmSet(elmbell,"client"); datas != nil {
			for _, set := range datas.Values() {
				jone.SetElmValue(elmcli, set.Val.Clone(), set.Key)
			}
		}
		jone.SetElmValue(elmoffer, elmcli.Id, "clientid")
	}

	jone.SetElmValue(elmoffer, rule.ItSession.IdMember, "memberid")

	at := int(time.Now().Unix())
	jone.SetElmValue(elmoffer, at,"date")
	jone.SetElmValue(elmbell, at,"datedone")
	atbell := jone.CalculateElmInt(elmbell, "date")
	if atbell == 0 { atbell = at }
	notify := fmt.Sprintf("Контакт №%03d", int(elmbell.Id))
	notify += ", тел. " + phone
	repo.AddHistorySet(rule, elmoffer,"date", atbell, "what=contact", "bellid", elmbell.Id, "notify", notify)
	notify = "Создана из контакта " + jone.CalculateElmSid(elmbell)
	repo.AddHistorySet(rule, elmoffer,"date", at+1, "what=create", "bellid", elmbell.Id, "notify", notify)
	repo.AddHistorySet(rule, elmoffer,"date", at+2, "what=status", "status", "pass")
	//rule.SetPagePush(fmt.Sprintf("offerstaff%d", int(elmoffer.Id)))
	path := fmt.Sprintf("/%s%d/offershow%d?_tp=1", target, int(elmoffer.Id), int(elmoffer.Id))
	rule.SetResponse(path,"_function_lik_window_part")
}

//	Запись изменений
func (it *BellControl) BellCmdWrite(rule *repo.DataRule, data lik.Seter) bool {
	ok := false
	oldtargid := lik.IDB(0)
	if it.Fun == fancy.FunAdd {
		if elm := jone.TableBell.GetElm(likbase.StrToIDB(it.Sel)); elm != nil {
			oldtargid = elm.GetIDB("targetid")
		}
	}
	if elm := it.ListCmdElm(rule); elm != nil {
		if elm.GetIDB("receptorid") == 0 {
			elm.SetValue(rule.ItSession.IdMember, "receptorid")
		}
		clid := lik.IDBFromXS(data.GetString("u_clientid"))
		client := jone.TableClient.GetElm(clid)
		if client == nil {
			phone := lik.StringFromXS(data.GetString("p_clientid__phone1"))
			if client = repo.SearchClient(phone); client == nil {
				client = jone.TableClient.CreateElm()
			}
		}
		if client != nil {
			elm.SetValue(client.Id, "clientid")
		} else if elm.GetSet("clientid") == nil {
			elm.SetValue(lik.BuildSet(), "clientid")
		}
		data.SetItem(nil, "u_clientid")
		it.ListCmdUpdate(rule, elm, data)
		it.Sel = likbase.IDBToStr(elm.Id)
		rule.SetLocateId(elm.Id)
		if targid := elm.GetIDB("targetid"); targid > 0 && targid != oldtargid {
			it.bellAddToContact(rule, elm.Id, elm.GetIDB("clientid"), targid, 0)
		}
		if it.Fun == fancy.FunMod {
			it.CmdFormDo = fmt.Sprintf("show_%d", int(elm.Id))
		}
		rule.SetPagePart(0, fmt.Sprintf("/bell%d", int(elm.Id)))
		ok = true
	}
	it.Fun = fancy.FunNo
	return ok
}

//	Добавление контакта в заявку
func (it *BellControl) bellAddToContact(rule *repo.DataRule, idbell lik.IDB, idclient lik.IDB, idoffer lik.IDB, at int) {
	if offer := jone.TableOffer.GetElm(idoffer); offer != nil {
		text := ""
		if bell := jone.TableOffer.GetElm(idbell); bell != nil {
			text += "Контакт " + jone.CalculateElmSid(bell)
		}
		if client := jone.TableOffer.GetElm(idclient); client != nil {
			text += ", тел. " + client.GetString("phone1")
		}
		conts := repo.GetHistory(offer, "contact")
		for _,cont := range conts {
			if cont.GetIDB("bellid") == idbell && cont.GetString("notify") == text {
				//if ato := cont.GetInt("date"); at - ato < 3600 {
				at = 0
				break
			}
		}
		if at > 0 {
			repo.AddHistorySet(rule, offer, "what=contact",
				"bellid", idbell, "notify", text, "date", at)
		}
	}
}

//	Проверка готовности к заявке
func (it *BellControl) bellTestDiagnosis(rule *repo.DataRule, info lik.Seter) []string {
	diagns := []string{}
	if jone.CalculateString(info,"clientid/phone1") == "" {
		diagns = append(diagns, "Не указан телефон клиента")
	}
	//if repo.CalculateString(info,"isitrealtor") == "" {
	//	diagns = append(diagns, "Не указан признак риэлтора")
	//}
	if jone.CalculateString(info,"target") == "" {
		diagns = append(diagns, "Не указана цель заявки")
	}
	if jone.CalculateIDB(info,"receptorid") == 0 {
		diagns = append(diagns, "Не указан принявший сотрудник")
	}
	if jone.CalculateIDB(info,"memberid") == 0 {
		diagns = append(diagns, "Не указан ответственный риэлтор")
	}
	return diagns
}

//	Переход к заявке
func (it *BellControl) bellGoOffer(rule *repo.DataRule, id lik.IDB) {
	rule.SetPagePush(fmt.Sprintf("offershow%d", int(id)))
}

//	Вход в строку
func (it *BellControl) BellEnterRow(rule *repo.DataRule) {
	if sel := rule.Shift(); sel != "" {
		it.Sel = sel
	}
	if it.Sel != "" {
		if elm := jone.GetElm(it.Part, likbase.StrToIDB(it.Sel)); elm != nil {
			parm := fmt.Sprintf("%s_%s", it.Main, it.Zone)
			idmem := jone.CalculateElmIDB(elm, "memberid")
			idrec := jone.CalculateElmIDB(elm, "receptorid")
			if rule.IAmAdmin() || idrec == 0 ||
				repo.ProbeItMy(rule, "member", idmem) ||
				repo.ProbeItMy(rule, "member", idrec) ||
				rule.IAmManager() && repo.ProbeItDep(rule, "member", idmem)  ||
				rule.IAmManager() && repo.ProbeItDep(rule, "member", idrec) {
				parm += "_edit"
			} else {
				parm += "_show"
			}
			rule.SetResponse(parm, "_function_fancy_trio_form")
		}
	}
}

