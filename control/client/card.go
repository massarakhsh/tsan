//	Модуль кабинета клиента
package client

import (
	"bitbucket.org/961961/tsan/control/controls"
	"bitbucket.org/961961/tsan/repo"
	"bitbucket.org/shaman/lik"
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

