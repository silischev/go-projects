package main

import (
	"reflect"
)

func i2s(data interface{}, out interface{}) error {
	r := reflect.ValueOf(out).Elem()
	iType := r.Type()

	dataMap, _ := data.(map[string]interface{})

	for i := 0; i < r.NumField(); i++ {
		/*log.Println(iType.Field(i).Name)
		log.Println(r.Field(i).Type().Name())*/

		switch r.Field(i).Type().Name() {
		case "int":
			tmp := dataMap[iType.Field(i).Name].(float64)
			r.Field(i).SetInt(int64(tmp))
		case "string":
			r.Field(i).SetString(dataMap[iType.Field(i).Name].(string))
		case "bool":
			r.Field(i).SetBool(dataMap[iType.Field(i).Name].(bool))
		}
	}

	return nil
}
