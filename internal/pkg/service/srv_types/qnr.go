package srv_types

import (
	"mio/internal/app/mp2c/controller/api/api_types"
	"mio/internal/pkg/model"
)

type GetQnrSubjectDTO struct {
	QnrId int64
}

type GetQnrOptionDTO struct {
	SubjectIds []model.LongID
}

type GetQnrUserDTO struct {
	UserId int64
	OpenId string
}

type CreateQnrAnswerDTO struct {
	QnrId     int64
	SubjectId model.LongID
	UserId    model.LongID
	Answer    string
}

type AddQnrAnswerDTO struct {
	OpenId string
	Answer []api_types.GetQnrTypeAnswer
	UserId int64
}
