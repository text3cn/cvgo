package clog

import (
	"fmt"
	"log"
	"os"
	"reflect"
)

func P(data any, outputColor ...string) {
	color := ColorBlack
	if len(outputColor) > 0 {
		color = outputColor[0]
	}
	_log := log.New(os.Stderr, "", log.Lshortfile|log.LstdFlags)
	calldepth := 4
	// 类型
	rVal := reflect.ValueOf(data)
	rKind := rVal.Kind().String()
	rType := rVal.Type().String()
	typeStr := fmt.Sprintf(yellowBG+black+" %v -> %v "+reset, rKind, rType)

	// 正文
	fmt.Print(color)
	switch rKind {
	case "array", "slice":
		switch rType {
		case "[]int":
			_log.Output(calldepth, typeStr)
			fmt.Print(color)
			array := data.([]int)
			for k, v := range array {
				fmt.Println("[", k, "]", "=>", v)
			}
		case "[]int64":
			_log.Output(calldepth, typeStr)
			fmt.Print(color)
			array := data.([]int64)
			for k, v := range array {
				fmt.Println("[", k, "]", "=>", v)
			}
		case "[]string":
			_log.Output(calldepth, typeStr)
			fmt.Print(color)
			array := data.([]string)
			for k, v := range array {
				fmt.Println("[", k, "]", "=>", v)
			}
		case "[]uint8":
			_log.Output(calldepth, typeStr)
			fmt.Print(color)
			fmt.Printf("%v", string(rVal.Bytes()))
		case "[]map[string]interface {}":
			_log.Output(calldepth, typeStr)
			fmt.Print(color)
			array := data.([]map[string]interface{})
			for k, v := range array {
				fmt.Println("[", k, "]", "=>", "map (")
				for key, val := range v {
					fmt.Println("      ", "[", key, "]", "=>", val)
				}
				fmt.Println("),")
			}
		case "[]interface {}":
			_log.Output(calldepth, typeStr)
			fmt.Print(color)
			array := data.([]interface{})
			for k, v := range array {
				fmt.Println("\t", "[", k, "]", "=>", v)
			}
		default:
			str := fmt.Sprintf(color+" %v "+reset, data)
			_log.Output(4, typeStr+str)
			fmt.Print(color)
		}
	case "map":
		switch rType {
		case "map[string]interface {}":
			_log.Output(calldepth, typeStr)
			fmt.Print(color)
			for k, v := range data.(map[string]interface{}) {
				fmt.Println("\t", "[", k, "]", "=>", v)
			}
		default:
			str := fmt.Sprintf(color+" %v "+reset, data)
			_log.Output(4, typeStr+str)
			fmt.Print(color)
		}
	case "struct":
		_log.Output(calldepth, typeStr)
		fmt.Print(color)
		t := reflect.TypeOf(data)
		v := reflect.ValueOf(data)
		for k := 0; k < t.NumField(); k++ {
			fmt.Println("\t", "[", t.Field(k).Name, "]", "=>", v.Field(k).Interface())
		}
	default:
		str := fmt.Sprintf(color+" %v "+reset, data)
		_log.Output(4, typeStr+str)
		fmt.Print(color)
	}
	// 重置样式
	fmt.Printf(reset)
}
