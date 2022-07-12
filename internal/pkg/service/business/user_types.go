package business

type CreateBusinessUserParam struct {
	Uid        string `json:"uid"`
	Mobile     string `json:"mobile"`
	BCompanyId int    `json:"bCompanyId"`
}

type CreateBusinessUserImportParam struct {
	Mobile        string `json:"mobile"`
	Name          string `json:"Name"`
	BDepartId     int    `json:"bDepartId"`
	BCompanyId    int    `json:"bCompanyId"`
	TelephoneCode string `json:"telephoneCode"`
}

type GetBusinessUserParam struct {
	Mobile     string `json:"mobile"`
	BCompanyId int    `json:"bCompanyId"`
}
