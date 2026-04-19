package vestaboard

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) newRequest(method, url string, body any) (*http.Request, error) {
	var buf io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		buf = bytes.NewReader(data)
	}
	req, err := http.NewRequest(method, url, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Vestaboard-Token", c.token)
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func (c *Client) do(method, url string, body any) ([]byte, error) {
	req, err := c.newRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("vestaboard: HTTP %d: %s", resp.StatusCode, string(data))
	}
	return data, nil
}

func (c *Client) get(url string) ([]byte, error) {
	return c.do(http.MethodGet, url, nil)
}

func (c *Client) post(url string, body any) ([]byte, error) {
	return c.do(http.MethodPost, url, body)
}

func (c *Client) put(url string, body any) ([]byte, error) {
	return c.do(http.MethodPut, url, body)
}
