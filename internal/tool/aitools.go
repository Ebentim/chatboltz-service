package tool

import (
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/shared"
)

// Define a function tool
var Tools = []openai.ChatCompletionToolParam{
	{
		Function: shared.FunctionDefinitionParam{
			Name:        "get_weather",
			Description: openai.String("Get weather for a location"),
			Parameters: (map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"location": map[string]interface{}{
						"type":        "string",
						"description": "City name",
					},
				},
				"required": []string{"location"},
			}),
		},
	},
}

// Use with tools
// completion, err := provider.GetCompletionWithTools(ctx, conv, "You are a helpful assistant", "What's the weather in Paris?", "", tools)

// var Gtls = []genai.FunctionDeclaration{}
// var Shoes = []anthropic.ToolParam{}
