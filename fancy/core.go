package fancy

import (
	"bitbucket.org/shaman/lik"
	"strings"
)

//	Дескриптор объекта Fancy
type FancyCore struct {
	Class      string
	Width      int
	Height     int
	Parameters lik.Seter
	Events     lik.Lister
}

//	Очистка объекта
func (it *FancyCore) FancyClear() {
	it.Class = ""
	//it.Width= 0
	//it.Height = 0
	it.Parameters = lik.BuildSet()
	it.Events = lik.BuildList()
}

//	Установка значения параметра
func (it *FancyCore) SetParameter(value interface{}, path string) {
	it.Parameters.SetItem(value, path)
}

//	Кстановка размеров объекта
func (it *FancyCore) SetSize(width int, height int) {
	it.Width = width
	it.Height = height
}

//	Добавить обработку события
func (it *FancyCore) AddEventAction(name string, proc string) {
	for ne := 0; ne < it.Events.Count(); ne++ {
		if event := it.Events.GetSet(ne); event != nil {
			if old := event.GetString(name); old != "" {
				if old != proc {
					event.SetItem(proc, name)
				}
				return
			}
		}
	}
	it.Events.AddItemSet(name, proc)
}

//	Создать параметры ядра
func (it *FancyCore) FillCore(code lik.Seter) {
	if !code.IsItem("width") {
		if it.Width > 0 {
			code.SetItem(it.Width, "width")
		} else if it.Width < 0 {
			code.SetItem("100%", "width")
		} else {
			code.SetItem("fit", "width")
		}
	}
	if !code.IsItem("height") {
		if it.Height > 0 {
			code.SetItem(it.Height, "height")
		} else if it.Height < 0 {
			code.SetItem("100%", "height")
		} else {
			code.SetItem("fit", "height")
		}
	}
	if !code.IsItem("cls") && it.Class != "" {
		code.SetItem(it.Class, "cls")
	}
	if !code.IsItem("defaults/type") {
		code.SetItem("string", "defaults/type")
	}
	if !code.IsItem("events") {
		code.SetItem(it.Events, "events")
	}
}

//	Разделить цифры
func StringateNumber(data string) string {
	text := strings.ReplaceAll(data, " ", "")
	ceil := text
	frac := ""
	mode := ""
	if match := lik.RegExParse(text,"^(\\d*)(,|\\.|)(\\d*)(.*)"); match != nil {
		ceil = match[1]
		frac = match[3]
		mode = match[4]
	}
	for len(ceil) < 10 {
		ceil = "0" + ceil
	}
	for len(frac) < 3 {
		frac = frac + "0"
	}
	return ceil + "." + frac + mode
}

