# Task: USR-04 Deep link generation for editors

**Priority**: 🔥 High
**Status**: [Pending]

## 🎯 Objective
Implement deep link generation system for IDEs and editors to enable seamless connection to remote workspaces with proper authentication and context.

## 🛠 Implementation Details

### **Deep Link Protocols**
1. **VS Code Remote SSH**: `vscode://vscode-remote/ssh-remote+{host}{workspace_path}`
   - SSH host configuration with user credentials
   - Workspace path resolution
   - Extension recommendations

2. **Cursor**: `cursor://file/{workspace_path}` with SSH context
   - Integration with Cursor's remote development features
   - MCP server configuration injection

3. **JetBrains IDEs**: Custom protocol handlers
   - Gateway configuration for remote development
   - Project structure detection

### **Link Generation Service** (`pkg/user/deeplinks.go`)
- **generateVSCodeLink()**: VS Code remote SSH links
- **generateCursorLink()**: Cursor-specific deep links
- **generateJetBrainsLink()**: IntelliJ/CLion/etc. remote links
- Context-aware link generation based on workspace contents

### **Authentication Integration**
- SSH key injection into link parameters
- Temporary authentication tokens for web-based editors
- Secure parameter encoding and validation

### **Workspace Context Detection**
- Language/framework detection from files
- Recommended extensions based on tech stack
- IDE-specific configuration generation
- Project structure analysis

### **CLI Integration** (`cmd/vendatta/workspace.go`)
- `vendatta workspace link --ide vscode`: Generate IDE-specific links
- `vendatta workspace open --ide cursor`: Open workspace in specified IDE
- Link sharing and export functionality

## 🧪 Proof of Work
- [ ] Deep link generation for major IDEs
- [ ] Workspace context analysis
- [ ] Authentication parameter handling
- [ ] CLI integration for link generation
- [ ] Cross-platform link opening
- [ ] Integration with user workspace assignments
