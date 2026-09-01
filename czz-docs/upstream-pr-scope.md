# 上游 PR 范围记录

本文件只是记录，不随 PR 提交。**czz-docs/ 整个目录永远不进上游 PR**：上游 `Syngnat/GoNavi` 没有这个目录，这里放的是本机工作文档。

记录时间 2026-09-01。基线为 `upstream/dev` 的 `89f9ad71`，本地分支 `czz-dev`，合并基 `9ef337cd`，合并提交 `44e5e3e3`。

## 1. 总量与三种口径

| 口径 | 文件 | 行数 | 说明 |
| --- | --- | --- | --- |
| 相对 `upstream/dev` 全量 | 86 | +11486 / −1439 | 含本机文档与 SQL 例程安全 |
| 去掉 `czz-docs/` | 82 | +9379 / −1439 | 本机文档 2107 行不进 PR |
| 再去掉 SQL 例程安全 | 81 | +9375 / −1439 | 只剩 AI 供应商主题 |

## 2. 想提的主题：AI 供应商功能

用户口径的三块，对应到实际文件：

- **页面样式改造** — `AISettingsProvidersSection.{tsx,css}`、`AISettingsModal.tsx`、`AISettingsSidebar.tsx`、`useAIProviderLayout.ts`、`AIProviderModelSelect.tsx`、`common/ResizableDraggableModal.{tsx,css}`、`common/tooltipTiming.ts`、`App.{tsx,css}`。
- **更多 CLI 与地域候选** — `aiSettingsModalConfig.tsx`、`utils/aiProviderPresets.ts`、`utils/aiProviderEndpoints.ts`、`internal/ai/provider/{cursor_cli,grok_cli,claude_cli,codex_cli,codebuddy_cli}.go`、`shared/i18n/*.json` 六语言。
- **NVM 条件下的识别与查询** — `internal/ai/provider/{cli_lookup,cli_capabilities,cli_model_catalog}.go` 与对应测试、`internal/ai/service/local_cli_detection_cache_test.go`、`codebuddy_nvm_test.go`。让应用自身完成 CLI 发现与模型目录查询，不再依赖额外脚本。

## 3. 明确排除

- `czz-docs/`（4 文件 +2107）——本机研究稿、评估记录与任务卡。永久排除。
- SQL 例程安全（9 文件 +590 /−14）——`internal/ai/safety/{classifier,guard}.go`、`internal/app/{sql_inspect,sql_sanitize,headless_safety}.go`、`internal/mcpserver/service.go` 及各自测试。

### SQL 例程安全是什么

与供应商 UI 无关，是一条独立的安全收紧。它新增了方言无关的 `IsRoutineSQL`，把两类语句判定为 `SQLOpRoutine`：

- **例程调用**：`CALL` / `EXEC` / `EXECUTE`，以及 SQL Server 的裸过程调用。
- **例程部署**：对 `PROCEDURE` / `FUNCTION` / `TRIGGER` / `ROUTINE` 做 `CREATE` / `ALTER` / `DROP`。

该类型在**任何权限级别下都不放行**，`PermissionFull` 也不放宽，并在三个执行面同步：MCP `execute_sql`、Headless CLI、AI 对话守卫。适合作为单独的安全 PR，与供应商功能分开评审。

## 4. 拆分时的已知障碍

`internal/ai/types.go` 同时承载两个主题：供应商相关类型与 `SQLOpRoutine` 枚举常量。因此「只提 AI 供应商」不是纯粹的文件级排除——

- 保留整份 `types.go` 时，`SQLOpRoutine` 常量会以未被引用的形式进入 PR（可编译，但语义上属于另一主题）。
- 真正干净的拆分需要对 `types.go` 做 hunk 级取舍，并在拆分后单独跑一次 `go build ./internal/... ./cmd/...` 确认无悬空引用。

本记录**未**执行该拆分，也未建立任何 PR 分支。

## 5. 落库前的核验（2026-09-01）

| 项目 | 结果 |
| --- | --- |
| `npx tsc --noEmit` | 通过 |
| `npx vitest run` | 5025/5027；两条失败为 `main.browserMock` 既有与上游自带的 `testPolicy` 基线漏登 |
| `go build ./internal/... ./cmd/...` | 通过 |
| `go vet ./internal/ai/... ./internal/app/... ./internal/mcpserver/` | 通过 |
| `go test ./internal/ai/safety/ ./internal/mcpserver/` | 通过 |
| `go test ./internal/ai/provider/ ./internal/ai/service/` | 通过 |
| 密钥/个人路径扫描 | 未检出令牌、私钥或 `/Users/...` 绝对路径 |
| 新增 `TODO` / `FIXME` / `console.log` | 0 处 |

## 6. 需要推送时的做法

不改写本地 `czz-dev` 的 16 个提交，另起一条压平分支：

```
git switch -c feat/ai-provider-management upstream/dev
git merge --squash czz-dev
git restore --staged --worktree czz-docs <SQL 例程安全的 9 个文件>
# 按 §4 处理 internal/ai/types.go，然后 go build 复核
git commit
git push -u origin feat/ai-provider-management
```

PR 的 base 是 `Syngnat/GoNavi` 的 `dev`（上游默认分支），head 是 `CzzRef/GoNavi` 的该分支。截至本记录，用户明确要求**暂不推送**。
