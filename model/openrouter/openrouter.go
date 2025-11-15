package openrouter

import (
	"adk-openrouter/internal/converters"
	"context"
	"iter"

	"github.com/revrost/go-openrouter"
	"google.golang.org/adk/model"
)

type openRouterModel struct {
	name   string
	client *openrouter.Client
}

func NewModel(ctx context.Context, modelName string, apiKey string) (model.LLM, error) {
	client := openrouter.NewClient(apiKey)
	return &openRouterModel{
		name:   modelName,
		client: client,
	}, nil
}

func (m *openRouterModel) Name() string {
	return m.name
}

func (m *openRouterModel) GenerateContent(ctx context.Context, req *model.LLMRequest, stream bool) iter.Seq2[*model.LLMResponse, error] {
	if stream {
		return m.generateStream(ctx, req)
	}
	return func(yield func(*model.LLMResponse, error) bool) {
		yield(m.generate(ctx, req))
	}
}

func (m *openRouterModel) generate(ctx context.Context, req *model.LLMRequest) (*model.LLMResponse, error) {
	or := converters.OpenRouterRequestFromLLMRequest(req, m.name, false)
	openRouterResp, err := m.client.CreateChatCompletion(ctx, or)
	if err != nil {
		return nil, err
	}
	res := converters.OpenRouterResponseToLLMResponse(&openRouterResp)
	return res, nil
}

func (m *openRouterModel) generateStream(ctx context.Context, req *model.LLMRequest) iter.Seq2[*model.LLMResponse, error] {
	or := converters.OpenRouterRequestFromLLMRequest(req, m.name, true)
	stream, err := m.client.CreateChatCompletionStream(ctx, or)
	if err != nil {
		return func(yield func(*model.LLMResponse, error) bool) {
			yield(nil, err)
		}
	}

	return func(yield func(*model.LLMResponse, error) bool) {
		if err != nil {
			yield(nil, err)
		}
		defer stream.Close()

		for {
			resp, err := stream.Recv()
			if err != nil {
				break
			}
			yield(converters.OpenRouterStreamResponseToLLMResponse(&resp), nil)
		}
	}
}
