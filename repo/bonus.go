package repo

import (
	"github.com/massarakhsh/tsan/jone"
	"github.com/massarakhsh/lik"
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

