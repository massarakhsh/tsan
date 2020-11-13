//	Модуль клиентского окна
package window

import (
	"github.com/massarakhsh/tsan/control"
	"github.com/massarakhsh/tsan/repo"
	"github.com/massarakhsh/lik"
	"strings"
)

//	Дескриптор клиентского окна
type ClientBox struct {
	control.DataControl       //	Основан на общем контроллере
	IsCollect	bool         //	Признак коллекции
	IsForm       bool         //	Признак формы
	ItShowGrid   DealShowGrid //	Интерфейс заполнения таблицы
	ItShowPage   DealShowPage //	Интерфейс заполнения страницы
	Class        string       //	Класс стиля
	Width        int          //	Ширина окна
	Height       int          //	Высота окна
	Title        string       //	Заголовок окна
	Parameters   lik.Seter    //	Параметры
	Events       lik.Lister   //	События
	Columns      lik.Lister   //	Колонки
	Cmds         lik.Lister   //	Элементы заголовков
	Titles       lik.Lister   //	Элементы инструментов
	Context      lik.Lister   //	Контекстное меню
	IsLockRemote bool         //	Блокировка удаленного управления
	IdSelected   string       //	Выбранный элемент
}

//	Дескриптор обработки событий
type dealBoxExecute struct {
	It	*ClientBox
}
func (it *dealBoxExecute) Run(rule *repo.DataRule, cmd string, data lik.Seter) {
	it.It.BoxExecute(rule, cmd, data)
}

//	Интерфейс и дескриптор заполнения окна
type DealShowGrid interface {
	Run(rule *repo.DataRule)
}
type dealShowGrid struct {
	It	*ClientBox
}
func (it *dealShowGrid) Run(rule *repo.DataRule) {
	it.It.CollectShowGrid(rule)
}

//	Интерфейс и дескриптор заполнения страницы
type DealShowPage interface {
	Run(rule *repo.DataRule) lik.Lister
}
type dealShowPage struct {
	It	*ClientBox
}
func (it *dealShowPage) Run(rule *repo.DataRule) lik.Lister {
	return it.It.CollectShowPage(rule)
}

//	Очистка дескриптора
func (it *ClientBox) BoxClear() {
	it.Class = ""
	it.Parameters = lik.BuildSet()
	it.Events = lik.BuildList()
	it.Columns = lik.BuildList()
	it.Cmds = lik.BuildList()
	it.Titles = lik.BuildList()
	it.Context = lik.BuildList()
}

//	Инициализация дескриптора
func (it *ClientBox) BoxInitialize(rule *repo.DataRule, frame string, id lik.IDB, mode string) {
	it.ControlInitializeZone(frame, id, mode)
	it.ItExecute = &dealBoxExecute{it}
	it.ItShowGrid = &dealShowGrid{it}
	it.ItShowPage = &dealShowPage{it}
	it.BoxClear()
}

//	Обработчик событий
func (it *ClientBox) BoxExecute(rule *repo.DataRule, cmd string, data lik.Seter) {
	if cmd == "showgrid" {
		it.boxShowGrid(rule)
	} else if cmd == "showpage" {
		it.boxShowPage(rule)
	} else if cmd == "rowselect" {
		it.boxRowSelect(rule)
	} else {
		it.ControlExecute(rule, cmd, data)
	}
}

//	Установка параметра
func (it *ClientBox) SetParameter(value interface{}, path string) {
	it.Parameters.SetItem(value, path)
}

//	Установка размеров
func (it *ClientBox) SetSize(width int, height int) {
	it.Width = width
	it.Height = height
}

//	Добавление обработчика события
func (it *ClientBox) AddEventAction(name string, proc string) {
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

//	Добавить разделитель в заголовок
func (it *ClientBox) AddCmdSep(rule *repo.DataRule, ord int) {
	it.AddCmdTitle(rule, ord, "|")
}

//	Добавить наименование в заголовок
func (it *ClientBox) AddCmdTitle(rule *repo.DataRule, ord int, text string, args ...interface{}) {
	it.AddCmdText(rule, ord, text, "", args...)
}

//	Добавить текствовую команду в заголовок
func (it *ClientBox) AddCmdText(rule *repo.DataRule, ord int, text string, cmd string, args ...interface{}) {
	item := lik.BuildSet(args...)
	if cmd != "" {
		item.SetItem("button", "type")
		item.SetItem(text, "text")
		item.SetItem(cmd, "cmd")
		item.SetItem("function_fancy_grid_cmd", "handler")
	} else if text != "" {
		item.SetItem("text", "type")
		item.SetItem(text, "text")
	}
	it.AddCmdItem(rule, ord, item)
}

//	Добавить пиктограмму - команду в заголовок
func (it *ClientBox) AddCmdImg(rule *repo.DataRule, ord int, text string, imgcls string, cmd string, args ...interface{}) {
	item := lik.BuildSet(args...)
	item.SetItem("button", "type")
	item.SetItem(text, "tip")
	item.SetItem("img"+imgcls, "imageCls")
	if match := lik.RegExParse(cmd, "(.*)/(.*)"); match != nil {
		item.SetItem(match[1], "cmd")
		item.SetItem("function_" + match[2], "handler")
	} else if cmd != "" {
		item.SetItem(cmd, "cmd")
		item.SetItem("function_fancy_grid_cmd", "handler")
	}
	item.SetItem(!strings.Contains(imgcls, "add"), "disabled")
	it.AddCmdItem(rule, ord, item)
}

//	Добавить элемент в заголовок
func (it *ClientBox) AddCmdItem(rule *repo.DataRule, ord int, item lik.Seter) {
	var pos int
	var itis bool
	for pos = it.Cmds.Count(); pos > 0; pos-- {
		if elm := it.Cmds.GetSet(pos - 1); elm != nil {
			if eord := elm.GetInt("_ord"); eord <= ord {
				itis = (eord == ord)
				break
			}
		}
	}
	if itis && item == nil {
		it.Cmds.DelItem(pos - 1)
	} else if itis {
		item.SetItem(ord, "_ord")
		it.Cmds.SetItem(item, pos - 1)
	} else if item != nil {
		item.SetItem(ord, "_ord")
		it.Cmds.InsertItem(item, pos)
	}
}

//	Отобразить объект окна клиента
func ShowCollect(rule *repo.DataRule, it *ClientBox) {
	FancyShowCode(rule, it)
}

//	Отобразить страницу окна клиента
func ShowPage(rule *repo.DataRule, it *ClientBox, rows lik.Lister) {
	FancyShowPage(rule, it, rows)
}
