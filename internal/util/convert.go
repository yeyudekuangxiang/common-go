package util

import (
	"encoding/json"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

func MapTo(data interface{}, v interface{}) error {
	bs, err := json.Marshal(data)
	if err != nil {
		return errors.WithStack(err)
	}
	return json.Unmarshal(bs, v)
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
