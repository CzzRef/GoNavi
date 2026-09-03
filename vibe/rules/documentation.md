# Documentation Rules

Tool: tool-neutral
Date: 2026-09-03

## Purpose

把 CodeNote 文档层级映射到本仓库。跨项目规则仍在宿主内核；本文件只写本仓路径。

## Authoritative Sources

| Layer | Location | Role |
| --- | --- | --- |
| Project adapter | [README.md](README.md) | 本仓入口 |
| Process hub | [../specs/PROJECT_STATUS.md](../specs/PROJECT_STATUS.md) | 当前主线、任务索引、未验项 |
| New process docs | [../specs/](../specs/) | `vibe/specs/<yyMMdd>/<HHmm-task-id>/` |
| Task cards | 上述 dated 目录里的 `task-card.md` | Standard non-requirement 的唯一过程 owner |
| Requirement Spec | 同目录 `raw-requirement.md` + `spec.md` | 文档位置、PR 留痕等需求增量 |
| Knowledge | [../knowledge/README.md](../knowledge/README.md) | 界面约定、核验通路、上游 PR 历史 |
| Legacy | [../../czz-docs/README.md](../../czz-docs/README.md) | 已归档指针与 MCP 研究稿；**不是**当前 task card / PR 范围 |

## Project Mapping

- **task card 只写在 `vibe/specs/`。** 禁止在 `czz-docs/` 新建或继续当作当前权威。
- 新过程文档目录名：`vibe/specs/<yyMMdd>/<HHmm-task-id>/`。日期用 Asia/Shanghai 日历，时间用创建该目录时的实测 `HHmm`。
- 上游 PR 的节奏化历史写在 [../knowledge/upstream-pr/](../knowledge/upstream-pr/README.md)，不写进 task card 当第二份范围表。
- 可复用界面规则写在 [../knowledge/ai-provider-ui-conventions.md](../knowledge/ai-provider-ui-conventions.md)，不在 PR 记录里复制像素。
- 不建 `vibe/ai-db/`。
- 需求 Manifest 尚未单独立项；本轮文档位置需求的 canonical 就是本文件与 `upstream-pr` 索引。

## Closeout

- 改了当前主线、task card、PR 状态或未验项时更新 [PROJECT_STATUS.md](../specs/PROJECT_STATUS.md)。
- 可复用结论升到 knowledge / 项目规则，不把过程稿当成当前行为。
