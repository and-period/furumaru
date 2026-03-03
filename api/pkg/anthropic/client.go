package anthropic

import (
	"context"
	"encoding/json"

	sdk "github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/packages/ssestream"
)

type Client struct {
	client sdk.Client
	model  sdk.Model
}

type Message struct {
	Role    string
	Content []ContentBlock
}

type ContentBlock struct {
	Type string // "text", "tool_use", "tool_result"
	// text
	Text string
	// tool_use
	ID    string
	Name  string
	Input string // JSON string
	// tool_result
	ToolUseID string
	Content   string
}

type Tool struct {
	Name        string
	Description string
	InputSchema map[string]interface{}
}

type StreamParams struct {
	System   string
	Messages []Message
	Tools    []Tool
}

func NewClient(apiKey, model string) *Client {
	client := sdk.NewClient(option.WithAPIKey(apiKey))
	return &Client{
		client: client,
		model:  sdk.Model(model),
	}
}

func (c *Client) StreamMessage(ctx context.Context, params *StreamParams) *ssestream.Stream[sdk.MessageStreamEventUnion] {
	messages := make([]sdk.MessageParam, 0, len(params.Messages))
	for _, msg := range params.Messages {
		blocks := make([]sdk.ContentBlockParamUnion, 0, len(msg.Content))
		for _, cb := range msg.Content {
			switch cb.Type {
			case "text":
				blocks = append(blocks, sdk.NewTextBlock(cb.Text))
			case "tool_use":
				var input interface{}
				if err := json.Unmarshal([]byte(cb.Input), &input); err != nil {
					input = map[string]interface{}{}
				}
				blocks = append(blocks, sdk.NewToolUseBlock(cb.ID, input, cb.Name))
			case "tool_result":
				blocks = append(blocks, sdk.NewToolResultBlock(cb.ToolUseID, cb.Content, false))
			}
		}
		messages = append(messages, sdk.MessageParam{
			Role:    sdk.MessageParamRole(msg.Role),
			Content: blocks,
		})
	}

	tools := make([]sdk.ToolUnionParam, 0, len(params.Tools))
	for _, t := range params.Tools {
		tools = append(tools, sdk.ToolUnionParam{
			OfTool: &sdk.ToolParam{
				Name:        t.Name,
				Description: sdk.String(t.Description),
				InputSchema: sdk.ToolInputSchemaParam{
					Properties: t.InputSchema,
				},
			},
		})
	}

	msgParams := sdk.MessageNewParams{
		Model:     c.model,
		MaxTokens: 4096,
		Messages:  messages,
		Tools:     tools,
	}

	if params.System != "" {
		msgParams.System = []sdk.TextBlockParam{
			{Text: params.System},
		}
	}

	return c.client.Messages.NewStreaming(ctx, msgParams)
}
