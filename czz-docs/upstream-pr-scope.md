# 上游 PR 范围记录

本文件只是记录，不随 PR 提交。**czz-docs/ 整个目录永远不进上游 PR**：上游 `Syngnat/GoNavi` 没有这个目录，这里放的是本机工作文档。

刷新于 2026-09-01 21:05。基线为最新 `upstream/dev` 的 `c6fb9251`（已含 [#1131](https://github.com/Syngnat/GoNavi/pull/1131) 与 [#1132](https://github.com/Syngnat/GoNavi/pull/1132)）。本地 `czz-dev` 仍以旧基线 `57658ed1` 为合并点，相对上游 `8 35`：上游多 8 个提交（含 1131 squash 与 1132），本地多 35 个历史提交。下一份 PR 必须从最新 `upstream/dev` 另起压平分支，不能把 `czz-dev` 直接对准 `dev`。

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

该类型曾在任何权限级别下都不放行。**2026-09-01 已按产品决定回滚**：MCP / Headless / 应用按连接既有 SQL 权限执行例程，不再无条件拒绝。历史排除记录保留，避免再被拣进 PR。

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

| PR / 分支 | 状态 | 说明 |
| --- | --- | --- |
| [#1130](https://github.com/Syngnat/GoNavi/pull/1130) `feat/ai-provider-management` | 已关闭 | 首版，72 文件 +9097/−1422。因测试范围需重新界定而关闭 |
| [#1131](https://github.com/Syngnat/GoNavi/pull/1131) `feat/ai-provider-management-v2` | **已合入** `dev` | squash `692e17d7`，merge `d9a081a2`，72 文件 +9087/−1422。与首版只差 §3.5 里剔掉的那条脆弱用例 |
| `feat/ai-provider-management-v3` | 已推送，未开 PR | 在 v2 上做了 §7.5 的生成文件噪声归一。1131 已合入带空白噪声的 `models.ts`，v3 不再用来替换它 |

三条分支都保留，不删除也不 force-push。#1130 的历史与讨论完整可对照。

## 7.5 生成文件的噪声归一

`frontend/wailsjs/go/models.ts` 相对上游有 694 行改动，`git diff -w` 实测**真实内容只有 34 行**（`CLICapabilityView` 类与 `disabledModels` / `customModels` 字段）。其余 660 行是空行风格差异：本机 wails 在空行写 `\t`，上游写空。

处理方式是只把「整行仅由空白构成」的行归一化，不动任何有内容的行；归一后 `npx tsc --noEmit`、`npx vitest run`（5034/5036）与 `npm run build` 均通过。

另两个绑定文件本就无噪声：`Service.d.ts` +6、`Service.js` +12，内容是 `AIGetCLICapabilities` / `AIGetCLIModelCatalog` / `AIListCLIModels` 三个桥接方法。

**注意**：若下次 `wails build` 不带 `-skipbindings` 重新生成绑定，噪声会回流，提 PR 前需重做这一步。#1131 已合入带空白噪声的版本，下一份 PR 不要再带 `models.ts` 除非有真实字段增量。

## 8. 尚未做的

- `czz-dev` 已合并最新 `upstream/dev`（`436a29c8`）。
- 已隐藏底栏抽屉的实机观感仍待用户确认。
- 下一份上游 PR 压平分支为 `feat/ai-provider-ui-followup`，只含 9A+9B，等明确「开 PR」再提。

## 9. #1131 合入后的下一份 PR 候选

相对 squash 头 `692e17d7`（即已合入的 #1131）的工作区增量，去掉 `czz-docs/`、`.codemark/`、`build/evidence/` 后是 37 文件。其中仍须永久排除：

| 排除 | 理由 |
| --- | --- |
| `czz-docs/` | 本机工作文档，上游没有该目录 |
| `.codemark/`、`build/evidence/` | 本机核验产物 |
| `.gitignore` 的 `.agents/`、`frontend/package.json.md5` | 本机工具噪声 |
| `AISettingsProvidersSection.test.tsx` 读 CSS 断言 `left: 31px` | #1131 已按 §3.5 剔除，不要回流 |

压平分支从最新 `upstream/dev`（`436a29c8`）另起，不要 `czz-dev` → `dev` 直接开 PR。

### 9A 设置页跟进（推荐下一份）

折叠栏箭头、顶栏三列、已隐藏底栏可拖高、提示不截鼠标、目录最小宽 168→128。约 15 文件：

- `frontend/src/components/ai/AIProviderModelSelect.tsx` 与对应测试
- `frontend/src/components/ai/AISettingsProvidersSection.tsx` / `.css` / `mounted.test.tsx`
- `frontend/src/components/ai/useAIProviderLayout.ts` 与对应测试
- `frontend/src/components/common/tooltipTiming.ts` 与对应测试
- `shared/i18n/{zh-CN,zh-TW,en-US,ja-JP,de-DE,ru-RU}.json` 的 `search_short` / `resize_hidden`

### 9B CLI 流式与空闲续命

已提交 `caf98c1c`，8 文件：`cli_idle_watchdog.go` + 测试，以及 `codex_cli` / `cursor_cli` / `grok_cli` 的流式与超时改动。与 9A 一并进入 `feat/ai-provider-ui-followup`。

### 9C SQL 例程安全 —— 已回滚，不进 PR

2026-09-01 产品决定：CodeNote「不得直接执行存储过程」约束的是 Agent 绕过产品打库；经 MCP 或 GoNavi 本身发出的语句按连接权限执行。`SQLOpRoutine` 无条件禁止已从 `czz-dev` 撤回到与 `upstream/dev` 一致，**不要再放进上游 PR**。
