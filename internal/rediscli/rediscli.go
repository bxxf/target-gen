package rediscli

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type UpstashClient struct {
	baseURL string
	token   string
}

func NewUpstashClient(baseURL string, token string) *UpstashClient {
	return &UpstashClient{baseURL: baseURL, token: token}
}

type BrandCountryResponse struct {
	Result string `json:"result"`
}

func (c *UpstashClient) GetCountries(brand string) ([]string, error) {
	url := fmt.Sprintf("%s/get/BRAND:%s", c.baseURL, brand)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Add Authorization header
	req.Header.Add("Authorization", c.token)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response BrandCountryResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	countries := strings.Split(response.Result, ",")

	return countries, nil
}
