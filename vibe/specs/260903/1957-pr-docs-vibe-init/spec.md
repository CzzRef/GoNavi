# Standard Requirement Spec: PR 记录迁入 vibe

Tool: grok
Date: 2026-09-03
Status: `confirmed`
Documentation level: `standard requirement`

Raw source: [raw-requirement.md](raw-requirement.md)
Canonical target: 本仓尚无独立 Requirement Manifest；落点权威为 [documentation.md](../../../rules/documentation.md) 与 [upstream-pr/README.md](../../../knowledge/upstream-pr/README.md)

> 本 `spec.md` 是本条 Standard requirement 的过程 owner。不另建 Task Card。

## Task Documentation Sync Group

- Group key: `dsg:gonavi:260903-pr-docs-vibe-init`
- Group owner: this `spec.md`
- Predecessor key: `gonavi-ai-provider-management`（owner 曾是 `czz-docs/ai-provider-management-task-card.md`）
- Git document prefixes: `vibe/`、`AGENTS.md`、`czz-docs/README.md` 与退役桩
- Durable document members: 本目录 raw/spec/changes、`vibe/rules/`、`vibe/specs/README.md`、`vibe/specs/PROJECT_STATUS.md`、`vibe/knowledge/`、迁入的两张 task card、退役 `czz-docs` 指针
- Declared code/config dependencies: 无产品代码。Skill 指针在 CodeNote `gonavi-verify-build-restart/SKILL.md`
- Linked authorities: CodeNote 文档层级、`github/rules.md` §3、session-title、communication §8
- Excluded: 前端/后端产品 diff、核验 `.app`、`build/evidence/`
- Lookup contract: `get --lookup-only` returns `present/freshness=unchecked`
- Active-group lifecycle: 本 key 替换旧 `gonavi-ai-provider-management` 作为**文档落点** Work Unit；供应商 UI 产品线改由两张 vibe task card 各自成组，不把新旧 Work Unit 并进同一 block

```json documentation-sync-group-v1
{
  "schema": "documentation-sync-group-v1",
  "group_key": "dsg:gonavi:260903-pr-docs-vibe-init",
  "group_owner": "vibe/specs/260903/1957-pr-docs-vibe-init/spec.md",
  "documents": [
    "vibe/specs/260903/1957-pr-docs-vibe-init/raw-requirement.md",
    "vibe/specs/260903/1957-pr-docs-vibe-init/spec.md",
    "vibe/specs/260903/1957-pr-docs-vibe-init/changes.md",
    "AGENTS.md",
    "vibe/rules/README.md",
    "vibe/rules/documentation.md",
    "vibe/specs/README.md",
    "vibe/specs/PROJECT_STATUS.md",
    "vibe/knowledge/upstream-pr/README.md",
    "vibe/specs/260901/0000-ai-provider-management/task-card.md",
    "vibe/specs/260903/0000-provider-editor-compact/task-card.md",
    "czz-docs/README.md"
  ],
  "dependencies": [],
  "validators": [],
  "git_scope_prefixes": [
    "vibe/",
    "AGENTS.md",
    "czz-docs/"
  ]
}
```

## Requirement Delta

- Add: 按 CodeNote 规则初始化本仓 `vibe/`（`rules` / `specs` / `knowledge`）；PR 历史目录 `vibe/knowledge/upstream-pr/`；每条记录标题 `序号-日期-方式-概要`；内部必须含生命周期、GitHub README、skills/规则（含全局 PR 规则同步）。
- Modify: 既有 `czz-docs/upstream-pr-scope.md` 重构进编号记录 + 索引；task card / 界面约定 / 核验通路的当前权威改到 `vibe/`；索引链接全部改指新位置。
- Remove / supersede: `czz-docs/` 作为 **task card 与 PR 范围的当前权威**。目录可留 MCP 研究稿与归档桩。
- Confirmed facts: 用户明确「不应该在 czz-docs」= 授权物理迁移 task card，不是只加 `archived-linked` 而不搬文件。
- Pending decisions: 无。不开上游 PR、不提交本轮文档，除非另授。
- Acceptance criteria: `vibe/rules` 可发现；两张 task card 在 `vibe/specs/<yyMMdd>/…/task-card.md`；至少三条 PR 史；`czz-docs` 旧路径是指针；全局 §3 在每条 PR 记录的 README 节可对照。

## Requirement Change Review

- Scan scope: 当前请求 → 本 raw/spec → 本仓 documentation 落点 → 供应商 task card / PR 范围表 → 项目 Skill 指针。
- Visible changes:
  - added: `vibe/` 初始化；`upstream-pr` 节奏史；本 spec。
  - changed: 索引与相对链接。
  - removed: 无文件销毁；MCP 三份研究稿仍在 `czz-docs/`。
  - superseded: `czz-docs/ai-provider-management-task-card.md`、`czz-docs/260903-provider-editor-compact-task-card.md`、`czz-docs/upstream-pr-scope.md`、`czz-docs/ai-provider-ui-conventions.md`、`czz-docs/gonavi-verify-build-restart.md` 的**当前权威**角色。
- Conflict classification: `supersede`
- Conflict evidence: 用户指定 task card 不在 `czz-docs`；旧 group_owner 路径失效。
- Decision status: `explicit-current-request`
- Decision source: 本会话用户消息全文。
- User-facing review shown: 本表。
- Post-sync rescan: `pass`（旧路径改为归档桩；Skill 指针改到 `vibe/knowledge/…`）。

## Prior Task Overlap

- Relationship: `supersedes` 文档落点；产品线 `continuation` 于两张已迁 task card
- Document governance: 旧 `gonavi-ai-provider-management` 只记为 predecessor，不把其 documents 并进本 block
- Execution logic verification: 不改产品代码
- Traceability: `supersede` 落点 + `delta-only` 产品线

## Canonical Merge

- Base project version: 无 Manifest 版本号
- Result: [documentation.md](../../../rules/documentation.md) 成为 task card / spec 路径的项目映射
- Merge status: 本仓 canonical 即 documentation + upstream-pr 索引；无独立 R 版本可 bump

## Implementation Sync

- Changed logic: 无运行时逻辑
- Authoritative current document: [PROJECT_STATUS.md](../../PROJECT_STATUS.md)
- Archive: 旧 `czz-docs` 路径 `archived-linked` 桩；内容 `git mv` 到 vibe
- CodeNote 登记：`vibe/knowledge/project-index.json` 的 `gonavi.routes` 改为 `AGENTS.md` / `vibe/rules/README.md` / `vibe/specs/PROJECT_STATUS.md`（仓库相对路径，无 Home）

## Out of scope

- 不开/不改 GitHub 上的 #1131/#1134
- 不把 `vibe/` 打进上游 PR
- 不改供应商 UI 代码
- 不在可提交文件写入 Home 绝对路径
- 不为 GoNavi 建 `vibe/ai-db/`
- 不把本仓文档提交进上游 PR（本 fork 本地提交另授）

## Verification

- 目录存在、索引互链、旧路径桩可读
- 仓库内 `czz-docs/*.md` 不再被当作 task card owner（桩文件除外）
- Skill「相关规则」不再指向 `czz-docs/` 的当前权威路径
- 第十六轮 r59 实机：2026-09-03 用户确认通过（记在收缩卡 / 总卡 / PR 03，不升本 spec 的产品范围）
