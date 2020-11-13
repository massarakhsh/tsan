package avito

import (
	"github.com/massarakhsh/tsan/jone"
	"github.com/massarakhsh/tsan/one"
	"github.com/massarakhsh/lik"
	"github.com/massarakhsh/lik/likapi"
	"github.com/massarakhsh/lik/likbase"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

//	Запрос отчета об экспорте
func (it *TaskAvito) RequestReport() lik.Seter {
	var answer lik.Seter
	if one.AvitoToken == "" {
		one.AvitoToken = it.getAvitoToken()
	}
	//rid := GetAvitoLastId()
	//if rid == "" {
	//	return nil
	//}
	url := fmt.Sprintf("https://api.avito.ru/autoload/v1/accounts/%s/reports/last_report", one.AvitoId)
	headers := lik.BuildSet("Authorization", "Bearer " + one.AvitoToken)
	if answer = likapi.GetHttpRequest(url, headers); answer == nil {
		it.Pause = time.Minute * 5
	} else if answer.IsItem("error") {
		one.AvitoToken = ""
		it.Pause = time.Minute * 5
	} else if status := answer.GetString("report/status"); strings.Contains(status, "noprocess") {
		it.Pause = time.Minute * 5
	} else if ads := answer.GetList("report/ads"); ads != nil {
		it.AnswerAt = int(time.Now().Unix())
		fin := answer.GetString("report/finished_at")
		if at,err := time.Parse(time.RFC3339, fin); err == nil {
			it.AnswerAt = int(at.Unix())
		}
		for nad := 0; nad < ads.Count(); nad++ {
			if ad := ads.GetSet(nad); ad != nil {
				if id := ad.GetIDB("ad_id"); id > 0 {
					if elm := jone.TableOffer.GetElm(id); elm != nil {
						it.Pause = time.Minute * 15
						it.receiveAnswerOffer(ad, elm)
					}
				}
			}
		}
	}
	return answer
}

//	Обработка ответа с отчетом
func (it *TaskAvito) receiveAnswerOffer(ad lik.Seter, elm *likbase.ItElm) {
	pot := elm.GetSet("export/avito")
	if pot == nil {
		pot = lik.BuildSet()
		jone.SetElmValue(elm, pot, "export/avito")
	}
	if url := ad.GetString("url"); url != "" {
		if old := pot.GetString("url"); old != url {
			jone.SetElmValue(elm, url, "export/avito/url")
		}
		it.storeAnswerLink(ad, elm, url)
	}
	if status := ad.GetString("statuses/general/value"); status == "success" {
		jone.SetElmValue(elm, nil, "export/avito/diagnosis")
	} else if help := ad.GetString("statuses/general/help"); help != "" {
		jone.SetElmValue(elm, help, "export/avito/diagnosis")
	}
	if messages := ad.GetList("messages"); true {
		it.storeAnswerMessages(ad, elm, messages)
	}
}

//	Запомнить ссылки из отчета
func (it *TaskAvito) storeAnswerLink(ad lik.Seter, elm *likbase.ItElm, url string) {
	found := false
	lpic := jone.CalculateElmList(elm,"objectid/picture")
	if lpic == nil {
		lpic = lik.BuildList()
		jone.SetElmValue(elm, lpic, "objectid/picture")
	} else {
		for np := 0; np < lpic.Count(); np++ {
			if link := lpic.GetSet(np); link != nil && link.GetString("media") == "link" {
				if old := link.GetString("url"); strings.Contains(old, "avito") {
					found = true
					if url != old {
						link.SetItem(url, "url")
						elm.OnModify()
					}
					break
				}
			}
		}
	}
	if !found {
		link := lpic.AddItemSet("id", fmt.Sprintf("%d", 100000000 + rand.Intn(900000000)))
		link.SetItem("link", "media")
		link.SetItem("Авито", "comment")
		link.SetItem(url, "url")
		elm.OnModify()
	}
}

//	Запомнить сообщения из отчета
func (it *TaskAvito) storeAnswerMessages(ad lik.Seter, elm *likbase.ItElm, mess lik.Lister) {
	query := "proto='public' AND scope='avito' AND offer_id=? AND deleted_at IS NULL"
	oldmess := one.SelectMessage(query, int(elm.Id))
	if mess != nil {
		for nm := 0; nm < mess.Count(); nm++ {
			if msg := mess.GetSet(nm); msg != nil {
				if body := msg.GetString("description"); body != "" && body != "null"{
					found := false
					for no := 0; no < len(oldmess); no++ {
						if oldmess[no].Body == body {
							found = true
							oldmess[no].Body = ""
							break
						}
					}
					if !found {
						message := &one.Message{Proto: "public", Scope: "avito"}
						message.OfferId = elm.Id
						message.Body = body
						message.TimeAt = it.AnswerAt
						message.Save()
					}
				}
			}
		}
	}
	for no := 0; no < len(oldmess); no++ {
		if oldmess[no].Body != "" {
			oldmess[no].Delete()
		}
	}
}

//	Получить новый токен
func (it *TaskAvito) getAvitoToken() string {
	token := ""
	url := "https://api.avito.ru/token/?grant_type=client_credentials"
	url += fmt.Sprintf("&client_id=%s", one.AvitoClient)
	url += fmt.Sprintf("&client_secret=%s", one.AvitoSecret)
	if answer := likapi.GetHttpRequest(url, nil); answer != nil {
		token = answer.GetString("access_token")
	}
	return token
}

//	Определить номер последнего отчета
func (it *TaskAvito) setAvitoLastId() string {
	rid := ""
	url := fmt.Sprintf("https://api.avito.ru/autoload/v1/accounts/%s/reports", one.AvitoId)
	headers := lik.BuildSet("Authorization", "Bearer " + one.AvitoToken)
	if answer := likapi.GetHttpRequest(url, headers); answer == nil {
	} else if reps := answer.GetList("reports"); reps != nil {
		for nr := 0; nr < reps.Count(); nr++ {
			if rep := reps.GetSet(nr); rep != nil {
				if rep.GetString("finished_at") != "" {
					rid = rep.GetString("id")
					break
				}
			}
		}
	}
	return rid
}

