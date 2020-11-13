package controls

import (
	"github.com/massarakhsh/tsan/control"
	"github.com/massarakhsh/tsan/fancy"
	"github.com/massarakhsh/tsan/jone"
	"github.com/massarakhsh/tsan/repo"
	"github.com/massarakhsh/tsan/show"
	"github.com/massarakhsh/lik"
	"github.com/massarakhsh/lik/likbase"
	"github.com/massarakhsh/lik/likdom"
	"github.com/massarakhsh/lik/likpdf"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
)

//	Дескриптор печати
type PrintControl struct {
	control.DataControl
	fancy.DataFancy
	Pdf       likpdf.PDFFiler	//	Генератор PDF
	Elm       *likbase.ItElm	//	Объект для экспорта
	Format		string
	IsMap     	bool				//	Имеется карта
	IsScheme  	bool				//	Имеется план - схема
	ImgPhotos 	[]string			//	Список фотографий
}

//	Интерфейс событий печати
type dealPrintExecute struct {
	It	*PrintControl
}
func (it *dealPrintExecute) Run(rule *repo.DataRule, cmd string, data lik.Seter) {
	it.It.PrintExecute(rule, cmd, data)
}

//	Конструктор контроллера печати
func BuildPrint(rule *repo.DataRule, id lik.IDB) *PrintControl {
	it := &PrintControl{}
	it.ControlInitialize("print", id)
	it.Sel = likbase.IDBToStr(it.IdMain)
	it.Fun = fancy.FunShow
	it.FancyInitialize("print", "offer","print")
	it.ItExecute = &dealPrintExecute{it}
	return it
}

//	Обработчик событий печати
func (it *PrintControl) PrintExecute(rule *repo.DataRule, cmd string, data lik.Seter) {
	if cmd == "print" || cmd == "all" {
		it.PrintExecute(rule, rule.Shift(), data)
	} else if cmd == "showform" {
		it.cmdShowForm(rule)
	} else if cmd == "printshort" {
		it.cmdPrintPdf(rule, "short")
	} else if cmd == "printlong" {
		it.cmdPrintPdf(rule, "long")
	} else if cmd == "printlist" {
		it.cmdPrintPdf(rule, "list")
	}
}

//	Отображение контроллера
func (it *PrintControl) BuildShow(rule *repo.DataRule, sx int, sy int) likdom.Domer {
	it.SetSize(sx, sy)
	return show.BuildFancyForm(it.Main,"editor")
}

//	Отображение формы
func (it *PrintControl) cmdShowForm(rule *repo.DataRule) {
	it.FormClear()
	it.RunFieldsFill(rule)
	it.Form.Title = "Печать в файл PDF"
	it.Form.Tools.AddItemSet("text=Отмена", "handler=function_fancy_form_cancel")
	it.Form.Items.AddItemSet("type=html", "value", "Выберите что печатать<br><br>")
	it.Form.Items.AddItemSet("type=html", "value",
		it.printCmdWhat(rule, "printshort", "Краткая презентация"))
	it.Form.Items.AddItemSet("type=html", "value",
		it.printCmdWhat(rule, "printlong", "Полная презентация"))
	cnt := rule.ItPage.CountCollect("offer")
	it.Form.Items.AddItemSet("type=html", "value",
		it.printCmdWhat(rule, "printlist", fmt.Sprintf("Список вариантов (%d)", cnt)))
	it.Form.SetSize(400,0)
	it.ShowForm(rule)
}

//	Выбор режима презентации
func (it *PrintControl) printCmdWhat(rule *repo.DataRule, code string, text string) string {
	return show.LinkTextProc("cmd", text, "print_as('" + code + "')").ToString()
}

//	Изготовление PDF-файла
func (it *PrintControl) cmdPrintPdf(rule *repo.DataRule, mode string) {
	it.Format = mode
	if mode == "long" || mode == "list" {
		it.Pdf = likpdf.Create("L,A4")
	} else {
		it.Pdf = likpdf.Create("P,A4")
	}
	if mode == "short" {
		it.printPrepareData(rule)
		it.printBuildShowShort(rule)
	} else if mode == "long" {
		it.printPrepareData(rule)
		it.printBuildShowLong(rule)
	} else {
		it.printBuildShowList(rule)
	}
	path := fmt.Sprintf("var/pdf/%09d.pdf", rand.Int31n(1000000000))
	it.Pdf.SaveToFile(path)
	rule.SetResponse("/" + path, "_function_goto_file")
}

//	Подготовка данных
func (it *PrintControl) printPrepareData(rule *repo.DataRule) {
	pmap := ""
	pscheme := ""
	photos := []string{}
	level := len(rule.ItPage.Locates)
	it.Elm = jone.TableOffer.GetElm(rule.ItPage.Locates[level - 1].GetId())
	if item := jone.CalculateElm(it.Elm,"objectid/picture"); item != nil {
		if lpic := item.ToList(); lpic != nil {
			for npic := 0; npic < lpic.Count(); npic++ {
				if pic := lpic.GetSet(npic); pic != nil {
					if pic.GetString("media") == "photo" {
						if url := pic.GetString("url"); url == "" {
						} else if tp := pic.GetString("imagetype"); tp == "map" {
							if pmap == "" {
								pmap = url
							}
						} else if tp == "scheme" {
							if pscheme == "" {
								pscheme = url
							}
						} else if pic.GetString("album") != "nn" {
							photos = append(photos, url)
						}
					}
				}
			}
		}
	}
	if pmap == "" {
		pmap = it.printMakeMap(rule)
	}
	it.ImgPhotos = []string{}
	if pmap != "" {
		it.IsMap = true
		it.ImgPhotos = append(it.ImgPhotos, pmap)
	}
	if pscheme != "" {
		it.IsScheme = true
		it.ImgPhotos = append(it.ImgPhotos, pscheme)
	}
	it.ImgPhotos = append(it.ImgPhotos, photos...)
}

//	Вывод короткой презентации
func (it *PrintControl) printBuildShowShort(rule *repo.DataRule) {
	it.Pdf.AddPage()
	it.printBuildLogo(rule,0,0, 0.66, 0.1)
	it.printBuildMember(rule, 0.66, 0, 1, 0.1)
	it.printBuildTitle(rule, 0, 0.1, 1, 0.15)
	it.printBuildCharacter(rule,0, 0.15, 0.5, 0.5, false)
	it.printBuildMap(rule, 0.5, 0.15, 1, 0.5)
	for nimg := 1; nimg < 5 && nimg < len(it.ImgPhotos); nimg++ {
		x := 0 + 0.5 * float64((nimg-1) % 2)
		y := 0.5 + 0.25 * float64((nimg-1) / 2)
		it.printBuildPhoto(rule, x, y, x+0.5, y+0.25, it.ImgPhotos[nimg])
	}
}

//	Вывод полной презентации
func (it *PrintControl) printBuildShowLong(rule *repo.DataRule) {
	mimg := len(it.ImgPhotos)
	nimg := 0
	for npage := 0; npage < 2 || nimg < mimg; npage++ {
		it.Pdf.AddPage()
		it.printFullTop(rule)
		npos := 0
		if npage == 0 {
			it.printBuildCharacter(rule, 0, 0.2, 0.5, 1, true)
			it.printBuildMap(rule, 0.5, 0.1, 1, 0.5)
			if it.IsMap && it.IsScheme && mimg > 2 {
				it.printBuildPhoto(rule, 0.5, 0.5, 1, 1, it.ImgPhotos[2])
			} else if mimg > 1 {
				it.printBuildPhoto(rule, 0.5, 0.5, 1, 1, it.ImgPhotos[1])
			}
			nimg = 2
			npos = 4
		} else if npage == 1 {
			it.printBuildDefine(rule, 0, 0.2, 0.5, 0.5)
			if it.IsMap && it.IsScheme && mimg > 2 {
				it.printBuildPhoto(rule, 0.5, 0.1, 1, 0.5, it.ImgPhotos[1])
			} else if mimg > 2 {
				it.printBuildPhoto(rule, 0.5, 0.1, 1, 0.5, it.ImgPhotos[2])
			}
			nimg = 3
			npos = 2
		}
		for pos := npos; pos < 4 && nimg < mimg; pos++ {
			x := 0 + 0.5 * float64(pos % 2)
			y := 0.1 + 0.45 * float64(pos / 2)
			it.printBuildPhoto(rule, x, y, x+0.5, y+0.45, it.ImgPhotos[nimg])
			nimg++
		}
	}
}

//	Вывод верхней панели с заголовками
func (it *PrintControl) printFullTop(rule *repo.DataRule) {
	it.printBuildLogo(rule,0,0, 0.5, 0.1)
	it.printBuildTitle(rule, 0.5, 0, 0.75, 0.1)
	it.printBuildMember(rule, 0.75, 0, 1, 0.1)
}

//	Вывод логотипа компании
func (it *PrintControl) printBuildLogo(rule *repo.DataRule, x1,y1,x2,y2 float64) {
	px1,py1,px2,py2 := it.Pdf.ToPad(x1, y1, x2, y2, 0.1)
	dir,_ := os.Getwd()
	it.Pdf.DrawImage(px1, py1, px2, py2, dir + "/images/tsan.jpg", "")
}

//	Вывод панели риэлтора
func (it *PrintControl) printBuildMember(rule *repo.DataRule, x1,y1,x2,y2 float64) {
	px1,py1,px2,py2 := it.Pdf.ToPad(x1, y1, x2, y2, 0.1)
	px1,py1,px2,py2 = it.Pdf.ToRatio(px1, py1, px2, py2, 0.25)
	var operator *likbase.ItElm
	if it.Format == "short" || it.Format == "long" {
		operator = jone.TableMember.GetElm(it.Elm.GetIDB("memberoid"))
	} else {
		operator = rule.GetMember()
	}
	dx1 := px1 * 0.75 + px2 * 0.25
	dx2 := dx1 * 0.95 + px2 * 0.05
	dy1 := py1 * 0.66 + py2 * 0.34
	dy2 := py1 * 0.34 + py2 * 0.66
	if photo := operator.GetString("photo"); photo != "" {
		dir, _ := os.Getwd()
		if _, err := os.Stat(dir + photo); err == nil {
			it.Pdf.DrawImage(px1, py1, dx2, py2, dir+photo, "")
		}
	}
	fio := operator.GetString("namely") + " " + operator.GetString("family")
	it.Pdf.DrawText(dx2, py1, px2, dy1, fio, "F8,#000")
	phone := operator.GetString("prophone")
	if phone == "" {
		phone = operator.GetString("phone")
	}
	it.Pdf.DrawText(dx2, dy1, px2, dy2, show.PhoneToFormat(phone), "F8,#000")
	it.Pdf.DrawText(dx2, dy2, px2, py2, operator.GetString("email"), "F8,#000")
}

//	Вывод заголовка окна
func (it *PrintControl) printBuildTitle(rule *repo.DataRule, x1,y1,x2,y2 float64) {
	title := fmt.Sprintf("№%d", int(it.Elm.Id))
	if realty := jone.CalculateElmString(it.Elm,"objectid/realty"); realty == "room" {
		title += ", комната"
	} else if realty == "flat" {
		title += ", "
		if rooms := jone.CalculateElmString(it.Elm,"objectid/define/rooms"); rooms == "1" {
			title += "1-комнатная квартира"
		} else if rooms == "2" {
			title += "2-комнатная квартира"
		} else if rooms == "3" {
			title += "3-комнатная квартира"
		} else if rooms == "4" {
			title += "4-комнатная квартира"
		} else if rooms == "5" {
			title += "5-комнатная квартира"
		} else if rooms == "1e" {
			title += "квартира-студия"
		} else {
			title += "квартира"
		}
	}
	if sub := jone.CalculateElmTranslate(it.Elm,"objectid/address/subcity"); sub != "" {
		title += ", " + sub
	}
	it.Pdf.DrawText(x1, y1, x2, y2, title, "F16,C,M,#000")
}
//	Вывод текстового окна
func (it *PrintControl) printBuildCharacter(rule *repo.DataRule, x1,y1,x2,y2 float64, islong bool) {
	xm := x1*0.5 + x2*0.5
	hline := 0.025
	keyset := "print_short"
	if islong {
		keyset = "print_long"
		hline = 0.035
	}
	target := jone.CalculateElmTranslate(it.Elm, "target")
	realty := jone.CalculateElmTranslate(it.Elm, "objectid/realty")
	if ent := repo.GenStruct.FindEnt(keyset); ent != nil {
		content := ent.GetContent()
		ct := 0
		for _,pt := range content {
			tags := pt.GetInt("tags")
			if it.FancyProbeTags(target, realty, tags) {
				y := y1 + float64(ct) * hline
				if y + hline > y2 {
					break
				}
				ttl := pt.GetString("name")
				part := pt.GetString("part")
				txt := jone.CalculateElmTranslate(it.Elm, part)
				if part == "cost" {
					txt = show.CashToFormat(txt)
				}
				it.Pdf.DrawText(x1, y, xm, y + hline, ttl, "F12,L,M,#888")
				it.Pdf.DrawText(xm, y, x2, y + hline, txt, "F12,L,M,#000")
				ct++
			}
		}
	}
}

//	Вывод описания объекта
func (it *PrintControl) printBuildDefine(rule *repo.DataRule, x1,y1,x2,y2 float64) {
	text := jone.CalculateElmString(it.Elm,"objectid/definition")
	it.Pdf.DrawText(x1, y1, x2, y2, text, "F12,L,T,#000")
}

//	Вывод карты
func (it *PrintControl) printBuildMap(rule *repo.DataRule, x1,y1,x2,y2 float64) {
	ym := y1 * 0.95 + y2 * 0.05
	if addr := jone.CalculateElmSet(it.Elm, "objectid/address"); addr != nil {
		text := ""
		if val := jone.CalculateString(addr, "street"); val != "" {
			text = val
		}
		if val := jone.CalculateString(addr, "home"); val != "" {
			text += ", д." + val
		}
		if val := jone.CalculateString(addr, "build"); val != "" {
			text += "-" + val
		}
		it.Pdf.DrawText(x1, y1, x2, ym, text, "F12,C,M,#000")
	}
	if len(it.ImgPhotos) > 0 {
		it.printBuildPhoto(rule, x1,ym,x2,y2, it.ImgPhotos[0])
	}
}

//	Изготовление и вывод карты
func (it *PrintControl) printMakeMap(rule *repo.DataRule) string {
	dap := repo.BuildMap(rule, it.Elm)
	if dap == nil { return "" }
	sx,sy := 400,400
	if it.Format == "long" {
		sy = 200
	}
	url := "https://static-maps.yandex.ru/1.x/"
	url += fmt.Sprintf("?size=%d,%d&l=map&ll=%.6f,%.6f&z=%d", sx,sy, dap.CenterY, dap.CenterX, int(dap.Zoom))
	if dap.Points != nil && dap.Points.Count() > 1 {
		url += fmt.Sprintf("&pt=%.6f,%.6f,comma", dap.Points.GetFloat(1), dap.Points.GetFloat(0))
	}
	response, e := http.Get(url)
	if e != nil {
		return ""
	}
	path := fmt.Sprintf("var/img/%09d.png", rand.Int31n(1000000000))
	os.MkdirAll("var/img", os.ModePerm)
	file, err := os.Create(path)
	if err != nil {
		response.Body.Close()
		return ""
	}
	_, err = io.Copy(file, response.Body)
	if err != nil {
		response.Body.Close()
		file.Close()
		return ""
	}
	response.Body.Close()
	file.Close()
	return "/" + path
}

//	Вывод изображения
func (it *PrintControl) printBuildPhoto(rule *repo.DataRule, x1,y1,x2,y2 float64, name string) {
	px1,py1,px2,py2 := it.Pdf.ToPad(x1, y1, x2, y2, 0.08)
	it.Pdf.DrawRect(px1, py1, px2, py2, "#ccc")
	if name != "" {
		dir, _ := os.Getwd()
		if _, err := os.Stat(dir + name); err == nil {
			px1, py1, px2, py2 = it.Pdf.ToPad(x1, y1, x2, y2, 0.1)
			it.Pdf.DrawImage(px1, py1, px2, py2, dir+name, "")
		}
	}
}

//	Вывод списка объектов
func (it *PrintControl) printBuildShowList(rule *repo.DataRule) {
	content := []lik.Seter{}
	if ent := repo.GenStruct.FindEnt("print_list"); ent != nil {
		content = ent.GetContent()
	}
	it.Pdf.AddPage()
	it.printBuildLogo(rule,0,0, 0.5, 0.1)
	it.printBuildMember(rule, 0.75, 0, 1, 0.1)
	list := []*likbase.ItElm{}
	for key,val := range rule.ItPage.Session.Collect {
		if match := lik.RegExParse(key, "^offer(\\d+)"); match !=nil && val {
			if elm := jone.TableOffer.GetElm(lik.IDB(lik.StrToInt(match[1]))); elm != nil {
				list = append(list, elm)
			}
		}
	}
	target := ""
	realty := ""
	if len(list) > 0 {
		target = jone.CalculateElmTranslate(list[0], "target")
		realty = jone.CalculateElmTranslate(list[0], "objectid/realty")
	}
	yc := 0.15
	for nl := -1; nl < len(list); nl++ {
		var elm *likbase.ItElm
		if nl >= 0 {
			elm = list[nl]
		}
		y1,y2 := yc, yc + 0.06
		xc := 0.0
		for _,field := range content {
			if tags := field.GetInt("tags"); it.FancyProbeTags(target, realty, tags) {
				part := field.GetString("part")
				text := field.GetString("name")
				xd := field.GetFloat("width") * 0.001
				opt := "M,C,F9"
				if part == "_" {
					xd = 0.025
					if elm != nil {
						text = fmt.Sprintf("%d.", 1+nl)
					}
				} else if part == "id" {
					if elm != nil {
						text = fmt.Sprintf("%d.", int(elm.Id))
					}
				} else {
					if elm != nil {
						text = jone.CalculateElmTranslate(elm, part)
					}
				}
				if xd == 0 {
					xd = 0.02
				}
				x1, x2 := xc, xc+xd
				if x2 > 1.0 {
					break
				}
				it.Pdf.DrawRect(x1, y1, x2, y2, "#888")
				it.Pdf.DrawText(x1, y1, x2, y2, text, opt+",#000")
				xc = x2
			}
		}
		yc = y2
	}
}

