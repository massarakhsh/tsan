package jone

import (
	"github.com/massarakhsh/lik"
	"github.com/massarakhsh/lik/likbase"
)

//	Изменение объекта по разделу и индексу
//	part, id - раздел и индекс
//	val - устанавливаемое значение
//	path - путь изменения
func SetPartIdValue(part string, id lik.IDB, val interface{}, path string) {
	if elm := GetElm(part, id); elm != nil {
		SetElmValue(elm, val, path)
	}
}

//	Изменение объекта по указателю
//	elm - указатель на объект
//	val - устанавливаемое значение
//	path - путь изменения
func SetElmValue(elm *likbase.ItElm, val interface{}, path string) {
	if elm == nil {
		return
	}
	if elm.Info == nil {
		elm.Info = lik.BuildSet()
		elm.OnModify()
	}
	if SetInfoValue(elm.Info, val, path) {
		elm.OnModify()
	}
}

//	Изменение объекта по интерфейсу
//	info - интерфейс
//	val - устанавливаемое значение
//	path - путь изменения
func SetInfoValue(info lik.Itemer, val interface{}, path string) bool {
	modify := false
	if info == nil {
	} else if name, ext := lik.GetFirstExt(path); name == "" && ext == "" {
	} else if name == "" {
		if SetInfoValue(info, val, ext) {
			modify = true
		}
	} else if imap := info.ToSet(); imap != nil {
		if ext == "" {
			if imap.SetItem(val, name) {
				modify = true
			}
		} else if match := lik.RegExParse(name, "^(.*)id$"); match != nil {
			part := match[1]
			id := lik.IDB(0)
			if item := imap.GetItem(name); item != nil && item.IsSet() {
				if SetInfoValue(item, val, ext) {
					modify = true
				}
			} else if table := GetTable(part); table != nil {
				if item != nil {
					id = lik.IDB(item.ToInt())
				}
				if id == 0 {
					elm := table.CreateElm()
					id = elm.Id
					imap.SetItem(id, name)
					modify = true
				}
				if elm := table.GetElm(id); elm != nil {
					SetElmValue(elm, val, ext)
				}
			}
		} else if item := imap.GetItem(name); item != nil {
			if SetInfoValue(item, val, ext) {
				modify = true
			}
		} else if lik.RegExCompare(ext, "^\\d+") {
			modify = true
			item := lik.BuildList()
			imap.SetItem(item, name)
			SetInfoValue(item, val, ext)
		} else if lik.RegExCompare(ext, "^idx(\\d*)") {
			modify = true
			item := lik.BuildList()
			imap.SetItem(item, name)
			SetInfoValue(item, val, ext)
		} else {
			modify = true
			item := lik.BuildSet()
			imap.SetItem(item, name)
			SetInfoValue(item, val, ext)
		}
	} else if ilist := info.ToList(); ilist != nil {
		if num, ok := lik.StrToIntIf(name); ok {
			if ext == "" {
				if ilist.SetItem(val, num) {
					modify = true
				}
			} else if item := ilist.GetItem(num); item != nil {
				if SetInfoValue(item, val, ext) {
					modify = true
				}
			} else if lik.RegExCompare(ext, "^\\d+") {
				modify = true
				item := lik.BuildList()
				ilist.SetItem(item, num)
				SetInfoValue(item, val, ext)
			} else if lik.RegExCompare(ext, "^idx(\\d*)") {
				modify = true
				item := lik.BuildList()
				ilist.SetItem(item, num)
				SetInfoValue(item, val, ext)
			} else {
				modify = true
				item := lik.BuildSet()
				ilist.SetItem(item, num)
				SetInfoValue(item, val, ext)
			}
		} else if match := lik.RegExParse(name, "^idx(\\d*)$"); match != nil {
			idx := lik.StrToInt(match[1])
			num, item := ScanItemIdx(ilist, idx)
			if num >= 0 && ext == "" {
				if ilist.SetItem(val, num) {
					modify = true
				}
			} else if ext == "" {
				if ilist.SetItem(val, 0) {
					modify = true
				}
			} else if item != nil {
				if SetInfoValue(item, val, ext) {
					modify = true
				}
			} else {
				modify = true
				item := lik.BuildSet("idx", idx)
				ilist.AddItems(item)
				SetInfoValue(item, val, ext)
			}
		}
	}
	return modify
}

