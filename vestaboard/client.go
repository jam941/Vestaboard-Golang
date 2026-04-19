package vestaboard

import (
	"encoding/json"
	"net/http"
)

const (
	defaultBaseURL = "https://cloud.vestaboard.com"
	defaultVbmlURL = "https://vbml.vestaboard.com"
)

type Client struct {
	token      string
	baseURL    string
	vbmlURL    string
	httpClient *http.Client
}

func New(token string) *Client {
	return newWithURLs(token, defaultBaseURL, defaultVbmlURL)
}

func newWithURLs(token, baseURL, vbmlURL string) *Client {
	return &Client{
		token:      token,
		baseURL:    baseURL,
		vbmlURL:    vbmlURL,
		httpClient: &http.Client{},
	}
}


func (c *Client) GetMessage() (*MessageResult, error) {
	data, err := c.get(c.baseURL + "/")
	if err != nil {
		return nil, err
	}

	var envelope struct {
		CurrentMessage struct {
			ID     string      `json:"id"`
			Layout BoardLayout `json:"layout"`
		} `json:"currentMessage"`
	}
	if err := json.Unmarshal(data, &envelope); err != nil {
		return nil, err
	}

	return &MessageResult{
		ID:     envelope.CurrentMessage.ID,
		Layout: envelope.CurrentMessage.Layout,
	}, nil
}
