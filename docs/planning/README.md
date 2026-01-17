# Project Planning

## 📅 Migrating to Sprint-Based Planning

As of January 17, 2026, the project is transitioning from **milestone-based planning** to **sprint-based planning** to enable more frequent delivery (every 2 weeks vs. every 2-3 months).

**For Active Sprint Work**: See [`docs/sprints/`](../sprints/)  
**For Planning Framework**: See [`docs/sprints/SPRINT_FRAMEWORK.md`](../sprints/SPRINT_FRAMEWORK.md)  
**For Migration Details**: See [`docs/sprints/MIGRATION.md`](../sprints/MIGRATION.md)

---

## 🏷 Task Type Taxonomy

All tasks are organized by type to clarify ownership and dependencies:

| Code | Name | Description |
| :--- | :--- | :--- |
| **INF** | Infrastructure | Docker, LXC, Worktrees, Networking, Providers |
| **COR** | Core/Control | Orchestration logic, config parsing, lifecycle, coordination |
| **AGT** | Agent Integration | SSH, Agent Scaffold sync, AI agent configs, transport |
| **CLI** | CLI/UX | Command structure, output formatting, scaffolding, messaging |
| **USR** | User Management | User registry, SSH keys, deep links, port discovery |
| **VFY** | Verification | E2E testing, coverage, CI integration, validation |
| **CFG** | Configuration | Config parsing, validation, template management |

---

## 📅 Active Development

**Current Milestone**: M3 - Provider-Agnostic Remote Nodes  
**Current Status**: 20% Complete 🔴  
**Critical Path**: Coordination Server + Node Agents + Transport Layer  
**Timeline**: 6-8 weeks to completion  
**Target Completion**: March 2026

### M3 Phase Breakdown
- **M3.1 (Sprints 1-2)**: Coordination server foundation → 60% complete
- **M3.2 (Sprints 3-4)**: All providers remote support → 85% complete
- **M3.3 (Sprints 5-6)**: Production polish & testing → 100% complete

For active sprint work, see [`docs/sprints/active/`](../sprints/active/)

---

## 📋 M3 Implementation Tracking

All M3 planning documents are now organized in [`./M3/`](./M3/):

- **[M3_IMPLEMENTATION_STATUS.md](./M3/M3_IMPLEMENTATION_STATUS.md)** - Component-by-component status  
- **[M3_ARCHITECTURAL_CORRECTION.md](./M3/M3_ARCHITECTURAL_CORRECTION.md)** - Architectural lessons learned
- **[M3_ROADMAP.md](./M3/M3_ROADMAP.md)** - Phased implementation plan  
- **[M3_COORDINATION_SERVER_PLAN.md](./M3/M3_COORDINATION_SERVER_PLAN.md)** - Critical component spec
- **[M3_VERIFICATION_FINDINGS.md](./M3/M3_VERIFICATION_FINDINGS.md)** - Architecture verification results
- **[M3_VERIFICATION_REPORT.md](./M3/M3_VERIFICATION_REPORT.md)** - Initial verification report
- **[M3_FINAL_VERIFICATION_SUMMARY.md](./M3/M3_FINAL_VERIFICATION_SUMMARY.md)** - Final verification summary
- **[M3_COORDINATION_IMPLEMENTATION.md](./M3/M3_COORDINATION_IMPLEMENTATION.md)** - Implementation notes
- **[M3_DOCS_TIDY_COMPLETE.md](./M3/M3_DOCS_TIDY_COMPLETE.md)** - Documentation completion status
- **[M3_NEXT_STEPS.md](./M3/M3_NEXT_STEPS.md)** - Prioritized next steps

---

## 🏆 Past Milestones

Completed milestones are archived for reference:

- **[M1: CLI MVP](./past-sprints/M1_MVP.md)** - ✅ COMPLETED  
  - Delivered: Working Docker+Worktree + Agent Integration
  
- **[M2: Alpha](./past-sprints/M2_ALPHA.md)** - ✅ COMPLETED  
  - Delivered: Namespaced Plugins, UV-style Locking, Remote Configs

---

## 📝 Legacy Task Files

Individual task specifications from the old planning system are archived in [`./tasks/`](./tasks/) for reference.

**For new tasks**: Add to sprint planning in [`docs/sprints/`](../sprints/) instead.

---

## 📖 Planning Resources

- **Sprint Framework**: [`docs/sprints/SPRINT_FRAMEWORK.md`](../sprints/SPRINT_FRAMEWORK.md) - Sprint methodology guide
- **Migration Guide**: [`docs/sprints/MIGRATION.md`](../sprints/MIGRATION.md) - How old tasks map to new sprints
- **Backlog**: [`docs/sprints/backlog.md`](../sprints/backlog.md) - Unscheduled work and future planning
- **Master Spec**: [`docs/specs/m3.md`](../specs/m3.md) - M3 complete specification

---

**Planning Structure Updated**: January 17, 2026
