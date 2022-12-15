package recycle

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestKind(t *testing.T) {
	m1 := make(map[int]interface{}, 0)
	m1[1] = 1.3
	m1[3] = map[string]float64{"床": 130.8}
	s1 := "hello"
	kind := reflect.ValueOf(s1).Kind()
	fmt.Println(kind)
}

func TestRune(t *testing.T) {
	if forName, ok := recyclePointForName[5]; ok {
		v := reflect.ValueOf(forName)
		if v.Kind() == reflect.Map {
			it := v.MapRange()
			for it.Next() {
				if strings.ContainsAny(it.Key().String(), "沙发") {
					fmt.Println(it.Key().String(), it.Value().Float())
				}
			}
		}
	}
}
