// Базовые модули контроллеров
package controls

import (
	"github.com/massarakhsh/tsan/control"
	"github.com/massarakhsh/tsan/repo"
	"github.com/massarakhsh/lik/likdom"
)

//	Дескриптор пустого окна
type EmptyControl struct {
	control.DataControl
	Title string
}

//	Конструктор дескриптора
func BuildEmpty(rule *repo.DataRule, title string) *EmptyControl {
	it := &EmptyControl{}
	it.Title = title
	it.ControlInitialize( "empty", 0)
	return it
}

//	Отображение окна
func (it *EmptyControl) BuildShow(rule *repo.DataRule, sx int, sy int) likdom.Domer {
	div := likdom.BuildDivClass("fill")
	div.BuildTableClass("fill").BuildTrTdClass("fill").BuildString(it.Title)
	return div
}

