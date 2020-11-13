package promo

import (
	"bitbucket.org/961961/tsan/control"
	"bitbucket.org/961961/tsan/jone"
	"bitbucket.org/961961/tsan/repo"
	"bitbucket.org/961961/tsan/show"
	"bitbucket.org/shaman/lik"
	"bitbucket.org/shaman/lik/likdom"
	"time"
)

//	Дескриптор управления продвижением
type ManageControl struct {
	control.DataControl
}

//	Интерфейс команд
type dealManageExecute struct {
	It	*ManageControl
}
func (it *dealManageExecute) Run(rule *repo.DataRule, cmd string, data lik.Seter) {
	it.It.ManageExecute(rule, cmd, data)
}

//	Создание дескриптора
func BuildManage(rule *repo.DataRule, main string, id lik.IDB) *ManageControl {
	it := &ManageControl{}
	it.ControlInitializeZone(main, id, "manage")
	it.ItExecute = &dealManageExecute{it}
	return it
}

//	Отображение окна
func (it *ManageControl) BuildShow(rule *repo.DataRule, sx int, sy int) likdom.Domer {
	elm := jone.TableOffer.GetElm(it.GetId())
	tbl := likdom.BuildTableClass("manage")
	if row := tbl.BuildTr(); row != nil {
		row.BuildTd().BuildString("")
		td := row.BuildTd()
		td.BuildString("Публикация рекламы<br>")
		if status := jone.CalculateElmString(elm, "status"); status != jone.ItActive {
			stex := "Нет статуса активности"
			if status != "" {
				stex = jone.SystemStringTranslate("status", status)
			}
			td.BuildItemClass("b", "atten").BuildString(stex)
		}
	}
	tbl.BuildTrTd("colspan=2").BuildString("<hr>")
	if row := tbl.BuildTr(); row != nil {
		row.BuildTd().BuildString("Готовность<br>(риэлтор)")
		td := row.BuildTd()
		ok := jone.CalculateElmBool(elm, "export/ready")
		if ok {
			td.BuildItemClass("b", "good").BuildString("Готова к выгрузке")
		} else {
			td.BuildItemClass("b", "atten").BuildString("Не готова к выгрузке")
		}
		if memrd := jone.CalculatePartIdText("member", jone.CalculateElmIDB(elm, "export/readyid")); memrd != "" {
			td.BuildString("<br>" + memrd)
		}
		if rule.RightElmOffer(elm, "edit") {
			td.BuildString("<br>")
			if ok {
				td.AppendItem(show.LinkTextCmd("txtcmd", "Снять готовность", "promo", "manage", "noready"))
			} else {
				td.AppendItem(show.LinkTextCmd("txtcmd", "Установить готовность", "promo", "manage", "ready"))
			}
		}
	}
	tbl.BuildTrTd("colspan=2").BuildString("<hr>")
	if row := tbl.BuildTr(); row != nil {
		row.BuildTd().BuildString("Подтверждение<br>(менеджер)")
		td := row.BuildTd()
		ok := jone.CalculateElmBool(elm, "export/confirm")
		if ok {
			td.BuildItemClass("b", "good").BuildString("Готовность подтверждена")
		} else {
			td.BuildItemClass("b", "atten").BuildString("Готовность не подтверждена")
		}
		if memrd := jone.CalculatePartIdText("member", jone.CalculateElmIDB(elm, "export/confirmid")); memrd != "" {
			td.BuildString("<br>" + memrd)
		}
		if rule.RightElmOffer(elm, "promo") {
			td.BuildString("<br>")
			if ok {
				td.AppendItem(show.LinkTextCmd("txtcmd", "Отозвать подтверждение", "promo", "manage", "noconfirm"))
			} else {
				td.AppendItem(show.LinkTextCmd("txtcmd", "Подтвердить", "promo", "manage", "confirm"))
			}
		}
	}
	tbl.BuildTrTd("colspan=2").BuildString("<hr>")
	if row := tbl.BuildTr(); row != nil {
		row.BuildTd().BuildString("Публикация<br>(менеджер<br>по рекламе)")
		td := row.BuildTd()
		ok := jone.CalculateElmBool(elm, "export/enable")
		if ok {
			td.BuildItemClass("b", "good").BuildString("Разрешена к публикации")
		} else {
			td.BuildItemClass("b", "atten").BuildString("Публикация не разрешена")
		}
		if memrd := jone.CalculatePartIdText("member", jone.CalculateElmIDB(elm, "export/enableid")); memrd != "" {
			td.BuildString("<br>" + memrd)
		}
		if rule.ICanRole(jone.ItAdvert) {
			td.BuildString("<br>")
			if ok {
				td.AppendItem(show.LinkTextCmd("txtcmd", "Снять с публикации", "promo", "manage", "noenable"))
			} else {
				td.AppendItem(show.LinkTextCmd("txtcmd", "Опубликовать", "promo", "manage", "enable"))
			}
		}
	}
	tbl.BuildTrTd("colspan=2").BuildString("<hr>")
	/*if row := tbl.BuildTr(); row != nil {
		row.BuildTd().BuildString("Проверка")
		td := row.BuildTd()
		td.AppendItem(show.LinkTextCmd("txtcmd", "Тест", "promo", "manage", "test"))
	}*/
	return tbl
}

//	Команды управления продвижением
func (it *ManageControl) ManageExecute(rule *repo.DataRule, cmd string, data lik.Seter) {
	elm := jone.TableOffer.GetElm(it.GetId())
	at := int(time.Now().Unix())
	if cmd == "manage" {
		it.ManageExecute(rule, rule.Shift(), data)
	} else if match := lik.RegExParse(cmd, "(no|)(ready|confirm|enable)"); match != nil {
		ok := (match[1] == "")
		part := match[2]
		jone.SetElmValue(elm, ok, "export/" + part)
		jone.SetElmValue(elm, rule.ItSession.IdMember,"export/" + part + "id")
		repo.AddHistorySet(rule, elm, "date", at, "what=promo", "promo", part, "action", lik.IfString(ok, "yy", "nn"))
		rule.OnChangeData()
	} else if cmd == "test" {
		//taskexport.CalculeTestElm(it.IdMain)
		rule.OnChangeData()
	} else {
		it.ControlExecute(rule, cmd, data)
	}
}

