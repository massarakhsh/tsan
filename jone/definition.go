package jone

import (
	"github.com/massarakhsh/lik"
	"github.com/massarakhsh/lik/likbase"
	"fmt"
)

const (
	ItAll      = "all"
	ItActive   = "active"
	ItError    = "error"
	ItReady    = "ready"
	ItCancel   = "cancel"
	ItDone     = "done"
	ItDep      = "dep"
	ItMy       = "my"
	ItAdmin    = "admin"
	ItAdvert   = "advert"
	ItManager  = "manager"
	ItRealtor  = "realtor"
	ItDispatch = "dispatch"
	DoTune		= "tune"
	DoCall		= "call"
	DoSecond	= "second"
	DoRent		= "rent"
	DoNew		= "new"
	DoVilla		= "villa"
	DoArea		= "area"
)

//	Определит наименование объекта
func DefinitionObject(elm *likbase.ItElm) string {
	value := fmt.Sprintf("ID%03d", int(elm.Id))
	if val := CalculateElmString(elm, "define/rooms"); val != "" {
		value += ", тип:" + val
	}
	return value
}

//	Определит наименование контакта
func DefinitionBell(elm *likbase.ItElm) string {
	value := fmt.Sprintf("ID%03d", int(elm.Id))
	return value
}

//	Определит наименование заявки
func DefinitionOffer(elm *likbase.ItElm) string {
	value := fmt.Sprintf(" №%03d", int(elm.Id))
	return value
}

//	Определит наименование сотрудника
func DefinitionMember(elm *likbase.ItElm) string {
	value := CalculateElmString(elm,"family")
	if val := CalculateElmString(elm,"paterly"); val != "" {
		value = lik.SubString(val, 0, 1) + "." + value
	}
	if val := CalculateElmString(elm,"namely"); val != "" {
		value = lik.SubString(val, 0, 1) + "." + value
	}
	if value == "" {
		value = CalculateElmString(elm,"login")
	}
	return value
}

//	Определит наименование подразделения
func DefinitionDepart(elm *likbase.ItElm) string {
	value := CalculateElmString(elm,"name")
	return value
}

//	Определит наименование клиента
func DefinitionClient(elm *likbase.ItElm) string {
	value := CalculateElmString(elm,"namely")
	if val := CalculateElmString(elm,"paterly"); val != "" {
		value += " " + val
	}
	if val := CalculateElmString(elm,"family"); val != "" {
		value += " " + val
	}
	return value
}

//	Нормализовать номер телефона
func NormalizePhone(phone string) string {
	pho := lik.RegExReplace(phone, "\\D", "")
	if lenpho := len(pho); lenpho == 11 {
		if pho[0] == '7' || pho[0] == '8' {
			pho = pho[1:]
		}
	}
	return pho
}

//	Нормализовать плавающее число
func NormalizeFloat(val float64, nd int) float64 {
	snd := fmt.Sprintf("%d", nd)
	return lik.StrToFloat(fmt.Sprintf("%." + snd + "f", val))
}

//	Собрать адрес
func MakeAddress(info lik.Seter) string {
	text := ""
	if info != nil {
		if info := info.GetString("city"); info != "" {
			text = info
		} else {
			text = "Рязань"
		}
		if info := info.GetString("street"); info != "" {
			if text != "" {
				text += ", "
			}
			text += info
		}
		if info := info.GetString("home"); info != "" {
			if text != "" {
				text += ", "
			}
			text += "д." + info
		}
		if info := info.GetString("build"); info != "" {
			text += "-" + info
		}
	}
	return text
}

//	Собрать координаты
func MakePoint(points lik.Lister) (float64,float64) {
	cx, cy := 0.0, 0.0
	if points != nil {
		cx, cy = BuildCenterList(points)
		cx = NormalizeFloat(cx, 6)
		cy = NormalizeFloat(cy, 6)
	}
	return cx, cy
}

//	Посчитать центр точек
func BuildCenterList(points lik.Lister) (float64, float64) {
	xp, yp, np := 0.0, 0.0, 0
	if points != nil {
		for n := 0; n < points.Count(); n++ {
			if pt := points.GetSet(n); pt != nil {
				xp += pt.GetFloat("x")
				yp += pt.GetFloat("y")
				np++
			}
		}
	}
	if np > 0 {
		xp /= float64(np)
		yp /= float64(np)
	}
	return xp, yp
}