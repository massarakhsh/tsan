package member

import (
	"github.com/massarakhsh/tsan/control"
	"github.com/massarakhsh/tsan/fancy"
	"github.com/massarakhsh/tsan/jone"
	"github.com/massarakhsh/tsan/repo"
	"github.com/massarakhsh/tsan/show"
	"github.com/massarakhsh/lik"
	"github.com/massarakhsh/lik/likdom"
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

type BonusControl struct {
	control.DataControl
	fancy.TableFancy
}

type dealBonusExecute struct {
	It	*BonusControl
}
func (it *dealBonusExecute) Run(rule *repo.DataRule, cmd string, data lik.Seter) {
	it.It.DoExecute(rule, cmd, data)
}

type dealBonusGridFill struct {
	It	*BonusControl
}
func (it *dealBonusGridFill) Run(rule *repo.DataRule) {
	it.It.DoGridFill(rule)
}

type dealBonusPageFill struct {
	It	*BonusControl
}
func (it *dealBonusPageFill) Run(rule *repo.DataRule) lik.Lister {
	return it.It.DoPageFill(rule)
}

func BuildBonus(rule *repo.DataRule, main string, id lik.IDB) *BonusControl {
	it := &BonusControl{ }
	it.ControlInitializeZone(main, id, "bonus")
	it.TableInitialize(rule,it.Frame,"member", it.Mode)
	it.ItExecute = &dealBonusExecute{it}
	it.ItGridFill = &dealBonusGridFill{it}
	it.ItPageFill = &dealBonusPageFill{it}
	return it
}

func (it *BonusControl) BuildShow(rule *repo.DataRule, sx int, sy int) likdom.Domer {
	it.SetSize(sx, sy)
	return show.BuildFancyGrid(it.Main,"bonus")
}

func (it *BonusControl) DoExecute(rule *repo.DataRule, cmd string, data lik.Seter) {
	if cmd == "bonus" {
		it.DoExecute(rule, rule.Shift(), data)
	} else {
		it.TableExecute(rule, cmd, data)
	}
}

func (it *BonusControl) DoGridFill(rule *repo.DataRule) {
	it.TableGridFill(rule)
	it.GridBuildColumns(rule, true, true)
	it.Grid.SetSize(it.Grid.Width, it.Grid.Height)
	title := fmt.Sprintf("Текущий баланс бонусов: %d", repo.CalculeBonus(it.IdMain))
	it.AddCommandItem(rule, 80, lik.BuildSet("type=text", "text", title))
	it.AddCommandImg(rule, 910, "Открыть", "toshow", "show")
}

func (it *BonusControl) DoPageFill(rule *repo.DataRule) lik.Lister {
	rows := it.TablePageFill(rule)
	what := "bonuses.created_at AS date," +
		"bonuses_list.title AS title," +
		"bonuses_list.value AS value"
	from := "bonuses INNER JOIN bonuses_list ON bonuses.bonuses_list_id=bonuses_list.id"
	where := fmt.Sprintf("bonuses.members_id=%d", int(it.IdMain))
	order := "bonuses.created_at"
	if list := jone.DB.GetListElm(what, from, where, order); list != nil {
		summa := 0
		for ne := 0; ne < list.Count(); ne++ {
			if elm := list.GetSet(ne); elm != nil {
				summa += elm.GetInt("value")
				row := lik.BuildSet("id", 1 + ne)
				if match := lik.RegExParse(elm.GetString("date"), "(\\d\\d\\d\\d).(\\d\\d).(\\d\\d)"); match != nil {
					row.SetItem(match[3]+"/"+match[2]+"/"+match[1], "date")
				}
				row.SetItem(elm.GetString("title"), "title")
				row.SetItem(elm.GetString("value"), "value")
				row.SetItem(summa, "summa")
				rows.AddItems(row)
			}
		}
	}
	if rows.Count() <= 0 {
		rows.AddItemSet("id=1", "name=НЕТ")
	}
	return rows
}

