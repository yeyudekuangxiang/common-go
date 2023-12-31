package validator

import (
	"github.com/gin-gonic/gin/binding"
	zhongwen "github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/pkg/errors"
	"log"
	"reflect"
	"sync"
)

type defaultValidator struct {
	once     sync.Once
	validate *validator.Validate
}

//类型校验
var _ binding.StructValidator = &defaultValidator{}

func (v *defaultValidator) ValidateStruct(obj interface{}) error {

	if kindOfData(obj) == reflect.Struct {

		v.lazyinit()

		if err := v.validate.Struct(obj); err != nil {
			return err
		}
	}

	return nil
}

func (v *defaultValidator) Engine() interface{} {
	v.lazyinit()
	return v.validate
}

func (v *defaultValidator) lazyinit() {
	v.once.Do(func() {
		v.validate = validator.New()
		v.validate.SetTagName("binding")
		zh := zhongwen.New()
		uni := ut.New(zh, zh)
		found := false
		trans, found = uni.GetTranslator("zh")
		if !found {
			log.Println("未找到validator中文翻译")
		}
		err := zh_translations.RegisterDefaultTranslations(v.validate, trans)

		if err != nil {
			log.Fatal(err)
		}
		// add any custom validations etc. here
	})
}

func kindOfData(data interface{}) reflect.Kind {

	value := reflect.ValueOf(data)
	valueType := value.Kind()

	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}

var trans ut.Translator

func NewValidator() binding.StructValidator {
	return new(defaultValidator)
}
func Translate(err error) []string {
	errStrList := make([]string, 0)
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return errStrList
	}

	for _, err := range errs {
		errStrList = append(errStrList, err.Translate(trans))
	}
	return errStrList
}
func IsValidationErrors(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(validator.ValidationErrors)
	return ok
}
func TranslateError(err error) error {
	if IsValidationErrors(err) {
		return errors.New(Translate(err)[0])
	}
	return err
}
