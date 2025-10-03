# Tasks: [FEATURE NAME]

**Input**: Design documents from `/specs/[###-feature-name]/`
**Prerequisites**: plan.md (required), research.md, data-model.md, contracts/

## Execution Flow (main)
```
1. Load plan.md from feature directory
   → If not found: ERROR "No implementation plan found"
   → Extract: tech stack, libraries, structure
2. Load optional design documents:
   → data-model.md: Extract entities → model tasks
   → contracts/: Each file → contract test task
   → research.md: Extract decisions → setup tasks
3. Generate tasks by category:
   → Setup: project init, dependencies, linting
   → Tests: contract tests, integration tests
   → Core: models, services, CLI commands
   → Integration: DB, middleware, logging
   → Polish: unit tests, performance, docs
4. Apply task rules:
   → Different files = mark [P] for parallel
   → Same file = sequential (no [P])
   → If changes impact existing logic: add regression tests BEFORE implementation; otherwise tests may follow implementation
5. Number tasks sequentially (T001, T002...)
6. Generate dependency graph
7. Create parallel execution examples
8. Validate task completeness:
   → All contracts have tests?
   → All entities have models?
   → All endpoints implemented?
9. Return: SUCCESS (tasks ready for execution)
```

## Format: `[ID] [P?] Description`
- **[P]**: Can run in parallel (different files, no dependencies)
- Include exact file paths in descriptions

## Path Conventions
- **Single project**: `src/`, `tests/` at repository root
- **Web app**: `backend/src/`, `frontend/src/`
- **Mobile**: `api/src/`, `ios/src/` or `android/src/`
- Paths shown below assume single project - adjust based on plan.md structure

## Phase 3.1: Setup
- [ ] T001 Create project structure per implementation plan
- [ ] T002 Initialize [language] project with [framework] dependencies
- [ ] T003 [P] Configure linting and formatting tools

 
## Phase 3.2: Impact Analysis & Regression Safety (before risky changes)

**CRITICAL: If changes impact existing stable logic, add/adjust regression tests BEFORE implementation.**

**Constitution Compliance: 测试策略—重构驱动与回归保护**
- [ ] T004 [P] Identify affected behaviors and expected user-observable outcomes (list in tasks.md)
- [ ] T005 [P] Add/adjust regression tests for affected areas in tests/integration and tests/unit
- [ ] T006 [P] Contract tests for new/changed endpoints in tests/contract/* (as applicable)
- [ ] T007 [P] Performance smoke for critical paths: audio ops <20ms, keyboard <10ms

 
## Phase 3.3: Core Implementation (ONLY after tests are failing)

**Constitution Compliance: 代码质量与架构分离 - 清晰的层次边界**
- [ ] T010 [P] User model in src/models/user.py
- [ ] T011 [P] UserService CRUD in src/services/user_service.py  
- [ ] T012 [P] CLI --create-user in src/cli/user_commands.py
- [ ] T013 POST /api/users endpoint (backend only, no frontend system calls)
- [ ] T014 GET /api/users/{id} endpoint (backend only, no frontend system calls)
- [ ] T015 Input validation with user-friendly error messages (用户体验一致性)
- [ ] T016 Error handling with i18n support (用户体验一致性)
- [ ] T017 Cross-platform compatibility checks (跨平台兼容性)

 
## Phase 3.4: Integration

**Constitution Compliance: 性能与响应性要求 + 跨平台兼容性**
- [ ] T018 Connect UserService to DB with connection pooling
- [ ] T019 Auth middleware with platform-specific secure storage
- [ ] T020 Structured logging (Go slog or standard log package)
- [ ] T021 CORS and security headers
- [ ] T022 Resource usage monitoring (memory <100MB, CPU <5%)
- [ ] T023 Cross-platform file path handling

 
## Phase 3.5: Polish & Quality Gates

**Constitution Compliance: All principles verification**
- [ ] T024 [P] Unit tests for validation in tests/unit/test_validation.py
- [ ] T025 [P] Coverage health check (report and discuss hotspots; no global hard threshold)
- [ ] T026 [P] UI accessibility (ARIA) compliance tests
- [ ] T027 [P] I18n completeness verification (no missing translation keys)
- [ ] T028 Performance benchmark tests (<3s startup, resource limits)
- [ ] T029 Cross-platform integration tests (Windows/macOS/Linux)
- [ ] T030 [P] Update docs/api.md with architecture diagrams
- [ ] T031 Code review checklist validation
- [ ] T032 Remove duplication and architectural violations
- [ ] T033 Run manual-testing.md across all platforms

 
## Dependencies

- Tests (T004-T007) before implementation (T008-T014)
- T008 blocks T009, T015
- T016 blocks T018
- Implementation before polish (T019-T023)

 
## Parallel Example

```bash
Task: "Contract test POST /api/users in tests/contract/test_users_post.py"
Task: "Contract test GET /api/users/{id} in tests/contract/test_users_get.py"
Task: "Integration test registration in tests/integration/test_registration.py"
Task: "Integration test auth in tests/integration/test_auth.py"
```

 
## Notes

- [P] tasks = different files, no dependencies
- If regression tests were added for impacted logic, verify they fail before the fix
- Commit after each task
- Avoid: vague tasks, same file conflicts

 
 
## Task Generation Rules

*Applied during main() execution*

1. **From Contracts**:
   - Each contract file → contract test task [P]
   - Each endpoint → implementation task
   
2. **From Data Model**:
   - Each entity → model creation task [P]
   - Relationships → service layer tasks
   
3. **From User Stories**:
   - Each story → integration test [P]
   - Quickstart scenarios → validation tasks

4. **Ordering**:
   - Setup → Tests → Models → Services → Endpoints → Polish
   - Dependencies block parallel execution

 
 
## Validation Checklist

*GATE: Checked by main() before returning*

- [ ] All contracts have corresponding tests
- [ ] All entities have model tasks
- [ ] All tests come before implementation
- [ ] Parallel tasks truly independent
- [ ] Each task specifies exact file path
- [ ] No task modifies same file as another [P] task