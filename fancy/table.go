package fancy

import (
	"github.com/massarakhsh/tsan/repo"
	"github.com/massarakhsh/tsan/show"
	"github.com/massarakhsh/lik"
	"github.com/massarakhsh/lik/likdom"
	"fmt"
	"strings"
)

//	Предопределенные номера позиций
const (
	OrdFilter = 100
	OrdColumn = 200
	OrdSegment = 300
	OrdRealty = 400
	OrdLocate = 500
	OrdStatus = 600
	OrdSearch = 700
	OrdUnit = 800
	OrdPhone = 900
)

//	Дескриптор таблицы
type TableFancy struct {
	DataFancy						//	Ядро
	ItGridFill  DealTableGridFill	//	Интерфейс таблицы
	ItPageFill  DealTablePageFill	//	Интерфейс страницы
	ItFormFill  DealTableFormFill	//	Интерфейс формы
	//ItSelectRow DealTableSelectRow	//	Выбрать строку
	//ItEnterRow  DealTableEnterRow	//	Войти в строку
	IsLockRemote	bool			//	Блокировка удаленного управления
}

//	Интерфейс заполнения таблицы
type DealTableGridFill interface {
	Run(rule *repo.DataRule)
}
type dealTableGridFill struct {
	It	*TableFancy
}
func (it *dealTableGridFill) Run(rule *repo.DataRule) {
	it.It.TableGridFill(rule)
}

//	Интерфейс заполнения строницы
type DealTablePageFill interface {
	Run(rule *repo.DataRule) lik.Lister
}
type dealTablePageFill struct {
	It	*TableFancy
}
func (it *dealTablePageFill) Run(rule *repo.DataRule) lik.Lister {
	return it.It.TablePageFill(rule)
}

//	Интерфейс заполнения формы
type DealTableFormFill interface {
	Run(rule *repo.DataRule)
}
type dealTableFormFill struct {
	It	*TableFancy
}
func (it *dealTableFormFill) Run(rule *repo.DataRule) {
	it.It.TableFormFill(rule)
}

//	Инициализация таблицы
func (it *TableFancy) TableInitialize(rule *repo.DataRule, main string, part string, zone string) {
	it.FancyInitialize(main, part, zone)
	it.ItGridFill = &dealTableGridFill{it}
	it.ItPageFill = &dealTablePageFill{it}
	it.ItFormFill = &dealTableFormFill{it}
}

//	Выполнения команд
func (it *TableFancy) TableExecute(rule *repo.DataRule, cmd string, data lik.Seter) {
	if cmd == "showgrid" {
		it.GridSync.Lock()
		it.cmdShowGrid(rule)
		it.GridSync.Unlock()
	} else if cmd == "showpage" {
		it.cmdShowPage(rule)
	} else if cmd == "showform" {
		parm := rule.Shift()
		if strings.HasPrefix(parm, "_") {
			parm = parm[1:]
		} else {
			parm = parm
		}
		if parm == "show" {
			it.Fun = FunShow
		} else if parm == "mod" {
			it.Fun = FunMod
		} else if parm == "add" {
			it.Fun = FunAdd
		} else if parm == "edit" {
			it.Fun = FunEdit
		} else if parm == "del" {
			it.Fun = FunDel
		} else {
			it.Fun = FunNo
		}
		if it.Fun != FunNo {
			it.cmdShowForm(rule)
		}
	} else if cmd == "tab" {
		it.FormFixTab(rule.Shift())
		rule.SetResponse(it.Form.Tab,"tab")
	} else if cmd == "cancel" {
		//it.Sel = ""
		it.Fun = FunNo
	} else if cmd == "rowselect" {
		it.TableSelectRow(rule, rule.Shift())
	} else if lik.RegExCompare(cmd,"col(size|show|drag)") {
		index := rule.Shift()
		prm := lik.StrToInt(rule.Shift())
		if !it.ChangeColumn(rule, cmd, index, prm) {
			rule.OnChangeData()
		}
	} else if match := lik.RegExParse(cmd,"^to(show|add|mod|edit|enter|delete)"); match != nil {
		mode := match[1]
		if mode == "enter" {
			mode = "show"
		}
		if top := rule.Top(); top != "" {
			it.TableSelectRow(rule, top)
		}
		parm := fmt.Sprintf("%s_%s_%s", it.Main, it.Zone, mode)
		rule.SetResponse(parm, "_function_fancy_trio_form")
	}
}

//	Отображение страницы
func (it *TableFancy) BuildShow(rule *repo.DataRule, sx int, sy int) likdom.Domer {
	it.SetSize(sx, sy)
	return show.BuildFancyGrid(it.Main, it.Zone)
}

//	Отображение таблицы
func (it *TableFancy) cmdShowGrid(rule *repo.DataRule) {
	it.GridClear()
	it.RunFieldsFill(rule)
	it.RunGridFill(rule)
	if it.IsLockRemote && !it.Grid.Parameters.IsItem("data/items") {
		items := it.RunPageFill(rule)
		it.Grid.SetParameter(items, "data/items")
	}
	it.GridShow(rule)
}

//	Отображение страницы
func (it *TableFancy) cmdShowPage(rule *repo.DataRule) {
	rows := it.RunPageFill(rule)
	it.GridPage(rule, rows)
}

//	Отображение формы
func (it *TableFancy) cmdShowForm(rule *repo.DataRule) {
	it.FormClear()
	it.RunFormFill(rule)
	if it.Form.Tools.Count() == 0 {
		it.AddTitleToolText(rule, "Закрыть", "function_fancy_form_cancel")
	}
	if it.Form.Items.Count() == 0 {
		it.Form.Items.AddItemSet("type=html", "value=Нет полей")
	}
	it.ShowForm(rule)
}

//	Отобразить форму
func (it *TableFancy) RunGridFill(rule *repo.DataRule) {
	if it.ItGridFill != nil {
		it.ItGridFill.Run(rule)
	}
}
func (it *TableFancy) TableGridFill(rule *repo.DataRule) {
	it.Grid.SetSize(it.Sx, it.Sy)
	if !it.IsLockRemote {
		it.Grid.SetParameter(true, "data/remoteSort")
		it.Grid.SetParameter(true, "data/remoteFilter")
		it.Grid.SetParameter(true, "data/remotePage")
		url := rule.BuildUrl("/front/" + it.Main + "/" + it.Zone + "/showpage")
		it.Grid.SetParameter(url, "data/proxy/url")
		it.Grid.SetParameter("rest", "data/proxy/type")
		it.Grid.SetParameter("POST", "data/proxy/methods/read")
		it.Grid.SetParameter(it.Part, "data/proxy/params/part")
		it.Grid.SetParameter("items", "data/proxy/reader/root")
		it.Grid.SetParameter("function_fancy_before_request", "data/proxy/beforeRequest")
		it.Grid.SetParameter("function_fancy_after_request", "data/proxy/afterRequest")
	}
	//it.Grid.SetParameter(30,"paging", "pageSize")
}

//	Отобразить страницу
func (it *TableFancy) RunPageFill(rule *repo.DataRule) lik.Lister {
	if it.ItPageFill != nil {
		return it.ItPageFill.Run(rule)
	} else {
		return lik.BuildList()
	}
}
func (it *TableFancy) TablePageFill(rule *repo.DataRule) lik.Lister {
	return lik.BuildList()
}

func (it *TableFancy) RunFormFill(rule *repo.DataRule) {
	if it.ItFormFill != nil {
		it.ItFormFill.Run(rule)
	}
}
func (it *TableFancy) TableFormFill(rule *repo.DataRule) {
}

func (it *TableFancy) TableSelectRow(rule *repo.DataRule, sel string) {
	if sel != it.Sel {
		it.Sel = sel
	}
}

