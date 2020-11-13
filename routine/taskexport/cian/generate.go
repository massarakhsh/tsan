package cian

import (
	"bitbucket.org/961961/tsan/jone"
	"bitbucket.org/961961/tsan/show"
	"bitbucket.org/shaman/lik"
	"bitbucket.org/shaman/lik/likbase"
)

//	Подготовка отчета
func (it *TaskCian) PrepareReport() bool {
	it.Report = it.AddContent(nil)
	_, offers := it.AddSetContent(it.Report, "feed")
	it.SetValue(offers, "feed_version", 2)
	it.Offers = offers
	return true
}

//	Запись файла
func (it *TaskCian) WriteFiles() {
	it.WriteFileXML(it.Report)
}

//	Добавление в отчёт ЦИАН заявки
//	elm - добавляемая заявка
func (it *TaskCian) AppendOffer(elm *likbase.ItElm) string {
	card,content := it.AddSetContent(nil, "object")
	if diag := it.appendOfferMain(elm, content); diag != "" {
		return diag
	}
	if diag := it.appendOfferPhoto(elm, content); diag != "" {
		return diag
	}
	if diag := it.appendOfferBargain(elm, content); diag != "" {
		return diag
	}
	if diag := it.appendOfferBuilding(elm, content); diag != "" {
		return diag
	}
	it.Offers.AddItems(card)
	return ""
}

//	Добавление основных сведений
func (it *TaskCian) appendOfferMain(elm *likbase.ItElm, content lik.Lister) string {
	it.SetValue(content, "Category", "flatSale")
	idu := elm.GetIDB("idu")
	if idu == 0 { idu = elm.Id }
	it.SetValue(content, "ExternalId", idu)
	if info := jone.CalculateElmString(elm, "objectid/definition"); info != "" {
		it.SetValue(content, "Description", info)
	} else {
		return "Отсутствует описание объекта"
	}
	if info := jone.CalculateElmString(elm,"objectid/define/rooms"); info != "" {
		rc := 0
		if rms := lik.StrToInt(info); rms >= 1 && rms <= 5 {
			rc = rms
		} else if rms > 5 {
			rc = 6
		} else if info == "1e" {
			rc = 9
		}
		if rc > 0 {
			it.SetValue(content, "FlatRoomsCount", rc)
		} else {
			return "Отсутствует число комнат"
		}
	}
	if info := jone.CalculateElmString(elm,"objectid/define/floor"); info != "" {
		it.SetValue(content, "FloorNumber", info)
	} else {
		return "Отсутствует этаж"
	}
	if info := it.MakeAddress(elm); info != "" {
		it.SetValue(content, "Address", info)
	} else {
		return "Отсутствует адрес объекта"
	}
	if info := jone.CalculateElmString(elm,"objectid/define/square"); info != "" {
		it.SetValue(content, "TotalArea", info)
	} else {
		return "Отсутствует площадь"
	}
	phone := jone.CalculateElmString(elm,"memberid/prophone");
	if phone == "" {
		phone = jone.CalculateElmString(elm,"memberid/phone");
	}
	if phone != "" {
		_,ph := it.AddSetContent(content, "Phones")
		_,phs := it.AddSetContent(ph, "PhoneSchema")
		it.SetValue(phs, "CountryCode", "+7")
		it.SetValue(phs, "Number", phone)
	} else {
		return "Отсутствует телефон"
	}
	return ""
}

//	Информация о торге
func (it *TaskCian) appendOfferBargain(elm *likbase.ItElm, main lik.Lister) string {
	_,content := it.AddSetContent(main, "BargainTerms")
	info := jone.CalculateElmString(elm, "cost")
	if info == "" {
		return "Отсутствует цена"
	}
	it.SetValue(content, "Price", info)
	return ""
}

//	Информация о здании
func (it *TaskCian) appendOfferBuilding(elm *likbase.ItElm, main lik.Lister) string {
	_,content := it.AddSetContent(main, "Building")
	if info := jone.CalculateElmString(elm,"objectid/define/floortotal"); info != "" {
		it.SetValue(content, "FloorsCount", info)
	} else {
		return "Отсутствует этажность"
	}
	return ""
}

//	Добавление фотографий
func (it *TaskCian) appendOfferPhoto(elm *likbase.ItElm, main lik.Lister) string {
	_,content := it.AddSetContent(main, "Photos")
	if lpic := jone.CalculateElmList(elm,"objectid/picture"); lpic != nil {
		nph := 0
		for np := 0; np < lpic.Count(); np++ {
			if pic := lpic.GetSet(np); pic != nil &&
				pic.GetString("media") == "photo" &&
				pic.GetString("album") != "nn" &&
				pic.GetString("promot") != "nn" {
				if url := pic.GetString("url"); url != "" {
					_,pho := it.AddSetContent(content, "PhotoSchema")
					nph++
					it.SetValue(pho, "FullUrl", show.UrlToString(url))
					it.SetValue(pho, "IsDefault", nph == 1)
				}
			}
		}
	}
	return ""
}

