package pocketsmith

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

var ErrNotFound = errors.New("not found")

type Client struct {
	token string
}

func NewClient(token string) *Client {
	return &Client{token: token}
}

func (c *Client) doAndDecode(req *http.Request, responseType any) error {
	req.Header.Add("accept", "application/json")
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

	//b := bytes.Buffer{}
	//b.ReadFrom(reader)
	//fmt.Println("Body: ", string(b.Bytes()))
	//reader.Seek(0, 0)

	var apiError *ApiError
	if err := json.NewDecoder(reader).Decode(&apiError); err == nil {
		if apiError != nil {
			return apiError
		}
	}

	if responseType == nil {
		return nil
	}

	reader.Seek(0, 0)
	return json.NewDecoder(reader).Decode(responseType)
}
