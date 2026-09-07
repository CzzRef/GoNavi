# PROJECT_STATUS

Tool: grok
Updated: 2026-09-07

GoNavi 过程枢纽。不复制 CodeNote 内核。

## 当前主线

1. **DBX 供应商交互复现** — 设置中心 `ai-providers` 为 list/edit、品牌下拉（合作方空列）、拍平字段；表单候选项按内容定宽。过程 owner：[1616-dbx-provider-parity/task-card.md](260907/1616-dbx-provider-parity/task-card.md)。实现在 child `codex/260907-dbx-provider-parity`，待合入 `czz-dev`。
2. **供应商设置页第十六轮（编辑收缩）** — r59 实机已通过；上游 [#1155](https://github.com/Syngnat/GoNavi/pull/1155) 已开（压平分支 `feat/ai-provider-editor-compact`，不是 `czz-dev`）。该页主交互已被上一行替换；紧凑并排降为非目标。
3. **文档落点迁入 vibe** — 需求：[1957-pr-docs-vibe-init/spec.md](260903/1957-pr-docs-vibe-init/spec.md)。PR 史在 [../pr/](../pr/README.md)；MCP 研究稿在 [../knowledge/mcp-agent/](../knowledge/mcp-agent/README.md)；`czz-docs/` 仅退役 README。

## 已合入上游

| 序号 | 记录 | GitHub |
| --- | --- | --- |
| 01 | [01-260901-squash-供应商目录与CLI发现](../pr/01-260901-squash-供应商目录与CLI发现.md) | [#1131](https://github.com/Syngnat/GoNavi/pull/1131) MERGED；[#1130](https://github.com/Syngnat/GoNavi/pull/1130) CLOSED |
| 02 | [02-260902-followup-设置页折叠与CLI续命](../pr/02-260902-followup-设置页折叠与CLI续命.md) | [#1134](https://github.com/Syngnat/GoNavi/pull/1134) MERGED |

## 未验 / 未开

- 已接入芯片行拖拽排序仍未做。
- 真实模型回复、Windows/Linux 实机、签名发布包。
- [#1155](https://github.com/Syngnat/GoNavi/pull/1155) 待上游评审；`vibe/` 与 `czz-docs/` 未进该 PR。

## 入口

- 项目规则：[../rules/README.md](../rules/README.md)
- 知识：[../knowledge/README.md](../knowledge/README.md)
- 上游 PR 索引：[../pr/README.md](../pr/README.md)
- 退役 `czz-docs`：[../../czz-docs/README.md](../../czz-docs/README.md)
