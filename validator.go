package verify

import (
	"errors"

	cn "github.com/go-playground/locales/zh_Hans_CN"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	cn_translations "github.com/go-playground/validator/v10/translations/zh"
)

func NewValidator() *Validator {
	cn_translator := cn.New()
	obj := &Validator{
		uni:      ut.New(cn_translator, cn_translator), //nolint:typecheck
		validate: validator.New(),
	}
	obj.trans, _ = obj.uni.GetTranslator("zh")
	if err := cn_translations.RegisterDefaultTranslations(obj.validate, obj.trans); err != nil { //nolint:typecheck
		panic(err)
	}
	obj.register()
	return obj
}

type Validator struct {
	trans    ut.Translator
	uni      *ut.UniversalTranslator
	validate *validator.Validate
}

func (c *Validator) register() {
	if err := c.validate.RegisterValidation("phone", ValidatePhone); err != nil {
		panic(err)
	}
	if err := c.validate.RegisterTranslation("phone", c.trans, func(ut ut.Translator) error {
		return ut.Add("phone", "{0} 手机号码有误", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("phone", fe.Field())
		return t
	}); err != nil {
		panic(err)
	}

	if err := c.validate.RegisterValidation("objectId", ValidateObjectID); err != nil {
		panic(err)
	}
	if err := c.validate.RegisterTranslation("objectId", c.trans, func(ut ut.Translator) error {
		return ut.Add("objectId", "{0} 不符合ObjectID格式", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("objectId", fe.Field())
		return t
	}); err != nil {
		panic(err)
	}
}

func (c *Validator) Validate(i interface{}) error {
	var errString string
	if es, ok := c.validate.Struct(i).(validator.ValidationErrors); ok {
		for _, e := range es {
			errString += " " + e.Translate(c.trans)
		}
		return errors.New(errString)
	}
	return nil
}

func ValidatePhone(fl validator.FieldLevel) bool {
	return IsValidMobilePhoneNumber(fl.Field().String())
}

func ValidateObjectID(fl validator.FieldLevel) bool {
	return IsObjectID(fl.Field().String())
}
