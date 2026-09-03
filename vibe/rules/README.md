# GoNavi 项目入口

Tool: tool-neutral
Hosts: any

本目录只维护本仓库特化路由。共享 primary / additive / evidence 由宿主已加载的 CodeNote Rule Kernel 选择；这里不复制算法。

本仓库**没有** `intent-note-gate`。

## Required Reads

1. 宿主已注入的 CodeNote 内核（新任务一次）
2. [documentation.md](documentation.md)（Standard / Controlled、文档治理、task card 位置）
3. [PROJECT_STATUS.md](../specs/PROJECT_STATUS.md)（进行中的供应商 UI / 上游 PR）
4. [project.md](project.md)（SQL 判定、三执行面、fork 与上游 PR 边界）
5. [workflow.md](workflow.md)（核验包、测试命令）

## Project Routes

- 过程枢纽：[PROJECT_STATUS.md](../specs/PROJECT_STATUS.md)
- 任务索引：[../specs/README.md](../specs/README.md)
- 可复用知识：[../knowledge/README.md](../knowledge/README.md)
- 上游 PR 历史：[../knowledge/upstream-pr/README.md](../knowledge/upstream-pr/README.md)
- 供应商界面约定：[../knowledge/ai-provider-ui-conventions.md](../knowledge/ai-provider-ui-conventions.md)
- 核验通路：[../knowledge/gonavi-verify-build-restart.md](../knowledge/gonavi-verify-build-restart.md)
- 项目 Skill：`gonavi-verify-build-restart`（CodeNote `Skills/projects/gonavi/`，本仓 `.agents/skills/` 与 `.claude/skills/` 为软链）

## Hard Gates

- 见 [project.md](project.md)。不要为 GoNavi 建空的 `vibe/ai-db/`。
- `vibe/` 与 `czz-docs/` 不进上游 `Syngnat/GoNavi` 的 PR。

## Git 快捷排除

- 遵循全局 GitHub commit scope。
- 当前任务目录之外的 `vibe/specs/<yyMMdd>/<HHmm-task-id>/` 默认快捷排除。
- `czz-docs/` 里未点名的研究稿默认排除。
- `vibe/specs/PROJECT_STATUS.md`、本目录 README、上游 PR 索引属于歧义 owner，不得只凭文件名排除。
