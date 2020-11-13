package fancy

import (
	"github.com/massarakhsh/tsan/jone"
	"github.com/massarakhsh/tsan/repo"
	"github.com/massarakhsh/lik"
	"github.com/massarakhsh/lik/likbase"
	"regexp"
	"strings"
)

//	Дескриптор объекта таблицы
type FancyGrid struct {
	FancyCore				//	Ядро объекта
	Columns    lik.Lister	//	Коллекция колонок
	Tops       lik.Lister	//	Коллекция кнопок
	Context    	lik.Lister	//	Меню контекста
}

//	Интерфейс проверки колонки
type DealColumnProbe interface {
	Run(rule *repo.DataRule, column lik.Seter) bool
}
type dealColumnProbe struct {
	It	*DataFancy
}
func (it *dealColumnProbe) Run(rule *repo.DataRule, column lik.Seter) bool {
	return true
}

//	Очистка дескриптора
func (it *DataFancy) GridClear() {
	it.Grid.FancyClear()
	it.Grid.Columns = lik.BuildList()
	it.Grid.Tops = lik.BuildList()
	it.Grid.Context = lik.BuildList()
}

//	Добавить события таблицы
func (it *DataFancy) AppendGridEvents(rule *repo.DataRule) {
	it.Grid.AddEventAction("selectrow", "function_fancy_grid_rowselect")
	it.Grid.AddEventAction("contextmenu", "function_fancy_grid_rclick")
	it.Grid.AddEventAction("celldblclick", "function_fancy_grid_dblclick")
	it.Grid.AddEventAction("columnresize", "function_fancy_col_size")
	it.Grid.AddEventAction("columndrag", "function_fancy_col_drag")
}

//	Добавить разделитель
func (it *DataFancy) AddSeparate(rule *repo.DataRule, ord int) {
}

//	Добавить кнопку с изображением
func (it *DataFancy) AddCommandImg(rule *repo.DataRule, ord int, text string, cmd string, img string) {
	dis := (img == "add") && false;
	it.AddCommandItem(rule, ord, lik.BuildSet(
		"type=button", "tip="+text, "imageCls=img"+img, "disabled", dis,
		"handler=function_fancy_grid_cmd", "cmd", cmd,
	))
}

//	Добавить кнопку
func (it *DataFancy) AddCommandItem(rule *repo.DataRule, ord int, item lik.Seter) {
	var pos int
	var itis bool
	for pos = it.Grid.Tops.Count(); pos > 0; pos-- {
		if elm := it.Grid.Tops.GetSet(pos - 1); elm != nil {
			if eord := elm.GetInt("_ord"); eord <= ord {
				itis = (eord == ord)
				break
			}
		}
	}
	if itis && item == nil {
		it.Grid.Tops.DelItem(pos - 1)
	} else if itis {
		item.SetItem(ord, "_ord")
		it.Grid.Tops.SetItem(item, pos - 1)
	} else if item != nil {
		item.SetItem(ord, "_ord")
		it.Grid.Tops.InsertItem(item, pos)
	}
}

//	Построить список колонок
func (it *DataFancy) GridBuildColumns(rule *repo.DataRule, numerate bool, editable bool) {
	if numerate {
		clm := lik.BuildSet("type=order", "index=c", "title=№", "width=40", "locked=true", "draggable=false")
		if it.RunColumnProbe(rule, clm) {
			it.Grid.Columns.AddItems(clm)
		}
	}
	type CLMN struct {
		Column	lik.Seter
		PosOrg	int
		PosNew	int
	}
	clmns := []*CLMN{}
	clmnmap := make(map[string]*CLMN)
	for _,column := range it.ListFields {
		tags := column.GetInt("tags")
		if (tags & jone.TagGrid) == 0 { continue }
		clm := lik.BuildSet()
		clm.SetItem(true, "sortable")
		part := column.GetString("part")
		key := part
		if match := lik.RegExParse(part,"([^/]+)$"); match != nil {
			key = match[1]
			clm.SetItem(key, "key")
		}
		index := column.GetString("index")
		if index == "" {
			index = strings.Replace(part, "/", "__", -1)
		}
		clm.SetItem(index, "index")
		clm.SetItem(column.GetString("name"), "title")
		clm.SetItem((column.GetInt("tags")&jone.TagHide) != 0, "hidden")
		clm.SetItem(true, "ellipsis")
		if part != "id" && part != "pic" {
			clm.SetItem(editable, "draggable")
		}
		if width := column.GetInt("width"); width > 0 {
			clm.SetItem(width, "width")
		}
		align := column.GetString("align")
		if align == "" { align = "left" }
		clm.SetItem(align, "cellAlign")
		format := column.GetString("format")
		//clm.SetItem(lik.BuildSet("header=true", "headerNote=true", "emptyText=Шо", "fn=function_fancy_column_filter"),"filter")
		if format == "h" {
			clm.SetItem("center", "cellAlign")
			clm.SetItem("checkbox", "type")
		} else if format == "r" {
			clm.SetItem("tree", "type")
		} else if format == "o" {
			clm.SetItem("order", "type")
			clm.SetItem(true, "locked")
		} else if format == "g" {
			clm.SetItem("image", "type")
			clm.SetItem("gridimg", "cls")
		} else if format == "p" {
			clm.SetItem("function_fancy_render_phone", "render")
		} else if format == "d" || format == "t" {
			clm.SetItem("function_fancy_render_ymd", "render")
		} else if format == "n" || format == "m" {
			clm.SetItem("right", "cellAlign")
			clm.SetItem("number", "type")
			if !strings.Contains(column.GetString("part"), "year") {
				clm.SetItem("function_fancy_render_number_"+format, "render")
			}
		} else if format == "b" || format == "c" || format == "l" {
			dict,_ := it.FancyBuildDictVals(rule, key, format, nil)
			clm.SetItem("combo", "type")
			clm.SetItem(true, "multiSelect")
			clm.SetItem(160, "minListWidth")
			clm.SetItem(true, "itemCheckBox")
			lds := lik.BuildList()
			lds.AddItems("")
			for nd := 0; nd < dict.Count(); nd++ {
				if name := dict.GetSet(nd).GetString("name"); name != "" {
					lds.AddItems(name)
				}
			}
			clm.SetItem(lds, "data")
		}
		if it.RunColumnProbe(rule, clm) {
			clmn := &CLMN{ Column: clm, PosOrg: len(clmns), PosNew: 0 }
			clmns = append(clmns, clmn)
			if index != "" {
				clmnmap[index] = clmn
			}
		}
	}
	if it.RuleFilter != nil {
		if fcols := it.RuleFilter.GetList("cols"); fcols != nil {
			poses := make([]bool, 1000)
			for nf := 0; nf < fcols.Count(); nf++ {
				if fcol := fcols.GetSet(nf); fcol != nil {
					index := strings.ReplaceAll(fcol.GetString("part"), "/", "__")
					if clmn,_ := clmnmap[index]; clmn != nil {
						if width := fcol.GetInt("width"); width > 0 {
							clmn.Column.SetItem(width, "width")
						}
						if tags := fcol.GetInt("tags"); (tags & jone.TagTune) != 0 {
							tago := clmn.Column.GetInt("tags")
							if (tags & jone.TagHide) != 0 {
								tago |= jone.TagHide
							} else {
								tago &= 0xffffff ^ jone.TagHide
							}
							clmn.Column.SetItem(tago, "tags")
						}
						if order := fcol.GetInt("order"); order > 0 {
							clmn.PosNew = order
							if order < len(poses) {
								poses[order] = true
							}
						}
					}
				}
			}
			pc := 0
			for nc := 0; nc < len(clmns); nc++ {
				clmn := clmns[nc]
				if clmn.PosNew == 0 {
					for pc = pc+1; pc < len(poses); pc++ {
						if !poses[pc] {
							break
						}
					}
					clmn.PosNew = pc
				}
			}
		}
	}
	for true {
		var minco *CLMN
		for nc := 0; nc < len(clmns); nc++ {
			clmn := clmns[nc]
			if clmn.Column != nil {
				if minco == nil || clmn.PosNew < minco.PosNew {
					minco = clmn
				}
			}
		}
		if minco == nil {
			break
		}
		it.Grid.Columns.AddItems(minco.Column)
		minco.Column = nil
	}
}

//	Настроить фильтры
func (it *DataFancy) GridTuneFilter(rule *repo.DataRule) {
	if it.RuleFilter != nil {
		datas := it.FancyFilterDecode(rule)
		for nc := 0; nc < len(datas.Conds); nc++ {
			key := datas.Conds[nc].Key
			opr := datas.Conds[nc].Opr
			val := datas.Conds[nc].Val
			if opr == "|" {
				chs := lik.BuildList()
				for _, ch := range strings.Split(val, ",") {
					if len(ch) > 0 {
						chs.AddItems(ch)
					}
				}
				if chs.Count() > 0 {
					it.Grid.SetParameter(chs, "state/filters/"+key+"/"+opr)
				}
			} else {
				if opr == "" {
					opr = "="
				}
				it.Grid.SetParameter(val, "state/filters/"+key+"/"+opr)
			}
		}
		if datas.SortKey != "" {
			dir := "ASC"
			if !datas.SortDir {
				dir = "DESC"
			}
			srt := lik.BuildList(lik.BuildSet("key", datas.SortKey, "dir", dir))
			it.Grid.SetParameter(srt, "state/sorters")
		}
	}
}

//	Проверка колонки
func (it *DataFancy) RunColumnProbe(rule *repo.DataRule, column lik.Seter) bool {
	return it.ItColumnProbe.Run(rule, column)
}

//	Заполнение строки
func (it *DataFancy) GridElmRow(rule *repo.DataRule, elm *likbase.ItElm) lik.Seter {
	if elm == nil {
		return nil
	}
	id := elm.Id
	return it.GridInfoRow(rule, likbase.IDBToStr(id), elm.Info)
}

//	Изготовить строку
func (it *DataFancy) GridInfoRow(rule *repo.DataRule, id string, info lik.Seter) lik.Seter {
	row := lik.BuildSet("id", id)
	if it.ListFields != nil {
		for _, field := range it.ListFields {
			index := field.GetString("index")
			part := field.GetString("part")
			format := field.GetString("format")
			value := it.RunCalculate(rule, info, part, format, false)
			if index == "id" {
			} else if value == nil {
				row.SetItem("", index)
			} else if match := regexp.MustCompile("([^/]+)id$").FindStringSubmatch(part); match != nil {
				text := jone.CalculatePartIdText(match[1], lik.IDB(value.ToInt()))
				row.SetItem(text, index)
			} else if format != "b" && format != "c" && format != "l" {
				row.SetItem(value, index)
			} else if match := lik.RegExParse(part,"([^/]+)$"); match != nil {
				text := jone.SystemTranslate(match[1], value)
				row.SetItem(text, index)
			} else {
				row.SetItem(value, index)
			}
		}
	}
	return row
}

//	Добавить элемент в контекстное меню
func (it *DataFancy) AddContextMenu(text string, side string, part string) {
	if text != "" {
		set := lik.BuildSet("text", text)
		if side != "" {
			set.SetItem(side, "sideText")
		}
		if part != "" {
			set.SetItem(part, "part")
			set.SetItem("function_fancy_grid_context", "handler")
		}
		it.Grid.Context.AddItems(set)
	} else {
		it.Grid.Context.AddItems(lik.BuildItem("-"))
	}
}

//	Установить видимость колонки
func (it *DataFancy) GridColumnVisible(rule *repo.DataRule, index string, visible bool) {
	if visible {
		it.ChangeColumn(rule,"colshow", index, 1)
	} else {
		it.ChangeColumn(rule,"colshow", index, 0)
	}
}

//	Изменить колонку
func (it *DataFancy) ChangeColumn(rule *repo.DataRule, cmd string, index string, prm int) bool {
	ok := false
	my := it.FancyFilterSeek(rule) >= 0
	if rule.IAmAdmin() && !my {
		ok = it.ChangeColumnTune(rule, cmd, index, prm)
	} else if !my {
		it.FancyFilterSaveAs(rule,"Новый")
		it.ChangeColumnFilter(rule, cmd, index, prm)
		ok = false
	} else {
		it.ChangeColumnFilter(rule, cmd, index, prm)
		ok = true
	}
	if cmd == "colshow" || cmd == "colsize" { ok = false }
	rule.SaveMemberParam()
	return ok
}

//	Настроить колонку прототипа
func (it *DataFancy) ChangeColumnTune(rule *repo.DataRule, cmd string, index string, prm int) bool {
	ok := true
	part := index
	if field := it.SeekPartField(part); field != nil {
		part = field.GetString("part")
	}
	if ent := repo.GenStruct.FindEnt(it.GetParter()); ent != nil {
		if content := jone.CalculateElmList(ent.It, "content"); content != nil {
			if idx, prt := ent.FindPart(part); idx >= 0 {
				if cmd == "colsize" {
					prt.SetItem(prm, "width")
				} else if cmd == "colshow" {
					tags := prt.GetInt("tags")
					if prm > 0 {
						tags &= jone.TagHide ^ 0xffffff
					} else {
						tags |= jone.TagHide
					}
					prt.SetItem(tags, "tags")
				} else if cmd == "coldrag" {
					if pold,col := it.seekColumnIndex(rule, index); col != nil {
						dst := prm - (pold - 1)
						content.DelItem(idx)
						content.InsertItem(prt, idx + dst)
						ok = false
					}
				}
				ent.SaveToBase()
			}
		}
	}
	return ok
}

//	Настроить колонку фильтра
func (it *DataFancy) ChangeColumnFilter(rule *repo.DataRule, cmd string, index string, prm int) bool {
	ok := true
	it.FancyFilterSeek(rule)
	part := index
	field := it.SeekPartField(part)
	if field == nil { return ok }
	cols := it.RuleFilter.GetList("cols")
	if cols == nil {
		cols = lik.BuildList()
		it.RuleFilter.SetItem(cols, "cols")
	}
	part = field.GetString("part")
	var set lik.Seter
	for nf := 0; nf < cols.Count(); nf++ {
		if flt := cols.GetSet(nf); flt != nil && part == flt.GetString("part") {
			set = flt
			break
		}
	}
	if set == nil {
		set = lik.BuildSet("part", part,
			"width", field.GetInt("width"),
			"tags", field.GetInt("tags"),
		)
		cols.AddItems(set)
	}
	if cmd == "colsize" {
		set.SetItem(prm, "width")
	} else if cmd == "colshow" {
		tags := set.GetInt("tags") | jone.TagTune
		if prm > 0 {
			tags &= jone.TagHide ^ 0xffffff
		} else {
			tags |= jone.TagHide
		}
		set.SetItem(tags, "tags")
		ok = false
	} else if cmd == "coldrag" {
		if pold,col := it.seekColumnIndex(rule, index); col != nil {
			dst := prm - (pold - 1)
			_ = dst
			for nf := 0; nf < cols.Count(); nf++ {
				if flt := cols.GetSet(nf); flt != nil {
					if part == flt.GetString("part") {
						flt.SetItem(prm + 1, "order")
					} else if ord := flt.GetInt("order"); ord > 0 {
						if ord-1 < pold && ord-1 >= prm {
							flt.SetItem(ord + 1, "order")
						} else if ord-1 > pold && ord-1 < prm {
							flt.SetItem(ord - 1, "order")
						}
					}
				}
			}
			ok = false
		}
		/*if pold,col := it.seekColumnIndex(rule, index); col != nil {
			for nc := 1; nc < it.Grid.Columns.Count(); nc++ {
				if col := it.Grid.Columns.GetSet(nc); col != nil {
					if col.GetString("index") == index {
						set.SetItem(prm, "order")
						//dst := prm - (nc - 1)
						//content.DelItem(idx)
						//content.InsertItem(prt, idx + dst)
						ok = false
						break
					}
				}
			}
		}*/
	}
	return ok
}

//	Найти колонку
func (it *DataFancy) seekColumnIndex(rule *repo.DataRule, index string) (int,lik.Seter) {
	pos := -1
	var column lik.Seter
	for nc := 1; nc < it.Grid.Columns.Count(); nc++ {
		if col := it.Grid.Columns.GetSet(nc); col != nil {
			if col.GetString("index") == index {
				pos = nc
				column = col
				break
			}
		}
	}
	return pos,column
}

//	Найти строку
func (it *DataFancy) GridSeekSel(rule *repo.DataRule, list lik.Lister) int {
	cur := -1
	if list != nil {
		for nc := 0; nc < list.Count(); nc++ {
			if row := list.GetSet(nc); row != nil {
				if row.GetString("id") == it.Sel {
					cur = nc
					break
				}
			}
		}
	}
	return cur
}

//	Отобразить таблицу
func (it *DataFancy) GridShow(rule *repo.DataRule) {
	code := it.Grid.Parameters.Clone().ToSet()
	if true {
		it.AppendGridEvents(rule)
	}
	if !code.IsItem("columns") {
		if it.Grid.Columns.Count() > 0 {
			code.SetItem(it.Grid.Columns, "columns")
		} else {
			code.SetItem(lik.BuildSet("title=(Нет колонок)"),"columns")
		}
	}
	if !code.IsItem("tbar") && it.Grid.Tops.Count() > 0 {
		for n := 0; n < it.Grid.Tops.Count(); n++ {
			if item := it.Grid.Tops.GetSet(n); item != nil && item.GetString("type") == "" {
				it.Grid.Tops.SetItem("|", n)
			}
		}
		code.SetItem(it.Grid.Tops, "tbar")
	}

	if !code.IsItem("contextmenu") && it.Grid.Context.Count() > 0 {
		code.SetItem(it.Grid.Context, "contextmenu")
	}

	data := code.GetSet("data")
	if data == nil {
		data = lik.BuildSet()
		code.SetItem(data, "data")
	}

	it.Grid.Class = "show-fancy-grid"
	it.Grid.FillCore(code)
	it.ShowResult(rule, code)
}

//	Отобразить страницу
func (it *DataFancy) GridPage(rule *repo.DataRule, rows lik.Lister) {
	rule.SetResponse(rows, "items")
	if rows != nil {
		rule.SetResponse(rows.Count(), "totalCount")
		if cur := it.GridSeekSel(rule, rows); cur >= 0 {
			rule.SetResponse(cur, "likSelect")
		}
	}
	rule.SetResponse(true,"success")
}

