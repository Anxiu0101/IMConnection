package model

import (
	"IMConnection/conf"
	"IMConnection/pkg/e"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
)

// Response 基础序列化器
type Response struct {
	Code  int         `json:"code"`
	Data  interface{} `json:"data"`
	Msg   string      `json:"msg"`
	Error string      `json:"error"`
}

//DataList 带有总数的Data结构
type DataList struct {
	Item  interface{} `json:"item"`
	Total uint        `json:"total"`
}

//TokenData 带有token的Data结构
type TokenData struct {
	User         interface{} `json:"user"`
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
}

//TrackedErrorResponse 有追踪信息的错误反应
type TrackedErrorResponse struct {
	Response
	TrackID string `json:"track_id"`
}

// BuildListResponse 带有总数的列表构建器
func BuildListResponse(items interface{}, total uint) Response {
	return Response{
		Code: e.Success,
		Data: DataList{
			Item:  items,
			Total: total,
		},
		Msg: "ok",
	}
}

// ErrorResponse 返回错误信息
func ErrorResponse(err error) Response {
	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, error := range ve {
			field := conf.T(fmt.Sprintf("Field.%s", error.Field))
			tag := conf.T(fmt.Sprintf("Tag.Valid.%s", error.Tag))
			return Response{
				Code:  e.Error,
				Msg:   fmt.Sprintf("%s%s", field, tag),
				Error: fmt.Sprint(err),
			}
		}
	}
	if _, ok := err.(*json.UnmarshalTypeError); ok {
		return Response{
			Code:  e.Error,
			Msg:   "JSON类型不匹配",
			Error: fmt.Sprint(err),
		}
	}
	return Response{
		Code:  e.InvalidParams,
		Msg:   "参数错误",
		Error: fmt.Sprint(err),
	}
}
