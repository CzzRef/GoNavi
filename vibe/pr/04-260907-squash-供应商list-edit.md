# 04-260907-squash-供应商list-edit

Status: `opened`
GitHub: [#1178](https://github.com/Syngnat/GoNavi/pull/1178)
Branch: `feat/ai-provider-list-edit`（从当时最新 `upstream/dev` `77116266` 另起 squash `5c1dcfb7`）
日期: 260907
方式: squash
过程稿: [DBX 供应商交互复现](../specs/260907/1616-dbx-provider-parity/task-card.md)

## 1. 生命周期

#1155 已 MERGED。本条把设置中心 `ai-providers` 主路径改成 DBX 式 list/edit，并拍平表单字段；合作方列只留空占位。

压平：主目录 `czz-dev` 先合入 `upstream/dev`（#1147），再合入 child `codex/260907-dbx-provider-parity`。开 PR 时从 `77116266` 另起分支，只检出该主题文件（不含 `vibe/`、不含 Oracle/#1147 回退、不含 `.agents/` gitignore）。

实测：47 文件，+1079/−1498。相关 vitest 27/27。核验包 r63 待实机确认。

## 2. GitHub README

正文已提交到 [#1178](https://github.com/Syngnat/GoNavi/pull/1178)，结构对齐全局 `github/rules.md` §3。

## 3. Skills 与规则

- 约定：[ai-provider-ui-conventions.md](../knowledge/ai-provider-ui-conventions.md) §0。
- 核验包 r61–r63 走 `gonavi-verify-build-restart`；产物不入库。
- 永久排除见 [README.md](README.md)。未把 `vibe/` 提交进上游。

## 4. 范围快照（开 PR 实测）

| 口径 | 文件 | 行数 | 说明 |
| --- | --- | --- | --- |
| 可进 PR | 47 | +1079/−1498 | `5c1dcfb7` on `77116266` |
| 永久排除 | `vibe/`、`czz-docs/`、核验产物、工具噪声、合作方预设 | — | Skip TLS / 代理 / MCP 安装条不在本条 |
