//	Модуль кабинета клиента
package client

import (
	"github.com/massarakhsh/tsan/control/controls"
	"github.com/massarakhsh/tsan/repo"
	"github.com/massarakhsh/lik"
)

//	Дескриптор кабинета
type ClientCard struct {
	controls.FrameControl //	На окне формата 1+4
}

//	Конструктор дескриптора
func BuildClientCard(rule *repo.DataRule, id lik.IDB) *ClientCard {
	it := &ClientCard{}
	it.ControlInitialize("clientcard", id)
	it.SetLayoutFour(320, 50, 50)
	it.AddControl("L", BuildEditor(rule, it.Frame, id))
	it.AddControl("LU", controls.BuildEmpty(rule, "Свободно"))
	it.AddControl("RU", controls.BuildEmpty(rule, "Свободно"))
	it.AddControl("LD", controls.BuildEmpty(rule, "Свободно"))
	it.AddControl("RD", controls.BuildEmpty(rule, "Свободно"))
	return it
}

