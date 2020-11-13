package fancy

import (
	"bitbucket.org/961961/tsan/jone"
	"bitbucket.org/961961/tsan/repo"
	"bitbucket.org/shaman/lik"
	"bitbucket.org/shaman/lik/likbase"
	"fmt"
	"regexp"
	"strings"
	"time"
)

//	Дескриптор объекта формы
type FancyForm struct {
	FancyCore				//	Ядро
	Title     string		//	Заголовок
	Tools     lik.Lister	//	Инструменты
	Tabs      lik.Lister	//	Закладки
	Items     lik.Lister	//	Элементы
	SingleCho	bool
	Tab		int				//	Текущая закладка
}

//	Очистка дескриптора
func (it *DataFancy) FormClear() {
	it.Form.FancyClear()
	it.Form.Title = ""
	it.Form.Tools = lik.BuildList()
	it.Form.Tabs = lik.BuildList()
	it.Form.Items = lik.BuildList()
	it.Form.SingleCho = false
}

//	Установка заголовка
func (it *DataFancy) SetTitle(rule *repo.DataRule, fun FunData, text string) {
	it.Form.Title = text
	if fun == FunAdd || fun == FunMod || fun == FunShow {
		it.Form.Class = "edit-fancy-form"
	} else if fun == FunDel {
		it.Form.Class = "delete-fancy-form"
	} else {
		it.Form.Class = "show-fancy-form"
	}
}

//	Добавить текст инструменты
func (it *DataFancy) AddTitleToolText(rule *repo.DataRule, text string, deal string) {
	it.Form.Tools.AddItemSet("text="+text, "handler="+deal)
}

//	Заполнить поля из объекта
func (it *DataFancy) FormElmFill(rule *repo.DataRule, elm *likbase.ItElm, part string) lik.Lister {
	var info lik.Seter
	if elm != nil { info = elm.Info }
	return it.FormInfoFill(rule, info, part)
}
func (it *DataFancy) FormInfoFill(rule *repo.DataRule, info lik.Seter, part string) lik.Lister {
	if fields := it.ListFields; fields != nil {
		return it.FormFillFields(rule, info, part, fields)
	}
	return lik.BuildList()
}

//	Заполнить коллекцию полей
func (it *DataFancy) FormElmCollect(rule *repo.DataRule, elm *likbase.ItElm, part string, key string) lik.Lister {
	var info lik.Seter
	if elm != nil { info = elm.Info }
	return it.FormInfoCollect(rule, info, part, key)
}

//	Заполнить поля из структуры
func (it *DataFancy) FormInfoCollect(rule *repo.DataRule, info lik.Seter, part string, key string) lik.Lister {
	if ent := repo.GenStruct.FindEnt(key); ent != nil {
		if content := ent.GetContent(); content != nil {
			return it.FormFillFields(rule, info, part, content)
		}
	}
	return lik.BuildList()
}

//	Заполнить список полей
func (it *DataFancy) FormFillFields(rule *repo.DataRule, info lik.Seter, part string, fields []lik.Seter) lik.Lister {
	list := lik.BuildList()
	for _, field := range fields {
		if (field.GetInt("tags") & jone.TagForm) != 0 {
			if it.RunFieldProbe(rule, field) {
				it.FormAppendField(rule, info, part, field, list)
			}
		}
	}
	return list
}

//	Добавить список полей
func (it *DataFancy) FormAppendField(rule *repo.DataRule, info lik.Seter, path string, field lik.Seter, list lik.Lister) {
	format := field.GetString("format")
	label := field.GetString("name")
	pt := field.GetString("part")
	if path != "" {
		path += "/" + pt
	} else {
		path = pt
	}
	value := it.RunCalculate(rule, info, path, format, true)
	tags := field.GetInt("tags")
	if strings.HasSuffix(path, "squarerooms") && it.IsEdit() {
		it.FormAppendFieldRooms(rule, info, label, path, format, tags, value, list)
	} else if /*part == "date" &&*/ format == "t" {
		it.FormAppendFieldDateTime(rule, info, label, path, format, tags, value, list)
	} else {
		list.AddItems(it.FormBuildItem(rule, label, path, format, tags, value))
	}
}

//	Добавить поля площади комнат
func (it *DataFancy) FormAppendFieldRooms(rule *repo.DataRule, info lik.Seter, label string,
	part string, format string, tags int, value lik.Itemer, list lik.Lister) {
	prefix := ""
	if match := lik.RegExParse(part,"^(.*)squarerooms$"); match != nil {
		prefix = match[1]
	}
	rooms := 1
	if srm := jone.CalculateString(info, "objectid/define/rooms"); srm != "" {
		if match := regexp.MustCompile("^(\\d+)").FindStringSubmatch(srm); match != nil {
			if nrm := lik.StrToInt(match[1]); nrm > rooms {
				rooms = nrm
			}
		}
	}
	var lval []string
	if value != nil {
		lval = strings.Split(value.ToString(), "+")
	}
	if len(lval) > rooms {
		rooms = len(lval)
	}
	for r := 0; r < rooms; r++ {
		rlabel := fmt.Sprintf("Комната №%d", r+1)
		rpart := prefix + fmt.Sprintf("rooms_%d", r)
		rval := ""
		if r < len(lval) {
			rval = lval[r]
		}
		list.AddItems(it.FormBuildItem(rule, rlabel, rpart, format, tags, lik.BuildItem(rval)))
	}
}

//	Добавить поля даты и времени
func (it *DataFancy) FormAppendFieldDateTimeNew(rule *repo.DataRule, info lik.Seter, label string,
	part string, format string, tags int, value lik.Itemer, list lik.Lister) {
	itemdt := it.FormBuildItem(rule, label, part, format, tags, value)
	_ = itemdt
	vtm := ""
	if value != nil {
		if mtm := regexp.MustCompile("(\\d+:\\d+)").FindStringSubmatch(value.ToString()); mtm != nil {
			vtm = mtm[1]
		}
	}
	itemtm := it.FormBuildItem(rule,"Время","time","w", tags, lik.BuildItem(vtm))
	_ = itemtm
	list.AddItemSet("type=line", "items", lik.BuildList(
		lik.BuildSet("type=html", "value=DTDT"),
		//itemdt,
		))
}

//	Добавить поля даты и времени
func (it *DataFancy) FormAppendFieldDateTime(rule *repo.DataRule, info lik.Seter, label string,
	part string, format string, tags int, value lik.Itemer, list lik.Lister) {
	list.AddItems(it.FormBuildItem(rule, label, part, format, tags, value))
	vtm := ""
	if value != nil {
		if mtm := regexp.MustCompile("(\\d+:\\d+)").FindStringSubmatch(value.ToString()); mtm != nil {
			vtm = mtm[1]
		}
	}
	list.AddItems(it.FormBuildItem(rule,"Время","time","w", tags, lik.BuildItem(vtm)))
}

//	Зафиксировать номер закладки
func (it *DataFancy) FormFixTab(parm string) {
	if iparm,ok := lik.StrToIntIf(parm); ok {
		it.Form.Tab = iparm
	}
}

//	Построить поле со значением
func (it *DataFancy) FormBuildItem(rule *repo.DataRule, label string, part string, format string, tags int, value lik.Itemer) lik.Seter {
	item := lik.BuildSet()
	item.SetItem(label, "label")
	//item.SetItem("left", "labelAlign")
	if format == "" { format = "u" }
	item.SetItem(format+"_"+strings.ReplaceAll(part, "/", "__"), "name")
	editable := it.IsEdit() && (tags &jone.TagEdit) != 0
	item.SetItem(editable, "editable")
	if (tags & jone.TagMust) != 0 {
		item.SetItem("must", "cls")
	} else if it.IsEdit() && !editable {
		item.SetItem("readonly", "cls")
	}
	key := part
	if match := lik.RegExParse(part,"([^/]+)$"); match != nil {
		key = match[1]
		if !it.IsEdit() {
			item.SetItem(key, "key")
		}
	}
	if format == "d" || format == "t" {
		item.SetItem("date", "type")
		item.SetItem(lik.BuildSet("read=Y/m/d", "write=d/m/Y", "edit=d/m/Y"), "format")
	} else if format == "w" {
		item.SetItem("function_fancy_input_time", "format/inputFn")
	} else if format == "n" || format == "m" {
		item.SetItem("function_fancy_input_number", "format/inputFn")
	} else if format == "p" {
		item.SetItem("function_fancy_input_telephone", "format/inputFn")
	} else if format == "h" {
		item.SetItem("checkbox","type")
	} else if (format == "b" || format == "c" || format == "l") {
		if it.IsEdit() {
			dict,vals := it.FancyBuildDictVals(rule, key, format, value)
			item.SetItem("combo", "type")
			item.SetItem("part", "valueKey")
			item.SetItem("name", "displayKey")
			item.SetItem(dict, "data")
			if format == "l" {
				if !it.Form.SingleCho {
					item.SetItem(true, "multiSelect")
					value = vals
				}
			} else if vals.Count() == 0 {
				value = nil
			}
		} else if value != nil {
			value = lik.BuildItem(jone.SystemTranslate(key, value))
		}
	} else if match := regexp.MustCompile("([^/]+)id$").FindStringSubmatch(part); match != nil {
		link := match[1]
		if table := jone.GetTable(link); table != nil {
			if it.IsEdit() {
				dict,found := it.FancyBuildChoose(rule, table, value)
				item.SetItem("combo", "type")
				item.SetItem("name", "displayKey")
				item.SetItem("id", "valueKey")
				item.SetItem(dict, "data")
				if !found { value = nil }
			} else if value != nil && value.ToInt() > 0 {
				if elm := jone.GetElm(link, lik.IDB(value.ToInt())); elm != nil {
					value = lik.BuildItem(jone.CalculateElmText(elm))
				} else {
					value = lik.BuildItem("(удалён)")
				}
			}
		}
	}
	if value == nil {
		item.SetItem("", "value")
	} else {
		item.SetItem(value, "value")
	}
	return item
}

//	Обновить поля формы
func (it *DataFancy) UpdateElmData(rule *repo.DataRule, elm *likbase.ItElm, data lik.Seter) {
	if elm.Info == nil {
		elm.Info = lik.BuildSet()
		elm.OnModify()
	}
	if modes := it.UpdateInfoData(rule, elm.Info, data); modes != nil {
		if elm.Table.Part == "offer" {
			history := lik.BuildSet("what=modify")
			history.SetItem(rule.ItSession.IdMember, "memberid")
			history.SetItem(modes, "changes")
			repo.AddHistory(rule, elm, history)
		}
		elm.OnModify()
	}
}

//	Обновить поля формы структуры
func (it *DataFancy) UpdateInfoData(rule *repo.DataRule, info lik.Seter, data lik.Seter) lik.Lister {
	var modes lik.Lister
	for _, set := range data.Values() {
		if key := set.Key; key != "id" {
			value := lik.StringFromXS(set.Val.ToString())
			if strings.HasPrefix(key, "t_") && value != "" {
				if stm := data.GetString("w_time"); stm != "" {
					value += " " + lik.StringFromXS(stm)
				}
			} else if match := lik.RegExParse(key,"^(.+)rooms_0$"); match != nil {
				path := match[1]
				key = path + "squarerooms"
				value = ""
				for r := 0; r < 10; r++ {
					if xval := data.GetString(fmt.Sprintf("%srooms_%d", path, r)); xval != "" {
						val := lik.StringFromXS(xval)
						val = strings.ReplaceAll(val, " ", "")
						val = strings.ReplaceAll(val, ",", ".")
						if val != "" {
							if value != "" {
								value += "+"
							}
							value += val
						}
					}
				}
			} else if strings.Contains(key, "rooms_") {
				continue
			} else if key == "w_time" {
				continue
			} else if lik.RegExCompare(key,"^(\\w)__") {
				continue
			}
			if it.UpdateInfoKeyVal(rule, info, key, value) {
				if modes == nil {
					modes = lik.BuildList()
				}
				modes.AddItems(key)
			}
		}
	}
	return modes
}

//	Обновить поле ключом и значением
func (it *DataFancy) UpdateInfoKeyVal(rule *repo.DataRule, info lik.Seter, key string, valstr string) bool {
	modify := false
	format := "u"
	if match := lik.RegExParse(key,"^(\\w)?_(.*)"); match != nil {
		if match[1] != "" { format = match[1] }
		key = match[2]
	}
	path := strings.ReplaceAll(key, "__", "/")
	value := lik.BuildItem(valstr)
	if valstr == "" {
		value = nil
	} else if format == "p" {
		valstr = jone.NormalizePhone(valstr)
		value = lik.BuildItem(valstr)
	} else if format == "d" || format == "t" {
		if match := lik.RegExParse(valstr,"(\\d+)\\D+(\\d+)\\D+(\\d\\d\\d\\d)"); match != nil {
			year := lik.StrToInt(match[3])
			month := lik.StrToInt(match[2])
			day := lik.StrToInt(match[1])
			hour := 0
			minute := 0
			if mtm := lik.RegExParse(valstr,"(\\d+):(\\d+)"); mtm != nil {
				hour = lik.StrToInt(mtm[1])
				minute = lik.StrToInt(mtm[2])
				if hour > 23 {
					hour = 23
					minute = 59
				}
				if minute > 59 {
					minute = 59
				}
			}
			dt := time.Date(year, time.Month(month), day, hour, minute, 0, 0, time.Local)
			value = lik.BuildItem(int(dt.Unix()))
		}
	} else if format == "t" {
		if match := lik.RegExParse(valstr,"(\\d+)\\D+(\\d+)"); match != nil {
			hour := lik.StrToInt(match[1])
			minute := lik.StrToInt(match[2])
			if hour > 23 {
				hour = 23
				minute = 59
			}
			if minute > 59 {
				minute = 59
			}
			dt := hour*60 + minute
			value = lik.BuildItem(dt)
		}
	} else if ival, ok := lik.StrToIntIf(strings.ReplaceAll(valstr, " ", "")); ok {
		value = lik.BuildItem(ival)
	} else if fval, ok := lik.StrToFloatIf(strings.ReplaceAll(valstr, " ", "")); ok {
		value = lik.BuildItem(fval)
	} else if valstr == "false" {
		value = lik.BuildItem(false)
	} else if valstr == "true" {
		value = lik.BuildItem(true)
	}
	if strings.Contains(key, "password") {
		if valstr != "" {
			valmd5 := lik.GetMD5Hash(valstr)
			if jone.SetInfoValue(info, valmd5, path) {
				modify = true
			}
		}
	} else {
		if jone.SetInfoValue(info, value, path) {
			modify = true
		}
	}
	return modify
}

//	Отображение формы
func (it *DataFancy) ShowForm(rule *repo.DataRule) {
	code := it.Form.Parameters.Clone().ToSet()
	if !code.IsItem("title") {
		title := lik.BuildSet()
		title.SetItem(it.Form.Title, "text")
		if it.Form.Tools.Count() > 0 {
			title.SetItem(it.Form.Tools,"tools")
		}
		code.SetItem(title,"title")
	}
	if !code.IsItem("tabs") && it.Form.Tabs.Count() > 0 {
		code.SetItem(it.Form.Tabs, "tabs")
		if !code.IsItem("activeTab") && it.Form.Tab >= 0 && it.Form.Tab < it.Form.Tabs.Count() {
			code.SetItem(it.Form.Tab, "activeTab")
		}
	}
	if !code.IsItem("items") {
		if it.Form.Items.Count() > 0 {
			code.SetItem(it.Form.Items, "items")
		} else {
			code.SetItem(lik.BuildSet("type=html", "html=(Нет полей)"),"items")
		}
	}
	if it.Form.Width == 0 {
		it.Form.Width = 560
	}
	if it.Form.Height == 0 {
		//it.Form.Height = "fit"
	}
	if it.Form.Class == "" {
		it.Form.Class = "show-fancy-form"
	}
	code.SetItem(true, "isform")
	it.Form.FillCore(code)
	it.ShowResult(rule, code)
}
