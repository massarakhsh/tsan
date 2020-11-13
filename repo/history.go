package repo

import (
	"github.com/massarakhsh/tsan/jone"
	"github.com/massarakhsh/lik"
	"github.com/massarakhsh/lik/likbase"
	"time"
)

//	Получить новый идентификатор события
func ReserveHistory(elm *likbase.ItElm) int {
	histidx := jone.CalculateElmInt(elm,"histidx") + 1
	jone.SetElmValue(elm, histidx, "histidx")
	return histidx
}

//	Добавить событие
func AddHistorySet(rule *DataRule, elm *likbase.ItElm, values ...interface{}) {
	AddHistory(rule, elm, lik.BuildSet(values...))
}

//	Добавить событие
func AddHistory(rule *DataRule, elm *likbase.ItElm, data lik.Seter) {
	histidx := jone.CalculateElmInt(elm,"histidx")
	history := jone.CalculateElmList(elm,"history")
	if history == nil {
		history = lik.BuildList()
		jone.SetElmValue(elm, history, "history")
	}
	if data != nil {
		if data.GetInt("date") <= 0 {
			data.SetItem(int(time.Now().Unix()), "date")
		}
		if data.GetIDB("memberid") <= 0 {
			data.SetItem(rule.ItSession.IdMember, "memberid")
		}
		if data.GetInt("idx") <= 0 {
			histidx++
			jone.SetElmValue(elm, histidx, "histidx")
			data.SetItem(histidx, "idx")
		}
		history.AddItems(data)
		elm.OnModify()
		SortHistory(elm)
	}
}

//	Получить список событий
func GetHistory(elm *likbase.ItElm, what string) []lik.Seter {
	return ExtractHistory(jone.CalculateElmList(elm,"history"), what)
}

//	Получить список событий
func ExtractHistory(history lik.Lister, what string) []lik.Seter {
	list := []lik.Seter{}
	if history != nil {
		for n := 0; n < history.Count(); n++ {
			if hist := history.GetSet(n); hist != nil {
				if what == "" || what == "history" || hist.GetString("what") == what {
					list = append(list, hist)
				}
			}
		}
	}
	return list
}

//	Сортировать события
func SortHistory(elm *likbase.ItElm) {
	if history := jone.CalculateElmList(elm,"history"); history != nil {
		if hs := history.Count(); hs > 1 {
			for hb := hs - 1; hb > 0; hb-- {
				for ha := 0; ha < hb; ha++ {
					if he2 := history.GetSet(ha + 1); he2 != nil {
						if he1 := history.GetSet(ha); he1 != nil {
							if he1.GetInt("date") > he2.GetInt("date") {
								history.SwapItem(ha, ha+1)
								elm.OnModify()
							}
						} else {
							history.SwapItem(ha, ha+1)
							elm.OnModify()
						}
					}
				}
			}
			for ha := hs - 1; ha >= 0; ha-- {
				if he := history.GetSet(ha); he == nil || he.GetInt("date") == 0 {
					history.DelItem(ha)
					elm.OnModify()
				} else {
					break
				}
			}
		}
	}
}

