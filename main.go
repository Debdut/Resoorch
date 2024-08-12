package main

import (
	"strings"

	"github.com/neurosnap/sentences"
)

// ExtractRelevantParts finds the most relevant sentences in the text
func GetSentences(text string) []string {
	// Tokenize sentences
	storage := sentences.NewStorage()
	tokenizer := sentences.NewSentenceTokenizer(storage)
	tokenizedSentences := tokenizer.Tokenize(text)

	var sentences []string
	for _, sentence := range tokenizedSentences {
		sentences = append(sentences, sentence.Text)
	}

	return sentences
}

// ChunkText breaks the text into chunks of maxChunkSize
func ChunkText(text string, maxChunkSize int) []string {
	// Split text into paragraphs
	paragraphs := strings.Split(text, "\n")

	var chunks []string
	for _, paragraph := range paragraphs {
		sentences := GetSentences(paragraph)
		var currentChunk []string
		currentChunkSize := 0

		for _, sentence := range sentences {
			sentenceSize := len(strings.Fields(sentence))
			if currentChunkSize+sentenceSize <= maxChunkSize {
				currentChunk = append(currentChunk, sentence)
				currentChunkSize += sentenceSize
			} else {
				chunks = append(chunks, strings.Join(currentChunk, " "))
				currentChunk = []string{sentence}
				currentChunkSize = sentenceSize
			}
		}

		if len(currentChunk) > 0 {
			chunk := strings.TrimSpace(strings.Join(currentChunk, " "))
			if chunk != "" {
				chunks = append(chunks, chunk)
			}
		}
	}

	return chunks
}

func main() {
	text := `Well fuck that bitch.
	 Python is a high-level, interpreted programming language known for its simplicity and readability. 
	It was created by Guido van Rossum and first released in 1991. Python supports multiple programming paradigms, including procedural, object-oriented, and functional programming. It has a large and comprehensive standard library, often described as having "batteries included". Python is widely used in various fields such as web development, data analysis, artificial intelligence, and scientific computing.
	
	Its syntax allows programmers to express concepts in fewer lines of code than would be possible in languages such as C++ or Java. Python's design philosophy emphasizes 
	code readability with its notable use of significant whitespace.`

	chunks := ChunkText(text, 512)
	for i, chunk := range chunks {
		println(i, "[[", chunk, "]]")
	}
}
