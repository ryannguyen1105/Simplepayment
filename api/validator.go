package api

import (
	"github.com/go-playground/validator/v10"
	"github.com/ryannguyen1105/Simplepayment/util"
)

var validStatus validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if status, ok := fieldLevel.Field().Interface().(string); ok {
		return util.IsSupportedStatus(status)
	}
	return false
}
