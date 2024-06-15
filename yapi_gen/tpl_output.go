package main

import (
	"fmt"
	"log"
	"os"
	"text/template"
)

var (
	structTpl = `{{range  $struct := .}}
        type {{$struct.StructName}} struct { {{range  $field := $struct.Fields}}
                        {{$field.FieldName}}    {{$field.FieldType}} ` + "`json:\"{{$field.JsonName}}\"`" + `
                {{- end}}
        }
        {{- end}}`
)

func structTplOutput(list []*StructNode) error {
	if len(list) < 1 {
		return nil
	}
	tmpl, err := template.New("structTpl").Parse(structTpl)
	if err != nil {
		log.Printf("structTplOutput template.New err:%s\n", err)
		return err
	}
	err = tmpl.Execute(os.Stdout, list)
	if err != nil {
		log.Printf("structTplOutput tmpl.Execute err:%s\n", err)
	}
	fmt.Println()
	return nil
}
