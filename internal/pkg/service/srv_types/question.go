package srv_types

import (
	"mio/internal/app/mp2c/controller/api/api_types"
	"mio/internal/pkg/model"
	"mio/internal/pkg/model/entity/question"
)

type GetQuestionSubjectDTO struct {
	QuestionId int64
}

type GetQuestionUserDTO struct {
	UserId     int64
	OpenId     string
	QuestionId int64
	Nickname   string
}

type CreateQuestionAnswerDTO struct {
	QuestionId int64
	SubjectId  model.LongID
	UserId     model.LongID
	Answer     string
	Carbon     float64
}

type DeleteQuestionAnswerDTO struct {
	QuestionId int64
	UserId     int64
}

type AddQuestionAnswerDTO struct {
	OpenId     string
	Answer     []api_types.GetQuestionTypeAnswer
	UserId     int64
	QuestionId int64
}

type AddUserCarbonInfoDTO struct {
	CarbonYear         string               `json:"carbonYear"`
	CarbonToday        string               `json:"carbonToday"`
	CarbonDay          string               `json:"carbonDay"`
	CarbonCompletion   string               `json:"carbonCompletion"`
	CarbonClassify     []UserCarbonClassify `json:"carbonClassify"`
	CompareWithCountry string               `json:"comparisonWithCountry"`
	CompareWithGlobal  string               `json:"compareWithGlobal"`
	UserGroup          string               `json:"userGroup"`
	UserGroupTips      string               `json:"userGroupDesc"`
	Level              int8                 `json:"level"`
}

type AddUserCarbonInfoV2DTO struct {
	CarbonYear          string                                    `json:"carbonYear"`
	CarbonToday         string                                    `json:"carbonToday"`
	CarbonDay           string                                    `json:"carbonDay"`
	CarbonYearValue     float64                                   `json:"carbonYearValue"`
	CarbonTodayValue    float64                                   `json:"carbonTodayValue"`
	CarbonDayValue      float64                                   `json:"carbonDayValue"`
	CarbonClassify      []UserCarbonClassify                      `json:"carbonClassify"`
	UserGroup           UserGroup                                 `json:"userGroup"`
	User                User                                      `json:"user"`
	CarbonCompletion    []CarbonCompletion                        `json:"carbonCompletion"`
	CarbonGroup         []CarbonGroup                             `json:"carbonGroup"`
	TodayCarbonClassify []api_types.CarbonTransactionClassifyList `json:"todayCarbonClassify"`
	CarbonCountry       float64                                   `json:"carbonCountry"`
	CarbonGlobal        float64                                   `json:"carbonGlobal"`
}

type CarbonGroup struct {
	Key   string  `json:"key"`
	Value float64 `json:"value"`
}
type UserGroup struct {
	Group     string `json:"group"`
	GroupTips string `json:"groupDesc"`
	Level     int8   `json:"level"`
}
type User struct {
	Nickname string `json:"nickname"`
	Uid      int64  `json:"uid"`
}

type CarbonCompletion struct {
	Key string  `json:"key"`
	Val float64 `json:"val"`
}
type UserCarbonClassify struct {
	CategoryId   question.QuestionCategoryType `json:"categoryId"`
	CategoryName string                        `json:"categoryName"`
	Carbon       string                        `json:"carbon"`
	CarbonValue  float64                       `json:"carbonValue"`
}
