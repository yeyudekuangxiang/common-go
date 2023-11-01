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
	captchaString, err := GetCaptchaAudio()
	if err != nil {
		return
	}
	if err != nil {
		return
	}
	var _ = captchaString.B64s
	log.Println("id =", captchaString.Id)
	log.Println("VerifyValue =", store.Get(captchaString.Id, true))
	result := VerifyCaptcha(captchaString.Id, store.Get(captchaString.Id, true))
	log.Println("result =", result)
}
