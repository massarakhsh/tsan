package staff

import (
	"github.com/massarakhsh/tsan/control/controls"
	"github.com/massarakhsh/tsan/jone"
	"github.com/massarakhsh/tsan/repo"
	"github.com/massarakhsh/lik"
	"github.com/massarakhsh/lik/likbase"
)

//	Дескриптор списка цен
type TubeCost struct {
	controls.TubeControl
}

//	Конструктор дескриптора цены
func BuildCost(rule *repo.DataRule, main string, id lik.IDB) *TubeCost {
	it := &TubeCost{}
	it.Self = it
	it.TubeInitialize(rule, main, "cost", id,"Цена")
	return it
}

//	Выполнение команд
func (it *TubeCost) RunTubeExecute(rule *repo.DataRule, cmd string, data lik.Seter) {
	it.TubeExecute(rule, cmd, data)
}

//	Заполнение таблицы
func (it *TubeCost) RunTubeGridFill(rule *repo.DataRule) {
	it.TubeGridFill(rule)
	it.AddCommandImg(rule, 950, "Добавить", "toadd", "add")
}

//	Завершение редактирования
func (it *TubeCost) RunTubeFinalEdit(rule *repo.DataRule, elm *likbase.ItElm) {
	if list := repo.GetHistory(elm,"cost"); len(list) > 0 {
		hist := list[len(list)-1]
		jone.SetElmValue(elm, hist.GetInt("cost"), "cost")
	}
}

