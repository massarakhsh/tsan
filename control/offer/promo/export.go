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
	"strings"
)

//	Дескриптор окна площадок экспорта
type ExportControl struct {
	control.DataControl
	fancy.TableFancy
}

//	Интерфейс команд
type dealExportExecute struct {
	It	*ExportControl
}
func (it *dealExportExecute) Run(rule *repo.DataRule, cmd string, data lik.Seter) {
	it.It.ExportExecute(rule, cmd, data)
}

//	Создание таблицы
type dealExportGridFill struct {
	It	*ExportControl
}
func (it *dealExportGridFill) Run(rule *repo.DataRule) {
	it.It.ExportGridFill(rule)
}

//	Создание страницы
type dealExportPageFill struct {
	It	*ExportControl
}
func (it *dealExportPageFill) Run(rule *repo.DataRule) lik.Lister {
	return it.It.ExportPageFill(rule)
}

//	Создание формы
type dealExportFormFill struct {
	It	*ExportControl
}
func (it *dealExportFormFill) Run(rule *repo.DataRule) {
	it.It.ExportFormFill(rule)
}

//	Конструктор дескриптора
func BuildExport(rule *repo.DataRule, main string, id lik.IDB) *ExportControl {
	it := &ExportControl{ }
	it.ControlInitializeZone(main, id, "export")
	it.TableInitialize(rule, main,"offer","export")
	it.ItExecute = &dealExportExecute{it}
	it.ItGridFill = &dealExportGridFill{it}
	it.ItPageFill = &dealExportPageFill{it}
	it.ItFormFill = &dealExportFormFill{it}
	return it
}

//	Исполнение команд
func (it *ExportControl) ExportExecute(rule *repo.DataRule, cmd string, data lik.Seter) {
	if cmd == "all" || cmd == "export" {
		it.ExportExecute(rule, rule.Shift(), data)
	} else if cmd == "ready" || cmd == "use" {
		it.ExportCmdTrigger(rule, cmd)
	} else {
		it.TableExecute(rule, cmd, data)
	}
}

//	Отображение экспорта
func (it *ExportControl) BuildShow(rule *repo.DataRule, sx int, sy int) likdom.Domer {
	it.SetSize(sx, sy)
	return show.BuildFancyGrid(it.Main,"export")
}

//	Создание окна
func (it *ExportControl) ExportGridFill(rule *repo.DataRule) {
	it.TableGridFill(rule)
	it.GridBuildColumns(rule, true, true)
	for _,pot := range it.Grid.Columns.Values() {
		if sot := pot.ToSet(); sot != nil {
			if sot.GetString("type") == "checkbox" {
				avi := true
				if strings.HasSuffix(sot.GetString("index"), "use") && !rule.IAmManager() && !rule.IAmAdmin() {
					avi = false
				}
				sot.SetItem(avi, "editable")
			}
		}
	}
	it.AddCommandItem(rule, 80, lik.BuildSet("type=text", "text", "Публикация на площадках"))
	it.Grid.AddEventAction("cellclick", "function_fancy_grid_mark")
}

//	Создание страницы
func (it *ExportControl) ExportPageFill(rule *repo.DataRule) lik.Lister {
	rows := it.TablePageFill(rule)
	if elm := jone.GetElm("offer", it.IdMain); elm != nil {
		if list := repo.GetExportList(); list != nil {
			for nc := 0; nc < list.Count(); nc++ {
				if exp := list.GetSet(nc); exp != nil {
					part := exp.GetString("part")
					adv := elm.GetSet("export/"+part)
					if row := it.GridInfoRow(rule, lik.IntToStr(nc+1), adv); row != nil {
						row.SetItem(exp.GetString("name"), "name")
						rows.AddItems(row)
					}
				}
			}
		}
	}
	if rows.Count() <= 0 {
		rows.AddItemSet("id=1", "nam=НЕТ")
	}
	return rows
}

//	Команда триггера
func (it *ExportControl) ExportCmdTrigger(rule *repo.DataRule, cmd string) {
	numex := lik.StrToInt(rule.Shift())
	on := lik.StrToInt(rule.Shift()) > 0
	if elm := jone.GetElm("offer", it.IdMain); elm != nil {
		if parts := repo.GetExportParts(); parts != nil && numex > 0 && numex <= len(parts) {
			part := parts[numex - 1]
			info := elm.GetSet("export/"+part)
			if info == nil {
				info = lik.BuildSet()
				elm.SetValue(info, "export/"+part)
			}
			if on {
				info.SetItem(true, cmd)
			} else {
				info.SetItem(false, cmd)
			}
			elm.OnModify()
		}
	}
}

//	Создание формы
func (it *ExportControl) ExportFormFill(rule *repo.DataRule) {
	var data lik.Seter
	numex := lik.StrToInt(it.Sel)
	if elm := jone.GetElm("offer", it.IdMain); elm != nil {
		if parts := repo.GetExportParts(); parts != nil && numex > 0 && numex <= len(parts) {
			part := parts[numex-1]
			data = elm.GetSet("export/" + part)
		}
	}
	it.Form.Items = it.FormInfoFill(rule, data, "")
	it.SetTitle(rule, it.Fun, "Карточка публикации")
	it.AddTitleToolText(rule, "Закрыть", "function_fancy_form_cancel")
}

