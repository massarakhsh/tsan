package repo

import (
	"bitbucket.org/961961/tsan/jone"
	"bitbucket.org/shaman/lik"
)

//	Структура сделки
type DealData struct {
	IdDeal     lik.IDB	//	Индекс сделки
	IdMain     lik.IDB	//	Индекс основной заявки
	IdPair     lik.IDB	//	Индекс другой заявки
	TargetMain string	//	Цель основной заявки (sale,buy,...)
	TargetPair string	//	Цель другой заявки
	IsFrom     bool		//	Признак, что основная заявка - продавец
	IdFrom     lik.IDB	//	Индекс заявки продавца
	IdTo       lik.IDB	//	Индекс заявки покупателя
}

//	Конструктор структуры сделки
func SeekOfferDeal(id lik.IDB) DealData {
	deal := DealData{}
	if elm := jone.TableOffer.GetElm(id); elm != nil {
		deal.IdMain = id
		deal.TargetMain = elm.GetString("target")
		if deal.TargetMain == "sale" {
			deal.IsFrom = true
			deal.IdFrom = id
			deal.TargetPair = "buy"
		} else if deal.TargetMain == "buy" {
			deal.IsFrom = false
			deal.IdTo = id
			deal.TargetPair = "sale"
		}
		for idd, elm := range jone.TableDeal.Elms {
			if deal.IsFrom && elm.GetIDB("saleid") == id ||
				!deal.IsFrom && elm.GetIDB("buyid") == id {
				deal.IdDeal = idd
				if deal.IsFrom {
					deal.IdPair = elm.GetIDB("buyid")
					deal.IdTo = deal.IdPair
				} else {
					deal.IdPair = elm.GetIDB("saleid")
					deal.IdFrom = deal.IdPair
				}
				break
			}
		}
	}
	return deal
}

