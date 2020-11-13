package window

import (
	"github.com/massarakhsh/tsan/one"
	"github.com/massarakhsh/tsan/repo"
	"github.com/massarakhsh/lik"
)

//	Отображение окна клиента, построенного на объекте FancyGrid
func FancyShowCode(rule *repo.DataRule, it *ClientBox) {
	code := it.Parameters.Clone().ToSet()

	if !code.IsItem("title") && (it.Title != "" || it.Titles.Count() > 0) {
		title := lik.BuildSet()
		title.SetItem(it.Title, "text")
		if it.Titles.Count() > 0 {
			title.SetItem(it.Titles,"tools")
		}
		code.SetItem(title,"title")
	}
	if !code.IsItem("tbar") && it.Cmds.Count() > 0 {
		for n := 0; n < it.Cmds.Count(); n++ {
			if item := it.Cmds.GetSet(n); item != nil && item.GetString("type") == "" {
				it.Cmds.SetItem("|", n)
			}
		}
		code.SetItem(it.Cmds, "tbar")
	}

	if it.IsCollect {
		if !code.IsItem("columns") {
			if it.Columns.Count() > 0 {
				code.SetItem(it.Columns, "columns")
			} else {
				code.SetItem(lik.BuildSet("title=(Нет колонок)"), "columns")
			}
		}

		if !code.IsItem("contextmenu") && it.Context.Count() > 0 {
			code.SetItem(it.Context, "contextmenu")
		}

		data := code.GetSet("data")
		if data == nil {
			data = lik.BuildSet()
			code.SetItem(data, "data")
		}

		if it.Class == "" {
			it.Class = "show-fancy-grid"
		}
	}

	if it.IsForm {
	}

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
	if !code.IsItem("likMain") {
		code.SetItem(it.Frame, "likMain")
	}
	if !code.IsItem("likZone") {
		code.SetItem(it.Mode, "likZone")
	}
	if !code.IsItem("license") {
		code.SetItem(one.FancyLicense, "license")
	}
	rule.SetResponse(code, "fancy")
}

//	Отображение страницы, построенной на объекте fancyGrid
func FancyShowPage(rule *repo.DataRule, it *ClientBox, rows lik.Lister) {
	rule.SetResponse(rows, "items")
	if rows != nil {
		rule.SetResponse(rows.Count(), "totalCount")
		//if cur := it.GridSeekSel(rule, rows); cur >= 0 {
		//	rule.SetResponse(cur, "likSelect")
		//}
	}
	rule.SetResponse(true,"success")
}

