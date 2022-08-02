package router

import (
	"github.com/gin-gonic/gin"
)

func Router(router *gin.Engine) {
	router.GET("/ping", func(context *gin.Context) {
		context.String(200, "pong")
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
                    }
                })
                break
            default:
                alert("暂不支持此跳转方式")
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
	apiRouter(router)
	adminRouter(router)
	openRouter(router)
	pugcRouter(router)
	BusinessRouter(router)
}
