package repo

import (
	"bitbucket.org/961961/tsan/jone"
	"bitbucket.org/shaman/lik"
	"bitbucket.org/shaman/lik/likapi"
	"bitbucket.org/shaman/lik/likbase"
	"fmt"
	"os"
	"sync"
	"time"
)

const MAX_CONTROL = 16

//	Дескриптор сессии
type DataSession struct {
	likapi.DataSession
	Sync       sync.Mutex //	Семафор
	Params     lik.Seter  //	Параметры сессии
	WaitLogin  bool       //	Признак ожидания логина
	ExitLogin  bool       //	Признак выхода из логина
	IdMember   lik.IDB    //	Идентификатор сотрудника
	RoleMember string     //	Роль сотрудника
	Collect    map[string]bool
}

//	Дескриптор страницы
type DataPage struct {
	likapi.DataPage
	Sync       sync.Mutex   //	Семафор
	Session    *DataSession //	Указатель на сессию
	Params     lik.Seter    //	Параметры страницы
	Locates    []Controller //	Контроллеры стека
	Sheets     []Controller //	Скрытые контроллеры
	PathLast   string       //	Последний путь
	PathNeed   string       //	Нужный путь
	PathClient string       //	Путь на клиенте
	Tune       struct {     //	Параметры настройки
		ItGen  string   //	Текущмй генератор
		ItKey  string   //	Текущий ключ
		List   []string //	Список чего-то
		Copied bool     //	Скопировано
		Cuted  bool     //	Вырезано
	}
	Mask struct { //	Параметры подбора
		Id    lik.IDB         //	Идентификатор заявки
		Ignor map[string]bool //	Отключенные требования
	}
}

//	Интерфейс страницы
type DataPager interface {
	likapi.DataPager
	GetItPage() *DataPage //	Получить указатель на страницу
}

//	Дескриптор запроса
type DataRule struct {
	likapi.DataDrive
	Sync         sync.Mutex   //	Семафор
	ItPage       *DataPage    //	Указатель на страницу
	ItSession    *DataSession //	Указатель на сессию
	Title        string       //	Заголовок страницы
	PathCommand  string       //	Команда на пути
	ResultFormat bool         //	Признак форматирования
	IsJson       bool         //	Это JSON
	IsXml        bool         //	Это XML
	IsCsv        bool         //	Это CSV
	IsTechno     bool         //	Технологический режим
	IsChangePage bool         //	Изменена страницыа
	IsChangeData bool         //	Изменены данные
	IsNeedReload bool         //	Необходима перезагрузка страниц
	IsNeedGoPath bool         //	Необходимо изменить путь
}

//	Интерфейс запроса
type DataRuler interface {
	likapi.DataDriver
	Page() *DataPage
}

var (
	HostPid    int    = 0               //	Идентификатор процесса
	HostPort   int    = 80              //	Порт хоста
	HostServ   string = "localhost" //	Адрес сервера базы
	HostBase   string = "tsan"          //	Имя базы
	HostUser   string = "tsan"          //	Логин базы
	HostPass   string = "tsan"          //	Пароль базы
	HostSignal string = ""              //	Посылаемый сигнал

	TimeStart   time.Time      //	Время старта программы
	TimeBin     time.Time      //	Время сборки файла
	ToPause     bool           //	Признак перехода в паузу
	ToTerminate bool           //	Признак перехода в останов
	ChanSignal  chan os.Signal //	Канал приема сигналов
)

//	Запуск страницы
func StartPage(uri string) *DataPage {
	session := &DataSession{}
	session.Uri = uri
	session.Params = lik.BuildSet()
	session.Collect = make(map[string]bool)
	page := &DataPage{Session: session}
	page.Params = lik.BuildSet()
	page.Self = page
	session.StartToPage(page)
	return page
}

//	Клонирование страницы
func ClonePage(from *DataPage) *DataPage {
	page := &DataPage{Session: from.Session}
	page.Self = page
	page.Params = from.Params.Clone().ToSet()
	page.PathNeed = from.PathNeed
	from.ContinueToPage(page)
	return page
}

//	Подключение запроса к страницу
func BindRule(page *DataPage) *DataRule {
	rule := &DataRule{ItPage: page, ItSession: page.Session}
	rule.Page = page
	if ToPause && !rule.IAmShaman() {
		rule.IsTechno = true
	}
	return rule
}

//	Получить указатель на страницу
func (page *DataPage) GetItPage() *DataPage {
	return page
}

//	Собрать идентификатор объекта
func (page *DataPage) stringCollect(part string, id lik.IDB) string {
	return fmt.Sprintf("%s%d", part, id)
}

//	Добавить объект в коллекцию
func (page *DataPage) AppendCollect(part string, id lik.IDB) {
	page.Session.Collect[page.stringCollect(part, id)] = true
}

//	Удалить объект из коллекции
func (page *DataPage) RemoveCollect(part string, id lik.IDB) {
	delete(page.Session.Collect, fmt.Sprintf("%s%d", part, id))
}

//	Проверка наличия в коллекции
func (page *DataPage) ProbeCollect(part string, id lik.IDB) bool {
	probe, _ := page.Session.Collect[fmt.Sprintf("%s%d", part, id)]
	return probe
}

//	Посчитать размер коллекции
func (page *DataPage) CountCollect(part string) int {
	count := 0
	for key, ok := range page.Session.Collect {
		if ok && lik.RegExCompare(key, "^"+part) {
			count++
		}
	}
	return count
}

//	Очистка сессии
func (rule *DataRule) ClearSession() {
	rule.ItSession.Params = lik.BuildSet()
	rule.ItPage.Params = lik.BuildSet()
	rule.ItPage.Locates = []Controller{}
	rule.ItPage.Sheets = []Controller{}
	rule.ItPage.Session.Collect = make(map[string]bool)
}

//	Проверка, что страница изменена
func (rule *DataRule) IsItChangePage() bool {
	return rule.IsChangePage
}

//	Установить изменение страницы
func (rule *DataRule) OnChangePage() {
	rule.IsChangePage = true
}

//	Стереть изменение страницы
func (rule *DataRule) OffChangePage() {
	rule.IsChangePage = false
}

//	Проверить, что данные изменены
func (rule *DataRule) IsItChangeData() bool {
	return rule.IsChangeData
}

//	Установить изменение данных
func (rule *DataRule) OnChangeData() {
	rule.IsChangeData = true
}

//	Стереть изменение данных
func (rule *DataRule) OffChangeData() {
	rule.IsChangeData = false
}

//	Получить дескриптор оператора
func (rule *DataRule) GetMember() *likbase.ItElm {
	var elm *likbase.ItElm
	if rule.ItSession.IdMember > 0 {
		elm = jone.TableMember.GetElm(rule.ItSession.IdMember)
	}
	return elm
}

//	Установить идентификатор текущего объекта
func (rule *DataRule) SetLocateId(id lik.IDB) {
	if lev := len(rule.ItPage.Locates); lev > 0 {
		rule.ItPage.Locates[lev-1].SetId(id)
		rule.ItPage.PathNeed = rule.BuildFullPath("")
		rule.ItPage.PathLast = rule.ItPage.PathNeed
	}
}

//	Занести путь в стек
func (rule *DataRule) SetPagePush(part string) {
	rule.ItPage.PathNeed = rule.BuildFullPath("") + "/" + part
}

//	Снять путь из стека
func (rule *DataRule) SetPageExit(level int) {
	rule.ItPage.PathNeed = rule.BuildLevelPath("", level)
}

//	Установить путь на стеке
func (rule *DataRule) SetPagePart(level int, part string) {
	rule.ItPage.PathNeed = rule.BuildLevelPath("", level) + "/" + part
}

//	ПОстроить полный путь
func (rule *DataRule) BuildFullPath(parm string) string {
	return rule.BuildLevelPath(parm, len(rule.ItPage.Locates))
}

//	Построить путь указанной длины
func (rule *DataRule) BuildLevelPath(parm string, levmax int) string {
	path := ""
	for lev := 0; lev < levmax && lev < len(rule.ItPage.Locates); lev++ {
		loc := rule.ItPage.Locates[lev]
		path += "/" + loc.GetMode()
		if id := loc.GetId(); id != 0 {
			path += fmt.Sprintf("%d", int(id))
		}
	}
	if parm != "" {
		path += "?" + parm
	}
	return path
}

//	Определить глубину стека
func (rule *DataRule) GetLevel() int {
	return len(rule.ItPage.Locates)
}

//	Зарегистрировать контроллер
func (rule *DataRule) RegControl(loc Controller) {
	rule.ItPage.Sync.Lock()
	if loc != nil && loc.IsStable() {
		size := len(rule.ItPage.Sheets)
		pos := 0
		for pos < size {
			if loc == rule.ItPage.Sheets[pos] {
				break
			}
			pos++
		}
		if pos >= size {
			rule.ItPage.Sheets = append(rule.ItPage.Sheets, loc)
		}
		if pos > 0 {
			for p := pos; p > 0; p-- {
				rule.ItPage.Sheets[p] = rule.ItPage.Sheets[p-1]
			}
			rule.ItPage.Sheets[0] = loc
		}
	}
	rule.ItPage.Sync.Unlock()
}

//	Определить сегмент
func (rule *DataRule) SeekSegment(segment string) string {
	if operator := rule.GetMember(); operator == nil {
		segment = ""
	} else {
		if segment == "" {
			segment = rule.ItPage.Params.GetString("segment")
		}
		if segment == "" {
			segment = rule.GetMemberParamString("segment")
		}
		if segment == jone.DoCall {
		} else if rule.IAmAdmin() {
		} else if segment == "system" {
			segment = ""
		} else if !operator.GetBool("do_" + segment) {
			segment = ""
		}
		if segment == "" {
			if operator.GetString("role") == jone.ItDispatch {
				segment = jone.DoCall
			} else if operator.GetBool("do_" + jone.DoSecond) {
				segment = jone.DoSecond
			} else if operator.GetBool("do_" + jone.DoNew) {
				segment = jone.DoNew
			} else if operator.GetBool("do_" + jone.DoVilla) {
				segment = jone.DoVilla
			} else if operator.GetBool("do_" + jone.DoRent) {
				segment = jone.DoRent
			} else if operator.GetBool("do_" + jone.DoArea) {
				segment = jone.DoArea
			}
		}
		rule.SetMemberParam(segment, "segment")
	}
	rule.ItPage.Params.SetItem(segment, "segment")
	return segment
}

//	Записать в лог ошибку
func (rule *DataRule) SayError(text string) {
	lik.SayError(rule.GetIP() + ": " + text)
}

//	Записать в лог предупреждение
func (rule *DataRule) SayWarning(text string) {
	lik.SayWarning(rule.GetIP() + ": " + text)
}

//	Записать в лог информацию
func (rule *DataRule) SayInfo(text string) {
	lik.SayInfo(rule.GetIP() + ": " + text)
}
