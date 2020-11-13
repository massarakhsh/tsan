package repo

import (
	"bitbucket.org/961961/tsan/jone"
	"bitbucket.org/shaman/lik"
	"bitbucket.org/shaman/lik/likbase"
)

//	Найти клиента по телефону
func SearchClient(phone string) *likbase.ItElm {
	if pho := jone.NormalizePhone(phone); pho != "" {
		for _, elm := range jone.TableClient.Elms {
			if pho == jone.CalculateElmString(elm,"phone1") {
				return elm
			} else if pho == jone.CalculateElmString(elm,"phone2") {
				return elm
			}
		}
	}
	return nil
}

//	Найти контакт по телефону
func SearchBell(phone string) *likbase.ItElm {
	if pho := jone.NormalizePhone(phone); pho != "" {
		for _, elm := range jone.TableBell.Elms {
			if pho == jone.CalculateElmString(elm,"clientid/phone1") {
				return elm
			} else if pho == jone.CalculateElmString(elm,"clientid/phone2") {
				return elm
			}
		}
	}
	return nil
}

//	Определить актуальный список типов недвижимости
func GetListRealty(rule *DataRule) []lik.Seter {
	var reals []lik.Seter
	if ent := GenDiction.FindEnt("realty"); ent != nil {
		segment := rule.GetMemberParamString("context/segment")
		if list := ent.It.GetList("content"); list != nil {
			for nc := list.Count()-1; nc >= 0; nc-- {
				elm := list.GetSet(nc)
				part := elm.GetString("part")
				if segment == "second" {
					if part != "flat" && part != "room" { continue }
				} else if segment == "new" {
					if part != "flat" && part != "room" { continue }
				}
				reals = append(reals, elm)
			}
		}
	}
	return reals
}

