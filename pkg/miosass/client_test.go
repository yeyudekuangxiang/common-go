package miosass

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestAutoLogin(t *testing.T) {
	c := Client{
		Domain:    "https://apidev-welfare.miotech.com",
		AppKey:    "w9x8j7Bs",
		AccessKey: "f507664a45b7eb08a3fbf86e02a50ebbcb4d617b",
	}
	loginUrl := c.AutoLogin(AutoLoginParam{
		Name:      "user",
		Credits:   100,
		HeadImg:   "https://avatar.png",
		Uid:       "opjqwpoejqwe",
		Timestamp: 1677470178663,
	})
	wantUrl := "https://apidev-welfare.miotech.com/api/mp2c/login/mio/login?appKey=w9x8j7Bs&credits=100&headImg=https%3A%2F%2Favatar.png&name=user&sign=879930f65938f557e3479d398af64fb7&timestamp=1677470178663&uid=opjqwpoejqwe"
	assert.Equal(t, wantUrl, loginUrl)
}
func TestCertNum(t *testing.T) {
	c := Client{
		Domain:    "https://apidev-welfare.miotech.com",
		AppKey:    "w9x8j7Bs",
		AccessKey: "f507664a45b7eb08a3fbf86e02a50ebbcb4d617b",
	}
	certNumResp, err := c.CertificateCount(CertificateCountParam{
		UserId: "oy_BA5FK1t3dEwrMZndhlUoI2-HY",
	})
	log.Printf("%+v %+v\n", certNumResp, err)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, certNumResp.IsSuccess())
}
func TestExchangeCertificate(t *testing.T) {
	c := Client{
		Domain:    "https://apidev-welfare.miotech.com",
		AppKey:    "w9x8j7Bs",
		AccessKey: "f507664a45b7eb08a3fbf86e02a50ebbcb4d617b",
	}
	exchangeResp, err := c.ExchangeCertificate(ExchangeCertificateParam{
		sign:    sign{},
		SkuId:   1613004259638587394,
		UserId:  "oy_BA5FK1t3dEwrMZndhlUoI2-HY",
		HeadImg: "http://aa.avara.png",
		Name:    "qepqwpeqpwe",
	})
	log.Printf("%+v\n", exchangeResp)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, exchangeResp.IsSuccess())
}
