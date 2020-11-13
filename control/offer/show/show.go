//	Модуль отображение заявки
package show

import (
	"github.com/massarakhsh/tsan/control/controls"
	"github.com/massarakhsh/tsan/control/offer/staff"
	"github.com/massarakhsh/tsan/jone"
	"github.com/massarakhsh/tsan/repo"
	"github.com/massarakhsh/lik"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

//	Дескриптор заявки
type ShowControl struct {
	controls.FrameControl
	TabId   int
}

//	Конструктор дескриптора
func BuildOfferShow(rule *repo.DataRule, id lik.IDB) *ShowControl {
	it := &ShowControl{}
	rule.CachePartIdPush("offer", id)
	it.ControlInitialize("offershow", id)
	it.AddControl("L", BuildCharacter(rule, it.Mode, it.IdMain))
	if target := jone.CalculatePartIdString("offer", it.IdMain, "target"); target == "sale" {
		it.SetLayoutLRR(340, 50, 66)
		it.AddControl("LU", BuildGallery(rule, it.Mode, it.IdMain))
		it.AddControl("RU", staff.BuildMap(rule, it.Mode, it.IdMain, false))
		it.AddControl("RD", BuildDescriptor(rule, it.Mode, id))
	} else {
		it.SetLayoutOne(340)
		it.AddControl("LU", controls.BuildChooseOffer(rule, it.Mode, it.IdMain))
	}
	it.TabId = 0
	return it
}

