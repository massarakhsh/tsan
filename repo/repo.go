// Инструментальный интерфейс базы данных.
//
//	//  Из таблицы заявок получить заявку по ее ID
//	offer := repo.TableOffer.GetElm(id)
//	//  Получить цену заявки как целое число
//	cost := repo.CalculateElmInt(offer, "cost")
//	//  Получить идентификатор ответственного риэлтора как целое число
//	memberid := repo.CalculateElmInt(offer, "memberid")
//	//  Получить этаж из описания объекта как целое число
//	flat := repo.CalculateElmInt(offer, "objectid/define/flat")
//	//  Получить фамилию клиента заявки
//	family := repo.CalculateElmInt(offer, "clientid/family")
package repo

import (
	"bitbucket.org/961961/tsan/jone"
	"bitbucket.org/961961/tsan/one"
	"bitbucket.org/shaman/lik"
	"bitbucket.org/shaman/lik/likbase"
)

//	Запустить и синхронизировать систему
func GoIt(serv string, base string, user string, pass string) {
	SystemInitialize()
	OneSynchronize()
}

//	Синхронизировать систему
func OneSynchronize() {
	go goOneSynchronizeDepart()
	go goOneSynchronizeMember()
}

//	Синхронизировать подразделения
func goOneSynchronizeDepart() {
	if list := one.SelectDepart(nil); list != nil {
		for _, it := range list {
			if elm := jone.TableDepart.GetElm(lik.IDB(it.ID)); elm != nil {
				SynchronizeDepartOne(&it, elm)
			}
		}
	}
}

//	Синхронизировать сотрудников
func goOneSynchronizeMember() {
	if list := one.SelectMember(nil); list != nil {
		for _, it := range list {
			if elm := jone.TableMember.GetElm(lik.IDB(it.ID)); elm != nil {
				SynchronizeMemberOne(&it, elm)
			}
		}
	}
}

//	Синхронизировать запись подразделения
func SynchronizeDepartOne(it *one.Depart, elm *likbase.ItElm) {
	modify := false
	if val := elm.GetString("name"); it.Name != val {
		it.Name = val
		modify = true
	}
	if val := elm.GetString("notes"); it.Notes != val {
		it.Notes = val
		modify = true
	}
	if val := elm.GetInt("departid"); it.UpDepartId != val {
		it.UpDepartId = val
		modify = true
	}
	if modify {
		it.Save()
	}
}

//	Синхронизировать запись сотрудника
func SynchronizeMemberOne(it *one.Member, elm *likbase.ItElm) {
	modify := false
	if val := elm.GetString("family"); it.Family != val {
		it.Family = val
		modify = true
	}
	if val := elm.GetString("namely"); it.Namely != val {
		it.Namely = val
		modify = true
	}
	if val := elm.GetString("paterly"); it.Paterly != val {
		it.Paterly = val
		modify = true
	}
	if val := elm.GetString("phone"); it.Phone != val {
		it.Phone = val
		modify = true
	}
	if val := elm.GetString("prophone"); it.ProPhone != val {
		it.ProPhone = val
		modify = true
	}
	if val := elm.GetString("roles"); it.Role != val {
		it.Role = val
		modify = true
	}
	if val := elm.GetString("email"); it.Email != val {
		it.Email = val
		modify = true
	}
	if val := elm.GetString("photo"); it.Photo != val {
		it.Photo = val
		modify = true
	}
	if val := elm.GetString("pin"); it.Pin != val {
		it.Pin = val
		modify = true
	}
	if val := elm.GetString("notes"); it.Notes != val {
		it.Notes = val
		modify = true
	}
	if val := elm.GetInt("departid"); it.DepartId != val {
		it.DepartId = val
		modify = true
	}
	if modify {
		it.Save()
	}
}

