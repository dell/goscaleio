// Copyright © 2019 - 2026 Dell Inc. or its subsidiaries. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package api_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	api "github.com/dell/goscaleio/api"
	"github.com/stretchr/testify/require"
)

type captureObserver struct {
	mu           sync.Mutex
	observations []api.RequestObservation
}

func (o *captureObserver) ObserveRequest(obs api.RequestObservation) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.observations = append(o.observations, obs)
}

func (o *captureObserver) snapshot() []api.RequestObservation {
	o.mu.Lock()
	defer o.mu.Unlock()
	out := make([]api.RequestObservation, len(o.observations))
	copy(out, o.observations)
	return out
}

func TestClient_RequestObserver_RecordsSuccessfulRequest(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("unexpected method: got %s want %s", r.Method, http.MethodGet)
		}
		if r.URL.Path != "/api/version" {
			t.Errorf("unexpected path: got %s want %s", r.URL.Path, "/api/version")
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("5.0.0"))
	}))
	defer srv.Close()

	client, err := api.New(context.Background(), srv.URL, api.ClientOptions{Timeout: time.Second})
	require.NoError(t, err)

	observer := &captureObserver{}
	client.SetRequestObserver(observer)

	_, err = client.DoAndGetResponseBody(context.Background(), http.MethodGet, "/api/version", nil, nil, "")
	require.NoError(t, err)

	observations := observer.snapshot()
	require.Len(t, observations, 1)
	require.Equal(t, "/api/version", observations[0].Endpoint)
	require.Equal(t, http.MethodGet, observations[0].Method)
	require.Equal(t, http.StatusOK, observations[0].StatusCode)
	require.NoError(t, observations[0].Err)
	require.Greater(t, observations[0].Duration, time.Duration(0))
}

func TestClient_RequestObserver_RecordsTransportError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	client, err := api.New(context.Background(), srv.URL, api.ClientOptions{Timeout: time.Second})
	require.NoError(t, err)

	observer := &captureObserver{}
	client.SetRequestObserver(observer)

	srv.Close()

	_, err = client.DoAndGetResponseBody(context.Background(), http.MethodGet, "/api/version", nil, nil, "")
	require.Error(t, err)

	observations := observer.snapshot()
	require.Len(t, observations, 1)
	require.Equal(t, "/api/version", observations[0].Endpoint)
	require.Equal(t, http.MethodGet, observations[0].Method)
	require.Equal(t, 0, observations[0].StatusCode)
	require.Error(t, observations[0].Err)
}

func TestObserver_PanicIsContained(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	client, err := api.New(context.Background(), srv.URL, api.ClientOptions{Timeout: time.Second})
	require.NoError(t, err)

	client.SetRequestObserver(RequestObserverFunc(func(api.RequestObservation) {
		panic(errors.New("boom"))
	}))

	_, err = client.DoAndGetResponseBody(context.Background(), http.MethodGet, "/api/version", nil, nil, "")
	require.NoError(t, err)
}

// RequestObserverFunc adapts a function to the RequestObserver interface.
type RequestObserverFunc func(api.RequestObservation)

// ObserveRequest implements api.RequestObserver.
func (f RequestObserverFunc) ObserveRequest(obs api.RequestObservation) {
	f(obs)
}
