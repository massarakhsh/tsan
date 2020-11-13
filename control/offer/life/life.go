//	Модуль истории заявки
package life

import (
	"bitbucket.org/961961/tsan/control/controls"
	"bitbucket.org/961961/tsan/repo"
	"bitbucket.org/shaman/lik"
)

//	Дескриптор окна истории
type ShowLife struct {
	controls.FrameControl
}

//	Конструктор дескриптора
func BuildOfferLife(rule *repo.DataRule, id lik.IDB) *ShowLife {
	it := &ShowLife{}
	it.ControlInitialize("offerlife", id)
	it.SetLayoutLR(0, 50)
	it.AddControl("LU", BuildTubeHistory(rule, it.Mode, it.IdMain))
	it.AddControl("RU", controls.BuildEmpty(rule, "Подробности"))
	return it
}

