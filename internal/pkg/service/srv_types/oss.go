package srv_types

import (
	"encoding/json"
	"time"
)

type AssumeRoleParam struct {
	Scheme          string
	Method          string
	RoleArn         string
	RoleSessionName string
	DurationSeconds time.Duration
	Policy          StsPolicy
}
type GetOssPolicyTokenParam struct {
	//有效期
	ExpireTime time.Duration
	//单位 B
	MaxSize     int64
	UploadDir   string
	CallbackUrl string
	MimeTypes   []string
	MaxAge      int //缓存有效期 单位秒
}
type OssPolicyConfig struct {
	Expiration string        `json:"expiration"`
	Conditions []interface{} `json:"conditions"`
}

//AddContentLength 上传文件大小 单位 B 1GB=1024MB，1MB=1024KB，1KB=1024B
func (c *OssPolicyConfig) AddContentLength(maxLength int64) {
	c.Conditions = append(c.Conditions, []interface{}{
		"content-length-range",
		1,
		maxLength,
	})
}

// AddBucket 添加bucket限制
func (c *OssPolicyConfig) AddBucket(bucket string) {
	c.Conditions = append(c.Conditions, map[string]string{
		"bucket": bucket,
	})
}

func (c *OssPolicyConfig) AddMaxAge(maxAge int) {
	/*c.Conditions = append(c.Conditions, []interface{}{
		"not-in",
		"Cache-Control",
		[]string{"no-cache"},
	})*/
}

//AddContentType []string{"image/jpg","image/png"}
func (c *OssPolicyConfig) AddContentType(contentTypes []string) {
	if len(contentTypes) == 0 {
		return
	}
	c.Conditions = append(c.Conditions, []interface{}{
		"in",
		"$content-type",
		contentTypes,
	})
}

//AddStartWithKey 上传文件目录 mp2c/user/avatar
func (c *OssPolicyConfig) AddStartWithKey(dir string) {
	c.Conditions = append(c.Conditions, []interface{}{
		"starts-with",
		"$key",
		dir,
	})
}

type OssCallbackParam struct {
	CallbackUrl      string `json:"callbackUrl"`
	CallbackBody     string `json:"callbackBody"`
	CallbackBodyType string `json:"callbackBodyType"`
}

type OssPolicyToken struct {
	AccessKeyId string `json:"AccessKeyId"`
	Host        string `json:"host"`
	Expire      int64  `json:"expire"`
	Signature   string `json:"signature"`
	Policy      string `json:"policy"`
	Directory   string `json:"dir"`
	Callback    string `json:"callback"`
}

//StsPolicy https://help.aliyun.com/document_detail/138809.html
type StsPolicy struct {
	//Version 1
	Version   string      `json:"Version"`
	Statement []Statement `json:"Statement"`
}

func (s StsPolicy) String() string {
	data, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}
	return string(data)
}

//Statement https://help.aliyun.com/document_detail/93738.htm?spm=a2c4g.11186623.0.0.921b2f860HjrNF#concept-xg5-51g-xdb
type Statement struct {
	Effect    string                 `json:"Effect"`
	Action    []string               `json:"Action"`
	Resource  []string               `json:"Resource"`
	Condition map[string]interface{} `json:"Condition,omitempty"`
}

// SetEffect Allow、Deny
func (s *Statement) SetEffect(effect string) {
	s.Effect = effect
}

//AddAction Action元素用于描述允许或拒绝的特定操作，是必选元素。Action元素的取值：云服务所定义的API操作名称。
//
//Action元素的格式：<ram-code>:<action-name>。
//
//ram-code：云服务的RAM代码。更多信息，请参见支持RAM的云服务的RAM代码列 https://help.aliyun.com/document_detail/28630.htm?spm=a2c4g.11186623.0.0.4d7936e4QpIgYY#concept-ofk-yt2-xdb。
//action-name：相关的API操作接口名称。
func (s *Statement) AddAction(action string) {
	s.Action = append(s.Action, action)
}
func (s *Statement) SetAction(actions []string) {
	s.Action = actions
}

//AddResource Resource元素用于描述被授权的一个或多个对象，是必选元素。Resource元素的取值：云服务所定义的资源ARN（Aliyun Resource Name）。
//
//Resource元素的格式遵循阿里云ARN统一规范。具体如下：
//
//acs：Alibaba Cloud Service的首字母缩写，表示阿里云的公共云平台。
//ram-code：云服务RAM代码。更多信息，请参见支持RAM的云服务的RAM代码列。
//region：地域信息。对于全局资源（无需指定地域就可以访问的资源），该字段置空。更多信息，请参见地域和可用区。
//account-id：阿里云账号ID。例如：123456789012****。
//relative-id：与服务相关的资源描述部分，其语义由具体云服务指定。这部分的格式支持树状结构（类似文件路径）。以OSS为例，表示一个OSS对象的格式为：relative-id = “mybucket/dir1/object1.jpg”
//Resource元素的值严格区分大小写，请按照云服务提供的鉴权文档使用准确的资源ARN。此外，如果将Resource元素的值置空，将会导致非预期的授权行为。因此，在创建和更新权限策略时，请检查策略内容，避免将Resource元素置空。
//"Resource": [
//  "acs:ecs:*:*:instance/inst-001",
//  "acs:ecs:*:*:instance/inst-002",
//  "acs:oss:*:*:mybucket",
//  "acs:oss:*:*:mybucket/*"
//]
func (s *Statement) AddResource(resource string) {
	s.Resource = append(s.Resource, resource)
}
func (s *Statement) SetResource(resources []string) {
	s.Resource = resources
}

//SetCondition Condition元素用于指定授权生效的限制条件，是可选元素。Condition元素也称为条件块（Condition Block），它是由一个或多个条件子句构成。一个条件子句由条件操作类型、条件关键字和条件值组成。
//条件关键字的名称（key）严格区分大小写，条件值（value）是否区分大小写取决于您使用的条件运算符。例如：针对字符串类型的条件关键字，如果使用StringEquals运算符，则会将策略内容中的值和请求中的值进行匹配，区分大小写。如果使用StringEqualsIgnoreCase运算符，则会将策略内容中的值和请求中的值进行匹配，忽略大小写。
func (s *Statement) SetCondition(key string, val map[string]interface{}) {
	s.Condition[key] = val
}
