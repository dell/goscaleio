/*
 *
 * Copyright © 2021-2024 Dell Inc. or its subsidiaries. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Copyright © 2025 Dell Inc. or its subsidiaries. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//      http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package log

import (
	"fmt"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

// checkLogs is a string used in our logging function to store the msg being logged
// used for verifying doLog is logging when expected
var checkLogs string

// logFunction is our logging function passed to DoLog()
// it will print the msg, and set checkLogs equal to the msg. We can then check the value of checkLogs
// to verify that DoLogs worked correctly
func logFunction(msg string, args ...any) {
	fmt.Println(msg, args)
	checkLogs = msg
}

func TestDoLog(t *testing.T) {
	type args struct {
		l   func(msg string, args ...any)
		msg string
	}
	tests := []struct {
		name     string
		args     args
		debug    bool
		expected string
	}{
		{
			name: "Test DoLog with debug true",
			args: args{
				l:   logFunction,
				msg: "test log",
			},
			debug:    true,
			expected: "test log",
		},
		{
			name: "Test DoLog with debug false",
			args: args{
				l:   logFunction,
				msg: "test log",
			},
			debug:    false,
			expected: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.debug {
				SetLogLevel(slog.LevelDebug)
			} else {
				SetLogLevel(slog.LevelInfo)
			}
			DoLog(tt.args.l, tt.args.msg)
			assert.Equal(t, tt.expected, checkLogs)
			// reset to empty string
			checkLogs = ""
		})
	}
}

func TestSetLogLevel(t *testing.T) {
	tests := []struct {
		name  string
		level slog.Level
	}{
		{
			name:  "Test SetLogLevel to debug",
			level: slog.LevelDebug,
		},
		{
			name:  "Test SetLogLevel to info",
			level: slog.LevelInfo,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetLogLevel(tt.level)
			assert.Equal(t, tt.level, logLevel.Level())
			if tt.level == slog.LevelDebug {
				assert.True(t, debug)
			} else {
				assert.False(t, debug)
			}
		})
	}
}
