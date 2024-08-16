// 公共请求参数
const configs = {}

configs.server = 'http://localhost:9000'

configs.headers = {
  Authorization: "token",
}

var str
if (console) {
  var _console = {
    log: console.log,
  }
  console.log = function (attr) {
    _console.log(attr)
    str = JSON.stringify(attr, null, 4)
  }
}
