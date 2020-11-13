package repo

import (
	"bitbucket.org/961961/tsan/jone"
	"bitbucket.org/shaman/lik"
	"fmt"
)

//	Посчитать бонусы сотрудника
func CalculeBonus(id lik.IDB) int {
	what := "SUM(bonuses_list.value)"
	from := "bonuses INNER JOIN bonuses_list ON bonuses.bonuses_list_id=bonuses_list.id"
	where := fmt.Sprintf("bonuses.members_id=%d", int(id))
	sql := jone.DB.PrepareSql(what, from, where, "")
	summ,_ := jone.DB.CalculeInt(sql)
	return summ
}

