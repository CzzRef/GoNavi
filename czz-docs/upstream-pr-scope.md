# 上游 PR 范围记录

本文件只是记录，不随 PR 提交。**czz-docs/ 整个目录永远不进上游 PR**：上游 `Syngnat/GoNavi` 没有这个目录，这里放的是本机工作文档。

刷新于 2026-09-01 13:45。基线为 `upstream/dev` 的 `57658ed1`，本地 `czz-dev` 为 `49d9ed2f`（含 28 个本地独有提交，已推送到 `origin/czz-dev`），`upstream/dev...czz-dev` 计数为 `0 28`，上游无遗留提交未合入。

## 1. 从 104 文件收敛到 72 文件

| 口径 | 文件 | 行数 | 说明 |
| --- | --- | --- | --- |
| A 全量 | 104 | +13171 / −1439 | `czz-dev` 相对 `upstream/dev` 的全部差异 |
| B 去掉 `czz-docs/` | 97 | +10791 / −1439 | 本机文档 7 文件 +2380 |
| C 再去掉 SQL 例程安全 | 88 | +10201 / −1425 | 该主题 9 文件 +590 / −14 |
| D 再去掉本机产物 | 74 | — | `.codemark/` 与 `build/evidence/` 共 14 文件 +1093 |
| **E 再去掉工具噪声 = 实际 PR** | **72** | **+9097 / −1422** | `.gitignore` 与 `frontend/package.json.md5` 共 2 文件 |

## 2. PR 分支已在本地建好

分支 `feat/ai-provider-management`，从 `upstream/dev` 起，单个提交 `eee677f0`，**尚未推送**。

复现方式：

```
git switch -c feat/ai-provider-management upstream/dev
git merge --squash czz-dev
# 按 §3 退回排除项，按 §4 处理两个混主题文件
git commit
```

本地 `czz-dev` 的 28 个提交原样保留，不改写。

## 3. 排除项与理由

- **`czz-docs/`（7 文件 +2380）** —— 研究稿、评估记录、任务卡与本文件。上游没有该目录，永久排除。
- **SQL 例程安全（9 文件 +590 / −14）** —— 见 §5，独立主题，建议单独开 PR。
- **`.codemark/` 与 `build/evidence/`（14 文件 +1093）** —— CodeMark 插件的本机状态与 r8/r9 的核验截图。**这两处已被并行会话提交进 `czz-dev`**，与任务卡自述的「构建目录、截图、本机证据与独立预览不纳入提交」相悖；PR 里必须排除，`czz-dev` 上是否回滚另议。
- **`.gitignore` 与 `frontend/package.json.md5`（2 文件）** —— 前者是本机代理发现目录的忽略项，后者是 wails 生成的本机哈希，对上游都是噪声。

## 4. 两个混主题文件的拆分

这两个文件同时承载供应商与 SQL 两个主题，必须按 hunk 取舍，不能整份保留或整份排除：

| 文件 | 保留 | 丢弃 |
| --- | --- | --- |
| `internal/ai/types.go` | 供应商字段与 `CLICapabilityView`（`@@ -98,23 +98,50 @@`） | `SQLOpRoutine` 枚举（`@@ -274,7 +301,12 @@`） |
| `frontend/src/types.ts` | `disabledModels` / `customModels` / `effort` | `operationType` 加入 `"routine"`（`@@ -885,7 +893,9 @@`） |

实测 `SQLOpRoutine` 只被那 9 个 SQL 文件引用，因此丢弃该 hunk 不会留下悬空引用。

## 5. SQL 例程安全是什么

与供应商 UI 无关，是一条独立的安全收紧。新增方言无关的 `IsRoutineSQL`，把两类语句判定为 `SQLOpRoutine`：

- **例程调用**：`CALL` / `EXEC` / `EXECUTE`，以及 SQL Server 的裸过程调用。
- **例程部署**：对 `PROCEDURE` / `FUNCTION` / `TRIGGER` / `ROUTINE` 做 `CREATE` / `ALTER` / `DROP`。

该类型在**任何权限级别下都不放行**，`PermissionFull` 也不放宽，并在三个执行面同步：MCP `execute_sql`、Headless CLI、AI 对话守卫。

## 6. PR 分支上的核验（2026-09-01）

在 `feat/ai-provider-management` 上实测，不是在 `czz-dev` 上：

| 项目 | 结果 |
| --- | --- |
| `go build ./internal/... ./cmd/...` | 退出码 0 —— 证明 §4 的 hunk 取舍没有留下悬空引用 |
| `go vet ./internal/ai/... ./internal/app/... ./internal/mcpserver/` | 退出码 0 |
| `go test ./internal/ai/... ./internal/mcpserver/` | 全部通过 |
| `go test ./internal/app/` | 仅 `TestFetchReleaseByURLFallsBackToCacheOn403` 失败，属项目入口已记录的预存在失败（5 个测试共享 `updateReleaseCache` 全局状态） |
| `npx tsc --noEmit` | 退出码 0 |
| `npx vitest run` | 5035/5037，与 `czz-dev` 完全一致；两条失败在基线分支上同样复现 |
| 密钥/个人路径扫描 | 未检出令牌、私钥或 `/Users/...` 绝对路径 |
| 新增 `TODO` / `FIXME` / `console.log` | 0 处 |

## 7. 尚未做的

- **未推送**，也未创建 PR。推送目标应为 `origin`（`CzzRef/GoNavi`）的 `feat/ai-provider-management`，PR base 为 `Syngnat/GoNavi` 的 `dev`。
- `frontend/wailsjs/` 三个生成文件带 −330 行重排，上游通常自行生成，是否包含值得在开 PR 前单独确认。
- SQL 例程安全的独立 PR 未建分支。
