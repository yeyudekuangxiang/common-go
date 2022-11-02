package util

import (
	"encoding/json"
	"github.com/pkg/errors"
	"log"
	"reflect"
	"strconv"
	"strings"
)

func MapTo(data interface{}, v interface{}, strict ...bool) error {
	bs, err := json.Marshal(data)
	if err != nil {
		return errors.WithStack(err)
	}
	err = json.Unmarshal(bs, v)
	if len(strict) > 1 && strict[0] {
		return errors.WithStack(err)
	}
	switch err.(type) {
	case *json.UnsupportedValueError, *json.UnmarshalTypeError, *json.UnsupportedTypeError:
		log.Printf("MapTo normal err %+v %+v \n", err, data)
		err = nil
	}
	return errors.WithStack(err)
}

func StrToArrayInt(str string, sep string) ([]int, error) {
	list := make([]int, 0)
	if len(str) == 0 {
		return list, nil
	}
	strs := strings.Split(str, sep)

	for _, item := range strs {
		data, err := strconv.Atoi(item)
		if err != nil {
			return list, err
		}
		list = append(list, data)
	}
	return list, nil
}

func InterfaceToString(data interface{}) string {
	var key string
	switch data.(type) {
	case string:
		key = data.(string)
	case int:
		key = strconv.Itoa(data.(int))
	case int64:
		it := data.(int64)
		key = strconv.FormatInt(it, 10)
	case float64:
		it := data.(float64)
		key = strconv.FormatFloat(it, 'f', -1, 64)
	}
	return key
}

func Map2SliceE(i interface{}) ([]interface{}, []interface{}, error) {
	kind := reflect.TypeOf(i).Kind()
	if kind != reflect.Map {
		return nil, nil, errors.New("the input is not a map")
	}
	m := reflect.ValueOf(i)
	keys := m.MapKeys()
	slK, slV := make([]interface{}, 0, len(keys)), make([]interface{}, 0, len(keys))
	for _, k := range keys {
		slK = append(slK, k.Interface())
		v := m.MapIndex(k)
		slV = append(slV, v.Interface())
	}
	return slK, slV, nil
}
