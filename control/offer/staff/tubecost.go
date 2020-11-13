package staff

import (
	"bitbucket.org/961961/tsan/control/controls"
	"bitbucket.org/961961/tsan/jone"
	"bitbucket.org/961961/tsan/repo"
	"bitbucket.org/shaman/lik"
	"bitbucket.org/shaman/lik/likbase"
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

