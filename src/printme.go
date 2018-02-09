package main

import (
	"bytes"
	"fmt"
	"log"
)

func processOne(key string, keyType interface{}, prefix string, w *bytes.Buffer) {
	w.WriteString(prefix)
	if key != "" {
		w.WriteString(fmt.Sprintf(`"%s": `, key))
	}
	switch value := keyType.(type) {
	case string:
		w.WriteString(fmt.Sprintf(`"%s"`, value))
	case bool:
		w.WriteString(fmt.Sprintf(`%t`, value))
	case int:
		w.WriteString(fmt.Sprintf(`%d`, value))
	case float64:
		w.WriteString(fmt.Sprintf(`%f`, value))
	case map[string]interface{}, []interface{}, []map[string]string, []string, map[string]string:
		printme(value, prefix+"  ", w)
	default:
		log.Printf("1. Unrecognized type %T  %v\n", value, value)
	}
	w.WriteString("\n")
}

func printme(in interface{}, prefix string, w *bytes.Buffer) {
	switch rangeType := in.(type) {
	case map[string]interface{}:
		w.WriteString("{\n")
		for key, keyType := range rangeType {
			processOne(key, keyType, prefix, w)
		}
		w.WriteString(prefix + "}")
	case map[string]string:
		w.WriteString("{\n")
		for key, keyType := range rangeType {
			processOne(key, keyType, prefix, w)
		}
		w.WriteString(prefix + "}")
	case []interface{}:
		w.WriteString("[\n")
		for _, keyType := range rangeType {
			processOne("", keyType, prefix, w)
		}
		w.WriteString(prefix + "]")
	case []string:
		w.WriteString("[\n")
		for _, keyType := range rangeType {
			processOne("", keyType, prefix, w)
		}
		w.WriteString(prefix + "]")
	case []map[string]string:
		w.WriteString("[\n")
		for _, keyType := range rangeType {
			processOne("", keyType, prefix, w)
		}
		w.WriteString(prefix + "]")
	default:
		log.Printf("2. Unrecognized type %T for key %v\n", rangeType, rangeType)
	}
}

func xmain() {

	me := map[string]interface{}{
		"name": "ali",
		"age":  30,
		"things": map[string]interface{}{
			"class":    "2nd",
			"students": 30,
			"items":    []string{"one", "two", "three"},
			"children": []map[string]string{
				map[string]string{"a": "b"},
				map[string]string{"c": "d"},
			},
		},
	}

	var data bytes.Buffer

	printme(me, "", &data)

	fmt.Println(data.String())
}
