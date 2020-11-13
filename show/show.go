package show

import (
	"github.com/massarakhsh/lik"
	"github.com/massarakhsh/lik/likdom"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

//	Текстовая команда с вызовом процедуры
func LinkTextProc(cls string, text string, proc string) likdom.Domer {
	return LinkTextIdProc(cls, text,"", proc)
}

//	Текстовая команда с идентификатором
func LinkTextIdProc(cls string, text string, id string, proc string) likdom.Domer {
	div := likdom.BuildDivClassId(cls, id, "onclick", proc)
	div.BuildString(text)
	return div
}

//	Текстовая команда restFULL
func LinkTextCmd(cls string, text string, main string, zone string, cmd string) likdom.Domer {
	if cmd != "" {
		return LinkTextProc(cls, text, fmt.Sprintf("front_get('/%s/%s/%s')", main, zone, cmd))
	} else {
		return likdom.BuildString(text)
	}
}

//	Разместить таблицу FancyGrid
func BuildFancyGrid(main string, zone string) likdom.Domer {
	id := fmt.Sprintf("id_%d", 100000 + rand.Int31n(900000))
	return likdom.BuildItem("div","id", id, "remain", main, "rezone", zone, "redraw=fancy_redraw_grid")
}

//	Разместить форму FancyForm
func BuildFancyForm(main string, zone string) likdom.Domer {
	id := fmt.Sprintf("id_%d", 100000 + rand.Int31n(900000))
	return likdom.BuildItem("div","id", id, "remain", main, "rezone", zone, "redraw=fancy_redraw_form")
}

//	Построить команду вызова кода JS
func BuildRunScript(code string, del int) likdom.Domer {
	item := likdom.BuildItem("script")
	if del > 0 {
		code = fmt.Sprintf("setTimeout(function(){ %s }, %d);", code, del)
	}
	script := fmt.Sprintf("jQuery(document).ready(function() { %s });", code)
	item.BuildString(script)
	return item
}

//	Изготовление штампа времени в текстовом формате
func TimeToString(dt int) string {
	date := ""
	if dt > 0 {
		dtm := time.Unix(int64(dt), 0)
		date = dtm.Format("2006-01-02")
		date += "T" + dtm.Format("15:04:05")
		date += "+03:00"
	}
	return date
}

//	Перевод телефона в международный формат
func PhoneToString(phone string) string {
	normal := ""
	if !strings.HasPrefix(phone,"+") {
		normal += "+7"
	}
	normal += phone
	return normal
}

//	Перевод денег в полный формат
func CashToFormat(cash string) string {
	normal := ""
	for len(cash) > 0 {
		if len(normal) > 0 {
			normal = " " + normal
		}
		if match := lik.RegExParse(cash, "(.*)(\\d\\d\\d)$"); match != nil {
			normal = match[2] + normal
			cash = match[1]

		} else {
			normal = cash + normal
			cash = ""
		}
	}
	return normal
}

//	Перевод телефона в полный формат
func PhoneToFormat(phone string) string {
	digits := lik.RegExReplace(phone, "\\D", "")
	if match := lik.RegExParse(digits, "^(4912)(\\d\\d)(\\d\\d)(\\d\\d)$"); match != nil {
		phone = "+7" + "(" + match[1] + ")" + match[2] + "-" + match[3] + "-" + match[4]
	} else if match := lik.RegExParse(digits, "^(\\d\\d\\d)(\\d\\d\\d)(\\d\\d)(\\d\\d)$"); match != nil {
		phone = "+7" + "(" + match[1] + ")" + match[2] + "-" + match[3] + "-" + match[4]
	}
	return phone
}

//	Перевод ссылки в полный формат
func UrlToString(url string) string {
	normal := ""
	if strings.Contains(url, "http") {
		normal = url
	} else {
		normal = "http://rltweb.ru" + url
	}
	return normal
}

