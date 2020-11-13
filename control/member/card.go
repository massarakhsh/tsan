//	Модуль кабинета сотрудника
package member

import (
	"bitbucket.org/961961/tsan/control/controls"
	"bitbucket.org/961961/tsan/repo"
	"bitbucket.org/shaman/lik"
)

//	Дескриптор кабинета
type MemberCard struct {
	controls.FrameControl //	На окне формата 1+4
}

//	Конструктор дескриптора
func BuildMemberCard(rule *repo.DataRule, id lik.IDB) *MemberCard {
	it := &MemberCard{}
	it.ControlInitialize("membercard", id)
	it.SetLayoutFour(320, 50, 50)
	it.AddControl("L", BuildEditor(rule, it.Frame, id))
	it.AddControl("LU", BuildManage(rule, it.Frame, id))
	it.AddControl("RU", BuildBonus(rule, it.Frame, id))
	it.AddControl("LD", controls.BuildEmpty(rule, "Свободно"))
	it.AddControl("RD", controls.BuildEmpty(rule, "Свободно"))
	return it
}

