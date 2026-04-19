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


// SendText does not utilize styling of any sort
func (c *Client) SendText(text string, forced bool) (*SendResult, error) {
	body := map[string]any{"text": text}
	if forced {
		body["forced"] = true
	}
	data, err := c.post(c.baseURL+"/", body)
	if err != nil {
		return nil, err
	}
	var result SendResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SendCharacters sends a raw character-code layout (array of arrays of ints) to the board.
func (c *Client) SendCharacters(layout BoardLayout, forced bool) (*SendResult, error) {
	body := map[string]any{"characters": layout}
	if forced {
		body["forced"] = true
	}
	data, err := c.post(c.baseURL+"/", body)
	if err != nil {
		return nil, err
	}
	var result SendResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
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
