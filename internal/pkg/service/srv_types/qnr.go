package srv_types

type GetQnrSubjectDTO struct {
	QnrId int64
	Scene string `json:"scene" form:"scene" binding:"" alias:"轮播图场景"`
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
	Answer string
	UserId int64
}
