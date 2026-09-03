# PROJECT_STATUS

Tool: grok
Updated: 2026-09-03

GoNavi 过程枢纽。不复制 CodeNote 内核。

## 当前主线

1. **供应商设置页第十六轮（编辑收缩）** — 代码已在 `czz-dev`；r59 实机于 2026-09-03 由用户确认通过。过程 owner：[0000-provider-editor-compact/task-card.md](260903/0000-provider-editor-compact/task-card.md)。上游 PR 未开，节奏记录：[03-260903-pending-编辑收缩展示](../knowledge/upstream-pr/03-260903-pending-编辑收缩展示.md)。
2. **文档落点迁入 vibe** — 本轮需求：[1957-pr-docs-vibe-init/spec.md](260903/1957-pr-docs-vibe-init/spec.md)。task card 与 PR 范围不再以 `czz-docs/` 为当前权威。

## 已合入上游

| 序号 | 记录 | GitHub |
| --- | --- | --- |
| 01 | [01-260901-squash-供应商目录与CLI发现](../knowledge/upstream-pr/01-260901-squash-供应商目录与CLI发现.md) | [#1131](https://github.com/Syngnat/GoNavi/pull/1131) MERGED；[#1130](https://github.com/Syngnat/GoNavi/pull/1130) CLOSED |
| 02 | [02-260902-followup-设置页折叠与CLI续命](../knowledge/upstream-pr/02-260902-followup-设置页折叠与CLI续命.md) | [#1134](https://github.com/Syngnat/GoNavi/pull/1134) MERGED |

## 未验 / 未开

- 已接入芯片行拖拽排序仍未做。
- 真实模型回复、Windows/Linux 实机、签名发布包。
- 下一份上游 PR（03 编辑收缩）必须从最新 `upstream/dev` 另起压平分支，**不要** `czz-dev` → `dev`。`vibe/` 与 `czz-docs/` 永久排除。

## 入口

- 项目规则：[../rules/README.md](../rules/README.md)
- 知识：[../knowledge/README.md](../knowledge/README.md)
- 上游 PR 索引：[../knowledge/upstream-pr/README.md](../knowledge/upstream-pr/README.md)
- 退役 `czz-docs`：[../../czz-docs/README.md](../../czz-docs/README.md)
