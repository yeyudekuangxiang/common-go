目前支持的markdown语法是如下的子集：

标题 （支持1至6级标题，注意#与文字中间要有空格）
# 标题一
## 标题二
### 标题三
#### 标题四
##### 标题五
###### 标题六
加粗
**bold** 

链接
[这是一个链接](http://work.weixin.qq.com/api/doc) 

行内代码段（暂不支持跨行）
`code`
引用
> 引用文字
字体颜色(只支持3种内置颜色)
<font color="info">绿色</font>
<font color="comment">灰色</font>
<font color="warning">橙红色</font>




### 模版卡片类型  
### 文本通知模版卡片  
请求参数

参数	类型	必须	说明  
msgtype	String	是	消息类型，此时的消息类型固定为template_card  
template_card	Object	是	具体的模版卡片参数  

template_card的参数说明  
参数	类型	必须	说明  
card_type	String	是	模版卡片的模版类型，文本通知模版卡片的类型为text_notice  
source	Object	否	卡片来源样式信息，不需要来源样式可不填写  
source.icon_url	String	否	来源图片的url  
source.desc	String	否	来源图片的描述，建议不超过13个字  
source.desc_color	Int	否	来源文字的颜色，目前支持：0(默认) 灰色，1 黑色，2 红色，3 绿色  
main_title	Object	是	模版卡片的主要内容，包括一级标题和标题辅助信息  
main_title.title	String	否	一级标题，建议不超过26个字。模版卡片主要内容的一级标题main_title.title和二级普通文本sub_title_text必须有一项填写  
main_title.desc	String	否	标题辅助信息，建议不超过30个字  
emphasis_content	Object	否	关键数据样式  
emphasis_content.title	String	否	关键数据样式的数据内容，建议不超过10个字  
emphasis_content.desc	String	否	关键数据样式的数据描述内容，建议不超过15个字  
quote_area	Object	否	引用文献样式，建议不与关键数据共用  
quote_area.type	Int	否	引用文献样式区域点击事件，0或不填代表没有点击事件，1 代表跳转url，2 代表跳转小程序  
quote_area.url	String	否	点击跳转的url，quote_area.type是1时必填  
quote_area.appid	String	否	点击跳转的小程序的appid，quote_area.type是2时必填  
quote_area.pagepath	String	否	点击跳转的小程序的pagepath，quote_area.type是2时选填  
quote_area.title	String	否	引用文献样式的标题  
quote_area.quote_text	String	否	引用文献样式的引用文案  
sub_title_text	String	否	二级普通文本，建议不超过112个字。模版卡片主要内容的一级标题main_title.title和二级普通文本sub_title_text必须有一项填写  
horizontal_content_list	Object[]	否	二级标题+文本列表，该字段可为空数组，但有数据的话需确认对应字段是否必填，列表长度不超过6  
horizontal_content_list.type	Int	否	链接类型，0或不填代表是普通文本，1 代表跳转url，2 代表下载附件，3 代表@员工  
horizontal_content_list.keyname	String	是	二级标题，建议不超过5个字  
horizontal_content_list.value	String	否	二级文本，如果horizontal_content_list.type是2，该字段代表文件名称（要包含文件类型），建议不超过26个字  
horizontal_content_list.url	String	否	链接跳转的url，horizontal_content_list.type是1时必填  
horizontal_content_list.media_id	String	否	附件的media_id，horizontal_content_list.type是2时必填  
horizontal_content_list.userid	String	否	被@的成员的userid，horizontal_content_list.type是3时必填  
jump_list	Object[]	否	跳转指引样式的列表，该字段可为空数组，但有数据的话需确认对应字段是否必填，列表长度不超过3  
jump_list.type	Int	否	跳转链接类型，0或不填代表不是链接，1 代表跳转url，2 代表跳转小程序  
jump_list.title	String	是	跳转链接样式的文案内容，建议不超过13个字  
jump_list.url	String	否	跳转链接的url，jump_list.type是1时必填  
jump_list.appid	String	否	跳转链接的小程序的appid，jump_list.type是2时必填  
jump_list.pagepath	String	否	跳转链接的小程序的pagepath，jump_list.type是2时选填  
card_action	Object	是	整体卡片的点击跳转事件，text_notice模版卡片中该字段为必填项  
card_action.type	Int	是	卡片跳转类型，1 代表跳转url，2 代表打开小程序。text_notice模版卡片中该字段取值范围为[1,2]  
card_action.url	String	否	跳转事件的url，card_action.type是1时必填  
card_action.appid	String	否	跳转事件的小程序的appid，card_action.type是2时必填  
card_action.pagepath	String	否	跳转事件的小程序的pagepath，card_action.type是2时选填  
