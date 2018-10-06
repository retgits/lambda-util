// Package util implements utility methods
package util

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPGet(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusForbidden)
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", r.Header["Content-Type"][0])
		io.WriteString(w, `{"hello":"world"}`)

	}))
	defer ts.Close()

	resp, err := HTTPGet(ts.URL, "application/json", nil)
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, 200)
}

func TestHTTPPost(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusForbidden)
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", r.Header["Content-Type"][0])
		io.WriteString(w, `{"hello":"world"}`)

	}))
	defer ts.Close()

	resp, err := HTTPPost(ts.URL, "application/json", "", nil)
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, 200)
}

func TestHTTPPatch(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PATCH" {
			w.WriteHeader(http.StatusForbidden)
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", r.Header["Content-Type"][0])
		io.WriteString(w, `{"hello":"world"}`)

	}))
	defer ts.Close()

	resp, err := HTTPPatch(ts.URL, "application/json", "", nil)
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, 200)
}

func TestHTTPGetWithHeader(t *testing.T) {
	basicAuthUserPass := "username:password"
	basicAuthString := fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(basicAuthUserPass)))

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Authorization"][0] != basicAuthString {
			w.WriteHeader(http.StatusForbidden)
		} else {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", r.Header["Content-Type"][0])
			io.WriteString(w, `{"hello":"world"}`)
		}
	}))
	defer ts.Close()

	httpHeader := http.Header{"Authorization": {basicAuthString}}
	resp, err := HTTPGet(ts.URL, "application/json", httpHeader)
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, 200)

	httpHeader = http.Header{"Authorization": {"bla!"}}
	resp, err = HTTPGet(ts.URL, "application/json", httpHeader)
	assert.Error(t, err)
	assert.Equal(t, resp.StatusCode, 403)
}

func TestHTTPPostWithData(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		if r.Method != "POST" || len(string(body)) < 5 {
			w.WriteHeader(http.StatusForbidden)
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", r.Header["Content-Type"][0])
		io.WriteString(w, `{"hello":"world"}`)

	}))
	defer ts.Close()

	resp, err := HTTPPost(ts.URL, "application/json", `{"hello":"world"}`, nil)
	assert.NoError(t, err)
	assert.Equal(t, resp.StatusCode, 200)
}
