package jone

import (
	"bitbucket.org/shaman/lik"
	"strings"
)

//	Сводный словарь
var Trans map[string]string

//	Трансляция кода и интерфейса
func SystemTranslate(key string, data lik.Itemer) string {
	text := ""
	if data != nil {
		text = SystemStringTranslate(key, data.ToString())
	}
	return text
}

//	Трансляция кода и строки
func SystemStringTranslate(key string, part string) string {
	text := part
	if part == "yy" {
		text = "да"
	} else if part == "nn" {
		text = "нет"
	} else if list := strings.Split(part, ","); len(list) > 1 {
		text = ""
		for _, txt := range list {
			txt := strings.TrimSpace(txt)
			if value, ok := Trans[key+"+"+txt]; ok {
				txt = value
			}
			if txt != "" {
				if text != "" {
					text += ", "
				}
				text += txt
			}
		}
	} else if value, ok := Trans[key+"+"+part]; ok {
		text = value
	}
	return text
}