package config

var Constants = struct {
	WxWorkBugRobotKey          string
	WxWorkGitlabRobotKey       string
	DuiBaActivityNoLoginH5Link string
	DuiBaActivityStaticH5Link  string
	DuiBaActivityInsideLink    string
	DuiBaActivityEwmLink       string
}{
	WxWorkBugRobotKey:          "e46174d1-9e2a-493c-ac9c-84f2423457f6",
	WxWorkGitlabRobotKey:       "3c14f523-ebb5-4716-9aa5-d865bfda48c5",
	DuiBaActivityNoLoginH5Link: "https://go-api.miotech.com/api/mp2c/duiba/h5?activityId=%s",                                              //免登录H5落地页
	DuiBaActivityStaticH5Link:  "https://cloud1-1g6slnxm1240a5fb-1306244665.tcloudbaseapp.com/duiba_bind_phone.html?activityId=%s&cid=s%", //静态网站 H5 跳小程序链接：                                              //静态网站 H5 跳小程序链接
	DuiBaActivityInsideLink:    "pages/duiba_v2/duiba/index?activityId=%s",                                                                //小程序内路径
	DuiBaActivityEwmLink:       "",                                                                                                        //二维码                                                                                                      //
}
