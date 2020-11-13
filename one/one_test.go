package one_test

import (
	"github.com/massarakhsh/tsan/jone"
	"github.com/massarakhsh/tsan/one"
	"github.com/massarakhsh/lik"
	"github.com/massarakhsh/lik/likbase"
)

func Example() {
	//  Получение записи сотрудника и переименование
	//  подразделения, в которое он входит
	var memberid lik.IDB
	member := one.GetMember(memberid)
	departid := lik.IDB(member.ID)
	if depart := one.GetDepart(departid); depart != nil {
		depart.Update("Name", "Отдел корпоративных клиентов")
	}
	//  Изменение всех или многих полей в записи заявки
	var offerid lik.IDB
	offer := one.GetOffer(offerid)
	offer.MemberId = memberid
	offer.Cost = cost
	//..............
	offer.Save()
}

