package repo

import (
	"github.com/massarakhsh/tsan/jone"
	"github.com/massarakhsh/lik"
	"github.com/massarakhsh/lik/likbase"
)

const (
	GEN_Total = 5
)

//	Дескриптор раздела
type SysGen struct {
	Gen		string		//	Ключ раздела
	Name	string		//	Имя раздела
	Ents	[]*SysEnt	//	Список сущностей
}

//	Дескриптор сущности
type SysEnt struct {
	Gen		*SysGen		//	Указатель на раздел
	Key		string		//	Ключ сущности
	It		*likbase.ItElm	//	Указатель на объект
	oldKey	string		//	Старый ключ при замене
}

//	Дескриптор элемента сущности
type Element struct {
	Part string			//	Ключ элемента
	Name string			//	Наименование элемента
}

var (
	SysGens	[GEN_Total]SysGen	//	Список корневых разделов
	SysIndex int				//

	KeyTable = "table"
	GenTable = &SysGens[0]
	KeyParam = "param"
	GenParam = &SysGens[1]
	KeyStruct = "struct"
	GenStruct = &SysGens[2]
	KeyDiction = "diction"
	GenDiction = &SysGens[3]
	KeyExtern = "extern"
	GenExtern = &SysGens[4]
)

var ListRootEnt = []Element{
	{KeyTable, "Таблица"},
	{KeyParam, "Параметры"},
	{KeyStruct, "Структура"},
	{KeyDiction, "Словарь"},
	{KeyExtern, "Окружение"},
}

var ListFormat = lik.BuildSet(
	"s=Строка",
	"c=Словарь",
	"b=Да/Нет",
	"n=Число",
	"m=Деньги",
	"p=Телефон",
	"d=Дата",
	"a=Текст",
	"t=Время",
	"g=Изображение",
	"h=Чекбокс",
	"l=Множество",
	"r=Дерево",
	// u - неопределен
)

//	Инициализация системы
func SystemInitialize() {
	SysIndex++
	SystemInitGen()
	SystemInitLoad()
	SystemInitTable()
	SystemLoadTrans()
}

//	Инициализация разделов
func SystemInitGen() {
	for ng,root := range(ListRootEnt) {
		SysGens[ng].Gen = root.Part
		SysGens[ng].Name = root.Name
		SysGens[ng].Ents = []*SysEnt{}
	}
}

//	Загрузка параметров
func SystemInitLoad() {
	all := jone.GetTable("param").Elms
	for id,it := range(all) {
		gen := jone.CalculateElmString(it,"gen")
		key := jone.CalculateElmString(it,"key")
		itgen := SystemFindGen(gen)
		if itgen == nil {
			key = ""
		} else if itgen.FindEnt(key) != nil {
			key = ""
		}
		if key != "" {
			itgen.BuildEnt(key, it)
		} else {
			jone.TableParam.DeleteElm(id)
		}
	}
}

//	Инициализация таблиц
func SystemInitTable() {
	for _,tbl := range(jone.ListTables) {
		key := tbl.Part
		name := tbl.Title
		ent := GenTable.FindEnt(key)
		modify := false
		if ent == nil {
			it := jone.TableParam.CreateElm()
			jone.SetElmValue(it, KeyTable,"gen")
			jone.SetElmValue(it, key,"key")
			jone.SetElmValue(it, name,"name")
			ent = GenTable.BuildEnt(key, it)
			modify = true
		}
		if modify {
			ent.SaveToBase()
		}
	}
}

//	Ариск раздела
func SystemFindGen(gen string) *SysGen {
	for ng := 0; ng < GEN_Total; ng++ {
		if SysGens[ng].Gen == gen { return &SysGens[ng] }
	}
	return nil
}

//	Поиск раздела и сущности
func SystemFindGenEnt(nmgen string, key string) *SysEnt {
	if gen := SystemFindGen(nmgen); gen != nil {
		return gen.FindEnt(key)
	}
	return nil
}

//	Поиск раздела
func (gen *SysGen) FindEnt(key string) *SysEnt {
	for _,ent := range (gen.Ents) {
		if ent.Key == key {
			return ent
		}
	}
	return nil
}

//	Поиск позиции в разделе
func (gen *SysGen) FindPart(key string, part string) (int, lik.Seter) {
	if itent := gen.FindEnt(key); itent != nil {
		return itent.FindPart(part)
	}
	return -1, nil
}

//	Выпор позиции по номеру
func (gen *SysGen) FindPartPos(key string, pos int) lik.Seter {
	if itent := gen.FindEnt(key); itent != nil {
		return itent.FindPartPos(pos)
	}
	return nil
}

//	Проверка и выбор уникального номера
func (gen *SysGen) GenUniqueEnt(key string) string {
	result := key
	if result == "" {
		result = "_"
	}
	for {
		if ent := gen.FindEnt(result); ent == nil { break }
		result += "_"
	}
	return result
}

//	Конструктор сущности
func (gen *SysGen) BuildEnt(key string, it *likbase.ItElm) *SysEnt {
	keyto := key
	if key != "" {
		keyto = gen.GenUniqueEnt(key)
	}
	ent := &SysEnt{ Gen: gen, Key: keyto, It: it, oldKey: key }
	gen.Ents = append(gen.Ents, ent)
	return ent
}

//	Найти =позицию в сущности
func (ent *SysEnt) FindPart(part string) (int, lik.Seter) {
	content := ent.GetContent()
	for ne,elm := range (content) {
		if elm.GetString("part") == part {
			return ne, elm
		}
	}
	return -1, nil
}

//	Найти позицию по номеру
func (ent *SysEnt) FindPartPos(pos int) lik.Seter {
	content := ent.GetContent()
	if pos >= 0 && pos < len(content) {
		return content[pos]
	}
	return nil
}

//	ПОлучить контекст сущности
func (ent *SysEnt) GetContent() []lik.Seter {
	content := []lik.Seter{}
	if list := jone.CalculateElmList(ent.It,"content"); list != nil {
		for ne := 0; ne < list.Count(); ne++ {
			if elm := list.GetSet(ne); elm != nil {
				content = append(content, elm)
			}
		}
	}
	return content
}

//	Сохранить сущность
func (ent *SysEnt) SaveToBase() string {
	if jone.CalculateElmString(ent.It,"gen") != ent.Gen.Gen {
		jone.SetElmValue(ent.It, ent.Gen.Gen,"gen")
	}
	key := ent.Key
	if key == "" {
		key = ent.Gen.GenUniqueEnt(key)
	} else if entkey := ent.Gen.FindEnt(key); entkey != nil && entkey != ent {
		key = ent.Gen.GenUniqueEnt(key)
	}
	ent.Key = key
	if jone.CalculateElmString(ent.It, "key") != key {
		jone.SetElmValue(ent.It, key,"key")
	}
	ent.It.OnModify()
	SystemLoadTrans()
	SysIndex++
	return "/" + ent.Gen.Gen + "/" + ent.Key
}

//	Удалить сущность
func (ent *SysEnt) DeleteFromBase() {
	ent.It.Delete()
	for ne := 0; ne < len(ent.Gen.Ents); ne++ {
		if ent.Gen.Ents[ne].Key == ent.Key {
			ent.Gen.Ents = append(ent.Gen.Ents[:ne], ent.Gen.Ents[ne+1:]...)
		}
	}
	SystemLoadTrans()
}

//	Выполить системную команду
func SystemExecuteCmd(cmd string) {
	if cmd == "reload" {
		SystemInitialize()
	} else if cmd == "clearbase" {
		jone.TableObject.Drop()
		jone.TableOffer.Drop()
		jone.TableBell.Drop()
		jone.TableDeal.Drop()
		jone.TableClient.Drop()
	} else if cmd == "clearparam" {
		jone.TableParam.Purge()
		SystemInitialize()
	} else if cmd == "loaddump" {
		LoadDump()
	}
}

//	Сохранить базу данныъ
func SystemSaveDataBase() lik.Seter {
	return SystemStoreBase(nil)
}

//	Восстановиь бызу данных
func SystemLoadDataBase(dump lik.Seter) {
	SystemRestoreBase(nil, dump)
	jone.StartBase()
	SystemInitialize()
}

//	Сохранить базу данных
func SystemStoreBase(tables []string) lik.Seter {
	dumpbase := lik.BuildSet()
	for _,table := range(jone.ListTables) {
		need := false
		if tables == nil {
			need = true
		} else {
			for _, part := range (tables) {
				if table.Part == part {
					need = true
					break
				}
			}
		}
		if need {
			dumptable := lik.BuildList()
			elms := table.GetListElm(true)
			for _,elm := range elms {
				dumptable.AddItemSet("id", elm.Id, "info", elm.Info)
			}
			dumpbase.SetItem(dumptable, table.Part)
		}
	}
	return dumpbase
}

//	Восстановить базу данных
func SystemRestoreBase(tables []string, dumpbase lik.Seter) {
	if tables == nil {
		tables = []string{}
		for _,table := range(jone.ListTables) {
			tables = append(tables, table.Part)
		}
	}
	for _,part := range(tables) {
		if table := jone.GetTable(part); table != nil {
			if dumptable := dumpbase.GetList(part); dumptable != nil {
				table.Drop()
				for ne := 0; ne < dumptable.Count(); ne++ {
					dumpelm := dumptable.GetSet(ne)
					id := dumpelm.GetIDB("id")
					info := dumpelm.GetSet("info")
					info.SetItem(nil, "idu")
					table.RestoreElm(id, info)
				}
			}
		}
	}
	SysIndex++
}

//	Сортировать список сущностей
func SystemSortEnts(gen string) []*likbase.ItElm {
	result := []*likbase.ItElm {}
	if gen := SystemFindGen(gen); gen != nil {
		for _, ent := range gen.Ents {
			result = append(result, ent.It)
		}
		me := len(result)
		for le := me - 2; le >= 0; le-- {
			for fe := 0; fe <= le; fe++ {
				if result[fe].GetString("name") > result[fe+1].GetString("name") {
					res := result[fe]
					result[fe] = result[fe+1]
					result[fe+1] = res
				}
			}
		}
	}
	return result
}

//	Загрузить таблицу трансляций
func SystemLoadTrans() {
	jone.Trans = make(map[string]string)
	for _,ent := range GenDiction.Ents {
		if content := ent.GetContent(); content != nil {
			key := ent.Key
			for _,pot := range content {
				part := pot.GetString("part")
				name := pot.GetString("name")
				jone.Trans[key+"+"+part] = name
			}
		}
	}
}

