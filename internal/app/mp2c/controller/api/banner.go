package api

import (
	"github.com/gin-gonic/gin"
	"mio/internal/pkg/util/apiutil"
)

var DefaultBannerController = BannerController{}

type BannerController struct {
}

func (BannerController) GetBannerList(ctx *gin.Context) (gin.H, error) {
	form := GetGetBannerListForm{}
	if err := apiutil.BindForm(ctx, &form); err != nil {
		return nil, err
	}
	switch form.Type {
	case "event":
		return gin.H{
			"list": []map[string]string{
				{
					"name":     "首页",
					"path":     "pages/home/index",
					"imageUrl": "https://resources.miotech.com/static/mp2c/images/general/home-banner-unicom_d8ahts.png",
				},
			},
		}, nil
	case "topic":
		return gin.H{
			"list": []map[string]string{
				{
					"name":     "banner1",
					"path":     "/pages/webview/index?url=https://mp.weixin.qq.com/s/NcP2QrU74to7xLbJJPoBOw",
					"imageUrl": "https://resources.miotech.com/static/mp2c/images/home/banner/oldthings.jpg",
				},
				{
					"name":     "banner2",
					"path":     "/pages/cool-mio/mio-list-tag/index?id=108",
					"imageUrl": "https://resources.miotech.com/static/mp2c/images/home/banner/banneresg.gif",
				},
				{
					"name":     "banner3",
					"path":     "/pages/webview-page/cool-mio/index",
					"imageUrl": "https://resources.miotech.com/static/mp2c/images/cool-mio/banner/reward.png",
				},
				{
					"name":     "banner4",
					"path":     "/pages/webview/index?url=https://mp.weixin.qq.com/s/kHncLffScxCcG8YMZLVgaQ",
					"imageUrl": "https://resources.miotech.com/static/mp2c/images/cool-mio/banner/cwei.jpg",
				},
			},
		}, nil
	}
	return gin.H{
		"list": make([]map[string]string, 0),
	}, nil
}
