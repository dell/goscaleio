// Copyright © 2024 Dell Inc. or its subsidiaries. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package api

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestIsBinOctetBody(t *testing.T) {
	tests := []struct {
		name     string
		header   http.Header
		expected bool
	}{
		{
			name: "Content-Type is binary/octet-stream",
			header: http.Header{
				"Content-Type": []string{"binary/octet-stream"},
			},
			expected: true,
		},
		{
			name: "Content-Type is application/json",
			header: http.Header{
				"Content-Type": []string{"application/json"},
			},
			expected: false,
		},
		{
			name:     "Content-Type is missing",
			header:   http.Header{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the function with the test case's header
			got := isBinOctetBody(tt.header)

			if got != tt.expected {
				t.Errorf("isBinOctetBody() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestDumpRequest(t *testing.T) {
	// Test case: Request with no body
	req, err := http.NewRequest("GET", "http://example.com/instances", nil)
	if err != nil {
		t.Fatal(err)
	}
	b, err := dumpRequest(req, !isBinOctetBody(req.Header))
	if err != nil {
		t.Fatal(err)
	}
	expected := "GET /instances HTTP/1.1\r\nHost: example.com\r\n\r\n"
	if string(b) != expected {
		t.Errorf("Expected %q, got %q", expected, string(b))
	}

	// Test case: Request with body
	req, err = http.NewRequest("POST", "http://example.com/instances", strings.NewReader("Hello, world!"))
	if err != nil {
		t.Fatal(err)
	}
	b, err = dumpRequest(req, true)
	if err != nil {
		t.Fatal(err)
	}
	expected = "POST /instances HTTP/1.1\r\nHost: example.com\r\n\r\nHello, world!"
	if string(b) != expected {
		t.Errorf("Expected %q, got %q", expected, string(b))
	}

	// Test case: Request with chunked transfer encoding
	req, err = http.NewRequest("PUT", "http://example.com/instances", strings.NewReader("Hello, world!"))
	if err != nil {
		t.Fatal(err)
	}
	req.TransferEncoding = []string{"chunked"}
	b, err = dumpRequest(req, true)
	if err != nil {
		t.Fatal(err)
	}
	expected = "PUT /instances HTTP/1.1\r\nHost: example.com\r\nTransfer-Encoding: chunked\r\n\r\nd\r\nHello, world!\r\n0\r\n\r\n"
	if string(b) != expected {
		t.Errorf("Expected %q, got %q", expected, string(b))
	}

	// Test case: Request with connection close
	req, err = http.NewRequest("DELETE", "http://example.com/instances", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Close = true
	b, err = dumpRequest(req, false)
	if err != nil {
		t.Fatal(err)
	}
	expected = "DELETE /instances HTTP/1.1\r\nHost: example.com\r\nConnection: close\r\n\r\n"
	if string(b) != expected {
		t.Errorf("Expected %q, got %q", expected, string(b))
	}
}

func TestWriteIndentedN(t *testing.T) {
	// Test case: Empty input
	var b bytes.Buffer
	err := WriteIndentedN(&b, []byte{}, 4)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if b.String() != "" {
		t.Errorf("Expected empty output, got: %s", b.String())
	}

	// Test case: Single line
	b.Reset()
	err = WriteIndentedN(&b, []byte("Hello, world!"), 4)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expected := "    Hello, world!"
	if b.String() != expected {
		t.Errorf("Expected: %s, got: %s", expected, b.String())
	}

	// Test case: Multiple lines
	b.Reset()
	err = WriteIndentedN(&b, []byte("Line 1\nLine 2\nLine 3"), 2)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expected = "  Line 1\n  Line 2\n  Line 3"
	if b.String() != expected {
		t.Errorf("Expected: %s, got: %s", expected, b.String())
	}
}

func TestLogResponse(_ *testing.T) {
	// Test case: Logging a successful response
	res := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(`{"message": "success"}`)),
	}
	logResponse(context.Background(), res, nil)

	// Test case: Logging a failed response
	res = &http.Response{
		StatusCode: http.StatusInternalServerError,
		Header:     http.Header{"Content-Type": []string{"text/plain"}},
		Body:       io.NopCloser(strings.NewReader("Internal server error")),
	}
	logResponse(context.Background(), res, nil)

	// Test case: Logging a response with binary content
	res = &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/octet-stream"}},
		Body:       io.NopCloser(bytes.NewReader([]byte{0x01, 0x02, 0x03})),
	}
	logResponse(context.Background(), res, nil)
}

func logChecker(level log.Level, msg string) func(func(args ...interface{}), string) {
	return func(lf func(args ...interface{}), message string) {
		if level != log.DebugLevel {
			// If the log level is not DebugLevel, call the logging function with an error message
			lf(fmt.Sprintf("Expected debug level, got %v", level))
		}
		if !strings.Contains(msg, "GOSCALEIO HTTP REQUEST") {
			// If the message does not contain the expected string, call the logging function with an error message
			lf(fmt.Sprintf("Expected request log, got %s", msg))
		}

		// You can log the original message if needed
		lf(message)
	}
}

func TestLogRequest(t *testing.T) {
	// Test case: Valid request
	req, err := http.NewRequest("GET", "https://example.com/instances", nil)
	if err != nil {
		t.Fatal(err)
	}
	logFunc := logChecker(log.DebugLevel, "GOSCALEIO HTTP REQUEST")
	logRequest(context.Background(), req, logFunc)

	// Test case: Error in dumpRequest
	req, err = http.NewRequest("GET", "https://example.com/instances", NewErrorReader(errors.New("simulated error while reading request body")))
	if err != nil {
		t.Fatal(err)
	}
	_, err = io.ReadAll(req.Body)
	if err == nil {
		t.Fatalf("Expected error when reading request body, got nil")
	}

	logRequest(context.Background(), req, logFunc)

	//Test case: Error in WriteIndented
	req, err = http.NewRequest("GET", "https://example.com/instances", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set(HeaderKeyContentType, "application/json")

	logRequest(context.Background(), req, logFunc)
}

// errorReader is a custom io.Reader that always returns an error when Read is called.
type ErrorReader struct {
	err error
}

// NewErrorReader creates an instance of errorReader with the specified error.
func NewErrorReader(err error) *ErrorReader {
	return &ErrorReader{err: err}
}

// Read always returns an error (simulating a read failure).
func (r *ErrorReader) Read(_ []byte) (n int, err error) {
	return 0, r.err
}

func TestDrainBody(t *testing.T) {
	// Test case: b is http.NoBody
	r1, r2, err := drainBody(http.NoBody)
	if r1 != http.NoBody {
		t.Errorf("Expected r1 to be http.NoBody, got %v", r1)
	}
	if r2 != http.NoBody {
		t.Errorf("Expected r2 to be http.NoBody, got %v", r2)
	}
	if err != nil {
		t.Errorf("Expected err to be nil, got %v", err)
	}

	// Test case: b is not http.NoBody
	b := io.NopCloser(strings.NewReader("test"))
	r1, r2, err = drainBody(b)
	if r1 == http.NoBody {
		t.Error("Expected r1 to not be http.NoBody")
	}
	if r2 == http.NoBody {
		t.Error("Expected r2 to not be http.NoBody")
	}
	if err != nil {
		t.Errorf("Expected err to be nil, got %v", err)
	}
	if _, err := r1.Read(make([]byte, 1)); err != nil {
		t.Errorf("Expected r1 to be readable, got %v", err)
	}
	if _, err := r2.Read(make([]byte, 1)); err != nil {
		t.Errorf("Expected r2 to be readable, got %v", err)
	}

	// Test case: b.ReadFrom returns an error
	b = io.NopCloser(errReader{})
	r1, r2, err = drainBody(b)
	if r1 != nil {
		t.Errorf("Expected r1 to be nil, got %v", r1)
	}
	if r2 == nil {
		t.Errorf("Expected r2 to be not nil, got %v", r2)
	}
	if err == nil {
		t.Error("Expected err to not be nil")
	}

	// Test case: b.Close returns an error
	b = io.NopCloser(strings.NewReader("test"))
	r1, r2, err = drainBody(b)
	if r1 == nil {
		t.Error("Expected r1 to not be nil")
	}
	if r2 == nil {
		t.Error("Expected r2 to not be nil")
	}
	if err != nil {
		t.Error("Expected err to be nil")
	}
}

type errReader struct{}

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("read error")
}
