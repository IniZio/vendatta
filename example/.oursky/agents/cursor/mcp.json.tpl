{
  "mcpServers": {
    "vibegear-mcp": {
      "type": "http",
      "url": "http://{{.Host}}:{{.Port}}",
      "headers": {
        "Authorization": "Bearer {{.AuthToken}}"
      }
    }
  }
}