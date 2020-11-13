package promo

import (
	"bitbucket.org/961961/tsan/control/controls"
	"bitbucket.org/961961/tsan/repo"
	"bitbucket.org/shaman/lik"
	"bitbucket.org/shaman/lik/likbase"
)

//	Дескриптор истории продвижений
type TubePromo struct {
	controls.TubeControl
}

//	Конструктор дескриптора
func BuildTubePromo(rule *repo.DataRule, main string, id lik.IDB) *TubePromo {
	it := &TubePromo{}
	it.Self = it
	it.TubeInitialize(rule, main,"promo", id,"История продвижений")
	return it
}

//	Выполнение команд
func (it *TubePromo) RunTubeExecute(rule *repo.DataRule, cmd string, data lik.Seter) {
	it.TubeExecute(rule, cmd, data)
}

//	Построение таблицы
func (it *TubePromo) RunTubeGridFill(rule *repo.DataRule) {
	it.TubeGridFill(rule)
}

//	Завершение редактирования
func (it *TubePromo) RunTubeFinalEdit(rule *repo.DataRule, elm *likbase.ItElm) {
}

//	Заполнение строки
func (it *TubePromo) RunTubeRowFill(rule *repo.DataRule, data lik.Seter, row lik.Seter) {
}

