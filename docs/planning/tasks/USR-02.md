# Task: USR-02 User registration API endpoint

**Priority**: 🔥 High
**Status**: [Pending]

## 🎯 Objective
Implement REST API endpoints for user registration, management, and SSH key operations in the coordination server to enable programmatic user onboarding.

## 🛠 Implementation Details

### **API Endpoints** (`pkg/coordination/api/user.go`)
1. **POST /api/v1/users**: User registration
   - Request: username, email, ssh_public_key
   - Response: user_id, status
   - Validation: email format, SSH key validity

2. **GET /api/v1/users**: List users (admin only)
   - Query parameters: limit, offset, search
   - Response: paginated user list

3. **GET /api/v1/users/{id}**: Get user details
   - Response: full user information including workspaces

4. **PUT /api/v1/users/{id}/ssh-key**: Update SSH key
   - Request: new_ssh_public_key
   - Response: updated key fingerprint

5. **DELETE /api/v1/users/{id}**: Remove user
   - Cascading cleanup of workspaces and assignments

### **Request/Response Models**
- UserRegistrationRequest struct
- UserResponse struct
- ErrorResponse struct for API errors
- Pagination metadata

### **Authentication & Authorization**
- API key authentication for admin operations
- User-specific endpoints with proper access control
- Rate limiting for registration endpoints
- Input sanitization and validation

### **Integration with Registry**
- Direct integration with user registry backend
- Atomic operations for consistency
- Error handling for registry failures
- Audit logging for user operations

## 🧪 Proof of Work
- [ ] Complete REST API endpoints
- [ ] Request/response validation
- [ ] Authentication middleware
- [ ] Integration tests with user registry
- [ ] API documentation (OpenAPI/Swagger)
- [ ] Rate limiting and security measures
