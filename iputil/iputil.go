package iputil

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
)

type httpClient interface {
	Do(request *http.Request) (response *http.Response, err error)
}

func parsePlainTextIP(reader io.Reader) (net.IP, error) {
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	ipAddr := net.ParseIP(string(content))
	if ipAddr == nil {
		return nil, fmt.Errorf("could not parse %s as an IP", content)
	}

	return ipAddr, nil
}

func PlainTextFetchMaker(client httpClient) func(endpoint string) *IPFetcher {
	return func(endpoint string) *IPFetcher {
		return &IPFetcher{
			client:    client,
			endpoint:  endpoint,
			parseFunc: parsePlainTextIP,
		}
	}
}

type IPFetcher struct {
	client httpClient
	endpoint string
	parseFunc func(reader io.Reader) (net.IP, error)
}

func (fetcher *IPFetcher) Get() (net.IP, error) {
	req, err := http.NewRequest("GET", fetcher.endpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := fetcher.client.Do(req)
	if err != nil {
		return nil, err
	}

	ipAddr, err := fetcher.parseFunc(resp.Body)
	if err != nil {
		_ = resp.Body.Close()
		return nil, err
	}

	return ipAddr, resp.Body.Close()
}
