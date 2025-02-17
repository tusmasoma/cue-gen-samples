package util

import (
	"html/template"

	"github.com/Masterminds/sprig"
)

func GetTmplFuncMap() template.FuncMap {
	funcMap := sprig.TxtFuncMap()
	myFuncMap := template.FuncMap{
		"sub": func(a, b int) int { return a - b },
		// 追加
	}
	for i := range myFuncMap {
		funcMap[i] = myFuncMap[i]
	}
	return funcMap
}
