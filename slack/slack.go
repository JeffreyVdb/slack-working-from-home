package slack

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"github.com/JeffreyVdb/slack-working-from-home/util"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const slackBaseURL = "https://slack.com/api/"

type SlackStatus struct {
	StatusText  string
	StatusEmoji string
}

type ProfileInfo struct {
	FirstName   string `json:"first_name,omitempty"`
	LastName    string `json:"last_name,omitempty"`
	Email       string `json:"email,omitempty"`
	StatusText  string `json:"status_text"`
	StatusEmoji string `json:"status_emoji"`
}

type httpClient interface {
	Do(request *http.Request) (response *http.Response, err error)
}

type Client struct {
	httpClient httpClient
	timeout    time.Duration
	apiToken   string
	baseURL    *url.URL
}

type ClientOption func(*Client) error

func Timeout(seconds int) ClientOption {
	return func(slackClient *Client) error {
		slackClient.timeout = time.Second * time.Duration(seconds)
		return nil
	}
}

func NewClient(token string, options ...ClientOption) (*Client, error) {
	baseURL, err := url.Parse(slackBaseURL)
	if err != nil {
		return nil, err
	}
	client := &Client{apiToken: token, baseURL: baseURL}

	for _, callback := range options {
		err := callback(client)
		if err != nil {
			return nil, err
		}
	}

	if client.httpClient == nil {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: util.IsEnvDefined("http_proxy") || util.IsEnvDefined("https_proxy")},
			Proxy:           http.ProxyFromEnvironment,
		}
		client.httpClient = &http.Client{Timeout: client.timeout, Transport: tr}
	}

	return client, nil
}

func (slackClient *Client) NewRequest(method, endpoint string, body interface{}) (*http.Request, error) {
	var buf io.ReadWriter
	relURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	relURL.Path = strings.TrimLeft(relURL.Path, "/")
	absURL := slackClient.baseURL.ResolveReference(relURL)
	if body != nil {
		buf = new(bytes.Buffer)
		err = json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, absURL.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+slackClient.apiToken)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	return req, nil
}

func (slackClient *Client) SetProfileStatus(status *SlackStatus) error {
	endpoint := "users.profile.set"
	payload := struct {
		Profile *ProfileInfo `json:"profile"`
	}{
		&ProfileInfo{
			StatusText:  status.StatusText,
			StatusEmoji: status.StatusEmoji,
		},
	}
	req, err := slackClient.NewRequest("POST", endpoint, payload)
	if err != nil {
		return err
	}

	resp, err := slackClient.httpClient.Do(req)
	if err != nil {
		return err
	}

	// TODO: Check if resp returns an error
	if resp.StatusCode != 200 {
		return errors.New("status code is not 200")
	}

	return nil
}
