package avito

import (
	"bitbucket.org/961961/tsan/jone"
	"bitbucket.org/961961/tsan/show"
	"bitbucket.org/shaman/lik"
	"bitbucket.org/shaman/lik/likbase"
	"time"
)

//	Подготовка отчета
func (it *TaskAvito) PrepareReport() bool {
	result := false
	if it.Set != nil {
		it.Report = it.AddContent(nil)
		it.Report.AddItems("<?xml version=\"1.0\" encoding=\"UTF-8\" standalone=\"yes\"?>")
		rf, offers := it.AddSetContent(it.Report, "Ads")
		rf.SetItem("Avito.ru", "target")
		rf.SetItem("3", "formatVersion")
		it.Offers = offers
		it.Offers.AddItemSet("_tag=generation-date", "_value", show.TimeToString(int(time.Now().Unix())))
		result = true
	}
	return result
}

//	Запись файла
func (it *TaskAvito) WriteFiles() {
	it.WriteFileXML(it.Report)
}

//	Добавление в отчёт Авито заявки
//	elm - добавляемая заявка
func (it *TaskAvito) AppendOffer(elm *likbase.ItElm) string {
	card,content := it.AddSetContent(nil, "Ad")
	idu := elm.GetIDB("idu")
	if idu == 0 { idu = elm.Id }
	it.SetValue(content, "Id", idu)
	if diag := it.appendOfferMain(elm, content); diag != "" {
		return diag
	}
	if diag := it.appendOfferSaler(elm, content); diag != "" {
		return diag
	}
	if diag := it.appendOfferAddress(elm, content); diag != "" {
		return diag
	}
	if diag := it.appendOfferDefine(elm, content); diag != "" {
		return diag
	}
	if diag := it.appendOfferParams(elm, content); diag != "" {
		return diag
	}
	it.Offers.AddItems(card)
	return ""
}

//	Основные сведения по заявке
func (it *TaskAvito) appendOfferMain(elm *likbase.ItElm, content lik.Lister) string {
	return ""
}

//	Сведения о продавце
func (it *TaskAvito) appendOfferSaler(elm *likbase.ItElm, content lik.Lister) string {
	member := jone.TableMember.GetElm(jone.CalculateElmIDB(elm, "memberid"))
	if member == nil {
		return "Отсутствует ответственный риэлтор"
	}
	info := jone.CalculateElmText(member)
	it.SetValue(content, "ManagerName", info)
	if info := member.GetString("prophone"); info != "" {
		it.SetValue(content, "ContactPhone", show.PhoneToString(info))
	} else if info := member.GetString("phone"); info != "" {
		it.SetValue(content, "ContactPhone", show.PhoneToString(info))
	} else {
		return "Отсутствует телефон риэлтора"
	}
	return ""
}

//	Сведения о местоположении
func (it *TaskAvito) appendOfferAddress(elm *likbase.ItElm, content lik.Lister) string {
	if text := it.MakeAddress(elm); text != "" {
		it.SetValue(content, "Address", text)
	} else {
		return "Отсутсвует адрес объекта"
	}
	return ""
}

//	Описание объекта
func (it *TaskAvito) appendOfferDefine(elm *likbase.ItElm, content lik.Lister) string {
	if info := jone.CalculateElmString(elm, "objectid/definition"); info != "" {
		text := it.TextToData(info)
		it.SetValue(content, "Description", text)
	} else {
		return "Отсутсвует описание объекта"
	}
	if cx, cy := it.MakePoint(elm); cx != 0 && cy != 0 {
		it.SetValue(content, "Latitude", cx)
		it.SetValue(content, "Longitude", cy)
	}
	return ""
}

//	Параметры объекта
func (it *TaskAvito) appendOfferParams(elm *likbase.ItElm, content lik.Lister) string {
	segment := jone.CalculateElmString(elm, "segment")
	target := jone.CalculateElmString(elm, "target")
	if segment == jone.DoRent && target == "sale" {
		it.SetValue(content, "OperationType", "Сдам")
	} else if target == "sale" {
		it.SetValue(content, "OperationType", "Продам")
	} else {
		return "Неверная цель заявки " + target
	}
	realty := jone.CalculateElmString(elm, "objectid/realty")
	if realty == "flat" {
		it.SetValue(content, "Category", "Квартиры")
		it.SetValue(content, "Status", "Квартира")
	} else if realty == "room" {
		it.SetValue(content, "Category", "Комнаты")
	} else if realty == "house" {
		it.SetValue(content, "Category", "Дома, дачи, коттеджи")
	} else if realty == "lot" {
		it.SetValue(content, "Category", "Земельные участки")
	} else {
		return "Неверный тип недвижимости " + realty
	}
	if info := jone.CalculateElmInt(elm, "cost"); info > 0 {
		it.SetValue(content, "Price", info)
	} else {
		return "Не указана цена объекта"
	}
	if info := jone.CalculateElmString(elm, "objectid/define/rooms"); info != "" {
		if info == "room" {
			info = "1"
		} else if info == "guest" {
			info = "1"
		} else if info == "malo" {
			info = "1"
		} else if info == "1e" {
			info = "Студия"
		} else if _,ok := lik.StrToIntIf(info); !ok {
			info = ""
		}
		if info != "" {
			it.SetValue(content, "Rooms", info)
		}
	}
	if info := jone.CalculateElmString(elm, "objectid/define/square"); info != "" {
		it.SetValue(content, "Square", info)
	}
	if info := jone.CalculateElmString(elm, "objectid/define/landarea"); info != "" {
		it.SetValue(content, "LandArea", info)
	}
	if info := jone.CalculateElmString(elm, "objectid/define/squarekitchen"); info != "" {
		it.SetValue(content, "KitchenSpace", info)
	}
	if info := jone.CalculateElmString(elm, "objectid/define/squareliving"); info != "" {
		it.SetValue(content, "LivingSpace", info)
	}
	if info := jone.CalculateElmInt(elm, "objectid/define/floor"); info > 0 {
		it.SetValue(content, "Floor", info)
	}
	if info := jone.CalculateElmInt(elm, "objectid/define/floortotal"); info > 0 {
		it.SetValue(content, "Floors", info)
	}
	if info := jone.CalculateElmString(elm, "objectid/define/wallmaterial"); info != "" {
		if realty == "flat" || realty == "room" {
			if info == "bricks" {
				it.SetValue(content, "HouseType", "Кирпичный")
			} else if info == "bricksmonolith" {
				it.SetValue(content, "HouseType", "Кирпичный")
			} else if info == "panel" {
				it.SetValue(content, "HouseType", "Панельный")
			} else if info == "monolith" {
				it.SetValue(content, "HouseType", "Монолитный")
			} else if info == "blocky" {
				it.SetValue(content, "HouseType", "Блочный")
			} else if info == "wood" {
				it.SetValue(content, "HouseType", "Деревянный")
			} else {
				it.SetValue(content, "HouseType", jone.SystemStringTranslate("wallmaterial", info))
			}
		} else if realty == "house" || realty == "land" {
			if info == "bricks" {
				it.SetValue(content, "WallsType", "Кирпич")
			} else {
				it.SetValue(content, "WallsType", jone.SystemStringTranslate("wallmaterial", info))
			}
		}
	}
	if info := jone.CalculateElmString(elm, "segment"); info == "second" {
		it.SetValue(content, "MarketType", "Вторичка")
	} else if info == "new" {
		it.SetValue(content, "MarketType", "Новостройка")
	}
	it.SetValue(content, "PropertyRights", "Посредник")
	it.appendOfferPhotos(elm, content)
	return ""
}

//	Добавление фотографий
func (it *TaskAvito) appendOfferPhotos(elm *likbase.ItElm, content lik.Lister) string {
	if lpic := jone.CalculateElmList(elm,"objectid/picture"); lpic != nil {
		_,images := it.AddSetContent(content, "Images")
		nph := 0
		for np := 0; np < lpic.Count(); np++ {
			if pic := lpic.GetSet(np); pic != nil &&
				pic.GetString("media") == "photo" &&
				pic.GetString("album") != "nn" &&
				pic.GetString("promot") != "nn" {
				if url := pic.GetString("url"); url != "" {
					nph++
					photo := it.SetValue(images, "Image", "")
					photo.SetItem(show.UrlToString(url), "url")
				}
			}
		}
	}
	return ""
}

