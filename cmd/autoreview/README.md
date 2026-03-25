# Auto Review (`cmd/autoreview`)

AI-powered code review tool. It analyzes the latest commit or your current
working directory changes using Google's Gemini API and provides a summary and
an LGTM status.

## Usage

Build and run the tool:

```bash
bazel run //cmd/autoreview
```

Display a help message:

```bash
bazel run //cmd/autoreview -- --help
```

Pass command line options:

```bash
bazel run //cmd/autoreview -- [options]
```

## Authentication Setup

The tool uses Google Cloud application-default credentials.
Run the following commands to authenticate via Google Cloud:

```bash
gcloud auth application-default set-quota-project skia-infra-corp
gcloud auth application-default login
```

See more go/gemini-api/authentication.

## MCP Server (`cmd/mcp`)

The MCP (Model Context Protocol) server exposes tools that allow AI assistants
(like Jetski, Gemini CLI, etc.) to run repository-specific commands. It wraps
`cmd/autoreview` to let the AI assistant trigger a code review on the codebase
itself.

### Gemini CLI Configuration

To configure Gemini CLI to use this MCP server, create or update your
`~/.gemini/settings.json` by adding the `buildbot` field to the `mcpServers`
field as shown below (set the correct path to the `buildbot` repo on your
machine).

```json
{
  "mcpServers": {
    "buildbot": {
      "command": "/full/path/to/buildbot/cmd/autoreview/mcp/run.sh"
    }
  }
}
```

Restart Gemini CLI to use the `buildbot` MCP server. Verify by entering
`/mcp list` in the Gemini CLI window.

### Jetski Configuration (Recommended)

To configure Jetski to use this MCP server, create or update your
`~/.gemini/jetski/mcp_config.json`:

```json
{
  "mcpServers": {
    "buildbot": {
      "command": "/full/path/to/buildbot/cmd/autoreview/mcp/run.sh"
    }
  }
}
```

Reload Jetski by `Ctrl(Cmd)+Shift+P` -> `Kill Language Server and Reload Window`
to start using the `buildbot` MCP server.
