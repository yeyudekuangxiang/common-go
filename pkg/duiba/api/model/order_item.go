package model

type OrderItem struct {
	Title               string `json:"title" form:"title" alias:"title" binding:"required"`
	IsSelf              string `json:"isSelf" form:"isSelf" alias:"isSelf"`
	SmallImage          string `json:"smallImage" form:"smallImage" alias:"smallImage"`
	MerchantCode        string `json:"merchantCode" form:"merchantCode" alias:"merchantCode"`
	PerCredit           IntStr `json:"perCredit" form:"perCredit" alias:"perCredit"`
	PerPrice            string `json:"perPrice" form:"perPrice" alias:"perPrice"`
	Quantity            IntStr `json:"quantity" form:"quantity" alias:"quantity"`
	Code                string `json:"code" form:"code" alias:"code"`
	Password            string `json:"password" form:"password" alias:"password"`
	CardBeginTime       string `json:"cardBeginTime" form:"cardBeginTime" alias:"cardBeginTime"`
	CardEndTime         string `json:"cardEndTime" form:"cardEndTime" alias:"cardEndTime"`
	DeliveryCompanyNo   string `json:"deliveryCompanyNo" form:"deliveryCompanyNo" alias:"deliveryCompanyNo"`
	DeliveryCompanyName string `json:"deliveryCompanyName" form:"deliveryCompanyName" alias:"deliveryCompanyName"`
	DuibaSupplyPrice    string `json:"duibaSupplyPrice" form:"duibaSupplyPrice" alias:"duibaSupplyPrice"`
}
