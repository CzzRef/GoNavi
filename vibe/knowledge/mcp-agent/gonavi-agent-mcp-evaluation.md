# GoNavi Agent 接入与执行层评估方案

> 评估对象：[gonavi-agent-mcp-deepResearch.md](./gonavi-agent-mcp-deepResearch.md) 提出的 "Execution Plane + MCP/CLI/HTTP 多适配器" 集成蓝图。
> 本文档回答三个问题：①现在通过 MCP 或 CLI 能直接用到什么程度；②蓝图落地前后各阶段如何评估验收；③调研设想哪些地方不齐全、哪些可以优化改造。
> 事实基线均在 2026-08-28 于分支 `claude/gonavi-agent-mcp-evaluation-77148e`（与 `dev` 同源）实测，证据见附录 A。

---

## 0. 结论速览

1. **数据面今天就可直接使用，且不止 Codex。** GoNavi 已内置 5 个 Agent 客户端的 MCP 一键安装器（Claude Code、Codex、OpenCode、DeepSeek Harness、Grok Build）并带安装状态检测；stdio 与 Streamable HTTP 双通道、`--schema-only` 降级、bearer token、行数截断都已就位。"其他 ID 的智能体连接"不是待建能力，而是待评估覆盖度的既有能力——调研文档只围绕 Codex 展开，这是它的第一个缺口。
2. **CLI 不是 MCP 的替代方案，而是同一内核的第二适配器 + 评估工具。** 现有 CLI 已有稳定 exit code 契约（0/2/3/4/5/6/7）与 JSONL 输出，可立刻作为无 MCP 能力的 Agent 的交互通道，也是 G1 阶段对比评估 MCP 效率的基准。
3. **执行面（project_run 等）尚不存在，调研蓝图方向正确**（Execution Core 单一真源 + 多适配器 + 不可变 Plan + digest 审批），但存在若干设想不齐全处，最关键的三个：**manifest 信任链未闭环（Agent 可自改审批档位）**、**MCP 工具同步性契约未定义（审批等待会撞上 Agent tool timeout）**、**桌面 GUI 生命周期与 durable run 的关系未定义**。详见 §4。
4. 评估采用 **G0→G3 四道门禁**：G0/G1 零开发、现在即可执行；G2 对应调研文档 PR A–E 的原型验收；G3 对应生产化。每道门有明确用例编号与通过标准（§3）。

---

## 1. 事实基线核验

调研文档的论断逐条对照源码核验。结论分三类：✅ 确认、🔧 修正、➕ 调研未覆盖的补充事实。

| # | 调研论断 | 代码证据 | 结论 |
|---|---|---|---|
| F1 | MCP 注册 11 个 schema 工具，`execute_sql` 由 SchemaOnly 门控（不注册而非拒绝） | `internal/mcpserver/server.go:13-14`（SchemaOnly 注释）、`server.go:89-90`（`if !options.SchemaOnly` 才 AddTool） | ✅ |
| F2 | `execute_sql` 默认每结果集 50 行、上限 200 行 | `internal/mcpserver/service.go:19-20`（`defaultMaxRowsPerResult = 50`、`maxRowsPerResultLimit = 200`） | ✅ |
| F3 | HTTP MCP 默认 `127.0.0.1:8765` `/mcp`，bearer token 必需，非 loopback 需显式 `--allow-non-loopback` | `internal/mcpserver/run.go:23-24`、`run.go:264`（bearerTokenAuthHandler + LimitRequestBody）、`run.go:298,324` | ✅ |
| F4 | CLI 稳定 exit code 0/2/3/4/5/6/7 | `internal/cli/cli.go:25-33`（ExitSuccess…ExitUnknownOutcome 常量块） | ✅ |
| F5 | CLI 子命令：mcp、list-connections、connection add/import、query、export、batch/exec-file、audit | `internal/cli/cli.go:143-188`（dispatch）、`cli.go:737-745`（mcp stdio/http/remote-config） | ✅ |
| F6 | Codex installer 写 `~/.codex/config.toml`，固定 server id `gonavi`，会修复自管旧配置、不动用户无关条目 | `internal/ai/service/claude_code_mcp.go:22,54,130,163-173` | ✅ |
| F7 | `codex_cli.go` provider 只是 credentialed model transport，不能当执行器放开 | `internal/ai/provider/codex_cli.go`（read-only sandbox、approval never 等声明） | ✅ |
| F8 | requesttrace 已有 RequestID/Cancellation/RetryCount 等结构化字段，容量有界、只存脱敏摘要 | `internal/requesttrace/store.go:36-84,235-275` | ✅ |
| F9 | "Codex 集成现状"仅描述 Codex 一键安装 | `internal/ai/service/` 实际有 **AIInstallClaudeCodeMCP**（`claude_code_mcp.go:95`）、**AIInstallCodexMCP**（`:130`）、**AIInstallOpenCodeMCP**（`open_code_mcp.go:27`）、**AIInstallDeepSeekHarnessMCP**（`deepseek_harness_mcp.go:31`）、**AIInstallGrokBuildMCP**（`grok_build_mcp.go:27`），以及统一状态 API `AIGetMCPClientInstallStatuses`（`claude_code_mcp.go:79`） | ➕ 调研严重低估了多 Agent 接入面的现成度 |
| F10 | Web Server 可作远程审批 UI 基础 | `internal/webserver/server.go:917`（`text/event-stream`）、`:555-566`（SSE 队列有界） | ✅ SSE 通道确已存在 |
| F11 | AI 路线图建议 tool registry / mcp server config / runtime bridge 三层拆分 | `AI_EXTENSIONS_ROADMAP.md` §2 | ✅ |
| F12 | 长任务先例 | `internal/syncjob/`、`internal/importjob/` 为既有任务域包 | ➕ 调研提到"可参考"但未评估其状态机可否抽取复用（见 GAP-12） |

**核验结论**：调研文档的事实层可信，可作为评估基线；其"现状缺口表"（缺 Project Catalog / Executor / DAG / Approval / 通用 Run）与代码一致——`internal/` 下不存在 execution/projectcatalog/approval/worker 任何一包。

---

## 2. 使用路径评估：MCP 直接用，还是给一个 CLI？

用户目标拆成两个时段回答：**现在（执行面未开发）** 与 **执行面落地后**。

### 2.1 路径清单与当前可用度

| 路径 | 现在可用？ | 现在能做什么 | 执行面落地后 |
|---|---|---|---|
| A. MCP stdio 直连（一键安装） | ✅ 5 类 Agent 均可 | schema 上下文 + 受控 SQL | + project_plan/run/status 等 6 工具 |
| B. MCP Streamable HTTP | ✅（token + loopback 默认） | 同上，适合容器/远程 Agent | 同上 |
| C. CLI 交互（Agent 经 Bash 调用） | ✅ | query/export/batch/audit，JSONL + exit code | + `project`/`run` 子命令 |
| D. Execution HTTP API + SSE | ✖ 不存在 | — | Code Note 等服务级调用主通路 |
| E. Codex App Server（反向：驱动 Codex） | ✖ 未集成 | — | Code Note 富会话驱动 Codex |

### 2.2 评估维度打分（1–5，5 优）

针对 "Agent 使用 GoNavi" 的三条正向路径：

| 维度 | A/B MCP | C CLI | D HTTP API |
|---|---:|---:|---:|
| Agent 接入成本（有 MCP 能力的宿主） | 5（一键安装） | 3（需写调用约定/skill） | 2 |
| 无 MCP 能力宿主的可达性 | 1 | 5 | 4 |
| Token/上下文效率 | 4（截断 + 结构化） | 3（stdout 需自行裁剪） | — |
| 长任务/异步语义 | 2（工具调用天然同步，见 GAP-06） | 3（`--no-wait` 可补） | 5（SSE） |
| 服务间（Code Note）契约稳定性 | 2 | 2 | 5 |
| CI 适配 | 2 | 5 | 4 |
| 安全边界表达（scope/审批） | 3 | 3 | 5 |

### 2.3 结论

**不做二选一。** 三条路径共享同一 Execution Core（调研文档已论证），分工是：

- **MCP（A/B）**：Agent 默认通路——现在即用于数据面，执行面落地后加 6 个 project 工具；
- **CLI（C）**：三重身份——①无 MCP 宿主与 CI 的正式通路；②G1 评估中与 MCP 做效率对照的基准；③执行面开发期间的**先行交付面**（`gonavi project run` 比 MCP 工具先可用、先可测，符合调研的 PR C→PR D 顺序）；
- **HTTP API（D）**：Code Note 服务级契约，不给 Agent 直用。

对用户"或者提供一个 CLI 来使用交互"的直接回答：**CLI 要提供，但作为同核适配器而非独立方案**；且在执行面未开发的窗口期，"Agent 经 Bash 调 CLI 白名单命令"就是可用的过渡执行通道（受 Claude Code/Codex 自身审批管控），不必等 Execution Plane。

---

## 3. 分阶段评估方案（G0–G3 门禁）

每道门 = 前置条件 + 用例集 + 采集指标 + 通过标准。未过门不进入下一阶段开发投入。

### 3.1 G0 — 现状连通性评估（零开发，本周可执行）

**目的**：把调研的"已具备"变成实测矩阵，同时为 G2 建立性能/成本基线。

| 用例 | 内容 | 通过标准 |
|---|---|---|
| CON-01 | Codex 一键安装 → stdio 发现全部 11+1 工具 | 工具清单与 `server.go` 注册一致 |
| CON-02 | Claude Code 同上 | 同上 |
| CON-03 | OpenCode / DeepSeek Harness / Grok Build 各安装一次并调用 `get_connections` | 安装状态 API 与实际一致；无互相覆盖 |
| CON-04 | HTTP 模式：正确 token 200、错误 token 401、无 `--allow-non-loopback` 绑外网失败 | 三项全部符合 |
| CON-05 | `--schema-only` 启动：`execute_sql` 在工具列表中**不存在**（而非报错） | 工具缺席 |
| CON-06 | 查询 >50 行表：默认截断到 50，`maxRowsPerResult=500` 被钳到 200，truncated 标志可见 | 与 `service.go:19-20,990-1008` 行为一致 |
| CON-07 | Agent 中断长查询：cancellation 传播、requesttrace 记录 outcome | trace 出现 forwarded/not_accepted 之一，无悬挂 |
| CON-08 | HTTP session 空闲 30 分钟后 Agent 重连行为 | 可恢复或明确报错，无静默失败 |

**采集指标**：每工具调用 P50/P95 延迟、单次调用返回字节数、每 Agent 每任务 token 消耗（作为 G2 对照基线）。

### 3.2 G1 — CLI 交互契约评估（零开发）

**目的**：验证 CLI 作为 Agent 通道的充分性，并产出"MCP vs CLI"对照数据，为执行面 CLI 子命令设计提供输入。

| 用例 | 内容 | 通过标准 |
|---|---|---|
| CLI-01 | `query --format jsonl` 输出被 Agent 稳定解析 100 次 | 0 解析失败；stdout/stderr 分离 |
| CLI-02 | 构造 7 种结局验证 exit code 映射 | 与 `cli.go:25-33` 一一对应 |
| CLI-03 | `batch/exec-file` 中途失败：语义（中止/继续/exit code）明确且文档一致 | 无歧义 |
| CLI-04 | secret 只经 env 注入，`ps`/audit log 中无值 | 零泄漏 |
| CLI-05 | 同一任务分别经 MCP 与 CLI 完成，对比 token 成本、轮次、错误可诊断性 | 产出对照报告（无通过线，供决策） |
| CLI-06 | Windows 下 CLI-01/02 复测 | 行为一致 |

**门禁判定**：CLI-01…04、06 全过 → CLI 契约可承载执行面扩展；CLI-05 报告归档为基线。

### 3.3 G2 — Execution Plane 原型门（对应调研 PR A–E）

**前置**：调研文档 PR A（manifest/plan）、PR B（worker/durable state）、PR C（CLI）、PR D（MCP 工具）、PR E（approval）合入。

**功能用例**（在调研各 PR 验收清单之上补充）：

| 用例 | 内容 | 通过标准 |
|---|---|---|
| EXE-01 | 同一 manifest 在 macOS 与 Windows 生成 planDigest | **逐字节一致**（要求 canonical 序列化，见 OPT-05） |
| EXE-02 | DAG 成环 manifest | plan 阶段拒绝，零进程执行 |
| EXE-03 | run 进行中杀掉 GoNavi 进程再重启 | run 状态可查，处于 UNKNOWN_OUTCOME 或恢复中，绝不静默丢失 |
| EXE-04 | UNKNOWN_OUTCOME 的 step | 不自动重试；CLI exit 7 |
| EXE-05 | 端到端：Codex 生成脚本 → plan → 一键审批 → 执行 → Codex 读 status/artifacts → 修复失败再来一轮 | 全程无人工拼命令；审批 UI 展示 digest/risk |
| EXE-06 | run 中途 cancel | step 进程组被终止（含子进程，Windows 用 Job Object 验证） |
| EXE-07 | 相同 idempotencyKey 重复 project_run | 去重，返回同一 run |
| EXE-08 | artifact 元数据经 MCP 返回、内容经 API 拉取 | MCP 响应体积 < 4KB |

**安全红队用例**（一票否决）：

| 用例 | 攻击 | 期望 |
|---|---|---|
| SEC-01 | step workdir 含 `../..` 或绝对路径逃逸 project root | plan 阶段拒绝 |
| SEC-02 | 审批后把被批脚本换成同路径恶意文件（symlink/覆盖） | worker 执行前 hash 复验失败，run 终止 |
| SEC-03 | **Agent 提交 commit 把 manifest 中 `approval: two_person` 改成 `auto`** | 见 GAP-02：manifest 档位降级必须触发人工确认，本用例必须被挡 |
| SEC-04 | ApprovalGrant nonce 重放 | 第二次拒绝 |
| SEC-05 | 过期 grant | 拒绝 |
| SEC-06 | step 输出中回显 secret 值 | 日志/MCP 响应/artifact 清单中被 redact |
| SEC-07 | 仓库文件中植入提示注入，诱导 Agent 调 `project_run deploy-prod` | Agent 无 approve 权限 → 卡在 WAITING_APPROVAL，审批 UI 明示来源 |
| SEC-08 | args 中含 shell 元字符（`;`、`&&`、`$()`） | 原样作为 argv 传递，不被解释 |
| SEC-09 | schema-only / ExecutionEnabled=false 构型 | 6 个 project 工具全部**不注册** |
| SEC-10 | 两个 Agent 并发 run 同一 `concurrency.group` 独占锁 workflow | 串行化，第二个排队或拒绝 |

**采集指标**：plan 延迟 P95、approval_wait_seconds、unknown_outcome_total、run 成功率、G0 基线对比的 Agent token 成本变化。

**门禁判定**：功能用例全过 + 安全用例零失败 + EXE-05 由至少两种 Agent（Codex、Claude Code）各走通一次。

### 3.4 G3 — 生产化门（对应 OIDC/RBAC、容器 Worker、K8s）

要点用例（详表在 G2 通过后按当时范围制定）：scope 矩阵逐格验证（重点：持 `runs.plan` 的身份调 approve 必须 403）、容器 Worker 逃逸测试（read-only root、no-new-privileges、egress deny-by-default）、control plane 与 worker 分离后 GUI 退出不影响 run、审计链路以单一 run_id 贯穿 Codex OTel → GoNavi → Worker、24h 浸泡测试无泄漏。

### 3.5 评估资产管理

- 用例编号（CON/CLI/EXE/SEC-\*）进入 `docs/` 下的执行面测试文档，结果表每次门禁评审归档一份快照；
- G0/G1 建议直接做成可重复脚本（Agent skill 或 make target），G2 起并入 CI。
- **G0 首轮执行结果（2026-08-28）**：[gonavi-agent-mcp-eval-G0-results.md](./gonavi-agent-mcp-eval-G0-results.md) —— 7/8 通过、CON-08 顺延，新增发现 OBS-1…5（驱动无头启用缺 CLI 通道、execute_sql 输出非结构化等）。

---

## 4. 调研设想的缺口补全

按风险排序。每项：缺口 → 后果 → 建议。

**GAP-01 多 Agent 接入面被低估（事实性缺口）**
调研通篇以 Codex 为唯一 Agent 样本，但代码已有 5 类客户端安装器与统一状态 API（F9）。后果：评估矩阵、身份模型、回归范围都会漏掉 4/5 的既有接入面。建议：G0 起所有用例跑五客户端矩阵；执行面的 per-client 能力档位（见 GAP-03）以 installer 的 client id 为天然键。

**GAP-02 manifest 信任链未闭环（最高风险设计缺口）**
调研用 `ExpectedManifestDigest` 防读取后替换（TOCTOU），但没有防"**Agent 正当地修改 manifest 本身**"：Codex 有仓库写权限，一次 commit 就能把 `deploy-prod` 的 `approval: two_person` 改成 `auto`，随后的 plan 完全合法。建议：GoNavi 本地库为每个项目 pin 一份 manifest digest（首次注册即信任锚点）；digest 变化时，**降低审批档位/放宽 network/新增 secret 引用的差异必须触发人工确认**，其余差异可自动接受并记审计。这是 SEC-03 的判定依据。

**GAP-03 单 bearer token 无身份区分，但不必等 OIDC**
调研把身份问题整体推给 OIDC/RBAC（G3）。中间态很便宜：installer 本来就按客户端分别写配置，可直接**每客户端发独立 token**，服务端维护 token→(client id, scopes) 映射。这样 G2 阶段就能实现"Codex 无 approve 权限"而无需 IdP。

**GAP-04 MCP 工具的同步性契约未定义**
`project_run` 若同步等待审批+执行，会撞上 Agent 宿主的 tool timeout（Codex 默认 startup 60s，工具调用也有上限），表现为 Agent 侧超时而 run 实际继续——制造 UNKNOWN_OUTCOME。建议：**所有执行类工具立即返回 runId + 当前状态**；另提供 `project_wait`（有界轮询，如最长 30s）供 Agent 显式等待。此契约应写进 PR D 的验收。

**GAP-05 桌面 GUI 生命周期 vs durable run**
调研的 durable state 解决"重启后状态还在"，但没回答"**用户关掉 Wails 应用时，正在跑的 run 怎么办**"。MVP 是单进程，GUI 退出即 worker 消失。建议：MVP 明确语义——退出时有活跃 run 则弹确认，强退后 run 标记 UNKNOWN_OUTCOME，重启做恢复扫描（EXE-03 覆盖）；中期提供 `gonavi execution-server` headless 常驻模式（复用现有 CLI/mcp-server 的 headless 先例），GUI 只作为它的审批/观察面。

**GAP-06 审批到达性与超时策略**
"一键审批"假设用户在屏幕前。缺：审批通知通道（webserver SSE 已可推、可加系统通知）、审批超时策略（建议 expire→auto-deny 并通知 Agent，绝不 auto-approve）、以及审批被拒时给 Agent 的结构化理由（便于其修正后重提）。

**GAP-07 Windows/跨平台进程语义缺失**
sandbox 一节全是 POSIX 假设。Windows 需要：Job Object 保证 cancel 杀进程树（EXE-06）、路径规范化处理 junction/8.3 短名、`.bat/.cmd` 不经 shell 解释的调用方式。GoNavi 是三平台桌面产品，这不是边缘情况。

**GAP-08 版本协商与契约演进**
三个面都缺版本策略：MCP 工具 schema（新增字段对旧 Agent 的兼容）、CLI JSONL 输出（建议每行带 `schemaVersion` 或在 header 行声明）、manifest `apiVersion`（v1alpha1→v1beta1 的迁移与拒绝规则）。建议在 PR A 的 ADR 里一并定义。

**GAP-09 资源配额与滥用防护**
Agent 失败重试循环可能无限触发 plan/run。缺：per-client 并发 run 上限、plan 速率限制、artifact 磁盘配额与保留期（GC 策略、单 artifact 大小上限）。建议 MVP 就给保守默认值（如每 client 并发 2、artifact 保留 7 天）。

**GAP-10 dry-run 与 plan 预览的信息设计**
调研有 `DryRun` 字段但没定义返回什么。Agent 自检需要结构化 risk 报告：将写哪些路径、访问哪些域名、引用哪些 secret、审批档位及原因。这也是审批 UI 的数据源，应作为 Plan 的一等输出而非 UI 私有拼装。

**GAP-11 Codex App Server 协议成熟度风险**
该协议仍在演进。建议：pin 协议版本 + 从官方 schema 生成类型 + 明确降级路径（App Server 不可用时 Code Note 退回 "MCP 状态轮询 + 无富 diff" 模式），并把此降级作为 Code Note adapter 的验收用例。

**GAP-12 与 syncjob/importjob 的关系未评估**
仓库已有两个长任务域。建议 PR A 之前做一次半天的评估：其状态机/进度上报/取消模式哪些可抽为共享 `jobcore`，避免第三套 job 语义；结论无论复用与否都写进 ADR。

**GAP-13 审计查询面割裂**
sqlaudit 与未来 executionaudit 若各一套导出，"一个 run_id 贯穿所有动作"就要用户自己拼。建议统一审计导出入口（CLI `audit export` 扩类型过滤），存储可分表。

**GAP-14 调研缺评估方法论**
调研只有实施路线与验收标准表，没有阶段化的评估/红队方案——本文档 §3 即补此缺口，应与调研文档互为配套：蓝图变更时同步更新门禁用例。

---

## 5. 优化改造建议（按优先级）

| # | 建议 | 理由 | 优先级 |
|---|---|---|---|
| OPT-01 | **MVP 再切一刀：先交付"顺序 workflow"（有序 step 列表），DAG 并行推后** | 调研 MVP 25–40 人日中 DAG/状态机占大头；顺序执行已满足 "Codex 脚本一键批准执行" 主场景，Walking Skeleton 可压到 ~12–15 人日先验证审批闭环价值 | P0 |
| OPT-02 | 执行类能力用**构型开关不注册**（`ExecutionEnabled=false` → 工具/子命令整体缺席），复制 SchemaOnly 模式 | 已被 `server.go:89` 验证的最强隔离手段；SEC-09 依赖它 | P0 |
| OPT-03 | per-client token + scope 先行（GAP-03），OIDC 后置到 G3 | G2 即可测"Agent 不能自批"，成本一两天 | P0 |
| OPT-04 | MCP 执行工具全部异步化 + `project_wait`（GAP-04） | 避免制度性 UNKNOWN_OUTCOME | P0 |
| OPT-05 | planDigest 用 canonical 序列化（如 JCS/确定性 CBOR），禁止直接 hash `encoding/json` 默认输出 | EXE-01 跨平台一致性的前提；map 序 & 路径分隔符都会咬人 | P0 |
| OPT-06 | manifest 信任锚点 + 档位降级审批（GAP-02） | 补最大安全洞 | P0 |
| OPT-07 | run/step trace 直接扩展 `requesttrace` schema 而非新建关联体系 | `store.go` 已有容量上界与脱敏纪律，复用即继承 | P1 |
| OPT-08 | 审批事件复用 webserver SSE 通道推送，MVP 不新建事件总线 | `server.go:917` 通道现成，有界队列语义已处理 | P1 |
| OPT-09 | CLI 执行子命令提供 `--wait/--no-wait/--timeout`，exit code 复用现表 | CI 需要阻塞语义，Agent 需要非阻塞语义，一套命令两用 | P1 |
| OPT-10 | 评估基线脚本化：G0/G1 矩阵做成可重复 skill/make target，结果 JSONL 归档 | 每次加固（近期 PR 频繁）后可回归，防契约漂移 | P1 |
| OPT-11 | 审批默认带 TTL，过期 auto-deny + 通知（GAP-06） | 防"周五发起周一还挂着"的僵尸授权 | P1 |
| OPT-12 | 文档漂移清理：以 `service.go:19-20`（50/200）为真源同步各 MCP README 表述 | 调研已点名，顺手关闭 | P2 |
| OPT-13 | jobcore 抽取评估（GAP-12）作为 PR A 前置半日任务 | 花半天避免三套 job 语义 | P2 |

---

## 6. 风险与回退

| 风险 | 触发信号 | 回退路径 |
|---|---|---|
| 执行面工期超支 | G2 门禁两次未过 | 停在 OPT-01 的顺序执行 MVP；Agent 继续走 "CLI 白名单 + 宿主自身审批" 过渡通道 |
| Codex App Server 协议破坏性变更 | 类型生成失败/握手不兼容 | GAP-11 的 MCP 轮询降级模式 |
| MCP 工具契约变更影响存量 Agent | G0 回归矩阵失败 | 工具只增不改语义；破坏性变更走新工具名 |
| 审批疲劳导致用户放宽档位 | approval_wait 指标持续走高 | 用指标驱动把高频低风险 workflow 移入 `auto` 白名单，而非全局放宽 |

---

## 附录 A — 证据清单

| 证据 | 位置 |
|---|---|
| MCP 工具注册与 SchemaOnly 门控 | `internal/mcpserver/server.go:13-14,34-90` |
| 行数截断 50/200 | `internal/mcpserver/service.go:19-20,990-1008` |
| HTTP 默认地址/token/loopback | `internal/mcpserver/run.go:23-24,264,289-324` |
| CLI exit code 契约 | `internal/cli/cli.go:25-33` |
| CLI 子命令 dispatch | `internal/cli/cli.go:143-188,737-745` |
| 五客户端安装器与状态 API | `internal/ai/service/claude_code_mcp.go:79,95,130`、`open_code_mcp.go:27`、`deepseek_harness_mcp.go:31`、`grok_build_mcp.go:27` |
| Codex config.toml 路径与固定 id | `internal/ai/service/claude_code_mcp.go:22,54` |
| requesttrace 结构化字段 | `internal/requesttrace/store.go:36-84,235-275` |
| webserver SSE | `internal/webserver/server.go:555-566,917` |
| AI 扩展路线三层拆分 | `AI_EXTENSIONS_ROADMAP.md` §2 |

## 附录 B — 用例编号索引

- CON-01…08：G0 连通性（§3.1）
- CLI-01…06：G1 CLI 契约（§3.2)
- EXE-01…08：G2 功能（§3.3）
- SEC-01…10：G2 安全红队（§3.3）
