package pocketsmith

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var ErrNotFound = errors.New("not found")

type ApiError struct {
	Err string `json:"error"`
}

func (a ApiError) Error() string {
	return fmt.Sprintf("Pocketsmith API Error: %s", a.Err)
}

type Client struct {
	token string
}

func NewClient(token string) *Client {
	return &Client{token: token}
}

func (c *Client) doAndDecode(req *http.Request, responseType any) error {
	req.Header.Add("accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("X-Developer-Key", c.token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bodyBuf := new(bytes.Buffer)
	bodyBuf.ReadFrom(resp.Body)

	reader := bytes.NewReader(bodyBuf.Bytes())

	// b := bytes.Buffer{}
	// b.ReadFrom(reader)
	// fmt.Println("Body: ", string(b.Bytes()))
	// reader.Seek(0, 0)

	var apiError ApiError
	if err := json.NewDecoder(reader).Decode(&apiError); err == nil {
		if apiError.Err != "" {
			return apiError
		}
	}

	reader.Seek(0, 0)
	if responseType != nil {
		return json.NewDecoder(reader).Decode(responseType)
	}

	return nil
}
