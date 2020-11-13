//	Модуль работы с объектами FancyGrid и FancyForm
package fancy

import (
	"bitbucket.org/961961/tsan/jone"
	"bitbucket.org/961961/tsan/one"
	"bitbucket.org/961961/tsan/repo"
	"bitbucket.org/shaman/lik"
	"bitbucket.org/shaman/lik/likbase"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
)

//	Код режима отображения формы
type FunData int

const FunNo = FunData(0)
const FunShow = FunData(1)
const FunMod = FunData(2)
const FunAdd = FunData(3)
const FunDel = FunData(4)
const FunEdit = FunData(5)

//	Дескриптор объекта
type DataFancy struct {
	Main          string				//	Режим фрейма
	Part          string				//	Раздел объекта
	Zone          string				//	Зона объекта
	Sx, Sy        int					//	Размер объекта
	ItFieldProbe  DealFieldProbe		//	Интерфейс проверки поля
	ItCalculate   DealFancyCalculate	//	Интерфейс вычислений
	ItFieldsFill  DealFancyFieldsFill	//	Интерфейс списка полей
	ItColumnProbe DealColumnProbe		//	Интерфейс проверки колонки
	ListFields    []lik.Seter			//	Список полей
	RuleFilter 	  	lik.Seter			//
	GridSync	sync.Mutex				//	Синхронизатор
	Grid       FancyGrid				//	Объект таблицы
	Sel			string					//	Текущий элемент
	Fun			FunData					//	Текущая функция
	Form       FancyForm				//	Объект формы
}

//	Интерфейс вычислений
type DealFancyCalculate interface {
	Run(rule *repo.DataRule, info lik.Seter, part string, format string, isform bool) lik.Itemer
}
type dealFancyCalculate struct {
	It	*DataFancy
}
func (it *dealFancyCalculate) Run(rule *repo.DataRule, info lik.Seter, part string, format string, isform bool) lik.Itemer {
	return it.It.FancyCalculate(rule, info, part, format, isform)
}

//	Интерфейс списка полей
type DealFancyFieldsFill interface {
	Run(rule *repo.DataRule)
}
type dealFancyFieldsFill struct {
	It	*DataFancy
}
func (it *dealFancyFieldsFill) Run(rule *repo.DataRule) {
	it.It.FancyFieldsFill(rule)
}

//	Интерфейс проверки поля
type DealFieldProbe interface {
	Run(rule *repo.DataRule, field lik.Seter) bool
}
type dealFieldProbe struct {
	It	*DataFancy
}
func (it *dealFieldProbe) Run(rule *repo.DataRule, field lik.Seter) bool {
	return it.It.FancyFieldProbe(rule, field)
}

//	Инициализация объекта
func (it *DataFancy) FancyInitialize(main string, part string, zone string) {
	it.Main = main
	it.Part = part
	it.Zone = zone
	it.ItCalculate = &dealFancyCalculate{it}
	it.ItFieldsFill = &dealFancyFieldsFill{it}
	it.ItFieldProbe = &dealFieldProbe{it}
	it.ItColumnProbe = &dealColumnProbe{it}
	it.Grid.FancyClear()
	it.Form.FancyClear()
}

//	Очистка объекта
func (it *DataFancy) FieldsClear() {
	it.ListFields = []lik.Seter{}
}

//	Проверка режима редактиования
func (it *DataFancy) IsEdit() bool {
	return it.Fun == FunMod || it.Fun == FunAdd || it.Fun == FunEdit
}

//	Проверка режима создания
func (it *DataFancy) IsCreate() bool {
	return it.Fun == FunAdd
}

//	Вычисление имени раздела
func (it *DataFancy) GetParter() string {
	code := it.Part + "_"
	if it.Main == "sale" && it.Zone == "all" {
		code += it.Main
	} else if it.Main == "buy" && it.Zone == "all" {
		code += it.Main
	} else if it.Zone != "" {
		code += it.Zone
	} else {
		code += "all"
	}
	return code
}

//	Установка размеров объекта
func (it *DataFancy) SetSize(sx int, sy int) {
	it.Sx = sx
	it.Sy = sy
}

//	Проверка поля
func (it *DataFancy) RunFieldProbe(rule *repo.DataRule, field lik.Seter) bool {
	return it.ItFieldProbe.Run(rule, field)
}
func (it *DataFancy) FancyFieldProbe(rule *repo.DataRule, field lik.Seter) bool {
	return field != nil
}

//	Проверка по тегам
func (it *DataFancy) FancyProbeTags(target string, realty string, tags int) bool {
	if target != "" && (tags & (jone.TagSale|jone.TagBuy)) != 0 {
		if target == "sale" && (tags&jone.TagSale) == 0 {
			return false
		}
		if target == "buy" && (tags&jone.TagBuy) == 0 {
			return false
		}
	}
	if realty != "" && (tags & (jone.TagFlat|jone.TagHouse|jone.TagLand)) != 0 {
		if (realty == "flat" || realty == "room") && (tags&jone.TagFlat) == 0 {
			return false
		}
		if realty == "house" && (tags&jone.TagHouse) == 0 {
			return false
		}
		if realty == "land" && (tags&jone.TagLand) == 0 {
			return false
		}
	}
	return true
}

//	Вычисление списка полей
func (it *DataFancy) RunFieldsFill(rule *repo.DataRule) {
	it.ItFieldsFill.Run(rule)
}

//	Вычисление списка полей
func (it *DataFancy) FancyFieldsFill(rule *repo.DataRule) {
	it.FieldsClear()
	it.FancyFillAppendEnt(rule, repo.GenStruct.FindEnt(it.GetParter()))
	it.FancyFieldsFilter(rule)
	if len(it.ListFields) == 0 {
		it.ListFields = []lik.Seter{lik.BuildSet("part=no", "name=Нет колонок")}
	}
}

//	Применение фильтров
func (it *DataFancy) FancyFieldsFilter(rule *repo.DataRule) {
	if it.FancyFilterSeek(rule); it.RuleFilter != nil {
		nos := len(it.ListFields)
		if scols := it.RuleFilter.GetList("cols"); scols != nil {
			for nc := scols.Count()-1; nc >= 0; nc-- {
				found := false
				if scol := scols.GetSet(nc); scol != nil {
					part := scol.GetString("part")
					for ns := 0; ns < nos; ns++ {
						col := it.ListFields[ns]
						if part == col.GetString("part") {
							if item := scol.GetItem("tags"); item != nil {
								stags := item.ToInt()
								tags := col.GetInt("tags")
								if (stags & jone.TagHide) != 0 {
									tags |= jone.TagHide
								} else {
									tags &= jone.TagHide ^ 0xffffff
								}
								col.SetItem(tags, "tags")
							}
							found = true
							break
						}
					}
				}
				if !found {
					scols.DelItem(nc)
					rule.SaveMemberParam()
				}
			}
		}
	}
}

//	Заполнение полей из сущности
func (it *DataFancy) FancyFillAppendEnt(rule *repo.DataRule, ent *repo.SysEnt) {
	if ent != nil && ent.It != nil {
		if list := jone.CalculateElmList(ent.It,"content"); list != nil {
			for ne := 0; ne < list.Count(); ne++ {
				if item := list.GetSet(ne); item != nil {
					if it.RunFieldProbe(rule, item) {
						clone := item.Clone().ToSet()
						part := clone.GetString("part")
						if part == "idu" { part = "id" }
						index := strings.ReplaceAll(part, "/","__")
						clone.SetItem(index,"index")
						it.ListFields = append(it.ListFields, clone)
					}
				}
			}
		}
	}
}

//	Построение словаря значений
func (it *DataFancy) FancyBuildDictVals(rule *repo.DataRule, part string, format string, value lik.Itemer) (lik.Lister,lik.Lister) {
	dict := lik.BuildList()
	vals := lik.BuildList()
	if format == "b" {
		dict.AddItemSet("part=yy", "name=да")
		dict.AddItemSet("part=nn", "name=нет")
		if value != nil && (value.ToString() == "yy" || value.ToString() == "nn") {
			vals.AddItems(value.ToString())
		}
	} else if ent := repo.GenDiction.FindEnt(part); ent != nil {
		if content := jone.CalculateElmList(ent.It,"content"); content != nil {
			words := []string{}
			if value != nil {
				str := value.ToString()
				mstr := strings.Split(strings.ReplaceAll(str, " ", ""), ",")
				if len(mstr) > 1 {
					words = mstr
				} else if str != "" {
					words = []string{str}
				}
			}
			for np := 0; np < content.Count(); np++ {
				word := content.GetSet(np)
				part := word.GetString("part")
				name := word.GetString("name")
				dict.AddItemSet("part", part, "name", name)
				for _,wd := range words {
					if wd == part {
						vals.AddItems(part)
					}
				}
			}
		}
	} else if part == "format" {
		for _, pa := range repo.ListFormat.Values() {
			part := pa.Key
			name := pa.Val
			dict.AddItemSet("part", part, "name", name)
			if value != nil && value.ToString() == part {
				vals.AddItems(part)
			}
		}
	}
	return dict, vals
}

//	ПОстроение объекта выбора
func (it *DataFancy) FancyBuildChoose(rule *repo.DataRule, table *likbase.ItTable, value lik.Itemer) (lik.Lister,bool) {
	coll := make(map[string]lik.Seter)
	keys := []string{}
	found := false
	for id, it := range table.Elms {
		var name string
		if table.Part == "member" {
			name = jone.CalculateElmString(it, "family") + " " +
					jone.CalculateElmString(it, "namely") + " " +
					jone.CalculateElmString(it, "paterly");
			if len(name) <= 2 {
				name = jone.CalculateElmString(it, "login")
			}
		} else {
			name = jone.CalculateElmText(it)
		}
		item := lik.BuildSet("id", int(id), "name", name)
		coll[name] = item
		keys = append(keys, name)
		if value != nil && value.ToInt() == int(id) { found = true }
	}
	sort.Strings(keys)
	dict := lik.BuildList()
	for _, key := range keys {
		dict.AddItems(coll[key])
	}
	return dict, found
}

//	Запуск вычислений
func (it *DataFancy) RunCalculate(rule *repo.DataRule, info lik.Seter, part string, format string, isform bool) lik.Itemer {
	return it.ItCalculate.Run(rule, info, part, format, isform)
}
func (it *DataFancy) FancyCalculate(rule *repo.DataRule, info lik.Seter, part string, format string, isform bool) lik.Itemer {
	if info == nil {
		return nil
	}
	var value = jone.Calculate(info, part)
	if value == nil {
	} else if strings.Contains(part, "password") {
		value = lik.BuildItem("")
	} else if format == "h" {
		value = lik.BuildItem(value.ToBool())
	} else if format == "d" || format == "t" {
		sval := ""
		if ival := value.ToInt(); ival > 0 {
			layout := "2006/01/02"
			if format == "t" {
				layout += " 15:04"
			}
			sval = time.Unix(int64(ival), 0).Format(layout)
			if match := regexp.MustCompile("^(.*) 00:00").FindStringSubmatch(sval); match != nil {
				sval = match[1]
			}
		}
		value = lik.BuildItem(sval)
	}
	return value
}

//	Поиск поля
func (it *DataFancy) SeekPartField(part string) lik.Seter {
	if it.ListFields != nil {
		for _, field := range it.ListFields {
			if field.GetString("part") == part {
				return field
			}
			if field.GetString("index") == part {
				return field
			}
		}
	}
	return nil
}

//	Поиск элемента
func (it *DataFancy) FindItem(list lik.Lister, val string) (int,lik.Seter) {
	pos := -1
	var item lik.Seter
	if list != nil {
		for ni := 0; ni < list.Count(); ni++ {
			if pot := list.GetSet(ni); pot != nil {
				if strings.HasSuffix(pot.GetString("name"), val) ||
					strings.HasSuffix(pot.GetString("index"), val) {
					pos = ni
					item = pot
					break
				}
			}
		}
	}
	return pos,item
}

//	Отображение результата на клиенте
func (it *DataFancy) ShowResult(rule *repo.DataRule, code lik.Seter) {
	if !code.IsItem("likMain") {
		code.SetItem(it.Main, "likMain")
	}
	if !code.IsItem("likPart") {
		code.SetItem(it.Part, "likPart")
	}
	if !code.IsItem("likZone") {
		code.SetItem(it.Zone, "likZone")
	}
	if !code.IsItem("license") {
		code.SetItem(one.FancyLicense, "license")
	}
	rule.SetResponse(code, "fancy")
}

