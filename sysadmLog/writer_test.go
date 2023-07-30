package sysadmLog_test

import (
	"log"
	"net/http"

	"sysadm/sysadmLog"
)

func ExampleLogger_Writer_httpServer() {
	logger := sysadmLog.New()
	w := logger.Writer()
	defer w.Close()

	srv := http.Server{
		// create a stdlib log.Logger that writes to
		// sysadmLog.Logger.
		ErrorLog: log.New(w, "", 0),
	}

	if err := srv.ListenAndServe(); err != nil {
		logger.Fatal(err)
	}
}

func ExampleLogger_Writer_stdlib() {
	logger := sysadmLog.New()
	logger.Formatter = &sysadmLog.JSONFormatter{}

	// Use sysadmLog for standard log output
	// Note that `log` here references stdlib's log
	// Not sysadmLog imported under the name `log`.
	log.SetOutput(logger.Writer())
}
