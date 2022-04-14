package v1

import (
	"IMConnection/conf"
	"IMConnection/model"
	"IMConnection/pkg/e"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
)

// ErrorResponse 返回错误信息
func ErrorResponse(err error) model.Response {
	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, error := range ve {
			field := conf.T(fmt.Sprintf("Field.%s", error.Field))
			tag := conf.T(fmt.Sprintf("Tag.Valid.%s", error.Tag))
			return model.Response{
				Code:  e.Error,
				Msg:   fmt.Sprintf("%s%s", field, tag),
				Error: fmt.Sprint(err),
			}
		}
	}
	if _, ok := err.(*json.UnmarshalTypeError); ok {
		return model.Response{
			Code:  e.Error,
			Msg:   "JSON类型不匹配",
			Error: fmt.Sprint(err),
		}
	}
	return model.Response{
		Code:  e.InvalidParams,
		Msg:   "参数错误",
		Error: fmt.Sprint(err),
	}
}
