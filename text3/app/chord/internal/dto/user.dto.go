package dto

import (
	"text3/common/dto"
)

/*
	type Example struct {
	    Field  int   `json:"field"  swaggertype:"integer" validate:"required" description:"描述" example:"示例"`
	    Field2 []int `json:"field2" swaggertype:"array,number"` // 字段注释
	}
*/

type GetUserInfoReq struct{}

type GetUserInfoRes struct {
	dto.BaseRes
}
