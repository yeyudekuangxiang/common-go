package srv_types

type GetQnrSubjectDTO struct {
	QnrId int64
	Scene string `json:"scene" form:"scene" binding:"" alias:"轮播图场景"`
}

type GetQnrOptionDTO struct {
	SubjectIds []int64
}
