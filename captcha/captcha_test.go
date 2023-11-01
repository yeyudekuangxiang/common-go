package captcha

import (
	"log"
	"testing"
)

/*
func TestMyCaptcha(t *testing.T) {
	id, b64s, err := GetCaptchaString()
	if err != nil {
		return
	}
	var _ = b64s
	log.Println("id =", id)
	log.Println("VerifyValue =", store.Get(id, true))


	result := VerifyCaptcha(id, store.Get(id, true))
	log.Println("result =", result)
}
*/

func TestMyCaptcha(t *testing.T) {
	captchaString, err := GetCaptchaString()
	if err != nil {
		return
	}
	if err != nil {
		return
	}
	var _ = captchaString.B64s
	log.Println("id =", captchaString.Id)
	a2 := store.Get(captchaString.Id, true)
	log.Println("VerifyValue =", a2)
	println("a2", a2)

	a1 := store.Get(captchaString.Id, true)
	println("a1", a1)
	result := VerifyCaptcha(captchaString.Id, a1)
	log.Println("result =", result)

	result2 := VerifyCaptcha("eRVja1QEbEnrVM7JMnDY", "uc8m")
	log.Println("result2 =", result2)
}
