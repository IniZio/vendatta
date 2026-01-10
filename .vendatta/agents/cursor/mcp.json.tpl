{
  "mcpServers": {
    "vibegear": {
      "command": "vendatta",
      "args": ["agent", "{{.ProjectName}}"],
      "env": {
        "MCP_PORT": "{{.Port}}"
      }
    }
  }
}
