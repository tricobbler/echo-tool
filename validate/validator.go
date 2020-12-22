package validate

import (
	zhongwen "github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
	"strings"
)

var (
	V     *validator.Validate
	trans ut.Translator
)

type (
	validateErrList []string
	customValidator struct {
		validator *validator.Validate
	}
)

func (cv *customValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func NewCustomValidator() *customValidator {
	return &customValidator{validator: V}
}

func init() {
	zh := zhongwen.New()
	uni := ut.New(zh, zh)
	trans, _ = uni.GetTranslator("zh")

	V = validator.New()
	V.RegisterTagNameFunc(func(field reflect.StructField) string {
		label := field.Tag.Get("label")
		if label == "" {
			return field.Name
		}
		return label
	})
	zh_translations.RegisterDefaultTranslations(V, trans)
}

func Translate(errs error) validateErrList {
	var errList validateErrList
	for _, e := range errs.(validator.ValidationErrors) {
		// can translate each error one at a time.
		errList = append(errList, e.Translate(trans))
	}
	return errList
}

func (v validateErrList) One() string {
	if len(v) == 0 {
		return ""
	}

	return v[0]
}

func (v validateErrList) All() string {
	return strings.Join(v, "|")
}
