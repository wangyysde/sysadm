/* =============================================================
* @Author:  Wayne Wang <net_use@bzhy.com>
*
* @Copyright (c) 2023 Bzhy Network. All rights reserved.
* @HomePage http://www.sysadm.cn
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at:
* http://www.apache.org/licenses/LICENSE-2.0
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and  limitations under the License.
* @License GNU Lesser General Public License  https://www.sysadm.cn/lgpl.html
 */

package sysadmLog

import (
	"fmt"
	"github.com/wangyysde/sysadmServer"
	"os"
	"strings"
)

type LogFields map[string]interface{}

type LoggerConfig struct {
	// DefaultLogger is a instance of logger. We use defaultLogger to log log message
	// when AccessLogger  or ErrorLogger is nil
	DefaultLogger *Logger

	// AccessLogger is a instance of logrus
	// and this instance if for logging access log
	AccessLogger *Logger

	// ErrorLogger is a instance of logrus
	// and this instance if for logging error log
	ErrorLogger *Logger

	// Kind specifies the format of the log where be log to
	// kind is one of text or json
	Kind string

	// AccessLogFile records the path of log file for access
	// if the access log and error log log into difference files
	AccessLogFile string

	// file descriptor for access logger
	accessFP *os.File

	// ErrorLogFile records the path of log file for error
	// if the access log and error log log into difference files
	ErrorLogFile string

	// file descriptor for error logger
	errorFP *os.File

	// Level specifies which level log will be logged
	Level string

	// SplitAccessAndError identify if log access log and error log
	// into difference io.Writer.
	// Logs will be log into  AccessLogger if SplitAccessAndError is false
	// otherwise access logs  will be log into AccessLogger and error logs will be log into ErrorLogger.
	SplitAccessAndError bool

	// Specifies the format of the log timestamp,like: "2021/09/02 - 15:04:05"
	TimeStampFormat string

	// Flag for whether to log caller info (off by default)
	ReportCaller bool

	// SkipPaths is a url path array which logs are not written.
	// Optional.
	SkipPaths []string
}

type Sysadmerror struct {
	ErrorNo    int
	ErrorLevel int
	ErrorMsg   string
}

var Levels = []string{
	"trace",
	"debug",
	"info",
	"warning",
	"error",
	"fatal",
	"panic",
}

var TimestampFormat = map[string]string{
	"ANSIC":       "Mon Jan _2 15:04:05 2006",
	"UnixDate":    "Mon Jan _2 15:04:05 MST 2006",
	"RubyDate":    "Mon Jan 02 15:04:05 -0700 2006",
	"RFC822":      "02 Jan 06 15:04 MST",
	"RFC822Z":     "02 Jan 06 15:04 -0700", // 使用数字表示时区的RFC822
	"RFC850":      "Monday, 02-Jan-06 15:04:05 MST",
	"RFC1123":     "Mon, 02 Jan 2006 15:04:05 MST",
	"RFC1123Z":    "Mon, 02 Jan 2006 15:04:05 -0700", // 使用数字表示时区的RFC1123
	"RFC3339":     "2006-01-02T15:04:05Z07:00",
	"RFC3339Nano": "2006-01-02T15:04:05.999999999Z07:00",
	"Kitchen":     "3:04PM",
	"Stamp":       "Jan _2 15:04:05",
	"StampMilli":  "Jan _2 15:04:05.000",
	"StampMicro":  "Jan _2 15:04:05.000000",
	"StampNano":   "Jan _2 15:04:05.000000000",
	"DateTime":    "2006-01-02 15:04:05",
}

// NewSysadmLogger new a new instance of sysadmLogger configuraton for webserver
func NewSysadmLogger() *LoggerConfig {
	defaultLogger := New()
	defaultLogger.Out = sysadmServer.DefaultWriter

	logger := &LoggerConfig{
		DefaultLogger:       defaultLogger,
		AccessLogger:        nil,
		ErrorLogger:         nil,
		Kind:                "text",
		AccessLogFile:       "",
		accessFP:            nil,
		ErrorLogFile:        "",
		errorFP:             nil,
		Level:               "debug",
		SplitAccessAndError: false,
		TimeStampFormat:     TimestampFormat["RFC3339"],
		ReportCaller:        true,
		SkipPaths:           nil,
	}
	logger.SetLoggerLevel("")
	logger.SetLoggerKind("")

	return logger
}

// define text formatter with default values
var textFormatter = &TextFormatter{
	ForceColors:               false,
	DisableColors:             false,
	ForceQuote:                false,
	DisableQuote:              true,
	EnvironmentOverrideColors: true,
	DisableTimestamp:          false,
	FullTimestamp:             true,
	TimestampFormat:           TimestampFormat["RFC3339"],
	DisableSorting:            true,
	DisableLevelTruncation:    true,
	PadLevelText:              false,
	QuoteEmptyFields:          true,
}

// define json formatter with default values
var jsonFormatter = &JSONFormatter{
	TimestampFormat:   TimestampFormat["RFC3339"],
	DisableTimestamp:  false,
	DisableHTMLEscape: true,
}

// SetLogLevel  set the value  of LoggerConfig.Level to "debug" if the value of it is ""
// Then set the the levels for DefaultLoger, AccessLogger and ErrorLogger
// This function should be called during LoggerConfig.Level is setting and a new logger is initating.
func (l *LoggerConfig) SetLoggerLevel(lvl string) {
	lvl = strings.TrimSpace(lvl)

	if _, err := ParseLevel(l.Level); err == nil {
		l.Level = lvl
	} else {
		if strings.TrimSpace(l.Level) == "" {
			l.Level = "debug"
		}
	}

	loggerLevel, _ := ParseLevel(l.Level)

	if l.DefaultLogger != nil {
		l.DefaultLogger.SetLevel(loggerLevel)
	}

	if l.AccessLogger != nil {
		l.AccessLogger.SetLevel(loggerLevel)
	}

	if l.ErrorLogger != nil {
		l.ErrorLogger.SetLevel(loggerLevel)
	}

}

// setLoggerKind  set the value  of l.Kind to "text" if the value of it is ""
// Then setLoggerKind  sets the the formatter for DefaultLoger, AccessLogger and ErrorLogger
//
//	This function should be called during l.kind is setting and a new logger is initating.
func (l *LoggerConfig) SetLoggerKind(kind string) {
	kind = strings.TrimSpace(strings.ToLower(kind))
	if kind == "text" || kind == "json" {
		l.Kind = kind
	}

	if strings.TrimSpace(l.Kind) == "" {
		l.Kind = "text"
	}

	if logger := l.DefaultLogger; logger != nil {
		if strings.ToLower(l.Kind) == "text" {
			formatter := *textFormatter
			formatter.DisableColors = false
			logger.SetFormatter(&formatter)
		} else {
			formatter := *jsonFormatter
			logger.SetFormatter(&formatter)
		}
	}

	if logger := l.AccessLogger; logger != nil {
		if strings.ToLower(l.Kind) == "text" {
			logger.SetFormatter(textFormatter)
		} else {
			logger.SetFormatter(jsonFormatter)
		}
	}

	if logger := l.ErrorLogger; logger != nil {
		if strings.ToLower(l.Kind) == "text" {
			logger.SetFormatter(textFormatter)
		} else {
			logger.SetFormatter(jsonFormatter)
		}
	}
}

// CloseAccessLogger close access logger out file descriptor and set access logger to nil
func (l *LoggerConfig) CloseAccessLogger() {
	if l.accessFP != nil {
		fp := l.accessFP
		_ = fp.Close()
		l.accessFP = nil
	}
	l.AccessLogger = nil
}

// CloseErrorLogger close error logger out file descriptor and set error logger to nil
func (l *LoggerConfig) CloseErrorLogger() {
	if l.errorFP != nil {
		fp := l.errorFP
		_ = fp.Close()
		l.errorFP = nil
	}
	l.ErrorLogger = nil
}

// SetAccessLogFile set file to access log file and then open the file for access logger output.
// the caller should call l.CloseAccessLogger after called this method.
func (l *LoggerConfig) SetAccessLogFile(file string) error {
	if strings.TrimSpace(file) == "" {
		return fmt.Errorf("The length of access log file should be larger 1")
	}

	file = strings.TrimSpace(file)
	fp, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return fmt.Errorf("Open log file %s error: %s", file, fmt.Sprintf("%s", err))
	}

	logger := l.AccessLogger
	if logger == nil {
		logger = New()
		l.AccessLogger = logger
		l.SetLoggerLevel("")
		l.SetLoggerKind("")
	}
	logger.Out = fp
	oldFp := l.accessFP
	if fp != nil {
		_ = oldFp.Close()
	}
	l.accessFP = fp

	return nil
}

// SetErrorLogFile set file to error log file and then open the file for error logger output.
// the caller should call l.CloseErrorLogger after called this method.
func (l *LoggerConfig) SetErrorLogFile(file string) error {
	if strings.TrimSpace(file) == "" {
		return fmt.Errorf("The length of error log file should be larger 1")
	}

	file = strings.TrimSpace(file)
	fp, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return fmt.Errorf("Open log file %s error: %s", file, fmt.Sprintf("%s", err))
	}

	logger := l.ErrorLogger
	if logger == nil {
		logger = New()
		l.ErrorLogger = logger
		l.SetLoggerLevel("")
		l.SetLoggerKind("")
	}
	logger.Out = fp
	oldFp := l.errorFP
	if fp != nil {
		_ = oldFp.Close()
	}
	l.errorFP = fp

	return nil
}

// SetAccessLoggerWithFp set fp(writer) to access logger
// the function which close the fp should be called by defer following this function calling
func (l *LoggerConfig) SetAccessLoggerWithFp(fp *os.File) error {
	if fp == nil {
		return fmt.Errorf("can not set nil(writer) to logger")
	}

	logger := l.AccessLogger
	if logger == nil {
		logger = New()
		l.AccessLogger = logger
		l.SetLoggerLevel("")
		l.SetLoggerKind("")
	}
	logger.Out = fp
	l.accessFP = fp

	return nil
}

// SetErrorLoggerWithFp set fp(writer) to error logger
// the function which close the fp should be called by defer following this function calling
func (l *LoggerConfig) SetErrorLoggerWithFp(fp *os.File) error {
	if fp == nil {
		return fmt.Errorf("can not set nil(writer) to logger")
	}

	logger := l.ErrorLogger
	if logger == nil {
		logger = New()
		l.ErrorLogger = logger
		l.SetLoggerLevel("")
		l.SetLoggerKind("")
	}
	logger.Out = fp
	l.errorFP = fp

	return nil
}

// SetIsSplitLog set IsSplitLog  to Logger configuration
func (l *LoggerConfig) SetIsSplitLog(isSplit bool) error {
	if isSplit {
		if l.AccessLogger == nil || l.ErrorLogger == nil {
			return fmt.Errorf("you try to set SplitAccessAndError to true, but access log or error logger have not opened")
		}
	} else {
		if l.AccessLogger == nil {
			return fmt.Errorf("you try to set SplitAccessAndError to false, but access logger have not opened. All log message will be log to defaultOutput and error logger")
		}
	}
	l.SplitAccessAndError = isSplit

	return nil
}

// SetTimestampFormat set timeStampFormat to the LoggerConfig and then set it to all the Loggers
// The value of format is one of constants of time package and "DateTime"
func (l *LoggerConfig) SetTimestampFormat(format string) error {
	if strings.TrimSpace(format) == "" {
		return fmt.Errorf("The length of format should be larger 1")
	}

	for k, v := range TimestampFormat {
		if strings.ToLower(k) == strings.ToLower(format) {
			l.TimeStampFormat = v
			textFormatter.TimestampFormat = v
			jsonFormatter.TimestampFormat = v
			l.SetLoggerKind("")
			return nil
		}
	}

	return fmt.Errorf("The TimeStampFormat(%s) is invalid.", format)
}

// if isDisable is false, then timestamp message will be added to log messages automatically.
// Otherwise no timestamp will be added.
func (l *LoggerConfig) DisableTimestamp(isDisable bool) {
	textFormatter.DisableTimestamp = isDisable
	jsonFormatter.DisableTimestamp = isDisable
	l.SetLoggerKind("")
}

// SetReportCaller sets ReportCaller of LoggerConfig to true or false.
// if the value of ReportCaller is true, then callers name and the file path information which the caller in will be
// added to log messages.
func (l *LoggerConfig) SetReportCaller(isReportCaller bool) {
	l.ReportCaller = isReportCaller

	if l.DefaultLogger != nil {
		l.DefaultLogger.ReportCaller = isReportCaller
	}

	if l.AccessLogger != nil {
		l.AccessLogger.ReportCaller = isReportCaller
	}

	if l.ErrorLogger != nil {
		l.ErrorLogger.ReportCaller = isReportCaller
	}

}

// WriteLog write message to the logger. logLevel will be set to "error" when its value is ""
func (l *LoggerConfig) WriteLog(logger *Logger, message string, logLevel string) {
	if logger == nil {
		return
	}

	if strings.TrimSpace(logLevel) == "" {
		logLevel = "error"
	}

	found := false

	for _, value := range Levels {
		if strings.ToLower(strings.TrimSpace(logLevel)) == value {
			found = true
		}
	}

	if !found {
		logLevel = "error"
	}

	switch strings.ToLower(logLevel) {
	case "trace":
		logger.Trace(message)
	case "debug":
		logger.Debug(message)
	case "info":
		logger.Info(message)
	case "warning":
		logger.Warn(message)
	case "error":
		logger.Error(message)
	case "fatal":
		logger.Fatal(message)
	case "panic":
		logger.Panic(message)
		panic("")
	}
}

// WriteLogWithFields write fields to the logger. logLevel will be set to "info" when its value is ""
// the value of logLevel will be set to "error" if fields["ErrorMessage"] is not ""
// fields["ErrorMessage"] will write to logger when  fields["ErrorMessage"] and  fields["Message"] are not ""
// the value of fields["Message"] will be igored.
func (l *LoggerConfig) WriteLogWithFields(fields LogFields, logLevel string) {
	if strings.TrimSpace(logLevel) == "" {
		logLevel = "info"
	}

	found := false
	for _, value := range Levels {
		if strings.ToLower(strings.TrimSpace(logLevel)) == value {
			found = true
		}
	}

	errormsg := ""
	if v, ok := fields["ErrorMessage"]; ok {
		errormsg = v.(string)
	}

	if strings.TrimSpace(errormsg) != "" && !found {
		logLevel = "error"
	} else if !found {
		logLevel = "info"
	}

	switch strings.ToLower(logLevel) {
	case "trace":
		logger := l.DefaultLogger
		if logger != nil {
			logger.WithFields(Fields(fields)).Trace("")
		}
		logger = l.AccessLogger
		if logger != nil {
			logger.WithFields(Fields(fields)).Trace("")
		}
	case "debug":
		logger := l.DefaultLogger
		if logger != nil {
			logger.WithFields(Fields(fields)).Debug("")
		}
		logger = l.AccessLogger
		if logger != nil {
			logger.WithFields(Fields(fields)).Debug("")
		}
	case "info":
		logger := l.DefaultLogger
		if logger != nil {
			logger.WithFields(Fields(fields)).Info("")
		}
		logger = l.AccessLogger
		if logger != nil {
			logger.WithFields(Fields(fields)).Info("")
		}
	case "warning":
		logger := l.DefaultLogger
		if logger != nil {
			logger.WithFields(Fields(fields)).Warn("")
		}
		logger = l.ErrorLogger
		if logger != nil {
			logger.WithFields(Fields(fields)).Warn("")
		}
	case "error":
		logger := l.DefaultLogger
		if logger != nil {
			logger.WithFields(Fields(fields)).Error("")
		}
		logger = l.ErrorLogger
		if logger != nil {
			logger.WithFields(Fields(fields)).Error("")
		}
	case "fatal":
		logger := l.DefaultLogger
		if logger != nil {
			logger.WithFields(Fields(fields)).Fatal("")
		}
		logger = l.ErrorLogger
		if logger != nil {
			logger.WithFields(Fields(fields)).Fatal("")
		}
	case "panic":
		logger := l.DefaultLogger
		if logger != nil {
			logger.WithFields(Fields(fields)).Panic("")
		}
		logger = l.ErrorLogger
		if logger != nil {
			logger.WithFields(Fields(fields)).Panic("")
		}
		panic("")
	}
}

// build a handlerFunc for sysadmServer to log access log .
func (l *LoggerConfig) BuildWebServerLogger() sysadmServer.HandlerFunc {

	notlogged := l.SkipPaths
	var skip map[string]struct{}

	if length := len(notlogged); length > 0 {
		skip = make(map[string]struct{}, length)
		for _, path := range notlogged {
			skip[path] = struct{}{}
		}
	}

	return func(c *sysadmServer.Context) {
		fields := make(LogFields)
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log only when path is not being skipped
		if _, ok := skip[path]; !ok {

			fields["Request"] = c.Request
			fields["Keys"] = c.Keys
			fields["ClientIP"] = c.ClientIP()
			fields["Method"] = c.Request.Method
			fields["StatusCode"] = c.Writer.Status()
			fields["ErrorMessage"] = c.Errors.ByType(sysadmServer.ErrorTypePrivate).String()
			fields["BodySize"] = c.Writer.Size()
			if raw != "" {
				path = path + "?" + raw
			}
			fields["Path"] = path

			l.WriteLogWithFields(fields, "info")
		}
	}
}

// Log log message to DefaultLogger, AccessLogger and ErrorLogger.
// if SplitAccessAndError is true, then log info log message to AccessLogger and other log message to ErrorLogger.
// otherwise  all log message will be log to AccessLogger.  If  SplitAccessAndError is true but ErrorLogger is nil
// then all log message will be log to AccessLogger and a additional warning log message will be log to AccessLogger.
// If SplitAccessAndError is false but AccessLogger is nil then all log message will be log to ErrorLogger and a additional
// warning log message will be log to ErrorLogger.
func (l *LoggerConfig) Log(message string, logLevel string) {
	if strings.TrimSpace(logLevel) == "" {
		logLevel = "error"
	}

	found := false
	for _, value := range Levels {
		if strings.ToLower(logLevel) == strings.ToLower(value) {
			found = true
		}
	}

	if !found {
		logLevel = "error"
	}

	logger := l.DefaultLogger
	if logger != nil {
		l.WriteLog(logger, message, logLevel)
	}

	accessLogger := l.AccessLogger
	errorLogger := l.ErrorLogger

	if strings.ToLower(logLevel) == "info" {
		if accessLogger != nil {
			l.WriteLog(accessLogger, message, logLevel)
		}
	} else {
		if l.SplitAccessAndError {
			if errorLogger != nil {
				l.WriteLog(errorLogger, message, logLevel)
			} else if accessLogger != nil {
				l.WriteLog(accessLogger, message, logLevel)
				logLevel = "warning"
				l.WriteLog(accessLogger, "SplitAccessAndError has be set to true, but  error log file has not be set.", logLevel)
			}
		} else {
			if accessLogger != nil {
				l.WriteLog(accessLogger, message, logLevel)
			} else if errorLogger != nil {
				l.WriteLog(errorLogger, message, logLevel)
				logLevel = "warning"
				l.WriteLog(errorLogger, "SplitAccessAndError has be set to false, but  access log file has not be set.", logLevel)
			}
		}
	}

}

// Logf log message to DefaultLogger, AccessLogger and ErrorLogger.
// if SplitAccessAndError is true, then log info log message to AccessLogger and other log message to ErrorLogger.
// otherwise  all log message will be log to AccessLogger.  If  SplitAccessAndError is true but ErrorLogger is nil
// then all log message will be log to AccessLogger and a additional warning log message will be log to AccessLogger.
// If SplitAccessAndError is false but AccessLogger is nil then all log message will be log to ErrorLogger and a additional
// warning log message will be log to ErrorLogger.
func (l *LoggerConfig) Logf(logLevel string, format string, a ...interface{}) {
	if strings.TrimSpace(logLevel) == "" {
		logLevel = "error"
	}

	found := false
	for _, value := range Levels {
		if strings.ToLower(logLevel) == strings.ToLower(value) {
			found = true
		}
	}

	if !found {
		logLevel = "error"
	}

	logMsg := fmt.Sprintf(format, a...)
	logger := l.DefaultLogger
	if logger != nil {
		l.WriteLog(logger, logMsg, logLevel)
	}

	accessLogger := l.AccessLogger
	errorLogger := l.ErrorLogger

	if strings.ToLower(logLevel) == "info" {
		if accessLogger != nil {
			l.WriteLog(accessLogger, logMsg, logLevel)
		}
	} else {
		if l.SplitAccessAndError {
			if errorLogger != nil {
				l.WriteLog(errorLogger, logMsg, logLevel)
			} else if accessLogger != nil {
				l.WriteLog(accessLogger, logMsg, logLevel)
				logLevel = "warning"
				l.WriteLog(accessLogger, "SplitAccessAndError has be set to true, but  error log file has not be set.", logLevel)
			}
		} else {
			if accessLogger != nil {
				l.WriteLog(accessLogger, logMsg, logLevel)
			} else if errorLogger != nil {
				l.WriteLog(errorLogger, logMsg, logLevel)
				logLevel = "warning"
				l.WriteLog(errorLogger, "SplitAccessAndError has be set to false, but  access log file has not be set.", logLevel)
			}
		}
	}

}

// Get the string of error level
// Return the string of the error level if the level was found, otherwise return ""
func GetLevelString(level int) string {
	if level < 0 || level > 6 {
		return ""
	}

	return Levels[level]
}

// Return the error level of err if err is not nil, otherwise return 1("debug")
func GetLevelStringInError(err Sysadmerror) string {
	if err == (Sysadmerror{}) {
		return ""
	}

	return GetLevelString(err.ErrorLevel)
}

// Get the index of error level.
// Return int of index if found otherwise return 1 which is the index of "debug"
func GetLevelNum(level string) int {

	for key := range Levels {
		if strings.ToLower(Levels[key]) == strings.ToLower(level) {
			return key
		}
	}
	return 1
}

// Return the error level of err if err is not nil, otherwise return 1("debug")
func GetLevelNumInError(err Sysadmerror) int {
	if err == (Sysadmerror{}) {
		return 0
	}

	return err.ErrorLevel
}

// Return the error no of err if err is not nil, otherwise return 0
func GetErrorNo(err Sysadmerror) int {
	if err == (Sysadmerror{}) {
		return 0
	}
	return err.ErrorNo
}

// log log messages to logfile or stdout
func (l *LoggerConfig) LogErrors(errs []Sysadmerror) {

	for _, e := range errs {
		lvl := GetLevelStringInError(e)
		no := e.ErrorNo
		l.Logf(lvl, "erroCode: %d Msg: %s", no, e.ErrorMsg)
	}
}

// NewErrorWithNumLevel create a new instance of Sysadmerror with errno,errLevel(int) and errMsg
// return Sysadmerror
func NewErrorWithNumLevel(errno int, errLevel int, errMsg string, args ...interface{}) Sysadmerror {
	errmsg := fmt.Sprintf(errMsg, args...)
	if errLevel < 0 || errLevel > 6 {
		errLevel = 1
	}

	err := Sysadmerror{
		ErrorNo:    errno,
		ErrorLevel: errLevel,
		ErrorMsg:   errmsg,
	}

	return err
}

// NewErrorWithStringLevel create a new instance of Sysadmerror with errno,errLevel(string) and errMsg
// return Sysadmerror. ErrorLevel will be set to 1("debug") if errLevel was not found in Levels
func NewErrorWithStringLevel(errno int, errLevel string, errMsg string, args ...interface{}) Sysadmerror {
	errmsg := fmt.Sprintf(errMsg, args...)
	level := GetLevelNum(errLevel)

	err := Sysadmerror{
		ErrorNo:    errno,
		ErrorLevel: level,
		ErrorMsg:   errmsg,
	}

	return err
}

// Get the maxLevels in []Sysadmerror
// return -1 if the length of []Sysadmerror less 1
// otherwise return the maxLevels in the []Sysadmerror
func GetMaxLevel(errs []Sysadmerror) int {
	if len(errs) < 1 {
		return -1
	}

	maxLevel := 0
	for _, v := range errs {
		l := v.ErrorLevel
		if l > maxLevel {
			maxLevel = l
		}
	}

	return maxLevel
}
