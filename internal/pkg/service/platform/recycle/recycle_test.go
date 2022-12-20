package recycle

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"
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
	if forName, ok := recyclePointOfName[5]; ok {
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

func TestMonth(t *testing.T) {
	month := time.Now().Format("01")
	days := time.Now().AddDate(0, 1, -time.Now().Day()).Day()
	fmt.Println(month)
	fmt.Println(days)
}
