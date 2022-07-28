package entity

type CertificateStock struct {
	ID            int64  `json:"id"`
	CertificateId string `json:"certificateId"`
	Code          string `json:"code"`
	Used          bool   `json:"used"`
}
