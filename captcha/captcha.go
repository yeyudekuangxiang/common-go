package captcha

import (
	"github.com/mojocn/base64Captcha"
	"image/color"
)

//configJsonBody json request body.
type configJsonBody struct {
	Id            string
	CaptchaType   string
	VerifyValue   string
	DriverAudio   *base64Captcha.DriverAudio
	DriverString  *base64Captcha.DriverString
	DriverChinese *base64Captcha.DriverChinese
	DriverMath    *base64Captcha.DriverMath
	DriverDigit   *base64Captcha.DriverDigit
}

var store = base64Captcha.DefaultMemStore

// https://captcha.mojotv.cn/ 调试配置
// base64Captcha create return id, b64s, err

type Captcha struct {
	B64s string
	Id   string
	Code string
}

func GetCaptchaString() (*Captcha, error) {
	var param = configJsonBody{
		Id:          "",
		CaptchaType: "string",
		VerifyValue: "",
		DriverAudio: &base64Captcha.DriverAudio{},
		DriverString: &base64Captcha.DriverString{
			Length:          4,
			Height:          60,
			Width:           240,
			ShowLineOptions: 2,
			NoiseCount:      0,
			Source:          "1234567890qwertyuioplkjhgfdsazxcvbnm",
		},
		DriverChinese: &base64Captcha.DriverChinese{},
		DriverMath:    &base64Captcha.DriverMath{},
		DriverDigit:   &base64Captcha.DriverDigit{},
	}
	return getCaptcha(param)
}

func GetCaptchaChinese() (*Captcha, error) {
	var param = configJsonBody{
		Id:           "",
		CaptchaType:  "chinese",
		VerifyValue:  "",
		DriverAudio:  &base64Captcha.DriverAudio{},
		DriverString: &base64Captcha.DriverString{},
		DriverChinese: &base64Captcha.DriverChinese{
			Height:          60,
			Width:           320,
			ShowLineOptions: 0,
			NoiseCount:      0,
			Source:          "设想,你在,处理,消费者,的音,频输,出音,频可,能无,论什,么都,没有,任何,输出,或者,它可,能是,单声道,立体声,或是,环绕立,体声的,,不想要,的值",
			Length:          2,
			Fonts:           []string{"wqy-microhei.ttc"},
			BgColor:         &color.RGBA{R: 125, G: 125, B: 0, A: 118},
		},
		DriverMath:  &base64Captcha.DriverMath{},
		DriverDigit: &base64Captcha.DriverDigit{},
	}
	return getCaptcha(param)
}

func GetCaptchaMath() (*Captcha, error) {
	var param = configJsonBody{
		Id:            "",
		CaptchaType:   "math",
		VerifyValue:   "",
		DriverAudio:   &base64Captcha.DriverAudio{},
		DriverString:  &base64Captcha.DriverString{},
		DriverChinese: &base64Captcha.DriverChinese{},
		DriverMath: &base64Captcha.DriverMath{
			Height:          60,
			Width:           240,
			ShowLineOptions: 0,
			NoiseCount:      0,
			Fonts:           []string{"wqy-microhei.ttc"},
			BgColor:         &color.RGBA{R: 0, G: 0, B: 0, A: 0},
		},
		DriverDigit: &base64Captcha.DriverDigit{},
	}
	return getCaptcha(param)
}

func GetCaptchaAudio() (*Captcha, error) {
	var param = configJsonBody{
		Id:          "",
		CaptchaType: "audio",
		VerifyValue: "",
		DriverAudio: &base64Captcha.DriverAudio{
			Length:   6,
			Language: "zh",
		},
		DriverString:  &base64Captcha.DriverString{},
		DriverChinese: &base64Captcha.DriverChinese{},
		DriverMath:    &base64Captcha.DriverMath{},
		DriverDigit:   &base64Captcha.DriverDigit{},
	}
	return getCaptcha(param)
}

func GetCaptchaDigit() (*Captcha, error) {
	var param = configJsonBody{
		Id:            "",
		CaptchaType:   "digit",
		VerifyValue:   "",
		DriverAudio:   &base64Captcha.DriverAudio{},
		DriverString:  &base64Captcha.DriverString{},
		DriverChinese: &base64Captcha.DriverChinese{},
		DriverMath:    &base64Captcha.DriverMath{},
		DriverDigit: &base64Captcha.DriverDigit{
			Height:   80,
			Width:    240,
			Length:   5,
			MaxSkew:  0.7,
			DotCount: 80,
		},
	}
	return getCaptcha(param)
}

func getCaptcha(param configJsonBody) (*Captcha, error) {

	// 	{
	// 		ShowLineOptions: [],
	// 		CaptchaType: "string",
	// 		Id: '',
	// 		VerifyValue: '',
	// 		DriverAudio: {
	// 			Length: 6,
	// 			Language: 'zh'
	// 		},
	// 		DriverString: {
	// 			Height: 60,
	// 			Width: 240,
	// 			ShowLineOptions: 0,
	// 			NoiseCount: 0,
	// 			Source: "1234567890qwertyuioplkjhgfdsazxcvbnm",
	// 			Length: 6,
	// 			Fonts: ["wqy-microhei.ttc"],
	// 			BgColor: {R: 0, G: 0, B: 0, A: 0},
	// 		},
	// 		DriverMath: {
	// 			Height: 60,
	// 			Width: 240,
	// 			ShowLineOptions: 0,
	// 			NoiseCount: 0,
	// 			Length: 6,
	// 			Fonts: ["wqy-microhei.ttc"],
	// 			BgColor: {R: 0, G: 0, B: 0, A: 0},
	// 		},
	// 		DriverChinese: {
	// 			Height: 60,
	// 			Width: 320,
	// 			ShowLineOptions: 0,
	// 			NoiseCount: 0,
	// 			Source: "设想,你在,处理,消费者,的音,频输,出音,频可,能无,论什,么都,没有,任何,输出,或者,它可,能是,单声道,立体声,或是,环绕立,体声的,,不想要,的值",
	// 			Length: 2,
	// 			Fonts: ["wqy-microhei.ttc"],
	// 			BgColor: {R: 125, G: 125, B: 0, A: 118},
	// 		},
	// 		DriverDigit: {
	// 			Height: 80,
	// 			Width: 240,
	// 			Length: 5,
	// 			MaxSkew: 0.7,
	// 			DotCount: 80
	// 		}
	// 	},
	// 	blob: "",
	// 	loading: false
	// }

	var driver base64Captcha.Driver

	//create base64 encoding captcha
	switch param.CaptchaType {
	case "audio":
		driver = param.DriverAudio
	case "string":
		driver = param.DriverString.ConvertFonts()
	case "math":
		driver = param.DriverMath.ConvertFonts()
	case "chinese":
		driver = param.DriverChinese.ConvertFonts()
	default:
		driver = param.DriverDigit
	}
	c := base64Captcha.NewCaptcha(driver, store)
	id, VerifyValue, err := c.Generate()

	if err != nil {
		return nil, err
	}
	return &Captcha{
		B64s: VerifyValue,
		Id:   id,
		Code: store.Get(id, true),
	}, nil
}

// base64Captcha verify

func VerifyCaptcha(id, VerifyValue string) bool {
	return store.Verify(id, VerifyValue, true)
}
