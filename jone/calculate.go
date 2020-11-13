package jone

import (
	"github.com/massarakhsh/lik"
	"github.com/massarakhsh/lik/likbase"
	"fmt"
	"regexp"
)

//	Получить строку по разделу и индексу
func CalculatePartIdString(part string, id lik.IDB, path string) string {
	return infoToString(CalculatePartId(part, id, path))
}
//	Получить целое по разделу и индексу
func CalculatePartIdInt(part string, id lik.IDB, path string) int {
	return infoToInt(CalculatePartId(part, id, path))
}
//	Получить индекс по разделу и индексу
func CalculatePartIdIDB(part string, id lik.IDB, path string) lik.IDB {
	return infoToIDB(CalculatePartId(part, id, path))
}
//	Получить булевское по разделу и индексу
func CalculatePartIdBool(part string, id lik.IDB, path string) bool {
	return infoToBool(CalculatePartId(part, id, path))
}
//	Получить список по разделу и индексу
func CalculatePartIdList(part string, id lik.IDB, path string) lik.Lister {
	return infoToList(CalculatePartId(part, id, path))
}
//	Получить структуру по разделу и индексу
func CalculatePartIdSet(part string, id lik.IDB, path string) lik.Seter {
	return infoToSet(CalculatePartId(part, id, path))
}
//	Получить и транслировать строку по разделу и индексу
func CalculatePartIdTranslate(part string, id lik.IDB, path string) string {
	return CalculateElmTranslate(GetElm(part, id), path)
}
//	Получить интерфейс по разделу и индексу
func CalculatePartId(part string, id lik.IDB, path string) lik.Itemer {
	return CalculateElm(GetElm(part, id), path)
}

//	Получить строку по объекту
func CalculateElmString(elm *likbase.ItElm, path string) string {
	return infoToString(CalculateElm(elm, path))
}
//	Получить целое по объекту
func CalculateElmInt(elm *likbase.ItElm, path string) int {
	return infoToInt(CalculateElm(elm, path))
}
//	Получить булевское по объекту
func CalculateElmBool(elm *likbase.ItElm, path string) bool {
	return infoToBool(CalculateElm(elm, path))
}
//	Получить индекс по объекту
func CalculateElmIDB(elm *likbase.ItElm, path string) lik.IDB {
	return infoToIDB(CalculateElm(elm, path))
}
//	Получить плавающее по объекту
func CalculateElmFloat(elm *likbase.ItElm, path string) float64 {
	return infoToFloat(CalculateElm(elm, path))
}
//	Получить структуру по объекту
func CalculateElmSet(elm *likbase.ItElm, path string) lik.Seter {
	return infoToSet(CalculateElm(elm, path))
}
//	Получить список по объекту
func CalculateElmList(elm *likbase.ItElm, path string) lik.Lister {
	return infoToList(CalculateElm(elm, path))
}
//	Получить и транслировать строку по объекту
func CalculateElmTranslate(elm *likbase.ItElm, path string) string {
	value := ""
	if elm == nil {
	} else if path == "id" || path == "/id" {
		value = likbase.IDBToStr(elm.Id)
	} else if elm.Info != nil {
		value = CalculateTranslate(elm.Info, path)
	}
	return value
}

//	Получить строку на интерфейсе
func CalculateString(info lik.Itemer, path string) string {
	return infoToString(Calculate(info, path))
}
//	Получить целое на интерфейсе
func CalculateInt(info lik.Itemer, path string) int {
	return infoToInt(Calculate(info, path))
}
//	Получить индекс на интерфейсе
func CalculateIDB(info lik.Itemer, path string) lik.IDB {
	return infoToIDB(Calculate(info, path))
}
//	Получить булевское на интерфейсе
func CalculateBool(info lik.Itemer, path string) bool {
	return infoToBool(Calculate(info, path))
}
//	Получить список на интерфейсе
func CalculateList(info lik.Itemer, path string) lik.Lister {
	return infoToList(Calculate(info, path))
}
//	Получить структуру на интерфейсе
func CalculateSet(info lik.Itemer, path string) lik.Seter {
	return infoToSet(Calculate(info, path))
}
//	Получить и транслировать строку на интерфейсе
func CalculateTranslate(info lik.Itemer, path string) string {
	return infoToTranslate(Calculate(info, path), path)
}

//	Получить интерфейс на объекте
func CalculateElm(elm *likbase.ItElm, path string) lik.Itemer {
	var value lik.Itemer
	if elm == nil {
	} else if path == "id" || path == "/id" {
		value = lik.BuildItem(elm.Id)
	} else if elm.Info != nil {
		value = Calculate(elm.Info, path)
	}
	return value
}

//	Получить номер объекта
func CalculateElmSid(elm *likbase.ItElm) string {
	value := ""
	if elm != nil {
		value += fmt.Sprintf("%03d", int(elm.Id))
	}
	return value
}

//	Получить текст на объекте
func CalculateElmText(elm *likbase.ItElm) string {
	value := ""
	if elm == nil {
	} else if part := elm.Table.Part; part == "object" {
		value = DefinitionObject(elm)
	} else if part == "bell" {
		value = DefinitionBell(elm)
	} else if part == "offer" {
		value = DefinitionOffer(elm)
	} else if part == "member" {
		value = DefinitionMember(elm)
	} else if part == "depart" {
		value = DefinitionDepart(elm)
	} else if part == "client" {
		value = DefinitionClient(elm)
	}
	if value == "" && elm != nil {
		value = fmt.Sprintf("#%d", int(elm.Id))
	}
	return value
}

//	Получить интерфейс на интерфейсе
func Calculate(info lik.Itemer, path string) lik.Itemer {
	var value lik.Itemer
	if info == nil {
	} else if name,ext := lik.GetFirstExt(path); name == "" && ext == "" {
		value = info
	} else if name == "" {
		value = Calculate(info, ext)
	} else if iset := info.ToSet(); iset != nil {
		if item := iset.GetItem(name); item != nil {
			if ext == "" {
				value = item
			} else if item.IsSet() {
				value = Calculate(item, ext)
			} else if match := regexp.MustCompile("^(.+)id$").FindStringSubmatch(name); match == nil {
				value = Calculate(item, ext)
			} else if table := GetTable(match[1]); table == nil {
				value = Calculate(item, ext)
			} else if elm := table.GetElm(lik.IDB(item.ToInt())); elm != nil {
				value = CalculateElm(elm, ext)
			}
		}
	} else if ilist := info.ToList(); ilist != nil {
		if ival,ok := lik.StrToIntIf(name); ok {
			if item := ilist.GetItem(ival); item != nil {
				if ext == "" {
					value = item
				} else {
					value = Calculate(item, ext)
				}
			}
		} else if match := regexp.MustCompile("^idx(\\d*)$").FindStringSubmatch(name); match != nil {
			idx := lik.StrToInt(match[1])
			if nl,item := ScanItemIdx(ilist, idx); nl >= 0 {
				if ext == "" {
					value = item
				} else {
					value = Calculate(item, ext)
				}
			}
		}
	}
	return value
}

//	Посчитать текст на разделе и идентификаторе
func CalculatePartIdText(part string, id lik.IDB) string {
	value := ""
	if elm := GetElm(part, id); elm != nil {
		value = CalculateElmText(elm)
	}
	return value
}

//	Преобразовать интерфейс в строку
func infoToString(info lik.Itemer) string {
	if info != nil {
		return info.ToString()
	} else {
		return ""
	}
}

//	Преобразовать интерфейс в целое
func infoToInt(info lik.Itemer) int {
	if info != nil {
		return info.ToInt()
	} else {
		return 0
	}
}

//	Преобразовать интерфейс в действительное
func infoToFloat(info lik.Itemer) float64 {
	if info != nil {
		return info.ToFloat()
	} else {
		return 0
	}
}

//	Преобразовать интерфейс в идентификатор
func infoToIDB(info lik.Itemer) lik.IDB {
	if info != nil {
		return lik.IDB(info.ToInt())
	} else {
		return 0
	}
}

//	Преобразовать интерфейс в булевское
func infoToBool(info lik.Itemer) bool {
	if info != nil {
		return info.ToBool()
	} else {
		return false
	}
}

//	Преобразовать интерфейс в список
func infoToList(info lik.Itemer) lik.Lister {
	if info != nil {
		return info.ToList()
	} else {
		return nil
	}
}

//	Преобразовать интерфейс в структуру
func infoToSet(info lik.Itemer) lik.Seter {
	if info != nil {
		return info.ToSet()
	} else {
		return nil
	}
}

//	Преобразовать и транслировать строку
func infoToTranslate(info lik.Itemer, path string) string {
	value := infoToString(info)
	if match := lik.RegExParse(path,"([^/]+)$"); match != nil {
		value = SystemStringTranslate(match[1], value)
	}
	return value
}

//	Установить элемент по индексу
func ScanItemIdx(list lik.Lister, idx int) (int,lik.Seter) {
	if list != nil {
		if ns := list.Count(); ns > 0 {
			for nl := 0; nl < ns; nl++ {
				if item := list.GetSet(nl); item != nil && idx == item.GetInt("idx") {
					return nl,item
				}
			}
		}
	}
	return -1,nil
}

