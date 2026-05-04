package tool

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/urfave/cli/v2"
	"go.skia.org/infra/go/auth"
	"go.skia.org/infra/go/httputils"
	"go.skia.org/infra/go/skerr"
	"golang.org/x/oauth2/google"
)

const (
	mcpServerURL            = "https://mcp-skia.luci.app/sse"
	mcpServerOverrideEnvVar = "SK_MCP_SERVER_OVERRIDE"
)

func createCommandsForMCPTools(ctx context.Context) ([]*cli.Command, error) {
	mcpClient, err := initMCP(ctx)
	if err != nil {
		return nil, skerr.Wrap(err)
	}

	tools, err := mcpClient.ListTools(ctx, mcp.ListToolsRequest{})
	if err != nil {
		return nil, skerr.Wrap(err)
	}

	var commands []*cli.Command
	for _, tool := range tools.Tools {
		toolName := tool.Name
		cmd := &cli.Command{
			Name:        toolName,
			Usage:       tool.Description,
			Description: tool.Description,
			Flags:       getFlagsFromSchema(tool.InputSchema),
			Action: func(c *cli.Context) error {
				return callMCPTool(c, toolName)
			},
		}
		commands = append(commands, cmd)
	}

	return commands, nil
}

func getFlagsFromSchema(schema mcp.ToolInputSchema) []cli.Flag {
	var flags []cli.Flag
	for name, prop := range schema.Properties {
		propMap, ok := prop.(map[string]interface{})
		if !ok {
			continue
		}

		description := ""
		if d, ok := propMap["description"]; ok {
			description = fmt.Sprintf("%v", d)
		}

		required := false
		for _, req := range schema.Required {
			if req == name {
				required = true
				break
			}
		}

		propType, _ := propMap["type"].(string)
		switch propType {
		case "boolean":
			flags = append(flags, &cli.BoolFlag{
				Name:  name,
				Usage: description,
			})
		case "number", "integer":
			flags = append(flags, &cli.IntFlag{
				Name:     name,
				Usage:    description,
				Required: required,
			})
		default:
			// Just default to strings.
			flags = append(flags, &cli.StringFlag{
				Name:     name,
				Usage:    description,
				Required: required,
			})
		}
	}
	return flags
}

func callMCPTool(ctx *cli.Context, toolName string) error {
	mcpClient, err := initMCP(ctx.Context)
	if err != nil {
		return skerr.Wrap(err)
	}

	args := make(map[string]interface{})
	for _, flagName := range ctx.FlagNames() {
		if ctx.IsSet(flagName) {
			args[flagName] = ctx.Value(flagName)
		}
	}

	res, err := mcpClient.CallTool(ctx.Context, mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      toolName,
			Arguments: args,
		},
	})
	if err != nil {
		return skerr.Wrap(err)
	}
	var textContent strings.Builder
	for _, content := range res.Content {
		textContent.WriteString(content.(mcp.TextContent).Text)
		textContent.WriteString("\n")
	}
	if res.IsError {
		return fmt.Errorf("tool reported an error: %s", textContent.String())
	} else {
		fmt.Println(textContent.String())
	}
	return nil
}

func initMCP(ctx context.Context) (*client.Client, error) {
	ts, err := google.DefaultTokenSource(ctx, auth.ScopeUserinfoEmail)
	if err != nil {
		return nil, skerr.Wrap(err)
	}
	c := httputils.DefaultClientConfig().WithTokenSource(ts).WithoutRetries().Client()

	mcpURL := os.Getenv(mcpServerOverrideEnvVar)
	if mcpURL == "" {
		mcpURL = mcpServerURL
	}
	mcpClient, err := client.NewSSEMCPClient(mcpURL, transport.WithHTTPClient(c))
	if err != nil {
		return nil, skerr.Wrap(err)
	}

	if err := mcpClient.Start(ctx); err != nil {
		return nil, skerr.Wrap(err)
	}
	_, err = mcpClient.Initialize(ctx, mcp.InitializeRequest{
		Params: mcp.InitializeParams{
			ProtocolVersion: mcp.LATEST_PROTOCOL_VERSION,
		},
	})
	if err != nil {
		return nil, skerr.Wrap(err)
	}
	return mcpClient, nil
}
