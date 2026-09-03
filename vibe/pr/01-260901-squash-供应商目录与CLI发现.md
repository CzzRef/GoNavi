# 01-260901-squash-供应商目录与CLI发现

Status: `merged`
GitHub: [#1131](https://github.com/Syngnat/GoNavi/pull/1131) MERGED；首版 [#1130](https://github.com/Syngnat/GoNavi/pull/1130) CLOSED
Branch: `feat/ai-provider-management-v2`（首版 `feat/ai-provider-management`；噪声归一 `feat/ai-provider-management-v3` 未开 PR）
日期: 260901
方式: squash
过程稿: [供应商管理总卡](../specs/260901/0000-ai-provider-management/task-card.md)

## 1. 生命周期

### 1.1 从本机全量差收敛

相对当时 `upstream/dev`（记录刷新于 2026-09-01 21:05，基线 `c6fb9251`）本机 `czz-dev` 曾是 104 文件。收敛口径：

| 口径 | 文件 | 行数 | 说明 |
| --- | --- | --- | --- |
| A 全量 | 104 | +13171 / −1439 | `czz-dev` 相对 `upstream/dev` |
| B 去掉 `czz-docs/` | 97 | +10791 / −1439 | 本机文档 |
| C 再去掉 SQL 例程安全 | 88 | +10201 / −1425 | 独立主题，当时拟另开，后回滚 |
| D 再去掉本机产物 | 74 | — | `.codemark/` 与 `build/evidence/` |
| **E 再去掉工具噪声 = 实际 PR** | **72** | **+9097 / −1422**（#1131 为 +9087/−1422） | `.gitignore` 与 `frontend/package.json.md5` |

压平方式：从 `upstream/dev` 另起分支，`merge --squash czz-dev`，再按上表退回排除项，并对两个混主题文件按 hunk 取舍。本地 `czz-dev` 历史提交不改写。

### 1.2 混主题文件按 hunk 取舍

| 文件 | 保留 | 丢弃 |
| --- | --- | --- |
| `internal/ai/types.go` | 供应商字段与 `CLICapabilityView` | `SQLOpRoutine` 枚举 |
| `frontend/src/types.ts` | `disabledModels` / `customModels` / `effort` | `operationType` 加入 `"routine"` |

`SQLOpRoutine` 当时只被那 9 个 SQL 文件引用，丢弃 hunk 不留悬空引用。2026-09-01 产品决定后，例程无条件禁止已从 `czz-dev` 回滚，后续 PR 也不再拣回。

### 1.3 测试随主题走，只剔脆弱用例

曾被质疑「上游会不会不关心测试」。实测相反：当时 `upstream/dev` 有 521 个 `_test.go` 与 540 个 `.test.ts(x)`；`.github/workflows/backend-tests.yml` 在对 `dev` 的 PR 上跑 `go test ./...`。

| 做法 | 后果 |
| --- | --- |
| 全部剔除测试改动 | `internal/ai/provider` 编译失败：`buildCodexCLIEnv` 参数个数对不上 |
| 只剔除新增测试文件 | `internal/ai/service` 编译失败：`undefined: newProviderManagementTestService` |

72 文件里测试是「新增 18 + 改动既有 12」。那 12 个是上游用例因源码契约变化必须同步。真正剔除的是读 `.css` 文本并断言 `left: 31px` 的像素用例——它绕开 `testPolicy`（只扫 `.ts/.tsx`），样式一格式化就误报。

### 1.4 #1130 关闭，#1131 合入

- [#1130](https://github.com/Syngnat/GoNavi/pull/1130) `feat/ai-provider-management`：72 文件 +9097/−1422。因测试范围需重新界定于 2026-09-01 关闭（UTC `2026-09-01T05:53:45Z`）。讨论与 diff 保留对照。
- [#1131](https://github.com/Syngnat/GoNavi/pull/1131) `feat/ai-provider-management-v2`：与首版只差剔掉那条脆弱 CSS 像素用例。squash `692e17d7`，merge `d9a081a2`，72 文件 +9087/−1422。合入 UTC `2026-09-01T06:39:48Z`。

三条本地分支都保留，不删除也不 force-push。

### 1.5 生成文件噪声（v3，未用来替换已合入的 #1131）

`frontend/wailsjs/go/models.ts` 相对上游曾有 694 行改动，`git diff -w` 真实内容约 34 行（`CLICapabilityView` 与 `disabledModels` / `customModels`）。其余是空行风格：本机 wails 在空行写 `\t`，上游写空。`feat/ai-provider-management-v3` 做了空白行归一并推送，但 #1131 已合入带空白噪声的版本，**v3 不再开 PR 去替换它**。下次 `wails build` 若不带 `-skipbindings` 会把噪声打回来；后续 PR 不要再带 `models.ts`，除非有真实字段增量。

### 1.6 压平分支上的核验（当时）

在 `feat/ai-provider-management-v2` 上，不是在 `czz-dev` 上：

| 项目 | 结果 |
| --- | --- |
| `go build ./internal/... ./cmd/...` | 退出码 0 |
| `go vet ./internal/ai/... ./internal/app/... ./internal/mcpserver/` | 退出码 0 |
| `go test ./internal/ai/... ./internal/mcpserver/` | 通过 |
| `go test ./internal/app/` | 仅 `TestFetchReleaseByURLFallsBackToCacheOn403` 失败（预存在） |
| `npx tsc --noEmit` | 退出码 0 |
| `npx vitest run` | 5035/5037，两条失败在基线 `dev` 同样复现 |
| 密钥/个人路径扫描 | 未检出令牌、私钥或 `/Users/...` 绝对路径 |

## 2. GitHub README

#1130 与 #1131 正文相同（英文）。以下为实际上游正文：

Reworks the AI provider settings page and the CLI backends behind it, so providers can be discovered, configured and queried entirely from within GoNavi instead of relying on external scripts.

**Provider settings page**

- Searchable add flow without an endpoint-first step; the provider catalog is collapsible, resizable, independently scrollable, and individual candidates can be hidden and restored.
- Saved providers render as chips carrying default selection, edit, and a hover-revealed delete that still confirms before removing.
- Field hints collapse into one icon per heading and keep their copy in the accessible name, instead of full-width note blocks that pushed the form out of view on short panes.
- Error reveal after a rejected save scrolls only the editor's own container. `scrollIntoView` also scrolls `overflow: hidden` ancestors, which then have no scrollbar to put them back.
- The model picker separates default selection from enable/disable management and reserves a fixed popup height, so switching tabs cannot resize the popup against its trigger.
- Save and save-as share one split button; singleton CLI presets omit save-as entirely because they reuse a single machine login, and a copy keeps a renamed draft rather than always appending a suffix.

**CLI backends and candidates**

- Adds Cursor and Grok CLI providers; extends Claude, Codex and CodeBuddy.
- Resolves CLI commands through a shared lookup that handles NVM installs and shims, backed by a detection cache.
- Capability and model-catalog probes let the app enumerate models itself. Effort levels are projected from each CLI's own value range rather than duplicated in the frontend, since the three CLIs disagree on both range and rejection semantics.

**Shared modal**

Restores vertical resizing for the shared draggable modal. Callers pass an inline `height` through antd `styles.content`, which a normal author rule cannot override, so only width was adjustable.

**Scope and verification**

Six locales updated. This branch is a single squashed commit on top of `dev`; it deliberately excludes an unrelated SQL routine-safety change and all local working documents.

- `go build ./internal/... ./cmd/...` — clean
- `go vet ./internal/ai/... ./internal/app/... ./internal/mcpserver/` — clean
- `go test ./internal/ai/... ./internal/mcpserver/` — pass
- `go test ./internal/app/` — only `TestFetchReleaseByURLFallsBackToCacheOn403` fails; it shares `updateReleaseCache` global state with four other tests and also fails on `dev`
- `npx tsc --noEmit` — clean
- `npx vitest run` — 5035/5037; both failures reproduce unchanged on `dev`

Not covered: real model responses, Windows/Linux desktops, and signed release packaging.

对照全局 GitHub 规则 §3：Summary、关键行为、核验、未覆盖项已写；关键文件清单偏弱（用主题段落代替路径表）；无 UI 截图（独立核验包不入库）。#1131 相对 #1130 的差异（剔 CSS 像素用例）**没有**写进正文，只存在于本条生命周期。

## 3. Skills 与规则

- 界面约定当时写在本机 `czz-docs/ai-provider-ui-conventions.md`，现迁到 [ai-provider-ui-conventions.md](../knowledge/ai-provider-ui-conventions.md)。本 PR **不提交**该文件。
- 核验通路 [gonavi-verify-build-restart.md](../knowledge/gonavi-verify-build-restart.md) 与项目 Skill `gonavi-verify-build-restart` 同轮沉淀。Skill 物理来源在 CodeNote `Skills/projects/gonavi/`；本仓 `.agents/skills/` 与 `.claude/skills/` 为软链，gitignore。
- 全局：commit 排除本机文档/产物符合 `github/rules.md` §2；PR 正文基本覆盖 §3 的 Summary / 用户可见行为 / 风险（例程排除）/ 核验。
- 本仓压平规则见 [README.md](README.md)。

## 4. 范围快照（合入口径）

| 口径 | 文件 | 行数 | 说明 |
| --- | --- | --- | --- |
| 可进 PR | 72 | +9087/−1422 | #1131 squash |
| 永久排除 | `czz-docs/`（当时含 task card 与本范围表）、`.codemark/`、`build/evidence/`、`.gitignore` `.agents/`、`package.json.md5`、CSS 像素用例 | — | 现另排除整个 `vibe/` |
