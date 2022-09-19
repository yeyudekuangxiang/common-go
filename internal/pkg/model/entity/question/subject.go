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
		return "衣1"
	case CATEGORY_FOOD:
		return "食2"
	case CATEGORY_LIVE:
		return "住3"
	case CATEGORY_USE:
		return "用4"
	case CATEGORY_WALK:
		return "行4"
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
