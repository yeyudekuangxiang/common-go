package starcharge

import (
	"testing"
)

func TestBikeCard(t *testing.T) {
	c := Client{
		Domain:     "https://openapi.hellobike.com/bike/activity",
		Version:    "",
		AESSecret:  "",
		SigSecret:  "",
		Token:      "",
		OperatorID: "",
	}
	//bizId := time.Now().Format("20060102150405") + c.rand()
	resp, _ := c.QueryEquipAuth(QueryEquipAuthParam{
		EquipBizSeq: "",
		ConnectorID: "",
	})

	println(resp.EquipAuthSeq)

}
