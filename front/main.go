package front

import (
	"bitbucket.org/961961/tsan/control"
	"bitbucket.org/961961/tsan/control/bell"
	"bitbucket.org/961961/tsan/control/client"
	"bitbucket.org/961961/tsan/control/controls"
	"bitbucket.org/961961/tsan/control/deal"
	"bitbucket.org/961961/tsan/control/depart"
	"bitbucket.org/961961/tsan/control/member"
	"bitbucket.org/961961/tsan/control/offer"
	"bitbucket.org/961961/tsan/control/offer/deal"
	"bitbucket.org/961961/tsan/control/offer/files"
	"bitbucket.org/961961/tsan/control/offer/life"
	"bitbucket.org/961961/tsan/control/offer/promo"
	"bitbucket.org/961961/tsan/control/offer/show"
	"bitbucket.org/961961/tsan/control/offer/staff"
	"bitbucket.org/961961/tsan/control/tune"
	"bitbucket.org/961961/tsan/jone"
	"bitbucket.org/961961/tsan/repo"
	"bitbucket.org/shaman/lik"
	"bitbucket.org/shaman/lik/likbase"
	"bitbucket.org/shaman/lik/likdom"
	"fmt"
	"strings"
)

//	Добавить основную страницу
func mainPage(rule *repo.DataRule, pater likdom.Domer) {
	pater.AppendItem(mainPageGen(rule))
}

//	Построить основную страницу
func mainPageGen(rule *repo.DataRule) likdom.Domer {
	rule.OffChangePage()
	rule.OffChangeData()
	div := likdom.BuildDivClassId("main_page","main_page")
	tbl := div.BuildTableClass("fill")
	if len(rule.ItPage.Locates) == 0 {
		mainSetPath(rule, "/logon")
	} else if segment := rule.SeekSegment(""); segment == jone.DoTune && rule.ItPage.Locates[0].GetMode() != jone.DoTune {
		mainSetPath(rule, "/tune")
	} else if segment != jone.DoTune && rule.ItPage.Locates[0].GetMode() == jone.DoTune {
		if segment == jone.DoCall {
			mainSetPath(rule, "/bell")
		} else {
			mainSetPath(rule, "/sale")
		}
	}
	for nm := 0; nm < len(rule.ItPage.Locates); nm++ {
		td := tbl.BuildTrTdClass("main_menu")
		td.AppendItem(MenuBuildMain(rule, nm))
	}
	mainData(rule, tbl.BuildTrTdClass("main_data"))
	return div
}

//	Передать обновления на клиента
func mainStore(rule *repo.DataRule) {
	if rule.IsItChangePage() {
		rule.StoreItem(mainPageGen(rule))
	} else if rule.IsItChangeData() {
		rule.StoreItem(mainGenData(rule))
	}
}

//	Добавить данные на страницу
func mainData(rule *repo.DataRule, pater likdom.Domer) {
	rule.ItSession.WaitLogin = !rule.IsLogin()
	pater.AppendItem(mainGenData(rule))
}

//	Построить поля данных
func mainGenData(rule *repo.DataRule) likdom.Domer {
	rule.OffChangeData()
	div := likdom.BuildDivClassId("main_view", "main_view")
	if _,ctrl := GetFront(rule); ctrl != nil {
		rule.RegControl(ctrl)
		sx,sy := rule.ItPage.GetSize()
		div.AppendItem(ctrl.BuildShow(rule, sx -control.BD, sy - 30 - 30 * rule.GetLevel()))
	}
	return div
}

//	Установить путь restFull
func mainSetPath(rule *repo.DataRule, path string) {
	if rule.IsTechno {
		path = "/cabinet"
	}
	rule.ItPage.PathLast = path
	rule.ItPage.PathNeed = path
	rule.OnChangeData()
	var level int
	for level = 0; path != ""; level++ {
		match := lik.RegExParse(path,"^/([^/\\d]*)(\\d*)(.*)")
		if match == nil { break }
		main := match[1]
		id := likbase.StrToIDB(match[2])
		if !mainSetLevel(rule, level, main, id) {
			break
		}
		path = match[3]
	}
	if match := lik.RegExParse(path,"^/_(.*)"); match != nil {
		rule.PathCommand = match[1]
	} else if match := lik.RegExParse(path,"^/(.*)"); match != nil {
		rule.PathCommand = match[1]
	} else {
		rule.PathCommand = ""
	}
	if level == 0 {
		mainSetLevel(rule, level, "bell", 0)
		level++
		rule.OnChangePage()
	}
	if len(rule.ItPage.Locates) > level {
		rule.ItPage.Locates = rule.ItPage.Locates[:level]
		rule.OnChangePage()
	}
}

//	Установить уровень
func mainSetLevel(rule *repo.DataRule, level int, mode string, id lik.IDB) bool {
	if level < len(rule.ItPage.Locates) && rule.ItPage.Locates[level].GetMode() != mode {
		rule.ItPage.Locates = rule.ItPage.Locates[:level]
		rule.OnChangePage()
	}
	var ctrl control.Controller
	if level < len(rule.ItPage.Locates) {
		ctrl = rule.ItPage.Locates[level].(control.Controller)
	} else if lev,loc := FindControl(rule, mode); lev < 0 && loc != nil {
		rule.ItPage.Locates = append(rule.ItPage.Locates, loc)
		ctrl = loc
	} else if loc := mainBuildControl(rule, mode, id); loc != nil {
		rule.ItPage.Locates = append(rule.ItPage.Locates, loc)
		ctrl = loc
	}
	if ctrl != nil {
		ctrl.SetId(id)
		ctrl.RunUpdate(rule)
		//rule.ItPage.PathLast = rule.BuildFullPath("")
		rule.RegControl(ctrl)
	}
	return ctrl != nil
}

//	Построить контроллер на уровне
func mainBuildControl(rule *repo.DataRule, mode string, id lik.IDB) control.Controller {
	var ctrl control.Controller
	rule.OnChangePage()
	rule.OnChangeData()
	if !rule.IsLogin() {
		ctrl = controls.BuildLogon(rule, id)
	} else if mode == "logon" {
		ctrl = controls.BuildLogon(rule, id)
	} else if mode == "tune" {
		ctrl = tune.BuildTune(rule)
	} else if mode == "bell" {
		ctrl = bell.BuildListBell(rule, id)
	} else if match := lik.RegExParse(mode, "^(offer|sale|buy|lease|rent)(.*)"); match != nil {
		if match[2] == "show" {
			ctrl = show.BuildOfferShow(rule, id)
		} else if match[2] == "files" {
			ctrl = files.BuildOfferFiles(rule, id)
		} else if match[2] == "staff" {
			ctrl = staff.BuildOfferStaff(rule, id)
		} else if match[2] == "promo" {
			ctrl = promo.BuildOfferPromo(rule, id)
		} else if match[2] == "life" {
			ctrl = life.BuildOfferLife(rule, id)
		} else if match[2] == "deal" {
			ctrl = deal.BuildOfferDeal(rule, id)
		} else {
			ctrl = offer.BuildListOffer(rule, mode, id)
		}
	} else if match := lik.RegExParse(mode, "^client(.*)"); match != nil {
		if match[1] == "card" {
			ctrl = client.BuildClientCard(rule, id)
		} else {
			ctrl = client.BuildClientList(rule, id)
		}
	} else if match := lik.RegExParse(mode, "^member(.*)"); match != nil {
		if match[1] == "card" {
			ctrl = member.BuildMemberCard(rule, id)
		} else {
			ctrl = member.BuildMemberList(rule, id)
		}
	} else if mode == "depart" {
		ctrl = depart.BuildListDepart(rule, id)
	} else if mode == "deal" {
		ctrl = listdeal.BuildListDeal(rule, id)
	} else if mode != "" && !strings.HasPrefix(mode, "_") {
		fmt.Println("Unknown part '", mode, "'")
	}
	if ctrl != nil {
		rule.RegControl(ctrl)
	}
	return ctrl
}

