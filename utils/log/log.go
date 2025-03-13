package log

import (
	"fmt"
	"log/slog"
	"os"
	"runtime"

	"github.com/VanLavr/Diploma-fin/utils/errors"
)

func ErrorWrapper(err error, errType errors.ErrorType, desc string, params ...any) error {
	var (
		funcName  string
		paramsStr string
		pattern   string
	)

	pc, _, line, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		funcName = details.Name()
	}
	if len(params) != 0 && len(params)%2 == 0 {
		for i := 0; i < len(params); i += 2 {
			paramsStr = fmt.Sprintf("%s[%s]=%v", paramsStr, params[i], params[i+1])
		}
	}
	pattern = fmt.Sprintf(">>func[%s]#%v type[%s]", funcName, line, errType)
	if desc != "" {
		pattern = pattern + fmt.Sprintf(" desc[%s]", desc)
	}
	if paramsStr != "" {
		pattern = pattern + fmt.Sprintf(" params: %s", paramsStr)
	}
	return fmt.Errorf("%s: %w", pattern, err)
}

func GetMethodName() string {
	var funcName string
	pc, _, _, ok := runtime.Caller(1)
	details := runtime.FuncForPC(pc)
	if ok && details != nil {
		funcName = details.Name()
	}

	return funcName + "()"
}

var Logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
	AddSource: true,
	Level:     slog.LevelDebug,
}))
