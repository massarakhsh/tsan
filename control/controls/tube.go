package controls

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

//	Обработчик команд
type dealTubeExecute struct {
	It	*TubeControl
}
func (it *dealTubeExecute) Run(rule *repo.DataRule, cmd string, data lik.Seter) {
	it.It.Self.RunTubeExecute(rule, cmd, data)
}

//	Обработчик отображения таблицы
type dealTubeGridFill struct {
	It	*TubeControl
}
func (it *dealTubeGridFill) Run(rule *repo.DataRule) {
	it.It.Self.RunTubeGridFill(rule)
}

//	Обработчик отображения страницы
type dealTubePageFill struct {
	It	*TubeControl
}
func (it *dealTubePageFill) Run(rule *repo.DataRule) lik.Lister {
	return it.It.TubePageFill(rule)
}

//	Обработчик отображения строки
type DealTubeRow interface {
	Run(rule *repo.DataRule, data lik.Seter, row lik.Seter)
}
type dealTubeRow struct {
	It	*TubeControl
}
func (it *dealTubeRow) Run(rule *repo.DataRule, data lik.Seter, row lik.Seter) {
	it.It.TubeRowFill(rule, data, row)
}

//	Обработчик отображения формы
type dealTubeFormFill struct {
	It	*TubeControl
}
func (it *dealTubeFormFill) Run(rule *repo.DataRule) {
	it.It.TubeFormFill(rule)
}

//	Интерфейс заполнения таблицы
type DealTubeFill interface {
	CmdTubeFill(rule *repo.DataRule)
}

//	Интерфейс подготовки к редактированию
type DealPreEdit interface {
	PreEdit(rule *repo.DataRule, elm *likbase.ItElm, items lik.Lister)
}

//	Дескриптор окна событий
type TubeControl struct {
	control.DataControl
	fancy.TableFancy
	Self       ShowTuber
	Title      string
	ItTubeFill DealTubeFill
	ItPreEdit  DealPreEdit
	ItTubeRow  DealTubeRow
}

//	Интерфейс окна событий
type ShowTuber interface {
	RunTubeExecute(rule *repo.DataRule, cmd string, data lik.Seter)
	RunTubeGridFill(rule *repo.DataRule)
	RunTubeFinalEdit(rule *repo.DataRule, elm *likbase.ItElm)
}

//	Инициализация окна событий
func (it *TubeControl) TubeInitialize(rule *repo.DataRule, main string, mode string, id lik.IDB, title string) {
	it.Title = title
	it.ControlInitializeZone(main, id, mode)
	it.TableInitialize(rule, main,"offer", mode)
	it.ItGridFill = &dealTubeGridFill{it}
	it.ItPageFill = &dealTubePageFill{it}
	it.ItFormFill = &dealTubeFormFill{it}
	it.ItTubeRow = &dealTubeRow{it}
	it.ItExecute = &dealTubeExecute{it}
}

//	Обработка событий
func (it *TubeControl) TubeExecute(rule *repo.DataRule, cmd string, data lik.Seter) {
	if cmd == "delete" {
		it.cmdDelete(rule)
	} else if cmd == "write" {
		it.cmdWrite(rule, data)
	} else {
		it.TableExecute(rule, cmd, data)
	}
}

//	Отображение событий
func (it *TubeControl) buildShowTube(rule *repo.DataRule, pater likdom.Domer, sx int, sy int) {
	it.SetSize(sx, sy)
	pater.AppendItem(show.BuildFancyGrid(it.Main, it.Zone))
}

//	Отображение окна событий
func (it *TubeControl) TubeGridFill(rule *repo.DataRule) {
	it.TableGridFill(rule)
	it.GridBuildColumns(rule, true, true)
	it.Grid.SetSize(it.Grid.Width, it.Grid.Height)
	it.AddCommandItem(rule, 80, lik.BuildSet("type=text", "text", it.Title+"&nbsp; "))
	it.AddCommandImg(rule, 910, "Открыть", "toshow","show")
}

//	Отображение страницы событий
func (it *TubeControl) TubePageFill(rule *repo.DataRule) lik.Lister {
	list := it.tubeInfoBill(rule)
	return it.tubeGenRows(rule, list, 1000)
}

//	Отбор событий
func (it *TubeControl) tubeInfoBill(rule *repo.DataRule) []lik.Seter {
	if elm := jone.GetElm(it.Part, it.IdMain); elm != nil {
		return repo.GetHistory(elm, it.Zone)
	}
	return []lik.Seter{}
}

//	Генерация списка строк
func (it *TubeControl) tubeGenRows(rule *repo.DataRule, list []lik.Seter, limit int) lik.Lister {
	rows := it.TablePageFill(rule)
	for nelm := 0; nelm < len(list) && (limit == 0 || nelm < limit); nelm++ {
		item := list[nelm]
		idx := item.GetString("idx")
		if row := it.GridInfoRow(rule, idx, item); row != nil {
			it.RunRowFill(rule, item, row)
			rows.AddItems(row)
		}
	}
	return rows
}

//	Заполнение строки
func (it *TubeControl) RunRowFill(rule *repo.DataRule, data lik.Seter, row lik.Seter) {
	it.ItTubeRow.Run(rule, data, row)
}

//	Встроенное заполнение строки
func (it *TubeControl) TubeRowFill(rule *repo.DataRule, data lik.Seter, row lik.Seter) {
	if description := row.GetString("description"); description == "" {
		if what := data.GetString("what"); what == "create" {
			description = it.TubeRowCreate(rule, data, row)
		} else if what == "modify" {
			description = it.TubeRowModify(rule, data, row)
		} else if what == "contact" {
			description = it.TubeRowContact(rule, data, row)
		} else if what == "promo" {
			description = it.TubeRowPromo(rule, data, row)
		}
		row.SetItem(description, "description")
	}
}

//	Строка создания позиции
func (it *TubeControl) TubeRowCreate(rule *repo.DataRule, data lik.Seter, row lik.Seter) string {
	description :=  row.GetString("notify")
	if description == "" {
		description = "Создание"
	}
	return description
}

//	Строка изменения позиции
func (it *TubeControl) TubeRowModify(rule *repo.DataRule, data lik.Seter, row lik.Seter) string {
	description := ":"
	if changes := data.GetList("changes"); changes != nil {
		ches := lik.BuildSet()
		for nm := 0; nm < changes.Count(); nm++ {
			part := changes.GetString(nm)
			if lik.RegExCompare(part,"objectid__define") {
				ches.SetItem(true, "define")
			} else if lik.RegExCompare(part,"objectid__address") {
				ches.SetItem(true, "address")
			} else if lik.RegExCompare(part,"objectid") {
				ches.SetItem(true, "object")
			} else if lik.RegExCompare(part,"clientid") {
				ches.SetItem(true, "client")
			} else if lik.RegExCompare(part,"contractid") {
				ches.SetItem(true, "contract")
			} else {
				ches.SetItem(true, "other")
			}
		}
		if ches.GetBool("define") {
			description += ", характеристики"
		}
		if ches.GetBool("address") {
			description += ", адрес"
		}
		if ches.GetBool("object") {
			description += ", объект"
		}
		if ches.GetBool("client") {
			description += ", клиент"
		}
		if ches.GetBool("contract") {
			description += ", контракт"
		}
		if ches.GetBool("other") {
			description += ", прочее"
		}
	}
	return description
}

//	Строка указания контакта
func (it *TubeControl) TubeRowContact(rule *repo.DataRule, data lik.Seter, row lik.Seter) string {
	description := ""
	if idbell := data.GetIDB("bellid"); idbell != 0 {
		if bell := jone.TableBell.GetElm(idbell); bell != nil && !it.IsEdit() {
			description = jone.CalculateElmTranslate(bell, "target")
			text := fmt.Sprintf("№%03d", int(idbell))
			proc := fmt.Sprintf("bell_bind(%d)", int(idbell))
			code := show.LinkTextProc("cmd", text, proc)
			row.SetItem(code.ToString(), "bell")
		}
	}
	if description == "" {
		description = jone.CalculateString(data, "notify")
	}
	return description
}

//	Строка продвижения
func (it *TubeControl) TubeRowPromo(rule *repo.DataRule, data lik.Seter, row lik.Seter) string {
	description := jone.CalculateTranslate(data, "promo")
	if act := data.GetString("action"); act == "yy" || act == "true" {
		description += ". Установлено"
	} else if act == "nn" || act == "false" {
		description += ". Снято"
	}
	return description
}

//	Вывод формы
func (it *TubeControl) TubeFormFill(rule *repo.DataRule) {
	elm := jone.GetElm(it.Part, it.IdMain)
	idx := lik.StrToInt(it.Sel)
	if it.Fun == fancy.FunAdd {
		idx = repo.ReserveHistory(elm)
	}
	basic := fmt.Sprintf("history/idx%d", idx)
	it.Form.Items = it.FormElmFill(rule, elm, basic)
	it.Form.Items.AddItemSet("type=hidden", "name=id", fmt.Sprintf("value=%d", idx))
	it.Form.Items.AddItemSet("type=hidden", fmt.Sprintf("name=history__idx%d__what", idx), fmt.Sprintf("value=%s", it.Zone))
	name := ""
	if it.Zone == "cost" {
		name = "цены"
	} else if it.Zone == "status" {
		name = "статуса"
	} else if it.Zone == "contact" {
		name = "контакта"
	}
	if it.Fun == fancy.FunShow {
		it.SetTitle(rule, it.Fun, "Карточка "+name)
		if rule.IAmAdmin() {
			it.AddTitleToolText(rule, "Изменить", "function_fancy_form_toedit")
			it.AddTitleToolText(rule, "Удалить", "function_fancy_form_todelete")
		}
		it.AddTitleToolText(rule, "Закрыть", "function_fancy_form_cancel")
	} else if it.Fun == fancy.FunAdd {
		it.SetTitle(rule, it.Fun, "Создание "+name)
		it.AddTitleToolText(rule, "Сохранить", "function_fancy_form_write")
		it.AddTitleToolText(rule, "Отменить", "function_fancy_form_cancel")
		if _,item := it.FindItem(it.Form.Items,"date"); item != nil {
			item.SetItem(time.Now().Format("2006/01/02"), "value")
		}
		if _,item := it.FindItem(it.Form.Items,"time"); item != nil {
			item.SetItem(time.Now().Format("15:04"), "value")
		}
	} else if it.Fun == fancy.FunMod || it.Fun == fancy.FunEdit {
		it.SetTitle(rule, it.Fun, "Редактирование "+name)
		it.AddTitleToolText(rule, "Сохранить", "function_fancy_form_write")
		it.AddTitleToolText(rule, "Отменить", "function_fancy_form_cancel")
		if _,item := it.FindItem(it.Form.Items,"date"); item != nil {
			item.SetItem(false, "editable")
		}
		if _,item := it.FindItem(it.Form.Items,"time"); item != nil {
			item.SetItem(false, "editable")
		}
	} else if it.Fun == fancy.FunDel {
		it.SetTitle(rule, it.Fun, "Удаление "+name)
		it.AddTitleToolText(rule, "Действительно удалить?", "function_fancy_real_delete")
		it.AddTitleToolText(rule, "Отменить", "function_fancy_form_cancel")
	}
}

//	Команда записи изменений
func (it *TubeControl) cmdWrite(rule *repo.DataRule, data lik.Seter) {
	elm := jone.GetElm(it.Part, it.IdMain)
	it.UpdateElmData(rule, elm, data)
	it.finalEdit(rule, elm)
	elm.OnModify()
	rule.OnChangeData()
}

//	Команда удаления
func (it *TubeControl) cmdDelete(rule *repo.DataRule) {
	elm := jone.GetElm(it.Part, it.IdMain)
	if num, _ := it.FindListData(rule, lik.StrToInt(it.Sel)); num >= 0 {
		if list := jone.CalculateElmList(elm,"history"); list != nil {
			list.DelItem(num)
			elm.OnModify()
		}
	}
	it.finalEdit(rule, elm)
	rule.OnChangeData()
}

//	Завершение редактирования
func (it *TubeControl) finalEdit(rule *repo.DataRule, elm *likbase.ItElm) {
	repo.SortHistory(elm)
	it.Self.RunTubeFinalEdit(rule, elm)
}

//	Выбрать данные для списка
func (it *TubeControl) FindListData(rule *repo.DataRule, idx int) (int, lik.Seter) {
	if elm := jone.GetElm(it.Part, it.IdMain); elm != nil {
		if history := jone.CalculateElmList(elm,"history"); history != nil {
			for nh := 0; nh < history.Count(); nh++ {
				if event := history.GetSet(nh); event != nil {
					if idx == event.GetInt("idx") {
						return nh, event
					}
				}
			}
		}
	}
	return -1, nil
}
