package util

import (
	"github.com/pkg/errors"
	"mio/core/app"
	"os"
	"strings"
)

func GetAppConfig(key string) (string, error) {
	if key == "" {
		return "", errors.New("key can not be empty")
	}
	keys := strings.Split(key, ".")
	if len(keys) == 1 {
		return "", errors.New("need section")
	}
	section := app.Ini.Section(keys[0])
	if section == nil {
		return "", errors.New("not found section")
	}
	item := section.Key(keys[1])
	if item == nil {
		return "", errors.New("not found key")
	}
	return item.Value(), nil
}
func IsTesting() bool {
	return os.Getenv("TEST_ENV") != ""
}
