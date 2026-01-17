# Documentation Structure

This directory contains the complete documentation for the Nexus project.

## 📚 Quick Navigation

| Section | Purpose | Location |
|---------|---------|----------|
| **🎯 Specifications** | System design & requirements | `specs/` |
| **📅 Sprint Planning** | Active development cycles | `sprints/` |
| **🏗️ Milestones** | M4 (current) & reference docs | `planning/` |
| **📚 Guides** | Architecture & processes | `guides/` |

---

## Overview

```
docs/
├── specs/                       # Technical specifications (complete design)
│   ├── m3.md                   # M3: Provider-Agnostic Remote Nodes
│   ├── product/                # Product specifications
│   │   ├── configuration.md    # Configuration reference
│   │   ├── overview.md         # Product vision & goals
│   │   └── user_stories.md     # User requirements
│   ├── technical/              # Technical architecture
│   │   ├── agent-gateway.md    # Agent config generation
│   │   ├── architecture.md     # System design
│   │   ├── lifecycle.md        # Workspace lifecycle
│   │   ├── plugins.md          # Plugin system
│   │   └── roadmap.md          # Post-M3 vision
│   ├── testing/                # Quality assurance
│   │   ├── cases.md            # Test plan & cases
│   │   └── strategy.md         # Testing approach
│   ├── ux/                     # User experience
│   │   └── cli-ux.md           # CLI design
│   └── security.md             # Security specifications
│
├── sprints/                     # Sprint-based planning (ACTIVE)
│   ├── SPRINT_FRAMEWORK.md      # Framework & methodology
│   ├── MIGRATION.md             # M1/M2/M3 → Sprint mapping
│   ├── sprint-template.md       # Template for all sprints
│   ├── backlog.md               # Unscheduled work & ideas
│   ├── active/                  # Current & upcoming sprints
│   │   └── sprint-01.md        # Sprint 1 details (to be created)
│   └── completed/               # Finished sprints (archive)
│
├── planning/                    # Milestone documentation
│   ├── README.md               # Planning overview
│   ├── M4/                     # CURRENT: Staging Env & User Flow
│   │   ├── README.md           # M4 Overview
│   │   ├── M4_OVERVIEW.md      # Executive summary
│   │   ├── M4_USER_FLOW_SPECIFICATION.md
│   │   ├── M4_TECHNICAL_SPECIFICATION.md
│   │   ├── M4_IMPLEMENTATION_ROADMAP.md
│   │   ├── M4_QUICK_START_GUIDE.md
│   │   ├── api/                # API specifications
│   │   ├── checklists/         # Implementation checklists
│   │   ├── guides/             # Technical guides
│   │   └── specs/              # Configuration specs
│   ├── M3/                     # Reference: Provider-Agnostic Nodes
│   │   ├── M3_IMPLEMENTATION_STATUS.md
│   │   ├── M3_ARCHITECTURAL_CORRECTION.md
│   │   └── ... (other M3 docs)
│   ├── past-sprints/           # Completed milestones (archive)
│   │   ├── M1_MVP.md          # M1 (completed)
│   │   └── M2_ALPHA.md        # M2 (completed)
│   ├── tasks/                  # Legacy task files
│   │   ├── CLI-01.md
│   │   └── ... (old task specs)
│   └── TECHNICAL_DEBT.md       # Known issues & improvements
│
├── guides/                     # Architecture & process guides
│   ├── CONSOLIDATION_SUMMARY.md
│   ├── DELEGATION_WORKFLOW.md
│   └── M3_SPRINT_1_DELEGATION_EXAMPLE.md
```

---

## 📖 How to Use This Documentation

### 🎯 For Implementation (Active Development)
1. **Start here**: [`sprints/SPRINT_FRAMEWORK.md`](sprints/SPRINT_FRAMEWORK.md) - Understand sprint methodology
2. **Check current work**: [`sprints/active/`](sprints/active/) - See what's in progress
3. **Review specifications**: [`specs/m3.md`](specs/m3.md) - Master spec for current milestone
4. **Check backlog**: [`sprints/backlog.md`](sprints/backlog.md) - Upcoming work

### 📚 For Design & Architecture
- **Product Vision**: [`specs/product/overview.md`](specs/product/overview.md)
- **System Architecture**: [`specs/technical/architecture.md`](specs/technical/architecture.md)
- **User Experience**: [`specs/ux/cli-ux.md`](specs/ux/cli-ux.md)
- **Configuration**: [`specs/product/configuration.md`](specs/product/configuration.md)

### 🧪 For Testing & Quality
- **Test Strategy**: [`specs/testing/strategy.md`](specs/testing/strategy.md)
- **Test Cases**: [`specs/testing/cases.md`](specs/testing/cases.md)
- **Technical Debt**: [`planning/TECHNICAL_DEBT.md`](planning/TECHNICAL_DEBT.md)

### 📋 For M4 (Current Milestone)
- **M4 Planning**: [`planning/M4/`](planning/M4/) - Complete specification & roadmap
- **M4 Overview**: [`planning/M4/M4_OVERVIEW.md`](planning/M4/M4_OVERVIEW.md) - Executive summary
- **M4 Implementation**: [`planning/M4/M4_IMPLEMENTATION_ROADMAP.md`](planning/M4/M4_IMPLEMENTATION_ROADMAP.md) - Detailed plan
- **M4 User Flow**: [`planning/M4/M4_USER_FLOW_SPECIFICATION.md`](planning/M4/M4_USER_FLOW_SPECIFICATION.md) - Complete user journey
- **M4 Technical**: [`planning/M4/M4_TECHNICAL_SPECIFICATION.md`](planning/M4/M4_TECHNICAL_SPECIFICATION.md) - Architecture & APIs
- **M4 API Docs**: [`planning/M4/api/M4_API_SPECIFICATION.md`](planning/M4/api/M4_API_SPECIFICATION.md) - REST API reference
- **M4 Checklist**: [`planning/M4/checklists/M4_IMPLEMENTATION_CHECKLIST.md`](planning/M4/checklists/M4_IMPLEMENTATION_CHECKLIST.md) - Task list

### 📋 For Reference & History
- **M3 Details**: [`planning/M3/`](planning/M3/) - Provider-agnostic remote nodes
- **Completed Milestones**: [`planning/past-sprints/`](planning/past-sprints/) - M1, M2 archives
- **Old Task Files**: [`planning/tasks/`](planning/tasks/) - Legacy task specs
- **Guides**: [`guides/`](guides/) - Architecture & process documentation
- **Sprint-to-Milestone Mapping**: [`sprints/MIGRATION.md`](sprints/MIGRATION.md)

---

## 📊 Document Organization

### Specifications (`specs/`)
**Purpose**: Single source of truth for all system design decisions  
**Status**: Active - updated as implementation progresses  
**Contents**: Complete technical specifications, product requirements, architecture, testing strategy

### Sprints (`sprints/`)
**Purpose**: Timeboxed execution with regular feedback cycles  
**Status**: Active - Sprint 1 planned for January 20, 2026  
**Contents**: Sprint methodology, active sprint documents, completed sprints archive, backlog

### Planning (`planning/`)
**Purpose**: Historical reference and milestone tracking  
**Status**: Reference - M1/M2 archived, M3 active with legacy planning docs  
**Contents**: Milestone specifications, implementation status, M3 component plans, legacy tasks

---

## 🔄 Document Flow

```
Specs (Design) → Sprints (Execution) → Code (Implementation)
    ↓                  ↓                      ↓
M3.md          Sprint 1-6 docs         Source code
              Active sprint docs        Tests
              Sprint retrospectives     Commits
                                        PR descriptions

Feedback Loop: Code Review → Retrospectives → Planning Adjustments
```

---

## 🎯 Current Development Focus

**Current Milestone**: M4 - Staging Environment & Production User Flow  
**Status**: Planning Complete, Ready for Implementation  
**Timeline**: 4-6 weeks to completion (Late February 2026)

**M4 Phases**:
- **Phase 1 (Weeks 1-2)**: Coordination Server Foundation
- **Phase 2 (Weeks 2-3)**: GitHub CLI Integration
- **Phase 3 (Week 3)**: One-Line Install Script
- **Phase 4 (Weeks 4-5)**: Polish & Launch

**Previous Milestone**: M3 - Provider-Agnostic Remote Nodes (20% → 100%)

For M4 planning, see [`planning/M4/`](planning/M4/) → Start with [`M4_OVERVIEW.md`](planning/M4/M4_OVERVIEW.md)  
For active sprints, see [`sprints/active/`](sprints/active/)

---

## ✅ Documentation Maintenance

**Keep Updated**:
- Sprint documents (daily during sprint)
- `sprints/backlog.md` (weekly refinement)
- `specs/m3.md` (as implementation reveals changes)

**Archive When**:
- Sprint completes → Move to `sprints/completed/`
- Milestone completes → Move to `planning/past-sprints/`
- Task superseded → Move to `planning/tasks/` for reference

**No Longer Update**:
- Old milestone docs (reference only)
- Legacy task files (reference only)
- Completed sprint documents (archive only)

---

**Last Updated**: January 17, 2026  
**Current Focus**: M4 Specification & Planning Complete  
**Next Steps**: Begin M4 Phase 1 Implementation (Coordination Server)
