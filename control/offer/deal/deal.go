// Контроллер сделки.
//
// Состав контроллера:
//	BaseOn	controls.OneForeControl	//	Включает дексриптор стандартного окна из пяти частей
//	LeftUp	ManageControl			//	Дескриптор управления сделкой
//	RightUp	ChooseControl			//	Дексриптор выбора аявки для сделки
package deal

import (
	"github.com/massarakhsh/tsan/control/controls"
	"github.com/massarakhsh/tsan/repo"
	"github.com/massarakhsh/lik"
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

