package staff

import (
	"bitbucket.org/961961/tsan/control/controls"
	"bitbucket.org/961961/tsan/jone"
	"bitbucket.org/961961/tsan/repo"
	"bitbucket.org/shaman/lik"
	"bitbucket.org/shaman/lik/likbase"
	"fmt"
	"strings"
	"time"
)

//	Дескриптор списка статусов
type TubeStatus struct {
	controls.TubeControl
}

//	Конструктор дескриптора статусов
func BuildStatus(rule *repo.DataRule, main string, id lik.IDB) *TubeStatus {
	it := &TubeStatus{}
	it.Self = it
	it.TubeInitialize(rule,main, "status", id, "Статус заявки")
	return it
}

//	Выполнение команд
func (it *TubeStatus) RunTubeExecute(rule *repo.DataRule, cmd string, data lik.Seter) {
	if cmd == "status" {
		it.RunTubeExecute(rule, rule.Shift(), data)
	} else if cmd == "setstatus" {
		it.statusSetStatus(rule)
	} else {
		it.TubeExecute(rule, cmd, data)
	}
}

//	Заполнение таблицы
func (it *TubeStatus) RunTubeGridFill(rule *repo.DataRule) {
	it.TubeGridFill(rule)
	data := lik.BuildList()
	for _,svar := range []string{ "pass,Не активна", "active,Активна", "avance,Аванс", "done,Завершена" } {
		lvar := strings.Split(svar, ",")
		data.AddItems(lik.BuildSet("type=button", "text", lvar[1],
			"handler", fmt.Sprintf("function_set_status_offer(%s)", svar)))
	}
	it.AddCommandItem(rule, 950, lik.BuildSet(
		"type=button", "tip=Изменение статуса", "id=topofferstatus",
		"imageCls", "imgadd", "menu", data,
	))
}

//	Завершение редактирования
func (it *TubeStatus) RunTubeFinalEdit(rule *repo.DataRule, elm *likbase.ItElm) {
	if list := repo.GetHistory(elm,"status"); len(list) > 0 {
		hist := list[len(list)-1]
		jone.SetElmValue(elm, hist.GetString("status"), "status")
	}
}

//	Заполнение строки
func (it *TubeStatus) RunTubeRow(rule *repo.DataRule, data lik.Seter, row lik.Seter) {
}

//	Установка статуса
func (it *TubeStatus) statusSetStatus(rule *repo.DataRule) {
	status := rule.Shift()
	comment := lik.StringFromXS(rule.Shift())
	if elm := jone.GetElm(it.Part, it.IdMain); elm != nil {
		at := int(time.Now().Unix())
		repo.AddHistorySet(rule, elm, "date", at, "what=status", "status", status, "notify", comment)
		jone.SetElmValue(elm, status, "status")
		rule.OnChangeData()
	}
}

