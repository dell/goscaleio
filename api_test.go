package goscaleio

import (
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

func setupClient(t *testing.T, hostAddr string) *Client {
	os.Setenv("GOSCALEIO_ENDPOINT", hostAddr+"/api")
	client, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}
	// test ok
	_, err = client.Authenticate(&ConfigConnect{
		Username: "ScaleIOUser",
		Password: "password",
		Version:  "2.0",
	})
	if err != nil {
		t.Fatal(err)
	}
	return client
}

func requestAuthOK(resp http.ResponseWriter, req *http.Request) bool {
	_, pwd, _ := req.BasicAuth()
	if pwd == "" {
		resp.WriteHeader(http.StatusUnauthorized)
		resp.Write([]byte(`{"message":"Unauthorized","httpStatusCode":401,"errorCode":0}`))
		return false
	}
	return true
}

func handleAuthToken(resp http.ResponseWriter, req *http.Request) {
	if !requestAuthOK(resp, req) {
		return
	}
	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte(`"012345678901234567890123456789"`))
}

func TestClientVersion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(resp http.ResponseWriter, req *http.Request) {
			if req.RequestURI != "/api/version" {
				t.Fatal("Expecting endpoint /api/version got", req.RequestURI)
			}
			resp.WriteHeader(http.StatusOK)
			resp.Write([]byte(`"2.0"`))
		},
	))
	defer server.Close()
	hostAddr := server.URL
	os.Setenv("GOSCALEIO_ENDPOINT", hostAddr+"/api")
	client, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}
	ver, err := client.getVersion()
	if err != nil {
		t.Fatal(err)
	}
	if ver != "2.0" {
		t.Fatal("Expecting version string \"2.0\", got ", ver)
	}
}

func TestClientLogin(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(resp http.ResponseWriter, req *http.Request) {
			switch req.RequestURI {
			case "/api/version":
				resp.WriteHeader(http.StatusOK)
				resp.Write([]byte(`"2.0"`))
			case "/api/login":
				//accept := req.Header.Get("Accept")
				// check Accept header
				//if ver := strings.Split(accept, ";"); len(ver) != 2 {
				//	t.Fatal("Expecting Accept header to include version")
				//} else {
				//	if !strings.HasPrefix(ver[1], "version=") {
				//		t.Fatal("Header Accept must include version")
				//	}
				//}

				uname, pwd, basic := req.BasicAuth()
				if !basic {
					t.Fatal("Client only support basic auth")
				}

				if uname != "ScaleIOUser" || pwd != "password" {
					resp.WriteHeader(http.StatusUnauthorized)
					resp.Write([]byte(`{"message":"Unauthorized","httpStatusCode":401,"errorCode":0}`))
					return
				}
				resp.WriteHeader(http.StatusOK)
				resp.Write([]byte(`"012345678901234567890123456789"`))
			default:
				t.Fatal("Expecting endpoint /api/login got", req.RequestURI)
			}

		},
	))
	defer server.Close()
	hostAddr := server.URL
	os.Setenv("GOSCALEIO_ENDPOINT", hostAddr+"/api")
	client, err := NewClient()
	if err != nil {
		t.Fatal(err)
	}
	// test ok
	_, err = client.Authenticate(&ConfigConnect{
		Username: "ScaleIOUser",
		Password: "password",
		Endpoint: "",
		Version:  "2.0",
	})
	if err != nil {
		t.Fatal(err)
	}
	if client.GetToken() != "012345678901234567890123456789" {
		t.Fatal("Expecting token 012345678901234567890123456789, got", client.GetToken())
	}

	// test bad login
	_, err = client.Authenticate(&ConfigConnect{
		Username: "ScaleIOUser",
		Password: "badPassWord",
		Endpoint: "",
		Version:  "2.0",
	})
	if err == nil {
		t.Fatal("Expecting an error for bad Login, but did not")
	}
}

type stubTypeWithMetaData struct{}

func (s stubTypeWithMetaData) MetaData() http.Header {
	h := make(http.Header)
	h.Set("foo", "bar")
	return h
}

func Test_addMetaData(t *testing.T) {
	var tests = []struct {
		name           string
		givenHeader    map[string]string
		expectedHeader map[string]string
		body           interface{}
	}{
		{"nil header is a noop", nil, nil, nil},
		{"nil body is a noop", nil, nil, nil},
		{"header is updated", make(map[string]string), map[string]string{"Foo": "bar"}, stubTypeWithMetaData{}},
		{"header is not updated", make(map[string]string), map[string]string{}, struct{}{}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			addMetaData(tt.givenHeader, tt.body)

			switch {
			case tt.givenHeader == nil:
				if tt.givenHeader != nil {
					t.Errorf("(%s): expected %s, actual %s", tt.body, tt.expectedHeader, tt.givenHeader)
				}
			case tt.body == nil:
				if len(tt.givenHeader) != 0 {
					t.Errorf("(%s): expected %s, actual %s", tt.body, tt.expectedHeader, tt.givenHeader)
				}
			default:
				if !reflect.DeepEqual(tt.expectedHeader, tt.givenHeader) {
					t.Errorf("(%s): expected %s, actual %s", tt.body, tt.expectedHeader, tt.givenHeader)
				}
			}
		})
	}
}
