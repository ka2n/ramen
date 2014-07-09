package yo

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

type Client struct {
	URL      string
	Username string
	Token    string
	Debug    bool
}

const YO_ENDPOINT = "http://api.justyo.co/"

var DefaultClient = NewClient(YO_ENDPOINT, os.Getenv("YO_TOKEN"), os.Getenv("YO_USERNAME"))

func NewClient(url string, token string, username string) *Client {
	client := Client{URL: url, Token: token, Username: username}
	return &client
}

func (c *Client) APIRequest(path string, data url.Values) error {
	hc := http.DefaultClient
	req, err := c.NewRequest(path, data)
	if c.Debug {
		log.Println(req)
	}
	if err != nil {
		return err
	}

	res, err := hc.Do(req)
	if c.Debug {
		log.Println(res)
	}
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode/100 != 2 {
		return errors.New(fmt.Sprintf("yo request failed with StatusCode: %d", res.StatusCode))
	}
	return nil
}

func (c *Client) NewRequest(path string, data url.Values) (*http.Request, error) {
	data.Set("api_token", c.Token)
	req, err := http.NewRequest("POST", c.URL+path, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))
	return req, nil
}

func (c *Client) YoAll() error {
	if c.Debug {
		log.Print("Sending yo all")
	}
	err := c.APIRequest("yoall/", url.Values{})
	if err != nil {
		return err
	}
	if c.Debug {
		log.Print("Sent yo all")
	}
	return nil
}

func (c *Client) Yo(username string) error {
	if c.Debug {
		log.Printf("Sending yo to %s", username)
	}
	if username == "" {
		return errors.New("username required to Yo")
	}

	err := c.APIRequest("yo/", url.Values{"username": []string{username}})
	if err != nil {
		return err
	}
	if c.Debug {
		log.Printf("Sent yo to %s", username)
	}
	return nil
}
