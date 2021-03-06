package show

import (
	"github.com/massarakhsh/tsan/control"
	"github.com/massarakhsh/tsan/jone"
	"github.com/massarakhsh/tsan/repo"
	"github.com/massarakhsh/lik"
	"github.com/massarakhsh/lik/likdom"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"strings"
)

//	Дескриптор описания заявки
type DescriptorControl struct {
	control.DataControl
}

//	Конструктор дескриптора описания
func BuildDescriptor(rule *repo.DataRule, main string, id lik.IDB) *DescriptorControl {
	it := &DescriptorControl{}
	it.ControlInitializeZone(main, id, "descriptor")
	return it
}

//	Отображение описания заявки
func (it *DescriptorControl) BuildShow(rule *repo.DataRule, sx int, sy int) likdom.Domer {
	elm := jone.TableOffer.GetElm(it.GetId())
	div := likdom.BuildDivClassId("definition", "definition")
	if def := jone.CalculateElmString(elm,"objectid/definition"); def == "" {
		div.BuildString("Текстовое описание отсутствует")
	} else {
		text := strings.ReplaceAll(def,"\n", "<br>")
		div.BuildString(text)
	}
	return div
}

