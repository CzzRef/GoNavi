# G0 连通性矩阵 — 首轮执行结果

> 评估方案：[gonavi-agent-mcp-evaluation.md](./gonavi-agent-mcp-evaluation.md) §3.1。
> 执行时间：2026-08-28。执行环境：macOS (darwin/arm64)、Go 1.26.1，worktree 分支 `claude/gonavi-agent-mcp-evaluation-77148e`（版本串 `v0.8.9-0.20260828025834-9ef337cd81e4+dirty`）。
> 隔离方式：`GONAVI_DATA_ROOT` 指向 scratch 目录，未触碰真实用户数据与任何 Agent 全局配置；数据源为本地 SQLite（300 行种子表），无外部数据库依赖。

## 1. 结果矩阵

| 用例 | 结果 | 实测证据 |
|---|---|---|
| CON-01 Codex 接入 | ✅ 通过（范围见 §3） | `gonavi-mcp-server stdio` 与 `gonavi mcp stdio` 均暴露精确 12 工具（11 schema + execute_sql），serverInfo `gonavi-ai`；Codex installer 单测随 `go test ./internal/ai/service -run 'Install\|MCP'` 全绿（4.97s） |
| CON-02 Claude Code 接入 | ✅ 通过（范围见 §3) | 同上（stdio 契约同一入口；Claude Code installer 单测同批全绿） |
| CON-03 OpenCode / DeepSeek Harness / Grok Build | ✅ 通过（范围见 §3） | 三个 installer 单测同批全绿；stdio 契约同一入口 |
| CON-04 HTTP 认证与绑定 | ✅ 通过 | 正确 token→200 + `Mcp-Session-Id`；错误 token→401；无 token→401；`0.0.0.0` 无 `--allow-non-loopback`→启动即拒：`must bind to loopback …, got "0.0.0.0:18777"` |
| CON-05 schema-only 降级 | ✅ 通过 | `--schema-only` 服务器 tools/list 恰 11 个工具，`execute_sql` **缺席**（非报错）；全量服务器 12 个 |
| CON-06 行截断与钳制 | ✅ 通过 | 300 行表默认返回"仅显示前 50 行"+"结果已截断"；`maxRowsPerResult=500` 被钳到"仅显示前 200 行" |
| CON-07 取消传播 | ✅ 通过 | stdio 发起 2 亿次迭代递归 CTE，1.5s 后发 `notifications/cancelled`：**0.0s 内**返回结构化取消错误"SQLite 驱动代理 query 请求已取消：context canceled"，总耗时 1.51s，进程无悬挂 |
| CON-08 30 分钟会话超时 | ⏸ 顺延（部分） | 30 分钟空闲未实测；代码证据 `internal/mcpserver/run.go:181`（`SessionTimeout: 30 * time.Minute`）；代理测试：伪造 `Mcp-Session-Id` → 404（Agent 会话失效后可据此重新 initialize） |

补充实测（矩阵外）：

| 项 | 结果 |
|---|---|
| `remote-config` 生成 | 输出即贴即用的远程 Agent JSON 配置（streamable-http + Bearer），示例命令默认带 `--schema-only`（安全默认好），含 "OpenClaw" 云端 Agent 模板 |
| CLI `--request-trace` | stderr 输出完整结构化 trace（requestId/entry=cli/driverMode/事件时间线/cancellation/responseBytes），无头可观测性达标；本次 SQLite COUNT 查询 durationMs=67 |
| CLI JSONL 输出 | `result_set`/`row`/`summary` 分型行，stdout/stderr 分离，exit=0 |

## 2. 执行中的发现（进入评估记录）

**OBS-1【接入摩擦·中】SQLite 等 optional Go driver 无头启用无 CLI 通道。**
首次 `execute_sql` 报"SQLite 纯 Go 驱动未启用，请先在驱动管理中点击安装启用"。门控 = `<data-root>/drivers/<type>/installed.json` 标记 + 同目录 driver agent 可执行文件（`internal/db/driver_support.go:314-339,261-270`）。GUI 之外只能手工构建放置（本轮：`go build -tags gonavi_sqlite_driver ./cmd/optional-driver-agent`）。**建议：CLI 增加 `gonavi driver install <type>`**，否则纯无头/容器场景的 Agent 接入会卡在第一步。

**OBS-2【Agent 友好性·中】`execute_sql` MCP 返回 Markdown 文本，非结构化 JSON。**
行数、截断标志都嵌在中文文本里（"300 行，仅显示前 50 行"），Agent 需文本解析；与评估方案 GAP-10（结构化 risk/plan 输出）同向。**建议：content 增加结构化 JSON 块或提供输出模式参数**——执行面新工具应从第一天就结构化。

**OBS-3【隔离·低】日志不随 `GONAVI_DATA_ROOT`。**
错误提示指向 `~/.GoNavi/Logs/gonavi.log`，data root 已隔离而日志逃逸到真实用户目录。多实例/测试隔离时日志混流。**建议：日志目录跟随 data root 或提供独立 env。**

**OBS-4【小差异·低】CLI 与 MCP 的 SQLite 执行路径不同。**
CLI trace 报 `driverMode:"builtin"`（进程内执行），MCP 取消错误却来自"驱动代理"（agent IPC）。同一驱动两条运行路径，行为等价性（超时/取消/编码）值得在 G2 前确认一次。

**OBS-5【正向】版本可追溯、取消链路完整。**
serverInfo 版本串带 commit 与 `+dirty`；取消从 MCP 通知穿透到驱动代理进程并即时返回结构化错误——调研文档"Agent-friendly 错误契约"的判断得到实证。

## 3. 范围偏差声明

- CON-01/02/03 的"真宿主端到端"（真实 Codex/Claude Code 进程发现工具）未执行：不在本轮修改用户全局 Agent 配置（`~/.codex/config.toml` 等）。已覆盖：GoNavi 侧 stdio 协议契约（真宿主消费的就是它）+ 全部 5 个 installer 的单测。**残余风险低**；建议用户在 GUI 点一次安装后由任一宿主实测收尾。
- CON-08 的 30 分钟真实空闲等待未做（时间成本），以代码证据 + 伪会话行为代理；如需实测可用可配置超时或长挂测试补做。
- 延迟基线只做了抽样（CLI 67ms、取消响应 <50ms），未做系统性 P50/P95 采集；G2 对比前建议脚本化补齐（评估方案 OPT-10）。

## 4. 复现步骤（要点）

```text
1. go build ./cmd/gonavi ./cmd/gonavi-mcp-server
2. go build -tags gonavi_sqlite_driver -o <data-root>/drivers/sqlite/sqlite-driver-agent ./cmd/optional-driver-agent
   echo '{}' > <data-root>/drivers/sqlite/installed.json
3. sqlite3 g0.db 建 300 行表；GONAVI_DATA_ROOT=<data-root> gonavi connection add --name g0-sqlite --type sqlite --database g0.db
4. stdio：newline-delimited JSON-RPC（initialize → initialized → tools/list → tools/call）
5. HTTP：POST /mcp，Authorization: Bearer，捕获 Mcp-Session-Id 复用
6. 取消：tools/call 后发 notifications/cancelled(requestId)
```

探针脚本与原始输出存于本轮会话 scratchpad（`g0/mcp_stdio_probe.py`、`g0/results/*.json|raw|log`），未随文档归档。
