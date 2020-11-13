package repo

import (
	"bitbucket.org/961961/tsan/jone"
	"bitbucket.org/shaman/lik"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
)

//	ПОлучить параметр сессии
func GetSessionParam(rule *DataRule, path string) string {
	return rule.ItSession.Params.GetString(path)
}

//	Установить параметр сессии
func SetSessionParam(rule *DataRule, val interface{}, path string) {
	rule.ItSession.Params.SetItem(val, path)
}

//	Получить параметр сотрудника как строку
func (rule *DataRule) GetMemberParamString(path string) string {
	if item := rule.GetMemberParam(path); item != nil {
		return item.ToString()
	}
	return ""
}

//	Получить параметр сотрудника как целое
func (rule *DataRule) GetMemberParamInt(path string) int {
	if item := rule.GetMemberParam(path); item != nil {
		return item.ToInt()
	}
	return 0
}

//	Получить параметр сотрудника как структуру
func (rule *DataRule) GetMemberParamSet(path string) lik.Seter {
	if item := rule.GetMemberParam(path); item != nil {
		return item.ToSet()
	}
	return nil
}

//	Получить параметр сотрудника как список
func (rule *DataRule) GetMemberParamList(path string) lik.Lister {
	if item := rule.GetMemberParam(path); item != nil {
		return item.ToList()
	}
	return nil
}

//	Получить параметр сотрудника как интерфейс
func (rule *DataRule) GetMemberParam(path string) lik.Itemer {
	if operator := rule.GetMember(); operator != nil {
		return jone.CalculateElm(operator, "param/"+path)
	}
	return nil
}

//	Установить параметр сотрудника
func (rule *DataRule) SetMemberParam(val interface{}, path string) {
	rule.Sync.Lock()
	if operator := rule.GetMember(); operator != nil {
		jone.SetElmValue(operator, val, "param/"+path)
	}
	rule.Sync.Unlock()
}

//	Сохранить параметры сотрудника
func (rule *DataRule) SaveMemberParam() {
	if operator := rule.GetMember(); operator != nil {
		operator.OnModify()
	}
}

//	ПОлучить параметр системы
func GetSystemParam(path string) string {
	return ""
}

//	Установить параметр системы
func SetSystemParam(path string) {
}

//	Записать файл
func WriteFile(part string, id int, sfx string, data []byte) string {
	path := fmt.Sprintf("var/%s/%03d/%03d/%03d", part, id / 1000000, id % 1000000 / 1000, id % 1000)
	os.MkdirAll(path, os.ModePerm)
	path += fmt.Sprintf("/%09d.%s", rand.Int31n(1000000000), sfx)
	_ = ioutil.WriteFile(path, data, 0666)
	return path
}

//	Снять или установить пометку
func MarkElmSet(rule *DataRule, part string, id lik.IDB, val bool) {
	if val {
		rule.ItPage.AppendCollect(part, id)
	} else {
		rule.ItPage.RemoveCollect(part, id)
	}
}
