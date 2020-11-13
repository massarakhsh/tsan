package deal

import (
	"github.com/massarakhsh/tsan/control"
	"github.com/massarakhsh/tsan/jone"
	"github.com/massarakhsh/tsan/repo"
	"github.com/massarakhsh/tsan/show"
	"github.com/massarakhsh/lik"
	"github.com/massarakhsh/lik/likbase"
	"github.com/massarakhsh/lik/likdom"
	"fmt"
	"time"
)

//	Дескриптор управления сделкой
type ManageControl struct {
	control.DataControl		//	Включенный декскриптор общего контроллера
	IdPair	lik.IDB			//	Идентификатор привязанной сделки
}

//	Обработчик команд управления
type dealManageExecute struct {
	It	*ManageControl
}
func (it *dealManageExecute) Run(rule *repo.DataRule, cmd string, data lik.Seter) {
	it.It.ManageExecute(rule, cmd, data)
}

//	Конструктор дескриптора управления сделкой
func BuildManage(rule *repo.DataRule, main string, id lik.IDB) *ManageControl {
	it := &ManageControl{}
	deal := repo.SeekOfferDeal(id)
	it.IdPair = deal.IdPair
	it.ControlInitializeZone(main, id, "manage")
	it.ItExecute = &dealManageExecute{it}
	return it
}

//	Отображение
func (it *ManageControl) BuildShow(rule *repo.DataRule, sx int, sy int) likdom.Domer {
	tbl := likdom.BuildTableClass("fill", "id=offerdeal_all")
	if row := tbl.BuildTr(); row != nil {
		row.BuildTdClass("pair", control.MiniMax(sx/2, 0)...).AppendItem(it.BuildShowDeal(rule, sx/2, 0))
		row.BuildTdClass("pair", control.MiniMax(sx/2, 0)...)
	}
	if row := tbl.BuildTr(); row != nil {
		row.BuildTdClass("pair", control.MiniMax(sx/2, 0)...).AppendItem(it.BuildShowOffer(rule, sx/2, 0, true))
		row.BuildTdClass("pair", control.MiniMax(sx/2, 0)...).AppendItem(it.BuildShowOffer(rule, sx/2, 0, false))
	}
	if row := tbl.BuildTr(); row != nil {
		row.BuildTdClass("pair").BuildString("Требования")
		row.BuildTdClass("pair", control.MiniMax(sx/2, 0)...).AppendItem(it.BuildShowRequire(rule, sx/2, 0))
	}
	tbl.BuildTrTdClass("fill", "colspan=2")
	return tbl
}

//	Отображение сделки
func (it *ManageControl) BuildShowDeal(rule *repo.DataRule, sx int, sy int) likdom.Domer {
	deal := repo.SeekOfferDeal(it.IdMain)
	tbl := likdom.BuildTableClass("manage")
	if row := tbl.BuildTr(); row != nil {
		row.BuildTd().BuildString("Сделка")
		td := row.BuildTd()
		if delm := jone.TableDeal.GetElm(deal.IdDeal); delm != nil {
			td.BuildString(fmt.Sprintf("Сделка №%03d<br>", int(deal.IdDeal)))
			if rule.IAmAdmin() {
				td.AppendItem(show.LinkTextProc("txtcmd", "Удалить сделку", "deal_manage_delete()"))
				td.BuildString("<br>")
			}
		} else {
			td.BuildString("Не открыта<br>")
			td.AppendItem(show.LinkTextCmd("txtcmd", "Открыть проект сделки", "offerdeal", "manage", "create"))
			td.BuildString("<br>")
		}
	}
	return tbl
}

//	Отображение заявки
func (it *ManageControl) BuildShowOffer(rule *repo.DataRule, sx int, sy int, ismain bool) likdom.Domer {
	deal := repo.SeekOfferDeal(it.IdMain)
	var id lik.IDB
	var target string
	if ismain {
		id = it.IdMain
		target = deal.TargetMain
	} else {
		id = it.IdPair
		target = deal.TargetPair
	}
	tbl := likdom.BuildTableClass("manage", control.MiniMax(sx, sy)...)

	if row := tbl.BuildTr(); row != nil {
		td := row.BuildTdClass("tophead", "colspan=2")
		td.BuildString(jone.SystemStringTranslate("target", target))
	}
	elm := jone.TableOffer.GetElm(id)
	if elm != nil {
		tbl.BuildTrTd("colspan=2").BuildString("<hr>")
		if row := tbl.BuildTr(); row != nil {
			row.BuildTd().BuildString("Заявка")
			td := row.BuildTd()
			text := fmt.Sprintf("№%03d", int(elm.Id))
			if !ismain {
				path := fmt.Sprintf("/offershow%d?_tp=1", int(elm.Id))
				text = show.LinkTextProc("cmd", text, fmt.Sprintf("lik_window_part('%s')", path)).ToString()
			}
			td.BuildString(text + "<br>")
			td.BuildString(jone.CalculateElmTranslate(elm, "target") + "<br>")
			td.BuildString(jone.CalculateElmTranslate(elm, "segment") + "<br>")
			td.BuildString(jone.CalculateElmTranslate(elm, "status") + "<br>")
		}
		//if row := tbl.BuildTr(); row != nil {
		//	row.BuildTd().BuildString("Риэлтор")
		//	row.BuildTd().BuildString(jone.CalculatePartIdText("member", elm.GetIDB("memberid")))
		//}
	}

	if delm := jone.TableDeal.GetElm(deal.IdDeal); delm == nil {
	} else if ismain {
		tbl.BuildTrTd("colspan=2").BuildString("<hr>")
		if row := tbl.BuildTr(); row != nil {
			row.BuildTd().BuildString("Сделка")
			td := row.BuildTd()
			if status := delm.GetString("status"); status == "done" {
				td.BuildString("Закрыта<br>")
			} else if status == "active" {
				td.BuildString("Открыта<br>")
			} else if status == "avance" {
				td.BuildString("Получен аванс<br>")
			} else {
				td.BuildString("Подготавливается<br>")
			}
			if status := delm.GetString("status"); status != "done" && deal.IdPair != 0 {
				td.BuildString("Изменить статус:<br>")
				if status == "avance" {
					td.AppendItem(show.LinkTextCmd("txtcmd", "Возвращен аванс", it.Mode, "manage", "dealstatus/active"))
					td.BuildString("<br>")
				} else if status != "avance" {
					td.AppendItem(show.LinkTextCmd("txtcmd", "Получен аванс", it.Mode, "manage", "dealstatus/avance"))
					td.BuildString("<br>")
				}
				td.AppendItem(show.LinkTextCmd("txtcmd", "Закрыть сделку", it.Mode, "manage", "dealstatus/done"))
				td.BuildString("<br>")
			}
		}
	} else if it.IdPair != deal.IdPair {
		tbl.BuildTrTd("colspan=2").BuildString("<hr>")
		if row := tbl.BuildTr(); row != nil {
			row.BuildTdClass("attention").BuildString("Имеются изменения")
			td := row.BuildTd()
			if text := "В сделке: "; text != "" {
				if deal.IdPair == 0 {
					text += "заявка не указана"
				} else {
					text += fmt.Sprintf("№%03d", int(deal.IdPair))
				}
				td.BuildString(text + "<br>")
			}
			if text := "Выбрано: "; text != "" {
				if it.IdPair == 0 {
					text += "заявка не указана"
				} else {
					text += fmt.Sprintf("№%03d", int(it.IdPair))
				}
				td.BuildString(text + "<br>")
			}
			td.BuildString("<br>")
			td.AppendItem(show.LinkTextCmd("txtcmd", "Сохранить изменения", "offerdeal", "manage", "repoffer"))
			td.BuildString("<br>")
			td.AppendItem(show.LinkTextCmd("txtcmd", "Вернуть как было", "offerdeal", "manage", "oldoffer"))
			td.BuildString("<br>")
		}
	} else if deal.IdPair != 0 {
		tbl.BuildTrTd("colspan=2").BuildString("<hr>")
		if row := tbl.BuildTr(); row != nil {
			row.BuildTd().BuildString("Заявка")
			td := row.BuildTd()
			td.AppendItem(show.LinkTextCmd("txtcmd", "Стереть", "offerdeal", "manage", "deloffer"))
		}
	} else {
		tbl.BuildTrTd("colspan=2").BuildString("<hr>")
		if row := tbl.BuildTr(); row != nil {
			row.BuildTd().BuildString("Нет заявки")
			td := row.BuildTd()
			td.AppendItem(show.LinkTextCmd("txtcmd", "Создать для сделки", "offerdeal", "manage", "newoffer"))
		}
	}
	return tbl
}

//	Отображение требований запроса
func (it *ManageControl) BuildShowRequire(rule *repo.DataRule, sx int, sy int) likdom.Domer {
	tbl := likdom.BuildTableClass("manage", control.MiniMax(sx, sy)...)
	if elm := jone.TableOffer.GetElm(it.IdMain); elm != nil {
		if ent := repo.GenStruct.FindEnt("require"); ent != nil {
			if content := ent.GetContent(); content != nil {
				for _,pos := range content {
					part := pos.GetString("part")
					if val := jone.CalculateElmString(elm, "require/" + part); val != "" {
						if row := tbl.BuildTr(); row != nil {
							td := row.BuildTd()
							td.BuildString(part + ": " + val + "<br>")
						}
					}
				}
			}
		}
	}
	return tbl
}

//	Выполнение команд
func (it *ManageControl) ManageExecute(rule *repo.DataRule, cmd string, data lik.Seter) {
	if cmd == "manage" {
		it.ManageExecute(rule, rule.Shift(), data)
	} else if cmd == "create" {
		it.cmdCreateDeal(rule)
	} else if cmd == "delete" {
		it.cmdDeleteDeal(rule)
	} else if cmd == "choose" {
		it.IdPair = lik.StrToIDB(rule.Shift())
		rule.OnChangeData()
	} else if cmd == "repoffer" {
		if deal := repo.SeekOfferDeal(it.IdMain); deal.IdDeal != 0 {
			if delm := jone.TableDeal.GetElm(deal.IdDeal); delm != nil {
				if deal.IsFrom {
					delm.SetValue(it.IdPair, "buyid")
				} else {
					delm.SetValue(it.IdPair, "saleid")
				}
				rule.OnChangeData()
			}
		}
	} else if cmd == "oldoffer" {
		if deal := repo.SeekOfferDeal(it.IdMain); deal.IdDeal != 0 {
			it.IdPair = deal.IdPair
			rule.OnChangeData()
		}
	} else if cmd == "deloffer" {
		it.IdPair = 0
		rule.OnChangeData()
	} else if cmd == "newoffer" {
		it.cmdCreateOffer(rule)
	} else if cmd == "dealstatus" {
		it.cmdDealStatus(rule)
	}
}

//	Создание сделки
func (it *ManageControl) cmdCreateDeal(rule *repo.DataRule) {
	deal := repo.SeekOfferDeal(it.IdMain)
	if elm := jone.TableOffer.GetElm(it.GetId()); elm != nil && deal.IdDeal == 0 {
		delm := jone.TableDeal.CreateElm()
		jone.SetElmValue(delm, int(time.Now().Unix()), "date_open")
		jone.SetElmValue(delm, deal.IdFrom, "saleid")
		jone.SetElmValue(delm, deal.IdTo, "buyid")
		jone.SetElmValue(delm, "active", "status")
		notify := "Создана сделка " + jone.CalculateElmSid(delm)
		repo.AddHistorySet(rule, elm,"date", int(time.Now().Unix()), "what=deal", "dealid", delm.Id, "notify", notify)
	}
	rule.OnChangeData()
}

//	Удаление сделки
func (it *ManageControl) cmdDeleteDeal(rule *repo.DataRule) {
	deal := repo.SeekOfferDeal(it.IdMain)
	if delm := jone.TableDeal.GetElm(deal.IdDeal); delm != nil {
		delm.Delete()
	}
	rule.OnChangeData()
}

//	Создание сделки
func (it *ManageControl) cmdCreateOffer(rule *repo.DataRule) {
	deal := repo.SeekOfferDeal(it.IdMain)
	if elm := jone.TableOffer.GetElm(deal.IdMain); elm != nil {
		if delm := jone.TableDeal.GetElm(deal.IdDeal); delm != nil {
			at := int(time.Now().Unix())
			pair := jone.TableOffer.CreateElm()
			jone.SetElmValue(pair, deal.TargetPair, "target")
			jone.SetElmValue(pair, jone.CalculateElmString(elm, "segment"), "segment")
			jone.SetElmValue(pair, jone.CalculateElmString(elm, "objectid/realty"), "objectid/realty")
			jone.SetElmValue(pair, rule.ItSession.IdMember, "memberid")
			jone.SetElmValue(pair, "active", "status")
			notify := "Создана для заявки " + jone.CalculateElmSid(elm)
			repo.AddHistorySet(rule, pair,"date", at, "what=create", "offerid", elm.Id, "notify", notify)
			pair.OnModifyWait()
			it.IdPair = pair.Id
			if deal.IsFrom {
				delm.SetValue(it.IdPair, "buyid")
			} else {
				delm.SetValue(it.IdPair, "saleid")
			}
			notify = "Создана " + jone.CalculateElmTranslate(pair, "target") + " №" + jone.CalculateElmSid(pair)
			repo.AddHistorySet(rule, elm,"date", at, "what=deal", "offerid", pair.Id, "notify", notify)
			delm.OnModifyWait()
		}
	}
	rule.OnChangeData()
}

//	Изменение статуса сделки
func (it *ManageControl) cmdDealStatus(rule *repo.DataRule) {
	status := rule.Shift()
	deal := repo.SeekOfferDeal(it.IdMain)
	at := int(time.Now().Unix())
	if delm := jone.TableDeal.GetElm(deal.IdDeal); delm != nil {
		jone.SetElmValue(delm, status, "status")
		for side := 0; side < 2; side++ {
			var elm *likbase.ItElm
			if side == 0 {
				elm = jone.TableOffer.GetElm(deal.IdFrom)
			} else {
				elm = jone.TableOffer.GetElm(deal.IdTo)
			}
			if elm != nil {
				jone.SetElmValue(elm, status, "status")
				notify := "Изменён статус: " + jone.SystemStringTranslate("status", status)
				repo.AddHistorySet(rule, elm, "date", at, "what=status", "status", status, "notify", notify)
			}
		}
	}
	rule.OnChangeData()
}

