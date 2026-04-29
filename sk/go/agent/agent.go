package agent

import (
	"github.com/urfave/cli/v2"
	"go.skia.org/infra/sk/go/agent/tool"
	"go.skia.org/infra/sk/go/agent/workflow"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:  "agent",
		Usage: "Commands intended for use by AI agents",
		Subcommands: []*cli.Command{
			tool.Command(),
			workflow.Command(),
		},
	}
}
