package converters

import (
	"github.com/revrost/go-openrouter"
	"google.golang.org/adk/model"
	"google.golang.org/genai"
)

func OpenRouterResponseToLLMResponse(res *openrouter.ChatCompletionResponse) *model.LLMResponse {
	parts := make([]*genai.Part, 1)
	parts[0] = &genai.Part{
		Text: res.Choices[0].Message.Content.Text,
	}
	return &model.LLMResponse{
		Content: &genai.Content{
			Role:  "assistant",
			Parts: parts,
		},
	}
}

func OpenRouterStreamResponseToLLMResponse(res *openrouter.ChatCompletionStreamResponse) *model.LLMResponse {
	parts := make([]*genai.Part, 1)

	parts[0] = &genai.Part{
		Text: res.Choices[0].Delta.Content,
	}

	return &model.LLMResponse{
		Content: &genai.Content{
			Role:  "assistant",
			Parts: parts,
		},
	}
}

func OpenRouterRequestFromLLMRequest(req *model.LLMRequest, modelName string, stream bool) openrouter.ChatCompletionRequest {
	var messages []openrouter.ChatCompletionMessage
	for _, content := range req.Contents {
		for _, part := range content.Parts {
			messages = append(messages, openrouter.ChatCompletionMessage{
				Role: content.Role,
				Content: openrouter.Content{
					Text: part.Text,
				},
			})
		}
	}

	return openrouter.ChatCompletionRequest{
		Model:    modelName,
		Stream:   stream,
		Messages: messages,
	}
}
