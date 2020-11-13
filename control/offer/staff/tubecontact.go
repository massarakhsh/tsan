package staff

import (
	"github.com/massarakhsh/tsan/control/controls"
	"github.com/massarakhsh/tsan/jone"
	"github.com/massarakhsh/tsan/repo"
	"github.com/massarakhsh/tsan/show"
	"github.com/massarakhsh/lik"
	"github.com/massarakhsh/lik/likbase"
	"fmt"
)

//	Дескриптор списка контактов
type TubeContact struct {
	controls.TubeControl
}

//	Интерфейс заполнения формы
type dealContactFormFill struct {
	It	*TubeContact
}
func (it *dealContactFormFill) Run(rule *repo.DataRule) {
	it.It.ContactFormFill(rule)
}

//	Конструктор дескриптора списка контактов
func BuildContact(rule *repo.DataRule, main string, id lik.IDB) *TubeContact {
	it := &TubeContact{}
	it.Self = it
	it.TubeInitialize(rule, main, "contact", id, "Контакты")
	it.ItFormFill = &dealContactFormFill{it}
	return it
}

//	Выполнение команд
func (it *TubeContact) RunTubeExecute(rule *repo.DataRule, cmd string, data lik.Seter) {
	if true {
		it.TubeExecute(rule, cmd, data)
	}
}

//	Заполнение таблицы
func (it *TubeContact) RunTubeGridFill(rule *repo.DataRule) {
	it.TubeGridFill(rule)
	it.AddCommandImg(rule, 950, "Добавить", "toadd", "add")
}

//	Завершение редактирования
func (it *TubeContact) RunTubeFinalEdit(rule *repo.DataRule, elm *likbase.ItElm) {
}

//	Заполнение формы
func (it *TubeContact) ContactFormFill(rule *repo.DataRule) {
	it.TubeFormFill(rule)
	if _, data := it.FindListData(rule, lik.StrToInt(it.Sel)); data != nil {
		if idbell := data.GetIDB("bellid"); idbell != 0 {
			if bell := jone.TableBell.GetElm(idbell); bell != nil && !it.IsEdit() {
				segment := bell.GetString("segment")
				target := bell.GetString("target")
				fio := bell.GetString("client/namely") + " " +
					bell.GetString("client/paterly") + " " +
					bell.GetString("client/family")
				it.Form.Items.AddItemSet("type=string", "label=Клиент",
					"name=s__fio", "value", fio, "editable=false")
				it.Form.Items.AddItemSet("type=string", "label=Телефон",
					"name=s__phone", "value", bell.GetString("client/phone1"), "editable=false")
				it.Form.Items.AddItemSet("type=string", "label=Цель",
					"name=s__target", "key=target", "value", target, "editable=false")
				it.Form.Items.AddItemSet("type=string", "label=Сегмент",
					"name=s__segment", "key=segment", "value", bell.GetString("segment"), "editable=false")
				it.Form.Items.AddItemSet("type=string", "label=Тип недв.",
					"name=s__realty", "key=realty", "value", bell.GetString("realty"), "editable=false")
				text := "Открыть контакт"
				if segment == jone.DoRent && target == "sale" {
					text = "Контакт сдать в аренду"
				} else if segment == jone.DoRent && target == "buy" {
					text = "Контакт снять в аренду"
				} else if target == "sale" {
					text = "Контакт на продажу"
				} else if target == "buy" {
					text = "Контакт на покупку"
				}
				code := show.LinkTextProc("cmd", fmt.Sprintf("%s №%03d", text, int(idbell)), fmt.Sprintf("bind_bell(%d)", int(idbell)))
				it.Form.Items.AddItemSet("type=html", "value", code.ToString())
			}
		}
	}
}
