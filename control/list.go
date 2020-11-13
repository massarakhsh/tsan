// Контроллеры окна списков объектов.
//
// Контроллер списка
package control

import (
	"github.com/massarakhsh/tsan/fancy"
	"github.com/massarakhsh/tsan/jone"
	"github.com/massarakhsh/tsan/one"
	"github.com/massarakhsh/tsan/repo"
	"github.com/massarakhsh/lik"
	"github.com/massarakhsh/lik/likbase"
	"fmt"
	"strings"
	"time"
)

//	Объект контроллера списков объектов базы данных
type ListControl struct {
	DataControl                         //	На общем контроллере
	fancy.TableFancy                    //	Объект таблицы
	IsUpList         bool               //	Это - список высшего уровня
	ItListMakeSelect DealListMakeSelect //	Интерфейс выбора строк
	ItListMakeProbe  DealListMakeProbe  //	Интерфейс проверки строк
	ItListMakeSort   DealListMakeSort   //	Интерфейс сортировки
	ItListMakePage   DealListMakePage   //	Интерфейс отбора страницы
	ItListMakeBuild  DealListMakeBuild  //	Интерфейс зааполнения списка
	ItRowFill        DealRowFill        //	Интерфейс построения строки
	ItFormElm        DealFormElm        //	Интерфейс построения формы
	ItImport         DealListImport     //	Интерфейс импорта объекта
	Selector         fancy.DataFilter   //	Дескриптор селектора
	CmdSignDo        string             //	Команда индикации
	CmdGridDo        string             //	Команда списка
	TimeGridDo       time.Time          //	Время обновления списка
	CmdFormDo        string             //	Команда формы
	TimeFormDo       time.Time          //	Время обновления формы
}

//	Интерфейс заполнения списка
type dealListGridFill struct {
	It	*ListControl
}
func (it *dealListGridFill) Run(rule *repo.DataRule) {
	it.It.ListGridFill(rule)
}

//	Истерфейс заполнения страницы
type dealListPageFill struct {
	It	*ListControl
}
func (it *dealListPageFill) Run(rule *repo.DataRule) lik.Lister {
	return it.It.ListPageFill(rule)
}

//	Интерфейс заполнения формы
type DealListFormFill interface {
	Run(rule *repo.DataRule, elm *likbase.ItElm)
}
type dealListFormFill struct {
	It	*ListControl
}
func (it *dealListFormFill) Run(rule *repo.DataRule) {
	it.It.ListFormFill(rule)
}

//	Интерфейс команд
type dealExecute struct {
	It	*ListControl
}
func (it *dealExecute) Run(rule *repo.DataRule, cmd string, data lik.Seter) {
	it.It.ListExecute(rule, cmd, data)
}

//	Интерфейс маршализации
type dealMarshal struct {
	It	*ListControl
}
func (it *dealMarshal) Run(rule *repo.DataRule) {
	it.It.ListMarshal(rule)
}

//	Интерфейс отбора строк
type DealListMakeSelect interface {
	Run(rule *repo.DataRule) []*likbase.ItElm
}
type dealListMakeSelect struct {
	It	*ListControl
}
func (it *dealListMakeSelect) Run(rule *repo.DataRule) []*likbase.ItElm {
	return it.It.ListMakeSelect(rule)
}

//	Интерфейс проверки строк
type DealListMakeProbe interface {
	Run(rule *repo.DataRule, elm *likbase.ItElm) bool
}
type dealListMakeProbe struct {
	It	*ListControl
}
func (it *dealListMakeProbe) Run(rule *repo.DataRule, elm *likbase.ItElm) bool {
	return it.It.ListMakeProbe(rule, elm)
}

//	Интерфейс сортировки строк
type DealListMakeSort interface {
	Run(rule *repo.DataRule, list []*likbase.ItElm) []*likbase.ItElm
}
type dealListMakeSort struct {
	It	*ListControl
}
func (it *dealListMakeSort) Run(rule *repo.DataRule, list []*likbase.ItElm) []*likbase.ItElm {
	return it.It.ListMakeSort(rule, list)
}

//	Интерфейс отбора страницы
type DealListMakePage interface {
	Run(rule *repo.DataRule, list []*likbase.ItElm) []*likbase.ItElm
}
type dealListMakePage struct {
	It	*ListControl
}
func (it *dealListMakePage) Run(rule *repo.DataRule, list []*likbase.ItElm) []*likbase.ItElm {
	return it.It.ListRowsPage(rule, list)
}

//	Интерфейс заполнения строк
type DealListMakeBuild interface {
	Run(rule *repo.DataRule, list []*likbase.ItElm) lik.Lister
}
type dealListRowsBuild struct {
	It	*ListControl
}
func (it *dealListRowsBuild) Run(rule *repo.DataRule, list []*likbase.ItElm) lik.Lister {
	return it.It.ListMakeBuild(rule, list)
}

//	Интерфейс заполнения строки
type DealRowFill interface {
	Run(rule *repo.DataRule, elm *likbase.ItElm, row lik.Seter)
}
type dealRowFill struct {
	It	*ListControl
}
func (it *dealRowFill) Run(rule *repo.DataRule, elm *likbase.ItElm, row lik.Seter) {
	it.It.ListRowFill(rule, elm, row)
}

//	Интерфейс заполнения формы
type DealFormElm interface {
	Run(rule *repo.DataRule, elm *likbase.ItElm)
}
type dealFormElm struct {
	It	*ListControl
}
func (it *dealFormElm) Run(rule *repo.DataRule, elm *likbase.ItElm) {
	it.It.ListFormElm(rule, elm)
}

//	Интерфейс импорта объекта
type DealListImport interface {
	Run(rule *repo.DataRule, part string, sid string)
}
type dealListImport struct {
	It	*ListControl
}
func (it *dealListImport) Run(rule *repo.DataRule, part string, sid string) {
}

//	Инициализация списка
func (it *ListControl) ListInitialize(rule *repo.DataRule, mode string, part string) {
	it.Stable = true
	it.TableInitialize(rule, mode, part,"all")
	it.ItGridFill = &dealListGridFill{it}
	it.ItPageFill = &dealListPageFill{it}
	it.ItFormFill = &dealListFormFill{it}
	it.ItExecute = &dealExecute{it}
	it.ItMarshal = &dealMarshal{it}
	it.ItListMakeSelect = &dealListMakeSelect{it}
	it.ItListMakeProbe = &dealListMakeProbe{it}
	it.ItListMakeSort = &dealListMakeSort{it}
	it.ItListMakePage = &dealListMakePage{it}
	it.ItListMakeBuild = &dealListRowsBuild{it}
	it.ItRowFill = &dealRowFill{it}
	it.ItFormElm = &dealFormElm{it}
	it.ItImport = &dealListImport{it}
}

//	Маршализация
func (it *ListControl) ListMarshal(rule *repo.DataRule) {
	now := time.Now()
	if rule.PathCommand != "" {
		it.CmdFormDo = rule.PathCommand
		rule.PathCommand = ""
	}
	if it.CmdGridDo != "" && now.Sub(it.TimeGridDo) > 3 * time.Second {
		mode := it.CmdGridDo
		it.TimeGridDo = now
		rule.SetResponse(mode, "_function_fancy_grid_update")
	}
	if it.CmdFormDo != "" && now.Sub(it.TimeFormDo) > 3 * time.Second {
		mode := it.CmdFormDo
		if match := lik.RegExParse(it.CmdFormDo, "(.*)_(.*)"); match != nil {
			mode = match[1]
			it.Sel = match[2]
		}
		rule.SetResponse(it.Main+"_all_"+mode, "_function_fancy_trio_form")
	}
}

//	Обработка команд
func (it *ListControl) ListExecute(rule *repo.DataRule, cmd string, data lik.Seter) {
	if cmd == "all" || cmd == it.Zone {
		it.ListExecute(rule, rule.Shift(), data)
	} else if cmd == "select" {
		it.ItImport.Run(rule, it.Part, rule.Shift())
		it.FormFixTab(rule.Shift())
	} else if cmd == "filter" {
		nmf := rule.Shift()
		it.FancyFilterSet(rule, nmf, true)
		rule.OnChangeData()
	} else if cmd == "filtersave" {
		it.FancyFilterSave(rule)
		rule.OnChangeData()
	} else if cmd == "filtersaveas" {
		name := lik.StringFromXS(rule.Shift())
		it.FancyFilterSaveAs(rule, name)
		rule.OnChangeData()
	} else if cmd == "filterdelete" {
		it.FancyFilterDelete(rule)
		rule.OnChangeData()
	} else if lik.RegExCompare(cmd,"(segment|realty|locate|status)") {
		rule.SetMemberParam(rule.Shift(),"context/"+cmd)
		rule.OnChangeData()
	} else if cmd == "mark" {
		id := likbase.StrToIDB(rule.Shift())
		val := lik.StrToInt(rule.Shift())
		repo.MarkElmSet(rule, it.Part, id, val > 0)
	} else if cmd == "write" {
		if elm := it.ListCmdWrite(rule, data); elm != nil {
			rule.SetGoPart(fmt.Sprintf("/%s%d", it.Main, int(elm.Id)))
		}
	} else if cmd == "delete" {
		it.cmdElmDelete(rule)
	/*} else if cmd == "itshow" {
		if sid := rule.Shift(); sid != "" {
			it.RunSelectRow(rule, sid)
			it.CmdFormDo = fmt.Sprintf("%s_%s", "show", sid)
		}
	} else if cmd == "itedit" {
		if sid := rule.Shift(); sid != "" {
			it.RunSelectRow(rule, sid)
			it.CmdFormDo = fmt.Sprintf("%s_%s", "edit", sid)
		}*/
	} else {
		if cmd == "showform" && it.CmdFormDo != "" {
			it.CmdFormDo = ""
		}
		it.TableExecute(rule, cmd, data)
	}
}

//	Изготовление таблицы
func (it *ListControl) ListGridFill(rule *repo.DataRule) {
	it.TableGridFill(rule)
	it.GridBuildColumns(rule, true, true)
	it.GridTuneFilter(rule)
	it.Grid.SetParameter(true,"likFull")
	it.Grid.SetParameter(false, "data/proxy/autoLoad")
	it.Grid.SetParameter(true,"defaults/filter/header")
	it.ListGridCommandFilter(rule)
	it.ListGridCommandLocate(rule)
	it.ListGridCommandUnit(rule)
	it.ListGridCommandColumn(rule)
	it.ListGridCommandStatus(rule)
	it.Grid.Columns.InsertItem(lik.BuildSet("type=checkbox", "width=44", "index=mark",
		"title=Выбор", "cellTip=Отметить", "editable=true"), 1)
	it.Grid.AddEventAction("cellclick", "function_fancy_grid_mark")
	it.AddContextMenu("Открыть", "...", "toshow")
	it.AddContextMenu("Изменить", "...", "toedit")
	it.AddContextMenu("", "", "")
	it.AddContextMenu("Поставить/убрать отметку", "", "mark")
}

//	Кнопка и меню фильтров
func (it *ListControl) ListGridCommandFilter(rule *repo.DataRule) {
	data := rule.GetMemberParamList(it.GetParter()+"/filters")
	it.ListGridAddCombo(rule, fancy.OrdFilter, it.GetParter(),"filter", "* Фильтры","Всё", data)
	my := it.FancyFilterSeek(rule) >= 0
	data = lik.BuildList()
	data.AddItems(lik.BuildSet("type=button", "text=Сохранить фильтр", "handler=function_filter_save",
		"disabled", !my))
	data.AddItems(lik.BuildSet("type=button", "text=Сохранить как", "items", lik.BuildList(
		lik.BuildSet("type=string", "events", lik.BuildList(lik.BuildSet("input=function_symbol_value"))),
		lik.BuildSet("type=button", "text=Сохранить", "handler=function_filter_saveas"),
	)))
	data.AddItems(lik.BuildSet("type=button", "text=Удалить фильтр", "handler=function_filter_delete",
		"disabled", !my))
	it.AddCommandItem(rule, fancy.OrdFilter + 10, lik.BuildSet(
		"type=button", "tip=Управление фильтрами", "id=topfilters",
		"imageCls", "imgsel", "menu", data,
	))
	it.AddCommandItem(rule, fancy.OrdFilter + 20, lik.BuildSet())
}

//	Кнопка и меню колонок
func (it *ListControl) ListGridCommandColumn(rule *repo.DataRule) {
	data := lik.BuildList()
	for _,field := range it.ListFields {
		part := strings.ReplaceAll(field.GetString("part"), "/", "__")
		cls := "imgplus"
		if (field.GetInt("tags") & jone.TagHide) != 0 {
			cls = "imgno"
		}
		data.AddItems(lik.BuildSet("type=button",
			"imageCls", cls,
			"handler", fmt.Sprintf("function_fancy_column_show(%s)", part),
			"text", field.GetString("name")))
	}
	it.AddCommandItem(rule, fancy.OrdColumn, lik.BuildSet(
		"type=button", "tip=Выбор столбцов", "id=topcolumns",
		"imageCls", "imgset", "menu", data,
	))
	it.AddCommandItem(rule, fancy.OrdColumn+10, lik.BuildSet())
}

//	Кнопка и меню сегментов рынка
func (it *ListControl) ListGridCommandSegment(rule *repo.DataRule) {
	if segment := rule.SeekSegment(""); segment == "" || segment == jone.DoCall {
		list := lik.BuildList();
		if rule.ICanDo(jone.DoSecond) {
			list.AddItemSet("part", jone.DoSecond, "name=Вторичка")
		}
		if rule.ICanDo(jone.DoNew) {
			list.AddItemSet("part", jone.DoNew, "name=Новостройки")
		}
		if rule.ICanDo(jone.DoVilla) {
			list.AddItemSet("part", jone.DoVilla, "name=Загород")
		}
		if rule.ICanDo(jone.DoArea) {
			list.AddItemSet("part", jone.DoArea, "name=Участки")
		}
		if rule.ICanDo(jone.DoRent) {
			list.AddItemSet("part", jone.DoRent, "name=Аренда")
		}
		it.ListGridAddCombo(rule, fancy.OrdSegment, "context", "segment", "Сегмент", "Всё", list)
	} else {
		rule.SetMemberParam(segment, "context/segment")
	}
}

//	Кнопка и меню видов недвижимости
func (it *ListControl) ListGridCommandRealty(rule *repo.DataRule) {
	var list lik.Lister
	if ent := repo.GenDiction.FindEnt("realty"); ent != nil {
		segment := rule.GetMemberParamString("context/segment")
		if list = ent.It.GetList("content"); list != nil {
			for nc := list.Count()-1; nc >= 0; nc-- {
				part := list.GetSet(nc).GetString("part")
				if segment == "second" && part != "flat" && part != "room" {
					list.DelItem(nc)
				} else if segment == "new" && part != "flat" && part != "room" {
					list.DelItem(nc)
				}
			}
		}
	}
	it.ListGridAddCombo(rule, fancy.OrdRealty,"context","realty","Тип недв.","Всё", list)
}

//	Кнопка и меню локализации
func (it *ListControl) ListGridCommandLocate(rule *repo.DataRule) {
	list := lik.BuildList(
		lik.BuildSet("part", jone.ItMy, "name=Мои"),
		lik.BuildSet("part", jone.ItDep, "name=Отдел"),
		lik.BuildSet("part", jone.ItAll, "name=Все"),
	)
	it.ListGridAddCombo(rule, fancy.OrdLocate,"context","locate","Чьи","",list)
}

//	Кнопка и меню статуса
func (it *ListControl) ListGridCommandStatus(rule *repo.DataRule) {
	list := lik.BuildList()
	list.AddItemSet("part", jone.ItActive, "name=Активные")
	list.AddItemSet("part", jone.ItAll, "name=Все")
	it.ListGridAddCombo(rule, fancy.OrdStatus,"context","status","Статус","",list)
}

//	Кнопка и меню поиска
func (it *ListControl) ListGridCommandSearch(rule *repo.DataRule) {
	it.AddCommandItem(rule, fancy.OrdSearch- 1, lik.BuildSet())
	it.AddCommandItem(rule, fancy.OrdSearch, lik.BuildSet(
		"type=search",
		"width=150", "minListWidth=220",
		"emptyText=Найти",
		"paramsMenu=true",
		"paramsText=Где",
	))
}

//	Дополнительные кнопки
func (it *ListControl) ListGridCommandUnit(rule *repo.DataRule) {
	it.AddCommandItem(rule, fancy.OrdUnit- 1, lik.BuildSet())
	what := ""
	if it.Part == "bell" {
		what = "контакт"
	} else if it.Part == "object" {
		what = "объект"
	} else if it.Part == "offet" {
		what = "заявку"
	} else if it.Part == "deal" {
		what = "сделку"
	} else if it.Part == "member" {
		what = "сотрудника"
	} else if it.Part == "depart" {
		what = "подразделение"
	}
	it.AddCommandImg(rule, fancy.OrdUnit, "Открыть "+what, "toshow", "show")
	if it.IsUpList {
		it.AddCommandImg(rule, fancy.OrdUnit+ 1, "Создать "+what, "toadd", "add")
	}
}

//	Доавление кнопки на панель
func (it *ListControl) ListGridAddCombo(rule *repo.DataRule, ord int, part string, key string, title string, def string, list lik.Lister) {
	param := rule.GetMemberParamString(part+"/"+key)
	if param == "" {
		if key == "status" {
			param = jone.ItActive
		} else {
			param = jone.ItAll
		}
		rule.SetMemberParam(param, part+"/"+key)
	}
	data := lik.BuildList()
	if def != "" {
		data.AddItemSet("part=all", "name=Все")
	}
	if list != nil {
		for n := 0; n < list.Count(); n++ {
			if elm := list.GetSet(n); elm != nil {
				part := elm.GetString("part")
				name := elm.GetString("name")
				if name == "" {
					name = "(без имени)"
				}
				item := lik.BuildSet()
				item.SetItem(part, "part")
				item.SetItem(name, "name")
				data.AddItems(item)
			}
		}
	}
	if title != "" && !strings.HasPrefix(title,"*") {
		it.AddCommandItem(rule, ord, lik.BuildSet("type=text", "text", title+":", "cls=small"))
	}
	it.AddCommandItem(rule, ord+1, lik.BuildSet(
		"type=combo",
		"width=120", "minListWidth=120",
		"valueKey=part", "displayKey=name",
		"data", data,
		"value", param,
		"tip", title,
		"events", lik.BuildList(lik.BuildSet("change", "function_grid_command_" + key)),
	))
}

//	Изготовление страницы
func (it *ListControl) ListPageFill(rule *repo.DataRule) lik.Lister {
	it.CmdGridDo = ""
	params := rule.GetAllContext()
	it.FancyFilterFix(rule, params)
	return it.ListMakeRows(rule)
}

//	Изготовление строк
func (it *ListControl) ListMakeRows(rule *repo.DataRule) lik.Lister {
	it.Selector = it.FancyFilterDecode(rule)
	list := it.RunListMakeSelect(rule)
	list = it.RunListMakeSort(rule, list)
	list = it.RunListMakePage(rule, list)
	rows := it.RunListMakeBuild(rule, list)
	return rows
}
func printdiff(txt string, t1,t2 time.Time) {
	fmt.Printf("%s %4.2f\r\n", txt, float64(t2.Sub(t1)/time.Millisecond) / 1000.0)
}

//	Отбор строк
func (it *ListControl) RunListMakeSelect(rule *repo.DataRule) []*likbase.ItElm {
	if it.ItListMakeSelect == nil { return nil }
	return it.ItListMakeSelect.Run(rule)
}
func (it *ListControl) ListMakeSelect(rule *repo.DataRule) []*likbase.ItElm {
	var list []*likbase.ItElm
	elms := jone.GetTable(it.Part).Elms
	for _, elm := range elms {
		if it.RunListMakeProbe(rule, elm) {
			list = append(list, elm)
		}
	}
	return likbase.SortById(list)
}

//	Проверка строк
func (it *ListControl) RunListMakeProbe(rule *repo.DataRule, elm *likbase.ItElm) bool {
	if it.ItListMakeProbe == nil { return false }
	return it.ItListMakeProbe.Run(rule, elm)
}
func (it *ListControl) ListMakeProbe(rule *repo.DataRule, elm *likbase.ItElm) bool {
	accept := true
	if search := it.Selector.Search; search != "" {
		/*accept := false
		for _, ch := range strings.Split(search, ",") {
			if ch != "" {
				for _,set := range row.Values() {
					if key := set.Key; key != "id" {
						if strings.Contains(set.Val.ToString(), ch) {
							accept = true
							break
						}
					}
				}
			}
			if accept { break }
		}*/
	}
	if accept {
		for _, cond := range it.Selector.Conds {
			ok := false
			if value := it.RunCalculate(rule, elm.Info, cond.Key, "", false); value != nil {
				valstr := value.ToString()
				if match := lik.RegExParse(cond.Key, "([^/]+)id$"); match != nil {
					valstr = jone.CalculatePartIdText(match[1], lik.IDB(value.ToInt()))
				} else if match := lik.RegExParse(cond.Key, "([^/]+)$"); match != nil {
					valstr = jone.SystemStringTranslate(match[1], valstr)
				}
				valstr = strings.ToLower(valstr)
				if cond.Opr == "|" {
					val := strings.Replace(cond.Val, " ", "", -1)
					for _, ch := range strings.Split(val, ",") {
						if len(ch) > 0 {
							if valstr == ch {
								ok = true
								break
							}
						}
					}
				} else if cond.Opr == "=" && valstr == cond.Val {
					ok = true
				} else if cond.Opr == "!=" && valstr != cond.Val {
					ok = true
				} else if cond.Opr == "*" && strings.Contains(valstr, cond.Val) {
					ok = true
				} else if fval,isf := lik.StrToFloatIf(cond.Val); isf {
					if cond.Opr == "=" && value.ToFloat() == fval {
						ok = true
					} else if cond.Opr == "!=" && value.ToFloat() != fval {
						ok = true
					} else if cond.Opr == ">" && value.ToFloat() > fval {
						ok = true
					} else if cond.Opr == ">=" && value.ToFloat() >= fval {
						ok = true
					} else if cond.Opr == "<" && value.ToFloat() < fval {
						ok = true
					} else if cond.Opr == "<=" && value.ToFloat() <= fval {
						ok = true
					}
				} else if cond.Opr == ">" && valstr > cond.Val {
					ok = true
				} else if cond.Opr == ">=" && valstr >= cond.Val {
					ok = true
				} else if cond.Opr == "<" && valstr < cond.Val {
					ok = true
				} else if cond.Opr == "<=" && valstr <= cond.Val {
					ok = true
				} else {
					ok = strings.Contains(valstr, cond.Val)
				}
			}
			if !ok {
				accept = false
				break
			}
		}
	}
	return accept
}

//	Сортировка строк
func (it *ListControl) RunListMakeSort(rule *repo.DataRule, list []*likbase.ItElm) []*likbase.ItElm {
	if it.ItListMakeSort == nil { return list }
	return it.ItListMakeSort.Run(rule, list)
}

//	Дескриптор элемента сортировки
type listEnt struct {
	Key		string
	Elm		*likbase.ItElm
}

//	Стандартная сортировка
func (it *ListControl) ListMakeSort(rule *repo.DataRule, list []*likbase.ItElm) []*likbase.ItElm {
	collect := list
	if list != nil {
		sort,dir := "",true
		if it.Selector.SortKey != "" {
			sort = strings.ReplaceAll(it.Selector.SortKey, "__", "/")
			dir = it.Selector.SortDir
		}
		selist := []*listEnt{}
		for _,elm := range list {
			selm := &listEnt{ Elm: elm }
			if sort != "" {
				if value := it.RunCalculate(rule, elm.Info, sort, "", false); value != nil {
					selm.Key = value.ToString()
					if match := lik.RegExParse(sort,"([^/]+)$"); match != nil {
						selm.Key = jone.SystemStringTranslate(match[1], selm.Key)
					}
				}
			}
			selist = append(selist, selm)
		}
		selist = it.ListMakeSortField(selist, sort, dir)
		collect = []*likbase.ItElm{}
		for _, ent := range selist {
			collect = append(collect, ent.Elm)
		}
	}
	return collect
}

//	Сортировка по полю
func (it *ListControl) ListMakeSortField(list []*listEnt, sortkey string, sortup bool) []*listEnt {
	format := ""
	if field := it.SeekPartField(sortkey); sortkey != "" && field != nil {
		format = field.GetString("format")
	}
	size := len(list)
	for doze := 1; doze < size; doze *= 2 {
		trans := []*listEnt{}
		for pos := 0; pos < size; pos += doze * 2 {
			posa := pos
			enda := pos + doze
			if enda > size { enda = size }
			posb := enda
			endb := posb + doze
			if endb > size { endb = size }
			for posa < enda || posb < endb {
				first := true
				if posa >= enda {
					first = false
				} else if posb >= endb {
					first = true
				} else {
					keya := list[posa].Key
					keyb := list[posb].Key
					if format == "n" || format == "m" {
						keya = fancy.StringateNumber(keya)
						keyb = fancy.StringateNumber(keyb)
					}
					if keya < keyb {
						first = sortup
					} else if keya > keyb {
						first = !sortup
					} else if list[posa].Elm.Id >= list[posb].Elm.Id {
						first = sortup
					} else {
						first = !sortup
					}
				}
				if first {
					trans = append(trans, list[posa])
					posa++
				} else {
					trans = append(trans, list[posb])
					posb++
				}
			}
		}
		list = trans
	}
	return list
}

//	Отбор страницы
func (it *ListControl) RunListMakePage(rule *repo.DataRule, list []*likbase.ItElm) []*likbase.ItElm {
	if it.ItListMakePage == nil { return list }
	return it.ItListMakePage.Run(rule, list)
}
func (it *ListControl) ListRowsPage(rule *repo.DataRule, list []*likbase.ItElm) []*likbase.ItElm {
	return list
}

//	Построение списка
func (it *ListControl) RunListMakeBuild(rule *repo.DataRule, list []*likbase.ItElm) lik.Lister {
	if it.ItListMakeBuild == nil { return nil }
	return it.ItListMakeBuild.Run(rule, list)
}
func (it *ListControl) ListMakeBuild(rule *repo.DataRule, list []*likbase.ItElm) lik.Lister {
	rows := it.TablePageFill(rule)
	for nelm := 0; nelm < len(list); nelm++ {
		elm := list[nelm]
		if row := it.GridElmRow(rule, elm); row != nil {
			row.SetItem(it.Part, "idpart")
			row.SetItem(rule.ItPage.ProbeCollect(it.Part, elm.Id), "mark")
			it.RunRowFill(rule, elm, row)
			rows.AddItems(row)
		}
	}
	return rows
}

//	Заполнение строки
func (it *ListControl) RunRowFill(rule *repo.DataRule, elm *likbase.ItElm, row lik.Seter) {
	if it.ItRowFill != nil {
		it.ItRowFill.Run(rule, elm, row)
	}
}
func (it *ListControl) ListRowFill(rule *repo.DataRule, elm *likbase.ItElm, row lik.Seter) {
}

//	Выбор строки
func (it *ListControl) ListSelectRow(rule *repo.DataRule, sid string) {
	it.TableSelectRow(rule, sid)
	if sid != "" {
		rule.SetLocateId(likbase.StrToIDB(sid))
		rule.ItPage.PathLast = rule.BuildFullPath("")
		rule.ItPage.PathNeed = rule.ItPage.PathLast
	}
}

//	Заполнение формы
func (it *ListControl) ListFormFill(rule *repo.DataRule) {
	var elm *likbase.ItElm
	if it.Fun != fancy.FunAdd {
		elm = jone.GetElm(it.Part, likbase.StrToIDB(it.Sel))
		if elm == nil {
			return
		}
	}
	it.RunFormElm(rule, elm)
}

//	Заполнение формы из элемента
func (it *ListControl) RunFormElm(rule *repo.DataRule, elm *likbase.ItElm) {
	if it.ItFormElm != nil {
		it.ItFormElm.Run(rule, elm)
	}
}
func (it *ListControl) ListFormElm(rule *repo.DataRule, elm *likbase.ItElm) {
	name := ""
	if it.Part == "bell" {
		name = "контакта"
	} else if it.Part == "offer" {
		name = "заявки"
	} else if it.Part == "object" {
		name = "объекта"
	} else if it.Part == "member" {
		name = "сотрудника"
	} else if it.Part == "depart" {
		name = "подразделения"
	} else if it.Part == "client" {
		name = "клиента"
	}
	if it.Fun == fancy.FunShow {
		it.SetTitle(rule, it.Fun, "Карточка "+name)
		it.AddTitleToolText(rule, "Изменить", "function_fancy_form_toedit")
		if rule.IAmAdmin() {
			it.AddTitleToolText(rule, "Удалить", "function_fancy_form_todelete")
		}
		it.AddTitleToolText(rule, "Закрыть", "function_fancy_form_cancel")
	} else if it.Fun == fancy.FunAdd {
		it.SetTitle(rule, it.Fun, "Создание "+name)
		it.AddTitleToolText(rule, "Сохранить", "function_fancy_form_write")
		it.AddTitleToolText(rule, "Отменить", "function_fancy_form_cancel")
	} else if it.Fun == fancy.FunMod {
		it.SetTitle(rule, it.Fun, "Редактирование "+name)
		it.AddTitleToolText(rule, "Сохранить", "function_fancy_form_write")
		it.AddTitleToolText(rule, "Отменить", "function_fancy_form_toshow")
	} else if it.Fun == fancy.FunEdit {
		it.SetTitle(rule, it.Fun, "Редактирование "+name)
		it.AddTitleToolText(rule, "Сохранить", "function_fancy_form_write")
		it.AddTitleToolText(rule, "Отменить", "function_fancy_form_cancel")
	} else if it.Fun == fancy.FunDel {
		it.SetTitle(rule, it.Fun, "Удаление "+name)
		it.AddTitleToolText(rule, "Действительно удалить?", "function_fancy_real_delete")
		it.AddTitleToolText(rule, "Отменить", "function_fancy_form_toshow")
	}
	it.Form.Items = it.FormElmFill(rule, elm,"")
}

//	Удаление элемента
func (it *ListControl) cmdElmDelete(rule *repo.DataRule) {
	jone.DeleteElm(it.Part, likbase.StrToIDB(it.Sel))
	it.Sel = ""
	rule.SetGoPart("/"+it.Main)
}

//	Запись изменений
func (it *ListControl) ListCmdWrite(rule *repo.DataRule, data lik.Seter) *likbase.ItElm {
	elm := it.ListCmdElm(rule)
	it.ListCmdUpdate(rule, elm, data)
	return elm
}

//	Выбор или создание элемента
func (it *ListControl) ListCmdElm(rule *repo.DataRule) *likbase.ItElm {
	var elm *likbase.ItElm
	if it.Fun == fancy.FunMod || it.Fun == fancy.FunEdit {
		elm = jone.GetElm(it.Part, likbase.StrToIDB(it.Sel))
	} else if it.Fun == fancy.FunAdd {
		elm = jone.GetTable(it.Part).CreateElm()
	}
	if elm != nil {
		it.Sel = likbase.IDBToStr(elm.Id)
	}
	return elm
}

//	Обновление элемента
func (it *ListControl) ListCmdUpdate(rule *repo.DataRule, elm *likbase.ItElm, data lik.Seter) *likbase.ItElm {
	if elm == nil {
			return nil
		}
	if elm.Table.Part == "offer" && jone.CalculateElmIDB(elm,"objectid") == 0 {
		obj := jone.TableObject.CreateElm()
		jone.SetElmValue(elm, obj.Id, "objectid")
		jone.SetElmValue(elm, rule.ItSession.IdMember, "memberid")
	}
	it.UpdateElmData(rule, elm, data)
	if elm.Table.Part == "member" {
		if it,ok := one.GetMember(elm.Id); ok {
			repo.SynchronizeMemberOne(&it, elm)
		}
	} else if elm.Table.Part == "depart" {
		if it,ok := one.GetDepart(elm.Id); ok {
			repo.SynchronizeDepartOne(&it, elm)
		}
	}
	return elm
}

//	Обработка команд
func (it *ListControl) GoWindowMode(rule *repo.DataRule, mode string, sel string) {
	if sel != "" {
		path := fmt.Sprintf("/%s%s?_tp=1", mode, sel)
		rule.SetResponse(path, "_function_lik_window_part")
	}
}

