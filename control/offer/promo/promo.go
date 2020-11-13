//	Модуль окна продвижения
package promo

import (
	"github.com/massarakhsh/tsan/control/controls"
	"github.com/massarakhsh/tsan/control/message"
	"github.com/massarakhsh/tsan/repo"
	"github.com/massarakhsh/lik"
)

//	Дескриптор продвижения
type PromoControl struct {
	controls.FrameControl
}

//	Конструктор дескриптора
func BuildOfferPromo(rule *repo.DataRule, id lik.IDB) *PromoControl {
	it := &PromoControl{}
	it.ControlInitialize( "offerpromo", id)
	it.SetLayoutFour(400, 50, 50)
	it.AddControl("L", BuildManage(rule, it.Mode, id))
	it.AddControl("LU", BuildPromouter(rule, it.Mode, id))
	it.AddControl("RU", message.BuildList(rule, it.Mode, id))
	it.AddControl("LD", BuildExport(rule, it.Mode, id))
	it.AddControl("RD", BuildTubePromo(rule, it.Mode, id))
	return it
}

