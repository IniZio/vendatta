# M4: Staging Environment & Production User Flow - Executive Overview

**Milestone**: M4  
**Status**: Planning Complete, Ready for Implementation  
**Duration**: 4-6 weeks (late February 2026)  
**Investment**: Critical Path for General Availability  
**Target Users**: Designers, Developers, Engineers  

---

##  What is M4?

After M3 (technical foundation: coordination server), M4 delivers the **user experience** that makes Nexus viable for production use. It's the bridge between infrastructure and adoption.

### The Problem We're Solving

**Current State**:
- Developers can use Nexus with significant manual setup
- SSH keys require manual configuration
- Workspace creation is technical and fragmented
- Designer onboarding is complex and error-prone

**Target State**:
- One command starts a complete development environment
- Automatic GitHub authentication & SSH key management
- Services discoverable & immediately accessible
- Editor opens automatically with working connection
- fully automated beyond the initial install

---

## 💡 Core Concept

**Enable this workflow**:

```
User flow:
1. "I want to design for project X"
2. Run: curl https://... | bash
3. Authenticate with GitHub (browser, ~30 seconds)
4. Wait (~2 minutes) for environment setup
5. Editor opens automatically
6. Start designing

Total time: ~3 minutes
Manual steps: 0
Setup knowledge required: None
```

---

##  What Gets Built

### 1. Coordination Server (Core)
**What it does**: 
- Manages workspace lifecycle on remote nodes
- Allocates SSH ports & forwards connections
- Discovers running services & their ports
- Registers users (GitHub handle + SSH key)

**Technical**: Single-binary Go application, HTTP API, port 3001

**Why it matters**: Central management for multi-workspace, multi-user scenarios

### 2. Node Agent
**What it does**:
- Runs on LXC/Docker driver nodes
- Receives commands from coordination server
- Executes container operations
- Reports status back

**Technical**: Agent binary that connects via SSH/HTTP

**Why it matters**: Enables remote execution without direct SSH access to containers

### 3. CLI Commands
**New commands**:
- `nexus auth github` - GitHub login & token storage
- `nexus ssh setup` - SSH key generation & GitHub upload
- `nexus workspace create <repo>` - Fork/clone & create
- `nexus workspace connect` - Editor launch with SSH
- `nexus workspace services` - Service discovery

**Why it matters**: Provides clear, automated interface for users

### 4. One-Line Install Script
**What it does**:
```bash
curl https://nexus.example.com/install.sh | bash -s -- \
  --repo my-org/my-project --server staging.example.com
```

Automatically:
- Downloads nexus CLI
- Checks dependencies
- Initializes configuration
- Starts setup workflow

**Why it matters**: Removes all barriers to first-time use

---

##  User Experience Flow

### 7-Step Journey (3 Minutes)

```
Step 1: Install (20 seconds)
├─ Download & install CLI
├─ Check dependencies (git, ssh, gh)
└─ Initialize config

Step 2: GitHub Auth (30 seconds)
├─ Detect GitHub CLI
├─ Orchestrate `gh auth login` if needed
└─ Get username & tokens

Step 3: SSH Setup (15 seconds)
├─ Check for existing keys
├─ Generate new key if needed
└─ Upload to GitHub

Step 4: Repository Setup (20 seconds)
├─ Check ownership (owned vs. org)
├─ Fork if needed via `gh repo fork`
└─ Clone to local disk

Step 5: Workspace Creation (request)
├─ Send to coordination server
└─ Container starts creating (async)

Step 6: Workspace Initialization (60 seconds, polling)
├─ Container creation + SSH setup
├─ Repository clone + dependency install
├─ Service startup (with dependency ordering)
└─ Health checks pass

Step 7: Editor Launch (10 seconds)
├─ Detect editor (Cursor > VS Code > Vim)
├─ Generate SSH deep link
├─ Launch editor with remote connection
└─ Display summary

TOTAL: ~3 minutes from install script to editing code
```

---

##  Success Criteria

### For Users
-  First workspace takes < 3 minutes
-  All steps automated after install
-  Clear error messages with next steps
-  Services immediately accessible & discoverable
-  Editor opens automatically

### For Team
-  Coordination server achieves 99.9% uptime
-  SSH latency < 100ms
-  Container startup < 30 seconds
-  Service health checks 100% reliable
-  Supports 10+ concurrent workspaces

### For Business
-  Removes setup friction for new contributors
-  Enables frictionless designer onboarding
-  Creates reproducible development environment
-  Positions Nexus as "the easy development platform"

---

## 💰 Investment & ROI

### Development Cost
- **Engineering**: 4-6 weeks, ~1 senior engineer
- **Infrastructure**: 1 staging host (can be modest)
- **Testing**: 1-2 weeks iteration with beta testers

### Return
- **User Acquisition**: Frictionless onboarding removes major barrier
- **Product Viability**: Enables launch of general availability
- **Competitive Advantage**: "One command to develop" messaging
- **Support Load**: Clear error messages reduce support tickets

---

## 🗓 Implementation Timeline

### Phase 1: Foundation (Weeks 1-2)
- Coordination Server API & basic operations
- Node Agent infrastructure
- LXC driver integration
- Metadata storage (SQLite)

### Phase 2: Integration (Weeks 2-3)
- GitHub CLI orchestration
- SSH key management
- Repository fork/clone automation
- Workspace creation flow

### Phase 3: Automation (Week 3)
- One-line install script
- Platform support (macOS, Linux)
- Dependency detection & installation

### Phase 4: Launch (Week 4-5)
- E2E testing & bug fixes
- UX refinement
- Documentation & guides
- Rollout plan

---

##  What "Done" Looks Like

### For a New User

```
$ curl https://nexus.example.com/install.sh | bash -s -- \
    --repo my-org/my-project --server staging.example.com

📦 Installing nexus...
   Downloaded & installed to ~/.local/bin

⚙  Checking dependencies...
   git found
   ssh found
   curl found
  📦 Installing GitHub CLI...
     Installed

🔐 GitHub Authentication
  [Opens browser for OAuth]
   Authenticated as: alice

🔑 SSH Key Setup
   Detected existing key
   Key already on GitHub
   Registered with coordination server

📦 Repository Setup
   Checking ownership...
   Already your fork
   Cloning to ~/my-project

🔨 Workspace Initialization (estimated 60s)
   Container created
   SSH configured
   Dependencies installing...
   Services starting...

 Workspace Ready!
  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  Project: my-project (my-org/my-project)
  Status: Running 
  
  Services:
    • web   → http://localhost:23000 (npm dev)
    • api   → http://localhost:23001 (npm server)
    • db    → localhost:23002 (postgresql)
  
  SSH: ssh -p 2222 dev@staging.example.com
  
  Opening Cursor...
  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

 ready to work! Edit files in ~/my-project/src
```

**At this point**: Editor is open, services running, ready to code.

---

##  Architecture Overview

```
┌─────────────────────────────────────────┐
│     Staging Host (Ubuntu 22.04)         │
│                                         │
│  ┌──────────────────────────────────┐  │
│  │ Coordination Server              │  │
│  │ (Port 3001, REST API)            │  │
│  │                                  │  │
│  │ • Workspace CRUD                 │  │
│  │ • User registration              │  │
│  │ • SSH port forwarding            │  │
│  │ • Service discovery              │  │
│  └──────────────────────────────────┘  │
│                ↓                        │
│  ┌──────────────────────────────────┐  │
│  │ LXC Driver Node                  │  │
│  │ ┌──────────────────────────────┐ │  │
│  │ │ Node Agent                   │ │  │
│  │ │ • Cmd receiver & executor    │ │  │
│  │ │ • Container lifecycle        │ │  │
│  │ └──────────────────────────────┘ │  │
│  │ ┌──────────────────────────────┐ │  │
│  │ │ LXC Containers               │ │  │
│  │ │ ┌────────────────────────────┐│ │  │
│  │ │ │ workspace-abc123           ││ │  │
│  │ │ │ • SSH server               ││ │  │
│  │ │ │ • Services (web, api, db)  ││ │  │
│  │ │ │ • Repository               ││ │  │
│  │ │ └────────────────────────────┘│ │  │
│  │ └──────────────────────────────┘ │  │
│  └──────────────────────────────────┘  │
└─────────────────────────────────────────┘
         ↑ SSH:2222-2299 (forwarded)
         │
    Designer's Machine
    ├─ nexus CLI
    ├─ GitHub CLI
    ├─ Editor (Cursor/Code)
    └─ Terminal
```

---

## 🔐 Security & Privacy

### User Data
- **GitHub Authentication**: Uses official GitHub OAuth (no password stored)
- **SSH Keys**: Private keys stay on user's local machine only
- **Public Key**: Only fingerprint stored on coordination server

### Communication
- **SSH Forwarding**: Encrypted end-to-end
- **API Access**: HTTPS only, auth via GitHub token
- **Container Isolation**: Full LXC isolation per workspace

### Compliance
- No personal data collected beyond GitHub username
- SSH key management follows industry best practices
- Staging environment requires VPN/firewall access

---

##  Why This Matters

### For Designers
- **Instant Onboarding**: No setup knowledge required
- **Secure Access**: GitHub credentials, no password sharing
- **Isolated Environment**: Clean slate every workspace
- **Full IDE Features**: Remote SSH connection with IDE

### For Teams
- **Reproducible Environments**: Everyone gets same setup
- **Reduced Onboarding Time**: Hours → minutes
- **Reduced Support Load**: Clear error messages, self-service
- **Team Velocity**: Focus on work, not setup

### For the Company
- **Competitive Positioning**: "Develop anywhere, instantly"
- **Contributor Friction**: Minimal barrier to first contribution
- **Sustainability**: Platform requires minimal manual intervention
- **Scalability**: One coordination server manages many workspaces

---

##  Documentation Package

M4 includes comprehensive documentation for:

1. **Executive Stakeholders**: This overview + timeline
2. **Product Managers**: User flow specification + success metrics
3. **Developers**: Technical specification + API docs
4. **QA/Testing**: Checklist + test cases
5. **Operations**: Deployment guide + monitoring setup
6. **External Contributors**: Complete spec + contribution guide

All documents located in `/docs/planning/M4/`

---

##  Next Steps

### For Approval
1. Review this overview with stakeholders
2. Confirm timeline & resource allocation
3. Approve Phase 1 (Coordination Server) start

### For Implementation Team
1. Read full specifications (see README.md)
2. Review API contracts in detail
3. Prepare development environment
4. Begin Phase 1 sprint planning

### For Infrastructure
1. Provision staging host (Ubuntu 22.04)
2. Install LXC & configure network
3. Set up DNS for staging.example.com
4. Prepare for agent deployment

---

##  Key Differentiators

| Aspect | Traditional Setup | Nexus M4 |
|--------|-------------------|----------|
| Time to develop | 1 hour+ | 3 minutes |
| Manual steps | 10+ | 0 |
| SSH key setup | Manual copying | Automatic GitHub upload |
| Dependency install | Manual | Automatic |
| Service startup | Manual | Automatic with dependency ordering |
| Editor setup | SSH config files | Auto-launch with deep links |
| Cost to onboard | High (time) | Low (automated) |

---

##  Questions?

**Q: When does M4 start?**  
A: After M3 coordination server is merged and working.

**Q: Can I start working on M4 in parallel with M3?**  
A: Yes. Planning is complete. Can start Phase 1 once M3 is ~80% done.

**Q: What if LXC isn't available on the staging host?**  
A: Can use Docker as fallback. LXC is preferred for performance.

**Q: Is M4 staging-only or can it be deployed to production?**  
A: This specification is staging-focused. Production deployment is post-M4.

**Q: What happens to existing workspaces on staging when we deploy M4?**  
A: New system is additive. Existing workspaces can coexist or migrate.

---

##  Document Index

**Core Specifications**:
- `M4_OVERVIEW.md` (this file) - Executive summary
- `M4_USER_FLOW_SPECIFICATION.md` - Detailed 7-step journey
- `M4_TECHNICAL_SPECIFICATION.md` - Architecture & APIs
- `M4_IMPLEMENTATION_ROADMAP.md` - Phase-by-phase plan

**Supporting Documentation**:
- `api/` - OpenAPI specifications
- `checklists/` - Implementation & testing checklists
- `guides/` - Architecture guides & troubleshooting
- `specs/` - Configuration & protocol specs

---

**Document**: M4 Executive Overview  
**Version**: 1.0  
**Status**: Final, Ready for Review  
**Date**: January 17, 2026  
**Authors**: Planning Team  
**Distribution**: All stakeholders
