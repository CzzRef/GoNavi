<!-- codenote-local-context:conditional-v3 -->
# Project context

Project-owned conditional detail. Edit this local owner for project-specific facts; global policy stays in the compact core. Commands and inline paths are relative to the repository root unless their original text says otherwise. Read the sections relevant to the affected surface before material work.

## Project context from AGENTS.md

# GoNavi

本文件只做仓库入口，不复述 CodeNote 规则正文。共享内核由宿主全局指令加载；本仓库自己的路由在 [vibe/rules/README.md](<README.md>)。

- 过程枢纽：[vibe/specs/PROJECT_STATUS.md](<../specs/PROJECT_STATUS.md>)
- 任务索引：[vibe/specs/README.md](<../specs/README.md>)
- 上游 PR 历史：[vibe/pr/README.md](<../pr/README.md>)

`czz-docs/` 只留退役 README。当前权威只走 `vibe/`。

## Project context from CLAUDE.md

# GoNavi 项目入口

> **本文件是 link-only 的发现指针，不是规则正文。**
> 任何共享规则的正文都在 CodeNote 的 sole owner 里；本文件只负责路由 + 记录本仓库自己的事实。
> 若发现本文件复述了某条规则的算法、阈值或条款，那是漂移，应删掉并改成链接。

- 文件状态：`.gitignore:37`（`**/CLAUDE.md`）已忽略，本文件**只存在于本机**，不会提交、不会推送。
- 这是有意的：`origin` 是公开 fork `CzzRef/GoNavi`，`upstream` 是 `Syngnat/GoNavi`。带用户 Home 绝对路径的入口文件不能进版本库。

---

## 1. 每个新任务先加载

CodeNote 根：`/Users/gdkmjd/work/czz/CzzProj/CodeNote`

按 [Shared Always-Load Baseline](/Users/gdkmjd/work/czz/CzzProj/CodeNote/AiRef/VibePractice/Vibe_Rules/routing/README.md) 载入，每个任务一次：

| load_id | owner |
| --- | --- |
| routing-authority | [Vibe_Rules/routing/README.md](/Users/gdkmjd/work/czz/CzzProj/CodeNote/AiRef/VibePractice/Vibe_Rules/routing/README.md) |
| global-baseline | [Vibe_Rules/VibeAi.md](/Users/gdkmjd/work/czz/CzzProj/CodeNote/AiRef/VibePractice/Vibe_Rules/VibeAi.md) |
| project-entry | 本文件 |

固定 guard 通过 [Fixed Guards](/Users/gdkmjd/work/czz/CzzProj/CodeNote/AiRef/VibePractice/Vibe_Rules/routing/README.md) 解析，与领域路由无关地生效。本仓库**未**设置 `intent-note-gate`，该条件 guard 不激活。

会话命名遵循 [session-title.md](/Users/gdkmjd/work/czz/CzzProj/CodeNote/AiRef/VibePractice/Vibe_Rules/adapters/session-title.md)，格式与判定读该 owner，本文件不复述。

## 2. 本仓库是什么（决定路由的关键事实）

GoNavi 是**数据库客户端与 Agent 受控执行入口**，不是业务系统。它同时是：

- Wails/Go + React 桌面端
- Headless CLI
- MCP server（stdio + Streamable HTTP）
- Web server

这带来一个非常规的结论：**AI-DB 治理在本仓库的适用对象是「工具自身的行为」，不是某个业务库的数据。**

因此：

- 本仓库**没有** `vibe/ai-db/` 工作区，也不需要建——[AI-DB Governance §5](/Users/gdkmjd/work/czz/CzzProj/CodeNote/DevelopRef/调试工具/db/governance/README.md) 明确「不要为非 DB 项目建空的 ai-db 工作区」。
- 与 SQL 治理相关的工作在这里表现为**改判定代码与测试**，不是写 `sql.md` 候选。
- 当任务确实要对某个真实库出变更候选时，那属于该业务项目的 AI-DB 工作区，不属于 GoNavi。

## 3. SQL / AI 安全相关的 owner 路由

| 主题 | sole owner |
| --- | --- |
| 默认安全、触发条件、包布局、结果契约 | [AI-DB Governance](/Users/gdkmjd/work/czz/CzzProj/CodeNote/DevelopRef/调试工具/db/governance/README.md) |
| 交接文档结构与分阶段前置条件 | [SQL Handoff Prerequisites](/Users/gdkmjd/work/czz/CzzProj/CodeNote/DevelopRef/调试工具/db/governance/sql-handoff-prerequisites.md) |
| 规则/布局演进、`ADAPTIVE_REVIEW_REQUIRED` | [SQL Governance Adaptive Evolution](/Users/gdkmjd/work/czz/CzzProj/CodeNote/DevelopRef/调试工具/db/governance/sql-governance-adaptive-evolution.md) |
| 自动核验链 / 自动更新日志链（**候选提案，未升治理正文**） | [260828 提案](/Users/gdkmjd/work/czz/CzzProj/CodeNote/DevelopRef/调试工具/db/260828-sql-ai-verify-and-changelog-proposal.md) |
| GoNavi 统一执行层路线（**研究稿，非治理正文**） | [260826 研究稿](/Users/gdkmjd/work/czz/CzzProj/CodeNote/DevelopRef/调试工具/db/260826-db-gonavi-research.md) |

## 4. 本仓库自己的判定事实

以下是 GoNavi 独有的、无法从上面 owner 推导出来的事实。改动这些位置时必须同步本节。

### 4.1 单一判定来源

`appcore.InspectSQL` 是 MCP 与 CLI **共用**的语句判定来源，逐语句给出 `Keyword` / `ReadOnly`：

- [internal/app/sql_inspect.go](<../../internal/app/sql_inspect.go>) — 类型与 `InspectSQL`
- [internal/app/sql_sanitize.go](<../../internal/app/sql_sanitize.go>) — `isReadOnlySQLQuery`
- [internal/ai/safety/classifier.go](<../../internal/ai/safety/classifier.go>) — `ClassifySQL`

新增数据源或方言时，只读判定必须保持 **allowlist**（`default: return false`），不得改成 denylist。

### 4.2 三个执行面共用同一套判定

| 执行面 | 判定入口 |
| --- | --- |
| MCP `execute_sql` | [internal/mcpserver/service.go](<../../internal/mcpserver/service.go>) `ExecuteSQL` → `evaluateSQLSafety` → `isOperationAllowed` |
| Headless CLI | [internal/app/headless_safety.go](<../../internal/app/headless_safety.go>) `evaluateHeadlessSQLSafety` → `isHeadlessSQLOperationAllowed` |
| AI 对话面板 | [internal/ai/safety/guard.go](<../../internal/ai/safety/guard.go>) `Guard.Check` |

三处必须同进同退。只改其中一处就是漏洞——`AuthorizeMCPConnectionSQL` 是最后一道兜底，不是第一道。

### 4.3 例程执行走产品权限，不在执行层无条件禁止

2026-09-01 产品决定：CodeNote「Agent 不得直接执行 `CALL` / `EXEC` / 存储过程」约束的是 Agent **绕过产品去打库**。经 GoNavi MCP `execute_sql`、Headless CLI 或应用本身发出的语句，按连接既有 SQL 权限与写保护执行，**不再**用 `SQLOpRoutine` 在任何级别下无条件拒绝。

该无条件禁止已从 `czz-dev` 回滚，与 `upstream/dev` 对齐。不要在未重新授权的情况下加回三个执行面的例程闸门。

### 4.4 连接级保护与调用级断言是两层，不是二选一

- 连接级 `readOnly` 布尔 + 四类保护键（`restrictDataEdit` / `restrictStructureEdit` / `restrictScriptExecution` / `restrictDataImport`）：[internal/app/connection_readonly.go](<../../internal/app/connection_readonly.go>)
- CodeNote 的 `ExpectedDatabase` 是**调用级安全断言**，本仓库**尚未实现**（见 §5 G3）

两者互补：前者管「这条连接允不允许写」，后者管「这次调用打在不在预期的库上」。不要用其中一个替换另一个。

### 4.5 审计日志的性质边界

[internal/sqlaudit](<../../internal/sqlaudit>) 是 **GoNavi 本地**的 SHA-256 哈希链审计（`WeakValidation` 恒为真，自述只防误改不防攻击者）。

它**不是** CodeNote 要求的库内两级更新日志（`datafix_script_execution_log` + `datafix_script_operation_log`）。不要把它当成后者来汇报。

## 5. 与 CodeNote 的嵌套核验现状

完整核验见 [260828 提案 §8](/Users/gdkmjd/work/czz/CzzProj/CodeNote/DevelopRef/调试工具/db/260828-sql-ai-verify-and-changelog-proposal.md)。当前状态：

| # | 断裂点 | 状态 |
| --- | --- | --- |
| G1 | 仓库无项目入口 | **已闭合**（本文件，local-only；云端会话仍未覆盖，见 §6） |
| G2 | 存储过程在 Full 模式下可被 Agent 执行 | **按产品决定不在执行层闭合**（2026-09-01 回滚 `SQLOpRoutine`；MCP/应用按连接权限执行例程） |
| G3 | 无 `ExpectedDatabase` 身份闸门 | 未闭合 |
| G4 | DryRun 只覆盖 DataGrid 行改动 | 未闭合 |
| G5 | 审计日志无字段级 old→new | 未闭合 |
| G6 | 审计可被关闭 | 未闭合 |
| G7 | 无 `Reason` 强制 | 未闭合 |
| G8 | 需求追踪文档不是 handoff 契约 | 未闭合 |

`GRANT` / `REVOKE` 与 `DO $$`、`DECLARE`、`BEGIN…END` 匿名块仍归 `SQLOpOther`，完全模式下可执行。这与当前「经产品执行面按连接权限跑」的口径一致，不要在未授权时单独收紧。

## 6. 已知边界

- **云端会话未覆盖**：本文件是 gitignored 的，看不到用户 Home 的云端 worker 读不到它。routing 要求的「仓库指针」在公开 fork 上无法安全落地（会外泄 Home 路径）。需要云端覆盖时，应改为组织级 instruction surface，而不是往仓库里塞绝对路径。
- `docs/需求追踪/` 同样被 `.gitignore:35` 忽略，是本机文档，不进版本库。
- 根目录 `go build ./...` 需要 `frontend/dist`（`npm run build` 产物，已 gitignore）。只验证后端时用 `go build ./internal/... ./cmd/...`。
- `frontend/node_modules` 未安装时无法跑前端 typecheck 与测试。

## 7. 定向验证命令

改动 SQL 判定/安全策略后至少跑这三条：

```bash
go test ./internal/ai/safety/ ./internal/mcpserver/
```

```bash
go test ./internal/app/ -run 'TestInspectSQL|TestHeadlessSafety|TestIsReadOnlySQLQuery|TestEnsureReadOnlyConnectionAllows'
```

```bash
go vet ./internal/ai/... ./internal/app/... ./internal/mcpserver/
```

已知预存在失败（与 SQL 判定无关，勿误判为本次引入）：`internal/app` 全包运行时 `TestFetchReleaseByURLFallsBackToCacheOn403` 失败——5 个测试共享 `updateReleaseCache` 全局状态，单跑该测试通过。

## Project context from vibe/rules/README.md

本目录只维护本仓库特化路由。共享 primary / additive / evidence 由宿主已加载的 CodeNote Rule Kernel 选择；这里不复制算法。

本仓库**没有** `intent-note-gate`。

## Conditional task routes

1. 宿主已注入的 CodeNote 内核（新任务一次）
2. [documentation.md](<documentation.md>)（Standard / Controlled、文档治理、task card 位置）
3. [PROJECT_STATUS.md](<../specs/PROJECT_STATUS.md>)（进行中的供应商 UI / 上游 PR）
4. [project.md](<project.md>)（SQL 判定、三执行面、fork 与上游 PR 边界）
5. [workflow.md](<workflow.md>)（核验包、测试命令）

## Project Routes

- 过程枢纽：[PROJECT_STATUS.md](<../specs/PROJECT_STATUS.md>)
- 任务索引：[../specs/README.md](<../specs/README.md>)
- 可复用知识：[../knowledge/README.md](<../knowledge/README.md>)
- 上游 PR 历史：[../pr/README.md](<../pr/README.md>)
- 供应商界面约定：[../knowledge/ai-provider-ui-conventions.md](<../knowledge/ai-provider-ui-conventions.md>)
- 核验通路：[../knowledge/gonavi-verify-build-restart.md](<../knowledge/gonavi-verify-build-restart.md>)
- 项目 Skill：`gonavi-verify-build-restart`（CodeNote `Skills/projects/gonavi/`，本仓 `.agents/skills/` 与 `.claude/skills/` 为软链）

## Hard Gates

- 见 [project.md](<project.md>)。不要为 GoNavi 建空的 `vibe/ai-db/`。
- `vibe/` 与 `czz-docs/` 不进上游 `Syngnat/GoNavi` 的 PR。

## Git 快捷排除

- 遵循全局 GitHub commit scope。
- 当前任务目录之外的 `vibe/specs/<yyMMdd>/<HHmm-task-id>/` 默认快捷排除。
- `czz-docs/` 仅退役 README，未点名则排除。
- `vibe/specs/PROJECT_STATUS.md`、本目录 README、上游 PR 索引属于歧义 owner，不得只凭文件名排除。
