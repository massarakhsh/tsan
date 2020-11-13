// Контроллеры окон визуализации.
//
// Все контроллеры поддерживают интерфейс control.Controller
package control

import (
	"bitbucket.org/961961/tsan/repo"
	"bitbucket.org/shaman/lik"
	"bitbucket.org/shaman/lik/likdom"
	"fmt"
)

const BD = 8	//	Настраиваемый размер межоконного промежутка

//	Дескриптор окна контроллера
type DataControl struct {
	repo.DataControl			//	Ядро контроллера
	Controls  []Controller		//	Коллекция компонентов
	ItMarshal ControlMarshal	//	Интерфейс маршрутизации
	ItExecute ControlExecute	//	Интерфейс команд
	ItUpdate ControlUpdate		//	Интерфейс обновления
}

//	Интерфейс конроллера
type Controller interface {
	repo.Controller							//	Интерфейс ядра
	RunMarshal(rule *repo.DataRule)			//	Маршалинг
	RunExecute(rule *repo.DataRule, cmd string, data lik.Seter)		//	Команды
	RunUpdate(rule *repo.DataRule)				//	Обновления
	AddControl(mark string, control Controller)	//	Добавление компоненты
	FindControl(zone string) Controller			//	Найти компоненту
	BuildShow(rule *repo.DataRule, sx int, sy int) likdom.Domer	//	Отображение окна
}

//	Интерфейс маршалинга
type ControlMarshal interface {
	Run(rule *repo.DataRule)
}
type controlMarshal struct {
	It	*DataControl
}
func (it *controlMarshal) Run(rule *repo.DataRule) {
	it.It.ControlMarshal(rule)
}

//	Интерфейс команд
type ControlExecute interface {
	Run(rule *repo.DataRule, cmd string, data lik.Seter)
}
type controlExecute struct {
	It	*DataControl
}
func (it *controlExecute) Run(rule *repo.DataRule, cmd string, data lik.Seter) {
	it.It.ControlExecute(rule, cmd, data)
}

//	Интерфейс управления
type ControlUpdate interface {
	Run(rule *repo.DataRule)
}
type controlUpdate struct {
	It	*DataControl
}
func (it *controlUpdate) Run(rule *repo.DataRule) {
	//it.It.ControlUpdate()
}

//	Инициализация контроллера
func (it *DataControl) ControlInitialize(frame string, id lik.IDB) {
	it.ControlInitializeZone(frame, id, frame)
}

//	Инициализация контроллера в зоне
func (it *DataControl) ControlInitializeZone(frame string, id lik.IDB, mode string) {
	it.RepoInitialize(frame, id, mode)
	it.ItMarshal = &controlMarshal{it}
	it.ItExecute = &controlExecute{it}
	it.ItUpdate = &controlUpdate{it}
}

//	Маршализация
func (it *DataControl) RunMarshal(rule *repo.DataRule) {
	if it.ItMarshal != nil {
		it.ItMarshal.Run(rule)
	}
}

//	Команды
func (it *DataControl) RunExecute(rule *repo.DataRule, cmd string, data lik.Seter) {
	if it.ItExecute != nil {
		it.ItExecute.Run(rule, cmd, data)
	}
}

//	Обновление
func (it *DataControl) RunUpdate(rule *repo.DataRule) {
	if it.ItUpdate != nil {
		it.ItUpdate.Run(rule)
	}
}

//	Маршализация
func (it *DataControl) ControlMarshal(rule *repo.DataRule) {
	for _,control := range it.Controls {
		if control != nil {
			control.(Controller).RunMarshal(rule)
		}
	}
}

//	Команды
func (it *DataControl) ControlExecute(rule *repo.DataRule, cmd string, data lik.Seter) {
	if control := it.FindControl(cmd); control != nil {
		control.(Controller).RunExecute(rule, rule.Shift(), data)
	}
}

//	Сборка секции
func (it *DataControl) BuildSection(row likdom.Domer, sx int, sy int) (likdom.Domer,int,int) {
	td := row.BuildTdClass("section", MiniMax(sx, sy)...)
	return td, sx, sy
}

//	Определение границ окна
func MiniMax(sx int, sy int) []string {
	parms := []string{}
	if sx > 0 {
		parms = append(parms,
			fmt.Sprintf("width=%d", sx),
			fmt.Sprintf("max-width=%d", sx),
		)
	}
	if sy > 0 {
		parms = append(parms,
			fmt.Sprintf("height=%d", sy),
			fmt.Sprintf("max-height=%d", sy),
		)
	}
	return parms
}

//	Добавить компонент
func (it *DataControl) AddControl(mark string, control Controller) {
	if control != nil {
		control.SetMark(mark)
		it.Controls = append(it.Controls, control)
	}
}

//	Найти компонент
func (it *DataControl) FindControl(zone string) Controller {
	for _, control := range it.Controls {
		if control != nil && control.GetMode() == zone {
			return control
		}
	}
	for _, control := range it.Controls {
		if control != nil && control.GetMark() == zone {
			return control
		}
	}
	return nil
}

