package validator

import (
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
)

var (
	once sync.Once
	ins  *validator.Validate
)

// Instance returns the singleton instance of validator.
func Instance() *validator.Validate {
	if ins == nil {
		once.Do(func() {
			ins = validator.New()
			ins.RegisterTagNameFunc(func(fld reflect.StructField) string {
				strs := strings.SplitN(fld.Tag.Get("json"), ",", 2)
				if len(strs) == 0 {
					return ""
				}

				name := strs[0]

				if name == "-" {
					return ""
				}

				return name
			})
			ins.RegisterValidation("dob", dobValidator)
		})
	}

	return ins
}

func dobValidator(fl validator.FieldLevel) bool {
	layout := "02-01-2006"
	dob, err := time.Parse(layout, fl.Field().String())
	return err == nil && time.Now().After(dob)
}
