package staff

import (
	"bitbucket.org/961961/tsan/control"
	"bitbucket.org/961961/tsan/jone"
	"bitbucket.org/961961/tsan/repo"
	"bitbucket.org/961961/tsan/show"
	"bitbucket.org/shaman/lik"
	"bitbucket.org/shaman/lik/likdom"
	"fmt"
	"strings"
)

//	Дескриптор карты заявки
type MapControl struct {
	control.DataControl
	itIsControl bool
	itIsModify	bool
	itIsKeeping bool
	itIsRedraw	bool
}

//	Интерфейс маршализации
type dealMapMarshal struct {
	It	*MapControl
}
func (it *dealMapMarshal) Run(rule *repo.DataRule) {
	it.It.MapMarshal(rule)
}

//	Интерфейс команд
type dealMapExecute struct {
	It	*MapControl
}
func (it *dealMapExecute) Run(rule *repo.DataRule, cmd string, data lik.Seter) {
	it.It.MapExecute(rule, cmd, data)
}

//	Конструктор дескриптора карты
func BuildMap(rule *repo.DataRule, main string, id lik.IDB, itiscontrol bool) *MapControl {
	it := &MapControl{}
	it.itIsControl = itiscontrol
	it.ControlInitializeZone(main, id, "map")
	it.ItExecute = &dealMapExecute{it}
	it.ItMarshal = &dealMapMarshal{it}
	return it
}

//	Маршализация карты
func (it *MapControl) MapMarshal(rule *repo.DataRule) {
	if it.itIsRedraw {
		it.itIsRedraw = false
		rule.StoreItem(it.buildShowTools(rule))
		rule.SetResponse("now", "_function_map_redraw")
	}
}

//	Выполнение команд карты
func (it *MapControl) MapExecute(rule *repo.DataRule, cmd string, data lik.Seter) {
	if cmd == "showmap" {
		it.cmdShowMap(rule)
	} else if cmd == "toedit" {
		it.itIsModify = true
		it.itIsKeeping = true
		rule.StoreItem(it.buildShowTools(rule))
		rule.SetResponse(cmd,"_function_map_redraw")
	} else if cmd == "tocancel" {
		it.itIsModify = false
		rule.StoreItem(it.buildShowTools(rule))
		rule.SetResponse(cmd, "_function_map_redraw")
	} else if cmd == "write" {
		it.cmdWriteMap(rule, data)
		it.itIsModify = false
		it.itIsRedraw = true
	} else if cmd == "snap" {
		code := lik.StringFromXS(data.GetString("snap"))
		_ = code
	}
}

//	Отображение карты
func (it *MapControl) BuildShow(rule *repo.DataRule, sx int, sy int) likdom.Domer {
	spc := likdom.BuildSpace()
	height := sy
	if it.itIsControl {
		spc.AppendItem(it.buildShowTools(rule))
		height -= 20
	}
	spc.BuildDiv("id=map_data", fmt.Sprintf("style='width: %dpx; height: %dpx'", sx, height))
	//spc.BuildItem("canvas","id=map_data", fmt.Sprintf("style='width: %dpx; height: %dpx'", sx, height))
	code := fmt.Sprintf("ymaps.ready(function() { map_bind_zone('%s','%s','%s') });", it.GetMode(), "map", "map_data")
	spc.AppendItem(show.BuildRunScript(code,0))
	return spc
}

//	Отображение инструментов карты
func (it *MapControl) buildShowTools(rule *repo.DataRule) likdom.Domer {
	tbl := likdom.BuildItemClass("table","mapcmd", "id=maptools")
	row := tbl.BuildTr()
	if it.itIsModify {
		row.BuildTdClass("mapcmd mappass").AppendItem(show.LinkTextIdProc("", "&nbsp;Изменить&nbsp;", "mapedit", "map_setedit()"))
		row.BuildTdClass("mapcmd").AppendItem(show.LinkTextIdProc("", "&nbsp;Записать&nbsp;", "mapwrite", "map_setwrite()"))
		row.BuildTdClass("mapcmd").AppendItem(show.LinkTextIdProc("", "&nbsp;Стереть&nbsp;", "mapclear", "map_setclear()"))
		row.BuildTdClass("mapcmd").AppendItem(show.LinkTextIdProc("", "&nbsp;Отменить&nbsp;", "mapcancel", "map_setcancel()"))
		row.BuildTd().BuildString("&nbsp;<i>Правая кнопка - отметить</i>")
	} else {
		row.BuildTdClass("mapcmd").AppendItem(show.LinkTextIdProc("", "&nbsp;Изменить&nbsp;", "mapedit", "map_setedit()"))
	}
	row.BuildTd("width=100%")
	row.BuildTdClass("mapcmd").AppendItem(show.LinkTextIdProc("", "&nbsp;Фото&nbsp;", "mapsnap", "map_snap()"))
	return tbl
}

//	Отображение карты
func (it *MapControl) cmdShowMap(rule *repo.DataRule) {
	elm := jone.GetElm("offer", it.IdMain)
	dap := repo.BuildMap(rule, elm)
	if dap == nil { dap = repo.BuildMapDefault(rule) }
	rule.SetResponse(lik.BuildList(dap.CenterX, dap.CenterY),"center")
	rule.SetResponse(dap.Zoom, "zoom")
	if dap.Points != nil && dap.Points.Count() > 0 {
		rule.SetResponse(dap.Points,"points")
	}
	rule.SetResponse(it.itIsModify, "isedit")
	rule.SetResponse(!it.itIsKeeping, "isfull")
	it.itIsKeeping = false
}

//	Запись изменений карты
func (it *MapControl) cmdWriteMap(rule *repo.DataRule, data lik.Seter) {
	if elm := jone.GetElm("offer", it.IdMain); elm != nil {
		for _,parm := range []string {"zoom","centerx","centery"} {
			val := data.GetString(parm)
			fval := lik.StrToFloat(lik.StringFromXS(val))
			jone.SetElmValue(elm, fval, "objectid/map/"+parm)
		}
		points := lik.BuildList()
		if lpt := strings.Split(lik.StringFromXS(data.GetString("points")), ","); lpt != nil {
			for np := 0; np + 1 < len(lpt); np += 2 {
				x := lik.StrToFloat(lpt[np])
				y := lik.StrToFloat(lpt[np + 1])
				if x != 0 && y != 0 {
					points.AddItems(lik.BuildSet("x", x, "y", y))
				}
			}
		}
		if points.Count() > 0 {
			jone.SetElmValue(elm, points,"objectid/map/points")
		} else {
			jone.SetElmValue(elm, nil,"objectid/map/points")
		}
	}
}

