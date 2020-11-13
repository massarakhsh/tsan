package promo

import (
	"github.com/massarakhsh/tsan/control"
	"github.com/massarakhsh/tsan/fancy"
	"github.com/massarakhsh/tsan/jone"
	"github.com/massarakhsh/tsan/repo"
	"github.com/massarakhsh/tsan/show"
	"github.com/massarakhsh/lik"
	"github.com/massarakhsh/lik/likdom"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math/rand"
	"strings"
	"time"
)

//	Дескриптор специальных операций
type PromouterControl struct {
	control.DataControl
	fancy.TableFancy
}

//	Интерфейс команд
type dealPromouterExecute struct {
	It	*PromouterControl
}
func (it *dealPromouterExecute) Run(rule *repo.DataRule, cmd string, data lik.Seter) {
	it.It.PromouterExecute(rule, cmd, data)
}

//	Инерфейс отображения таблицы
type dealPromouterGridFill struct {
	It	*PromouterControl
}
func (it *dealPromouterGridFill) Run(rule *repo.DataRule) {
	it.It.PromouterGridFill(rule)
}

//	Интерфейс отображения формы
type dealPromouterFormFill struct {
	It	*PromouterControl
}
func (it *dealPromouterFormFill) Run(rule *repo.DataRule) {
	it.It.PromouterFormFill(rule)
}

//	Интерфейс отображения страницы
type dealPromouterPageFill struct {
	It	*PromouterControl
}
func (it *dealPromouterPageFill) Run(rule *repo.DataRule) lik.Lister {
	return it.It.PromouterPageFill(rule)
}

//	Дескриптор специального продвижения
func BuildPromouter(rule *repo.DataRule, main string, id lik.IDB) *PromouterControl {
	it := &PromouterControl{ }
	it.ControlInitializeZone(main, id, "promouter")
	it.TableInitialize(rule, main,"offer","promouter")
	it.ItExecute = &dealPromouterExecute{it}
	it.ItGridFill = &dealPromouterGridFill{it}
	it.ItPageFill = &dealPromouterPageFill{it}
	it.ItFormFill = &dealPromouterFormFill{it}
	return it
}

//	Команды спецопераций
func (it *PromouterControl) PromouterExecute(rule *repo.DataRule, cmd string, data lik.Seter) {
	if cmd == "all" || cmd == "promouter" {
		it.PromouterExecute(rule, rule.Shift(), data)
	} else if cmd == "write" {
		it.cmdWrite(rule, data)
	} else if lik.RegExCompare(cmd, "(enable)") {
		it.cmdEdit(rule, cmd)
	} else if cmd == "delete" {
		it.cmdDelete(rule)
	} else {
		it.TableExecute(rule, cmd, data)
	}
}

//	Отображение окна
func (it *PromouterControl) BuildShow(rule *repo.DataRule, sx int, sy int) likdom.Domer {
	it.SetSize(sx, sy)
	return show.BuildFancyGrid(it.Main,"promouter")
}

//	Оторбражение таблицы
func (it *PromouterControl) PromouterGridFill(rule *repo.DataRule) {
	it.TableGridFill(rule)
	it.GridBuildColumns(rule, true, true)
	for _,pot := range it.Grid.Columns.Values() {
		if sot := pot.ToSet(); sot != nil {
			if sot.GetString("type") == "checkbox" {
				avi := true
				if strings.HasSuffix(sot.GetString("index"), "enable") && !rule.IAmManager() && !rule.IAmAdmin() {
					avi = false
				}
				sot.SetItem(avi, "editable")
			}
		}
	}
	it.AddCommandItem(rule, 80, lik.BuildSet("type=text", "text", "Продвижения"))
	it.AddCommandImg(rule, 910, "Открыть", "toshow", "show")
	it.AddCommandImg(rule, 915, "Добавить", "toadd", "add")
	it.Grid.AddEventAction("cellclick", "function_fancy_grid_mark")
}

//	Отображение страницы
func (it *PromouterControl) PromouterPageFill(rule *repo.DataRule) lik.Lister {
	rows := it.TablePageFill(rule)
	if elm := jone.TableOffer.GetElm(it.IdMain); elm != nil {
		if list := jone.CalculateElmList(elm, "promo"); list != nil {
			for np := 0; np < list.Count(); np++ {
				if data := list.GetSet(np); data != nil {
					id := data.GetString("id")
					if id == "" {
						id = lik.IntToStr(rand.Intn(1000000))
						data.SetItem(id, "id")
						elm.OnModify()
					}
					if row := it.GridInfoRow(rule, id, data); row != nil {
						rows.AddItems(row)
					}
				}
			}
		}
	}
	return rows
}

//	Отображение формы
func (it *PromouterControl) PromouterFormFill(rule *repo.DataRule) {
	var data lik.Seter
	if it.Fun != fancy.FunAdd {
		_, data = it.seekData(rule, it.Sel)
	}
	it.Form.Items = it.FormInfoFill(rule, data, "")
	if it.Fun == fancy.FunShow {
		it.SetTitle(rule, it.Fun, "Карточка продвижения")
		if rule.IAmAdmin() {
			it.AddTitleToolText(rule, "Изменить", "function_fancy_form_toedit")
			it.AddTitleToolText(rule, "Удалить", "function_fancy_form_todelete")
		}
		it.AddTitleToolText(rule, "Закрыть", "function_fancy_form_cancel")
	} else if it.Fun == fancy.FunAdd {
		it.SetTitle(rule, it.Fun, "Создание карточки")
		it.AddTitleToolText(rule, "Сохранить", "function_fancy_form_write")
		it.AddTitleToolText(rule, "Отменить", "function_fancy_form_cancel")
		if _,item := it.FindItem(it.Form.Items,"date"); item != nil {
			item.SetItem(time.Now().Format("2006/01/02"), "value")
		}
		if _,item := it.FindItem(it.Form.Items,"time"); item != nil {
			item.SetItem(time.Now().Format("15:04"), "value")
		}
	} else if it.Fun == fancy.FunMod || it.Fun == fancy.FunEdit {
		it.SetTitle(rule, it.Fun, "Редактирование карточки")
		it.AddTitleToolText(rule, "Сохранить", "function_fancy_form_write")
		it.AddTitleToolText(rule, "Отменить", "function_fancy_form_cancel")
		if _,item := it.FindItem(it.Form.Items,"date"); item != nil {
			item.SetItem(false, "editable")
		}
		if _,item := it.FindItem(it.Form.Items,"time"); item != nil {
			item.SetItem(false, "editable")
		}
	} else if it.Fun == fancy.FunDel {
		it.SetTitle(rule, it.Fun, "Удаление карточки")
		it.AddTitleToolText(rule, "Действительно удалить?", "function_fancy_real_delete")
		it.AddTitleToolText(rule, "Отменить", "function_fancy_form_cancel")
	}
}

//	Запись изменений
func (it *PromouterControl) cmdWrite(rule *repo.DataRule, data lik.Seter) {
	if elm := jone.TableOffer.GetElm(it.IdMain); elm != nil {
		list := jone.CalculateElmList(elm,"promo")
		if list == nil {
			list = lik.BuildList()
			jone.SetElmValue(elm, list, "promo")
			elm.OnModify()
		}
		if list != nil && data != nil {
			var dt lik.Seter
			if it.Fun == fancy.FunAdd {
				dt = list.AddItemSet("id", rand.Intn(1000000))
				elm.OnModify()
			} else if it.Fun == fancy.FunMod || it.Fun == fancy.FunEdit {
				_,dt = it.seekData(rule, it.Sel)
			}
			if dt != nil && data != nil {
				if it.UpdateInfoData(rule, dt, data) != nil {
					elm.OnModify()
				}
			}
		}
		rule.OnChangeData()
	}
}

//	Удаление операции
func (it *PromouterControl) cmdDelete(rule *repo.DataRule) {
	if elm := jone.TableOffer.GetElm(it.IdMain); elm != nil {
		if list := jone.CalculateElmList(elm,"promo"); list != nil {
			if np,_ := it.seekData(rule, it.Sel); np >= 0 {
				list.DelItem(np)
				elm.OnModify()
			}
		}
	}
	rule.OnChangeData()
}

//	Игменение операции
func (it *PromouterControl) cmdEdit(rule *repo.DataRule, cmd string) {
	if elm := jone.TableOffer.GetElm(it.IdMain); elm != nil {
		if _, data := it.seekData(rule, rule.Shift()); data != nil {
			ival := lik.StrToInt(rule.Shift())
			if cmd == "enable" {
				if ival > 0 {
					data.SetItem(true, "enable")
				} else {
					data.SetItem(false, "enable")
				}
				if dt := data.GetInt("datedeal"); dt == 0 {
					data.SetItem(data.GetInt("dateneed"), "datedeal")
				}
				data.SetItem(rule.ItSession.IdMember, "memberid")
				elm.OnModify()
				rule.OnChangeData()
			}
		}
	}
}

//	Поиск операции
func (it *PromouterControl) seekData(rule *repo.DataRule, sid string) (int,lik.Seter) {
	if elm := jone.TableOffer.GetElm(it.IdMain); elm != nil {
		if list := jone.CalculateElmList(elm,"promo"); list != nil {
			for np := 0; np < list.Count(); np++ {
				if dt := list.GetSet(np); dt != nil {
					if dt.GetString("id") == sid {
						return np, dt
					}
				}
			}
		}
	}
	return -1,nil
}

