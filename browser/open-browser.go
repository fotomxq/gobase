package browser

import (
	"runtime"
	"os/exec"
	"github.com/pkg/errors"
)

var OpenBrowserCommands = map[string]string{
	"windows": "cmd /c start",
	"darwin":  "open",
	"linux":   "xdg-open",
}

//启动各类平台浏览器，并打开URL
// 如果失败则返回错误
//param url string  URL地址
//return error 错误
func Open(url string) error {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = errors.New("Unknown system type.")
	}
	if err != nil {
	}
	return err
}