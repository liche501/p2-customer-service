package logs

import (
	"fmt"
	"runtime"
	"strings"
)

const (
	Black = uint8(iota + 30)
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

var (
	Debug   = CustomerColor{Cyan}   //青色
	Warning = CustomerColor{Yellow} //黄色
	Error   = CustomerColor{Red}    //红色
	Succ    = CustomerColor{Green}
)

type CustomerColor struct {
	ColorNo uint8
}

func (c CustomerColor) Printf(format string, a ...interface{}) {
	s := fmt.Sprintf(format, a...)
	c.doPrint(s)
}
func (c CustomerColor) Println(a ...interface{}) {
	s := fmt.Sprint(a...)
	c.doPrint(s)
}

func (c CustomerColor) doPrint(a interface{}) {
	// isDebug := os.Getenv("WXSHOPDEBUG")
	// if len(isDebug) == 0 && c.ColorNo != Red {
	// 	return
	// }

	funcName, file, line, ok := runtime.Caller(2)

	if ok {
		base := "/best/wxshop"
		i := strings.Index(file, base)
		shoutFile := file[i+len(base):]

		fileInfo := fmt.Sprintf("-->%v:%d %s", shoutFile, line, runtime.FuncForPC(funcName).Name())
		colorFileInfo := fmt.Sprintf("\033[%dm%s\033[0m", Magenta, fileInfo)
		fmt.Printf("\033[%dm%v\033[0m", c.ColorNo, a)
		fmt.Println(colorFileInfo)

	} else {
		colorFileInfo := fmt.Sprintf("\033[%dm%v\033[0m", Red, "-->(Unknown debug info)")
		fmt.Printf("\033[%dm%v\033[0m", c.ColorNo, a)
		fmt.Println(colorFileInfo)
	}
}
