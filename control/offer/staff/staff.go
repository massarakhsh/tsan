//	Модуль служебного окна заявки
package staff

import (
	"bitbucket.org/961961/tsan/control/controls"
	"bitbucket.org/961961/tsan/control/offer/files"
	"bitbucket.org/961961/tsan/control/offer/life"
	"bitbucket.org/961961/tsan/jone"
	"bitbucket.org/961961/tsan/repo"
	"bitbucket.org/shaman/lik"
)

//	Дескриптор служебного окна
type StaffControl struct {
	controls.FrameControl
}

//	Конструктор дескриптора
func BuildOfferStaff(rule *repo.DataRule, id lik.IDB) *StaffControl {
	it := &StaffControl{}
	rule.CachePartIdPush("offer", id)
	it.ControlInitialize("offerstaff", id)
	it.SetLayoutFour(400, 50, 50)
	it.AddControl("L", BuildEditor(rule, it.Mode, id))
	it.AddControl("LU", BuildStatus(rule, it.Mode, id))
	it.AddControl("LD", BuildContact(rule, it.Mode, id))
	if target := jone.CalculatePartIdString("offer", id, "target"); target == "sale" {
		it.AddControl("RU", BuildMap(rule, it.Mode, id, true))
		it.AddControl("RD", BuildCost(rule, it.Mode, id))
	} else {
		it.AddControl("RU", files.BuildMedia(rule, it.Mode, "doc", id,"Документы"))
		it.AddControl("RD", life.BuildTubeHistory(rule, it.Mode, id))
	}
	return it
}

