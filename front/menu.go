package front

import (
	"github.com/massarakhsh/tsan/jone"
	"github.com/massarakhsh/tsan/repo"
	"github.com/massarakhsh/lik"
	"github.com/massarakhsh/lik/likdom"
	"fmt"
	"strings"
	"time"
)

//	Построить меню уровня
func MenuBuildMain(rule *repo.DataRule, nm int) likdom.Domer {
	tbl := likdom.BuildItem("table","class=main_menu")
	row := tbl.BuildTr()
	menuFillMain(rule, row, nm)
	return tbl
}

//	Заполнить поля меню
func menuFillMain(rule *repo.DataRule, row likdom.Domer, nm int) {
	segment := rule.SeekSegment("")
	operator := rule.GetMember()
	mode := rule.ItPage.Locates[nm].GetMode()
	if nm == 0 {
		if td := row.BuildTdClass("menu menu_logo"); td != nil {
			a := td.BuildItem("a", "class=menutop", "href='#'", "onclick", "lik_go_part('/')")
			a.BuildUnpairItem("img", "src='/images/logo-tsan.png'", "title='На главную страницу'")
		}
		if td := row.BuildTdClass("menu menu_logo"); td != nil {
			if operator == nil {
				td.BuildString("РИЭЛТОР")
			} else {
				a := td.BuildItem("a", "class=menutop", "href='#'")
				a.SetAttr("onclick", "click_segment()")
				if segment == jone.DoTune {
					a.BuildString("НАСТРОЙКИ")
				} else if segment == jone.DoCall {
					a.BuildString("КОЛЛЦЕНТР")
				} else if segment == jone.DoSecond {
					a.BuildString("ВТОРИЧКА")
				} else if segment == jone.DoRent {
					a.BuildString("АРЕНДА")
				} else if segment == jone.DoVilla {
					a.BuildString("ЗАГОРОД")
				} else if segment == jone.DoNew {
					a.BuildString("НОВОСТРОЙКИ")
				} else if segment == jone.DoArea {
					a.BuildString("УЧАСТКИ")
				} else {
					a.BuildString("ВЫБЕРИТЕ РЕЖИМ")
				}
			}
		}
		if operator != nil {
			menuItemImg(rule, row, nm, "cabinet", "/images/menumember.png")
			if td := row.BuildTdClass("menu menu_operator"); td != nil {
				a := td.BuildItem("a", "class=menutop", "href='#'")
				a.SetAttr("onclick", fmt.Sprintf("bind_cabinet(%d)", int(rule.ItSession.IdMember)))
				a.BuildString(jone.CalculateElmText(operator))
			}
			if td := row.BuildTdClass("menu menu_operator"); td != nil {
				a := td.BuildItem("a", "class=menutop", "title=Уровень доступа", "href='#'")
				a.SetAttr("onclick", "click_role()")
				if rule.IAmAdmin() {
					a.BuildString("Адм")
				} else if rule.IAmAdvert() {
					a.BuildString("Ркл")
				} else if rule.IAmManager() {
					a.BuildString("Мдж")
				} else if rule.IAmRealtor() {
					a.BuildString("Рлт")
				} else if rule.IAmDispatch() {
					a.BuildString("Дсп")
				} else {
					a.BuildString("Роль")
				}
			}
		}
	}
	if rule.IsTechno {
	} else if segment == jone.DoTune {
	} else if lik.RegExCompare(mode, "^offer\\D") {
		menuFillOffer(rule, row, nm)
	} else if lik.RegExCompare(mode, "^member\\D") {
		menuFillMember(rule, row, nm)
	} else {
		menuFillTop(rule, row, nm)
	}
	menuFillTool(rule, row, nm)
}

//	Заполнить меню верхнего уровня
func menuFillTop(rule *repo.DataRule, row likdom.Domer, nm int) {
	segment := rule.SeekSegment("")
	mode := rule.ItPage.Locates[nm].GetMode()
	if rule.IsLogin() {
		menuItemText(rule, row, nm, "bell", "Контакты")
		if segment == jone.DoRent {
			menuItemText(rule, row, nm, "sale", "Сдать")
			menuItemText(rule, row, nm, "buy", "Снять")
		} else if segment != "" {
			menuItemText(rule, row, nm, "sale", "Продажи")
			menuItemText(rule, row, nm, "buy", "Покупки")
		}
		menuItemText(rule, row, nm, "deal", "Сделки")
	}
	if rule.IAmManager() || rule.IAmAdmin() {
		menuItemText(rule, row, nm, "client", "Клиенты")
		menuItemText(rule, row, nm, "depart", "Подразделения")
	}
	if rule.IAmManager() || rule.IAmAdmin() || mode == "cabinet" {
		menuItemText(rule, row, nm, "member", "Сотрудники")
	}
}

//	Заполнить меню заявки
func menuFillOffer(rule *repo.DataRule, row likdom.Domer, nm int) {
	id := rule.ItPage.Locates[nm].GetId()
	elm := jone.TableOffer.GetElm(id)
	if elm != nil {
		target := elm.GetString("target")
		segment := elm.GetString("segment")
		td := row.BuildTdClass("menu menu_logo")
		text := fmt.Sprintf("Заявка №%d", int(elm.Id))
		if idu := elm.GetString("idu"); idu != "" {
			text += " / " + idu
		}
		if segment == jone.DoRent && target == "sale" {
			text += " на сдачу"
		} else if target == "sale" {
			text += " на продажу"
		} else if segment == jone.DoRent && target == "buy" {
			text += " на съём"
		} else if target == "buy" {
			text += " на покупку"
		}
		if realty := jone.CalculateElmString(elm, "objectid/realty"); realty == "flat" {
			text += " квартиры"
		} else if realty == "room" {
			text += " комнаты"
		} else if realty == "plan" {
			text += " участка"
		}
		if fio := jone.CalculatePartIdText("member", elm.GetIDB("memberid")); fio != "" {
			text += ", " + fio
		}
		if nm + 1 >= rule.GetLevel() {
			rule.Title = text
		}
		td.BuildString(text)
	}
	pit := fmt.Sprintf("%d", int(id))
	menuItemText(rule, row, nm, "offershow" + pit, "Общая")
	idm := jone.CalculateElmIDB(elm, "memberid")
	if rule.IAmAdmin() ||
		repo.ProbeItMy(rule, "member", idm) ||
 		rule.IAmManager() && repo.ProbeItDep(rule, "member", idm) {
		menuItemText(rule, row, nm, "offerstaff"+pit, "Служебная")
		if target := jone.CalculatePartIdString("offer", id, "target"); target == "sale" {
			menuItemText(rule, row, nm, "offerfiles"+pit, "Файлы")
			menuItemText(rule, row, nm, "offerpromo"+pit, "Продвижение")
			menuItemText(rule, row, nm, "offerlife"+pit, "История")
		}
		menuItemText(rule, row, nm, "offerdeal"+pit, "Сделка")
	}
}

//	Заполнить меню сотрудника
func menuFillMember(rule *repo.DataRule, row likdom.Domer, nm int) {
	id := rule.ItPage.Locates[nm].GetId()
	elm := jone.TableMember.GetElm(id)
	if elm != nil {
		td := row.BuildTdClass("menu menu_logo")
		text := jone.CalculateElmText(elm)
		if text == "" {
			text = fmt.Sprintf("Id%03d", int(elm.Id))
		}
		td.BuildString("Сотрудник " + text)
	}
	pit := fmt.Sprintf("%d", int(id))
	menuItemText(rule, row, nm, "membercard" + pit, "Личный кабинет")
}

//	Заполнить меню инструментов
func menuFillTool(rule *repo.DataRule, row likdom.Domer, nm int) {
	menuItemText(rule, row, nm,"","")
	if nm == 0 {
		now := time.Now()
		text := "<nobr><span id=srvtime title='Точное время'>" + now.Format("15:04") + "</span></nobr>"
		row.BuildTdClass("menu srvtime").BuildString(text)
		version := ""
		if diff := now.Sub(repo.TimeBin); diff > time.Hour * 48 {
			version = repo.TimeBin.Format("02/01/2006")
		} else if repo.TimeBin.Month() != now.Month() || repo.TimeBin.Day() != now.Day() {
			version = repo.TimeBin.Format("02/01 15:04")
		} else {
			version = repo.TimeBin.Format("15:04")
		}
		if version != "" {
			menuItemText(rule, row, nm, "project/Последние обновления", "&nbsp;<span class=rls>от " + version + "</span>")
		}
		if rule.IAmAdmin() {
			menuItemImg(rule, row, nm, "tune/Настройки", "/images/menuopt.png")
		}
	}
	if nm + 1 == rule.GetLevel() && !rule.IsTechno {
		menuItemImg(rule, row, nm, "print/Распечатать", "/images/menuprint.png")
	}
	menuItemImg(rule, row, nm, "exit/Закрыть окно", "/images/menuexit.png")
}

//	Добавить элемент с изображением
func menuItemImg(rule *repo.DataRule, row likdom.Domer, nm int, cmd string, pic string) {
	img := likdom.BuildUnpairItem("img","src", pic, "width=20", "height=20")
	menuItemText(rule, row, nm, cmd, img.ToString())
}

//	Добавить элемент с текстом
func menuItemText(rule *repo.DataRule, row likdom.Domer, nm int, cmd string, text string) {
	loc := rule.ItPage.Locates[nm]
	mode := loc.GetMode()
	cls := "menu"
	tip := ""
	if cmd != "" {
		if match := lik.RegExParse(cmd, "(.*)/(.*)"); match != nil {
			cmd = match[1]
			tip = match[2]
		}
		if cmd != "" {
			cls += " menu_act"
		}
		if mode != "" {
			if strings.HasPrefix(cmd, mode) || strings.Contains(cmd, "/"+mode) {
				if lik.RegExCompare(text, "<img") {
					cls += "" + " menu_pel"
				} else {
					cls += "" + " menu_sel"
				}
			}
		}
	} else if text == "" {
		cls += " menu_sep"
	}
	td := row.BuildTdClass(cls)
	td.BuildString(text)
	if cmd == "exit" && nm == 0 {
		td.SetAttr("onclick", fmt.Sprintf("click_exit()"))
	} else if cmd == "print" {
		td.SetAttr("onclick", fmt.Sprintf("click_print()"))
	} else if cmd == "project" {
		td.SetAttr("onclick", fmt.Sprintf("click_project()"))
	} else if cmd == "cabinet" {
		td.SetAttr("onclick", fmt.Sprintf("bind_cabinet(%d)", int(rule.ItSession.IdMember)))
	} else if cmd == "role" {
		td.SetAttr("onclick", fmt.Sprintf("click_role('%s')", rule.GetMember().GetString("role")))
	} else if cmd != "" {
		td.SetAttr("onclick", fmt.Sprintf("click_menu('%s','%s')", mode, cmd))
	}
	if tip != "" {
		td.SetAttr("title", tip)
	}
}

