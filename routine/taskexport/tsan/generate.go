package tsan

import (
	"github.com/massarakhsh/tsan/jone"
	"github.com/massarakhsh/tsan/show"
	"github.com/massarakhsh/lik"
	"github.com/massarakhsh/lik/likbase"
	"fmt"
	"strings"
)

//	Подготовка отчета
func (it *TaskTsan) PrepareReport() bool {
	it.Report = lik.BuildList()
	it.ListKeys = []string{
		"code", "name", "alias", "category", "category_primary",
		"author", "state", "comments", "frontpage", "metadata_title",
		"metadata_description", "metadata_keywords", "tags", "tags", "tags",
		"Offer-category", "Offer-building-name", "Offer-market-title", "Offer-rooms", "Offer-rooms-filter",
		"Offer-type-flat", "Offer-renovation", "Offer-ceiling-height", "Offer-balcony", "Offer-bathroom-unit",
		"Offer-floor", "Offer-floors", "Offer-floor-floors", "Offer-floor-between", "Offer-window-view",
		"Offer-building-type", "Offer-heating", "Offer-hot_water", "Offer-description", "Offer-lift",
		"Offer-haggle", "Offer-built-year", "Offer-mortgage", "Offer-cadastral-number", "Offer-video",
		"Price-value-old", "Price-value", "price_round", "area_round", "Price-value-metr",
		"Location-region", "Location-district", "Location-locality-name", "Location-sub-locality-name", "Location-street",
		"Location-dom-num", "Location-korp-num", "Location-coordinates", "Area-value", "Area-kitchen-space",
		"Area-living-space", "SalesAgent-name", "SalesAgent-phone", "SalesAgent-email", "SalesAgent-photo",
		"Offer-image",
	}
	list := lik.BuildList()
	for _,key := range it.ListKeys {
		if key == "Offer-image" {
			for nr := 0; nr < 21; nr++ {
				list.AddItems(key)
			}
		} else {
			list.AddItems(key)
		}
	}
	it.Report.AddItems(list)
	return true
}

//	Запись файла
func (it *TaskTsan) WriteFiles() {
	it.WriteFileCSV(it.Report, "^")
}

//	Добавление в отчёт ЦАН заявки
//	elm - добавляемая заявка
func (it *TaskTsan) AppendOffer(elm *likbase.ItElm) string {
	line := lik.BuildList()
	for _,key := range it.ListKeys {
		data := ""
		if key == "code" {
			idu := elm.GetIDB("idu")
			if idu == 0 { idu = elm.Id }
			data = fmt.Sprintf("%d", int(idu))
		} else if key == "name" {
			if rooms := jone.CalculateElmString(elm, "objectid/defime/rooms"); rooms != "" {
				if lik.StrToInt(rooms) > 0 {
					data += fmt.Sprintf("%s комн. квартира, ", rooms)
				}
			}
			if square := jone.CalculateElmFloat(elm, "objectid/defime/square"); square > 0 {
				data += fmt.Sprintf("%.2f кв.м., ", square)
			}
			data += fmt.Sprintf("арт. %d", int(elm.Id))
		} else if key == "alias" {
			data = fmt.Sprintf("%d", int(elm.Id))
		} else if key == "category" {
			data = "1|||1"
		} else if key == "category_primary" {
			data = "1|||1"
		} else if key == "author" {
			data = "Автозагрузка"
		} else if key == "state" {
			data = "1"
		} else if key == "comments" {
			data = "0"
		} else if key == "frontpage" {
			data = "1"
		} else if key == "metadata_title" {
			data = ""
		} else if key == "metadata_description" {
			data = ""
		} else if key == "metadata_keywords" {
			data = ""
		} else if key == "tags" {
			data = "bag@tsan.ru"

		} else if key == "tags" {
			data = ""	//"Островского ул."
		} else if key == "tags" {
			data = ""
		} else if key == "Offer-category" {
			data = jone.CalculateElmTranslate(elm, "objectid/realty")
		} else if key == "Offer-building-name" {
			data = ""
		} else if key == "Offer-market-title" {
			//data = "1 комн. квартира, 37.80 кв.м"
		} else if key == "Offer-rooms" {
			data = jone.CalculateElmTranslate(elm, "objectid/define/rooms")
		} else if key == "Offer-rooms-filter" {
			//data = "1"
		} else if key == "Offer-type-flat" {
			//data = ""
		} else if key == "Offer-renovation" {
			data = jone.CalculateElmTranslate(elm, "objectid/define/renovation")
		} else if key == "Offer-ceiling-height" {
			data = jone.CalculateElmTranslate(elm, "objectid/define/ceilingheight")
		} else if key == "Offer-balcony" {
			data = jone.CalculateElmTranslate(elm, "objectid/define/balcony")
		} else if key == "Offer-bathroom-unit" {
			//data = "совмещенный"
		} else if key == "Offer-floor" {
			data = jone.CalculateElmTranslate(elm, "objectid/define/floor")
		} else if key == "Offer-floors" {
			data = jone.CalculateElmTranslate(elm, "objectid/define/floortotal")
		} else if key == "Offer-floor-floors" {
			data = jone.CalculateElmTranslate(elm, "objectid/define/floor") + "/" +
				jone.CalculateElmTranslate(elm, "objectid/define/floortotal")
		} else if key == "Offer-floor-between" {
			//data = "Первый"
		} else if key == "Offer-window-view" {
			//data = ""
		} else if key == "Offer-building-type" {
			data = jone.CalculateElmTranslate(elm, "objectid/define/housetype")
		} else if key == "Offer-heating" {
			data = jone.CalculateElmTranslate(elm, "objectid/define/heatsystem")
		} else if key == "Offer-hot_water" {
			//data = ""
		} else if key == "Offer-description" {
			data = jone.CalculateElmString(elm, "objectid/definition")
			data = strings.Replace(data, "\r", "\\r", -1)
			data = strings.Replace(data, "\n", "\\n", -1)
		} else if key == "Offer-lift" {
			data = jone.CalculateElmTranslate(elm, "objectid/define/elevator")
		} else if key == "Offer-haggle" {
			//data = "нет"
		} else if key == "Offer-built-year" {
			//data = "1988"
		} else if key == "Offer-mortgage" {
			//data = "да"
		} else if key == "Offer-cadastral-number" {
			//data = "62:29:0070034:4278"
		} else if key == "Offer-video" {
			//data = ""
		} else if key == "Price-value-old" {
			//data = "1590000"
		} else if key == "Price-value" {
			//data = "1500000"
		} else if key == "price_round" {
			//data = ""
		} else if key == "area_round" {
			//data = "45"
		} else if key == "Price-value-metr" {
			//data = "42063,4920634921"
		} else if key == "Location-region" {
			data = jone.CalculateElmTranslate(elm, "objectid/address/region")
		} else if key == "Location-district" {
			data = jone.CalculateElmTranslate(elm, "objectid/address/district")
		} else if key == "Location-locality-name" {
			data = jone.CalculateElmTranslate(elm, "objectid/address/city")
		} else if key == "Location-sub-locality-name" {
			data = jone.CalculateElmTranslate(elm, "objectid/address/subcity")
		} else if key == "Location-street" {
			data = jone.CalculateElmTranslate(elm, "objectid/address/street")
		} else if key == "Location-dom-num" {
			data = jone.CalculateElmTranslate(elm, "objectid/address/home")
		} else if key == "Location-korp-num" {
			data = jone.CalculateElmTranslate(elm, "objectid/address/build")
		} else if key == "Location-coordinates" {
			if cx, cy := it.MakePoint(elm); cx != 0 && cy != 0 {
				data = fmt.Sprintf("%.6f,%.6f", cx, cy)
			}
		} else if key == "Area-value" {
			data = jone.CalculateElmTranslate(elm, "objectid/define/square")
		} else if key == "Area-kitchen-space" {
			data = jone.CalculateElmTranslate(elm, "objectid/define/squarekitchen")
		} else if key == "Area-living-space" {
			data = jone.CalculateElmTranslate(elm, "objectid/define/squareliving")
		} else if key == "SalesAgent-name" {
			data = jone.CalculatePartIdText("member", elm.GetIDB("memberid"))
		} else if key == "SalesAgent-phone" {
			data = jone.NormalizePhone(jone.CalculatePartIdString("member", elm.GetIDB("memberid"), "phone1"))
		} else if key == "SalesAgent-email" {
			//data = "bag@tsan.ru"
		} else if key == "SalesAgent-photo" {
			//data = ""
		} else if key == "Offer-image" {
			if lpic := jone.CalculateElmList(elm,"objectid/picture"); lpic != nil {
				nph := 0
				for np := 0; np < lpic.Count() && nph < 20; np++ {
					if pic := lpic.GetSet(np); pic != nil &&
						pic.GetString("media") == "photo" &&
						pic.GetString("album") != "nn" &&
						pic.GetString("promot") != "nn" {
						if url := pic.GetString("url"); url != "" {
							nph++
							line.AddItems(show.UrlToString(url))
						}
					}
				}
			}
		}
		line.AddItems(data)
	}
	it.Report.AddItems(line)
	return ""
}

