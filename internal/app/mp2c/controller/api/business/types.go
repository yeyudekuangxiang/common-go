package business

type RegisterBusinessUserForm struct {
	Mobile     string `json:"mobile" form:"mobile" binding:"required" alias:"手机号码"`
	BCompanyId int    `json:"bCompanyId" form:"bCompanyId" binding:"required"  alias:"企业id"`
}
