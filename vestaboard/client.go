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
	boardType  BoardType
	httpClient *http.Client
}


func New(token string) *Client {
	return newWithURLs(token, defaultBaseURL, defaultVbmlURL, BoardFlagship)
}


func NewNote(token string) *Client {
	return newWithURLs(token, defaultBaseURL, defaultVbmlURL, BoardNote)
}

func newWithURLs(token, baseURL, vbmlURL string, boardType BoardType) *Client {
	return &Client{
		token:      token,
		baseURL:    baseURL,
		vbmlURL:    vbmlURL,
		boardType:  boardType,
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


func (c *Client) GetTransition() (*TransitionInfo, error) {
	data, err := c.get(c.baseURL + "/transition")
	if err != nil {
		return nil, err
	}
	var result TransitionInfo
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// updates the transition animation and speed.
func (c *Client) SetTransition(t Transition, speed TransitionSpeed) (*TransitionInfo, error) {
	body := map[string]any{"transition": t, "transitionSpeed": speed}
	data, err := c.put(c.baseURL+"/transition", body)
	if err != nil {
		return nil, err
	}
	var result TransitionInfo
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// vbml -> board array
func (c *Client) Compose(req ComposeRequest) (BoardLayout, error) {
	if req.Style == nil {
		req.Style = &BoardStyle{}
	}
	if req.Style.Height == 0 {
		req.Style.Height = c.boardType.Rows()
	}
	if req.Style.Width == 0 {
		req.Style.Width = c.boardType.Cols()
	}
	data, err := c.post(c.vbmlURL+"/compose", req)
	if err != nil {
		return nil, err
	}
	var layout BoardLayout
	if err := json.Unmarshal(data, &layout); err != nil {
		return nil, err
	}
	return layout, nil
}

// Bundle of compose and send
func (c *Client) ComposeAndSend(req ComposeRequest, forced bool) (*SendResult, error) {
	layout, err := c.Compose(req)
	if err != nil {
		return nil, err
	}
	return c.SendCharacters(layout, forced)
}

// Nicely formats strings (this is like compose but simpler)
func (c *Client) FormatMessage(text string) (BoardLayout, error) {
	data, err := c.post(c.vbmlURL+"/format", map[string]any{"message": text})
	if err != nil {
		return nil, err
	}
	var layout BoardLayout
	if err := json.Unmarshal(data, &layout); err != nil {
		return nil, err
	}
	return layout, nil
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
