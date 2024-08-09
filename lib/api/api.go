package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type API struct {
	Name     string
	Codename string
	URL      string
	Key      string
}

func getKey(codename string) (string, error) {
	key := os.Getenv(strings.ToUpper(codename) + "_API_KEY")
	if key == "" {
		return "", fmt.Errorf("please set the " +
			strings.ToUpper(codename) + "_API_KEY environment variable")
	}
	return key, nil
}

func NewAPI(name string, codename string, url string) (*API, error) {
	key, err := getKey(codename)
	api := API{
		Name:     name,
		Codename: codename,
		URL:      url,
		Key:      key,
	}
	return &api, err
}

func (api *API) createRequest(method string, requestBody []byte) (*http.Request, error) {
	// create request object
	req, err := http.NewRequest(method, api.URL, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	// set headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func sendRequest(req *http.Request) ([]byte, int, error) {
	// create http client and send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()
	// read request body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("error reading response body: %v", err)
	}
	return body, resp.StatusCode, nil
}

func (api *API) CreateNSendPostRequest(requestBody []byte, headers [2]string) ([]byte, error) {
	// create request object
	req, err := api.createRequest("POST", requestBody)
	if err != nil {
		return nil, err
	}
	// set headers
	req.Header.Set(headers[0], headers[1])
	// send request
	responseBody, statusCode, err := sendRequest(req)
	if err != nil {
		return nil, err
	}
	// check if http status is ok
	if statusCode != http.StatusOK {
		return responseBody, fmt.Errorf(
			"response error status code: " +
				strconv.Itoa(statusCode))
	}
	return responseBody, nil
}
