package window

import (
	"github.com/massarakhsh/tsan/repo"
	"github.com/massarakhsh/tsan/show"
	"github.com/massarakhsh/lik"
	"github.com/massarakhsh/lik/likdom"
)

//	Инициализация дексриптора коллекции
func (it *ClientBox) CollectInitialize(rule *repo.DataRule, frame string, id lik.IDB, mode string) {
	it.IsCollect = true
	it.BoxInitialize(rule, frame, id, mode)
}

//	Конструктор дескриптора окна клиента
func (it *ClientBox) BuildShow(rule *repo.DataRule, sx int, sy int) likdom.Domer {
	it.SetSize(sx, sy)
	if it.IsCollect {
		return show.BuildFancyGrid(it.Frame, it.Mode)
	}
	return nil
}

//	Отображение окна коллекции
func (it *ClientBox) boxShowGrid(rule *repo.DataRule) {
	it.BoxClear()
	it.RunShowGrid(rule)
	ShowCollect(rule, it)
}

//	Отображение страницы коллекции
func (it *ClientBox) boxShowPage(rule *repo.DataRule) {
	rows := it.RunShowPage(rule)
	ShowPage(rule, it, rows)
}

//	Выбор строки коллекции
func (it *ClientBox) boxRowSelect(rule *repo.DataRule) {
	it.IdSelected = rule.Shift()
}

//	Точка входа отображения страницы коллекции
func (it *ClientBox) RunShowGrid(rule *repo.DataRule) {
	it.ItShowGrid.Run(rule)
	it.appendBoxEvents(rule, false)
}

//	Отображение окна общей коллекции
func (it *ClientBox) CollectShowGrid(rule *repo.DataRule) {
	if it.Columns.Count() == 0 {
		it.Columns.AddItemSet("index=id", "title=ID", "width=30")
	}
	if it.IsLockRemote {
		rows := it.RunShowPage(rule)
		it.SetParameter(rows, "data/items")
	}
}

//	Входная точка отображения страницы
func (it *ClientBox) RunShowPage(rule *repo.DataRule) lik.Lister {
	return it.ItShowPage.Run(rule)
}

//	Отображение общей страницы
func (it *ClientBox) CollectShowPage(rule *repo.DataRule) lik.Lister {
	return lik.BuildList()
}

//	Добавление колонки в таблицу
func (it *ClientBox) AddColumnItem(rule *repo.DataRule, datas ...interface{}) {
	it.Columns.AddItemSet(datas...)
}

//	Добавление событий, связанных с коллекциями
func (it *ClientBox) appendBoxEvents(rule *repo.DataRule, layout bool) {
	it.AddEventAction("selectrow", "function_fancy_grid_rowselect")
	it.AddEventAction("contextmenu", "function_fancy_grid_rclick")
	it.AddEventAction("celldblclick", "function_fancy_grid_dblclick")
	if rule.IAmShaman() {
		it.AddEventAction("columnresize", "function_fancy_col_size")
		it.AddEventAction("columndrag", "function_fancy_col_drag")
	}
}

