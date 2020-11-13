package jone_test

import (
	"bitbucket.org/961961/tsan/jone"
	"bitbucket.org/shaman/lik"
	"bitbucket.org/shaman/lik/likbase"
)

// Чтение данных
func Example_read() {
	var id lik.IDB
	// Из таблицы заявок получить заявку по ее ID
	offer := jone.TableOffer.GetElm(id)
	// Получить цену заявки как целое число
	cost := jone.CalculateElmInt(offer, "cost")
	// Получить идентификатор ответственного риэлтора как целое число
	memberid := jone.CalculateElmInt(offer, "memberid")
	// Получить этаж из описания объекта как целое число
	flat := jone.CalculateElmInt(offer, "objectid/define/flat")
	// Получить фамилию клиента заявки
	family := jone.CalculateElmInt(offer, "clientid/family")
	_, _, _, _ = cost, memberid, flat, family
}

// Запись данных
func Example_write() {
	var offer *likbase.ItElm
	jone.SetElmValue(offer, 2300000, "cost")
	// Установить цену заявки
}

