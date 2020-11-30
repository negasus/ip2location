package http

import (
	"bytes"
	"errors"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

func Test_getIP_query_arg(t *testing.T) {
	req := &http.Request{
		URL: &url.URL{
			RawQuery: "ip=4.0.0.0",
		},
		Header: map[string][]string{},
		Method: http.MethodPost,
		Body:   ioutil.NopCloser(bytes.NewReader([]byte("6.0.0.0"))),
	}
	req.Header.Add("X-IP2LOCATION-IP", "5.0.0.0")

	ip := getIP(req)

	assert.Equal(t, "4.0.0.0", ip)
}

func Test_getIP_query_header(t *testing.T) {
	req := &http.Request{
		URL:    &url.URL{},
		Header: map[string][]string{},
		Method: http.MethodPost,
		Body:   ioutil.NopCloser(bytes.NewReader([]byte("6.0.0.0"))),
	}
	req.Header.Add("X-IP2LOCATION-IP", "5.0.0.0")

	ip := getIP(req)

	assert.Equal(t, "5.0.0.0", ip)
}

func Test_getIP_body(t *testing.T) {
	req := &http.Request{
		URL:    &url.URL{},
		Header: map[string][]string{},
		Method: http.MethodPost,
		Body:   ioutil.NopCloser(bytes.NewReader([]byte("6.0.0.0"))),
	}

	ip := getIP(req)

	assert.Equal(t, "6.0.0.0", ip)
}

type errReader struct{}

func (*errReader) Read(_ []byte) (n int, err error) {
	return 0, errors.New("error read data")
}
func (*errReader) Close() error {
	return nil
}

func Test_getIP_error_body(t *testing.T) {
	req := &http.Request{
		URL:    &url.URL{},
		Header: map[string][]string{},
		Method: http.MethodPost,
		Body:   &errReader{},
	}

	ip := getIP(req)

	assert.Equal(t, "", ip)
}

func Test_getIP_no_ip(t *testing.T) {
	req := &http.Request{
		URL:    &url.URL{},
		Header: map[string][]string{},
		Method: http.MethodGet,
	}

	ip := getIP(req)

	assert.Equal(t, "", ip)
}
