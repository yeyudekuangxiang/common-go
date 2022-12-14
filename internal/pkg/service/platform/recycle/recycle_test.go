package recycle

import (
	"fmt"
	"reflect"
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
