package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/czc09/langchaingo/llms"
	"github.com/czc09/langchaingo/llms/ollama"
)

// readFile reads the content of a file given its path.
func readFile(filePath string) (string, error) {
	codeBytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("error reading file %s: %v", filePath, err)
	}
	return string(codeBytes), nil
}

// analyzeCode sends the code and system message to the LLM for analysis.
func analyzeCode(ctx context.Context, llm llms.Model, systemMessage, code string) error {
	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, systemMessage),
		llms.TextParts(llms.ChatMessageTypeHuman, code),
	}

	// Streaming function to handle LLM output
	streamingFunc := func(ctx context.Context, chunk []byte) error {
		fmt.Print(string(chunk))
		return nil
	}

	// Generate content with streaming
	_, err := llm.GenerateContent(ctx, content, llms.WithStreamingFunc(streamingFunc))
	return err
}

func main() {
	// Check command-line arguments
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <code-file>")
	}

	// Get the current working directory
	basePath, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current directory: %v\n", err)
	}

	// Read the system message from file
	systemFilePath := filepath.Join(basePath, "systemmessage.txt")
	systemMessage, err := readFile(systemFilePath)
	if err != nil {
		log.Fatalf("Failed to read system message: %v\n", err)
	}

	// Read the code from file
	codeFilePath := os.Args[1]
	code, err := readFile(codeFilePath)
	if err != nil {
		log.Fatalf("Failed to read code file: %v\n", err)
	}

	// Initialize the Ollama LLM
	llm, err := ollama.New(
		ollama.WithModel("llama3.2:1b"), // Replace with your desired model
		ollama.WithServerURL("http://localhost:11434"),
	)
	if err != nil {
		log.Fatalf("Failed to initialize LLM: %v\n", err)
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute) // Adjust timeout as needed
	defer cancel()

	// Analyze the code
	if err := analyzeCode(ctx, llm, systemMessage, code); err != nil {
		log.Fatalf("Error analyzing code: %v\n", err)
	}
}