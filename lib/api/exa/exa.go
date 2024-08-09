package exa

import (
	"encoding/json"
	"fmt"

	"github.com/debdut/Resoorch/lib/api"
)

type Exa struct {
	ID  int
	Api *api.API
}

type Request struct {
	Query      string `json:"query"`
	Type       string `json:"type"`
	NumResults int    `json:"numResults"`
	Contents   struct {
		Text bool `json:"text"`
	} `json:"contents"`
}

type SearchResult struct {
	Score         float64 `json:"score"`
	Title         string  `json:"title"`
	ID            string  `json:"id"`
	URL           string  `json:"url"`
	PublishedDate string  `json:"publishedDate"`
	Author        string  `json:"author"`
	Text          string  `json:"text"`
}

type Response struct {
	AutopromptString   string         `json:"autopromptString"`
	ResolvedSearchType string         `json:"resolvedSearchType"`
	Results            []SearchResult `json:"results"`
	RequestID          string         `json:"requestId"`
}

const Endpoint = "https://api.exa.ai/search"

func InitExa(id int) (*Exa, error) {
	api, err := api.NewAPI("Exa", "exa",
		Endpoint)
	if err != nil {
		return nil, err
	}
	return &Exa{ID: id, Api: api}, nil
}

func (exa *Exa) request(request *Request) (*Response, error) {
	// marshal body
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}
	// send request and response
	headers := [2]string{"X-API-KEY", exa.Api.Key}
	responseBody,
		err := exa.Api.CreateNSendPostRequest(requestBody, headers)
	if err != nil {
		return nil, err
	}
	// unmarshall response
	var response Response
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}
	return &response, nil
}

func (exa *Exa) Search(query string, numResults int) (*Response, error) {
	request := Request{
		Query:      query,
		Type:       "auto",
		NumResults: numResults,
		Contents: struct {
			Text bool `json:"text"`
		}{
			Text: true,
		},
	}
	return exa.request(&request)
}
