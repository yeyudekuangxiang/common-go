package common

import (
	mioContext "mio/internal/pkg/core/context"
	"mio/internal/pkg/service/platform/jhx"
)

type platformCommon struct {
	PlatformKey string `json:"platformKey"`
	Context     *mioContext.MioContext
}

func NewPlatFormCommon(platformKey string, ctx *mioContext.MioContext) *platformCommon {
	return &platformCommon{
		PlatformKey: platformKey,
		Context:     ctx,
	}
}

func (receiver platformCommon) SwitchService() interface{} {
	switch receiver.PlatformKey {
	case "jinhuaxing":
		return jhx.NewJhxService(receiver.Context)
	}
	return nil
}
