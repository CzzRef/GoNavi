# GoNavi 命令与核验

Tool: tool-neutral
Date: 2026-09-03

## 后端

改 SQL 判定 / 安全策略后至少：

```bash
go test ./internal/ai/safety/ ./internal/mcpserver/
go test ./internal/app/ -run 'TestInspectSQL|TestHeadlessSafety|TestIsReadOnlySQLQuery|TestEnsureReadOnlyConnectionAllows'
go vet ./internal/ai/... ./internal/app/... ./internal/mcpserver/
```

只验证后端时：`go build ./internal/... ./cmd/...`（根目录全量 `go build ./...` 需要 `frontend/dist`）。

已知预存在：`internal/app` 全包跑时 `TestFetchReleaseByURLFallsBackToCacheOn403` 可能失败（`updateReleaseCache` 全局状态）。

## 前端

```bash
cd frontend && npx tsc --noEmit
cd frontend && npx vitest run
```

全量 vitest 以当时记录为准。2026-09-03 合入 `upstream/dev` 后本机为 5144/5144。

## 实机核验包

需要用户在 macOS 桌面点的改动，走项目 Skill `gonavi-verify-build-restart` 与 [核验通路](../knowledge/gonavi-verify-build-restart.md)。停旧起新只跑 `restart.sh`，保持单实例。
