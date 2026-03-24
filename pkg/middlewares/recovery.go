package middlewares

import (
	"bytes"
	"fmt"
	"os"
	"runtime"

	"github.com/Paxtiny/oscar/pkg/core"
	"github.com/Paxtiny/oscar/pkg/errs"
	"github.com/Paxtiny/oscar/pkg/log"
	"github.com/Paxtiny/oscar/pkg/utils"
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

// Recovery logs error message when error occurs
func Recovery(c *core.WebContext) {
	defer func() {
		if err := recover(); err != nil {
			stack := stack(3)

			log.ErrorfWithExtra(c, string(stack), "System Error! because %s", err)
			utils.PrintJsonErrorResult(c, errs.ErrSystemError)
		}
	}()

	c.Next()
}

// The following code is from recovery.go of gin

func stack(skip int) []byte {
	buf := new(bytes.Buffer)
	var lines [][]byte
	var lastFile string

	for i := skip; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)

		if file != lastFile {
			data, err := os.ReadFile(file)

			if err != nil {
				continue
			}

			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}

		fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}

	return buf.Bytes()
}

func source(lines [][]byte, n int) []byte {
	n--

	if n < 0 || n >= len(lines) {
		return dunno
	}

	return bytes.TrimSpace(lines[n])
}

func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)

	if fn == nil {
		return dunno
	}

	name := []byte(fn.Name())

	if lastslash := bytes.LastIndex(name, slash); lastslash >= 0 {
		name = name[lastslash+1:]
	}

	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}

	name = bytes.Replace(name, centerDot, dot, -1)

	return name
}
