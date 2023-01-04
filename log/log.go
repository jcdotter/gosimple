// Copyright 2022 escend llc. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the gosimple LICENSE file.
// Author: jcdotter

// This package is an alternative to go 'log' pkg
// Includes simplified implementation with default:
//   log levels - TRACE, INFO, WARNING, ERROR, FATAL
//   log format - INFO: 2006-01-02 15:04:05.000 main.go:12 log message here
//   log file location - '../logs/file.log'
// with the ability to customize logging prior to first log record
// If a log directory is not provided in os.Setenv("GO_UTILS_LOG_PATH")
// or in SetLogDir, then a directory is created at '../logs'
// Fatal functions call os.Exit(1) after posting to log

package log

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// CONFIGS: settings for logging destination and format
var (
	__SESSION__    string                                                            // the unique id to the log session
	__HOST__       string                                                            // the source host server from os.GetEnv("HOST")
	__SERVICE__    string                                                            // the source servoce from os.GetEnv("SERVICE")
	__WRITER__     io.Writer                                                         // the writer used to post to log
	__DIR__        string                                                            // the directory path to log files
	__FILE__       string                                                            // the name of the current session log file
	__DELIM__      = " \t"                                                           // the delimeter between log line elements
	__JSON_FMT__   bool                                                              // if true, post log line in json format
	__TIME_FMT__   = `2006-01-02 15:04:05.000`                                       // the date format posted to log
	__TO_CONSOLE__ = true                                                            // if true, post logs to console
	__ACTIVE__     bool                                                              // if true, configs are locked
	__FORMAT__     = []int{LogLevel, LogDateTime, LogSession, LogSource, LogMessage} // the format for a log line
)

// Format Elements: elements included in log record
const (
	LogLevel      = iota // post log level to log
	LogDateTime          // post date to log
	LogFullSource        // post full source file path and line to log
	LogSource            // post only source file name and line to log
	LogSession           // post the log session to log
	LogHost              // post the source host server from os.GetEnv("HOST") to the log
	LogService           // post the source service from os.GetEnv("SERVICE") to the log
	LogMessage           // post the log message to the log
	LogJsonFmt           // post log line in json format
	LogStdFmt            // post log line in delimited format
)

var elNames = []string{
	LogLevel:      "level",
	LogDateTime:   "datetime",
	LogFullSource: "fullsource",
	LogSource:     "source",
	LogSession:    "session",
	LogHost:       "host",
	LogService:    "service",
	LogMessage:    "message",
}

// SetFormat configures the order
// and elements of a log record
// using the elements and their order provided
// as arguments to the function.
// Elements provied as log.Log<element>
func SetFormat(f ...int) {
	ft := []int{}
	l := len(elNames)
	for _, i := range f {
		if l > i {
			ft = append(ft, i)
		} else if i == LogJsonFmt {
			__JSON_FMT__ = true
		} else if i == LogStdFmt {
			__JSON_FMT__ = false
		}
	}
	if len(ft) > 0 {
		__FORMAT__ = ft
	}
}

// SetDateTimeFormat sets the format of the
// datetime stamp in the log record and
// uses the same formats as the go time pkg
func SetDateTimeFormat(f string) {
	if !__ACTIVE__ {
		_, err := time.Parse(string(f), string(f))
		if err != nil {
			panic("could not set log datetime format: invalid format")
		}
		__TIME_FMT__ = f
	}
}

// SetDelim sets the delimiter used to
// separated log record elements
func SetDelim(d string) {
	if !__ACTIVE__ {
		__DELIM__ = d
	}
}

// LogToConsole controls whether logs are
// written to the console during runtime
func LogToConsole(c bool) {
	if !__ACTIVE__ {
		__TO_CONSOLE__ = c
	}
}

// SetLogDir overides the env var GO_UTILS_LOG_PATH and
// sets the location of the log files to the path provided
func SetDir(d string) {
	if !__ACTIVE__ {
		if _, err := os.Stat(d); errors.Is(err, os.ErrNotExist) {
			panic("could not set custom log dir: " + d)
		}
		__DIR__ = d
	}
}

// SetLogFile overides the standard file naming and
// sets the name of the log file in the log directory
func SetFile(f string) {
	if !__ACTIVE__ {
		__FILE__ = f
	}
}

// SetLogWriter overides the standard writer with
// a custom provided io.writer
func SetWriter(w io.Writer) {
	if !__ACTIVE__ {
		__WRITER__ = w
	}
}

// SetHost overides the env var HOST and uses
// the host provided in log posts
func SetHost(h string) {
	if !__ACTIVE__ {
		__HOST__ = h
	}
}

// SetService overides the env var SERVICE and uses
// the service provided in log posts
func SetService(s string) {
	if !__ACTIVE__ {
		__SERVICE__ = s
	}
}

// LOG LEVELS: Level manages the logging levels
type Level uint

const (
	TRACE Level = iota
	INFO
	WARNING
	ERROR
	FATAL
)

var levelNames = []string{
	TRACE:   "TRACE",
	INFO:    "INFO",
	WARNING: "WARNING",
	ERROR:   "ERROR",
	FATAL:   "FATAL",
}

// Strings returns a string version of the logging Level
func (l Level) String() string {
	return levelNames[uint(l)]
}

// LevelByName returns logging Level for the provided string
func LevelByName(s string) Level {
	s = strings.ToUpper(s)
	var l Level
	for _, v := range levelNames {
		if v == s {
			return l
		}
	}
	return l
}

// Log records an entry to the log file
// and prints to console if log.LogToConsole(true)
// using the Level 'l' and 'msg' message provided
func Log(l Level, msg string) {
	dt := time.Now()
	if !__ACTIVE__ {
		activate()
	}
	_, fl, ln, _ := runtime.Caller(2)
	fs := fmt.Sprint(fl, ":", ln)
	logEls := map[string]string{
		"level":      levelNames[l],
		"datetime":   dt.Format(__TIME_FMT__),
		"session":    __SESSION__,
		"host":       __HOST__,
		"service":    __SERVICE__,
		"fullsource": fs,
		"source":     fs[strings.LastIndex(fs, "/")+1:],
		"message":    msg,
	}
	var r []byte
	if __JSON_FMT__ {
		r = buildJsonLog(logEls)
	} else {
		r = buildStdLog(logEls)
	}
	r = append(r, "\n"...)
	__WRITER__.Write(r)
}

// buildStdLog is a helper function to Log
// builds standard log format using elements in __FORMAT__
// separated by the __DELIM__
func buildStdLog(els map[string]string) []byte {
	var log string
	for i, el := range __FORMAT__ {
		if v := els[elNames[el]]; v != "" {
			if i > 0 {
				log += __DELIM__
			}
			log += v
		}
	}
	return []byte(log)
}

// buildJsonLog is a helper function to Log
// builds a json log format using the elements in __FORMAT__
func buildJsonLog(els map[string]string) []byte {
	log := map[string]string{}
	for _, el := range __FORMAT__ {
		if v := els[elNames[el]]; v != "" {
			log[elNames[el]] = v
		}
	}
	r, _ := json.Marshal(log)
	return r
}

// Trace is typically used for debugging
// it records a TRACE emtry to the log file
// and prints to console if log.LogToConsole(true)
// using 'msg' message provided and the stacktrace
func Trace(msg string) {
	Log(TRACE, msg)
}

// Info records an INFO entry to the log file
// and prints to console if log.LogToConsole(true)
// using 'msg' message provided
func Info(msg string) {
	Log(INFO, msg)
}

// Warning records a WARNING entry to the log file
// and prints to console if log.LogToConsole(true)
// using 'msg' message provided
func Warning(msg string) {
	Log(WARNING, msg)
}

// Error  records an ERROR entry to the log file
// and prints to console if log.LogToConsole(true)
// using 'msg' message provided
func Error(msg string) {
	Log(ERROR, msg)
}

// Fatal records a FATAL entry to the log file
// prints to console if log.LogToConsole(true)
// using 'msg' message provided
// and exits application using os.Exit(1)
func Fatal(msg string) {
	Log(FATAL, msg)
	os.Exit(1)
}

// Logf records an entry to the log file
// and prints to console if log.LogToConsole(true)
// using the Level 'l' and 'msg' message provided.
// Arguments are handled in the manner of fmt.Printf
func Logf(l Level, format string, a ...any) {
	Log(l, fmt.Sprintf(format, a...))
}

// Tracef is typically used for debugging
// it records a TRACE emtry to the log file
// and prints to console if log.LogToConsole(true)
// using 'msg' message provided and the stacktrace.
// Arguments are handled in the manner of fmt.Printf
func Tracef(format string, a ...any) {
	Logf(TRACE, format, a...)
}

// Infof records an INFO entry to the log file
// and prints to console if log.LogToConsole(true)
// using 'msg' message provided.
// Arguments are handled in the manner of fmt.Printf
func Infof(format string, a ...any) {
	Logf(INFO, format, a...)
}

// Warningf records a WARNING entry to the log file
// and prints to console if log.LogToConsole(true)
// using 'msg' message provided.
// Arguments are handled in the manner of fmt.Printf
func Warningf(format string, a ...any) {
	Logf(WARNING, format, a...)
}

// Error  records an ERROR entry to the log file
// and prints to console if log.LogToConsole(true)
// using 'msg' message provided.
// Arguments are handled in the manner of fmt.Printf
func Errorf(format string, a ...any) {
	Logf(ERROR, format, a...)
}

// Fatalf records a FATAL entry to the log file
// prints to console if log.LogToConsole(true)
// using 'msg' message provided
// and exits application using os.Exit(1).
// Arguments are handled in the manner of fmt.Printf
func Fatalf(format string, a ...any) {
	Logf(FATAL, format, a...)
	os.Exit(1)
}

// Read parses the active log file to a map
// and returns it for log evaluation
func Read() []map[string]any {
	m := []map[string]any{}
	cnt, err := os.ReadFile(__DIR__ + "/" + __FILE__)
	if err != nil {
		panic("could not read log file")
	}
	for _, ln := range strings.Split(string(cnt), "\n") {
		if len(ln) > 0 {
			l := map[string]any{}
			if __JSON_FMT__ {
				json.Unmarshal([]byte(ln), &l)
				m = append(m, l)
			} else {
				lVals := strings.Split(ln, __DELIM__)
				for i, el := range __FORMAT__ {
					if el == LogDateTime {
						v, _ := time.Parse(__TIME_FMT__, lVals[i])
						l[elNames[el]] = v
					} else if el == LogSource || el == LogFullSource {
						s := strings.Split(lVals[i], ":")
						l["file"] = s[0]
						l["line"], _ = strconv.Atoi(s[1])
					} else {
						l[elNames[el]] = lVals[i]
					}
				}
				m = append(m, l)
			}
		}
	}
	return m
}

// Activates logging session by
// by setting the session id, host, and service
// configuring the log writer and
// setting the __ACTIVE__ config to true
func activate() {
	if __SESSION__ == "" {
		initSession()
	}
	if __HOST__ == "" {
		h, exists := os.LookupEnv("HOST")
		if exists {
			__HOST__ = h
		}
	}
	if __SERVICE__ == "" {
		s, exists := os.LookupEnv("SERVICE")
		if exists {
			__SERVICE__ = s
		}
	}
	if __WRITER__ == nil {
		initWriter()
	}
	__ACTIVE__ = true
}

// generate and set the session id
func initSession() {
	__SESSION__ = time.Now().Format(`060102-150405`) + "-" + rGen(6)
}

func rGen(l int) string {
	b := make([]byte, l)
	rand.Seed(time.Now().UnixNano())
	for i := range b {
		n := rand.Intn(62)
		if n > 35 {
			n += 13
		} else if n > 9 {
			n += 7
		}
		b[i] = byte(48 + n)
	}
	return string(b)
}

// initWriter returns the writer for the log
// printing to both a log file and the console
// panics if it cannot access or create a log file
func initWriter() {
	var exists bool
	var err error
	if __DIR__ == "" { // set dir if not alread set
		__DIR__, exists = os.LookupEnv("GO_UTILS_LOG_PATH")
		if !exists { // if env var not set, use default dir
			__DIR__, err = filepath.Abs("../logs")
			if err != nil {
				panic("cannot initialize logger path")
			}
			os.Setenv("GO_UTILS_LOG_PATH", __DIR__)
		}
	}
	if _, err := os.Stat(__DIR__); errors.Is(err, os.ErrNotExist) {
		//create dir if it dones not already exist
		err := os.Mkdir(__DIR__, os.ModePerm)
		if err != nil {
			panic("could not initialize log directory: " + __DIR__)
		}
	}
	if __FILE__ == "" { // set default file name if not already set
		__FILE__ = __SESSION__ + ".log"
	}
	s := "/"
	if string(__FILE__[0]) == "/" {
		s = ""
	}
	if __WRITER__ == nil { // generate io writer if not already set
		file, err := os.OpenFile(__DIR__+s+__FILE__, os.O_RDWR|os.O_CREATE, os.ModePerm)
		if err != nil {
			panic("could not initatiate log file")
		}
		if __TO_CONSOLE__ {
			__WRITER__ = io.MultiWriter(os.Stdout, file)
		} else {
			__WRITER__ = io.Writer(file)
		}
	}
}
