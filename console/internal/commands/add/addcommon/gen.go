package addcommon

import (
	"cvgo/kit/filekit"
	"cvgo/kit/strkit"
)

func GenApidebug(apiDebugHtmlFilePath, requestPath, method string) {
	content := `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<script type="text/javascript" src="../jquery-3.6.0.min.js"></script>
	<script type="text/javascript" src="../config.js"></script>
</head>
<script>
let method = "` + strkit.Strtoupper(method) + `"
let url = configs.server + "/` + requestPath + `"
let data = {

}

$.ajax({
  type     : method,
  contentType: 'application/json',
  url      : url ,
  data     : method == "GET" ? data : JSON.stringify(data),
  headers  : configs.headers,
  async: false,
  success  :function(res) {
    console.log(res)
  },
  error    : function(err) {
    console.log(err)
  }
});
</script>

<body id="">
	<script>document.write(str)</script>
</body>
</html>
`
	filekit.FilePutContents(apiDebugHtmlFilePath, content)
}
