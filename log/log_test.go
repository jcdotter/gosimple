// Copyright 2022 escend llc. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the gosimple LICENSE file.
// Author: jcdotter

package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

var (
	format    = []int{LogHost, LogService, LogSession, LogSource, LogDateTime, LogLevel, LogMessage}
	fdatetime = `20060102 150405`
	delim     = "|"
	dir, dErr = filepath.Abs("../../logs")
	file      = "/test.log"
	buffer    = new(bytes.Buffer)
	host      = "test-server"
	service   = "test-service"
)

func deactivate() {
	__ACTIVE__ = false
}

func TestFormat(t *testing.T) {
	f := append([]int{LogJsonFmt}, format...)
	SetFormat(f...)
	if !reflect.DeepEqual(__FORMAT__, format) {
		t.Fatal("SetFormat did not update Log __FORMAT__")
	}
	if !__JSON_FMT__ {
		t.Fatal("SetFormat did not update Log __JSON_FMT__")
	}
	SetDateTimeFormat(fdatetime)
	if __TIME_FMT__ != fdatetime {
		t.Fatal("SetDateTimeFormat did not update Log __TIME_FMT__")
	}
	SetDelim(delim)
	if __DELIM__ != delim {
		t.Fatal("SetDelim did not update Log __DELIM__")
	}
	LogToConsole(false)
	if __TO_CONSOLE__ {
		t.Fatal("LogToConsole did not update Log __TO_CONSOLE__")
	}
	SetDir(dir)
	if dErr != nil || __DIR__ != dir {
		t.Fatal("SetDir did not update Log __DIR__")
	}
	SetFile(file)
	if __FILE__ != file {
		t.Fatal("SetFile did not update Log __FILE__")
	}
	SetWriter(buffer)
	if __WRITER__ == nil {
		t.Fatal("SetWriter did not update Log __WRITER__")
	}
	SetHost(host)
	if __HOST__ != host {
		t.Fatal("SetHost did not update Log __HOST__")
	}
	SetService(service)
	if __SERVICE__ != service {
		t.Fatal("SetService did not update Log __SERVICE__")
	}
}

func TestActivate(t *testing.T) {
	activate()
	if !__ACTIVE__ {
		t.Fatal("Log was not activated on activate()")
	}
	if __SESSION__ == "" {
		t.Fatal("Session was not set on activate()")
	}
	if __DIR__ != dir {
		t.Fatal("Activate overrode preconfigured __DIR__")
	}
	if __FILE__ != file {
		t.Fatal("Activate overrode preconfigured __FILE__")
	}
}

func TestLog(t *testing.T) {
	Logf(INFO, "test log Level: %s", levelNames[INFO])
	m := map[string]string{}
	err := json.Unmarshal(buffer.Bytes(), &m)
	if err != nil {
		t.Fatal("could not unmarshal json log")
	}
	fmt.Println(buffer)
	if v, ok := m["host"]; !ok || v != host {
		t.Fatal("log post host does not match __HOST__")
	}
	if v, ok := m["service"]; !ok || v != service {
		t.Fatal("log post service does not match __SERVICE__")
	}
	if v, ok := m["session"]; !ok || v != __SESSION__ {
		t.Fatal("log post session does not match __SESSION__")
	}
	if v, ok := m["level"]; !ok || v != levelNames[INFO] {
		t.Fatal("log post level does not match level provided")
	}
	if v, ok := m["source"]; !ok || v == "" {
		t.Fatal("log post source is missing")
	}
	if v, ok := m["datetime"]; !ok {
		t.Fatal("log post datetime is missing")
	} else if _, err := time.Parse(fdatetime, v); err != nil {
		t.Fatal("log post datetime does not match __TIME_FMT__")
	}
	if v, ok := m["message"]; !ok || v != "test log Level: INFO" {
		t.Fatal("log post message does not match the message provided")
	}
}
