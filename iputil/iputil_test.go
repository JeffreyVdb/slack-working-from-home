package iputil

import (
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"testing"
)

type IPResponseMock struct {
	content []byte
}

func (ipResponse *IPResponseMock) Read(buffer []byte) (int, error) {
	return copy(buffer, ipResponse.content), io.EOF
}

type HttpMock struct {
	alwaysReturn []byte
}

func (client *HttpMock) Do(req *http.Request) (*http.Response, error) {
	responseMock := &IPResponseMock{content: client.alwaysReturn}

	return &http.Response{
		StatusCode: http.StatusOK,
		Body: ioutil.NopCloser(responseMock),
	}, nil
}

func TestPlainTextFetchMaker(t *testing.T) {
	client := &HttpMock{alwaysReturn: []byte("192.1.1.1")}
	ipFetcher := PlainTextFetchMaker(client)
	ip, err := ipFetcher("some_service").Get()
	if err != nil {
		t.Errorf("No error should be raised when fetching ip")
		return
	}

	if ip == nil {
		t.Errorf("IP is expected to be 192.1.1.1, but it is nil")
		return
	}

	if !ip.Equal(net.IPv4(192, 1, 1, 1)) {
		t.Errorf("IP is expected to be 192.1.1.1, but it is %s", ip.String())
	}
}
