package member

import (
	"github.com/massarakhsh/tsan/control"
	"github.com/massarakhsh/tsan/jone"
	"github.com/massarakhsh/tsan/repo"
	"github.com/massarakhsh/tsan/show"
	"github.com/massarakhsh/lik"
	"github.com/massarakhsh/lik/likdom"
	"fmt"
	"strings"
)

//	Дескриптор управления записью сотрудника
type ManageControl struct {
	control.DataControl
	Command	string				//	Исполняемая команда
	ImageCode	[]byte
	ImageName	string
}

//	Интерфейс команд
type dealManageExecute struct {
	It	*ManageControl
}
func (it *dealManageExecute) Run(rule *repo.DataRule, cmd string, data lik.Seter) {
	it.It.cmdExecute(rule, cmd, data)
}

//	Конструктор дескриптора карточки
func BuildManage(rule *repo.DataRule, frame string, id lik.IDB) *ManageControl {
	it := &ManageControl{}
	it.ControlInitializeZone(frame, id, "manage")
	it.ItExecute = &dealManageExecute{it}
	return it
}

//	Обработка команд
func (it *ManageControl) cmdExecute(rule *repo.DataRule, cmd string, data lik.Seter) {
	if cmd == "manage" {
		it.cmdExecute(rule, rule.Shift(), data)
	} else if cmd == "append" {
		it.cmdAppend(rule)
	} else if cmd == "cancel" {
		it.cmdCancel(rule)
	} else if cmd == "upload" {
		it.imageUpload(rule)
	} else if cmd == "delete" {
		it.cmdDelete(rule)
	} else if cmd == "store" {
		it.imageStore(rule)
	} else {
		it.ControlExecute(rule, cmd, data)
	}
}

//	Отображение карточки
func (it *ManageControl) BuildShow(rule *repo.DataRule, sx int, sy int) likdom.Domer {
	//elm := repo.TableMember.GetElm(it.GetId())
	tbl := likdom.BuildTableClass("manage fill")
	if row := tbl.BuildTr(); row != nil {
		row.BuildTdClass("center", "width=50%").BuildString("Карточка")
		row.BuildTdClass("center", "width=50%").AppendItem(it.showPhoto(rule, sx / 2, sy))
	}
	return tbl
}

//	Отображение карточки
func (it *ManageControl) showPhoto(rule *repo.DataRule, sx int, sy int) likdom.Domer {
	if it.Command == "append" {
		return it.showPhotoAppend(rule, sx, sy)
	} else {
		return it.showPhotoAvatar(rule, sx, sy)
	}
}

func (it *ManageControl) showPhotoAppend(rule *repo.DataRule, sx int, sy int) likdom.Domer {
	div := likdom.BuildDiv()
	div.AppendItem(it.showPhotoTools(rule))
	url := fmt.Sprintf("/front/%s/%s/upload?_sp=%d&amp;_mf=1", it.Frame, it.Mode, rule.ItPage.GetPageId())
	div.BuildItem("form", "class=dropzone", "id=mediaDropzone", "action", url)
	script := "var options = { maxFiles: true };\n"
	script += "var myDropzone = new Dropzone(\"#mediaDropzone\", options);\n"
	div.BuildItem("script").BuildString("jQuery(function(){ " + script + " });")
	return div
}

func (it *ManageControl) showPhotoAvatar(rule *repo.DataRule, sx int, sy int) likdom.Domer {
	div := likdom.BuildDiv()
	div.AppendItem(it.showPhotoTools(rule))
	if url := jone.CalculatePartIdString("member", it.IdMain, "photo"); url != "" {
		a := div.BuildItem("a", "target=_blank", "href", url)
		img := a.BuildUnpairItem("img", "class=imgphoto", "src", url)
		img.SetAttr(fmt.Sprintf("height=%d", sy - 20))
	}
	return div
}

//	Отображение инструментов аватарки
func (it *ManageControl)showPhotoTools(rule *repo.DataRule) likdom.Domer {
	tbl := likdom.BuildItemClass("table","mapcmd", "id=maptools")
	row := tbl.BuildTr()
	if it.Command == "append" {
		row.BuildTdClass("mapcmd").AppendItem(show.LinkTextProc("", "&nbsp;Записать&nbsp;", "ava_control('store')"))
		row.BuildTdClass("mapcmd").AppendItem(show.LinkTextProc("", "&nbsp;Отменить&nbsp;", "ava_control('cancel')"))
	} else if url := jone.CalculatePartIdString("member", it.IdMain, "photo"); url != "" {
		row.BuildTdClass("mapcmd").AppendItem(show.LinkTextProc("", "&nbsp;Заменить&nbsp;", "ava_control('append')"))
		row.BuildTdClass("mapcmd").AppendItem(show.LinkTextProc("", "&nbsp;Удалить&nbsp;", "ava_control('delete')"))
	} else {
		row.BuildTdClass("mapcmd").AppendItem(show.LinkTextProc("", "&nbsp;Загрузить&nbsp;", "ava_control('append')"))
	}
	row.BuildTd("width=100%")
	return tbl
}

//	Добавление аватарки
func (it *ManageControl) cmdAppend(rule *repo.DataRule) {
	it.Command = "append"
	rule.OnChangeData()
}

//	Отмена добавления
func (it *ManageControl) cmdCancel(rule *repo.DataRule) {
	it.Command = ""
	it.ImageCode = nil
	rule.OnChangeData()
}

//	Стирание аватарки
func (it *ManageControl) cmdDelete(rule *repo.DataRule) {
	if elm := jone.TableMember.GetElm(it.IdMain); elm != nil {
		if url := elm.GetString("photo"); url != "" {
			elm.SetValue(nil, "photo")
			elm.OnModify()
		}
		rule.OnChangeData()
	}
}

//	Загрузка аватарки
func (it *ManageControl) imageUpload(rule *repo.DataRule) {
	if buffers := rule.GetBuffers(); buffers != nil {
		for key, val := range (buffers) {
			it.ImageName = key
			it.ImageCode = val
		}
	}
}

//	Запоминание аватарки
func (it *ManageControl) imageStore(rule *repo.DataRule) {
	if elm := jone.TableMember.GetElm(it.IdMain); elm != nil && it.ImageCode != nil {
		ext := ""
		if match := lik.RegExParse(it.ImageName,"\\.(\\w*)$"); match != nil {
			ext = strings.ToLower(match[1])
		}
		filepath := repo.WriteFile("member", int(it.IdMain), ext, it.ImageCode)
		elm.SetValue("/"+filepath,"photo")
		elm.OnModify()
		it.Command = ""
		it.ImageCode = nil
		rule.OnChangeData()
	}
}

