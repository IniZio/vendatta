# Task: USR-01 User struct and registry implementation

**Priority**: 🔥 High
**Status**: [Pending]

## 🎯 Objective
Implement core user management data structures and registry system to support user registration, authentication, and workspace assignment in the coordination server.

## 🛠 Implementation Details

### **User Data Structure**
1. **User Model** (`pkg/user/user.go`):
   - User struct with ID, username, email, SSH public key
   - Workspace assignments and permissions
   - Creation/update timestamps
   - JSON serialization tags for API responses

2. **Registry Interface** (`pkg/user/registry.go`):
   - User registration and lookup methods
   - SSH key validation and storage
   - Workspace assignment management
   - Thread-safe operations

### **Storage Backends**
- **Memory Registry**: In-memory implementation for development
- **File Registry**: JSON file-based persistence
- **Database Registry**: PostgreSQL/MySQL support (future)

### **SSH Key Management**
- Public key validation (format checking)
- Key fingerprint generation for identification
- Duplicate key detection
- Key rotation support

### **Integration Points**
- Coordination server user management
- SSH proxy authentication
- Workspace creation with user context
- Node agent user provisioning

## 🧪 Proof of Work
- [ ] User struct with all required fields
- [ ] Registry interface implementation
- [ ] SSH key validation logic
- [ ] Memory and file-based storage backends
- [ ] Unit tests (90%+ coverage)
- [ ] Integration with coordination server
