package srv_types

import "mio/internal/app/mp2c/controller/api/api_types"

type GetQnrSubjectDTO struct {
	QnrId int64
}

type GetQnrOptionDTO struct {
	SubjectIds []int64
}

type GetQnrUserDTO struct {
	UserId int64
	OpenId string
}

type CreateQnrAnswerDTO struct {
	QnrId     int64
	SubjectId int64
	UserId    int64
	Answer    string
}

type AddQnrAnswerDTO struct {
	OpenId string
	Answer []api_types.GetQnrTypeAnswer
	UserId int64
}
