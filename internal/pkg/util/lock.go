package util

import (
	"gitlab.miotech.com/miotech-application/backend/common-go/lock"
	"mio/config"
	"mio/internal/pkg/core/app"
)

var DefaultLock = lock.RedisDistributedLock{Prefix: config.RedisKey.Lock + "default_", Redis: app.Redis}
