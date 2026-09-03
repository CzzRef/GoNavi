# 02-260902-followup-设置页折叠与CLI续命

Status: `merged`
GitHub: [#1134](https://github.com/Syngnat/GoNavi/pull/1134) MERGED
Branch: `feat/ai-provider-ui-followup`（从当时最新 `upstream/dev` 另起，单提交）
日期: 260902
方式: followup
过程稿: [供应商管理总卡 · 第十四轮](../specs/260901/0000-ai-provider-management/task-card.md#第十四轮折叠栏箭头与顶栏间距)

## 1. 生命周期

### 1.1 主题从总卡拆出

#1131 合入后，工作区相对 squash 头 `692e17d7` 仍有设置页与 CLI 增量。当时拆成三块候选：

| 代号 | 主题 | 去向 |
| --- | --- | --- |
| 9A | 设置页跟进：折叠栏箭头、顶栏三列、已隐藏底栏可拖高、提示不截鼠标、目录最小宽 168→128 | **进本 PR** |
| 9B | CLI 流式与空闲续命（已提交 `caf98c1c`）：Grok/Cursor 真流式、Codex 思考边扫边推；stdout 重置 3 分钟空闲钟、硬上限 15 分钟 | **进本 PR** |
| 9C | SQL 例程安全 | **不进**。2026-09-01 产品决定回滚 `SQLOpRoutine` 无条件禁止 |

压平分支从当时最新 `upstream/dev`（记录里曾写 `436a29c8`，开 PR 前再 rebase 到 `50dc9839`）另起，**没有**把 `czz-dev` 直接对准 `dev`。

### 1.2 9A 实机来源

编辑已接入项时「认证与连接 / 更多设置」看不见箭头，点标题旁空白也不收展。WKWebView 给 `<summary>` 设 `display:flex` 会藏掉原生三角。改为内层满宽 flex 条 + 行尾 caret，整栏可点；ⓘ 仍 `stopPropagation`。

已隐藏改为目录底栏抽屉：有隐藏项才出现；顶部分隔条可拖高钉住。提示改为 `passThroughHintTooltip`（离场 0ms，浮层不截鼠标）。目录最小宽 128px。核验包当时写到 r45–r48；独立预览曾因单文件内联把 React 的 `$&` 展开而崩溃，已重建。

### 1.3 合入

标题：`opt(ai): 设置页折叠栏、隐藏抽屉与 CLI 空闲续命`。23 文件，+1166/−141。合入 UTC `2026-09-02T02:51:37Z`。

故意省略：当时的 `czz-docs/`、核验截图、`.gitignore` / `package.json.md5`、#1131 已丢掉的 CSS 像素断言。

## 2. GitHub README

实际上游正文：

## Summary

Follow-up to #1131 on the AI provider settings page and local CLI backends.

**Settings page**

- Editor disclosure rows (`认证与连接` / `更多设置`) keep a visible caret beside the title and are clickable across the whole header. WKWebView was hiding the native `<summary>` triangle when it was `display:flex`.
- Hidden catalog is a bottom drawer: it can be resized, a row click adds or edits (same as the main catalog), and the eye still only restores visibility.
- Hover hints on this page close as soon as the pointer leaves the trigger (`mouseLeaveDelay: 0`) and use `pointer-events: none`, so a neighbouring card can show its hint immediately.
- Catalog minimum width is 128px (one compact column). Search uses a shorter placeholder.

**CLI streaming / idle timeout**

- Grok and Cursor stream tokens for real; Codex thinking is forwarded while it is scanned.
- Any stdout resets a 3-minute idle clock; hard cap is 15 minutes. A hung CLI still dies.

## User-visible

- Provider settings: disclosures, hidden drawer, pass-through hints, narrower catalog.
- Long CLI replies no longer sit silent until a wall-clock timeout if they are still producing output.

## Key files

- `frontend/src/components/ai/AISettingsProvidersSection.tsx` / `.css` / `mounted.test.tsx`
- `frontend/src/components/ai/useAIProviderLayout.ts`
- `frontend/src/components/common/tooltipTiming.ts`
- `internal/ai/provider/cli_idle_watchdog.go`
- `internal/ai/provider/{codex,cursor,grok}_cli.go`
- six locale files: `search_short` / `resize_hidden` only

23 files, +1166 / −141 against current `dev`.

## Risk / compatibility

- Layout prefs still live in `gonavi.ai.providers.layout.v1` (now also `hiddenPaneHeight`). Clearing that key restores the catalog; it is not synced.
- Idle timeout only applies to these local CLI providers, not HTTP OpenAI-compatible backends.
- Stored-procedure / `CALL` policy is **not** in this PR. GoNavi continues to execute routines through MCP and the app according to the connection's existing SQL permission.

## Verification

- `npx vitest run` on the provider settings / tooltip files: 76/76
- `go test ./internal/ai/provider/ ./internal/ai/safety/ ./internal/mcpserver/`
- `go vet ./internal/ai/... ./internal/app/... ./internal/mcpserver/`
- Rebased onto current `dev` (`50dc9839`, download generation check) with a clean merge

Intentionally omitted: `czz-docs/`, verification screenshots, `.gitignore` / `package.json.md5`, and the CSS-text assertion that was dropped from #1131.

对照全局 GitHub 规则 §3：本条已按 Summary / 关键文件 / 用户可见 / 风险 / 核验 分节，比 #1131 更贴近清单。仍无 UI 截图。

## 3. Skills 与规则

- 折叠 caret、提示不截鼠标、目录最小宽写入 [ai-provider-ui-conventions.md](../knowledge/ai-provider-ui-conventions.md)。
- 核验停旧起新改为调用 Skill 脚本 `restart.sh`，见 [gonavi-verify-build-restart.md](../knowledge/gonavi-verify-build-restart.md)。
- 全局 §3 清单本条基本齐；例程策略在 Risk 节显式声明不在本 PR。
- 排除项仍含当时的 `czz-docs/`；现另排除 `vibe/`。

## 4. 范围快照（合入口径）

| 口径 | 文件 | 行数 | 说明 |
| --- | --- | --- | --- |
| 可进 PR | 23 | +1166/−141 | 9A+9B |
| 永久排除 | 本机文档、核验产物、工具噪声、CSS 像素用例、SQL 例程 | — | 9C 已回滚 |
