package util

import (
	"html/template"

	"github.com/Masterminds/sprig"
	"github.com/iancoleman/strcase"
)

func GetTmplFuncMap() template.FuncMap {
	funcMap := sprig.TxtFuncMap()
	myFuncMap := template.FuncMap{
		"sub":            func(a, b int) int { return a - b },
		"lowerCamelcase": strcase.ToLowerCamel,
		"upperCamelcase": strcase.ToCamel,
	}
	for i := range myFuncMap {
		funcMap[i] = myFuncMap[i]
	}
	return funcMap
}
