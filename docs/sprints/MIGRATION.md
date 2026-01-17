# Migration from Milestones to Sprint-Based Planning

## Executive Summary

This document outlines the migration strategy from **milestone-based planning** (M1, M2, M3) to a modern **sprint-based planning** system. The migration preserves all existing work while enabling more frequent delivery cycles (every 2 weeks vs. every 2-3 months).

**Migration Goal**: Achieve regular delivery cadence while maintaining clarity on milestone-level objectives.

---

## Current State vs Target State

### Current State: Milestone-Based
```
M1 (3-4 months) → M2 (2-3 months) → M3 (6-8 weeks)
```

### Target State: Sprint-Based with Milestone Alignment
```
Sprint 1  Sprint 2  Sprint 3  Sprint 4  Sprint 5  Sprint 6
  2w        2w        2w        2w        2w        2w
  └─ M3.1 (60%) ─┘
                  └─ M3.2 (85%) ──┘
                                  └─ M3.3 (100%) ─┘
```

---

## M3 Roadmap to Sprint Mapping

### Phase 1: M3.1 Foundation (20% → 60%)

#### Sprint 1: Coordination Server Core
- **Tasks**: Coordination server architecture, node registration, SSH pooling
- **Deliverable**: Working coordination server managing nodes

#### Sprint 2: Node Agents & Provider Dispatch
- **Tasks**: Node agent implementation, provider dispatch interface, CLI commands
- **Deliverable**: QEMU working via coordination server

**Milestone Status After**: 60% Complete ✅

---

### Phase 2: M3.2 Provider Integration (60% → 85%)

#### Sprint 3: Docker Remote Support
- **Tasks**: Docker execution via agents, port mapping
- **Deliverable**: Docker containers on remote nodes

#### Sprint 4: LXC Remote Support
- **Tasks**: LXC execution via agents, networking setup
- **Deliverable**: LXC containers on remote nodes

**Milestone Status After**: 85% Complete ✅

---

### Phase 3: M3.3 Polish (85% → 100%)

#### Sprint 5: Service Orchestration & User Management
- **Tasks**: Service dependencies, user registry, SSH key generation
- **Deliverable**: Complete user onboarding workflow

#### Sprint 6: Production Polish & Testing
- **Tasks**: Deep links, port discovery, E2E testing, security audit
- **Deliverable**: Production-ready remote development

**Milestone Status After**: 100% Complete ✅

---

## Task Type Preservation

All task types (INF, COR, AGT, CLI, USR, VFY, CFG) are preserved and distributed across sprints by type:

| Type | Current | Sprint Distribution | Total |
|------|---------|-------------------|-------|
| **COR** | 4 complete | COR-05,06 (S1), COR-07 (S2), COR-08 (S5) | 4 new |
| **INF** | 1 complete | INF-02 (S1), INF-03 (S2), INF-04-09 (S3-4) | 9 new |
| **AGT** | 1 complete | AGT-03 (S1), AGT-04 (S2) | 2 new |
| **CLI** | 3 complete | CLI-04 (S2), CLI-05 (S6) | 2 new |
| **USR** | 0 complete | USR-01,02,03 (S5), USR-04,05 (S6) | 5 new |
| **VFY** | 2 complete | VFY-03-09 (S1-6) | 7 new |

---

## Migration Process (5 Key Steps)

### Step 1: Directory Restructuring
1. Rename `docs/spec/` → `docs/specs/`
2. Create `docs/sprints/` with subdirectories
3. Move M1/M2 to `docs/planning/past-sprints/`
4. Move M3 docs to `docs/planning/M3/`

### Step 2: Framework Documentation
1. Create SPRINT_FRAMEWORK.md
2. Create sprint-template.md
3. Create backlog.md
4. Create MIGRATION.md (this file)

### Step 3: Reference Updates
1. Update all `docs/spec/` → `docs/specs/` references
2. Update README files in each directory
3. Verify no broken links

### Step 4: Initial Sprint Planning
1. Create Sprint 1 from template
2. Populate with M3.1 Phase 1 tasks
3. Conduct team kickoff

### Step 5: Ongoing Refinement
1. Run Sprint 1 and gather feedback
2. Adjust process based on learnings
3. Apply improvements to Sprint 2+

---

## M1 and M2 as Historical Reference

### Location
- **M1**: `docs/planning/past-sprints/M1_MVP.md`
- **M2**: `docs/planning/past-sprints/M2_ALPHA.md`

### Usage
- Reference for similar features
- Design pattern examples
- Performance baselines
- Lessons learned

---

## Success Indicators

**Sprint 1 (Week 1-2)**:
- Daily standups established
- Sprint document maintained
- Progress visible

**Sprint 2-3 (Week 3-6)**:
- Sprint 1 completed with 60-80% task completion
- Team comfortable with rhythm
- Retrospectives driving improvements

**Sprint 4+ (Week 7+)**:
- Consistent 80%+ completion rate
- Team autonomously managing sprints
- Regular delivery cadence established

---

## Migration Timeline

**Week of January 17, 2026:**
- Day 1-2: Prepare documentation and directories
- Day 2-3: Reorganize files and update references
- Day 4: Sprint 1 kickoff
- Week 2+: Execute sprints with regular rhythm

---

## Rollback Plan

If sprint-based planning isn't working after 4 sprints:
1. Conduct retrospective on sprint process itself
2. Identify specific issues
3. Attempt fixes (adjust duration, timing, scope)
4. Try 2 more sprints with improvements
5. Escalate if still problematic

**Note**: Give sprints at least 3-4 iterations before abandoning.

---

## Conclusion

This migration enables:
- ✅ Regular delivery every 2 weeks
- ✅ Frequent feedback and course correction
- ✅ Sustainable pace without crunch cycles
- ✅ Better risk management
- ✅ Improved team morale through visible progress

**Migration Start**: January 17, 2026  
**Sprint 1 Kickoff**: January 20, 2026  
**M3 Target Completion**: March 2026
