## 1. OpenSpec alignment

- [x] 1.1 Add proposal, design, tasks, and capability spec for `admin-ai-config-domain`.
- [x] 1.2 Keep the phase boundary explicit: phase 1 only realigns AI config API ownership and compatibility storage; provider/model table split remains out of scope.

## 2. Backend AI config domain realignment

- [x] 2.1 Add backend tests for AI-config-specific read/write service behavior using the existing `ai_config` storage shape.
- [x] 2.2 Add AI module endpoints for loading and saving AI config without routing the page through generic config-management APIs.
- [x] 2.3 Keep the first-phase persistence adapter compatible with existing `sys_config.key = ai_config` storage.
- [x] 2.4 Extend startup repair / built-in API metadata so AI config read-write endpoints are granted to admin/system_admin without overwriting customized built-in metadata.

## 3. Frontend AI config page switch-over

- [x] 3.1 Add or update frontend API wrappers so the AI config page calls AI module endpoints for load/save.
- [x] 3.2 Switch the existing AI config page to the new AI endpoints without regressing current drawer-based provider/model UX.
- [x] 3.3 Re-check role permission grouping so AI config related APIs appear under AI ownership rather than relying on config-management grouping for this workflow.

## 4. Verification

- [ ] 4.1 Run `cd server && go test ./...`.
- [x] 4.2 Run targeted frontend verification plus `cd web && npm run build`.
- [x] 4.3 Run `cd web && npm run typecheck`, or explicitly report the known existing `vue-tsc` environment issue if it still blocks verification.
