# 注意事项
## 1.gin form 时间绑定
如果传参方式为form或者get方式提交 则可使用下面方式绑定
```
type TestForm struct {
StartTime time.Time `form:"startTime"  time_format:"2006-01-02" time_utc:"false" time_location:"Asia/Shanghai"`
}
```
如果传参方式为json 则可使用下面方式绑定 默认接受格式为2006-01-02 15:04:05
```
type TestForm struct {
StartTime timeutils.Time `json:"startTime"`
}
```
如果想要兼容form和json 则只能使用string类型接收

```
type TestForm struct {
StartTime string `json:"startTime" form:"startTime" binding:"datetime=2006-01-02"`
}
```

2 float32 转换float64时数据异常

```
var a float32 = 100.123
fmt.Println(a)
fmt.Println(float64(a))
```
输出结果: 
```
100.123
100.12300109863281
```
使用int、decimal.Decimal替换float