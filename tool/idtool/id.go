package idtool

import (
	uuid "github.com/iris-contrib/go.uuid"
)

// UUID 生成uuid出现错误时会panic
func UUID() string {
	return uuid.Must(uuid.NewV4()).String()
}
