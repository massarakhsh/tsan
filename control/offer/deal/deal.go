// Контроллер сделки.
//
// Состав контроллера:
//	BaseOn	controls.OneForeControl	//	Включает дексриптор стандартного окна из пяти частей
//	LeftUp	ManageControl			//	Дескриптор управления сделкой
//	RightUp	ChooseControl			//	Дексриптор выбора аявки для сделки
package deal

import (
	"bitbucket.org/961961/tsan/control/controls"
	"bitbucket.org/961961/tsan/repo"
	"bitbucket.org/shaman/lik"
)

//	Дескриптор контроллера сделки
type DealControl struct {
	controls.FrameControl //	Включает дексриптор стандартного окна из пяти частей
}

//	Конструктор дескриптора контроллера
func BuildOfferDeal(rule *repo.DataRule, id lik.IDB) *DealControl {
	it := &DealControl{}
	it.ControlInitialize("offerdeal", id)
	it.SetLayoutLR(0, 40)
	it.AddControl("LU", BuildManage(rule, it.Mode, id))
	it.AddControl("RU", controls.BuildChooseOffer(rule, it.Mode, id))
	return it
}

