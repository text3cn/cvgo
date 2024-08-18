package api

import "cvgo/provider/httpserver"

func Index(ctx *httpserver.Context) {
	ctx.Resp.Text("Index")
}
