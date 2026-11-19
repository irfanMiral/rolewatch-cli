package cli

import (
	"os"
	"strings"

	"github.com/spf13/cobra"

	"rolewatch-cli/internal/client"
	"rolewatch-cli/internal/config"
	"rolewatch-cli/internal/mcpproxy"
)

func newMCPCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "mcp",
		Short: "Run RoleWatch's MCP server over stdio for a local AI host",
		Long: "mcp connects a local MCP host (Claude Code, Claude Desktop, Cursor) to RoleWatch's " +
			"hosted MCP server, signed in as the account `rolewatch auth login` stored — the same " +
			"tools, prompts and instructions as the remote connector, with no second sign-in. " +
			"Without a stored sign-in it serves the public catalogue tools only. The server also " +
			"tells the agent this CLI is on its shell, for results that belong on disk.\n\n" +
			"  claude mcp add rolewatch -- rolewatch mcp",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			r, err := config.ResolveOptional(os.Getenv)
			if err != nil {
				return err
			}
			base := r.APIURL
			if f, _ := cmd.Flags().GetString("api-url"); f != "" {
				base = f
			}
			relay := &mcpproxy.Relay{
				Endpoint:  strings.TrimRight(base, "/") + "/api/v1/mcp",
				Token:     r.Token,
				UserAgent: client.UserAgent,
			}
			return relay.Run(cmd.Context(), cmd.InOrStdin(), cmd.OutOrStdout())
		},
	}
}
