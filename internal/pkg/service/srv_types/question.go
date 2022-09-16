package srv_types

import (
	"mio/internal/app/mp2c/controller/api/api_types"
	"mio/internal/pkg/model"
)

type GetQuestionSubjectDTO struct {
	QuestionId int64
}

type GetQuestionOptionDTO struct {
	SubjectIds []model.LongID
}

type GetQuestionUserDTO struct {
	UserId     int64
	OpenId     string
	QuestionId int64
}

type CreateQuestionAnswerDTO struct {
	QuestionId int64
	SubjectId  model.LongID
	UserId     model.LongID
	Answer     string
	Carbon     float64
}

type AddQuestionAnswerDTO struct {
	OpenId     string
	Answer     []api_types.GetQuestionTypeAnswer
	UserId     int64
	QuestionId int64
}
