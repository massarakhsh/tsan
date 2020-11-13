package yandex

import (
	"github.com/massarakhsh/tsan/jone"
	"github.com/massarakhsh/tsan/show"
	"github.com/massarakhsh/lik"
	"github.com/massarakhsh/lik/likbase"
	"strings"
	"time"
)

//	Подготовка отчета
func (it *TaskYandex) PrepareReport() bool {
	it.Report = it.AddContent(nil)
	it.Report.AddItems("<?xml version=\"1.0\" encoding=\"UTF-8\"?>")
	rf, offers := it.AddSetContent(it.Report, "realty-feed")
	rf.SetItem("http://webmaster.yandex.ru/schemas/feed/realty/2010-06", "xmlns")
	it.Offers = offers
	it.Offers.AddItemSet("_tag=generation-date", "_value", show.TimeToString(int(time.Now().Unix())))
	return true
}

//	Запись файла
func (it *TaskYandex) WriteFiles() {
	it.WriteFileXML(it.Report)
}

//	Добавление в отчёт Яндекса заявки
//	elm - добавляемая заявка
func (it *TaskYandex) AppendOffer(elm *likbase.ItElm) string {
	card,content := it.AddSetContent(nil, "offer")
	idu := elm.GetIDB("idu")
	if idu == 0 { idu = elm.Id }
	card.SetItem(idu, "internal-id")
	if diag := it.appendOfferMain(elm, content); diag != "" {
		return diag
	}
	if diag := it.appendOfferSaler(elm, content); diag != "" {
		return diag
	}
	if diag := it.appendOfferDeal(elm, content); diag != "" {
		return diag
	}
	if diag := it.appendOfferObject(elm, content); diag != "" {
		return diag
	}
	it.Offers.AddItems(card)
	return ""
}

//	Основные сведения по заявке
func (it *TaskYandex) appendOfferMain(elm *likbase.ItElm, content lik.Lister) string {
	info := ""
	segment := jone.CalculateElmString(elm, "segment")
	target := jone.CalculateElmString(elm, "target")
	if segment == jone.DoRent && target == "sale" {
		info = "аренда"
	} else if target == "sale" {
		info = "продажа"
	} else {
		return "Неверная цель заявки " + target
	}
	it.SetValue(content, "type", info)
	info = "жилая"
	content.AddItemSet("_tag=property-type", "_value", info)
	realty := jone.CalculateElmString(elm, "objectid/realty")
	if realty == "flat" {
		info = "квартира"
	} else if realty == "room" {
		info = "комната"
	} else if realty == "house" {
		info = "дом"
	} else if realty == "lot" {
		info = "участок"
	} else {
		return "Неверный тип недвижимости " + realty
	}
	it.SetValue(content, "category", info)
	if cadastr := jone.CalculateElmString(elm, "objectid/define/cadnum"); cadastr != "" {
		it.SetValue(content, "cadastral-number", info)
	}
	creation, lastupdate := "", ""
	if tube := jone.CalculateElmList(elm, "history"); tube != nil && tube.Count() > 0 {
		if hist := tube.GetSet(0); hist != nil {
			if dt := hist.GetInt("date"); dt > 0 {
				creation = show.TimeToString(dt)
			}
		}
		if hist := tube.GetSet(tube.Count() - 1); hist != nil {
			if dt := hist.GetInt("date"); dt > 0 {
				lastupdate = show.TimeToString(dt)
			}
		}
	}
	if creation == "" {
		return "Не указана дата создания заявки"
	}
	it.SetValue(content, "creation-date", creation)
	it.SetValue(content, "last-update-date", lastupdate)
	address := jone.CalculateElmSet(elm, "/objectid/address")
	if address == nil {
		return "Не указан вдрес объекта"
	}
	_,locations := it.AddSetContent(content, "location")
	if info = address.GetString("country"); info == "" || info == "Россия" {
		it.SetValue(locations, "country", "Россия")
	} else {
		return "Неверное название страны " + info
	}
	if info = address.GetString("region"); info != "" {
		it.SetValue(locations, "region", info)
	}
	if info = address.GetString("district"); info != "" {
		it.SetValue(locations, "district", info)
		locations.AddItemSet("_tag=district", "_value", info)
	}
	if info = address.GetString("city"); info != "" {
		it.SetValue(locations, "locality-name", info)
	}
	if info = jone.CalculateTranslate(address,"subcity"); info != "" {
		it.SetValue(locations, "sub-locality-name", info)
	}
	if info = address.GetString("street"); info != "" {
		if home := address.GetString("home"); home != "" {
			info += " " + home
		}
		if build := address.GetString("build"); build != "" {
			info += "-" + build
		}
		it.SetValue(locations, "address", info)
	}
	if cx, cy := it.MakePoint(elm); cx != 0 && cy != 0 {
		it.SetValue(locations, "latitude", cx)
		it.SetValue(locations, "longitude", cy)
	}
	return ""
}

//	Сведения о продавце
func (it *TaskYandex) appendOfferSaler(elm *likbase.ItElm, content lik.Lister) string {
	member := jone.TableMember.GetElm(jone.CalculateElmIDB(elm, "memberid"))
	if member == nil {
		return "Отсутствует ответственный риэлтор"
	}
	_,sales := it.AddSetContent(content, "sales-agent")
	info := jone.CalculateElmText(member)
	it.SetValue(sales, "name", info)
	if info := member.GetString("prophone"); info != "" {
		it.SetValue(sales, "phone", show.PhoneToString(info))
	} else if info := member.GetString("phone"); info != "" {
		it.SetValue(sales, "phone", show.PhoneToString(info))
	} else {
		return "Отсутствует телефон риэлтора"
	}
	it.SetValue(sales, "category", "агентство")
	it.SetValue(sales, "organization", "ЦАН")
	it.SetValue(sales, "url", "961-961.ru")
	if info := member.GetString("email"); info != "" {
		it.SetValue(sales, "email", info)
	}
	//sales.AddItemSet("_tag=photo", "_value", "???")
	return ""
}

//	Сведения о сделке
func (it *TaskYandex) appendOfferDeal(elm *likbase.ItElm, content lik.Lister) string {
	info := jone.CalculateElmString(elm, "cost")
	if info == "" {
		return "Отсутствует цена"
	}
	_,price := it.AddSetContent(content, "price")
	it.SetValue(price, "value", info)
	it.SetValue(price, "currency", "RUB")
	if jone.CalculateElmString(elm, "segment") == jone.DoRent {
		it.SetValue(price, "period", "месяц")
	}
	return ""
}

//	Сведения об объекте
func (it *TaskYandex) appendOfferObject(elm *likbase.ItElm, content lik.Lister) string {
	define := jone.CalculateElmSet(elm, "objectid/define")
	if define == nil {
		return "Отсутствует описание объекта"
	}
	//realty := repo.CalculateElmString(elm, "objectid/realty")
	if info := define.GetString("square"); info != "" {
		_,area := it.AddSetContent(content, "area")
		it.SetValue(area, "value", info)
		it.SetValue(area, "unit", "кв.м")
	}
	if info := define.GetString("landarea"); info != "" {
		_,area := it.AddSetContent(content, "lot-area")
		it.SetValue(area, "value", info)
		it.SetValue(area, "unit", "сотка")
	}
	if info := define.GetString("squareliving"); info != "" {
		_,area := it.AddSetContent(content, "living-space")
		it.SetValue(area, "value", info)
		it.SetValue(area, "unit", "кв.м")
	}
	if info := define.GetString("squarekitchen"); info != "" {
		_,area := it.AddSetContent(content, "kitchen-space")
		it.SetValue(area, "value", info)
		it.SetValue(area, "unit", "кв.м")
	}
	it.appendOfferPhotos(elm, content)
	if info := define.GetString("renovation"); info != "" {
		if info == "needsrepair" {
			info = "требует ремонта"
		} else if info == "cosmetic" {
			info = "косметический"
		} else if info == "qualityrepair" {
			info = "евроремонт"
		} else if info == "designerrepair" {
			info = "дизайнерский"
		}
		if info != "" {
			it.SetValue(content, "renovation", info)
		}
	}
	if info := define.GetString("quality"); info != "" {
		if info == "goodrepair" {
			info = "хорошее"
		} else if info == "excellentrepair" {
			info = "отличное"
		} else if info == "designrepair" {
			info = "отличное"
		} else if info == "needsrepair" {
			info = "плохое"
		} else if info == "buildersrepair" {
			info = "нормальное"
		}
		if info != "" {
			it.SetValue(content, "quality", info)
		}
	}
	if info := jone.CalculateElmString(elm, "objectid/definition"); info != "" {
		it.SetValue(content, "description", info)
	}
	if info := define.GetString("rooms"); info != "" {
		if info == "room" {
			info = "1"
		} else if info == "guest" {
			info = "1"
		} else if info == "malo" {
			info = "1"
		} else if info == "1e" {
			it.SetValue(content, "studio", "1")
		}
		if ri := lik.StrToInt(info); ri > 0 {
			it.SetValue(content, "rooms", ri)
			it.SetValue(content, "rooms-offered", ri)
		}
	}
	if info := jone.CalculateTranslate(define,"floor"); info != "" {
		it.SetValue(content, "floor", info)
	}
	if info := jone.CalculateTranslate(define, "balcony"); info != "" {
		it.SetValue(content, "balcony", info)
	}
	if info := jone.CalculateTranslate(define,"restroom"); info != "" {
		it.SetValue(content, "bathroom-unit", info)
	}
	if info := jone.CalculateTranslate(define,"floortotal"); info != "" {
		it.SetValue(content, "floors-total", info)
	}
	if info := jone.CalculateTranslate(define,"housetype"); info != "" {
		it.SetValue(content, "building-type", info)
	}
	if info := jone.CalculateTranslate(define,"ceilingheight"); info != "" {
		it.SetValue(content, "ceiling-height", info)
	}
	if info := define.GetString("elevator"); info != "" {
		if info == "one" {
			info = "1"
		} else if info == "two" {
			info = "1"
		} else {
			info = ""
		}
		if info != "" {
			it.SetValue(content, "lift", info)
		}
	}
	if info := define.GetString("electricity"); info != "" {
		if strings.HasPrefix(info, "no") {
			info = ""
		} else {
			info = "1"
		}
		if info != "" {
			it.SetValue(content, "electricity-supply", info)
		}
	}
	if info := define.GetString("watersupply"); info != "" {
		if strings.HasPrefix(info, "no") {
			info = ""
		} else {
			info = "1"
		}
		if info != "" {
			it.SetValue(content, "water-supply", info)
		}
	}
	if info := define.GetString("gas"); info != "" {
		if strings.HasPrefix(info, "no") {
			info = ""
		} else {
			info = "1"
		}
		if info != "" {
			it.SetValue(content, "gas-supply", info)
		}
	}
	if info := define.GetString("sewage"); info != "" {
		if strings.HasPrefix(info, "no") {
			info = ""
		} else {
			info = "1"
		}
		if info != "" {
			it.SetValue(content, "sewerage-supply", info)
		}
	}
	if info := define.GetString("heatsystem"); info != "" {
		if strings.HasPrefix(info, "no") {
			info = ""
		} else {
			info = "1"
		}
		if info != "" {
			it.SetValue(content, "heating-supply", info)
		}
	}
	if info := define.GetString("heatsystem"); info != "" {
		if strings.HasPrefix(info, "no") {
			info = ""
		} else {
			info = "1"
		}
		if info != "" {
			it.SetValue(content, "toilet", info)
		}
	}
	return ""
}

//	Добавление фотографИЙ
func (it *TaskYandex) appendOfferPhotos(elm *likbase.ItElm, content lik.Lister) string {
	if lpic := jone.CalculateElmList(elm,"objectid/picture"); lpic != nil {
		nph := 0
		for np := 0; np < lpic.Count(); np++ {
			if pic := lpic.GetSet(np); pic != nil &&
				pic.GetString("media") == "photo" &&
				pic.GetString("album") != "nn" &&
				pic.GetString("promot") != "nn" {
				if url := pic.GetString("url"); url != "" {
					nph++
					photo := it.SetValue(content, "image", show.UrlToString(url))
					photo.SetItem(nph, "order")
				}
			}
		}
	}
	return ""
}

