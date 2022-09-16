package api_types

import "mio/internal/pkg/model"

type GetQuestionSubjectListForm struct {
	Scene string `json:"scene" form:"scene" binding:"oneof=home event topic" alias:"banner场景"`
}

type GetQuestionSubjectCreateDTO struct {
	Answer []GetQnrTypeAnswer `json:"answer"`
}

type GetQuestionTypeAnswer struct {
	Id     model.LongID `json:"id"`
	Answer string       `json:"answer"`
}

type QuestionListVo struct {
	Title string       `json:"title"`
	Desc  string       `json:"desc"`
	List  []QuestionVo `json:"list"`
}

type QuestionCategory struct {
	Id    int64  `json:"id"`
	Title string `json:"title"`
	Desc  string `json:"desc"`
}

type QuestionVo struct {
	ID        int64              `json:"id"`
	Title     string             `json:"title"`
	Type      int8               `json:"type"`
	Remind    string             `json:"remind"`
	IsHide    int8               `json:"isHide"`
	Option    []QuestionOptionVO `json:"option"`
	SubjectId model.LongID       `json:"subjectId"`
}

type QuestionOptionVO struct {
	ID             int64   `json:"id"`
	Title          string  `json:"title"`
	Remind         string  `json:"remind"`
	JumpSubject    int64   `json:"jumpSubject"`
	RelatedSubject string  `json:"relatedSubject"`
	Carbon         float64 `json:"carbon"`
}
