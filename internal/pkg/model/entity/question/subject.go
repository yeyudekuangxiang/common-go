package question

import "mio/internal/pkg/model"

type Subject struct {
	ID         int64
	SubjectId  model.LongID
	QuestionId int64
	CategoryId QuestionCategoryType
	Title      string
	Type       int8
	IsHide     int8
	Remind     string
	Sort       int8
}

type QuestionCategoryType string

// Text 获取积分类型的中文名称(给用户看的)
func (p QuestionCategoryType) Text() string {
	switch p {
	case CATEGORY_CLOTHES:
		return "衣"
	case CATEGORY_FOOD:
		return "食"
	case CATEGORY_LIVE:
		return "住"
	case CATEGORY_USE:
		return "用"
	case CATEGORY_WALK:
		return "行"
	}
	return "未知"
}

func (p QuestionCategoryType) DescText() string {
	switch p {
	case CATEGORY_CLOTHES:
		return "小贴士：选择天然材质，提高衣物的使用率，回收利用，减少碳排放"
	case CATEGORY_FOOD:
		return "小贴士：我国每年的粮食浪费导致近5000万吨的二氧化碳排放。节约粮食，光盘行动势在必行"
	case CATEGORY_LIVE:
		return "小贴士：节水节电，减少煤气等能源过度使用，可减少碳排放"
	case CATEGORY_USE:
		return "小贴士：减少一次性制品、洗发水等物品的使用，将有效降低碳排放"
	case CATEGORY_WALK:
		return "小贴士：一天不开车，按上下班来回共10公里及每1公升汽油行驶10公里估算，可以减碳2240克"
	}
	return "未知"
}

const (
	CATEGORY_CLOTHES QuestionCategoryType = "CLOTHES" //衣
	CATEGORY_FOOD    QuestionCategoryType = "FOOD"    //食
	CATEGORY_LIVE    QuestionCategoryType = "LIVE"    //住
	CATEGORY_USE     QuestionCategoryType = "USE"     //用
	CATEGORY_WALK    QuestionCategoryType = "WALK"    //行
)

var QuestionCategoryTypeMap = []QuestionCategoryType{
	CATEGORY_CLOTHES,
	CATEGORY_FOOD,
	CATEGORY_LIVE,
	CATEGORY_USE,
	CATEGORY_WALK,
}

func (Subject) TableName() string {
	return "question_subject"
}
