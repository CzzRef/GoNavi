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

## 3.5 测试必须随 PR 一起提，不能剥离

这条曾被质疑「上游会不会不关心测试」，实测答案是相反的：

- `upstream/dev` 里有 **521 个 `_test.go`** 与 **540 个 `.test.ts(x)`**；Go 源文件共 1052 个，近一半是测试。
- `.github/workflows/backend-tests.yml` 的触发是 `on: pull_request: branches: [dev]`，步骤为 `go test ./... -count=1 -timeout=30m`。**PR 一提交，全量后端测试就跑。**

剥离测试的实测后果：

| 做法 | 后果 |
| --- | --- |
| 全部剔除测试改动 | `internal/ai/provider` 编译失败：`not enough arguments in call to buildCodexCLIEnv`。我们改了该函数签名，上游旧测试仍按旧签名调用 |
| 只剔除新增的 18 个测试文件 | `internal/ai/service` 编译失败：`undefined: newProviderManagementTestService`。修改过的既有测试用了新文件里的共享辅助函数 |

本次 72 文件里测试的构成是「新增 18 + 改动既有 12」。那 12 个不是额外负担，而是**上游自己的用例因源码契约变化必须同步修改**；删掉它们等于把上游测试留成红的。

真正该剔除的是**不适合上游的测试**，本次剔了一条：读 `.css` 文本并断言 `left: 31px` 等像素字面量的用例。它绕开了仓库 `testPolicy` 守卫的正则（只匹配 `.ts/.tsx`），却会因任何一次样式格式化而误报。

想缩小 PR 时，正确手段是**按主题拆成多个各自带测试的 PR**，不是抽走测试。

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

在 `feat/ai-provider-management-v2` 上实测，不是在 `czz-dev` 上：

| 项目 | 结果 |
| --- | --- |
| `go build ./internal/... ./cmd/...` | 退出码 0 —— 证明 §4 的 hunk 取舍没有留下悬空引用 |
| `go vet ./internal/ai/... ./internal/app/... ./internal/mcpserver/` | 退出码 0 |
| `go test ./internal/ai/... ./internal/mcpserver/` | 全部通过 |
| `go test ./internal/app/` | 仅 `TestFetchReleaseByURLFallsBackToCacheOn403` 失败，属项目入口已记录的预存在失败（5 个测试共享 `updateReleaseCache` 全局状态） |
| `npx tsc --noEmit` | 退出码 0 |
| `npx vitest run` | 5035/5037，与 `czz-dev` 完全一致；两条失败在基线分支上同样复现，其中 `testPolicy` 标记的两个违规文件是上游自己的 |
| 密钥/个人路径扫描 | 未检出令牌、私钥或 `/Users/...` 绝对路径 |
| 新增 `TODO` / `FIXME` / `console.log` | 0 处 |

## 7. PR 现状

| PR | 分支 | 状态 | 说明 |
| --- | --- | --- | --- |
| [#1130](https://github.com/Syngnat/GoNavi/pull/1130) | `feat/ai-provider-management` | 已关闭 | 首版，72 文件 +9097/−1422。因测试范围需重新界定而关闭 |
| [#1131](https://github.com/Syngnat/GoNavi/pull/1131) | `feat/ai-provider-management-v2` | 开启中 | 72 文件 +9087/−1422，`MERGEABLE`。与首版只差 §3.5 里剔掉的那条脆弱用例，−10 行 |

两条分支都保留，不删除也不 force-push：#1130 的历史与讨论完整可对照。

## 8. 尚未做的

- `frontend/wailsjs/` 三个生成文件带 −330 行重排，上游通常自行生成，是否包含仍待确认。
- SQL 例程安全的独立 PR 未建分支。
- 若上游在评审期间前进，#1131 需按 §2 的步骤重做压平并 force-push，或另起 `-v3`。
