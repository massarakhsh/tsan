//	Контроллер выбора заявки для сделки.
//
//	Построен на базе списка заявок
package controls

import (
	"bitbucket.org/961961/tsan/control"
	"bitbucket.org/961961/tsan/fancy"
	"bitbucket.org/961961/tsan/jone"
	"bitbucket.org/961961/tsan/repo"
	"bitbucket.org/961961/tsan/show"
	"bitbucket.org/shaman/lik"
	"bitbucket.org/shaman/lik/likbase"
	"bitbucket.org/shaman/lik/likdom"
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"strings"
)

//	Дескриптор выбора заявки для сделки
type ChooseControl struct {
	control.DataControl					//	Дескриптор общего контроллера
	fancy.TableFancy					//	Дескриптор объекта FancyGrid
	IsFrom      bool					//	Признак заявки - источника объекта
	IsTo        bool					//	Признак заявки - получателя объекта
	ReqList			[]*ReqElm			//	Список условий подбора
	ReqMap			map[string]*ReqElm	//	Коллекция условий подбора
}

//	Дескриптор условия подбора
type ReqElm struct {
	Name	string			//	Наименование поля
	Part	string			//	Ключ поля получателя
	Path	string			//	Путь поля источника
	Value	string			//	Значение поля условия
	Use		bool			//	Признак учета поля
}

//	Обработчик команд выбора
type dealChooseExecute struct {
	It	*ChooseControl
}
func (it *dealChooseExecute) Run(rule *repo.DataRule, cmd string, data lik.Seter) {
	it.It.DoExecute(rule, cmd, data)
}

//	Обработчик списка полей
type dealChooseFieldsFill struct {
	It	*ChooseControl
}
func (it *dealChooseFieldsFill) Run(rule *repo.DataRule) {
	it.It.ChooseFieldsFill(rule)
}

//	Обраотчик отображения коллекции
type dealChooseGridFill struct {
	It	*ChooseControl
}
func (it *dealChooseGridFill) Run(rule *repo.DataRule) {
	it.It.DoGridFill(rule)
}

//	Обработчик отображения страницы
type dealChoosePageFill struct {
	It	*ChooseControl
}
func (it *dealChoosePageFill) Run(rule *repo.DataRule) lik.Lister {
	return it.It.DoPageFill(rule)
}

//	Конструктор дескриптора выбора заявки для сделки.
//
//	id - исходная заявка
func BuildChooseOffer(rule *repo.DataRule, main string, id lik.IDB) *ChooseControl {
	it := &ChooseControl{ }
	it.ControlInitializeZone(main, id, "choose")
	it.TableInitialize(rule, main,"offer","choose")
	it.ItFieldsFill = &dealChooseFieldsFill{it}
	it.ItExecute = &dealChooseExecute{it}
	it.ItGridFill = &dealChooseGridFill{it}
	it.ItPageFill = &dealChoosePageFill{it}
	return it
}

//	Отображение контроллера
func (it *ChooseControl) BuildShow(rule *repo.DataRule, sx int, sy int) likdom.Domer {
	it.calculateRequest(rule)
	tbl := likdom.BuildTableClass("fill", "id=offerchoose")
	ht := 100
	tbl.BuildTrTdClass("pair", control.MiniMax(sx, ht)...).AppendItem(it.buildShowRequire(rule, sx, ht - control.BD))
	tbl.BuildTrTdClass("pair", control.MiniMax(sx, sy - ht)...).AppendItem(it.buildShowChoose(rule, sx, sy - ht - control.BD))
	tbl.BuildTrTdClass("fill")
	return tbl
}

//	Перерасчет требований
func (it *ChooseControl) calculateRequest(rule *repo.DataRule) {
	if rule.ItPage.Mask.Id != it.IdMain {
		rule.ItPage.Mask.Id = it.IdMain
		rule.ItPage.Mask.Ignor = make(map[string]bool)
	}
	it.ReqList = []*ReqElm {
		{"Сегмент", "segment", "segment", "", true},
		{"Тип. объекта", "realty", "objectid/realty", "", true},
		{"Комнат", "rooms", "objectid/define/rooms", "", true},
		{"Район", "subcity", "objectid/address/subcity", "", true},
		{"Цена от", "pricefrom", "cost", "", true},
		{"Цена до", "priceto", "cost", "", true},
		{"Площадь от", "squarefrom", "objectid/define/square", "", true},
		{"Площадь до", "squareto", "objectid/define/square", "", true},
	}
	it.ReqMap = make(map[string]*ReqElm)
	elm := jone.TableOffer.GetElm(it.IdMain)
	for _,req := range it.ReqList {
		if part := req.Part; part != "" {
			it.ReqMap[part] = req
			msk,ok := rule.ItPage.Mask.Ignor[part]
			req.Use = !ok || !msk
			if target := elm.GetString("target"); target == "buy" || target == "rent" {
				req.Value = jone.CalculateElmString(elm, "require/" + part)
			} else {
				req.Value = "*"
			}
		}
	}
}

//	отображение раздела требований
func (it *ChooseControl) buildShowRequire(rule *repo.DataRule, sx int, sy int) likdom.Domer {
	tbl := likdom.BuildTableClass("manage", control.MiniMax(sx, sy)...)
	elm := jone.TableOffer.GetElm(it.IdMain)
	mz := len(it.ReqList)
	hz := (mz + 1) / 2
	for nh := 0; nh < hz; nh++ {
		if row := tbl.BuildTr(); row != nil {
			req := it.ReqList[nh]
			it.buildShowRequestTrio(rule, row, elm, req)
			if nh + hz < mz {
				req = it.ReqList[nh+hz]
			} else {
				req = nil
			}
			it.buildShowRequestTrio(rule, row, elm, req)
			row.BuildTd("width=100%")
		}
	}
	return tbl
}

//	Отображение позиции требований
func (it *ChooseControl) buildShowRequestTrio(rule *repo.DataRule, row likdom.Domer, elm *likbase.ItElm, req *ReqElm) {
	val := ""
	if req != nil {
		val = req.Value
	}
	if td := row.BuildTd(); td != nil {
		che := td.BuildUnpairItem("input", "type=checkbox", "value=on")
		if val != "" {
			path := "/" + it.Frame + "/" + it.Mode + "/switch/" + req.Part
			if req.Use {
				che.SetAttr("checked")
				path += "/0"
			} else {
				path += "/1"
			}
			che.SetAttr("onclick", fmt.Sprintf("front_get('%s')", path))
		} else {
			che.SetAttr("disabled=true")
		}
	}
	if td := row.BuildTd(); td != nil && req != nil {
		td.BuildString(req.Name)
		if val != "" {
			td.BuildString(": ")
			if strings.Contains(req.Part, "price") {
				val = show.CashToFormat(val)
			} else if !lik.RegExCompare(req.Part, "(from|to)$") {
				val = jone.SystemStringTranslate(req.Part, val)
			}
			td.BuildString("<b>"+val+"</b>")
		}
	}
}

//	Отображение списка вариантов
func (it *ChooseControl) buildShowChoose(rule *repo.DataRule, sx int, sy int) likdom.Domer {
	it.SetSize(sx, sy)
	return show.BuildFancyGrid(it.Main,"choose")
}

//	Обработчик событий
func (it *ChooseControl) DoExecute(rule *repo.DataRule, cmd string, data lik.Seter) {
	if cmd == "all" || cmd == "choose" {
		it.DoExecute(rule, rule.Shift(), data)
	} else if cmd == "itshow" {
		if sid := rule.Shift(); sid != "" {
			path := fmt.Sprintf("/offershow%s?_tp=1", sid)
			rule.SetResponse(path,"_function_lik_window_part")
		}
	} else if cmd == "mark" {
		id := likbase.StrToIDB(rule.Shift())
		val := lik.StrToInt(rule.Shift())
		repo.MarkElmSet(rule, it.Part, id, val > 0)
		rule.OnChangeData()
	} else if cmd == "switch" {
		part := rule.Shift()
		rule.ItPage.Mask.Ignor[part] = lik.StrToInt(rule.Shift()) == 0
		rule.OnChangeData()
	} else if cmd == "toenter" {
		if top := rule.Shift(); top != "" {
			it.Sel = top
		}
		if id := lik.StrToInt(it.Sel); id > 0 {
			rule.SetResponse(id, "_function_choose_offer_pair")
		}
	} else {
		it.TableExecute(rule, cmd, data)
	}
}

//	Заполенение списк полей
func (it *ChooseControl) ChooseFieldsFill(rule *repo.DataRule) {
	it.FieldsClear()
	it.ListFields = append(it.ListFields,
		lik.BuildSet("index=id", "part=id", "name=ID", "width=100", "tags", jone.TagGrid),
		lik.BuildSet("index=segment", "part=segment", "name=Сегмент", "width=100", "tags", jone.TagGrid),
		lik.BuildSet("index=realty", "part=realty", "name=Тип недв.", "width=150", "tags", jone.TagGrid),
		lik.BuildSet("index=rooms", "part=rooms", "name=Комнат", "width=150", "tags", jone.TagGrid),
		lik.BuildSet("index=square", "part=square", "name=Площадь", "width=150", "tags", jone.TagGrid),
		lik.BuildSet("index=price", "part=price", "name=Цена", "width=150", "tags", jone.TagGrid),
		lik.BuildSet("index=address", "part=address", "name=Адрес", "width=250", "tags", jone.TagGrid),
		)
}

//	Заполнение заголовка таблицы
func (it *ChooseControl) DoGridFill(rule *repo.DataRule) {
	it.TableGridFill(rule)
	it.GridBuildColumns(rule, false,true)
	elm := jone.TableOffer.GetElm(it.IdMain)
	target := jone.CalculateElmString(elm, "target")
	it.IsFrom = target == "sale"
	it.IsTo = target == "buy"
	title := "Список подбора вариантов"
	if it.IsFrom {
		title += " (покупки)"
	} else if it.IsTo {
		title += " (продажи)"
	}
	it.AddCommandItem(rule, 80, lik.BuildSet("type=text", "text", title))
	it.AddCommandImg(rule, 910, "Открыть", "toshow", "show")
	it.Grid.SetParameter("group", "grouping/by")
	it.Grid.SetParameter("Заявки: {text}:{number}", "grouping/tpl")
	it.AddContextMenu("Открыть", "...", "itshow")
	it.AddContextMenu("", "", "")
	it.AddContextMenu("Поставить/убрать отметку", "", "mark")
}

//	Заполнение страницы таблицы
func (it *ChooseControl) DoPageFill(rule *repo.DataRule) lik.Lister {
	rows := it.TablePageFill(rule)
	deal := repo.SeekOfferDeal(it.IdMain)
	pairs := make(map[lik.IDB]bool)
	//	Отмеченные заявки
	if true {
		for key, val := range rule.ItPage.Session.Collect {
			if match := lik.RegExParse(key, "^offer(\\d+)"); match != nil && val {
				if elm := jone.TableOffer.GetElm(lik.IDB(lik.StrToInt(match[1]))); elm != nil {
					if elm.GetString("target") == deal.TargetPair {
						pairs[elm.Id] = true
						it.chooseAddRow(rule, rows, elm, "Отмеченные", true)
					}
				}
			}
		}
	}
	//	Последние просмотренные
	if false {
		if list := rule.CachePartGet("offer"); list != nil {
			for _, id := range list {
				if _, ok := pairs[id]; !ok {
					if elm := jone.TableOffer.GetElm(id); elm != nil {
						if elm.GetString("target") == deal.TargetPair {
							pairs[elm.Id] = true
							it.chooseAddRow(rule, rows, elm, "Просмотренные", true)
						}
					}
				}
			}
		}
	}
	//	Подходящие
	if true {
		if elms := jone.TableOffer.GetListElm(false); elms != nil {
			for _, elm := range elms {
				if _, ok := pairs[elm.Id]; !ok {
					if elm.GetString("target") == deal.TargetPair {
						pairs[elm.Id] = true
						it.chooseAddRow(rule, rows, elm, "Удовлетворяющие требованиям", false)
					}
				}
			}
		}
	}
	return rows
}

//	Проверка и добавление строки в таблицу (или пропуск)
func (it *ChooseControl) chooseAddRow(rule *repo.DataRule, rows lik.Lister, elm *likbase.ItElm, group string, trust bool) bool {
	if it.IsFrom {
		return it.chooseAddRowTo(rule, rows, elm, group, trust)
	} else if it.IsTo {
		return it.chooseAddRowFrom(rule, rows, elm, group, trust)
	}
	return false
}

//	Проверка сооответствия фиксированного объекта заявкам на покупку из списка
func (it *ChooseControl) chooseAddRowTo(rule *repo.DataRule, rows lik.Lister, elm *likbase.ItElm, group string, trust bool) bool {
	if !trust {
		for _,req := range it.ReqList {
			if part := req.Part; part != "" && req.Use && req.Value != "" {
				if val := jone.CalculateElmString(elm, "require/" + part); val != "" {
					if strings.HasSuffix(part, "from") {
						if lik.StrToFloat(req.Value) < lik.StrToFloat(val) {
							return false
						}
					} else if strings.HasSuffix(part, "to") {
						if lik.StrToFloat(req.Value) > lik.StrToFloat(val) {
							return false
						}
					} else if !testInList(req.Value, val) {
						return false
					}
				}
			}
		}
	}
	row := lik.BuildSet("id", elm.Id, "group", group)
	row.SetItem(rule.ItPage.ProbeCollect(it.Part, elm.Id), "mark")
	for _, field := range it.ListFields {
		index := field.GetString("index")
		val := ""
		if index == "segment" || index == "realty" || index == "rooms" {
			val = jone.CalculateElmTranslate(elm, "require/" + index)
		} else if index == "address" {
			val = jone.CalculateElmTranslate(elm, "require/subcity")
		} else if index == "square" || index == "price" {
			if data := jone.CalculateElmString(elm, index + "from"); data != "" {
				val += "от " + data
			}
			if data := jone.CalculateElmString(elm, index + "to"); data != "" {
				if val != "" {
					val += ", "
				}
				val += "до " + data
			}
		} else {
			continue
		}
		row.SetItem(val, index)
	}
	row.SetItem(fmt.Sprintf("/offershow%d?_tp=1", int(elm.Id)), "pathopen")
	rows.AddItems(row)
	return false
}

//	Проверка, что вариант удовлетворяет условию
func testInList(val string, list string) bool {
	if elms := strings.Split(list, ","); elms != nil {
		for _,elm := range elms {
			if val == elm {
				return true
			}
		}
	}
	return false
}

//	Проверка соответствия фиксированной заявки на покупку объектам продажи из списка
func (it *ChooseControl) chooseAddRowFrom(rule *repo.DataRule, rows lik.Lister, elm *likbase.ItElm, group string, trust bool) bool {
	if !trust {
		for _,req := range it.ReqList {
			if part := req.Part; part != "" && req.Use && req.Value != "" {
				if val := jone.CalculateElmString(elm, req.Path); val != "" {
					if strings.HasSuffix(part, "from") {
						if lik.StrToFloat(req.Value) > lik.StrToFloat(val) {
							return false
						}
					} else if strings.HasSuffix(part, "to") {
						if lik.StrToFloat(req.Value) < lik.StrToFloat(val) {
							return false
						}
					} else if !testInList(val, req.Value) {
						return false
					}
				}
			}
		}
	}
	row := lik.BuildSet("id", elm.Id, "group", group)
	row.SetItem(rule.ItPage.ProbeCollect(it.Part, elm.Id), "mark")
	for _, field := range it.ListFields {
		index := field.GetString("index")
		val := ""
		if index == "segment" {
			val = jone.CalculateElmTranslate(elm, index)
		} else if index == "realty" {
			val = jone.CalculateElmTranslate(elm, "objectid/" + index)
		} else if index == "rooms" {
			val = jone.CalculateElmTranslate(elm, "objectid/define/" + index)
		} else if index == "address" {
			val = jone.MakeAddress(jone.CalculateElmSet(elm, "objectid/address"))
		} else if index == "square" {
			val = jone.CalculateElmString(elm, "objectid/define/" + index)
		} else if index == "price" {
			val = show.CashToFormat(jone.CalculateElmString(elm, "cost"))
		} else {
			continue
		}
		row.SetItem(val, index)
	}
	row.SetItem(fmt.Sprintf("/offershow%d?_tp=1", int(elm.Id)), "pathopen")
	rows.AddItems(row)
	return false
}

