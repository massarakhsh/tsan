//	Модуль списка заявок
package offer

import (
	"github.com/massarakhsh/tsan/control"
	"github.com/massarakhsh/tsan/fancy"
	"github.com/massarakhsh/tsan/jone"
	"github.com/massarakhsh/tsan/repo"
	"github.com/massarakhsh/tsan/show"
	"github.com/massarakhsh/lik"
	"github.com/massarakhsh/lik/likbase"
	"github.com/massarakhsh/lik/likdom"
	"fmt"
	"time"
)

//	Дескриптор заявки
type OfferControl struct {
	control.ListControl
}

//	Интерфейс вычисления поля
type dealOfferCalculate struct {
	It	*OfferControl
}
func (it *dealOfferCalculate) Run(rule *repo.DataRule, info lik.Seter, part string, format string, isform bool) lik.Itemer {
	return it.It.OfferCalculate(rule, info, part, format, isform)
}

//	Интерфейс вывода таблицы
type dealOfferGridFill struct {
	It	*OfferControl
}
func (it *dealOfferGridFill) Run(rule *repo.DataRule) {
	it.It.OfferGridFill(rule)
}

//	Интерфейс проверки поля
type dealOfferFieldProbe struct {
	It	*OfferControl
}
func (it *dealOfferFieldProbe) Run(rule *repo.DataRule, field lik.Seter) bool {
	return it.It.OfferFieldProbe(rule, field)
}

//	Интерфейс проверки строки
type dealOfferMakeProbe struct {
	It	*OfferControl
}
func (it *dealOfferMakeProbe) Run(rule *repo.DataRule, elm *likbase.ItElm) bool {
	return it.It.OfferMakeProbe(rule, elm)
}

//	Интерфейс вывода формы
type dealOfferElmForm struct {
	It	*OfferControl
}
func (it *dealOfferElmForm) Run(rule *repo.DataRule, elm *likbase.ItElm) {
	it.It.OfferElmForm(rule, elm)
}

//	Интерфейс вывода строки
type dealOfferRowFill struct {
	It	*OfferControl
}
func (it *dealOfferRowFill) Run(rule *repo.DataRule, elm *likbase.ItElm, row lik.Seter) {
	it.It.OfferRowFill(rule, elm, row)
}

//	Интерфейс команд
type dealOfferExecute struct {
	It	*OfferControl
}
func (it *dealOfferExecute) Run(rule *repo.DataRule, cmd string, data lik.Seter) {
	it.It.OfferExecute(rule, cmd, data)
}

//	Конструктор дескриптора
func BuildListOffer(rule *repo.DataRule, mode string, id lik.IDB) *OfferControl {
	it := &OfferControl{}
	it.ControlInitialize(mode, id)
	it.ListInitialize(rule, mode, "offer")
	it.ItCalculate = &dealOfferCalculate{it}
	it.ItGridFill = &dealOfferGridFill{it}
	it.ItFieldProbe = &dealOfferFieldProbe{it}
	it.ItListMakeProbe = &dealOfferMakeProbe{it}
	it.ItRowFill = &dealOfferRowFill{it}
	it.ItFormElm = &dealOfferElmForm{it}
	it.ItExecute = &dealOfferExecute{it}
	return it
}

//	Обработка команд
func (it *OfferControl) OfferExecute(rule *repo.DataRule, cmd string, data lik.Seter) {
	if cmd == "all" || cmd == it.Main {
		it.OfferExecute(rule, rule.Shift(), data)
	} else if cmd == "toshow" || cmd == "toenter" {
		if sel := rule.Top(); sel != "" {
			it.Sel = sel
		}
		it.GoWindowMode(rule, "offershow", it.Sel)
	} else {
		it.ListExecute(rule, cmd, data)
	}
}

//	Вывад таблицы
func (it *OfferControl) OfferGridFill(rule *repo.DataRule) {
	rule.SetMemberParam(it.GetMode(),"context/target")
	it.ListGridFill(rule)
	it.ListGridCommandSegment(rule)
	it.ListGridCommandRealty(rule)
}

//	Вывод формы
func (it *OfferControl) OfferElmForm(rule *repo.DataRule, elm *likbase.ItElm) {
	it.Form.Items = it.FormElmFill(rule, elm,"")
	if it.Fun == fancy.FunShow {
		name := "Заявка"
		target := jone.CalculateElmString(elm,"target")
		if target == "sale" {
			name += " на продажу"
		} else if target == "buy" {
			name += " на покупку"
		}
		it.SetTitle(rule, it.Fun, name)
		idmem := jone.CalculateElmIDB(elm, "memberid")
		if rule.IAmAdmin() ||
			repo.ProbeItMy(rule, "member", idmem) ||
			rule.IAmManager() && repo.ProbeItDep(rule, "member", idmem) {
			it.AddTitleToolText(rule, "Изменить", "function_fancy_form_toedit")
		}
		if rule.IAmAdmin() {
			it.AddTitleToolText(rule, "Удалить", "function_fancy_form_todelete")
		}
		it.AddTitleToolText(rule, "Закрыть", "function_fancy_form_cancel")
	} else if it.Fun == fancy.FunAdd {
		it.SetTitle(rule, it.Fun, "Создание заявки")
		it.AddTitleToolText(rule, "Сохранить", "function_fancy_form_write")
		it.AddTitleToolText(rule, "Отменить", "function_fancy_form_cancel")
	} else if it.Fun == fancy.FunMod || it.Fun == fancy.FunEdit {
		if target := elm.GetString("target"); target == "sale" {
			it.Form.SingleCho = true
		}
		it.SetTitle(rule, it.Fun, "Редактирование заявки")
		it.AddTitleToolText(rule, "Сохранить", "function_fancy_form_write")
		it.AddTitleToolText(rule, "Отменить", "function_fancy_form_cancel")
	} else if it.Fun == fancy.FunDel {
		it.SetTitle(rule, it.Fun, "Удаление заявки")
		it.AddTitleToolText(rule, "Действительно удалить?", "function_fancy_real_delete")
		it.AddTitleToolText(rule, "Отменить", "function_fancy_form_cancel")
	}
}

//	Проверка поля
func (it *OfferControl) OfferFieldProbe(rule *repo.DataRule, field lik.Seter) bool {
	if !it.FancyFieldProbe(rule, field) { return false }
	target := it.GetMode()
	realty := rule.GetMemberParamString("context/realty")
	tags := field.GetInt("tags")
	return it.FancyProbeTags(target, realty, tags)
}

//	Проверка строки
func (it *OfferControl) OfferMakeProbe(rule *repo.DataRule, elm *likbase.ItElm) bool {
	accept := (elm != nil)
	if accept {
		if tar := jone.CalculateElmString(elm,"target"); tar != "" && tar != it.Mode {
			accept = false
		}
	}
	if accept && it.Selector.ItSegment != "" && it.Selector.ItSegment != jone.ItAll {
		if seg := jone.CalculateElmString(elm,"segment"); seg != "" && seg != it.Selector.ItSegment {
			accept = false
		}
	}
	if accept && it.Selector.ItRealty != "" && it.Selector.ItRealty != jone.ItAll {
		if rea := jone.CalculateElmString(elm, "objectid/realty"); rea != "" && rea != it.Selector.ItRealty {
			accept = false
		}
	}
	if accept && it.Selector.ItStatus == "active" {
		if status := jone.CalculateElmString(elm,"status"); status != jone.ItActive {
			accept = false
		}
	}
	if accept && it.Selector.ItLocate == jone.ItMy {
		if !repo.ProbeItMy(rule, "member", jone.CalculateElmIDB(elm, "memberid")) {
			accept = false
		}
	} else if accept && it.Selector.ItLocate == jone.ItDep {
		if !repo.ProbeItDep(rule, "member", jone.CalculateElmIDB(elm, "memberid")) {
			accept = false
		}
	}
	if accept && !it.ListMakeProbe(rule, elm) {
		accept = false
	}
	return accept
}

//	Изготовление строки
func (it *OfferControl) OfferRowFill(rule *repo.DataRule, elm *likbase.ItElm, row lik.Seter) bool {
	it.ListRowFill(rule, elm, row)
	if row == nil {
		return false
	}
	status := jone.CalculateElmString(elm,"status")
	picatt := ""
	picyandex := ""
	picavito := ""
	piccian := ""
	picdomclick := ""
	pictsan := ""
	if elm.Id == 13 {
		picyandex += ""
	}
	expok := elm.GetString("status") == jone.ItActive &&
		elm.GetBool("export/ready") && elm.GetBool("export/confirm") &&elm.GetBool("export/enable")
	if list := repo.GetExportList(); list != nil {
		for nc := 0; nc < list.Count(); nc++ {
			if exp := list.GetSet(nc); exp != nil {
				part := exp.GetString("part")
				if adv := elm.GetSet("export/"+part); adv != nil && adv.GetBool("ready") {
					diag := adv.GetString("diagnosis")
					if diag != "" && picatt == "" {
						picatt = diag
					}
					if !expok {
					} else if part == "yandex" {
						if diag != "" {
							picyandex = diag
						} else {
							picyandex = "ok"
						}
					} else if part == "avito" {
						if diag != "" {
							picavito = diag
						} else {
							picavito = "ok"
						}
					} else if part == "cian" {
						if diag != "" {
							piccian = diag
						} else {
							piccian = "ok"
						}
					} else if part == "domclick" {
						if diag != "" {
							picdomclick = diag
						} else {
							picdomclick = "ok"
						}
					} else if part == "tsan" {
						if diag != "" {
							pictsan = diag
						} else {
							pictsan = "ok"
						}
					}
				}
			}
		}
	}
	pics := likdom.BuildSpace()
	if picatt != "" {
		pics.BuildUnpairItem("img", "src", "/images/control.png", "title", picatt)
	} else {
		pics.BuildUnpairItem("img", "src","/images/1616.png")
	}
	if picyandex == "ok" {
		pics.BuildUnpairItem("img", "src", "/images/yandex.png", "title=Яндекс")
	} else if picyandex != "" {
		pics.BuildUnpairItem("img", "src", "/images/yandexb.png", "title=Яндекс: "+picyandex)
	} else {
		pics.BuildUnpairItem("img", "src", "/images/1616.png")
	}
	if picavito == "ok" {
		pics.BuildUnpairItem("img", "src", "/images/avito.png", "title=Авито")
	} else if picavito != "" {
		pics.BuildUnpairItem("img", "src", "/images/avitob.png", "title=Авито: "+picavito)
	} else {
		pics.BuildUnpairItem("img", "src", "/images/1616.png")
	}
	if piccian == "ok" {
		pics.BuildUnpairItem("img", "src", "/images/cian.png", "title=ЦИАН")
	} else if piccian != "" {
		pics.BuildUnpairItem("img", "src", "/images/cianb.png", "title=ЦИАН: "+piccian)
	} else {
		pics.BuildUnpairItem("img", "src", "/images/1616.png")
	}
	if picdomclick == "ok" {
		pics.BuildUnpairItem("img", "src", "/images/domclick.png", "title=ДомКлик")
	} else if picdomclick != "" {
		pics.BuildUnpairItem("img", "src", "/images/domclickb.png", "title=ДомКлик: "+piccian)
	} else {
		pics.BuildUnpairItem("img", "src", "/images/1616.png")
	}
	if pictsan != "" {
		//pics.BuildUnpairItem("img", "src", "/images/tsan.png", "title=ЦАН: "+pictsan)
	} else {
		//pics.BuildUnpairItem("img", "src", "/images/1616.png")
	}
	row.SetItem(pics.ToString(), "pic")
	if status != "" {
		//row.SetItem(status, "status")
	}
	if deal := repo.SeekOfferDeal(elm.Id); deal.IdDeal > 0 {
		if dc := jone.CalculatePartIdInt("deal", deal.IdDeal, "date_close"); dc > 0 {
			text := time.Unix(int64(dc), 0).Format("2006/01/02")
			path := fmt.Sprintf("/offerdeal%d?_tp=1", int(elm.Id))
			cmd := show.LinkTextProc("cmd", text, fmt.Sprintf("lik_window_part('%s')", path))
			row.SetItem(cmd.ToString(), "deal")
		}
	}
	row.SetItem(fmt.Sprintf("/offershow%d?_tp=1", int(elm.Id)), "pathopen")
	return true
}

//	Вычисление поля
func (it *OfferControl) OfferCalculate(rule *repo.DataRule, info lik.Seter, part string, format string, isform bool) lik.Itemer {
	var value lik.Itemer
	if info == nil {
	} else if part == "photos" {
		cnt := ""
		if list := jone.CalculateList(info, "objectid/picture"); list != nil && list.Count() > 0 {
			cnt = lik.IntToStr(list.Count())
		}
		value = lik.BuildItem(cnt)
	} else if part == "coststart" {
		coststart := 0
		if listcost := repo.ExtractHistory(info.GetList("history"), "cost"); listcost != nil && len(listcost) > 0 {
			coststart = listcost[0].GetInt("cost")
		} else {
			coststart = info.GetInt("cost")
		}
		if coststart != 0 {
			value = lik.BuildItem(coststart)
		}
	} else {
		value = it.FancyCalculate(rule, info, part, format, isform)
	}
	return value
}

//	Вход в строку
func (it *OfferControl) OfferEnterRow(rule *repo.DataRule) {
	rule.SetPagePush(fmt.Sprintf("offershow%d", int(it.IdMain)))
}

