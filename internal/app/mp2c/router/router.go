package router

import (
	"github.com/gin-gonic/gin"
	"gitlab.miotech.com/miotech-application/backend/common-go/tool/httptool"
	"log"
	"mio/internal/pkg/core/app"
	"net/http"
	"strings"
)

func Router(router *gin.Engine) {
	router.GET("/ping", func(context *gin.Context) {
		if context.GetHeader("KEY") != "lvmio666" {
			context.Status(http.StatusServiceUnavailable)
			log.Println("ping key error")
			return
		}
		err := app.Ping(context)
		if err != nil {
			context.Status(http.StatusServiceUnavailable)
			log.Println("ping error", err)
			return
		}
		context.String(http.StatusOK, "pong")
	})

	router.GET("/", func(context *gin.Context) {
		context.String(200, "mio")
	})

	router.Any("MP_verify_pp3ZifoA3gboswNR.txt", func(context *gin.Context) {
		context.String(200, "pp3ZifoA3gboswNR")
	})
	router.Any("QUxp4PS6fh.txt", func(context *gin.Context) {
		context.String(200, "c636e427fa1d442771a93ff2885d6c15")
	})

	router.Any("pt04CfOnB5.txt", func(context *gin.Context) {
		context.String(200, "e2082465010d1787e6090c37ed629674")
	})

	router.Any("bwACP5dNsW.txt", func(context *gin.Context) {
		context.String(200, "c915276f12c8d1c2c20604b8d77072db")
	})

	router.Any("8LF0rq9WPf.txt", func(context *gin.Context) {
		context.String(200, "11d2f924ac51f2b502087d535c3c6b6e")
	})
	router.Any("vxqMmdbrL5.txt", func(context *gin.Context) {
		context.String(200, "fd8a7f90065f91224ceab80fcb0104ae")
	})

	router.Any("f316109564d78b58af9dfe4a38160d81.txt", func(context *gin.Context) {
		context.String(200, "d62cc36e9c20b9b294894f31bd5e620083dc4065")
	})

	router.GET("/duiba/jump.html", func(context *gin.Context) {
		context.Header("content-type", "text/html; charset=utf-8")
		context.Writer.WriteString(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>跳转中...</title>
</head>
<body>
<script type="text/javascript" src="https://res.wx.qq.com/open/js/jweixin-1.3.2.js"></script>
<script>
    window.onload = function () {
        try{
            const t = getQueryVariable("t")
            const path = decodeURIComponent(getQueryVariable("path"))
            gotoPath(t,path)
        }catch (e){
            alert("跳转失败")
            setTimeout(function (){
                history.back()
            },1000)
        }

    }

    function gotoPath(t,path){
        switch (t){
            case 'switchTab':
                wx.miniProgram.switchTab({
                    url:path,
                    fail:()=>{
                        alert("跳转失败")
                        setTimeout(function (){
                            history.back()
                        },1000)
                    },
					success:()=>{
						//此方法可以在跳转页点击返回时返回到真正的活动页
						history.back()
					}
                })
                break
            case 'reLaunch':
                wx.miniProgram.reLaunch({
                    url:path,
                    fail:()=>{
                        alert("跳转失败")
                        setTimeout(function (){
                            history.back()
                        },1000)
                    },
					success:()=>{
						//此方法可以在跳转页点击返回时返回到真正的活动页
						history.back()
					}
                })
                break
            case 'redirectTo':
                wx.miniProgram.redirectTo({
                    url:path,
                    fail:()=>{
                        alert("跳转失败")
                        setTimeout(function (){
                            history.back()
                        },1000)
                    },
					success:()=>{
						//此方法可以在跳转页点击返回时返回到真正的活动页
						history.back()
					}
                })
                break
            case 'navigateTo':
                wx.miniProgram.navigateTo({
                    url:path,
                    fail:()=>{
                        alert("跳转失败")
                        setTimeout(function (){
                            history.back()
                        },1000)
                    },
					success:()=>{
						//此方法可以在跳转页点击返回时返回到真正的活动页
						history.back()
					}
                })
                break
            case 'navigateBack':
                wx.miniProgram.navigateBack({
                    url:path,
                    fail:()=>{
                        alert("跳转失败")
                        setTimeout(function (){
                            history.back()
                        },1000)
                    },
					success:()=>{
						//此方法可以在跳转页点击返回时返回到真正的活动页
						history.back()
					}
                })
                break
            default:
                alert("暂不支持此跳转方式")
				setTimeout(function (){
					history.back()
				},1000)
                break
        }
    }
    function getQueryVariable(variable)
    {
        let query = window.location.search.substring(1);
        let vars = query.split("&");
        for (let i=0;i<vars.length;i++) {
            let pair = vars[i].split("=");
            if(pair[0] === variable){return pair[1];}
        }
        return false;
    }
</script>
</body>
</html>`)
	})
	router.GET("/mini/jump.html", func(context *gin.Context) {
		context.Header("content-type", "text/html; charset=utf-8")
		context.Writer.WriteString(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>跳转中...</title>
</head>
<body>
<script type="text/javascript" src="https://res.wx.qq.com/open/js/jweixin-1.3.2.js"></script>
<script>
    window.onload = function () {
        try{
            const t = getQueryVariable("t")
            const path = decodeURIComponent(getQueryVariable("path"))
            gotoPath(t,path)
        }catch (e){
            alert("跳转失败")
            setTimeout(function (){
                history.back()
            },1000)
        }

    }

    function gotoPath(t,path){
        switch (t){
            case 'switchTab':
                wx.miniProgram.switchTab({
                    url:path,
                    fail:()=>{
                        alert("跳转失败")
                        setTimeout(function (){
                            history.back()
                        },1000)
                    },
					success:()=>{
						//此方法可以在跳转页点击返回时返回到真正的活动页
						history.back()
					}
                })
                break
            case 'reLaunch':
                wx.miniProgram.reLaunch({
                    url:path,
                    fail:()=>{
                        alert("跳转失败")
                        setTimeout(function (){
                            history.back()
                        },1000)
                    },
					success:()=>{
						//此方法可以在跳转页点击返回时返回到真正的活动页
						history.back()
					}
                })
                break
            case 'redirectTo':
                wx.miniProgram.redirectTo({
                    url:path,
                    fail:()=>{
                        alert("跳转失败")
                        setTimeout(function (){
                            history.back()
                        },1000)
                    },
					success:()=>{
						//此方法可以在跳转页点击返回时返回到真正的活动页
						history.back()
					}
                })
                break
            case 'navigateTo':
                wx.miniProgram.navigateTo({
                    url:path,
                    fail:()=>{
                        alert("跳转失败")
                        setTimeout(function (){
                            history.back()
                        },1000)
                    },
					success:()=>{
						//此方法可以在跳转页点击返回时返回到真正的活动页
						history.back()
					}
                })
                break
            case 'navigateBack':
                wx.miniProgram.navigateBack({
                    url:path,
                    fail:()=>{
                        alert("跳转失败")
                        setTimeout(function (){
                            history.back()
                        },1000)
                    },
					success:()=>{
						//此方法可以在跳转页点击返回时返回到真正的活动页
						history.back()
					}
                })
                break
            default:
                alert("暂不支持此跳转方式")
				setTimeout(function (){
					history.back()
				},1000)
                break
        }
    }
    function getQueryVariable(variable)
    {
        let query = window.location.search.substring(1);
        let vars = query.split("&");
        for (let i=0;i<vars.length;i++) {
            let pair = vars[i].split("=");
            if(pair[0] === variable){return pair[1];}
        }
        return false;
    }
</script>
</body>
</html>`)
	})
	router.GET("/mini/jumpv2.html", func(context *gin.Context) {
		context.Header("content-type", "text/html; charset=utf-8")
		context.Writer.WriteString(`<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <title>跳转中...</title>
    <script type="text/javascript" src="https://res.wx.qq.com/open/js/jweixin-1.3.2.js"></script>
</head>

<body>
    <script>
        window.onload = function () {
            try {
                const t = getQueryVariable("t")
                const path = decodeURIComponent(getQueryVariable("path"))
                window.addEventListener('flutterInAppWebViewPlatformReady', () => {
                    const appPath = decodeURIComponent(getQueryVariable("appPath"))
                    goAppPath(t, appPath);
                }, false);
                if (isWechat()) {
                    gotoPath(t, path)
                } else {
                    console.log('不在微信环境中');
                }
            } catch (e) {
                console.error('appPath-----', e);
                fail()
            }

        }
        function goAppPath(t, path) {
            if (window.flutter_inappwebview) {
                switch (t) {
                    case 'switchTab':
                        window.flutter_inappwebview.callHandler('switchTab', path).then((res) => {
                            if (!res) {
                                fail();
                            }

                        }).catch(fail)
                        break
                    case 'redirectTo':
                        window.flutter_inappwebview.callHandler('redirectTo', {
                            path: path
                        }).then((res) => {
                            if (!res) {
                                fail();
                            }
                        }).catch(fail)
                        break
                    case 'navigateTo':
                        window.flutter_inappwebview.callHandler('navigateTo', {
                            path: path
                        }).then((res) => {
                            if (res) {
                                history.back()
                            } else {
                                fail();
                            }
                        }).catch(fail)
                        break
                    case 'navigateBack':
                        window.flutter_inappwebview.callHandler('navigateBack').then((res) => {
                            if (!res) {
                                fail();
                            }
                        }).catch(fail)
                        break
                    default:
                        alert("暂不支持此跳转方式")
                        setTimeout(function () {
                            history.back()
                        }, 1000)
                        break
                }
            }else{
                fail()
            }
        }

        function gotoPath(t, path) {
            switch (t) {
                case 'switchTab':
                    wx.miniProgram.switchTab({
                        url: path,
                        fail: fail,
                        success: () => {
                            //此方法可以在跳转页点击返回时返回到真正的活动页
                            history.back()
                        }
                    })
                    break
                case 'reLaunch':
                    wx.miniProgram.reLaunch({
                        url: path,
                        fail: fail,
                        success: () => {
                            //此方法可以在跳转页点击返回时返回到真正的活动页
                            history.back()
                        }
                    })
                    break
                case 'redirectTo':
                    wx.miniProgram.redirectTo({
                        url: path,
                        fail: fail,
                        success: () => {
                            //此方法可以在跳转页点击返回时返回到真正的活动页
                            history.back()
                        }
                    })
                    break
                case 'navigateTo':
                    wx.miniProgram.navigateTo({
                        url: path,
                        fail: fail,
                        success: () => {
                            //此方法可以在跳转页点击返回时返回到真正的活动页
                            //history.back()
                        }
                    })
                    break
                case 'navigateBack':
                    wx.miniProgram.navigateBack({
                        url: path,
                        fail:fail,
                        success: () => {
                            //此方法可以在跳转页点击返回时返回到真正的活动页
                            history.back()
                        }
                    })
                    break
                default:
                    alert("暂不支持此跳转方式")
                    setTimeout(function () {
                        history.back()
                    }, 1000)
                    break
            }
        }

        function getQueryVariable(variable) {
            let query = window.location.search.substring(1);
            let vars = query.split("&");
            for (let i = 0; i < vars.length; i++) {
                let pair = vars[i].split("=");
                if (pair[0] === variable) {
                    return pair[1];
                }
            }
            return false;
        }

        function isWechat() {
            const ua = navigator.userAgent.toLowerCase();
            return /micromessenger/i.test(ua) || /wechat/i.test(ua);
        }

        function fail() {
            alert("跳转失败")
            setTimeout(function () {
                history.back()
            }, 1000)
        }
    </script>
</body>

</html>`)
	})

	apple := `{
    "applinks":{
        "apps":[],
        "details":[
            {
                "appID":"5MLP8VB4J8.com.miotech.lvmio",
                "paths":["/app/*"]
            }
        ]
    }
}`
	router.GET("/apple-app-site-association", func(context *gin.Context) {
		context.Data(200, "application/octet-stream", []byte(apple))
	})
	router.GET("/.well-known/apple-app-site-association", func(context *gin.Context) {
		context.Data(200, "application/octet-stream", []byte(apple))
	})
	appHtml := getAppHtml()
	router.Any("/app/*path", func(context *gin.Context) {
		context.Header("content-type", "text/html")
		context.Writer.WriteString(appHtml)
	})
	apiRouter(router)
	adminRouter(router)
	openRouter(router)
	pugcRouter(router)
	BusinessRouter(router)
}
func getAppHtml() string {
	body, err := httptool.Get("https://applet-app.miotech.com/download.html")
	if err != nil {
		log.Printf("获取https://applet-app.miotech.com/download.html失败 %v", err)
		return ""
	}

	h := strings.ReplaceAll(string(body), `href="/`, `href="https://applet-app.miotech.com/`)
	h = strings.ReplaceAll(h, `src="/`, `src="https://applet-app.miotech.com/`)
	return h
}
