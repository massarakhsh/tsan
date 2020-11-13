//	Модуль кабинета сотрудника
package member

import (
	"github.com/massarakhsh/tsan/control/controls"
	"github.com/massarakhsh/tsan/repo"
	"github.com/massarakhsh/lik"
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

