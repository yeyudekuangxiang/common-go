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
		return "小贴士：尽量选择棉、麻质的环保面料衣服，适量购衣，够穿就好！"
	case CATEGORY_FOOD:
		return "小贴士：我国每年浪费的粮食高达3500万吨。节约粮食，光盘行动！"
	case CATEGORY_LIVE:
		return "小贴士：使用节能灯具电器、记得随手关灯。节能节电，从我做起！"
	case CATEGORY_USE:
		return "小贴士：去海边冲浪要选择海洋保护型防晒。力所能及，保护环境！"
	case CATEGORY_WALK:
		return "小贴士：精简车内物品减轻车重，节约能耗。绿色出行，始于足下！"
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
