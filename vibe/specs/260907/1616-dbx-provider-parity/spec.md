# DBX AI 供应商交互复现

Tool: grok
Date: 2026-09-07
Status: `implementing`
Product: GoNavi AI 设置 · 模型供应商
对照：本机 fork `/Users/gdkmjd/work/czz/GitFork/dbx` 设置 → AI（非仓内代码）
过程 owner：[task-card.md](task-card.md)

## 交互模型

设置中心 `ai-providers` 详情使用 **list | edit**，与 DBX `aiConfigListMode` 对齐。

- 列表：标题「AI 配置列表」+「新增配置」；空态虚线框；卡片含品牌图标、名称、提供商、默认徽章、设为默认 / 编辑 / 删除。
- 编辑：左「返回」+ 标题「新增配置 / 编辑配置」；左标签表单；页脚左测试、右取消/应用。
- `ai-providers-connected` 只进入列表，不单独渲染芯片；树 key 保留。

## 提供商下拉

- 触发器：品牌图标 + 英文名。
- 下拉约 32rem、两列。左列「内置支持」= 现有 `PROVIDER_PRESETS`。右列「优质赞助商」空占位。
- 无 Jalapeño / HuaLong、无赞助徽章、无外链。

## 字段

API（拍平，不默认塞进认证折叠）：配置名称、提供商、认证（仅 Claude / Anthropic 兼容 / 自定义 messages 显示 api-key|bearer）、API Key、Endpoint、自定义请求头、API 格式（单选项只读 Input）、默认模型（保留 `n/m 已启用`）、最大输出 Token、上下文窗口（配置可存，调用层不强制改协议）。

本机 CLI：可选可执行路径 + 环境变量；空则 PATH 发现。不搬 DBX MCP 安装条。

## 非目标

Skip TLS、代理 URL、OpenCode / Pi / Qoder 新后端、合作方预设、`dbx://`、把设置中心整树改成 DBX 左栏、`vibe/` 进上游 PR。
