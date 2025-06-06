package helloworld

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"go.skia.org/infra/mcp/common"
)

type HelloWorldService struct {
}

// Initialize the service with the provided arguments.
func (s HelloWorldService) Init(serviceArgs string) error {
	return nil
}

// GetTools returns the supported tools by the service.
func (s HelloWorldService) GetTools() []common.Tool {
	return []common.Tool{
		{
			Name:        "sayhello",
			Description: "Says hello to the caller.",
			Arguments: []common.ToolArgument{
				{
					Name:        "name",
					Description: "Name of the user",
					Required:    true,
				},
			},
			Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
				name, err := request.RequireString("name")
				if err != nil {
					return mcp.NewToolResultError(err.Error()), nil
				}

				return mcp.NewToolResultText(fmt.Sprintf("Hello, %s!", name)), nil
			},
		},
	}
}
