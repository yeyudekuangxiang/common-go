package sortutil

import (
	"reflect"
	"sort"
)

// Map 按照map的key生序排序 无法排序时将随机输出key
func Map(m interface{}, f func(key interface{})) {
	sortMap(m, f, func(key1, key2 interface{}) bool {
		key1V := reflect.ValueOf(key1)
		key2V := reflect.ValueOf(key2)
		switch key1V.Kind() {
		case reflect.String:
			return key1V.String() < key2V.String()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return key1V.Int() < key2V.Int()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return key1V.Uint() < key2V.Uint()
		case reflect.Float32, reflect.Float64:
			return key1V.Float() < key2V.Float()
		}
		return false
	})
}

// MapDesc 按照map的key降序排序 无法排序时将随机输出key
func MapDesc(m interface{}, f func(key interface{})) {
	sortMap(m, f, func(key1, key2 interface{}) bool {
		key1V := reflect.ValueOf(key1)
		key2V := reflect.ValueOf(key2)
		switch key1V.Kind() {
		case reflect.String:
			return key1V.String() > key2V.String()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return key1V.Int() > key2V.Int()
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return key1V.Uint() > key2V.Uint()
		case reflect.Float32, reflect.Float64:
			return key1V.Float() > key2V.Float()
		}
		return false
	})
}

func sortMap(m interface{}, f func(key interface{}), sortF func(key1, key2 interface{}) bool) {
	v := reflect.ValueOf(m)
	var keys []reflect.Value
	if v.Kind() == reflect.Map {
		keys = v.MapKeys()
	} else if v.Kind() == reflect.Ptr {
		if v.Elem().Kind() == reflect.Map {
			keys = v.Elem().MapKeys()
		} else {
			return
		}
	} else {
		return
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].Kind() != keys[2].Kind() {
			return false
		}
		switch keys[i].Kind() {
		case reflect.String,
			reflect.Int,
			reflect.Int8,
			reflect.Int16,
			reflect.Int32,
			reflect.Int64,
			reflect.Uint,
			reflect.Uint8,
			reflect.Uint16,
			reflect.Uint32,
			reflect.Uint64,
			reflect.Float32,
			reflect.Float64:
			return sortF(keys[i].Interface(), keys[j].Interface())
		}
		return false
	})

	for _, k := range keys {
		f(k.Interface())
	}
}
