package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/czc09/langchaingo/llms"
	"github.com/czc09/langchaingo/llms/ollama"
)

func checkArgs() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <code-file>")
		return
	}
}

func readSystemMessage() (string, error){
	// Read System message
	// Get the absolute path to the systemmessage.txt file
	basePath, err := os.Getwd() // Get the current working directory
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		return "", err
	}
	systemfilePath := filepath.Join(basePath, "systemmessage.txt")
	systemMessageByte, err := os.ReadFile(systemfilePath)

	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return "", err
	}

	return string(systemMessageByte), nil
}


func main() {

	checkArgs()

	systemmessage, err := readSystemMessage()
	
	if err != nil {
		fmt.Println("Failed to read system message:", err)
		return
	}

	// Read code from file
	filePath := os.Args[1]

	codeBytes, err := os.ReadFile(filePath)

	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	code := string(codeBytes)

	llm, err := ollama.New(ollama.WithModel("llama3.2:1b"), ollama.WithServerURL("http://localhost:11434"))

	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, systemmessage),
		llms.TextParts(llms.ChatMessageTypeHuman, code),
	}

	completion, err := llm.GenerateContent(ctx, content, llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
		fmt.Print(string(chunk))
		return nil
	}))

	if err != nil {
		log.Fatal(err)
	}
	_ = completion
}
