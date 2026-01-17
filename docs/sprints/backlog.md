# Project Backlog

**Last Updated**: January 17, 2026  
**Total Items**: [Backlog size]  
**Next Refinement**: [Date]

---

## Future Sprints (Sprint 7-12)

Work planned beyond M3 completion:

### Sprint 7: M4 Foundation - Multi-Machine Orchestration
**Estimated Duration**: 3 weeks  
**Success Criteria**: 
- [ ] Cross-machine workspace coordination working
- [ ] Provider synchronization between machines
- [ ] Distributed service discovery operational

### Sprint 8: Advanced Networking
**Estimated Duration**: 2 weeks  
**Success Criteria**:
- [ ] VPN/overlay networking for workspaces
- [ ] Multi-provider networking integration
- [ ] Performance optimization

### Sprint 9-12: Future Features
[Placeholder for future planning]

---

## Technical Debt

### High Priority

| ID | Issue | Impact | Est. Effort | Notes |
|----|----|--------|-------------|-------|
| **Debt-01** | QEMU provider mixed responsibilities | Architectural violations | 3-5 days | `execRemote()` should be transport layer |
| **Debt-02** | Service orchestration incomplete | Missing dependencies, health checks | 2-3 days | Need full dependency resolution |
| **Debt-03** | Error messages inconsistent | Poor user experience | 1-2 days | Standardize error format across CLI |

### Medium Priority

| ID | Issue | Impact | Est. Effort | Notes |
|----|----|--------|-------------|-------|
| **Debt-04** | Integration test coverage gaps | Risk of regressions | 3-4 days | Need multi-provider E2E tests |
| **Debt-05** | Configuration schema incomplete | Future extensibility blocked | 2-3 days | Need formal OpenAPI/JSON Schema |
| **Debt-06** | Code duplication in controller | Maintenance burden | 2 days | Extract common patterns |

### Low Priority

| ID | Issue | Impact | Est. Effort | Notes |
|----|----|--------|-------------|-------|
| **Debt-07** | LXC provider documentation incomplete | Developer onboarding friction | 1 day | Add LXC setup guide |
| **Debt-08** | Makefile incomplete | CI/CD friction | 1 day | Add all build targets |

---

## Ideas to Explore

### Research Spikes

| ID | Idea | Effort | Research Question |
|----|------|--------|-------------------|
| **Idea-01** | WebRTC for workspace access | 1 week | Can we support browser-based terminal? |
| **Idea-02** | Kubernetes integration | 2 weeks | How to coordinate K8s workspaces? |
| **Idea-03** | GPU support | 1 week | How to expose GPU to remote workspaces? |
| **Idea-04** | Persistent storage | 1 week | How to manage workspace data between sprints? |
| **Idea-05** | Team collaboration features | 2 weeks | How to enable workspace sharing? |

---

## Unscheduled Tasks from Existing Files

### From docs/planning/tasks/

#### User Management (USR)
- [ ] **USR-01**: User struct and registry implementation
- [ ] **USR-02**: User registration API endpoint
- [ ] **USR-03**: SSH key auto-generation CLI
- [ ] **USR-04**: Deep link generation for editors
- [ ] **USR-05**: Service port discovery and display

#### Configuration (CFG)
- [ ] **CFG-01**: Config Pull/Sync Commands (In Progress)
- [ ] **CFG-02**: Config Extraction to Plugins (Pending)

#### Verification (VFY)
- [ ] **VFY-03-09**: Integration and E2E tests (New for M3.1+)

---

## Sprint Planning Notes

### Dependencies to Watch
1. Coordination server (critical path) - blocks all remote functionality
2. Transport layer - blocks provider dispatch
3. Node agents - blocks provider execution on remote nodes

### Capacity Planning
- **Developer**: Typically 75-80% of week on sprint tasks (20-25% on interrupts/debt)
- **Buffer**: 10-15% per sprint for unknowns and learning

### Risk Considerations
- High complexity areas (coordination server) → allocate extra time
- First time in area → allocate mentoring/learning time
- External dependencies → identify and plan early mitigation

---

## Backlog Prioritization Framework

### Prioritization Factors
1. **Milestone Criticality**: Is this on critical path for current milestone?
2. **Dependency Impact**: Does this unblock other work?
3. **User Value**: How much value to users?
4. **Technical Health**: Does this reduce technical debt?
5. **Learning**: Is this a learning opportunity for team?

### Prioritization Formula
**Priority Score** = (Criticality × 40%) + (Impact × 30%) + (Value × 20%) + (Effort × 10%)

---

## Backlog Maintenance

### Weekly Refinement
- Review new items added to backlog
- Reassess priorities based on learnings
- Estimate unestimated items
- Identify dependencies

### Sprint Planning
- Select top-priority items
- Break into daily/weekly deliverables
- Estimate based on team capacity
- Identify risks

### Post-Sprint
- Update backlog with completed items
- Add discovered technical debt
- Reprioritize based on sprint learnings
- Extract action items from retrospectives

---

## Reference Documents

- **Sprint Framework**: `docs/sprints/SPRINT_FRAMEWORK.md`
- **Migration Guide**: `docs/sprints/MIGRATION.md`
- **Sprint Template**: `docs/sprints/sprint-template.md`
- **M3 Specification**: `docs/specs/m3.md`
- **M3 Planning**: `docs/planning/M3/`

---

**Backlog Created**: January 17, 2026
