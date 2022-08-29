package api_types

import "mio/internal/pkg/model"

type GetSubjectListForm struct {
	Scene string `json:"scene" form:"scene" binding:"oneof=home event topic" alias:"banner场景"`
}

type GetQnrSubjectCreateDTO struct {
	Answer []GetQnrTypeAnswer `json:"answer"`
}

type GetQnrTypeAnswer struct {
	Id     model.LongID `json:"id"`
	Answer string       `json:"answer"`
}

type QnrListVo struct {
	Title string  `json:"title"`
	List  []QnrVo `json:"list"`
}

type QnrCategory struct {
	Id    int64  `json:"id"`
	Title string `json:"title"`
}

type QnrVo struct {
	ID        int64        `json:"id"`
	Title     string       `json:"title"`
	Type      int8         `json:"type"`
	Remind    string       `json:"remind"`
	IsHide    int8         `json:"isHide"`
	Option    []OptionVO   `json:"option"`
	SubjectId model.LongID `json:"subjectId"`
}

type OptionVO struct {
	ID             int64  `json:"id"`
	Title          string `json:"title"`
	Remind         string `json:"remind"`
	JumpSubject    int64  `json:"jumpSubject"`
	RelatedSubject string `json:"relatedSubject"`
}
