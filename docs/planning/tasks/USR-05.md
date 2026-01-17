# Task: USR-05 Service port discovery and display

**Priority**: 🔥 High
**Status**: [Pending]

## 🎯 Objective
Implement automatic service port discovery and user-friendly display system to help users access running services in remote workspaces.

## 🛠 Implementation Details

### **Port Discovery Mechanisms**
1. **Container Inspection**: Docker/LXC port mapping detection
   - Runtime inspection of container configurations
   - Port binding analysis from docker inspect / lxc-info
   - Dynamic port allocation tracking

2. **Process Monitoring**: Service process port detection
   - Netstat/ss command execution on remote nodes
   - Process-to-port mapping via PID analysis
   - Listening socket enumeration

3. **Configuration-Based Discovery**: Declarative port definitions
   - Service configuration parsing from workspace YAML
   - Port mapping from mochi config files
   - Environment variable injection tracking

### **Port Display System** (`pkg/coordination/ports.go`)
- **PortRegistry**: Central port tracking and metadata storage
- **PortScanner**: Active port discovery on remote nodes
- **PortDisplay**: User-friendly formatting and presentation
- Integration with node agents for remote execution

### **User Interface Components**
1. **CLI Display**: `mochi workspace ports <name>`
   - Table format with service name, internal port, external port
   - Protocol detection (HTTP, HTTPS, TCP, UDP)
   - Direct access URLs for web services

2. **Web Dashboard**: Coordination server port view
   - Real-time port status updates
   - Service health indicators
   - Direct link generation for web services

### **Service Metadata**
- Service name and description
- Port protocol and purpose (web, api, db, etc.)
- Health check endpoints
- Access URLs with authentication
- Container/VM association

### **Integration Points**
- Workspace creation: Port allocation and registration
- Service startup: Automatic port detection
- Node agent reporting: Port status updates
- User authentication: Access control for port information

## 🧪 Proof of Work
- [ ] Port discovery from containers and processes
- [ ] CLI port display commands
- [ ] Service metadata tracking
- [ ] Real-time port status updates
- [ ] Integration with workspace lifecycle
- [ ] User-friendly URL generation
