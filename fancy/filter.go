package fancy

import (
	"github.com/massarakhsh/tsan/repo"
	"github.com/massarakhsh/lik"
	"fmt"
	"math/rand"
	"strings"
)

//	Дескриптор элемента фильтра
type DataFilterCond struct {
	Key, Opr, Val	string
}

//	Дескриптор фильтра
type DataFilter struct {
	ItSegment string	//	Сегмент рынка
	ItRealty  string	//	Тип объекта
	ItLocate  string	//	Локализация
	ItStatus  string	//	Статус
	Search    string	//	Поиск по контексту
	Conds     []DataFilterCond	//	Список элементов условия
	Total		int		//	Всего элементоа
	Start		int		//	Первый элемент
	Limit		int		//	Количество элементов
	Page		int		//	Размер страницы
	SortKey   string	//	Поле сортировки
	SortDir   bool		//	Направление сортировки
}

// Поиск фильтра в коллекции.
// part - ключ фильтра, name - наименование фильтра
//
// Результат - индекс в списке и интерфейс искомого фильтра
func (it *DataFancy) FancyFilterFind(rule *repo.DataRule, part string, name string) (int,lik.Seter) {
	if part != "" && part != "all" || name != "" {
		if filters := rule.GetMemberParamList(it.GetParter()+"/filters"); filters != nil {
			for nf := 0; nf < filters.Count(); nf++ {
				if filter := filters.GetSet(nf); filter != nil {
					if part != "" && part != filter.GetString("part") {
						continue
					}
					if name != "" && name != filter.GetString("name") {
						continue
					}
					return nf, filter
				}
			}
		}
	}
	return -1,nil
}

// Установка на текущий фильтр, загрузка при необходимости.
//
// Результат - индекс в списке
func (it *DataFancy) FancyFilterSeek(rule *repo.DataRule) int {
	nf := -1
	if it.RuleFilter != nil {
		part := it.RuleFilter.GetString("part")
		nf,_ = it.FancyFilterFind(rule, part,"")
	} else {
		part := rule.GetMemberParamString(it.GetParter()+"/filter")
		nf = it.FancyFilterSet(rule, part, false)
	}
	return nf
}

// Установка указанного фильтра.
//
// part - ключ фильтра, reset - признак сброса
func (it *DataFancy) FancyFilterSet(rule *repo.DataRule, part string, reset bool) int {
	nf,filter := it.FancyFilterFind(rule, part,"")
	if filter != nil {
		it.RuleFilter = filter.Clone().ToSet()
		rule.SetMemberParam(filter.GetString("part"), it.GetParter()+"/filter")
		rule.SetMemberParam(filter.GetString("segment"), "context/segment")
		rule.SetMemberParam(filter.GetString("realty"), "context/realty")
		rule.SetMemberParam(filter.GetString("locate"), "context/locate")
		rule.SetMemberParam(filter.GetString("status"), "context/status")
	} else {
		it.RuleFilter = lik.BuildSet()
		rule.SetMemberParam(nil, it.GetParter()+"/filter")
		if reset {
			rule.SetMemberParam(nil, "context/segment")
			rule.SetMemberParam(nil, "context/realty")
			rule.SetMemberParam(nil, "context/locate")
			rule.SetMemberParam(nil, "context/status")
		}
	}
	return nf
}

// Фиксация состояния фильтра
func (it *DataFancy) FancyFilterFix(rule *repo.DataRule, data lik.Seter) {
	if it.RuleFilter == nil {
		it.RuleFilter = lik.BuildSet()
	}
	segment := rule.GetMemberParamString("context/segment")
	realty := rule.GetMemberParamString("context/realty")
	locate := rule.GetMemberParamString("context/locate")
	active := rule.GetMemberParamString("context/status")
	it.RuleFilter.SetItem(segment, "segment")
	it.RuleFilter.SetItem(realty, "realty")
	it.RuleFilter.SetItem(locate, "locate")
	it.RuleFilter.SetItem(active, "status")
	if data != nil {
		filter := lik.BuildList()
		it.RuleFilter.SetItem(filter, "filter")
		if filters := lik.ListFromRequest(data.GetString("filter")); filters != nil {
			for nf := 0; nf < filters.Count(); nf++ {
				if set := filters.GetSet(nf); set != nil {
					key := strings.ReplaceAll(set.GetString("property"), "__", "/")
					opr := set.GetString("operator")
					if opr == "eq" {
						opr = "="
					} else if opr == "ne" {
						opr = "!="
					} else if opr == "gt" {
						opr = ">"
					} else if opr == "ge" {
						opr = ">="
					} else if opr == "lt" {
						opr = "<"
					} else if opr == "le" {
						opr = "<="
					} else if opr == "or" {
						opr = "|"
					} else if opr == "like" {
						opr = "*"
					}
					val := set.GetString("value")
					filter.AddItemSet("key", key, "opr", opr, "val", val)
				}
			}
		}
		sort := data.GetString("sort")
		it.RuleFilter.SetItem(sort, "sort")
		dir := data.GetString("dir") == "asc"
		it.RuleFilter.SetItem(dir, "dir")
		page := data.GetInt("page")
		it.RuleFilter.SetItem(page, "page")
		start := data.GetInt("start")
		it.RuleFilter.SetItem(start, "start")
		limit := data.GetInt("limit")
		it.RuleFilter.SetItem(limit, "limit")
	}
}

//	Сохранение фильтра
func (it *DataFancy) FancyFilterSave(rule *repo.DataRule) {
	if it.FancyFilterSeek(rule) >= 0 {
		name := it.RuleFilter.GetString("name")
		it.FancyFilterSaveAs(rule, name)
	}
}

//	Сокранения фильтра по имени
func (it *DataFancy) FancyFilterSaveAs(rule *repo.DataRule, name string) {
	it.FancyFilterFix(rule,nil)
	filter := it.RuleFilter.Clone().ToSet()
	list_filters := rule.GetMemberParamList(it.GetParter()+"/filters")
	if list_filters == nil {
		list_filters = lik.BuildList()
		rule.SetMemberParam(list_filters, it.GetParter()+"/filters")
	}
	if name == "" { name = "Новый" }
	filter.SetItem(name, "name")
	if nf,flt := it.FancyFilterFind(rule,"", name); flt != nil {
		filter.SetItem(flt.GetString("part"),"part")
		list_filters.SetItem(filter, nf)
	} else {
		filter.SetItem(fmt.Sprintf("f%09d", rand.Int31()), "part")
		list_filters.AddItems(filter)
	}
	rule.SaveMemberParam()
	it.FancyFilterSet(rule, filter.GetString("part"), true)
}

//	Удаление фильтра
func (it *DataFancy) FancyFilterDelete(rule *repo.DataRule) {
	if nf := it.FancyFilterSeek(rule); nf >= 0 {
		if list_filters := rule.GetMemberParamList(it.GetParter()+"/filters"); list_filters != nil {
			list_filters.DelItem(nf)
		}
		rule.SetMemberParam(nil, it.GetParter()+"/filter")
	}
	it.FancyFilterSet(rule, "", true)
}

//	Разбор и подготовка фильтра
func (it *DataFancy) FancyFilterDecode(rule *repo.DataRule) DataFilter {
	data := DataFilter{}
	if it.RuleFilter != nil {
		data.ItSegment = it.RuleFilter.GetString("segment")
		data.ItRealty = it.RuleFilter.GetString("realty")
		data.ItLocate = it.RuleFilter.GetString("locate")
		data.ItStatus = it.RuleFilter.GetString("status")
		data.Start = it.RuleFilter.GetInt("start")
		data.Limit = it.RuleFilter.GetInt("limit")
		data.Page = it.RuleFilter.GetInt("page")
		data.Total = it.RuleFilter.GetInt("total")
		if filter := it.RuleFilter.GetList("filter"); filter != nil {
			for nc := 0; nc < filter.Count(); nc++ {
				if set := filter.GetSet(nc); set != nil {
					key := set.GetString("key")
					opr := set.GetString("opr")
					val := strings.ToLower(set.GetString("val"))
					data.Conds = append(data.Conds, DataFilterCond{key, opr, val})
				}
			}
		}
		data.SortKey = it.RuleFilter.GetString("sort")
		data.SortDir = it.RuleFilter.GetBool("dir")
	}
	return data
}

