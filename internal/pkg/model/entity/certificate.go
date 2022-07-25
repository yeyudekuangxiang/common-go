package entity

type CertType string

const (
	CertTypeRandom CertType = "RANDOM"
	CertTypeRule   CertType = "RULE"
	CertTypeStock  CertType = "STOCK"
	CertTypeEmpty  CertType = "EMPTY"
)

type Certificate struct {
	ID            int64    `json:"id"`
	CertificateId string   `json:"certificateId"`
	Message       string   `json:"message"`
	Type          CertType `json:"type"`
	ProductItemId string   `json:"productItemId"`
}
