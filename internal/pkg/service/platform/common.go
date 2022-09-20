package platform

import mioContext "mio/internal/pkg/core/context"

type platformCommon struct {
	PlatformKey string `json:"platformKey"`
	Context     mioContext.MioContext
}

func NewPlatFormCommon(platformKey string, ctx mioContext.MioContext) *platformCommon {
	return &platformCommon{
		PlatformKey: platformKey,
		Context:     ctx,
	}
}
