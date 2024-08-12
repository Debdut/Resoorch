package gpt

import (
	"encoding/json"
	"fmt"

	"github.com/debdut/Resoorch/lib/api"
)

type GPT struct {
	ID  int
	Api *api.API
}

// structs for request
type Request struct {
	Model          string          `json:"model"`
	Messages       []Message       `json:"messages"`
	ResponseFormat *ResponseFormat `json:"response_format,omitempty"`
	MaxTokens      int             `json:"max_tokens,omitempty"`
	Temperature    float64         `json:"temperature,omitempty"`
	TopP           float64         `json:"top_p,omitempty"`
	N              int             `json:"n,omitempty"`
}

type Message struct {
	Role    string  `json:"role,omitempty"`
	Content string  `json:"content,omitempty"`
	Refusal *string `json:"refusal,omitempty"`
}

type ResponseFormat struct {
	Type       string      `json:"type,omitempty"`
	JSONSchema *JSONSchema `json:"json_schema,,omitempty"`
}

type JSONSchema struct {
	Name   string  `json:"name,omitempty"`
	Strict bool    `json:"strict,omitempty"`
	Schema *Schema `json:"schema,omitempty"`
}

type Schema struct {
	Type                 string              `json:"type,omitempty"`
	Properties           map[string]Property `json:"properties,omitempty"`
	Required             []string            `json:"required,omitempty"`
	AdditionalProperties bool                `json:"additionalProperties,omitempty"`
}

type Property struct {
	Type                 string              `json:"type,omitempty"`
	Items                *Property           `json:"items,omitempty"`
	Properties           map[string]Property `json:"properties,omitempty"`
	AdditionalProperties bool                `json:"additionalProperties"`
	Required             []string            `json:"required,omitempty"`
}

// Struct for response
type Response struct {
	ID       string       `json:"id"`
	Object   string       `json:"object"`
	Created  int64        `json:"created"`
	Model    string       `json:"model"`
	Choices  []Choice     `json:"choices"`
	Usage    UsageStats   `json:"usage"`
	SystemFP string       `json:"system_fingerprint"`
	Error    ErrorDetails `json:"error"`
}

type ErrorDetails struct {
	Message string  `json:"message"`
	Type    string  `json:"type"`
	Param   *string `json:"param"`
	Code    *string `json:"code"`
}

type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	LogProbs     *string `json:"logprobs"`
	FinishReason string  `json:"finish_reason"`
}

type UsageStats struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

const Model = "gpt-4o-2024-08-06"
const Endpoint = "https://api.openai.com/v1/chat/completions"

func InitGPT(id int) (*GPT, error) {
	api, err := api.NewAPI("GPT", "openai", Endpoint)
	if err != nil {
		return nil, err
	}
	return &GPT{ID: id, Api: api}, nil
}

func (gpt *GPT) request(request *Request) (*Response, error) {
	// marshal body
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}
	// send request and response
	headers := [2]string{"Authorization", "Bearer " + gpt.Api.Key}
	responseBody,
		err := gpt.Api.CreateNSendPostRequest(requestBody, headers)
	if err != nil {
		if responseBody == nil {
			return nil, err
		}
	}
	// // unmarshall response
	var response Response
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}
	return &response, err
}

func (gpt *GPT) Call(messages *[]Message, responseFormat *ResponseFormat) (*Response, error) {
	request := Request{
		Model:    Model,
		Messages: *messages,
	}
	if responseFormat != nil {
		request.ResponseFormat = responseFormat
	}
	return gpt.request(&request)
}
