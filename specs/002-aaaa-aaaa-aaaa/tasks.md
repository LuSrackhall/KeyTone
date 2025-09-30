# Tasks: 键音专辑签名系统

**Input**: Design documents from `/specs/002-aaaa-aaaa-aaaa/`
**Prerequisites**: plan.md (required), research.md, data-model.md, contracts/

## Phase 3.1: Setup
- [ ] T001 Ensure SDK project ready (Go 1.21+) and frontend dev env (Quasar) available
- [ ] T002 [P] Add frontend service bridge `frontend/src/services/sdk-bridge.ts` for signature APIs (reuse existing IPC pattern)
- [ ] T003 [P] Scaffold `SignatureDialog.vue` and route entry on Album page (button next to Delete)
- [ ] T004 Configure Playwright E2E test runner under `frontend/tests/e2e`

## Phase 3.2: Tests First (TDD)
- [ ] T005 [P] Contract test: GET sdk://signatures → `tests/contract/signatures_list.test.ts`
- [ ] T006 [P] Contract test: POST sdk://signatures → `tests/contract/signatures_create.test.ts`
- [ ] T007 [P] Contract test: DELETE sdk://signatures/{id} → `tests/contract/signatures_delete.test.ts`
- [ ] T008 [P] Contract test: POST sdk://signatures/{id}/export → `tests/contract/signatures_export.test.ts`
- [ ] T009 [P] Contract test: POST sdk://signatures/import → `tests/contract/signatures_import.test.ts`
- [ ] T010 [P] Contract test: GET sdk://albums/{albumId}/signatures → `tests/contract/albums_signatures_list.test.ts`
- [ ] T011 [P] Contract test: POST sdk://albums/{albumId}/sign → `tests/contract/albums_sign.test.ts`
- [ ] T012 [P] Contract test: POST sdk://albums/{albumId}/export → `tests/contract/albums_export.test.ts`
- [ ] T013 [P] E2E: Signature dialog create/export/import flow → `frontend/tests/e2e/signature_dialog.spec.ts`
- [ ] T014 [P] E2E: Album export with signature flow → `frontend/tests/e2e/album_sign_export.spec.ts`

## Phase 3.3: Core Implementation
- [ ] T015 [P] SDK: Signature entities and repository (load/save from global config; card image in sibling resources folder)
- [ ] T016 [P] SDK: Implement list/create/delete signature APIs
- [ ] T017 [P] SDK: Implement export (.ktsign) and import (no protect code input) APIs (with overwrite)
- [ ] T018 SDK: Album signing record (allowReexport record only in phase-1) and export with signatures
- [ ] T019 Frontend: Implement `SignatureDialog.vue` with fields (name, intro, card image)
- [ ] T020 Frontend: Add button next to Delete on Album page; open dialog; integrate list/create/delete/import/export via bridge
- [ ] T021 Frontend: Export flow integrate signature selection dialog and record preference

## Phase 3.4: Integration
- [ ] T022 [P] SDK: Structured logging and error mapping for signature operations
- [ ] T023 [P] Frontend: i18n strings for dialog and errors; accessibility labels
- [ ] T024 [P] Ensure no direct FS access in frontend except file chooser/save for .ktsign
- [ ] T025 [P] Cross-platform path handling in SDK; resources folder handling

## Phase 3.5: Polish & Quality Gates
- [ ] T026 [P] Unit tests coverage check (Go ≥85%, Frontend ≥80%)
- [ ] T027 [P] Performance: verify no UI regressions; E2E timings sane
- [ ] T028 [P] Update docs: quickstart.md confirmed; add screenshots if needed
- [ ] T029 [P] Code review against Constitution (architecture, UX, performance)

## Dependencies
- Setup (T001-T004) → Tests (T005-T014) → Implementation (T015-T021) → Integration (T022-T025) → Polish (T026-T029)
- SDK APIs (T016-T018) block Frontend integration (T019-T021)
- Contract tests (T005-T012) precede respective implementations

## Parallel Example
```bash
# Launch contract tests in parallel once scaffolds exist
# (different files, independent):
run test tests/contract/signatures_list.test.ts
run test tests/contract/signatures_create.test.ts
run test tests/contract/signatures_delete.test.ts
run test tests/contract/signatures_export.test.ts
run test tests/contract/signatures_import.test.ts
```
