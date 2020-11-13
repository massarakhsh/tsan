package files

import (
	"github.com/massarakhsh/tsan/control"
	"github.com/massarakhsh/tsan/fancy"
	"github.com/massarakhsh/tsan/jone"
	"github.com/massarakhsh/tsan/repo"
	"github.com/massarakhsh/tsan/show"
	"github.com/massarakhsh/lik"
	"github.com/massarakhsh/lik/likdom"
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"math/rand"
	"strings"
)

//	Дескриптор окна с файлами
type MediaControl struct {
	control.DataControl			//	Базовый дескриптор
	fancy.TableFancy			//	Объект окна
	Title string				//	Заголовок окна
	Command	string				//	Исполняемая команда
	Images	map[string][]byte	//	Список полученных файлов
}

//	Интерфейс исполнения команд
type dealMediaExecute struct {
	It	*MediaControl
}
func (it *dealMediaExecute) Run(rule *repo.DataRule, cmd string, data lik.Seter) {
	it.It.MediaExecute(rule, cmd, data)
}

//	Интерфейс списка полей
type dealMediaFieldsFill struct {
	It	*MediaControl
}
func (it *dealMediaFieldsFill) Run(rule *repo.DataRule) {
	it.It.MediaFieldsFill(rule)
}

//	Интерфейс создания таблицы
type dealMediaGridFill struct {
	It	*MediaControl
}
func (it *dealMediaGridFill) Run(rule *repo.DataRule) {
	it.It.MediaGridFill(rule)
}

//	Интерфейс создания страницы
type dealMediaPageFill struct {
	It	*MediaControl
}
func (it *dealMediaPageFill) Run(rule *repo.DataRule) lik.Lister {
	return it.It.MediaPageFill(rule)
}

//	Интерфейс создания формы
type dealMediaFormFill struct {
	It	*MediaControl
}
func (it *dealMediaFormFill) Run(rule *repo.DataRule) {
	it.It.MediaFormFill(rule)
}

//	Конструктор дескриптора
func BuildMedia(rule *repo.DataRule, main string, mode string, id lik.IDB, title string) *MediaControl {
	it := &MediaControl{ Title: title }
	it.ControlInitializeZone(main, id, mode)
	it.TableInitialize(rule, main,"offer", it.Mode)
	it.ItExecute = &dealMediaExecute{it}
	it.ItFieldsFill = &dealMediaFieldsFill{it}
	it.ItGridFill = &dealMediaGridFill{it}
	it.ItPageFill = &dealMediaPageFill{it}
	it.ItFormFill = &dealMediaFormFill{it}
	it.IsLockRemote = true
	return it
}

//	Испонения команд
func (it *MediaControl) MediaExecute(rule *repo.DataRule, cmd string, data lik.Seter) {
	if cmd == "cancel" {
		it.Command = ""
		rule.OnChangeData()
	} else if cmd == "append" {
		it.Command = "append"
		rule.OnChangeData()
	} else if cmd == "delete" {
		it.cmdDelete(rule)
		rule.OnChangeData()
	} else if cmd == "store" {
		it.storeImages(rule)
		it.Command = ""
		rule.OnChangeData()
	} else if cmd == "order" || cmd == "hide" {
		it.cmdPhotoEdit(rule, cmd, rule.Shift(), rule.Shift())
	} else if cmd == "write" {
		it.cmdWrite(rule, data)
	} else if cmd == "upload" {
		if buffers := rule.GetBuffers(); buffers != nil {
			if it.Images == nil {
				it.Images = make(map[string][]byte)
			}
			for key, val := range (buffers) {
				it.Images[key] = val
			}
		}
	} else {
		it.TableExecute(rule, cmd, data)
	}
}

//	Отображение окна
func (it *MediaControl) BuildShow(rule *repo.DataRule, sx int, sy int) likdom.Domer {
	it.SetSize(sx, sy)
	div := likdom.BuildDivClassId("roll_data", it.Zone + "_data")
	if it.Command == "append" {
		div.SetAttr("width", "100%")
		it.BuildDataAppend(rule, div)
	} else {
		div.AppendItem(show.BuildFancyGrid(it.Main, it.Zone))
	}
	return div
}

//	Окно добавления файлов
func (it *MediaControl) BuildDataAppend(rule *repo.DataRule, pater likdom.Domer) {
	pater.AppendItem(show.LinkTextCmd("cmd", "Записать эти файлы", it.Main, it.Zone,"store"))
	pater.BuildString("<br>")
	pater.AppendItem(show.LinkTextCmd("cmd", "Отменить добавление", it.Main, it.Zone,"cancel"))
	pater.BuildString("<br>")
	url := fmt.Sprintf("/front/%s/%s/upload?_sp=%d&amp;_mf=1", it.Main, it.Zone, rule.ItPage.GetPageId())
	pater.BuildItem("form", "class=dropzone", "id=mediaDropzone", "action", url)
	script := "var options = { addRemoveLinks: true };\n"
	script += "var myDropzone = new Dropzone(\"#mediaDropzone\", options);\n"
	pater.BuildItem("script").BuildString("jQuery(function(){ " + script + " });")
}

//	Заполнение полей
func (it *MediaControl) MediaFieldsFill(rule *repo.DataRule) {
	it.FancyFieldsFill(rule)
}

//	Создание таблицы
func (it *MediaControl) MediaGridFill(rule *repo.DataRule) {
	it.TableGridFill(rule)
	it.genListTitles(rule)
	it.genListColumns(rule)
	it.genListEvents(rule)
}

//	Создание страницы
func (it *MediaControl) MediaPageFill(rule *repo.DataRule) lik.Lister {
	return it.genListRows(rule)
}

//	Список кнопок
func (it *MediaControl) genListColumns(rule *repo.DataRule) {
	it.GridBuildColumns(rule, true, false)
	it.Grid.Columns.InsertItem(lik.BuildSet("type=rowdrag"), 1)
}

//	Список заголовок
func (it *MediaControl) genListTitles(rule *repo.DataRule) {
	it.Grid.Tops.AddItemSet(
		"type=text",
		"text", it.Title,
	)
	if it.Command == "" {
		it.AddCommandImg(rule,1010, "Открыть","toshow", "show")
		if it.Zone != "link" {
			it.AddCommandImg(rule, 1020, "Загрузить файлы", "append", "add")
		} else {
			it.AddCommandImg(rule, 1030, "Добавить ссылку", "toadd", "add")
		}
	} else if it.Command == "append" {
		it.Grid.Tops.AddItemSet(
			"type=button",
			"text=Записать в картотеку",
			fmt.Sprintf("handler=function_media_store(%s)", it.Zone),
		)
		it.Grid.Tops.AddItemSet(
			"type=button",
			"text=Отменить добавление",
			fmt.Sprintf("handler=function_media_cancel(%s)", it.Zone),
		)
	}
}

//	Список событий
func (it *MediaControl) genListEvents(rule *repo.DataRule) {
	it.Grid.AddEventAction("cellclick", "function_media_click")
	it.Grid.AddEventAction("dragrows", "function_fancy_grid_drag")
}

//	Список строк (файлов)
func (it *MediaControl) genListRows(rule *repo.DataRule) lik.Lister {
	rows := it.TablePageFill(rule)
	if elm := jone.TableOffer.GetElm(it.IdMain); elm != nil {
		if lpic := jone.CalculateElmList(elm,"objectid/picture"); lpic != nil {
			if it.Zone == "photo" {
				it.Grid.SetParameter(80, "cellHeight")
			}
			for np := 0; np < lpic.Count(); np++ {
				if pic := lpic.GetSet(np); pic != nil && pic.GetString("media") == it.Zone {
					sel := pic.GetString("id")
					if row := it.GridInfoRow(rule, sel, pic); row != nil {
						url := pic.GetString("url")
						if it.Zone == "photo" {
							img := fmt.Sprintf("<img class=imgedit src='%s'>", url)
							a := fmt.Sprintf("<a target=_blank href='%s'>%s</a>", url, img)
							row.SetItem(a, "photo")
						} else {
							a := fmt.Sprintf("<a target=_blank href='%s'>%s</a>", url, url)
							row.SetItem(a, "url")
						}
						rows.AddItems(row)
					}
				}
			}
		}
	}
	return rows
}

//	Заполнение формы
func (it *MediaControl) MediaFormFill(rule *repo.DataRule) {
	if it.Fun == fancy.FunAdd {
		it.Sel = ""
	}
	_,pic := it.seekPic(rule, it.Sel)
	it.Form.Items = it.FormInfoFill(rule, pic, "")
	if it.Fun == fancy.FunShow {
		it.SetTitle(rule, it.Fun, it.Title)
		it.AddTitleToolText(rule, "Изменить", "function_fancy_form_toedit")
		if true {
			it.AddTitleToolText(rule, "Удалить", "function_fancy_form_todelete")
		}
		it.AddTitleToolText(rule, "Закрыть", "function_fancy_form_cancel")
	} else if it.Fun == fancy.FunAdd {
		it.SetTitle(rule, it.Fun, "Создание внешней ссылки")
		it.AddTitleToolText(rule, "Сохранить", "function_fancy_form_write")
		it.AddTitleToolText(rule, "Отменить", "function_fancy_form_cancel")
	} else if it.Fun == fancy.FunMod || it.Fun == fancy.FunEdit {
		it.SetTitle(rule, it.Fun, "Редактирование")
		it.AddTitleToolText(rule, "Сохранить", "function_fancy_form_write")
		it.AddTitleToolText(rule, "Отменить", "function_fancy_form_cancel")
	} else if it.Fun == fancy.FunDel {
		it.SetTitle(rule, it.Fun, "Удаление")
		it.AddTitleToolText(rule, "Действительно удалить?", "function_fancy_real_delete")
		it.AddTitleToolText(rule, "Отменить", "function_fancy_form_cancel")
	}
}

//	Запись изменений
func (it *MediaControl) cmdWrite(rule *repo.DataRule, data lik.Seter) {
	if elm := jone.TableOffer.GetElm(it.IdMain); elm != nil {
		obj := jone.TableObject.GetElm(elm.GetIDB("objectid"))
		list := jone.CalculateElmList(elm,"objectid/picture")
		if list == nil {
			list = lik.BuildList()
			jone.SetElmValue(elm, list, "objectid/picture")
		}
		if list != nil && data != nil {
			url := lik.StringFromXS(data.GetString("s_url"))
			var pic lik.Seter
			if it.Fun == fancy.FunAdd && url != "" {
				pic = list.AddItemSet("id", fmt.Sprintf("%d", 100000000 + rand.Intn(900000000)))
			} else if it.Fun == fancy.FunMod || it.Fun == fancy.FunEdit {
				_,pic = it.seekPic(rule, it.Sel)
			}
			if pic != nil && data != nil {
				pic.SetItem(it.Zone,"media")
				it.UpdateInfoData(rule, pic, data)
			}
		}
		elm.OnModify()
		if obj != nil {
			obj.OnModify()
		}
		rule.OnChangeData()
	}
}

//	Удаление файла
func (it *MediaControl) cmdDelete(rule *repo.DataRule) {
	it.cmdPhotoEdit(rule, "delete", it.Sel, "1")
}

//	Запоминание файлов
func (it *MediaControl) storeImages(rule *repo.DataRule) {
	elm := jone.TableOffer.GetElm(it.IdMain)
	if elm != nil {
		if it.Images != nil {
			obj := jone.TableObject.GetElm(elm.GetIDB("objectid"))
			var list lik.Lister
			if item := jone.CalculateElm(elm,"objectid/picture"); item != nil {
				list = item.ToList()
			} else {
				list = lik.BuildList()
				jone.SetElmValue(elm, list,"objectid/picture")
			}
			for key, img := range (it.Images) {
				ext := ""
				if match := lik.RegExParse(key,"\\.(\\w*)$"); match != nil {
					ext = strings.ToLower(match[1])
				}
				filepath := repo.WriteFile("obj", int(it.IdMain), ext, img)
				pic := lik.BuildSet("media", it.Zone, "album=yy", "promot=yy")
				pic.SetItem("/"+filepath,"url")
				if match := lik.RegExParse(filepath,"(\\d+)\\."); match != nil {
					pic.SetItem(match[1],"id")
				}
				if match := lik.RegExParse(filepath, "^(.+)\\.jpg$"); match != nil {
					patht := match[1] + "t.jpg"
					show.MakeScaleJpg(filepath, 0, 64, patht)
				}
				list.AddItems(pic)
			}
			elm.OnModify()
			if obj != nil {
				obj.OnModify()
			}
			it.Images = nil
		}
	}
}

//	Редактирование элемента
func (it *MediaControl) cmdPhotoEdit(rule *repo.DataRule, cmd string, sid string, val string) {
	ival,ok := lik.StrToIntIf(val)
	if  sid != "" && ok {
		if elm := jone.TableOffer.GetElm(it.IdMain); elm != nil {
			if lpic := jone.CalculateElmList(elm,"objectid/picture"); lpic != nil {
				jone.SetElmValue(elm, lpic, "objectid/picture")
				for np := 0; np < lpic.Count(); np++ {
					pic := lpic.GetSet(np)
					if pic.GetString("id") == sid {
						if cmd == "video" || cmd == "model" || cmd == "scheme" {
							if ival > 0 {
								pic.SetItem(cmd, "what")
							} else {
								pic.DelItem("what")
							}
							elm.OnModify()
						} else if cmd == "hide" {
							pic.SetItem(ival, "hide")
							elm.OnModify()
						} else if cmd == "order" {
							lpic.DelItem(np)
							lpic.InsertItem(pic, ival)
							elm.OnModify()
						} else if cmd == "delete" {
							lpic.DelItem(np)
							it.Sel = ""
							elm.OnModify()
						}
						break
					}
				}
			}
			rule.OnChangeData()
		}
	}
}

//	Поиск элемента
func (it *MediaControl) seekPic(rule *repo.DataRule, sid string) (int,lik.Seter) {
	if elm := jone.TableOffer.GetElm(it.IdMain); elm != nil {
		if lpic := jone.CalculateElmList(elm,"objectid/picture"); lpic != nil {
			for np := 0; np < lpic.Count(); np++ {
				if pic := lpic.GetSet(np); pic != nil && pic.GetString("media") == it.Zone {
					if pic.GetString("id") == sid {
							return np, pic
						}
				}
			}
		}
	}
	return -1,nil
}

