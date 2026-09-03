# Changes：PR 文档迁入 vibe

> Inventory only. 判断在 spec，本文件只回答改了哪些文件。

## 1. 概览

| 批次 | 提交 | 文件数 | 核心说明 |
| --- | --- | --- | --- |
| 第一波 | `337ba525` 等 | 见 §2 | 初始化 vibe、迁移 task card、PR 史先落 knowledge |
| 第二波（本工作区） | 本回合分批提交 | 见 §2 续 | PR 史升为 `vibe/pr/`；MCP 迁入 knowledge；`czz-docs` 只留 README |

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
| `vibe/knowledge/upstream-pr/` | 曾新增 | 第一波节奏化 PR 史；第二波改为跳转桩 |
| `vibe/pr/` | 迁移 | 从 `vibe/knowledge/upstream-pr/` 升为独立目录（01–03 + 模板 + 索引） |
| `vibe/knowledge/mcp-agent/` | 迁移 | 原 `czz-docs` 三份 Agent/MCP 研究稿 |
| `czz-docs/README.md` | 改写 | 仅退役路由，指向 vibe |
| `czz-docs/*.md` 正文与桩 | 删除 | 不再保留归档桩 |
| CodeNote `gonavi-verify-build-restart/SKILL.md` | 改动 | 相关规则改指 `vibe/knowledge/…` 与 `vibe/pr/`（本仓外） |

## 3. 明确没做的

| 对象 | 核心说明 |
| --- | --- |
| 产品 UI / 测试 | 本轮纯文档治理 |
| 上游 GitHub PR | 本波不改 #1155 远端正文；01/02 已合入 |
| 本波提交/推送 | 本回合分批提交，不推送 |
