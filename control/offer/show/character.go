package show

import (
	"bitbucket.org/961961/tsan/control"
	"bitbucket.org/961961/tsan/fancy"
	"bitbucket.org/961961/tsan/jone"
	"bitbucket.org/961961/tsan/repo"
	"bitbucket.org/961961/tsan/show"
	"bitbucket.org/shaman/lik"
	"bitbucket.org/shaman/lik/likdom"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

//	Дескриптор характеристик объекта
type CharacterControl struct {
	control.DataControl
	fancy.TableFancy
}

//	Интерфейс команд
type dealCharExecute struct {
	It	*CharacterControl
}
func (it *dealCharExecute) Run(rule *repo.DataRule, cmd string, data lik.Seter) {
	it.It.TableExecute(rule, cmd, data)
}

//	Интерфейс заполнения таблицы
type dealCharGridFill struct {
	It	*CharacterControl
}
func (it *dealCharGridFill) Run(rule *repo.DataRule) {
	it.It.CharGridFill(rule)
}

//	Интерфейс заполнения страницы
type dealCharPageFill struct {
	It	*CharacterControl
}
func (it *dealCharPageFill) Run(rule *repo.DataRule) lik.Lister {
	return it.It.CharPageFill(rule)
}

//	Конструктор дескриптора характеристик объекта
func BuildCharacter(rule *repo.DataRule, main string, id lik.IDB) *CharacterControl {
	it := &CharacterControl{ }
	it.ControlInitializeZone(main, id, "char")
	it.TableInitialize(rule, main,"offer","char")
	it.ItExecute = &dealCharExecute{it}
	it.ItGridFill = &dealCharGridFill{it}
	it.ItPageFill = &dealCharPageFill{it}
	it.SeekLocate(rule)
	return it
}

//	Позиционирование
func (it *CharacterControl) SeekLocate(rule *repo.DataRule) bool {
	if elm := jone.GetElm("offer", it.IdMain); elm != nil {
		//target := elm.GetString("target")
	}
	return true
}

//	Отображение характеристик объекта
func (it *CharacterControl) BuildShow(rule *repo.DataRule, sx int, sy int) likdom.Domer {
	it.SetSize(sx, sy)
	return show.BuildFancyGrid(it.Main,"char")
}

//	Заполнение таблицы
func (it *CharacterControl) CharGridFill(rule *repo.DataRule) {
	it.TableGridFill(rule)
	it.Grid.Columns.AddItemSet("title=Поле", "index=name", "width=100")
	it.Grid.Columns.AddItemSet("title=Значение", "index=value", "width=200")
}

//	Заполнение страницы
func (it *CharacterControl) CharPageFill(rule *repo.DataRule) lik.Lister {
	rows := it.TablePageFill(rule)
	if elm := jone.TableOffer.GetElm(it.IdMain); elm != nil {
		target := elm.GetString("target")
		realty := elm.GetString("realty")
		if ent := repo.GenStruct.FindEnt(target + "_show"); ent != nil {
			if content := ent.It.GetList("content"); content != nil {
				for nc := 0; nc < content.Count(); nc++ {
					if col := content.GetSet(nc); col != nil {
						tags := col.GetInt("tags")
						if (tags &jone.TagShow) != 0 && it.FancyProbeTags(target, realty, tags) {
							part := col.GetString("part")
							format := col.GetString("format")
							if value := it.RunCalculate(rule, elm.Info, part, format, false); value != nil {
								row := lik.BuildSet()
								row.SetItem(col.GetString("name"), "name")
								if match := lik.RegExParse(part,"([^/]+)$"); match != nil {
									part = match[1]
									row.SetItem(part, "key")
								}
								dat := lik.BuildItem(jone.SystemTranslate(part, value))
								row.SetItem(dat, "value")
								rows.AddItems(row)
							}
						}
					}
				}
			}
		}
	}
	if rows.Count() <= 0 {
		rows.AddItemSet("id=1", "nam=НЕТ")
	}
	return rows
}

