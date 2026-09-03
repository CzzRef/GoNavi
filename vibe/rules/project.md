# GoNavi 项目约束

Tool: tool-neutral
Date: 2026-09-03

## 仓库是什么

Wails/Go + React 桌面端、Headless CLI、MCP server、Web server。AI-DB 治理的对象是**工具自身行为**，不是某个业务库。

## SQL 判定

`appcore.InspectSQL` 是 MCP 与 CLI 共用判定来源。只读判定保持 allowlist（`default: return false`）。

三个执行面必须同进同退：MCP `execute_sql`、Headless CLI、AI 对话面板。

2026-09-01 产品决定：经产品执行面发出的 `CALL` / `EXEC` / 存储过程按连接权限执行，**不要**用 `SQLOpRoutine` 无条件拒绝。该禁止已从 `czz-dev` 回滚，与 `upstream/dev` 对齐。

## Fork 与上游

- 工作分支：`czz-dev`
- 来源：`upstream` = `Syngnat/GoNavi`，默认合 `upstream/dev`
- 备份：`origin` = `CzzRef/GoNavi`
- **不要**用 `czz-dev` 直接对 `upstream/dev` 开 PR。从最新 `upstream/dev` 另起压平分支。
- 永久不进上游 PR：`vibe/`、`czz-docs/`、`.codemark/`、`build/evidence/`、`.gitignore` 的 `.agents/`、`frontend/package.json.md5`、读 CSS 断言像素的用例。

## 核验

独立核验程序走 [gonavi-verify-build-restart](../knowledge/gonavi-verify-build-restart.md)。编译与离线测试不构成实机通过。
