package cvgerr

var Success = 0
var Fail = 1

var AuthorizationFailed = NewApiError(1000, "Authorization failed.")
var ParseRequestParamsFailed = NewApiError(1001, "Parse request params failed.")
var SQL_ERR = 2001
