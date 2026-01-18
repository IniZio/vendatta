# Task: USR-03 SSH key auto-generation CLI

**Priority**: 🔥 High
**Status**: [Pending]

## 🎯 Objective
Implement CLI commands for automatic SSH key generation and management to streamline user onboarding and eliminate manual key handling.

## 🛠 Implementation Details

### **CLI Commands** (`cmd/nexus/user.go`)
1. **nexus user keygen**: Generate SSH key pair
   - Options: --type (rsa, ed25519), --bits, --output-dir
   - Automatic naming: id_ed25519_nexus_{timestamp}
   - Secure permissions (600 for private key)

2. **nexus user register**: Register user with generated key
   - Reads public key from standard location
   - Calls user registration API
   - Stores user ID and server connection info

3. **nexus user list-keys**: List user's SSH keys
   - Shows key fingerprints and creation dates
   - Integration with coordination server

4. **nexus user key-rotate**: Rotate SSH keys
   - Generates new key pair
   - Updates registration with new public key
   - Secure cleanup of old private keys

### **Key Generation Logic**
- Cryptographically secure random generation
- Standard SSH key formats (OpenSSH compatible)
- Key fingerprint calculation (MD5, SHA256)
- Private key encryption options (future)

### **Integration with User Registry**
- Automatic registration after key generation
- Key fingerprint tracking for uniqueness
- Server-side key validation
- Error handling for duplicate keys

### **Security Considerations**
- Private key never leaves local machine
- Secure file permissions enforcement
- Key backup recommendations
- Clear warnings about private key security

## 🧪 Proof of Work
- [ ] SSH key generation with multiple algorithms
- [ ] CLI commands for key management
- [ ] Integration with user registration API
- [ ] Secure file handling and permissions
- [ ] Cross-platform compatibility (Linux/macOS)
- [ ] Comprehensive error handling
