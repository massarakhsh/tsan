//	Модуль файлов заявки
package files

import (
	"bitbucket.org/961961/tsan/control/controls"
	"bitbucket.org/961961/tsan/jone"
	"bitbucket.org/961961/tsan/repo"
	"bitbucket.org/shaman/lik"
)

//	Дескриптор обозревателя файлов
type FilesControl struct {
	controls.FrameControl
}

//	Конструктор дескриптора
func BuildOfferFiles(rule *repo.DataRule, id lik.IDB) *FilesControl {
	it := &FilesControl{}
	it.ControlInitialize("offerfiles", id)
	it.SetLayoutLRR(0, 50, 50)
	it.AddControl("RU", BuildMedia(rule, it.Mode, "doc", id,"Документы"))
	if target := jone.CalculatePartIdString("offer", id, "target"); target == "sale" {
		it.AddControl("LU", BuildMedia(rule, it.Mode, "photo", id, "Изображения"))
		it.AddControl("RD", BuildMedia(rule, it.Mode, "link", id, "Внешние ссылки"))
	}
	return it
}

