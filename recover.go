package lee

import (
	"fmt"
	"runtime"
	"strings"
)

func trace(messages string)string  {
	var pcs [32]uintptr
	n := runtime.Callers(3,pcs[:])
	var str strings.Builder
	str.WriteString(messages+"\nTraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}


