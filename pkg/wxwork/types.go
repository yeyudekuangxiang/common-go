package wxwork

type MsgType string

const (
	MsgTypeText     MsgType = "text"
	MsgTypeMarkdown MsgType = "markdown"
	MsgTypeImage    MsgType = "image"
	MsgTypeFile     MsgType = "file"
	MsgTypeNews     MsgType = "news"
)

// Text
//    {
//    "msgtype": "text",
//    "text": {
//        "content": "广州今日天气：29度，大部分多云，降雨概率：60%",
//        "mentioned_list":["wangqing","@all"],
//        "mentioned_mobile_list":["13800001111","@all"]
//    }
//}/*
type Text struct {
	//文本内容，最长不超过2048个字节，必须是utf8编码
	Content string `json:"content"`
	//userid的列表，提醒群中的指定成员(@某个成员)，@all表示提醒所有人，如果开发者获取不到userid，可以使用mentioned_mobile_list
	MentionedList []string `json:"mentioned_list"`
	//手机号列表，提醒手机号对应的群成员(@某个成员)，@all表示提醒所有人
	MentionedMobileList []string `json:"mentioned_mobile_list"`
}

// Markdown
//   {
//    "msgtype": "text",
//    "text": {
//        "content": "广州今日天气：29度，大部分多云，降雨概率：60%",
//        "mentioned_list":["wangqing","@all"],
//        "mentioned_mobile_list":["13800001111","@all"]
//    }
//}/*
type Markdown struct {
	//markdown内容，最长不超过4096个字节，必须是utf8编码
	Content string `json:"content"`
}

//：图片（base64编码前）最大不能超过2M，支持JPG,PNG格式
type Image struct {
	//图片内容的base64编码
	Base64 string `json:"base64"`
	//图片内容（base64编码前）的md5值
	Md5 string `json:"md5"`
}

// News 图文类型
type News struct {
	Articles []struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Url         string `json:"url"`
		Picurl      string `json:"picurl"`
	} `json:"articles"`
}
type File struct {
	//消息类型，此时固定为file
	Msgtype MsgType `json:"msgtype"`
	File    struct {
		//文件id，通过下文的文件上传接口获取
		MediaId string `json:"media_id"`
	} `json:"file"`
}
