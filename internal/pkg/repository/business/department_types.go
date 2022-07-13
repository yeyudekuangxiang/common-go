package business

type GetDepartmentBy struct {
	ID         int `json:"id"`
	BCompanyId int `json:"bCompanyId"`
}

type GetDepartmentListBy struct {
	Ids []int
}
