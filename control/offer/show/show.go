//	Модуль отображение заявки
package show

import (
	"bitbucket.org/961961/tsan/control/controls"
	"bitbucket.org/961961/tsan/control/offer/staff"
	"bitbucket.org/961961/tsan/jone"
	"bitbucket.org/961961/tsan/repo"
	"bitbucket.org/shaman/lik"
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

