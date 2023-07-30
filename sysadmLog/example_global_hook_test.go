package sysadmLog_test

import (
	"os"

	"sysadm/sysadmLog"
)

var (
	mystring string
)

type GlobalHook struct {
}

func (h *GlobalHook) Levels() []sysadmLog.Level {
	return sysadmLog.AllLevels
}

func (h *GlobalHook) Fire(e *sysadmLog.Entry) error {
	e.Data["mystring"] = mystring
	return nil
}

func ExampleGlobalHook() {
	l := sysadmLog.New()
	l.Out = os.Stdout
	l.Formatter = &sysadmLog.TextFormatter{DisableTimestamp: true, DisableColors: true}
	l.AddHook(&GlobalHook{})
	mystring = "first value"
	l.Info("first log")
	mystring = "another value"
	l.Info("second log")
	// Output:
	// level=info msg="first log" mystring="first value"
	// level=info msg="second log" mystring="another value"
}
