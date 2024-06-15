package main

import (
	"encoding/json"
	"log"
	"strings"
)

type StructNode struct {
	StructName string
	Fields     []*Field
}

type Field struct {
	FieldName string
	FieldType string
	JsonName  string
}

func Exec(token string, id int) {
	yapiClient := Client{}
	info, err := yapiClient.QueryInterfaceInfo(token, id)
	if err != nil {
		return
	}
	if info.ResBodyType != "json" {
		log.Println("响应体不是json结构")
		return
	}
	_, err = execParseStruct(info.ResBody, "Response")
	if err != nil {
		log.Println("响应体生成失败")
	}
	if info.Method != "GET" && info.ReqBodyType != "json" {
		log.Println("非GET请求的请求参数不是json结构")
		return
	}
	if info.Method != "GET" && info.ReqBodyOther != "" {
		_, err = execParseStruct(info.ReqBodyOther, "Request")
		if err != nil {
			log.Println("请求体生成失败")
		}
	}

}

func execParseStruct(yapiJson string, rootName string) ([]*StructNode, error) {
	data := make(map[string]interface{})
	err := json.Unmarshal([]byte(yapiJson), &data)
	if err != nil {
		log.Printf("parseStruct json.Unmarshal err:%s\n", err)
		return nil, yapiStructParseError
	}
	list := []*StructNode{}
	parseStruct(data, rootName, &StructNode{}, 0, &list)
	err = structTplOutput(list)
	if err != nil {
		log.Printf("parseStruct structTplOutput err:%s\n", err)
		return nil, yapiStructParseError
	}
	return nil, nil
}

func parseStruct1(data map[string]interface{}, cur *StructNode, list *[]*StructNode) {
	if len(data) == 0 {
		return
	}
	for field, val := range data {

		valM, ok := val.(map[string]interface{})
		if !ok {
			return
		}
		t := valM["type"].(string)
		if t == objectType {
			properties := valM["properties"].(map[string]interface{})
			tmpStruct := &StructNode{}
			tmpStruct.StructName = field
			parseStruct1(properties, tmpStruct, list)
			f := &Field{
				FieldName: field,
				FieldType: field,
			}
			cur.Fields = append(cur.Fields, f)
		}
		if t == stringType || t == numberType || t == integerType || t == booleanType {
			f := &Field{
				FieldName: field,
				FieldType: GetGoType(t),
			}
			cur.Fields = append(cur.Fields, f)
		}
	}
	if cur != nil {
		*list = append(*list, cur)
	}
}

func parseStruct(data map[string]interface{}, fieldName string, cur *StructNode, arrayLayer int, list *[]*StructNode) {
	t := data["type"].(string)
	if t == objectType {
		f := &Field{
			FieldName: UnderscoreToUpperCamelCase(fieldName),
			FieldType: UnderscoreToUpperCamelCase(GetGoTypeOrDefault(t, fieldName)),
			JsonName:  fieldName,
		}
		f.FieldType = strings.Repeat("[]", arrayLayer) + "*" + f.FieldType
		cur.Fields = append(cur.Fields, f)
		properties := data["properties"].(map[string]interface{})
		tmp := &StructNode{
			StructName: UnderscoreToUpperCamelCase(fieldName),
		}
		for field, val := range properties {
			v := val.(map[string]interface{})
			parseStruct(v, field, tmp, 0, list)
		}
		*list = append(*list, tmp)
	}
	if t == stringType || t == numberType || t == integerType || t == booleanType {
		f := &Field{
			FieldName: UnderscoreToUpperCamelCase(fieldName),
			FieldType: GetGoTypeOrDefault(t, ""),
			JsonName:  fieldName,
		}
		if arrayLayer > 0 {
			f.FieldType = strings.Repeat("[]", arrayLayer) + f.FieldType
		}
		cur.Fields = append(cur.Fields, f)
	}
	if t == arrayType {
		items := data["items"].(map[string]interface{})
		parseStruct(items, fieldName, cur, arrayLayer+1, list)
	}
	return
}
