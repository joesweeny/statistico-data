package sportmonks

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
)

type Client struct {
	BaseUri string
	ApiKey string
}

func NewClient(uri, key string) *Client {
	return &Client{
		BaseUri: uri,
		ApiKey: key,
	}
}

func (c *Client) GetCountries() (*CountryResponse, error) {
	response, err := http.Get(c.BaseUri + "/countries?api_token=" + c.ApiKey)

	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		panic(err.Error())
	}

	res, err := parseCountries([]byte(body))

	return res, err
}

func parseCountries(body []byte) (*CountryResponse, error) {
	var r = new(CountryResponse)
	err := json.Unmarshal(body, &r)

	if err != nil {
		panic(err.Error())
	}

	return r, err
}
