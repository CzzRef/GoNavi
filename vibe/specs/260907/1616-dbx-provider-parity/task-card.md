# Task Card：DBX AI 供应商交互复现

Tool: grok
Date: 2026-09-07
Status: `implementing`
Product: GoNavi AI 设置 · 模型供应商
Spec: [spec.md](spec.md)
Predecessor: [供应商编辑收缩](../../260903/0000-provider-editor-compact/task-card.md) · [界面约定](../../../knowledge/ai-provider-ui-conventions.md)

## Task Documentation Sync Group

- Group key: `dsg:gonavi:260907-dbx-provider-parity`
- Group owner: this `task-card.md`

```json documentation-sync-group-v1
{
  "schema": "documentation-sync-group-v1",
  "group_key": "dsg:gonavi:260907-dbx-provider-parity",
  "group_owner": "vibe/specs/260907/1616-dbx-provider-parity/task-card.md",
  "documents": [
    "vibe/specs/260907/1616-dbx-provider-parity/task-card.md",
    "vibe/specs/260907/1616-dbx-provider-parity/spec.md",
    "vibe/knowledge/ai-provider-ui-conventions.md",
    "vibe/specs/PROJECT_STATUS.md"
  ],
  "dependencies": [
    "frontend/src/components/ai/AISettingsProvidersSection.tsx",
    "frontend/src/components/ai/AIProviderLogo.tsx",
    "frontend/src/components/ai/AIProviderPresetSelect.tsx",
    "frontend/src/components/ai/AIProviderConfigList.tsx",
    "frontend/src/components/ai/aiSettingsModalConfig.tsx",
    "frontend/src/components/AISettingsModal.tsx"
  ],
  "validators": [
    "frontend/src/components/ai/AISettingsProvidersSection.test.tsx",
    "frontend/src/components/ai/AISettingsProvidersSection.mounted.test.tsx",
    "frontend/src/components/ai/AIProviderLogo.test.tsx",
    "frontend/src/components/ai/AIProviderPresetSelect.test.tsx",
    "frontend/src/App.tool-center.test.ts"
  ],
  "git_scope_prefixes": [
    "vibe/specs/260907/1616-dbx-provider-parity/",
    "vibe/knowledge/ai-provider-ui-conventions.md",
    "vibe/specs/PROJECT_STATUS.md",
    "frontend/src/components/ai/",
    "frontend/src/components/AISettingsModal.tsx",
    "frontend/public/icons/ai/",
    "frontend/src/types.ts",
    "internal/ai/types.go",
    "shared/i18n/"
  ]
}
```

```json worktree-task-v1
{
  "schema": "worktree-task/v1",
  "task_id": "260907-dbx-provider-parity",
  "control_plane": "app-root",
  "target_branch": "czz-dev",
  "repositories": [
    {
      "repo_id": "gonavi",
      "base_sha": "9f463ba67c4ca820b6a1f2a925f09329d3af845e",
      "worktree_branch": "codex/260907-dbx-provider-parity",
      "task_owner": "vibe/specs/260907/1616-dbx-provider-parity/task-card.md",
      "head": "9f463ba67c4ca820b6a1f2a925f09329d3af845e",
      "upstream": null
    }
  ],
  "commit_mode": "verified-milestone",
  "push_mode": "current-message-only",
  "verification_state": "planned",
  "push_state": "not-authorized",
  "integration_state": "not-started",
  "next_action": "implement list-edit shell, brand select, and flattened provider form"
}
```

## 范围

P0：列表/表单交互、品牌图标下拉（合作方空列）、拍平字段、可选 CLI 路径覆盖、测试/约定/核验包。

不在本卡：Skip TLS、代理、新 CLI 后端、合作方预设、上游 PR。
