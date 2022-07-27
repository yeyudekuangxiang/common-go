package entity

type Address struct {
	ID          int64
	AddressId   string
	OpenId      string
	Detail      string
	IsDefault   bool
	PhoneNumber string
	Recipient   string
	Region      string
}
