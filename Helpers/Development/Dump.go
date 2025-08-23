package development

import (
	"fmt"
	"os"
	"reflect"
)

func Dump(values ...interface{}) {
	for i, v := range values {
		if reflect.ValueOf(v).Kind() == reflect.Struct ||
			reflect.ValueOf(v).Kind() == reflect.Slice ||
			reflect.ValueOf(v).Kind() == reflect.Map {
			fmt.Printf("[%d]: %+v\n", i, v)
		} else {
			fmt.Printf("[%d]: %v\n", i, v)
		}
	}
}

func Dd(values ...interface{}) {
	Dump(values...)
	os.Exit(1)
}
