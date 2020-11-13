package repo

import (
	"github.com/massarakhsh/lik"
)

//	Получить список разделов экспорта
func GetExportParts() []string {
	parts := []string{}
	if list := GetExportList(); list != nil {
		for nc := 0; nc < list.Count(); nc++ {
			if exp := list.GetSet(nc); exp != nil {
				part := exp.GetString("part")
				parts = append(parts, part)
			}
		}
	}
	return parts
}

//	Получить список дескрипторов экспорта
func GetExportList() lik.Lister {
	var list lik.Lister
	if ent := GenExtern.FindEnt("export"); ent != nil {
		list = ent.It.GetList("content")
	}
	return list
}

//	Получить прямой URL файла экспорта
func GetExternUrl(part string) string {
	url := ""
	if path := GetExternPath(part); path != "" {
		url = "http://rltweb.ru/" + path
	}
	return url
}

//	Получить прямой файл
func GetExternPath(part string) string {
	path := ""
	if _, ext := GenExtern.FindPart("export", part); ext != nil {
		path = "var/export/" + ext.GetString("signatura")
	}
	return path
}

