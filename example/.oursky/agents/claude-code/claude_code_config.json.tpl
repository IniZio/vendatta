{
  "mcpServers": {
    "vibegear-mcp": {
      "command": "npx",
      "args": ["-y", "mcp-remote", "http://{{.Host}}:{{.Port}}"],
      "env": {
        "MCP_AUTH_TOKEN": "{{.AuthToken}}"
      }
    }
  }
}