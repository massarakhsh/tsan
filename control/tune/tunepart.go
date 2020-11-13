package tune

import (
	"github.com/massarakhsh/tsan/fancy"
	"github.com/massarakhsh/tsan/jone"
	"github.com/massarakhsh/tsan/repo"
	"github.com/massarakhsh/lik"
	"github.com/massarakhsh/lik/likbase"
	"fmt"
	"regexp"
)

//	Интерфейс списка полей
type dealTuneFillFields struct {
	It	*TuneControl
}
func (it *dealTuneFillFields) Run(rule *repo.DataRule) {
	it.It.TuneFillFields(rule)
}

//	Интерфейс заполнения таблицы
type dealTuneGridFill struct {
	It	*TuneControl
}
func (it *dealTuneGridFill) Run(rule *repo.DataRule) {
	it.It.TuneGridFill(rule)
}

//	Интерфейс заполнения страницы
type dealTunePageFill struct {
	It	*TuneControl
}
func (it *dealTunePageFill) Run(rule *repo.DataRule) lik.Lister {
	return it.It.TunePageFill(rule)
}

//	Интерфейс заполнения формы
type dealTuneFormFill struct {
	It	*TuneControl
}
func (it *dealTuneFormFill) Run(rule *repo.DataRule) {
	it.It.TuneFormFill(rule)
}

//	Инициализация разделов
func (it *TuneControl) PartInitialize(rule *repo.DataRule) {
	it.TableInitialize(rule, "tune", "system", "part")
	it.ItGridFill = &dealTuneGridFill{it}
	it.ItPageFill = &dealTunePageFill{it}
	it.ItFormFill = &dealTuneFormFill{it }
	it.ItFieldsFill = &dealTuneFillFields{it}
	it.IsLockRemote = true
}

//	Команды разделов
func (it *TuneControl) PartExecute(rule *repo.DataRule, cmd string, data lik.Seter) {
	if match := regexp.MustCompile("^tag_(\\d+)$").FindStringSubmatch(cmd); match != nil {
		tag := lik.StrToInt(match[1])
		parm := rule.Shift()
		val := rule.Shift()
		it.PartCmdTrigger(rule, parm, tag, lik.StrToInt(val)>0)
	} else if cmd == "order" {
		sid := rule.Shift()
		pos := lik.StrToInt(rule.Shift())
		it.partCmdOrder(rule, sid, pos)
	} else if cmd == "rowselect" {
		//it.partRowSelecting(rule, parm, true)
		it.TableExecute(rule, cmd, data)
	} else if cmd == "rowdeselect" {
		//it.partRowSelecting(rule, parm, false)
		it.TableExecute(rule, cmd, data)
	} else if cmd == "delete" {
		it.PartCmdDelete(rule)
		rule.OnChangeData()
	} else if cmd == "write" {
		it.PartCmdWrite(rule, data)
		rule.OnChangeData()
	} else if cmd == "mark" {
		it.partCmdClipMark(rule, rule.Shift(), rule.Shift())
	} else if cmd == "clipcopy" {
		it.partCmdClipCopy(rule)
		rule.OnChangeData()
	} else if cmd == "clipcut" {
		it.partCmdClipCut(rule)
		rule.OnChangeData()
	} else if cmd == "clippaste" {
		it.partCmdClipPaste(rule)
		rule.OnChangeData()
	} else if cmd == "clipclear" {
		it.partCmdClipClear(rule)
		rule.OnChangeData()
	} else {
		it.TableExecute(rule, cmd, data)
	}
}

//	Заполнение полей
func (it *TuneControl) TuneFillFields(rule *repo.DataRule) {
	it.ListFields = []lik.Seter{}
	it.itTags = false
	tags := fmt.Sprintf("tags=%d", jone.TagGrid|jone.TagForm|jone.TagEdit)
	if it.itWhat == "gen" {
		it.ListFields = append(it.ListFields,
			lik.BuildSet("name=Имя", "part=name", "index=name", "format=s", "width=150", tags),
			lik.BuildSet("name=Ключ", "part=key", "index=key", "format=s", "width=200", tags),
		)
	} else if it.itWhat == "ent" {
		it.ListFields = append(it.ListFields,
			lik.BuildSet("name=Имя", "part=name", "index=name", "format=s", "width=150", tags),
			lik.BuildSet("name=Ключ", "part=part", "index=part", "format=s", "width=200", tags),
		)
		if it.itEnt == "export" {
			it.ListFields = append(it.ListFields,
				lik.BuildSet("name=Сигнатура", "part=signatura", "index=signatura", "format=s", "width=250", tags),
			)
		} else {
			it.itTags = true
			it.ListFields = append(it.ListFields,
				lik.BuildSet("name=Шир", "part=width", "index=width", "format=s", "width=50", tags),
				lik.BuildSet("name=Фмт", "part=format", "index=format", "format=c", "width=32", tags),
			)
		}
	}
}

//	Заполнение таблицы
func (it *TuneControl) TuneGridFill(rule *repo.DataRule) {
	it.TableGridFill(rule)
	it.Grid.SetParameter(false,"defaults/sortable")
	it.AddCommandItem(rule, 80, lik.BuildSet("type=text", "text", "Поля"))
	it.AddCommandItem(rule, 400, lik.BuildSet("type=search", "width=150", "emptyText=Найти"))
	it.AddCommandImg(rule, 910, "Открыть", "toshow", "show")
	it.Grid.AddEventAction("cellclick", "function_fancy_grid_mark")
	it.GridBuildColumns(rule, true, true)
	it.AddCommandImg(rule, 950, "Создать", "toadd", "add")
	it.AddCommandImg(rule, 960, "Удалить", "todel", "del")
	if it.itWhat == "ent" {
		it.Grid.Columns.InsertItem(lik.BuildSet("type=rowdrag"), 1)
		it.Grid.AddEventAction("dragrows", "function_fancy_grid_drag")
		//it.Grid.SetParameter("rows","selModel")
		//it.Grid.Columns.InsertItem(lik.BuildSet("type=select", "locked=true"), 1)
		it.Grid.Columns.InsertItem(lik.BuildSet("type=checkbox", "width=44", "index=mark",
			"title=Выбор", "cellTip=Выбор для операции", "editable=true"), 1)
		nsel := len(rule.ItPage.Tune.List)
		for nclp,clp := range([]string{ "copy", "cut", "paste", "clear" }) {
			text := ""
			disabled := "true"
			if clp == "copy" {
				text = "Копировать"
				disabled = "false"
			} else if clp == "cut" {
				text = "Вырезать"
				disabled = "false"
			} else if clp == "paste" {
				text = "Вставить"
				if nsel > 0 {
					text += fmt.Sprintf(" (%d)", nsel)
					disabled = "false"
				}
			} else if clp == "clear" {
				text = "Сброс"
				disabled = "false"
			}
			it.AddCommandItem(rule, 1500 + nclp,
				lik.BuildSet("type=button", "text="+text, "handler=function_fancy_clip_" + clp, "disabled="+disabled))
		}
		//it.Grid.AddEventAction("deselectrow", "function_fancy_grid_rowdeselect")
	}
	for nc := 1; nc < it.Grid.Columns.Count(); nc++ {
		it.Grid.Columns.GetSet(nc).SetItem("left","cellAlign")
	}
	if it.itTags {
		for _, opt := range ListOpt {
			if opt.Tag != 0 {
				title := opt.Title
				index := fmt.Sprintf("tag_%d", opt.Tag)
				it.Grid.Columns.AddItemSet("type=checkbox", "width=44", "index",
					index, "title", title, "cellTip", opt.Tip, "editable=true")
			} else {
				it.Grid.Columns.AddItemSet("width=4", "headerCls=separator")
			}
		}
	}
}

//	Заполнение страницы
func (it *TuneControl) TunePageFill(rule *repo.DataRule) lik.Lister {
	if it.itWhat == "gen" {
		return it.TuneRowsFillGen(rule)
	} else if it.itWhat == "ent" {
		return it.TuneRowsFillEnt(rule)
	}
	return it.TablePageFill(rule)
}

//	Заполнение строк генератора
func (it *TuneControl) TuneRowsFillGen(rule *repo.DataRule) lik.Lister {
	rows := it.TablePageFill(rule)
	if list := repo.SystemSortEnts(it.itGen); list != nil {
		for _, elm := range list {
			row := it.GridInfoRow(rule, likbase.IDBToStr(elm.Id), elm.Info)
			tags := elm.GetInt("tags")
			for _, opt := range ListOpt {
				if opt.Tag != 0 {
					index := fmt.Sprintf("tag_%d", opt.Tag)
					row.SetItem((tags&opt.Tag) > 0, index)
				}
			}
			rows.AddItems(row)
		}
	}
	return rows
}

//	Заполнение строк сущности
func (it *TuneControl) TuneRowsFillEnt(rule *repo.DataRule) lik.Lister {
	rows := it.TablePageFill(rule)
	it.listIdParts = []string{}
	if ent := repo.SystemFindGenEnt(it.itGen, it.itEnt); ent != nil {
		if content := ent.GetContent(); content != nil && len(content) > 0 {
			for ne, elm := range content {
				if imap := elm.ToSet(); imap != nil {
					it.listIdParts = append(it.listIdParts, lik.IntToStr(1+ne))
					row := it.GridInfoRow(rule, lik.IntToStr(1+ne), elm)
					//row.SetItem(true, "selected")
					tags := imap.GetInt("tags")
					for _, opt := range ListOpt {
						if opt.Tag != 0 {
							index := fmt.Sprintf("tag_%d", opt.Tag)
							row.SetItem((tags&opt.Tag) > 0, index)
						}
					}
					row.SetItem(it.partCmdClipTest(rule, imap.GetString("part")), "mark")
					rows.AddItems(row)
				}
			}
		}
	}
	return rows
}

//	Заполнение формы
func (it *TuneControl) TuneFormFill(rule *repo.DataRule) {
	title := "Параметр"
	var info lik.Seter
	if it.itWhat == "gen" {
		if it.itGen == repo.KeyTable {
			title = "Таблицы"
		} else if it.itGen == repo.KeyStruct {
			title = "Структуры"
		} else if it.itGen == repo.KeyDiction {
			title = "Словари"
		} else {
			title = "Раздел"
		}
		if gen := repo.SystemFindGen(it.itGen); gen != nil {
			if it.Fun != fancy.FunAdd {
				if elm := jone.GetElm(repo.KeyParam, likbase.StrToIDB(it.Sel)); elm != nil {
					info = elm.Info
				}
			}
		}
	} else if it.itWhat == "ent" {
		if ent := repo.SystemFindGenEnt(it.itGen, it.itEnt); ent != nil {
			if it.Fun != fancy.FunAdd {
				if poso := it.partSeekSid(it.Sel); poso >= 0 {
					info = ent.FindPartPos(poso)
				}
			}
		}
	}
	it.Form.Items = it.FormInfoFill(rule, info, "")
	if it.Fun == fancy.FunShow {
		it.AddTitleToolText(rule, "Изменить", "function_fancy_form_toedit")
		it.AddTitleToolText(rule, "Удалить", "function_fancy_form_todelete")
		it.AddTitleToolText(rule, "Закрыть", "function_fancy_form_cancel")
	} else if it.Fun == fancy.FunAdd {
		title += ". Создание"
		it.AddTitleToolText(rule, "Сохранить", "function_fancy_form_write")
		it.AddTitleToolText(rule, "Отменить", "function_fancy_form_cancel")
	} else if it.Fun == fancy.FunMod || it.Fun == fancy.FunEdit {
		title += ". Редактирование"
		it.AddTitleToolText(rule, "Сохранить", "function_fancy_form_write")
		it.AddTitleToolText(rule, "Отменить", "function_fancy_form_cancel")
	} else if it.Fun == fancy.FunDel {
		title += ". Удаление"
		it.AddTitleToolText(rule, "Действительно удалить?", "function_fancy_real_delete")
		it.AddTitleToolText(rule, "Отменить", "function_fancy_form_cancel")
	}
	it.SetTitle(rule, it.Fun, title)
	if it.itWhat == "ent" && it.itEnt == "export" && it.Fun == fancy.FunShow {
		url := repo.GetExternUrl(info.GetString("part"))
		it.Form.Items.AddItemSet("type=html", "value", fmt.Sprintf("Ссылка на файл выгрузки: <b><u>%s</u></b>", url))
	}
}

//	Команды триггера
func (it *TuneControl) PartCmdTrigger(rule *repo.DataRule, sid string, tag int, on bool) {
	if it.itWhat == "gen" {
		if elm := jone.GetElm(repo.KeyParam, likbase.StrToIDB(sid)); elm != nil {
			tags := elm.GetInt("tags")
			if on {
				tags |= tag
			} else {
				tags &= (0xffffff ^ tag)
			}
		}
	} else if it.itWhat == "ent" {
		if ent := repo.SystemFindGenEnt(it.itGen, it.itEnt); ent != nil {
			if col := ent.FindPartPos(it.partSeekSid(sid)); col != nil {
				tags := col.GetInt("tags")
				if on {
					tags |= tag
				} else {
					tags &= (0xffffff ^ tag)
				}
				col.SetItem(tags, "tags")
			}
			ent.SaveToBase()
		}
	}
}

//	Команда изменения порядка
func (it *TuneControl) partCmdOrder(rule *repo.DataRule, sid string, pos int) {
	if ent := repo.SystemFindGenEnt(it.itGen, it.itEnt); ent != nil {
		if content := jone.CalculateElmList(ent.It, "content"); content != nil {
			if poso := it.partSeekSid(sid); poso >= 0 {
				if col := ent.FindPartPos(poso); col != nil {
					listo := []string{}
					listo = append(listo, it.listIdParts[:poso]...)
					listo = append(listo, it.listIdParts[poso+1:]...)
					listn := []string{}
					for p := 0; p < len(listo); p++ {
						if p == pos {
							listn = append(listn, sid)
						}
						listn = append(listn, listo[p])
					}
					it.listIdParts = listn
					content.DelItem(poso)
					content.InsertItem(col, pos)
					ent.SaveToBase()
				}
			}
		}
	}
}

//	Позиционирование объекта
func (it *TuneControl) partSeekSid(sid string) int {
	pos := -1
	for poso := 0; poso < len(it.listIdParts); poso++ {
		if sid == it.listIdParts[poso] {
			pos = poso
			break
		}
	}
	return pos
}

//	Запись изменений
func (it *TuneControl) PartCmdWrite(rule *repo.DataRule, data lik.Seter) {
	if it.itWhat == "gen" {
		it.PartCmdWriteGen(rule, data)
	} else if it.itWhat == "ent" {
		it.PartCmdWriteEnt(rule, data)
	}
	rule.OnChangeData()
}

//	Запись изменений генератора
func (it *TuneControl) PartCmdWriteGen(rule *repo.DataRule, data lik.Seter) {
	var elm *likbase.ItElm
	if it.Fun == fancy.FunAdd {
		elm = jone.TableParam.CreateElm()
		elm.SetValue(it.itGen, "gen")
	} else {
		elm = jone.GetElm(repo.KeyParam, likbase.StrToIDB(it.Sel))
	}
	if elm != nil {
		if it.UpdateInfoData(rule, elm.Info, data) != nil {
			elm.OnModifyWait()
			repo.SystemInitialize()
		}
	}
}

//	Запись изменений сущности
func (it *TuneControl) PartCmdWriteEnt(rule *repo.DataRule, data lik.Seter) {
	if ent := repo.SystemFindGenEnt(it.itGen, it.itEnt); ent != nil {
		content := jone.CalculateElmList(ent.It, "content")
		if content == nil {
			content = lik.BuildList()
			ent.It.SetValue(content, "content")
		}
		var col lik.Seter
		if it.Fun == fancy.FunAdd {
			col = lik.BuildSet()
			content.AddItems(col)
		} else if poso := it.partSeekSid(it.Sel); poso >= 0 {
			col = ent.FindPartPos(poso)
		}
		if col != nil {
			if it.UpdateInfoData(rule, col, data) != nil {
				ent.SaveToBase()
			}
		}
	}
}

//	Команда удаления объекта
func (it *TuneControl) PartCmdDelete(rule *repo.DataRule) {
	if it.itWhat == "gen" {
		it.PartCmdDeleteGen(rule)
	} else if it.itWhat == "ent" {
		it.PartCmdDeleteEnt(rule)
	}
	rule.OnChangeData()
}

//	Команда удаления генератора
func (it *TuneControl) PartCmdDeleteGen(rule *repo.DataRule) {
	jone.TableParam.DeleteElm(likbase.StrToIDB(it.Sel))
	repo.SystemInitialize()
}

//	Команда удаления сущности
func (it *TuneControl) PartCmdDeleteEnt(rule *repo.DataRule) {
	if ent := repo.SystemFindGenEnt(it.itGen, it.itEnt); ent != nil {
		if content := jone.CalculateElmList(ent.It, "content"); content != nil {
			content.DelItem(lik.StrToInt(it.Sel)-1)
			ent.SaveToBase()
		}
	}
}

//	Работа со списком отмеченных
func (it *TuneControl) partCmdClipMark(rule *repo.DataRule, sid string, val string) {
	if ent := repo.SystemFindGenEnt(it.itGen, it.itEnt); ent != nil {
		if list := ent.It.GetList("content"); list != nil {
			if num := lik.StrToInt(sid); num > 0 {
				if pot := list.GetSet(num - 1); pot != nil {
					part := pot.GetString("part")
					mark := lik.StrToInt(val) > 0
					if rule.ItPage.Tune.ItGen != it.itGen || rule.ItPage.Tune.ItKey != it.itEnt {
						it.partCmdClipClear(rule)
					}
					present := false
					for np,pt := range rule.ItPage.Tune.List {
						if pt == part {
							if !mark {
								list := rule.ItPage.Tune.List[:np]
								list = append(list, rule.ItPage.Tune.List[np+1:]...)
								rule.ItPage.Tune.List = list
							}
							present = true
							break
						}
					}
					if !present && mark {
						rule.ItPage.Tune.List = append(rule.ItPage.Tune.List, part)
					}
				}
			}
		}
	}
}

//	Проверка наличия в отмеченных
func (it *TuneControl) partCmdClipTest(rule *repo.DataRule, part string) bool {
	result := false
	if rule.ItPage.Tune.ItGen == it.itGen && rule.ItPage.Tune.ItKey == it.itEnt {
		for _,pit := range rule.ItPage.Tune.List {
			if pit == part {
				result = true
				break
			}
		}
	}
	return result
}

//	Копирование отмеченных
func (it *TuneControl) partCmdClipCopy(rule *repo.DataRule) {
	if len(rule.ItPage.Tune.List) > 0 {
		rule.ItPage.Tune.Copied = true
		rule.ItPage.Tune.Cuted = false
	}
}

//	Вырезание отмеченных
func (it *TuneControl) partCmdClipCut(rule *repo.DataRule) {
	if len(rule.ItPage.Tune.List) > 0 {
		rule.ItPage.Tune.Copied = true
		rule.ItPage.Tune.Cuted = true
	}
}

//	Вставка отмеченных
func (it *TuneControl) partCmdClipPaste(rule *repo.DataRule) {
	if !rule.ItPage.Tune.Copied { return }
	if entfrom := repo.SystemFindGenEnt(rule.ItPage.Tune.ItGen, rule.ItPage.Tune.ItKey);
				entfrom != nil && len(rule.ItPage.Tune.List) > 0 {
		if contfrom := jone.CalculateElmList(entfrom.It, "content"); contfrom != nil {
			if entto := repo.SystemFindGenEnt(it.itGen, it.itEnt); entto != nil {
				contto := jone.CalculateElmList(entto.It, "content")
				if contto == nil {
					contto = lik.BuildList()
					entto.It.SetValue(contto, "content")
				}
				for _,part := range rule.ItPage.Tune.List {
					if pos,prm := entfrom.FindPart(part); prm != nil {
						if rule.ItPage.Tune.Cuted {
							contfrom.DelItem(pos)
						}
						contto.AddItems(prm.Clone())
					}
				}
				entfrom.SaveToBase()
				entto.SaveToBase()
				rule.OnChangeData()
			}
		}
	}
	if rule.ItPage.Tune.Cuted {
		it.partCmdClipClear(rule)
	}
}

//	Очистка отмеченных
func (it *TuneControl) partCmdClipClear(rule *repo.DataRule) {
	rule.ItPage.Tune.ItGen = it.itGen
	rule.ItPage.Tune.ItKey = it.itEnt
	rule.ItPage.Tune.List = []string{}
	rule.ItPage.Tune.Copied = false
	rule.ItPage.Tune.Cuted = false
}

