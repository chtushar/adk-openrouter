package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/cmd/launcher/adk"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/server/restapi/services"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/geminitool"
	"google.golang.org/genai"
)

const (
	modelName = "gemini-2.5-flash"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	ctx := context.Background()

	// Option 1: Use Gemini (default)
	model, err := gemini.NewModel(ctx, modelName, &genai.ClientConfig{
		APIKey: os.Getenv("GOOGLE_API_KEY"),
	})
	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}

	copyAgent, err := llmagent.New(llmagent.Config{
		Name:        "Marketing copy generation agent",
		Instruction: "You are a marketing copy generator. Generate compelling marketing copy based on the input",
		Model:       model,
		Tools: []tool.Tool{
			geminitool.GoogleSearch{},
		},
		OutputKey: "generated_marketing_copy",
	})

	if err != nil {
		log.Fatalf("Failed to generate marketing copy: %v", err)
		return
	}

	config := &adk.Config{
		AgentLoader: services.NewSingleAgentLoader(copyAgent),
	}

	l := full.NewLauncher()
	err = l.Execute(ctx, config, os.Args[1:])
	if err != nil {
		log.Fatalf("run failed: %v\n\n%s", err, l.CommandLineSyntax())
	}
}
