//	Модули коллекций
package collectproject

import (
	"bitbucket.org/961961/tsan/control/window"
	"bitbucket.org/961961/tsan/one"
	"bitbucket.org/961961/tsan/repo"
	"bitbucket.org/shaman/lik"
)

//	Дескриптор коллекции истории проекта
type Projectlist struct {
	window.ClientBox //	Основан на общей коллекции
}

//	Дескриптор обработки событий
type dealListExecute struct {
	It	*Projectlist
}
func (it *dealListExecute) Run(rule *repo.DataRule, cmd string, data lik.Seter) {
	it.It.ListExecute(rule, cmd, data)
}

//	Дескриптор заголовка таблицы
type dealShowGrid struct {
	It	*Projectlist
}
func (it *dealShowGrid) Run(rule *repo.DataRule) {
	it.It.ListShowGrid(rule)
}

//	Дескриптор заполнения страницы
type dealShowPage struct {
	It	*Projectlist
}
func (it *dealShowPage) Run(rule *repo.DataRule) lik.Lister {
	return it.It.ListShowPage(rule)
}

//	Конструктор дескриптора
func BuildProject(rule *repo.DataRule) *Projectlist {
	it := &Projectlist{}
	it.CollectInitialize(rule, "project", 0,"project")
	it.ItExecute = &dealListExecute{it}
	it.ItShowGrid = &dealShowGrid{it}
	it.ItShowPage = &dealShowPage{it}
	it.IsLockRemote = true
	return it
}

//	Обработчик событий
func (it *Projectlist) ListExecute(rule *repo.DataRule, cmd string, data lik.Seter) {
	if cmd == "project" {
		it.ListExecute(rule, rule.Shift(), data)
	} else {
		it.BoxExecute(rule, cmd, data)
	}
}

//	Заполнение таблицы
func (it *Projectlist) ListShowGrid(rule *repo.DataRule) {
	//it.Parameters.SetItem(64, "cellHeight")
	it.SetSize(700, 500)
	//it.SetParameter(48, "cellHeight")
	it.Title = "История обновлений"
	it.Titles.AddItemSet("text=Закрыть", "handler=function_fancy_form_cancel")
	it.AddColumnItem(rule, "type=text", "index=date", "title=Дата", "width=100", "autoHeight=true")
	it.AddColumnItem(rule, "type=text", "index=version", "title=Вер.", "width=50", "autoHeight=true")
	it.AddColumnItem(rule, "type=text", "index=what", "title=Что сделано", "width=545", "autoHeight=true")
	it.CollectShowGrid(rule)
}

//	Заполнение страницы
func (it *Projectlist) ListShowPage(rule *repo.DataRule) lik.Lister {
	rows := it.CollectShowPage(rule)
	for nm := 0; nm < len(one.LUpDt); nm++ {
		up := one.LUpDt[nm]
		row := lik.BuildSet("id", 1+nm)
		row.SetItem(up.Date, "date")
		row.SetItem(up.Ver, "version")
		row.SetItem(up.What, "what")
		rows.AddItems(row)
	}
	return rows
}

