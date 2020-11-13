package repo

import (
	"bitbucket.org/961961/tsan/jone"
	"bitbucket.org/shaman/lik"
	"bitbucket.org/shaman/lik/likbase"
)

//	Дескриптор карты
type DataMap struct {
	CenterX,CenterY float64		//	Координаты центра
	Zoom		float64			//	Фактор масштаба
	Points		lik.Lister		//	Список точек
}

//	Конструктор карты по-умолчанию
func BuildMapDefault(rule *DataRule) *DataMap {
	dap := &DataMap{}
	dap.Zoom = 10.0
	dap.CenterX, dap.CenterY = 54.60, 39.60
	return dap
}

//	Конструктор карты
func BuildMap(rule *DataRule, elm *likbase.ItElm) *DataMap {
	emap := jone.CalculateElmSet(elm, "objectid/map")
	if emap == nil { return nil }
	dap := BuildMapDefault(rule)
	if parm := emap.GetFloat("zoom"); parm > 0 {
		dap.Zoom = parm
	}
	if parm := emap.GetFloat("centerx"); parm != 0 {
		dap.CenterX = parm
	}
	if parm := emap.GetFloat("centery"); parm != 0 {
		dap.CenterY = parm
	}
	dap.Points = lik.BuildList()
	if lpt := emap.GetList("points"); lpt != nil {
		for npt := 0; npt < lpt.Count(); npt++ {
			if pt := lpt.GetSet(npt); pt != nil {
				dap.Points.AddItems(pt.GetFloat("x"), pt.GetFloat("y"))
			}
		}
	}
	return dap
}

