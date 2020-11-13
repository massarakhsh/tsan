// Интерфейс базы данных
package one

import (
	"fmt"
	"github.com/jinzhu/gorm"
)

//	Дескриптор базы данных
var ODB *gorm.DB

//	Инициализация базы данных
func GoIt(serv string, base string, user string, pass string) {
	args := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8&parseTime=True&loc=Local", user, pass, serv, base)
	if ODB,_ = gorm.Open("mysql", args); ODB == nil {
		return
	}
	InitializeDepart()
	InitializeMember()
	InitializeClient()
	InitializeBell()
	InitializeOffer()
	InitializeDeal()
	InitializeMessage()
	InitializeBonuses()
}

