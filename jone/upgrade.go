package jone

import (
	"github.com/massarakhsh/tsan/one"
	"github.com/massarakhsh/lik"
	"github.com/massarakhsh/lik/likbase"
	"time"
)

//	Версия старой базы данных
var BaseOldVersion string

//	Запуск обновления
func UpgradeStart() {
	if SysElm == nil {
		SysElm = TableSystem.CreateElm()
	}
	SysElm.SetValue(int(time.Now().Unix()), "timestart")
	BaseOldVersion = SysElm.GetString("version")
}

//	Останов обновления
func UpgradeStop() {
	if lik.CompareVersion(one.Version, BaseOldVersion) > 0 {
		SysElm.SetValue(one.Version,"version")
		SysElm.OnModify()
	}
}

//	Обновить объект
func UpgradeElm(elm *likbase.ItElm) {
	if lik.CompareVersion(BaseOldVersion,"2.6.1") < 0 && elm.Table.Part == "member" {
		if !CalculateElmBool(elm, "do_second") {
			SetElmValue(elm, true, "do_second")
		}
	}
}

//	Обновить таблицу
func UpgradeTable(table *likbase.ItTable) {
}

