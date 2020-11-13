package repo

import (
	"github.com/massarakhsh/tsan/jone"
	"github.com/massarakhsh/lik"
	"github.com/massarakhsh/lik/likbase"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)

//	Дескриптор загрузки заявок
type Offers struct {
	XMLName xml.Name    `xml:"realty-feed"`
	Offers  []LoadOffer `xml:"offer"`
}

//	Схема загрузки XML - заявки из CRM
type LoadOffer struct {
	XMLName      xml.Name     `xml:"offer"`
	InternalId   string       `xml:"internal-id,attr"`
	Type         string       `xml:"property-type"`
	Category     string       `xml:"category"`
	CreationDate string       `xml:"creation-date"`
	Location     LoadLocation `xml:"location"`
	Promo        string       `xml:"promo"`
	Rooms        string       `xml:"rooms"`
	Image        []string  	  `xml:"image"`
	Price		PriceBase	`xml:"price"`
	Area		ValueBase	`xml:"area"`
	Living		ValueBase	`xml:"living-space"`
	KitchenSpace ValueBase	`xml:"kitchen-space"`
	RoomSpace	[]ValueBase	`xml:"room-space"`
	Floor        string       `xml:"floor"`
	Floors       string       `xml:"floors-total"`
}

//	Схема загрузки места
type LoadLocation struct {
	XMLName    xml.Name `xml:"location"`
	Country		string	   `xml:"country"`
	Region		string	   `xml:"region"`
	LocalityName		string	   `xml:"locality-name"`
	Address		string	   `xml:"address"`
	Latitude		string	   `xml:"latitude"`
	Longitude		string	   `xml:"longitude"`
}

//	Схема цены
type PriceBase struct {
	Value		string	   `xml:"value"`
	Currency	string	   `xml:"currency"`
	Unit		string	   `xml:"unit"`
}

//	Схема значения
type ValueBase struct {
	Value		string	   `xml:"value"`
	Unit		string	   `xml:"unit"`
}

//	Загрузка файла XML
func LoadDump() {
	sid_objects := make(map[string]*likbase.ItElm)
	for _,obj := range(jone.TableObject.Elms) {
		if sid := jone.CalculateElmString(obj,"ids"); sid != "" {
			sid_objects[sid] = obj
		}
	}

	xmlFile, err := os.Open("topnlab.ru.xml")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer xmlFile.Close()
	byteValue, _ := ioutil.ReadAll(xmlFile)
	fmt.Printf("Size: %d\n", len(byteValue))

	var dump Offers
	xml.Unmarshal(byteValue, &dump)

	for i := 0; i < len(dump.Offers) /*&& i < 5*/; i++ {
		dumpelm := dump.Offers[i]
		ids := dumpelm.InternalId
		fmt.Println("Id: " + ids)
		var object *likbase.ItElm
		if obj := sid_objects[ids]; obj != nil {
			object = obj
		}
		if object == nil {
			object = jone.TableObject.CreateElm()
		}
		jone.SetElmValue(object, ids,"ids")
		jone.SetElmValue(object, dumpelm.Location.Region,"address/region")
		jone.SetElmValue(object, dumpelm.Location.LocalityName,"address/city")
		jone.SetElmValue(object, dumpelm.Location.Address,"address/street")
		jone.SetElmValue(object, dumpelm.Location.Latitude + " " + dumpelm.Location.Longitude,"point")
		jone.SetElmValue(object, dumpelm.Rooms,"define/rooms")
		jone.SetElmValue(object, dumpelm.Area.Value,"define/square")
		jone.SetElmValue(object, dumpelm.Living.Value,"define/squareliving")
		jone.SetElmValue(object, dumpelm.KitchenSpace.Value,"define/squarekitchen")
		jone.SetElmValue(object, dumpelm.Floor,"define/floor")
		jone.SetElmValue(object, dumpelm.Floors,"define/floortotal")
		if dumpelm.RoomSpace != nil && len(dumpelm.RoomSpace) > 0 {
			rooms := ""
			for _,rm := range(dumpelm.RoomSpace) {
				if rm.Value != "" {
					if rooms != "" { rooms += ", " }
					rooms += rm.Value
				}
			}
			jone.SetElmValue(object, rooms,"define/squarerooms")
		}
		if dumpelm.Image != nil && len(dumpelm.Image) > 0 {
			pictures := lik.BuildList()
			for _,img := range(dumpelm.Image) {
				if img != "" {
					pict := lik.BuildSet()
					pict.SetItem(img, "url")
					pictures.AddItems(pict)
				}
			}
			jone.SetElmValue(object, pictures, "picture")
		}
		id := object.Id
		offer := jone.TableOffer.CreateElm()
		jone.SetElmValue(offer, fmt.Sprint(int(id)),"objectid")
		price := dumpelm.Price.Value
		price += dumpelm.Price.Unit
		if dumpelm.Price.Currency != "" {
			price += "/" + dumpelm.Price.Currency
		}
		jone.SetElmValue(offer, price,"price")
	}
}

