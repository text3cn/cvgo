package api

import "github.com/textthree/cvgoweb"

func Index(ctx *httpserver.Context) {
	ctx.Resp.Text("Index")
}
