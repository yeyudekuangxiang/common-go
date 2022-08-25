package api_types

type GetSubjectListForm struct {
	Scene string `json:"scene" form:"scene" binding:"oneof=home event topic" alias:"banner场景"`
}

type GetQnrSubjectCreateDTO struct {
	Answer string `json:"answer"`
}

type QnrListVo struct {
	Title string  `json:"title"`
	List  []QnrVo `json:"list"`
}

type QnrVo struct {
	ID     int64      `json:"id"`
	Title  string     `json:"title"`
	Type   int8       `json:"type"`
	Remind string     `json:"remind"`
	IsHide int8       `json:"isHide"`
	Option []OptionVO `json:"option"`
}

type OptionVO struct {
	ID             int64  `json:"id"`
	Title          string `json:"title"`
	Remind         string `json:"remind"`
	JumpSubject    int64  `json:"jumpSubject"`
	RelatedSubject string `json:"relatedSubject"`
}
