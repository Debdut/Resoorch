package gpt

import (
	"fmt"
	"strings"

	"github.com/debdut/Resoorch/lib/api/exa"
)

// Repor represents a report with topics, reports, citations.
type ReportContainer struct {
	Topic     string     `json:"topic"`
	Reports   []Report   `json:"report"`
	Citations []Citation `json:"citations"`
}

// Report represents a report with paragraphs.
type Report struct {
	Paragraphs []Paragraph `json:"paragraph"`
}

type Paragraph struct {
	LineGroups []LineGroup `json:"linegroup"`
}

// Paragraph represents a paragraph with line groups and citations.
type LineGroup struct {
	Line      []string `json:"linegroup"`
	Citations []int    `json:"citations"`
}

// Citation represents a citation with a number, URL, and text.
type Citation struct {
	Num  int    `json:"num"`
	URL  string `json:"url"`
	Text string `json:"text"`
}

func generateInputString(exaResponse *exa.Response) string {
	// Set up the input string
	var sb strings.Builder
	for _, result := range exaResponse.Results {
		sb.WriteString(fmt.Sprintf("title: %s\n", result.Title))
		sb.WriteString(fmt.Sprintf("url: %s\n", result.URL))
		sb.WriteString(fmt.Sprintf("text: %s\n\n", result.Text))
	}
	return sb.String()
}

func generateTextReport(rc *ReportContainer) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Topic: %s\n\n", rc.Topic))
	for _, report := range rc.Reports {
		sb.WriteString("Report:\n")
		for _, paragraph := range report.Paragraphs {
			for _, lineGroup := range paragraph.LineGroups {
				lineGroupString := strings.Join(lineGroup.Line, " ")
				lineGroupString += " ["
				for _, citation := range lineGroup.Citations {
					lineGroupString += fmt.Sprintf(" [%d]", citation)
				}
				lineGroupString += "] "
				sb.WriteString(fmt.Sprintf("%s\n", lineGroupString))
			}
		}
	}
	sb.WriteString("\nCitations:\n")
	for _, citation := range rc.Citations {
		sb.WriteString(fmt.Sprintf("%d. %s\n", citation.Num, citation.Text))
	}
	sb.WriteString("\n\n")
	return sb.String()
}

func (g *GPT) GenerateReport(exaResponse *exa.Response) (string, error) {

	input := generateInputString(exaResponse)

	// Set up the messages for generating the report
	messages := []Message{
		{
			Role: "system",
			Content: `You are an advanced AI assistant tasked with generating a comprehensive report on a given topic. You will be provided with a topic and ten search result documents, each containing a title, URL, and text content. Your task is to synthesize this information into a coherent, well-structured report references at the bottom. The references must be one the documents provided. Your report should be in markdown format.
			Input Format:
Topic: [The main subject of the report]
Document 1: Title: [Title of the first document] URL: [URL of the first document] Text: [Relevant text excerpt from the first document]
Document 2: Title: [Title of the second document] URL: [URL of the second document] Text: [Relevant text excerpt from the second document]
... [Continued for all ten documents]
Instructions:
1. Carefully read and analyze all provided documents.
2. Identify the main themes, key points, and any conflicting information across the sources.
3. Synthesize the information into a cohesive report structure:
    * Introduction: Briefly introduce the topic and its significance.
    * Main Body: Organize the content into logical sections, covering different aspects of the topic.
    * Conclusion: Summarize the key findings and their implications.
4. Use information from multiple sources to support each main point.
5. Cite sources using in-text citations (Author, Year) or [1], [2], etc., corresponding to the document numbers.
6. Maintain an objective tone throughout the report.
7. Include a "References" section at the end, listing all the sources used.
Output Format:
Your report should follow this structure:
Title: [A descriptive title for the report]
1. Introduction [Introduce the topic and its relevance]
2. [Main Section 1] [Content synthesized from the sources]
3. [Main Section 2] [Content synthesized from the sources]
4. [Additional sections as needed]
5. Conclusion [Summarize key findings and their implications]
References:
1. [Document 2]
2. [Document 5] ...
Additional Guidelines:
* Aim for a report length of approximately 1000-1500 words.
* Use clear, concise language appropriate for a general audience.
* Highlight any areas where sources disagree, and present multiple viewpoints if applicable.
* If the provided information is insufficient for any key aspect of the topic, note this in the report.
Generate the report based on these instructions and the provided topic and documents.`,
		},
		{
			Role:    "user",
			Content: input,
		},
	}
	response, err := g.Call(&messages, nil)
	if err != nil {
		return "", fmt.Errorf("error calling GPT: %v", err)
	}
	return response.Choices[0].Message.Content, nil

	// // Set up the report container
	// reportContainer := ReportContainer{}
	// json.Unmarshal(resp.Choices[0].Message.Content, &reportContainer)
	// // Set up the response format
	// responseFormat := ResponseFormat{
	// 	Type: "json_schema",
	// 	JSONSchema: JSONSchema{
	// 		Name:   "report_response",
	// 		Strict: true,
	// 		Schema: Schema{
	// 			Type: "object",
	// 			Properties: map[string]Property{
	// 				"topic": {Type: "string"},
	// 				"report": {
	// 					Type: "array",
	// 					Items: &Property{
	// 						Type: "object",
	// 						Properties: map[string]Property{
	// 							"paragraph": {
	// 								Type: "array",
	// 								Items: &Property{
	// 									Type: "object",
	// 									Properties: map[string]Property{
	// 										"linegroup": {
	// 											Type: "array",
	// 											Items: &Property{
	// 												Type: "object",
	// 												Properties: map[string]Property{
	// 													"line": {
	// 														Type: "array",
	// 														Items: &Property{
	// 															Type: "string",
	// 														},
	// 													},
	// 													"citations": {
	// 														Type: "array",
	// 														Items: &Property{
	// 															Type: "integer",
	// 														},
	// 													},
	// 												},
	// 												Required:             []string{"line", "citations"},
	// 												AdditionalProperties: false,
	// 											},
	// 										},
	// 									},
	// 									Required:             []string{"linegroup"},
	// 									AdditionalProperties: false,
	// 								},
	// 							},
	// 						},
	// 						Required:             []string{"paragraph"},
	// 						AdditionalProperties: false,
	// 					},
	// 				},
	// 				"citations": {
	// 					Type: "array",
	// 					Items: &Property{
	// 						Type: "object",
	// 						Properties: map[string]Property{
	// 							"num":  {Type: "integer"},
	// 							"url":  {Type: "string"},
	// 							"text": {Type: "string"},
	// 						},
	// 						Required:             []string{"num", "url", "text"},
	// 						AdditionalProperties: false,
	// 					},
	// 				},
	// 			},
	// 			Required:             []string{"topic", "report", "citations"},
	// 			AdditionalProperties: false,
	// 		},
	// 	},
	// }

	// // Call the GPT API
	// response, err := g.CallWithResponseFormat(&messages, &responseFormat)
	// if err != nil {
	// 	return "", err
	// }

	// // Unmarshal the report container
	// var reportContainer ReportContainer
	// err = json.Unmarshal([]byte(response.Choices[0].Message.Content), &reportContainer)
	// if err != nil {
	// 	return "", err
	// }

	// fmt.Println(reportContainer)

	// return generateTextReport(&reportContainer), nil
}
