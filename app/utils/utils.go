package utils

import "reflect"

func GetStructFieldNames(s interface{}) []string {
    arr := []string{}
    val := reflect.ValueOf(s)
    for i := 0; i < val.Type().NumField(); i++ {
        arr = append(arr, val.Type().Field(i).Name)
    }
    return arr
}
