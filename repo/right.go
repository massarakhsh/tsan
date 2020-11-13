package repo

import (
	"github.com/massarakhsh/tsan/jone"
	"github.com/massarakhsh/lik"
	"github.com/massarakhsh/lik/likbase"
)

//	Проверка, что объект принадлежит сотруднику
func ProbeItMy(rule *DataRule, part string, id lik.IDB) bool {
	ok := false
	if part == "member" {
		if id == rule.ItSession.IdMember {
			ok = true
		}
	}
	return ok
}

//	Проверка, сто объект принадлежить подразделению сотрудника
func ProbeItDep(rule *DataRule, part string, id lik.IDB) bool {
	ok := false
	if id == rule.ItSession.IdMember {
		ok = true
	} else if depid := rule.GetMember().GetIDB("departid"); depid > 0 {
		if member := jone.TableMember.GetElm(id); member != nil {
			if member.GetIDB("departid") == depid {
				ok = true
			}
		}
	}
	return ok
}

//	Проверка доступности функциии
func (rule *DataRule) RightPartId(part string, id lik.IDB, fun string) bool {
	ok := false
	if elm := jone.GetElm(part, id); elm != nil {
		ok = rule.RightElm(elm, fun)
	}
	return ok
}

//	Проверка доступности функции
func (rule *DataRule) RightElm(elm *likbase.ItElm, fun string) bool {
	if elm == nil { return false }
	ok := false
	if rule.ItSession.RoleMember == jone.ItAdmin { return true }
	if part := elm.Table.Part; part == "bell" {
		ok = rule.RightElmBell(elm, fun)
	} else if part == "offer" {
		ok = rule.RightElmOffer(elm, fun)
	} else if part == "object" {
		ok = rule.RightElmObject(elm, fun)
	} else if part == "client" {
		ok = rule.RightElmClient(elm, fun)
	} else if part == "deal" {
		ok = rule.RightElmDeal(elm, fun)
	} else if part == "depart" {
		ok = rule.RightElmDepart(elm, fun)
	} else if part == "member" {
		ok = rule.RightElmMember(elm, fun)
	}
	return ok
}

//	Проверка функции на контакте
func (rule *DataRule) RightElmBell(elm *likbase.ItElm, fun string) bool {
	ok := false
	return ok
}

//	Проверка функции на заявке
func (rule *DataRule) RightElmOffer(elm *likbase.ItElm, fun string) bool {
	ok := false
	myid := rule.ItSession.IdMember
	myrole := rule.ItSession.RoleMember
	if myrole == jone.ItAdmin { return true }
	mydepart := jone.CalculatePartIdIDB("member", myid, "departid")
	toid := jone.CalculateElmIDB(elm, "memberid")
	todepart := jone.CalculatePartIdIDB("member", toid, "departid")
	if fun == "promo" {
		if myrole == jone.ItManager && mydepart != 0 && mydepart == todepart {
			ok = true
		}
	} else if fun == "edit" {
		if myid == toid {
			ok = true
		} else if myrole == jone.ItManager && mydepart != 0 && mydepart == todepart {
			ok = true
		}
	} else if fun == "" {
	}
	return ok
}

//	Проверка функции на объекте
func (rule *DataRule) RightElmObject(elm *likbase.ItElm, fun string) bool {
	ok := false
	return ok
}

//	Проверка функции на сделке
func (rule *DataRule) RightElmDeal(elm *likbase.ItElm, fun string) bool {
	ok := false
	return ok
}

//	Проверка функции на клиенте
func (rule *DataRule) RightElmClient(elm *likbase.ItElm, fun string) bool {
	ok := false
	return ok
}

//	Проверка функции на подразделении
func (rule *DataRule) RightElmDepart(elm *likbase.ItElm, fun string) bool {
	ok := false
	return ok
}

//	Проверка функции на сотруднике
func (rule *DataRule) RightElmMember(elm *likbase.ItElm, fun string) bool {
	ok := false
	return ok
}

