package entity

type Certificate struct {
	ID            int64  `json:"id"`
	CertificateId string `json:"certificateId"`
	Message       string `json:"message"`
	Type          string `json:"type"`
	ProductItemId string `json:"productItemId"`
}
