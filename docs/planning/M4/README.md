# M4: Staging Environment & Production User Flow

**Status**: Specification Complete, Ready for Implementation  
**Duration**: 4-6 weeks (Late February 2026)  
**Priority**: Critical path to general availability

---

## Overview

M4 delivers a production-ready staging environment with a complete user experience. Building on M3 (Provider-Agnostic Remote Nodes), M4 focuses on seamless developer onboarding and service accessibility.

### Objective

Single command to create a development-ready environment:

```bash
curl https://nexus.example.com/install.sh | bash -s -- \
  --repo my-org/my-project --server staging.example.com
```

Result: Workspace running, services started, editor open. ~3 minutes total.

---

## Documentation Structure

### Specifications (Read First)
These define what M4 will deliver:

| Document | Purpose | Audience | Time |
|----------|---------|----------|------|
| M4_OVERVIEW.md | Executive summary, vision, timeline | Stakeholders | 15 min |
| M4_USER_FLOW_SPECIFICATION.md | 7-step user journey, configuration | Product, Design | 45 min |
| M4_TECHNICAL_SPECIFICATION.md | Architecture, data models, APIs | Engineers | 60 min |
| M4_IMPLEMENTATION_ROADMAP.md | 4-phase plan, 60+ tasks, timeline | Teams | 30 min |
| M4_QUICK_START_GUIDE.md | Quick reference | All | 10 min |
| api/M4_API_SPECIFICATION.md | REST API complete reference | Backend | 45 min |
| checklists/M4_IMPLEMENTATION_CHECKLIST.md | Implementation tasks | QA/Teams | Reference |

### Execution Planning (Use for Implementation)
These guide how to build M4:

| Document | Purpose | Audience |
|----------|---------|----------|
| **execution/README.md** | ← Start here for implementation | Everyone |
| execution/EXECUTION_PLAN_SUMMARY.txt | Overview of entire execution plan | Leads |
| execution/M4_EXECUTION_SUMMARY.md | Current state + blockers | Planners |
| execution/M4_STRATEGIC_PLAN.md | Detailed 4-phase plan | Engineers |
| execution/DELEGATION_GUIDE.md | Task templates for subagents | Leads |
| execution/BLOCKER_RESOLUTION_TASKS.md | Blocker fix specifications | Systems Engineer |

---

## Architecture

```
User CLI
  -> Coordination Server (port 3001)
    -> LXC Driver Node
      -> LXC Daemon
        -> Container (with SSH + services)
```

**Components**:
- Coordination Server: Central workspace lifecycle management
- Node Agent: Remote command executor on driver nodes
- Container: Isolated development environment with services

---

## Implementation Timeline

| Phase | Duration | Deliverable | Success Criteria |
|-------|----------|-------------|------------------|
| 1 | Weeks 1-2 | Coordination Server foundation | API working, node integration |
| 2 | Weeks 2-3 | GitHub CLI integration | Full auth and workspace creation |
| 3 | Week 3 | Install script | One-line bootstrap |
| 4 | Weeks 4-5 | Production ready | All success metrics met |

---

## Success Criteria

**User Experience**:
- First workspace < 3 minutes
- All steps automated
- Services discoverable
- Editor launches automatically
- SSH latency < 100ms

**Technical**:
- Server uptime 99.9%
- Container startup < 30 seconds
- 10+ concurrent workspaces
- Full isolation between workspaces

**Coverage**:
- Node.js, Python, static sites
- Projects with/without .nexus/config.yaml

---

## Quick Start by Role

### If You're Planning Implementation (START HERE ⭐)
1. **execution/README.md** - Overview of execution documents
2. **execution/EXECUTION_PLAN_SUMMARY.txt** - 10-minute summary
3. **execution/M4_EXECUTION_SUMMARY.md** - Blocker analysis & timeline
4. **execution/M4_STRATEGIC_PLAN.md** - Detailed 4-phase plan

### If You're Understanding Requirements
1. **INDEX.md** - Navigation guide
2. **M4_OVERVIEW.md** - Executive summary
3. **M4_TECHNICAL_SPECIFICATION.md** - Architecture
4. **M4_USER_FLOW_SPECIFICATION.md** - User journey

### If You're Building Implementation
1. **execution/DELEGATION_GUIDE.md** - Get your task
2. **M4_API_SPECIFICATION.md** - API contracts
3. **M4_TECHNICAL_SPECIFICATION.md** - Design details
4. **checklists/M4_IMPLEMENTATION_CHECKLIST.md** - Verify completion

---

## Related Documentation

- **M3**: docs/planning/M3/ (Provider-agnostic remote nodes)
- **Project**: README.md, AGENTS.md (project root)
- **Implementation**: Use checklists/ folder for task tracking

---

**Last Updated**: January 17, 2026  
**Next Step**: Begin Phase 1 (Coordination Server development)
