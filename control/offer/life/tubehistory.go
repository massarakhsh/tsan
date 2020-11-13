package life

import (
	"github.com/massarakhsh/tsan/control/controls"
	"github.com/massarakhsh/tsan/repo"
	"github.com/massarakhsh/lik"
	"github.com/massarakhsh/lik/likbase"
)

//	Дескриптор окна истории
type TubeHistory struct {
	controls.TubeControl
}

//	Интерфейс заполнения строки
type dealHistoryRow struct {
	It	*TubeHistory
}
func (it *dealHistoryRow) Run(rule *repo.DataRule, data lik.Seter, row lik.Seter) {
	it.It.HistoryRowFill(rule, data, row)
}

//	Конструктор дескриптора
func BuildTubeHistory(rule *repo.DataRule, main string, id lik.IDB) *TubeHistory {
	it := &TubeHistory{}
	it.Self = it
	it.TubeInitialize(rule, main, "history", id,"История изменений")
	it.ItTubeRow = &dealHistoryRow{it}
	return it
}

//	Выполнение команд
func (it *TubeHistory) RunTubeExecute(rule *repo.DataRule, cmd string, data lik.Seter) {
	it.TubeExecute(rule, cmd, data)
}

//	Создание таблицы
func (it *TubeHistory) RunTubeGridFill(rule *repo.DataRule) {
	it.TubeGridFill(rule)
}

//	Завершение редактирования
func (it *TubeHistory) RunTubeFinalEdit(rule *repo.DataRule, elm *likbase.ItElm) {
}

//	Созадие строки
func (it *TubeHistory) HistoryRowFill(rule *repo.DataRule, data lik.Seter, row lik.Seter) {
	it.TubeRowFill(rule, data, row)
}

