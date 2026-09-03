# Changes：PR 文档迁入 vibe

> Inventory only. 判断在 spec，本文件只回答改了哪些文件。

## 1. 概览

| 批次 | 提交 | 文件数 | 核心说明 |
| --- | --- | --- | --- |
| 本轮工作区 | 未提交 | 见 §2 | 初始化 vibe、迁移 task card、重构 PR 史、退役 czz-docs 当前权威 |

## 2. 交付物清单

| 对象 | 类型 | 核心说明 |
| --- | --- | --- |
| `AGENTS.md` | 新增 | 无 Home 路径的仓库入口 |
| `vibe/rules/*` | 新增 | CodeNote 风格项目入口与文档/SQL/核验映射 |
| `vibe/specs/README.md` | 新增 | 任务索引 |
| `vibe/specs/PROJECT_STATUS.md` | 新增 | 过程枢纽 |
| `vibe/specs/260903/1957-pr-docs-vibe-init/` | 新增 | 本需求 raw + spec + changes |
| `vibe/specs/260901/0000-ai-provider-management/task-card.md` | 迁移 | 原 `czz-docs/ai-provider-management-task-card.md` |
| `vibe/specs/260903/0000-provider-editor-compact/task-card.md` | 迁移 | 原 `czz-docs/260903-provider-editor-compact-task-card.md` |
| `vibe/knowledge/ai-provider-ui-conventions.md` | 迁移 | 原 czz-docs 同名 |
| `vibe/knowledge/gonavi-verify-build-restart.md` | 迁移 | 原 czz-docs 同名 |
| `vibe/knowledge/upstream-pr/` | 新增 | 节奏化 PR 史 01–03 + 模板 |
| `czz-docs/README.md` | 新增 | 退役路由 |
| `czz-docs/*.md` 桩 | 改动 | 旧路径改指针 |
| CodeNote `gonavi-verify-build-restart/SKILL.md` | 改动 | 相关规则改指 vibe（本仓外） |

## 3. 明确没做的

| 对象 | 核心说明 |
| --- | --- |
| 产品 UI / 测试 | 本轮纯文档治理 |
| 上游 GitHub PR | 未开 03，未改 01/02 远端正文 |
| MCP 三份研究稿 | 仍留 `czz-docs/`，不是 task card |
