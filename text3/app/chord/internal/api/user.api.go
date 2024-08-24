package api

import (
	"chord/internal/dto"
	"github.com/textthree/cvgoweb"
	"text3/cvgerr"
)

// @Summary    GetUserInfo
// @Tags        user
// @Router      api/user/info [get]
// @Accept      json
// @Produce     json
// @Param       request body dto.GetUserInfoReq true " "
// @Success 	200 {object} dto.GetUserInfoRes
func GetUserInfo(ctx *httpserver.Context) {
	req := dto.GetUserInfoReq{}
	res := dto.GetUserInfoRes{}
	if err := ctx.Req.JsonScan(&req); err != nil {
		ctx.Resp.Json(cvgerr.ParseRequestParamsError())
		return
	}

	ctx.Resp.Json(res)
}
