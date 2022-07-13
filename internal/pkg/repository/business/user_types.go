package business

type GetUserBy struct {
	Mobile string //手机号精确匹配
	Uid    string
	ID     int64
}

type GetUserListBy struct {
	Ids []int64
	CId int
}

type GetUserImportBy struct {
	Mobile     string `json:"mobile"`
	BCompanyId int    `json:"bCompanyId"`
}
