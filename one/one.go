package one

import (
	"bitbucket.org/shaman/lik"
	"fmt"
	"github.com/jinzhu/gorm"
	"reflect"
	"time"
)

//	Ядро объекта базы данных
type One struct {
	gorm.Model
}

//	Интерфейс объекта базы данных
type Oner interface {
	Table() string
	Save() bool
	Update(datas... interface{}) bool
	Delete()
}

//	Создать объект
func (it *One) create(table string, itone interface{}) bool {
	return ODB.Table(table).Create(itone) != nil
}

//	Прочитать объект
func (it *One) read(table string, id lik.IDB, itone interface{}) bool {
	return ODB.Table(table).First(itone, int(id)) != nil
}

//	Сохранить объект
func (it *One) save(table string, itone interface{}) bool {
	if it.CreatedAt.Year() < 2000 {
		it.CreatedAt = time.Now()
	}
	return ODB.Table(table).Save(itone) != nil
}

//	Обновить объект
func (it *One) update(one Oner, datas []interface{}) bool {
	valit := reflect.ValueOf(one).Elem()
	ndm := len(datas)
	for nd := 0; nd < ndm; nd++ {
		data := datas[nd]
		switch key := data.(type) {
		case string:
			var val interface{}
			if match := lik.RegExParse(key, "(.+?)=(.*)"); match != nil {
				key = match[1]
				val = match[2]
			} else if nd + 1 < ndm {
				nd++
				val = datas[nd]
			} else {
				break
			}
			if field := valit.FieldByName(key); field.IsValid() {
				if typ := field.Type().Name(); typ == "string" {
					field.SetString(toString(val))
				} else if typ == "int" {
					field.SetInt(int64(toInt(val)))
				} else if typ == "float" {
					field.SetFloat(toFloat(val))
				}
			} else {
				fmt.Println("Update bad field: ", one.Table(), ": ", key)
				return false
			}
		default:
			fmt.Println("Update ERROR: ", one.Table(), ": ", data)
			return false
		}
	}
	return one.Save()
}

//	Удалить объект
func (it *One) delete(table string, itone interface{}) {
	ODB.Table(table).Delete(itone)
}

//	Интефейс в целое
func toInt(data interface{}) int {
	val := 0
	if data != nil {
		switch da := data.(type) {
		case int:
			val = da
		case uint:
			val = int(da)
		case byte:
			val = int(da)
		case string:
			val = lik.StrToInt(da)
		}
	}
	return val
}

//	Интерфейс в плавающее
func toFloat(data interface{}) float64 {
	val := 0.0
	if data != nil {
		switch da := data.(type) {
		case float64:
			val = da
		case int:
			val = float64(da)
		case uint:
			val = float64(da)
		case byte:
			val = float64(da)
		case string:
			val = lik.StrToFloat(da)
		}
	}
	return val
}

//	Интерфейс в строку
func toString(data interface{}) string {
	val := ""
	if data != nil {
		switch da := data.(type) {
		case string:
			val = da
		case float64:
			val = lik.FloatToStr(da)
		case int:
			val = lik.IntToStr(da)
		case uint:
			val = lik.IntToStr(int(da))
		case byte:
			val = lik.IntToStr(int(da))
		}
	}
	return val
}

