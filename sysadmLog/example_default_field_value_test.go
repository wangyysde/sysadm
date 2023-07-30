package sysadmLog_test

import (
	"os"

	"sysadm/sysadmLog"
)

type DefaultFieldHook struct {
	GetValue func() string
}

func (h *DefaultFieldHook) Levels() []sysadmLog.Level {
	return sysadmLog.AllLevels
}

func (h *DefaultFieldHook) Fire(e *sysadmLog.Entry) error {
	e.Data["aDefaultField"] = h.GetValue()
	return nil
}

func ExampleDefaultFieldHook() {
	l := sysadmLog.New()
	l.Out = os.Stdout
	l.Formatter = &sysadmLog.TextFormatter{DisableTimestamp: true, DisableColors: true}

	l.AddHook(&DefaultFieldHook{GetValue: func() string { return "with its default value" }})
	l.Info("first log")
	// Output:
	// level=info msg="first log" aDefaultField="with its default value"
}
