{
  "$schema": "https://opencode.ai/config.json",
  "model": "anthropic/claude-sonnet-4-5",
  "mcp": {
    "vibegear-mcp": {
      "type": "remote",
      "url": "http://{{.Host}}:{{.Port}}",
      "enabled": true,
      "headers": {
        "Authorization": "Bearer {{.AuthToken}}"
      }
    }
  },
  "rules": {{.RulesConfig}},
  "skills": {{.SkillsConfig}},
  "commands": {{.CommandsConfig}},
  "tools": {
    "mcp": true
  }
}