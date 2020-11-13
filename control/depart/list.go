//	Модуль списка подразделений
package depart

import (
	"github.com/massarakhsh/tsan/control"
	"github.com/massarakhsh/tsan/fancy"
	"github.com/massarakhsh/tsan/jone"
	"github.com/massarakhsh/tsan/repo"
	"github.com/massarakhsh/lik"
	"github.com/massarakhsh/lik/likbase"
)

//	Дескриптор списка подразделений
type DepartControl struct {
	control.ListControl
}

//	Интерфейс создания таблицы
type dealDepartGridFill struct {
	It	*DepartControl
}
func (it *dealDepartGridFill) Run(rule *repo.DataRule) {
	it.It.DepartGridFill(rule)
}

//	Интерфейс создания страницы
type dealDepartPageFill struct {
	It	*DepartControl
}
func (it *dealDepartPageFill) Run(rule *repo.DataRule) lik.Lister {
	return it.It.DepartPageFill(rule)
}

//	Интерфейс создания формы
type dealDepartElmForm struct {
	It	*DepartControl
}
func (it *dealDepartElmForm) Run(rule *repo.DataRule, elm *likbase.ItElm) {
	it.It.DepartElmForm(rule, elm)
}

//	Интерфейс команд
type dealDepartExecute struct {
	It	*DepartControl
}
func (it *dealDepartExecute) Run(rule *repo.DataRule, cmd string, data lik.Seter) {
	it.It.DepartExecute(rule, cmd, data)
}

//	Кнструктор дескриптора
func BuildListDepart(rule *repo.DataRule, id lik.IDB) *DepartControl {
	it := &DepartControl{}
	it.ControlInitialize("depart", id)
	it.ListInitialize(rule, "depart", "depart")
	it.ItGridFill = &dealDepartGridFill{it}
	it.ItPageFill = &dealDepartPageFill{it}
	it.ItFormElm = &dealDepartElmForm{it}
	it.ItExecute = &dealDepartExecute{it}
	it.IsLockRemote = true
	return it
}

//	Выполнение команд
func (it *DepartControl) DepartExecute(rule *repo.DataRule, cmd string, data lik.Seter) {
	if cmd == "all" || cmd == it.Main {
		it.DepartExecute(rule, rule.Shift(), data)
	} else {
		it.ListExecute(rule, cmd, data)
	}
}

//	Создание таблицы
func (it *DepartControl) DepartGridFill(rule *repo.DataRule) {
	it.ListGridFill(rule)
	if rule.IAmAdmin() {
		it.AddCommandImg(rule, fancy.OrdUnit+1, "Создать запись подразделения", "toadd", "add")
	}
}

//	Создание страницы
func (it *DepartControl) DepartPageFill(rule *repo.DataRule) lik.Lister {
	rows := it.TablePageFill(rule)
	elms := []*likbase.ItElm{}
	for _, elm := range jone.TableDepart.Elms {
		name := elm.GetString("name")
		used := false
		list := []*likbase.ItElm{}
		var pos int
		for pos = 0; pos < len(elms); pos++ {
			if !used && name < elms[pos].GetString("name") {
				list = append(list, elm)
				used = true
			}
			list = append(list, elms[pos])
		}
		if !used {
			list = append(list, elm)
			used = true
		}
		elms = list
	}
	nodes := make(map[lik.IDB]lik.Seter)
	doit := true
	for doit {
		doit = false
		for ne := 0; ne < len(elms); ne++ {
			if elm := elms[ne]; elm != nil {
				if idup := elm.GetIDB("departid"); idup == 0 {
					row := it.GridElmRow(rule, elm)
					row.SetItem(true, "leaf")
					rows.AddItems(row)
					nodes[elm.Id] = row
					elms[ne] = nil
					doit = true
				} else if rowup,_ := nodes[idup]; rowup != nil {
					row := it.GridElmRow(rule, elm)
					row.SetItem(true, "leaf")
					child := rowup.GetList("child")
					if child == nil {
						child = lik.BuildList()
						rowup.SetItem(child, "child")
						rowup.SetItem(false, "leaf")
						rowup.SetItem(true, "expanded")
					}
					child.AddItems(row)
					nodes[elm.Id] = row
					elms[ne] = nil
					doit = true
				}
			}
		}
	}
	for ne := 0; ne < len(elms); ne++ {
		if elm := elms[ne]; elm != nil {
			row := it.GridElmRow(rule, elm)
			row.SetItem(true, "leaf")
			rows.AddItems(row)
		}
	}
	return rows
}

//	Создание формы
func (it *DepartControl) DepartElmForm(rule *repo.DataRule, elm *likbase.ItElm) {
	it.ListFormElm(rule, elm)
}

//	Проверка строки таблицы
func (it *DepartControl) DepartListRowsProbe(rule *repo.DataRule, elm *likbase.ItElm) bool {
	accept := (elm != nil)
	if accept && !it.ListMakeProbe(rule, elm) {
		accept = false
	}
	return accept
}

