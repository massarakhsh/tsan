package show

import (
	"bitbucket.org/961961/tsan/control"
	"bitbucket.org/961961/tsan/jone"
	"bitbucket.org/961961/tsan/repo"
	"bitbucket.org/961961/tsan/show"
	"bitbucket.org/shaman/lik"
	"bitbucket.org/shaman/lik/likdom"
	"fmt"
	"os"
)

//	Дескриптор галереи изображений
type ShowGallery struct {
	control.DataControl
}

//	Интерфейс команд
type dealGalleryExecute struct {
	It	*ShowGallery
}
func (it *dealGalleryExecute) Run(rule *repo.DataRule, cmd string, data lik.Seter) {
	it.It.GalleryExecute(rule, cmd, data)
}

//	Конструктор галереи
func BuildGallery(rule *repo.DataRule, main string, id lik.IDB) *ShowGallery {
	it := &ShowGallery{}
	it.ControlInitializeZone(main, id, "gallery")
	it.ItExecute = &dealGalleryExecute{it}
	return it
}

//	Позиционирование
func (it *ShowGallery) SeekLocate(rule *repo.DataRule) bool {
	return false
}

//	Выполнение команд
func (it *ShowGallery) GalleryExecute(rule *repo.DataRule, cmd string, data lik.Seter) {
}

//	Отображение галереи
func (it *ShowGallery) BuildShow(rule *repo.DataRule, sx int, sy int) likdom.Domer {
	div := likdom.BuildItem("div")
	data := div.BuildItem("div","class=fotorama", "id=fotorama")
	hf := 64
	if elm := jone.TableOffer.GetElm(it.IdMain); elm != nil {
		data.SetAttr("data-nav='thumbs'",
			fmt.Sprintf("width=%d", sx),
			fmt.Sprintf("height=%d", sy-hf),
			fmt.Sprintf("data-thumbwidth=%d", hf),
			fmt.Sprintf("data-thumbheight=%d", hf),
			fmt.Sprintf("data-width=%d", sx),
			fmt.Sprintf("data-height=%d", sy-hf),
			fmt.Sprintf("data-maxwidth=%d", sx),
			fmt.Sprintf("data-maxheight=%d", sy-hf),
		)
		if item := jone.CalculateElm(elm,"objectid/picture"); item != nil {
			if lpic := item.ToList(); lpic != nil {
				for npic := 0; npic < lpic.Count(); npic++ {
					if pic := lpic.GetSet(npic); pic != nil {
						if pic.GetString("media") == "photo" && pic.GetString("album") != "nn" {
							if url := pic.GetString("url"); url != "" {
								a := data.BuildItem("a", "href", url)
								if match := lik.RegExParse(url, "^/(.+)\\.jpg$"); match != nil {
									urlt := match[1] + "t.jpg"
									if _, err := os.Stat(urlt); err != nil {
										show.MakeScaleJpg(match[1] + ".jpg", 0, hf, urlt)
									}
									a.BuildUnpairItem("img", "src", urlt)
								}
							}
						}
					}
				}
			}
		}
	}
	code := "var foto = $('#fotorama').fotorama();"
	code += "var rama = foto.data('fotorama');"
	code += fmt.Sprintf("rama.resize({ width: %d, height: %d });", sx, sy)
	div.AppendItem(show.BuildRunScript(code,250))
	return div
}

