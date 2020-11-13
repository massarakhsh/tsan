//	Модуль интерфейса отображения
package front

import (
	"bitbucket.org/961961/tsan/control"
	"bitbucket.org/961961/tsan/control/controls"
	"bitbucket.org/961961/tsan/control/project"
	"bitbucket.org/961961/tsan/one"
	"bitbucket.org/961961/tsan/repo"
	"bitbucket.org/shaman/lik"
	"bitbucket.org/shaman/lik/likdom"
	"fmt"
	"strings"
	"time"
)

//	Конструктор объекта страницы
func FrontPage(rule *repo.DataRule) (int, likdom.Domer) {
	rule.SayInfo("Page: " + rule.RootUrl)
	html := rule.InitializePage(one.Version)
	if head, _ := html.GetDataTag("head"); head != nil {
		urlmap := fmt.Sprintf("https://api-maps.yandex.ru/2.1/?apikey=%s&amp;lang=ru_RU", one.YandexMapKey)
		head.BuildString("<script type='text/javascript' src='" + urlmap + "'></script>")
		head.BuildString("<link rel='stylesheet' href='/lib/fancy.css'/>")
		head.BuildString("<script type='text/javascript' src='/lib/fancy.full.js'></script>")
		head.BuildString("<link rel='stylesheet' href='/lib/dropzone.css'/>")
		head.BuildString("<script type='text/javascript' src='/lib/dropzone.js'></script>")
		head.BuildString("<link rel='stylesheet' href='/lib/fotorama.css'/>")
		head.BuildString("<script type='text/javascript' src='/lib/fotorama.js'></script>")
		head.BuildString("<link rel='stylesheet' href='/js/styles.css?" + one.Version + "'/>")
		head.BuildString("<link rel='stylesheet' href='/js/fancy.css?" + one.Version + "'/>")
		head.BuildString("<script type='text/javascript' src='/js/script.js?" + one.Version + "'></script>")
		head.BuildString("<script type='text/javascript' src='/js/render.js?" + one.Version + "'></script>")
		head.BuildString("<script type='text/javascript' src='/js/fancy.js?" + one.Version + "'></script>")
		head.BuildString("<script type='text/javascript' src='/js/map.js?" + one.Version + "'></script>")
	}
	if body, _ := html.GetDataTag("body"); body != nil {
		if script := body.BuildItem("script"); script != nil {
			code := "script_start();\r\n"
			code += "Fancy.MODULESDIR = '/lib/modules/';\r\n"
			script.BuildString("jQuery(document).ready(function () { " + code + " });")
		}
		BuildPage(rule, body)
	}
	if head, _ := html.GetDataTag("head"); head != nil {
		if rule.Title == "" {
			rule.Title = "РИЭЛТОР. ЦАН"
		}
		head.BuildItem("title").BuildString(rule.Title)
	}
	return 200, html
}

//	Заполнение страницы
func BuildPage(rule *repo.DataRule, pater likdom.Domer) {
	if !rule.IsLogin() {
		rule.BindSession()
	}
	selpath := "/"
	if locs := rule.GetPath(); len(locs) > 0 {
		selpath += strings.Join(locs, "/")
	}
	mainSetPath(rule, selpath)
	mainPage(rule, pater)
}

//	Обработка команд
func FrontExecute(rule *repo.DataRule) lik.Seter {
	rule.IsJson = true
	start := time.Now()
	frontLog(rule)
	frontResizing(rule)
	if rule.IsShift("menu") {
		doMenu(rule)
	} else if rule.IsShift("system") {
		doSystem(rule)
	} else if rule.IsShift("database.json") {
		doSaveDatabase(rule)
	} else if rule.IsShift("database.load") {
		doLoadDatabase(rule)
	} else {
		doCommand(rule)
	}
	if rule.IsItChangePage() {
		frontRePage(rule)
	} else if rule.IsItChangeData() {
		frontReData(rule)
	} else if rule.ItPage.PathNeed != rule.ItPage.PathLast {
		frontReDraw(rule)
	}
	if rule.ItPage.PathClient != rule.ItPage.PathLast {
		setClientPath(rule)
	}
	_ = start
	//frontEndLog(rule, start)
	response := rule.GetAllResponse()
	if response == nil || response.Count() == 0 {
		response = lik.BuildSet("success=true")
	}
	return response
}

//	Обработка маршализации
func MarshalExecute(rule *repo.DataRule) lik.Seter {
	rule.IsJson = true
	frontResizing(rule)
	if !rule.ItPage.GetTrust() {
		if rule.BindSession() {
			rule.SetResponse(rule.ItPage.GetPageId(), "_page")
			frontDoReload(rule)
		}
	} else if rule.IsItChangePage() {
		frontDoReload(rule)
	} else if rule.IsNeedGoPath {
		frontDoGoPath(rule)
	} else if rule.IsNeedReload {
		frontDoReload(rule)
	} else if _,ctrl := GetFront(rule); ctrl != nil {
		ctrl.RunMarshal(rule)
	}
	if rule.ItPage.PathClient != rule.ItPage.PathLast {
		setClientPath(rule)
	}
	if rule.Title != "" {
		rule.SetTitle(rule.Title)
	}
	response := rule.GetAllResponse()
	if response != nil && response.Count() > 0 {
	} else {
		response = lik.BuildSet("success=true")
	}
	return response
}

//	Вывод в лог
func frontLog(rule *repo.DataRule) {
	path := "/" + strings.Join(rule.GetPath(), "/")
	if member := rule.GetMember(); member != nil {
		rule.SayInfo("Front: " + member.GetString("login") + " " + path)
		if strings.HasSuffix(path, "showpage") {
			rule.SayInfo(fmt.Sprintf("Filter: %d", len(rule.GetContext("filter"))))
		}
	}
}

//	Вывод результатов в лог
func frontEndLog(rule *repo.DataRule, start time.Time) {
	delay := fmt.Sprintf("Delay: %4.2f", float64(time.Now().Sub(start)/time.Millisecond) / 1000.0)
	rule.SayInfo(delay)
}

//	Контроль размера окна клиента
func frontResizing(rule *repo.DataRule) {
	if !rule.SeekPageSize() {
		rule.OnChangePage()
	}
}

//	Обработка команд меню
func doMenu(rule *repo.DataRule) {
	if level,control := getContext(rule); level >= 0 && control != nil {
		cmd := rule.Shift()
		if cmd == "exit" {
			rule.SetPageExit(level)
		} else {
			rule.SetPagePart(level, cmd)
		}
	}
}

//	Обработка стандартных команд
func doCommand(rule *repo.DataRule) {
	if _,ctrl := getContext(rule); ctrl != nil {
		cmd := rule.Shift()
		doControl(rule, ctrl, cmd)
	}
}

//	Обработка команд контроллера
func doControl(rule *repo.DataRule, ctrl control.Controller, cmd string) {
	if ctrl != nil {
		if rule.IsMethod("POST") || rule.IsMethod("post") {
			ctrl.RunExecute(rule, cmd, collectParms(rule, "up_"))
		} else {
			ctrl.RunExecute(rule, cmd,nil)
		}
	}
}

//	Определить контекст
func getContext(rule *repo.DataRule) (int, control.Controller) {
	main := rule.Shift()
	level,ctrl := FindControl(rule, main)
	if ctrl == nil {
		if main == "command" {
			var id lik.IDB
			if len(rule.ItPage.Locates) > 0 {
				id = rule.ItPage.Locates[0].GetId()
			}
			ctrl = controls.BuildCommand(rule, id)
			rule.RegControl(ctrl)
		} else if main == "print" {
			var id lik.IDB
			if len(rule.ItPage.Locates) > 0 {
				id = rule.ItPage.Locates[0].GetId()
			}
			ctrl = controls.BuildPrint(rule, id)
			rule.RegControl(ctrl)
		} else if main == "project" {
			ctrl = collectproject.BuildProject(rule)
			rule.RegControl(ctrl)
		} else {
			level, ctrl = GetFront(rule)
		}
	}
	return level, ctrl
}

//	Системные команды
func doSystem(rule *repo.DataRule) {
	if cmd := rule.Shift(); cmd == "nonono" {
	} else {
		repo.SystemExecuteCmd(cmd)
	}
	rule.OnChangePage()
}

//	Сохранить базу данных
func doSaveDatabase(rule *repo.DataRule) {
	system := repo.SystemSaveDataBase()
	rule.ResultFormat = true
	for _,set := range(system.Values()) {
		rule.SetResponse(set.Val, set.Key)
	}
}

//	Восстановить базу данных
func doLoadDatabase(rule *repo.DataRule) {
	files := rule.GetBuffers()
	for _,buf := range(files) {
		str := string(buf)
		if json := lik.SetFromRequest(str); json != nil {
			repo.SystemLoadDataBase(json)
		}
	}
	rule.IsJson = false
	rule.IsNeedReload = true
}

//	Обновить страницу
func frontRePage(rule *repo.DataRule) {
	rule.OnChangePage()
	frontReDraw(rule)
}

//	Обновить данные на странице
func frontReData(rule *repo.DataRule) {
	rule.OnChangeData()
	frontReDraw(rule)
}

//	Перерисовать что требуется
func frontReDraw(rule *repo.DataRule) {
	if rule.ItPage.PathNeed != rule.ItPage.PathLast {
		mainSetPath(rule, rule.ItPage.PathNeed)
	}
	mainStore(rule)
}

//	Перерисовать клиент
func frontDoReload(rule *repo.DataRule) {
	rule.SetResponse("", "_function_lik_reload")
}

//	Перейти на другую страницу
func frontDoGoPath(rule *repo.DataRule) {
	rule.SetResponse(rule.ItPage.PathNeed, "_function_lik_go_part")
}

//	Установить путь на клиенте
func setClientPath(rule *repo.DataRule) {
	rule.ItPage.PathClient = rule.ItPage.PathLast
	rule.PushOnPart(rule.ItPage.PathClient)
}

//	Клонировать страницу
func PageClone(rule *repo.DataRule, parm string) {
	rule.SetGoPart(rule.BuildFullPath(parm))
}

//	Собрать параметры по префиксу
func collectParms(rule *repo.DataRule, prefix string) lik.Seter {
	parms := lik.BuildSet()
	if context := rule.GetAllContext(); context != nil {
		for _,set := range(context.Values()) {
			if strings.HasPrefix(set.Key, prefix) && set.Val != nil {
				str := set.Val.ToString()
				parms.SetItem(str, set.Key[len(prefix):])
			}
		}
	}
	return parms
}

//	Найти контроллер
func FindControl(rule *repo.DataRule, main string) (int, control.Controller) {
	rule.ItPage.Sync.Lock()
	level := -1
	var ctrl control.Controller
	if ctrl == nil {
		for lev, loc := range rule.ItPage.Locates {
			if main == loc.GetMode() {
				level = lev
				ctrl = loc.(control.Controller)
				break
			}
		}
	}
	if ctrl == nil {
		for pos, loc := range rule.ItPage.Sheets {
			if main == loc.GetMode() {
				if pos > 0 {
					for p := pos; p > 0; p-- {
						rule.ItPage.Sheets[p] = rule.ItPage.Sheets[p-1]
					}
					rule.ItPage.Sheets[0] = loc
				}
				ctrl = loc.(control.Controller)
				break
			}
		}
	}
	rule.ItPage.Sync.Unlock()
	return level, ctrl
}

//	Найти контроллер первого плана
func GetFront(rule *repo.DataRule) (int, control.Controller) {
	rule.ItPage.Sync.Lock()
	level := -1
	var ctrl control.Controller
	if level = len(rule.ItPage.Locates); level > 0 {
		ctrl = rule.ItPage.Locates[level-1].(control.Controller)
	}
	rule.ItPage.Sync.Unlock()
	return level, ctrl
}

