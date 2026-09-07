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
    "frontend/src/utils/aiProviderKeyValue.test.ts",
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
    "frontend/src/utils/",
    "frontend/wailsjs/go/models.ts",
    "internal/ai/",
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
  "next_action": "await on-device r61 confirmation of list-edit, brand select, and flattened form"
}
```

## 范围

P0：列表/表单交互、品牌图标下拉（合作方空列）、拍平字段、可选 CLI 路径覆盖、测试/约定/核验包。

不在本卡：Skip TLS、代理、新 CLI 后端、合作方预设、上游 PR。

## 核验

离线（child worktree `codex/260907-dbx-provider-parity`，HEAD `fef06cf4` 之上未提交实现）：`npx tsc --noEmit` 退出码 0；相关 vitest 109/109；全量 vitest 5319 通过 / 29 失败（DataGrid / DriverManager / Snippet / 主题预设，与原目录同批基线，非本轮引入）；`go test ./internal/ai/provider ./internal/ai/service` 通过；`go build ./internal/... ./cmd/...` 通过。

r61 核验包：`build/bin/GoNavi-provider-settings-260907-r61`（68766338 字节，SHA-256 `807ff80dc954029485843967ac6f990c2ccc95ada96eea7711ecb6d082631bc5`），壳 `GoNavi-Provider-Verification-r61.app`，pid 65421 单实例。identifier `com.czz.gonavi.provider-verification.r61`（换 identifier 会清 WebView 布局偏好）。编译与离线测试不构成实机通过。

实机待确认：设置中心 → AI 配置列表 → 新增 → 品牌下拉两列滚动 → 测试/应用 → 返回。核验包在原目录 `build/bin/`（对照/验收面）；实现在 child worktree。
