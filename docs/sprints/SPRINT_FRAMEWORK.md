# Sprint-Based Planning Framework

## Overview

This framework describes how the project organizes work using **timeboxed sprints** instead of long-running milestones. Sprints provide regular cadence, predictable delivery, and frequent feedback opportunities.

**Key Benefits:**
- 📅 **Regular Cadence**: Consistent sprint rhythm (1-3 weeks)
- 🎯 **Clear Goals**: Each sprint has specific, achievable objectives
- 📊 **Visible Progress**: Daily standups and weekly metrics
- 🔄 **Frequent Feedback**: Retrospectives every sprint
- 📈 **Continuous Improvement**: Lessons applied immediately

---

## Sprint Lifecycle

Every sprint follows a consistent lifecycle: Planning → Execution → Review → Retrospective → Archival

### Phase 1: Planning (Day 1)
1. Define sprint goal(s) aligned with milestone
2. Select tasks from backlog based on capacity and priority
3. Break tasks into daily/weekly deliverables
4. Identify risks and mitigation strategies
5. Set success criteria and metrics targets

### Phase 2: Execution (Days 2-N)
1. Daily standups (15 min) - synchronous or async
2. Task progress tracking and status updates
3. Blockers identified and resolved quickly
4. Risk log maintained in real-time
5. Scope adjustments made transparently

### Phase 3: Review (Final Day)
1. Demo completed work to stakeholders
2. Validate against acceptance criteria
3. Capture metrics (completion %, quality, velocity)
4. Document learnings and adjustments
5. Prepare for retrospective

### Phase 4: Retrospective (Final Day + 1)
1. What went well? (celebrate wins)
2. What could improve? (honest reflection)
3. What will we do differently? (action items)
4. Document decisions for next sprint
5. Plan improvements implementation

### Phase 5: Archival (After Sprint Complete)
1. Move sprint document to `completed/` folder
2. Archive sprint summary in project history
3. Update backlog based on sprint results
4. Extract technical debt discoveries
5. Prepare for next sprint

---

## Task Type Taxonomy

All tasks are categorized by type to clarify dependencies and ownership:

| Code | Name | Description |
|------|------|-------------|
| **INF** | Infrastructure | Docker, LXC, Worktrees, networking, provider engines |
| **COR** | Core/Control | Orchestration logic, config parsing, lifecycle management |
| **AGT** | Agent Integration | SSH, Agent Scaffold sync, AI agent configs |
| **CLI** | CLI/UX | Command structure, output formatting, scaffolding |
| **USR** | User Management | User registry, SSH keys, deep links, port discovery |
| **VFY** | Verification | E2E testing, coverage, CI integration |
| **CFG** | Configuration | Config parsing, validation, template management |

---

## Daily Standups

Use this structure for synchronous or async standups:

```
[Date: YYYY-MM-DD]

@person-name
🟢 **Yesterday**: [What did you accomplish?]
⚡ **Today**: [What will you work on?]
🚧 **Blockers**: [Anything blocking progress?]
💯 **Confidence**: 8/10 [Scale 1-10]
```

---

## Sprint Metrics

Track these each sprint:

| Metric | Target | Purpose |
|--------|--------|---------|
| Task Completion Rate | ≥80% | Sprint execution health |
| Test Coverage | 85%+ | Code quality |
| Build Time | <5 min | Development velocity |
| Code Review Cycle | <1 day | Feedback speed |

---

## Risk Management

Maintain active risk log during sprint:

```
## Active Risks

| # | Risk | Mitigation | Owner | Status |
|---|------|-----------|-------|--------|
| R1 | Complexity high | Start with MVP | @dev | 🟡 Active |
```

---

## Sprint Retrospective

Conduct retrospective at sprint end (1-2 hours):

```
# Sprint [N] Retrospective

## ✅ What Went Well?
[Celebrate wins]

## 📈 What Could Improve?
[Honest reflection]

## 🎯 Action Items
[Concrete changes for next sprint]

## 📊 Metrics Summary
- Task Completion: [X%]
- Quality: [Defects, coverage]
```

---

## Framework Created: January 17, 2026
