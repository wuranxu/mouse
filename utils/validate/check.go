package validate

import (
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	translations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/wuranxu/mouse/exception"
	"log"
	"reflect"
	"strings"
)

var (
	check *validator.Validate
	trans ut.Translator
)

func Check(data interface{}, msg exception.Err) error {
	if err := check.Struct(data); err != nil {
		errs := err.(validator.ValidationErrors)
		translate := errs.Translate(trans)
		for _, v := range translate {
			return msg.New(v)
		}
	}
	return nil
}

func init() {
	cn := zh.New()
	uni := ut.New(cn, cn)

	// this is usually know or extracted from http 'Accept-Language' header
	// also see uni.FindTranslator(...)
	trans, _ = uni.GetTranslator("zh")
	check = validator.New()
	check.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})

	err := translations.RegisterDefaultTranslations(check, trans)
	if err != nil {
		log.Fatal("validate translate error: ", err)
	}
}
