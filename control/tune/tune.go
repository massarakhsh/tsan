//	Модуль окна настроек
package tune

import (
	"bitbucket.org/961961/tsan/control"
	"bitbucket.org/961961/tsan/fancy"
	"bitbucket.org/961961/tsan/jone"
	"bitbucket.org/961961/tsan/repo"
	"bitbucket.org/961961/tsan/show"
	"bitbucket.org/shaman/lik"
	"bitbucket.org/shaman/lik/likapi"
	"bitbucket.org/shaman/lik/likdom"
	"fmt"
	"strings"
)

//	Дескриптор окна настроек
type TuneControl struct {
	control.DataControl			//	На обычном контроллере
	fancy.TableFancy
	SxTree, SxPart int			//	Ширина окна дерева и разделов
	Path           string		//	Путь настроек
	itDiv          string		//	Имя раздела
	itGen          string		//	Имя генератора
	itEnt          string		//	Имя сущности
	itWhat         string		//	Текущее имя
	itTags			bool		//	Признак тегов
	listIdParts    []string		//	Список элементов
}

//	Интерфейс команд
type dealTuneExecute struct {
	It	*TuneControl
}
func (it *dealTuneExecute) Run(rule *repo.DataRule, cmd string, data lik.Seter) {
	it.It.TuneExecute(rule, cmd, data)
}

//	Интерфейс обновления окна
type dealTuneUpdate struct {
	It	*TuneControl
}
func (it *dealTuneUpdate) Run(rule *repo.DataRule) {
	it.It.TuneUpdate(rule)
}

//	Дескриптор опции
type EntOpt struct {
	Tag		int
	Title	string
	Tip		string
}

//	Список опций
var ListOpt = []EntOpt{
	{0,"",""},
	{jone.TagGrid,"ТБЛ","В таблице"},
	{jone.TagForm,"ФРМ","В форме"},
	{jone.TagMust,"ДОЛ","Должен быть"},
	{0,"",""},
	{jone.TagHide,"СКР","Скрытый"},
	{jone.TagShow,"ВИД","Видимый"},
	{jone.TagEdit,"РЕД","Редактируемый"},
	{0,"",""},
	{jone.TagSale,"ПРД","Продажи"},
	{jone.TagBuy,"ПОК","Покупки"},
	{0,"",""},
	{jone.TagFlat,"КВТ","Квартиры"},
	{jone.TagHouse,"ДОМ","Дома"},
	{jone.TagLand,"ЗЕМ","Замля"},
}

//	Конструктор дескриптора настроек
func BuildTune(rule *repo.DataRule) *TuneControl {
	it := &TuneControl{}
	it.ControlInitialize("tune", 0)
	it.ItExecute = &dealTuneExecute{it}
	it.ItUpdate = &dealTuneUpdate{it}
	it.PartInitialize(rule)
	return it
}

//	Обновление окна
func (it *TuneControl) TuneUpdate(rule *repo.DataRule) {
	path := lik.StringFromXS(rule.GetContext("path"))
	if path != "" && it.Path != path {
		it.Path = path
	}
	rule.OnChangeData()
	it.itDiv = ""
	it.itGen = ""
	it.itEnt = ""
	it.itWhat = ""
	it.itTags = false
	parts := lik.PathToNames(it.Path)
	if len(parts) > 0 {
		div := parts[0]
		parts = parts[1:]
		if repo.SystemFindGen(div) != nil {
			it.itDiv = "devsys"
			it.itWhat = "gen"
			it.itGen = div
			if len(parts) > 0 {
				it.itWhat = "ent"
				it.itEnt = parts[0]
				parts = parts[1:]
			}
		} else {
			it.itWhat = "div"
			it.itDiv = div
		}
	}
}

//	Выполнение команд
func (it *TuneControl) TuneExecute(rule *repo.DataRule, cmd string, data lik.Seter) {
	if cmd == "all" || cmd == "tune" {
		it.TuneExecute(rule, rule.Shift(), data)
	} else if cmd == "part" {
		it.PartExecute(rule, rule.Shift(), data)
	} else if cmd == "tree" {
		it.TreeExecute(rule, rule.Shift(), data)
	} else if cmd == "system" {
		repo.SystemExecuteCmd(rule.Shift())
		rule.SetGoPart("/tune")
		rule.OnChangeData()
	}
}

//	Отображение окна
func (it *TuneControl) BuildShow(rule *repo.DataRule, sx int, sy int) likdom.Domer {
	div := likdom.BuildDivClassId("roll_data", "tune_data")
	it.SxTree = 320
	it.SxPart = sx - it.SxTree
	row := div.BuildTableClass("fill").BuildTr()
	if td := row.BuildTdClass("section", control.MiniMax(it.SxTree, sy)...); true {
		td.AppendItem(it.buildDataTree(rule, it.SxTree -control.BD, sy -control.BD))
	}
	if td := row.BuildTdClass("section", control.MiniMax(it.SxPart, sy)...); true {
		td.AppendItem(it.buildDataTune(rule, it.SxPart -control.BD, sy -control.BD))
	}
	return div
}

//	Отображение окна раздела
func (it *TuneControl) buildDataTune(rule *repo.DataRule, sx int, sy int) likdom.Domer {
	tbl := likdom.BuildItem("table")
	tbl.BuildTrTd("colspan=2").BuildString("Путь: " + strings.ReplaceAll(it.Path, "__", "/"))
	tbl.BuildTrTd("colspan=2").BuildString("<hr>")
	if it.itWhat == "" {
		tbl.BuildTrTd("colspan=2").BuildString("Корень директории")
	} else if it.itDiv == "command" {
		it.showSystemCommands(rule, tbl)
	} else if it.itDiv == "session" {
		it.showSessions(rule, tbl, sx, sy - 80)
	//} else if it.itDiv == "system" && it.itWhat == "gen" {
	} else if it.itWhat == "gen" {
		it.showTableGen(tbl, sx, sy - 80)
	} else if it.itWhat == "ent" {
		it.showTableEnt(tbl, sx, sy - 80)
	}
	return tbl
}

//	Отображение списка команд
func (it *TuneControl) showSystemCommands(rule *repo.DataRule, tbl likdom.Domer) {
	for nc := 0; nc < 1000; nc++ {
		cmd,url,txt,prc,inp := "", "", "", "", false
		if cc := 0; nc == cc {
			cmd,txt = "reload", "Перегрузить базу данных"
		} else if cc++; nc == cc {
		} else if cc++; nc == cc {
			url,txt = likapi.BuildUrl(rule.ItPage.GetPageId(), "/front/database.json"), "Сохранить базу данных в файл"
		} else if cc++; nc == cc {
		} else if cc++; nc == cc {
			url = likapi.BuildUrl(rule.ItPage.GetPageId(), "/front/database.load?_mf=1")
			txt,inp = "Загрузить базу данных", true
		} else {
			break
		}
		tr := tbl.BuildTr()
		if cmd == "" && url == "" && prc == "" {
			tr.BuildTd("colspan=2").BuildString("<hr>")
		} else {
			tr.BuildTd().BuildUnpairItem("img", "src", rule.BuildUrl("/rast/pix/z0000.pix"))
			if inp {
				td := tr.BuildTd()
				td.BuildString(txt + "<br>")
				form := td.BuildItem("form", "method=post", "action", url, "enctype", "multipart/form-Data")
				form.BuildUnpairItem("input", "type=file", "name=file")
				form.BuildUnpairItem("input", "type=submit", "value", txt)
			} else if prc != "" {
				tr.BuildTd().AppendItem(show.LinkTextProc("syscmd", txt, prc))
			} else if cmd != "" {
				tr.BuildTd().AppendItem(show.LinkTextProc("syscmd", txt, fmt.Sprintf("front_get('/system/%s')", cmd)))
			} else {
				tr.BuildTd().BuildItem("a", "class=syscmd", "target=_blank", "href", url).BuildString(txt)
			}
		}
	}
}

//	Отображение окна генератора
func (it *TuneControl) showTableGen(tbl likdom.Domer, sx int, sy int) {
	if !it.lineShowGen(tbl) {
		return
	}
	it.showListPart(tbl, sx, sy)
}

//	Отображение окна сущности
func (it *TuneControl) showTableEnt(tbl likdom.Domer, sx int, sy int) {
	if !it.lineShowEnt(tbl) {
		return
	}
	it.showListPart(tbl, sx, sy)
}

//	Отображение списка разделов
func (it *TuneControl) showListPart(tbl likdom.Domer, sx int, sy int) {
	it.SetSize(sx, sy)
	div := tbl.BuildTrTd("colspan=2").BuildDivClass("listpart")
	div.AppendItem(show.BuildFancyGrid(it.Main,"part"))
}

//	Отображение сессий
func (it *TuneControl) showSessions(rule *repo.DataRule, tbl likdom.Domer, sx int, sy int) {
	tbl.BuildTrTd("colspan=2").BuildString("Сессии:")
	sessions := []likapi.DataSessioner{}
	for _, session := range likapi.ListSessions {
		sessions = append(sessions, session)
	}
	for ns := len(sessions) - 2; ns >= 0; ns-- {
		for ks := 0; ks <= ns; ks++ {
			if sessions[ks].GetSelf().AtLast.Before(sessions[ks+1].GetSelf().AtLast) {
				sessions[ks], sessions[ks+1] = sessions[ks+1], sessions[ks]
			}
		}
	}
	for _, session := range sessions {
		row := tbl.BuildTr()
		self := session.GetSelf()
		row.BuildTd().BuildString(self.AtLast.Format("15:04:05.&nbsp; "))
		td := row.BuildTd()
		td.BuildString(fmt.Sprintf("IP=%s", self.IP))
		td.BuildString(fmt.Sprintf(", pages=%d", self.CP))
		td.BuildString(fmt.Sprintf(", URI=%s", self.Uri))
		row = tbl.BuildTr()
		row.BuildTd()
		it.showSessionPages(rule, self.IdSession, row.BuildTd())
	}
}

//	Отображение страниц сессий
func (it *TuneControl) showSessionPages(rule *repo.DataRule, idsession int, td likdom.Domer) {
	tbl := td.BuildTable()
	pages := []likapi.DataPager{}
	for _, page := range likapi.ListPages {
		if page.GetSessionId() == idsession {
			pages = append(pages, page)
		}
	}
	for ns := len(pages) - 2; ns >= 0; ns-- {
		for ks := 0; ks <= ns; ks++ {
			if pages[ks].GetAt().Before(pages[ks+1].GetAt()) {
				pages[ks], pages[ks+1] = pages[ks+1], pages[ks]
			}
		}
	}
	for np, pg := range pages {
		row := tbl.BuildTr()
		if np >= 12 {
			row.BuildTd("colspan=2").BuildString(fmt.Sprintf("........... ещё %d", len(pages)-12))
			break
		}
		page := pg.(repo.DataPager).GetSelf()
		row.BuildTd().BuildString(page.GetAt().Format("15:04:05.&nbsp; "))
		//row.BuildTd().BuildString(page.RootUrl)
	}
}

//	Отображение части генератора
func (it *TuneControl) lineShowGen(tbl likdom.Domer) bool {
	if gen := repo.SystemFindGen(it.itGen); gen != nil {
		it.showTextCell(tbl, "Раздел:", gen.Name)
	} else {
		tbl.BuildTrTd("colspan=2").BuildString(fmt.Sprintf("Раздел '%s' не найден", it.itGen))
		return false
	}
	return true
}

//	Отображение части сущности
func (it *TuneControl) lineShowEnt(tbl likdom.Domer) bool {
	if !it.lineShowGen(tbl) {
		return false
	}
	ent := repo.SystemFindGenEnt(it.itGen, it.itEnt)
	if ent == nil {
		return false
	}
	name := jone.CalculateElmString(ent.It,"name")
	if name == "" {
		name = "(без имени)"
	}
	name += " (" + ent.Key + ")"
	it.showTextCell(tbl, "Объект:", name)
	return true
}

//	Отою\бражение ячейки
func (it *TuneControl) showTextCell(tbl likdom.Domer, title string, value string) {
	it.showItemsCell(tbl, title, likdom.BuildString(value))
}

//	Отображение строки с ячейками
func (it *TuneControl) showItemsCell(tbl likdom.Domer, title string, items ...likdom.Domer) {
	tr := tbl.BuildTr()
	tr.BuildTdClass("tune_title").BuildString(title)
	tr.BuildTdClass("tune_data").AppendItem(items...)
}

