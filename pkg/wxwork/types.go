package wxwork

type MsgType string

const (
	MsgTypeText     MsgType = "text"
	MsgTypeMarkdown MsgType = "markdown"
	MsgTypeImage    MsgType = "image"
	MsgTypeFile     MsgType = "file"
	MsgTypeNews     MsgType = "news"
	MsgTypeCard     MsgType = "template_card"
)

type IMessage interface {
	F4CD9CEE4B485C82938A4DD27E1704B0()
}

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

func (t Text) F4CD9CEE4B485C82938A4DD27E1704B0() {}

// Markdown
/*
{
    "msgtype": "markdown",
    "markdown": {
        "content": "实时新增用户反馈<font color=\"warning\">132例</font>，请相关同事注意。\n
         >类型:<font color=\"comment\">用户反馈</font>
         >普通用户反馈:<font color=\"comment\">117例</font>
         >VIP用户反馈:<font color=\"comment\">15例</font>"
    }
}
*/
type Markdown struct {
	//markdown内容，最长不超过4096个字节，必须是utf8编码
	Content string `json:"content"`
	//userid的列表，提醒群中的指定成员(@某个成员)，@all表示提醒所有人，如果开发者获取不到userid，可以使用mentioned_mobile_list
	MentionedList []string `json:"mentioned_list"`
	//手机号列表，提醒手机号对应的群成员(@某个成员)，@all表示提醒所有人
	MentionedMobileList []string `json:"mentioned_mobile_list"`
}

func (t Markdown) F4CD9CEE4B485C82938A4DD27E1704B0() {}

//：图片（base64编码前）最大不能超过2M，支持JPG,PNG格式
type Image struct {
	//图片内容的base64编码
	Base64 string `json:"base64"`
	//图片内容（base64编码前）的md5值
	Md5 string `json:"md5"`
}

func (t Image) F4CD9CEE4B485C82938A4DD27E1704B0() {}

// News 图文类型
type News struct {
	Articles []struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Url         string `json:"url"`
		Picurl      string `json:"picurl"`
	} `json:"articles"`
}

func (t News) F4CD9CEE4B485C82938A4DD27E1704B0() {}

type File struct {
	File struct {
		//文件id，通过下文的文件上传接口获取
		MediaId string `json:"media_id"`
	} `json:"file"`
}

func (t File) F4CD9CEE4B485C82938A4DD27E1704B0() {}

// Card 卡片类型
type CardText struct {
	CardType              string              `json:"card_type"`
	Source                Source              `json:"source"`
	MainTitle             MainTitle           `json:"main_title"`
	EmphasisContent       EmphasisContent     `json:"emphasis_content"`
	QuoteArea             QuoteArea           `json:"quote_area"`
	SubTitleText          string              `json:"sub_title_text"`
	HorizontalContentList []HorizontalContent `json:"horizontal_content_list"`
	JumpList              []Jump              `json:"jump_list"`
	CardAction            CardAction          `json:"card_action"`
}

func (t CardText) F4CD9CEE4B485C82938A4DD27E1704B0() {}

type Source struct {
	IconUrl   string `json:"icon_url"`
	Desc      string `json:"desc"`
	DescColor int    `json:"desc_color"`
}
type MainTitle struct {
	Title string `json:"title"`
	Desc  string `json:"desc"`
}
type EmphasisContent struct {
	Title string `json:"title"`
	Desc  string `json:"desc"`
}
type QuoteArea struct {
	Type      int    `json:"type"`
	Url       string `json:"url"`
	Appid     string `json:"appid"`
	Pagepath  string `json:"pagepath"`
	Title     string `json:"title"`
	QuoteText string `json:"quote_text"`
}
type HorizontalContent struct {
	Keyname string `json:"keyname"`
	Value   string `json:"value"`
	Type    int    `json:"type,omitempty"`
	Url     string `json:"url,omitempty"`
	MediaId string `json:"media_id,omitempty"`
	UserId  string `json:"userid,omitempty"`
}

type Jump struct {
	Type     int    `json:"type"`
	Url      string `json:"url,omitempty"`
	Title    string `json:"title"`
	Appid    string `json:"appid,omitempty"`
	Pagepath string `json:"pagepath,omitempty"`
}
type CardAction struct {
	Type     int    `json:"type"`
	Url      string `json:"url"`
	Appid    string `json:"appid"`
	Pagepath string `json:"pagepath"`
}
type CardNews struct {
	CardType  string    `json:"card_type"`
	Source    Source    `json:"source"`
	MainTitle MainTitle `json:"main_title"`

	CardImage             CardImage           `json:"card_image"`
	ImageTextArea         ImageTextArea       `json:"image_text_area"`
	QuoteArea             QuoteArea           `json:"quote_area"`
	VerticalContentList   []VerticalContent   `json:"vertical_content_list"`
	HorizontalContentList []HorizontalContent `json:"horizontal_content_list"`
	JumpList              []Jump              `json:"jump_list"`
	CardAction            CardAction          `json:"card_action"`
}

func (t CardNews) F4CD9CEE4B485C82938A4DD27E1704B0() {}

type CardImage struct {
	Url         string  `json:"url"`
	AspectRatio float64 `json:"aspect_ratio"`
}
type ImageTextArea struct {
	Type     int    `json:"type"`
	Url      string `json:"url"`
	Title    string `json:"title"`
	Desc     string `json:"desc"`
	ImageUrl string `json:"image_url"`
}
type VerticalContent struct {
	Title string `json:"title"`
	Desc  string `json:"desc"`
}
