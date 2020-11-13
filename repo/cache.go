package repo

import (
	"bitbucket.org/961961/tsan/jone"
	"bitbucket.org/shaman/lik"
	"sync"
	"time"
)

const SIZE_CACHE = 15

//	Дескриптор элемента кэша
type cacheGrid struct {
	At   time.Time
	What string
	Fase int
	Who  lik.IDB
	Sign string
	Data lik.Itemer
}

//	Семафор кэша
var cacheSync  sync.Mutex

//	Коллекция кэша
var cacheGrids []*cacheGrid

//	Получить список кэшированных элементов
func (rule *DataRule) CachePartGet(part string) []lik.IDB {
	queue := []lik.IDB{}
	rule.ItSession.Sync.Lock()
	if operator := rule.GetMember(); operator != nil {
		if cache := jone.CalculateElmList(operator,"cache/"+part); cache != nil {
			for nc := 0; nc < cache.Count() && nc < SIZE_CACHE; nc++ {
				if id := cache.GetItem(nc).ToInt(); id>0 {
					queue = append(queue, lik.IDB(id))
				}
			}
		}
	}
	rule.ItSession.Sync.Unlock()
	return queue
}

//	Записать элемент в кэш
func (rule *DataRule) CachePartIdPush(part string, id lik.IDB) {
	modify := false
	member := rule.GetMember()
	rule.ItSession.Sync.Lock()
	if member != nil && part != "" && id > 0 {
		cache := jone.CalculateElmSet(member,"cache")
		if cache == nil {
			modify = true
			cache = lik.BuildSet()
			member.Info.SetItem(cache, "cache")
		}
		cachepart := cache.GetList(part)
		if cachepart == nil {
			modify = true
			cachepart = lik.BuildList()
			cache.SetItem(cachepart, part)
		}
		count := cachepart.Count()
		for nc := 0; nc < count; nc++ {
			if chid := cachepart.GetIDB(nc); chid == id {
				if nc > 0 {
					modify = true
					cachepart.DelItem(nc)
					cachepart.InsertItem(id,0)
				}
				id = 0
				break
			}
		}
		if id > 0 {
			modify = true
			cachepart.InsertItem(id,0)
			for cachepart.Count() > SIZE_CACHE {
				cachepart.DelItem(cachepart.Count() - 1)
			}
		}
	}
	rule.ItSession.Sync.Unlock()
	if modify {
		member.OnModify()
	}
}
