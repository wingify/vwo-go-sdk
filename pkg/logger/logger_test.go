/*
 * Copyright 2020-2022 Wingify Software Pvt. Ltd.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package logger

import (
	"bufio"
	"bytes"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoggingBeforeInit(t *testing.T) {
	logLevel = 3
	old := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	os.Stderr = w
	// Reset
	initialize()

	info := "info log"
	warning := "warning log"
	errL := "error log"
	fatal := "fatal log"

	Info(info)
	Warning(warning)
	Error(errL)
	// We don't want os.Exit in a test
	defaultLogger.output(sFatal, 0, fatal)

	w.Close()
	os.Stderr = old

	var b bytes.Buffer
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		b.Write(scanner.Bytes())
	}

	out := b.String()

	for _, txt := range []string{info, warning, errL, fatal} {
		if !strings.Contains(out, txt) {
			t.Errorf("log output %q does not contain expected text: %q", out, txt)
		}
	}
}

func TestInit(t *testing.T) {
	var buf1 bytes.Buffer
	l1 := Init("test1", false, false, &buf1)
	if !reflect.DeepEqual(l1, defaultLogger) {
		t.Fatal("defaultLogger does not match logger returned by Init")
	}

	// Subsequent runs of Init shouldn't change defaultLogger.
	var buf2 bytes.Buffer
	l2 := Init("test2", false, false, &buf2)
	if !reflect.DeepEqual(l1, defaultLogger) {
		t.Error("defaultLogger should not have changed")
	}

	// Check log output.
	l1.Info("logger #1")
	l2.Info("logger #2")
	defaultLogger.Info("default logger")

	tests := []struct {
		out  string
		want int
	}{
		{buf1.String(), 2},
		{buf2.String(), 1},
	}

	for i, tt := range tests {
		got := len(strings.Split(strings.TrimSpace(tt.out), "\n"))
		if got != tt.want {
			t.Errorf("logger %d wrong number of lines, want %d, got %d", i+1, tt.want, got)
		}
	}
}

func TestSetLogLevel(t *testing.T) {
	SetLogLevel(2)
	assert.Equal(t, 2, logLevel)
}
