package tune

import (
	"bitbucket.org/961961/tsan/repo"
	"bitbucket.org/shaman/lik"
	"bitbucket.org/shaman/lik/likdom"
	"fmt"
	"regexp"
	"strings"
)

//	Команды работы с деревом
func (it *TuneControl) TreeExecute(rule *repo.DataRule, cmd string, data lik.Seter) {
	if cmd == "select" {
		it.cmdSelectItem(rule, rule.Shift())
		return
	} else if cmd == "switch" {
		it.cmdSwitchItem(rule, rule.Shift())
		return
	}
	rule.OnChangeData()
}

//	Построение раздела с деревом
func (it *TuneControl) buildDataTree(rule *repo.DataRule, sx int, sy int) likdom.Domer {
	div := likdom.BuildItem("div")
	it.buildToPath(rule, div, "", it.Path)
	return div
}

//	Построение раздела с деревом
func (it *TuneControl) buildToPath(rule *repo.DataRule, pater likdom.Domer, path string, next string) {
	tbl := pater.BuildTableClass("tree_info")
	list := it.LoadTree(path)
	nsel := -1
	//if path == "" { nsel = 0 }
	for n := 0; n < len(list); n++ {
		elm := list[n]
		pathelm := path + "/" + elm.Part
		if strings.HasPrefix(it.Path, pathelm) {
			nsel = n
		}
	}
	for n := 0; n < len(list); n++ {
		elm := list[n]
		pathelm := elm.Part
		if strings.Contains(elm.Part, "(") {
			pathelm += ""
		} else {
			pathelm = path + "/" + pathelm
		}
		childs := []repo.Element{}
		if elm.Part != "" {
			childs = it.LoadTree(pathelm)
		}
		issel := pathelm == it.Path || it.Path == "" && n == nsel
		isext := false
		pathnext := next
		if match := regexp.MustCompile("^/" + elm.Part + "(.*)").FindStringSubmatch(next); match != nil {
			pathnext = match[1]
			if len(childs) > 0 && pathnext == "" || strings.HasPrefix(pathnext, "/") {
				isext = true;
			}
		}
		row := tbl.BuildTr()
		dir1 := 0x1
		if elm.Part != "" {
			dir1 |= 0x2
		}
		if n+1 < len(list) {
			dir1 |= 0x4
		}
		row.BuildTdClass("tree_switch").AppendItem(it.GenSwitch(rule, pathelm, dir1))
		if childs == nil || len(childs) == 0 {
			row.BuildTdClass("tree_line", "colspan=2").AppendItem(it.GenLine(rule, pathelm, elm.Name, issel))
		} else {
			dir2 := 0xa
			if isext {
				dir2 |= 0x14
			} else if childs != nil && len(childs) > 0 {
				dir2 |= 0x20
			}
			row.BuildTdClass("tree_switch").AppendItem(it.GenSwitch(rule, pathelm, dir2))
			row.BuildTdClass("tree_line").AppendItem(it.GenLine(rule, pathelm, elm.Name, issel))
			row = tbl.BuildTr()
			ltd :=row.BuildTd()
			if (dir1&0x4) != 0 {
				ltd.SetAttr("background", rule.BuildUrl("/rast/pix/o.pix"));
			}
			row.BuildTdClass("tree_line", "colspan=2").AppendItem(it.GenExpand(rule, pathelm, isext, pathnext))
		}
	}
}

//	Создание переключателя в дереве
func (it *TuneControl) GenSwitch(rule *repo.DataRule, path string, dir int) likdom.Domer {
	sfx := it.buildPixName(dir)
	img := likdom.BuildUnpairItem("img", "src", rule.BuildUrl(sfx))
	if (dir & 0x30) != 0 {
		id := lik.StringToXS(path)
		img.SetAttr("id", "id_"+id+"_switch")
		img.SetAttr("class", "active")
		proc := fmt.Sprintf("tunetree_switch('%s','%s',%d)", it.Main, id, dir)
		img.SetAttr("onclick", proc)
	}
	return img
}

//	Вычисление имени пиктограммы переключателя
func (it *TuneControl) buildPixName(dir int) string {
	pix := ""
	if (dir & 0x10) != 0 {
		pix += "m"
	} else if (dir & 0x20) != 0 {
		pix += "p"
	} else {
		pix += "z"
	}
	if (dir & 0x1) != 0 {
		pix += "1"
	} else {
		pix += "0"
	}
	if (dir & 0x2) != 0 {
		pix += "1"
	} else {
		pix += "0"
	}
	if (dir & 0x4) != 0 {
		pix += "1"
	} else {
		pix += "0"
	}
	if (dir & 0x8) != 0 {
		pix += "1"
	} else {
		pix += "0"
	}
	return "/rast/pix/"+pix+".pix"
}

//	Создание строки с именем в дереве
func (it *TuneControl) GenLine(rule *repo.DataRule, path string, name string, issel bool) likdom.Domer {
	var code likdom.Domer
	if strings.HasSuffix(path,"/") && name == "" {
		code = likdom.BuildSpace()
	} else if strings.Contains(path, "(") {
		code = likdom.BuildDivClassId("tree_elm", "", "onclick", path)
	} else {
		id := lik.StringToXS(path)
		cls := "tree_elm"
		if issel {
			cls += " tree_sel"
		}
		proc := fmt.Sprintf("tunetree_select('%s','%s')", it.Main, id)
		code = likdom.BuildDivClassId(cls, "id_"+id+"_line", "onclick", proc)
	}
	code.BuildString(name)
	return code
}

//	Раздел расширения раскрываемого объекта
func (it *TuneControl) GenExpand(rule *repo.DataRule, path string, isext bool, next string) likdom.Domer {
	id := lik.StringToXS(path)
	expand := likdom.BuildDivClassId("fill", "id_"+id+"_expand")
	if isext {
		it.buildToPath(rule, expand, path, next)
	}
	return expand
}

//	Загрузка ветки дерева
func (it *TuneControl) LoadTree(path string) []repo.Element {
	parts := lik.PathToNames(path)
	locs := []repo.Element{}
	if len(parts) == 0 {
		for _, gen := range repo.SysGens {
			locs = append(locs, BuildTreeElm(gen.Gen, gen.Name))
		}
		locs = append(locs, BuildTreeElm("devmem", "Настройки оператора"))
		locs = append(locs, BuildTreeElm("devloc", "Настройки компьютера"))
		locs = append(locs, BuildTreeElm("command", "Команды"))
		locs = append(locs, BuildTreeElm("session", "Сессии"))
		locs = append(locs, BuildTreeElm("go_godoc()", "Сервер документации"))
	} else if root := parts[0]; root == "devmem" {
		if len(parts) == 1 {
		}
	} else if root == "devloc" {
		if len(parts) == 1 {
		}
	} else if gen := repo.SystemFindGen(root); gen != nil {
		if len(parts) == 1 {
			if elms := repo.SystemSortEnts(root); elms != nil {
				for _, elm := range elms {
					name := elm.GetString("name")
					if name == "" {
						name = "(без имени)"
					}
					locs = append(locs, BuildTreeElm(elm.GetString("key"), name))
				}
			}
		}
	}
	return locs
}

//	Команда переключения в дереве
func (it *TuneControl) cmdSwitchItem(rule *repo.DataRule, parm string) {
	path := lik.StringFromXS(parm)
	dir := lik.StrToInt(rule.Shift()) ^ 0x34
	rule.StoreItem(it.GenSwitch(rule, path, dir))
	rule.StoreItem(it.GenExpand(rule, path, (dir&0x10) != 0, ""))
}

//	Команда выбора в дереве
func (it *TuneControl) cmdSelectItem(rule *repo.DataRule, parm string) {
	it.Path = lik.StringFromXS(parm)
	it.RunUpdate(rule)
	rule.OnChangeData()
}

//	Конструктор элемента дерева
func BuildTreeElm(part string, text string) repo.Element {
	return repo.Element{ Part: part, Name: text }
}

