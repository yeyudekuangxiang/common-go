package duiba

import (
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"sort"

	"strings"
)

func MapTo(data interface{}, v interface{}) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:  v,
		TagName: "json",
	})
	if err != nil {
		return errors.WithStack(err)
	}
	return errors.WithStack(decoder.Decode(data))
}
func BuildQuery(params map[string]string) string {
	kList := make([]string, 0)
	for k := range params {
		kList = append(kList, k)
	}

	sort.Strings(kList)

	query := strings.Builder{}
	for _, k := range kList {
		if params[k] == "" {
			continue
		}
		query.WriteString(k)
		query.WriteString("=")
		query.WriteString(params[k])
		query.WriteString("&")
	}
	return strings.TrimRight(query.String(), "&")
}
