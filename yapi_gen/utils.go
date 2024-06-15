package main

import (
	"strconv"
	"strings"
)

// 下划线转大写驼峰
func UnderscoreToUpperCamelCase(s string) string {
	s = strings.Replace(s, "_", " ", -1)
	s = strings.Title(s)
	return strings.Replace(s, " ", "", -1)
}

func GetGoType(val string) string {
	switch val {
	case "string":
		return "string"
	case "number":
		return "float64"
	case "boolean":
		return "bool"
	case "integer":
		return "int"
	}
	return ""
}

func GetGoTypeOrDefault(val string, structName string) string {
	switch val {
	case "string":
		return "string"
	case "number":
		return "float64"
	case "boolean":
		return "bool"
	case "integer":
		return "int"
	}
	return structName
}

// json去除转义字符
func JsonRemoveEscaping(val string) (string, error) {
	newVal, err := strconv.Unquote(`"` + val + `"`)
	if err != nil {
		return "", err
	}
	return newVal, nil
}
