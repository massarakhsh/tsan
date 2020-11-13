//	Модуль задачи сопровождения площадки
package taskexport

import (
	"bitbucket.org/961961/tsan/jone"
	"bitbucket.org/961961/tsan/repo"
	"bitbucket.org/961961/tsan/routine"
	"bitbucket.org/shaman/lik"
	"bitbucket.org/shaman/lik/likbase"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

//	Дескриптор задачи экспорта
type TaskExport struct {
	routine.Task			//	Общая задача
	Self	TaskExpoter		//	Собственный интерфейс
	Pause 	time.Duration	//	Объявленная пауза
	AnswerAt int			//	Unix время ответа
	Set		lik.Seter		//	Структура настроек
	Offers	lik.Lister		//	Список заявок
	Report	lik.Lister		//	Накопленный отчет
}

//	Интерфейс задачи экспорта
type TaskExpoter interface {
	routine.Tasker
	PrepareReport() bool					//	Подготовка отчета
	RequestReport() lik.Seter				//	Запрос отчета
	AppendOffer(elm *likbase.ItElm) string	//	ДОбавление заявки
	WriteFiles()							//	Запись файла
}

//	Инициализация задачи
func (it *TaskExport) Initialize(name string, self TaskExpoter) {
	it.Name = name
	it.Self = self
	routine.RegisterTask(self)
	go it.run()
}

//	Основной процесс
func (it *TaskExport) run() {
	timeNextRequest := time.Now().Add(time.Second * 5)
	for !it.IsStoping() {
		it.seekSet()
		if time.Now().After(timeNextRequest) {
			it.Pause = time.Minute * 5
			go it.generateFile()
			go it.Self.RequestReport()
			time.Sleep(time.Second * 10)
			timeNextRequest = time.Now().Add(it.Pause)
		}
		time.Sleep(time.Second * 1)
	}
	it.OnStoped()
}

//	Генерация файла
func (it *TaskExport) generateFile() {
	if it.Self.PrepareReport() {
		path := "export/" + it.Name
		for _, elm := range jone.TableOffer.Elms {
			if jone.CalculateElmString(elm, "status") == jone.ItActive &&
				jone.CalculateElmBool(elm, "export/ready") &&
				jone.CalculateElmBool(elm, "export/confirm") &&
				jone.CalculateElmBool(elm, "export/enable") &&
				jone.CalculateElmBool(elm, path + "/ready") &&
				jone.CalculateElmBool(elm, path + "/use") {
				if diag := it.Self.AppendOffer(elm); diag != "" {
					jone.SetElmValue(elm, diag, path + "/diagnosis")
					jone.SetElmValue(elm, false, path + "/ready")
					jone.SetElmValue(elm, false, path + "/use")
				} else {
					//jone.SetElmValue(elm, nil, path + "/diagnosis")
				}
			}
		}
		it.Self.WriteFiles()
	}
}

//	Поиск дескриптора в настройках
func (it *TaskExport) seekSet() {
	if it.Set == nil {
		if ent := repo.GenExtern.FindEnt("export"); ent != nil {
			if content := ent.It.GetList("content"); content != nil {
				for nc := 0; nc < content.Count(); nc++ {
					if set := content.GetSet(nc); set != nil {
						if set.GetString("part") == it.Name {
							it.Set = set
							break
						}
					}
				}
			}
		}
	}
}

//	Добавление содержимого элемента в отчёт XML
//	set - элемент
//	результат - пустой список
func (it *TaskExport) AddContent(set lik.Seter) lik.Lister {
	content := lik.BuildList()
	if set != nil {
		set.SetItem(content, "_content")
	}
	return content
}

//	Добавление элемента с содержимым
//	set - элемент
//	результат - элемент и пустой список
func (it *TaskExport) AddSetContent(list lik.Lister, tag string) (lik.Seter,lik.Lister) {
	item := lik.BuildSet()
	if tag != "" {
		item.SetItem(tag, "_tag")
	}
	if list != nil {
		list.AddItems(item)
	}
	return item, it.AddContent(item)
}

//	Установка значения XML
func (it *TaskExport) SetValue(content lik.Lister, tag string, value interface{}) lik.Seter {
	item := content.AddItemSet("_tag", tag)
	if value != "" {
		item.SetItem(value, "_value")
	}
	return item
}

//	Собрать адрес
func (it *TaskExport) MakeAddress(elm *likbase.ItElm) string {
	text := jone.MakeAddress(jone.CalculateElmSet(elm, "objectid/address"))
	return text
}

//	Собрать координаты
func (it *TaskExport) MakePoint(elm *likbase.ItElm) (float64,float64) {
	cx, cy := jone.MakePoint(jone.CalculateElmList(elm, "objectid/map/points"))
	return cx, cy
}

//	Преобразование в формат CDATA
func (it *TaskExport) TextToData(text string) string {
	data := text
	data = strings.Replace(data, "\r", "", -1)
	data = strings.Replace(data, "\n", "<br/>", -1)
	return "<![CDATA[" + data + "]]>"
}

//	Экспорт отчёта в файл XML
//	Для моделирования XML - файла динамическими структурами библиотеки lik
//	приняты следующие соглашения:
//	Списки XML представлены списками lik.Lister
//	Элементы XML представлены структурами lik.Seter
//	Тег элемента представлен полем _tag
//	Значение элемента представлено полем _value
//	Содержимое элемента представлено полем _content
//	Поля, не начинающиеся с _, представляют параметры тега
//	Заголовок XML представлен полем _firstline
//	report - отчёт
func (it *TaskExport) WriteFileXML(report lik.Lister) {
	if dump := lik.XML_ListToString("", report); dump != "" {
		if path := repo.GetExternPath(it.Set.GetString("part")); path != "" {
			it.WriteFileString(path, dump)
		}
	}
}

//	Экспорт отчёта в файл CSV
func (it *TaskExport) WriteFileCSV(report lik.Lister, dlm string) {
	if dump := it.DumpCsv(report, dlm); dump != "" {
		if path := repo.GetExternPath(it.Set.GetString("part")); path != "" {
			it.WriteFileString(path, dump)
		}
	}
}

//	Запись файла
func (it *TaskExport) WriteFileString(path string, dump string) {
	if match := lik.RegExParse(path, "(.+)/[^/]*$"); match != nil {
		os.MkdirAll(match[1], os.ModePerm)
	}
	btdump := []byte(dump)
	_ = ioutil.WriteFile(path, btdump, 0666)
}

//	Дамп отчёта формата CSV в строку
//	Отчёт в формате CSV представляет собой простой список элементов,
//	каждый из который представляет список полей для вывода в одну строку.
//	report - отчёт
func (it *TaskExport) DumpCsv(list lik.Lister, dlm string) string {
	dump := list.ToCsv(dlm)
	return dump
}

