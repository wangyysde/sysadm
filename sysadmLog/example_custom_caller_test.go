package sysadmLog_test

import (
	"os"
	"path"
	"runtime"
	"strings"

	"sysadm/sysadmLog"
)

func ExampleJSONFormatter_CallerPrettyfier() {
	l := sysadmLog.New()
	l.SetReportCaller(true)
	l.Out = os.Stdout
	l.Formatter = &sysadmLog.JSONFormatter{
		DisableTimestamp: true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			s := strings.Split(f.Function, ".")
			funcname := s[len(s)-1]
			_, filename := path.Split(f.File)
			return funcname, filename
		},
	}
	l.Info("example of custom format caller")
	// Output:
	// {"file":"example_custom_caller_test.go","func":"ExampleJSONFormatter_CallerPrettyfier","level":"info","msg":"example of custom format caller"}
}
