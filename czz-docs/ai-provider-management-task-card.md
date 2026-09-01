# 供应商展示与测试优化

状态：正式页面已落实预览验收后的管理交互和隐藏目录。2026-09-01 第十二轮合并上游 `dev` 后，修复窄工作区遮罩裸露、设置中心弹窗只能横向缩放、模型管理弹窗高度跳动，并压低了启用/默认开关的重渲染量。全量前端 5025 项中 5023 通过，两条失败分别为既有基线与上游自带。用户已在 r18 确认弹窗不再跳动；其余实机观感项仍待核验，不能把编译和离线测试当作这些项目通过。

工作分支为 `czz-dev`，实施起点为 `4cc7493c`。第十二轮把上游 `Syngnat/GoNavi` 的 `dev`（`89f9ad71`）合并进来，合并基为 `9ef337cd`，合并提交为 `44e5e3e3`；`origin/dev` 未参与。本轮获得重启、Computer Use 核验、优化、提交和推送授权；交付目标是新建 `origin/czz-dev`，不合并或覆盖 `origin/dev`。提交结果以 Git 历史和远端引用为准。

## 需求与当前行为

来源为 CodeNote `claude/gonavi-db-operations-bd76cd` 的 `vibe/specs/260828/0933-db-tool-operations/raw-requirement.md` 及后续确认的供应商管理需求。CodeNote 保留 DB 操作规范与 SRM 试点权威，产品实现和验收记录由本任务维护。

- 保留 AI 分类导航。已接入条目置顶，按原保存顺序排列，支持搜索、换行、一键收展、紧凑/正常密度和悬浮详情；勾选默认项只表示后端已确认的当前供应商。
- 新增使用可搜索下拉，不再前置选择端点。页内不再重复左侧「模型供应商」标题；说明与添加下拉同一行且左对齐，添加框只比文案略宽。供应商目录可折叠，默认三列，可拖宽、独立滚动并填满可用高度；查找框与「供应商目录」同一行，收起目录时隐藏。窄窗口改用抽屉，编辑表单随宽度组合字段。设置中心进入 AI 页时保持与其他设置页相同的 1080×820 外壳，不因进入 AI 而自行放大。
- 显示名称可选，复用现有 `name`。已配置 CLI 直接编辑和复用本机登录，从新增下拉排除；历史重复配置仍可编辑，不自动删除。普通 API 可多份接入。
- 目录候选可隐藏到“已隐藏”分组。隐藏后从正常目录和新增下拉排除，分组自动收起；展开后可恢复到原顺序。隐藏不删除配置、不改变默认项、不禁用模型，已接入条目仍可操作。
- 隐藏分组固定在目录底部，展开后用与上方相同的目录卡片（恢复图标），不再占满一行列表。搜索提示隐藏结果但不自动展开，全部隐藏时仍保留恢复入口。隐藏/恢复不触发保存、删除、检查或默认切换。
- 模型、档位保持可见，登录和技术说明折叠；新建或异常时展开。模型候选共用于默认、常用及补全控件，可搜索和手填；目录失败不覆盖原值。
- 模型下拉区分默认选择和启用管理。停用只过滤本配置的候选，保存后生效；当前默认模型和 SQL 补全模型须先替换才能停用。自定义候选不等于常用模型，不改变账号权限或已有会话执行策略。
- 普通 API 可“另存为”：复制当前草稿和可保留凭证，生成独立 ID/名称，不复用原密钥引用，不改变原配置或默认项。CLI 保持单例，不支持另存第二份。
- 切换直接生效，检查按需执行。连续切换串行并合并为最后目标，后端持久化成功后才更新勾选；保存、检查的忙碌状态独立。
- 草稿修改后离开需确认；编辑会话、配置版本及请求编号隔离旧结果。修改参数、切换预设、关闭编辑后旧检查不得回写；恢复原值也不会恢复旧检查成功。
- 供应商独立加载，MCP 客户端探测在进入 MCP 页面后才执行。CLI 默认值由后端提供，仅填入新建且未编辑的字段；切换 CLI 清除旧档位。

布局偏好（含隐藏候选）只存于本机 `gonavi.ai.providers.layout.v1`，不写入供应商配置或凭证；清除该偏好会恢复完整目录，不跨设备同步。独立交互预览及其服务未更新，不用预览结果证明正式产品可用性。

不包含自定义分组、批量检查、成本统计、自动故障切换或数据库治理变更。

## 配置与检查契约

配置沿用 `id/name` 及现有协议、区域地址、密钥留存逻辑，仅新增可选 `disabledModels/customModels`，Go JSON 使用 `omitempty`。旧配置缺省时行为不变；模型目录不会自动写成常用模型列表。

`AITestProvider` 保留 `success/message` 并增加 `checkKind/modelVerified`。前端不将缺少范围的旧桥接成功当作可信检查；浏览器模拟桥接明确返回检查不可用。

| checkKind | 执行内容 | modelVerified |
| --- | --- | --- |
| `none` | 未执行、未知协议或前置失败，不能成功 | `false` |
| `endpoint` | 原有 HTTP 端点检查，不承诺模型回复 | `false` |
| `local-auth` | Codex/Claude/Cursor CLI 本机登录检查，不发聊天 | `false` |
| `model-list` | Grok 模型列表检查，15 秒截止及管道等待上限，不发聊天 | `false` |
| `model-response` | 既有 Claude CLI 代理或 CodeBuddy 最小探测，必须读到非空回复 | 仅成功时为 `true` |

Grok 未安装、未登录、空输出、零退出拒绝、非零退出和超时均失败。当前供应商、检查通过、模型响应成功分别表达不同事实，不能合为“已连接”。

`AISetActiveProvider` 返回持久化错误，未知 ID 拒绝写入，失败恢复原内存状态；既有 `Promise<void>` 绑定可传递失败。`AISaveProvider` 在密钥写入前校验模型偏好和 CLI 唯一性，元数据写入失败恢复供应商列表。

`AIGetCLIModelCatalog` 返回 `{models, source, stale}`，来源为 `cache/cli/aliases/none`。`AIListCLIModels` 保留数组兼容接口。Wails 接口和配置模型声明仅增量同步，不重新生成整份绑定。

## CLI 与 NVM 边界

- 共享解析按 LookPath → nvm default/newest → Unix login-shell 执行，Windows 只走 LookPath。供应商执行、模型枚举与 MCP 探测复用这条链。
- 子进程保留 PATH 中已选的 Node，再补 CLI 目录和 nvm Node；保持符号链接的调用名称，避免破坏 shim。CodeBuddy 使用解析后的实际路径，不再退回丢失路径信息的命令名称。
- 不改系统 PATH、LaunchAgent、shim 或 CLI 登录配置。nvm 原生包、env-node 第二跳及架构回退由临时目录/假进程验证；这些不是 Windows/Linux 实机证明。
- Codex 只读模型缓存中的标识及可见性，尊重 `CODEX_HOME`；超过 24 小时、异常时间、8 MiB 上限、无效数据或读取失败时降级为手填。不读取认证、会话或模型指令。
- Claude 使用 `sonnet/opus/haiku` 常用别名，CLI 自行解析版本与账号权限；模型留空跟随 CLI 默认。2026-08-31 的接口核对来源为本机帮助及[官方模型配置](https://code.claude.com/docs/en/model-config#model-aliases)，不把别名清单当作账号可用模型全集。
- Cursor CLI 明确解析 `cursor-agent`，不回退到其他供应商的 `agent` 或编辑器 `cursor`；与原 Cursor Cloud API 独立。登录 JSON 必须明确已认证，枚举失败允许手填，不提供独立档位。接口依据为 2026-08-31 的本机版本核对及 [Cursor 参数文档](https://cursor.com/docs/cli/reference/parameters)。
- Cursor 请求使用临时工作目录、标准输入、Ask/sandbox 和项目工具拒绝规则，有限时及输出上限；保留用户/团队/企业原生 Hooks、插件和策略。**这不等同 Codex 的完整运行隔离**，限制提示始终可见；不伪造逐字流或 token 用量。
- Devin 仍待确认实际 CLI 来源/启动方式。此前核对到的是 Devin/Windsurf 桌面启动包装，未确认非交互模型、枚举及工具边界；[Devin 公共 API](https://docs.devin.ai/api-reference/overview)的云端任务不能代替本机 CLI。未添加伪可用预设或启动云端任务。

## 第十二轮交互与缩放修正

2026-09-01 实机反馈三项：供应商目录右侧出现一块约 221×423 的灰板；设置中心弹窗只能左右拉宽，底边拖不动；模型管理弹窗在切换「选择默认／启用管理」和设为默认时上下跳动。

灰板不是留白。工作区实测宽 633px，低于当时写死的 660 窄屏断点，目录因此变成绝对定位抽屉、只覆盖左侧 380px，`gonavi-ai-provider-scrim` 的 `#0003` 遮罩在右侧裸露成灰板；抽屉外壳只有 0.98 不透明度，底层编辑区的空态文案还会透出来。断点改为由「目录最小 216 + 手柄 17 + 编辑区最小 320」推导的 553，抽屉在窄屏下铺满工作区，并在抽屉打开时隐藏底层编辑区。

弹窗竖向缩放失效，是因为调用方通过 antd `styles.content` 传入行内 `height`，压过了样式表里改高度的那条规则；横向能拉只是因为宽度写在 `.ant-modal` 上不冲突。改为在 `[data-has-resized-height='true']` 下以 `!important` 取回高度，并让每个手柄只钉自己那根轴——拖南边不再顺手冻住宽度。

模型管理弹窗跳动来自两个选项卡内容高度不同，而弹窗向上弹出，高度一变就只能靠移动顶边吸收。改为两个选项卡共用一个预留高度的盒子；预留值由 `MODEL_MANAGEMENT_BODY_HEIGHT` 单点持有，经 CSS 变量下发，样式表不再重复这个数字。

启用/默认开关的延迟按重渲染量处理：模型行抽成 `React.memo` 组件，行回调用 ref 承接最新闭包以保持引用稳定，切换一个模型不再重建整份列表。悬浮提示统一为进场 300ms、离场 150ms；模型行原本的原生 `title` 一并换成 antd Tooltip，避免系统级的「约一秒才出、一移开就没」。

## 第十一轮布局压缩

2026-09-01 实机空态截图：已隐藏项各占一行，分组下方留下大块空白；页内「模型供应商」与左侧导航重复，添加下拉挤在两行标题旁。

修正仅在供应商管理组件和样式中：去掉页内标题，说明与添加下拉同一行；隐藏项改为图标+名称单行并贴在目录底部。供应商面板用 class 撑满高度，`[hidden]` 必须盖过 `display:flex`，避免切到安全控制/上下文时目录仍叠在上面。隐藏/恢复语义、布局偏好键和保存契约不变。

定向回归覆盖隐藏卡片、折叠时无恢复入口、顶栏不再输出重复标题。原生核验交给独立 r11 程序，由用户在设置页空态和展开已隐藏时确认。

## 第十轮问题与修正

Computer Use 在真实 r9 桌面发现：旧 Claude 配置将一个候选停用后再启用，模型集合已恢复原状，页面却继续显示“未保存”，离开时仍要求放弃修改。

原因是旧配置缺省的可选字段与表单后来产生的空数组比较不等。修正仅在[草稿比较](../frontend/src/utils/aiProviderManagement.ts#L54-L73)中把空 `disabledModels/customModels` 与缺省值视为等价，不改变表单或保存负载。非空模型变更、默认模型、密钥及有序常用列表仍会被识别；检查版本仍失效。

新增三项回归覆盖比较不修改输入、两类可选字段恢复后可直接离开、真实修改仍受保护以及旧检查结果不会复活。修正后原生补验受系统锁屏阻挡，尚未将这一步计为通过。

## 当前验证记录

| 项目 | 本轮结果 | 证据边界 |
| --- | --- | --- |
| 前端全量回归（r19） | 5023/5025，548 文件 | 两条失败为既有 `main.browserMock` 与上游自带的 `testPolicy` 基线漏登，合并前后均复现，非本轮引入 |
| 上游合并 | 通过 | `upstream/dev` `89f9ad71` 合入；3 处冲突（`App.tsx` import 取并集、生成的 `models.ts` 取上游、`package.json.md5` 按合并结果重算）已解，合并后 `upstream/dev...HEAD` 计数为 `0 11` |
| 前端定向回归（r11 布局） | 43/43，3 文件 | `AISettingsProvidersSection` 挂载/静态与侧栏；覆盖隐藏卡片、顶栏去重、添加下拉仍可搜索。不等于实机观感 |
| 前端定向回归 | 211/211，15 文件 | 覆盖旧响应、快速切换、桥接/持久化失败、草稿、密钥、复制、模型启停、隐藏恢复及 CLI 去重；模拟桥接不等于真实模型 |
| TypeScript | 通过 | 无类型错误 |
| Go CLI 竞态回归 | 通过 | 解析/缓存/枚举/认证/执行的临时进程与路径用例，未发真实聊天 |
| Go 供应商服务竞态回归 | 通过 | 检查范围、CLI 唯一性、模型偏好、旧密钥、持久化回滚、MCP 路径探测；临时配置和本地 HTTP |
| 六语言目录 | 通过 | 前端完整性及 Go i18n 检查 |
| 原生 Computer Use，r9 | 部分通过 | 真实 Wails 页面读到 3 个既有 CLI；Claude 别名列表、模型连续启停、原生 Escape 关闭弹层、真实草稿取消/放弃保护可操作 |
| 隐藏目录，r9 原生 | 已观察通过的子项 | 隐藏 OpenAI 后目录 21→20，分组默认收起；展开可见恢复入口，重新进入后隐藏保持且收起；既有 3 条配置、原默认项及编辑草稿保持 |
| 原生补验 | 待解锁 | r10 修正后的启停恢复、恢复显示点击、完整新增去重搜索、拖拽/窄窗口、真实 CLI 检查未完成；没有用其他 UI 通道绕过锁屏 |
| 构建与重启 | 通过 | 第十二轮逐轮停旧起新至 r19；Vite 生产包 + `wails build -s -skipbindings -nosyncgomod`。用户已确认弹窗不再跳动，其余观感项未确认 |
| 未执行 | 保留 | 真实模型回复、Windows/Linux 实机、冷/热加载测量；没有连接真实数据库 |

本轮 Computer Use 在 r9 暴露问题后，先放弃了测试产生的模型草稿，未保存模型修改；原默认供应商仍为 Claude。OpenAI 的临时隐藏在 r9 测试实例中尚待通过界面恢复，恢复入口已经确认存在。程序重启到 r10 后未能继续读取锁屏后的界面，不推断新实例的偏好状态。

本轮在 r9 快照基础上仅修改草稿比较、两个测试文件及本记录；其他既有源码变更保留。新增源码和待推送祖先提交未检出真实令牌、私钥或新增个人绝对路径；既有模拟路径保留。构建目录、截图、本机证据与独立预览不纳入提交。

验收产物为 `build/bin/GoNavi-provider-settings-260901-r19`，100,298,866 字节，SHA-256 `11e7b3027e69f2f05cae3d15834c85714e835712921a7be30381ae04fcc6b4ec`。原生核验使用独立 `GoNavi-Provider-Verification-r19.app`；未覆盖安装版，也不是签名发布包。r8–r18 及早期产物保留。两条核验包边界须记住：第十二轮中途误点应用内更新，r14 核验包被官方 0.9.4 发布包整包替换（独立二进制未受影响）；核验包的 `CFBundleDisplayName` 自 r16 起才随轮次同步，r14/r15 的菜单仍显示 r13，判断运行版本应以进程路径为准。

## 历史证据与决定

| 阶段 | 保留结论 | 不延用的结论 |
| --- | --- | --- |
| 最初治理与 50 项基线 | DB 例程禁止执行、SQL/应用/MCP 守卫边界保持 | 基线通过不证明后续 UI |
| 第一至三轮 | 直接切换、可信检查、分类导航、默认勾选、CLI 去重 | 曾尝试的分类导航下拉已按用户反馈撤回 |
| 第四、五轮 | Claude 别名与 Cursor CLI 接入；Devin 单独保留未完成 | 不宣称 Devin 已接入或模型可回复 |
| 第六轮 | 协议矩阵、MiniMax 区域保持、百炼旧 Chat 配置兼容 | 端点前置入口被预览验收方向取代，正式页不再挂载旧引导 |
| 第七轮 | 应用内 NVM 解析；历史本机登录/目录检查成功且未验证回复 | 当时旧密钥失败的推测归因被第八轮源码证据纠正 |
| 第八轮 | 已验收预览移入正式组件；修复密钥变量遮蔽及缺省 custom 协议；204 项前端与两组 Go 竞态通过 | r8 未启动，浏览器 DOM 键盘事件不算原生键盘 |
| 第九轮 | 隐藏/恢复目录；208 项前端、六语言及构建通过 | r9 当轮未启动，桌面结论来自本轮实际操作 |

第八、九轮曾在正式 App/组件上使用模拟 Wails 桥接，验证 API 另存为、原配置/默认保持、模型保存、深浅主题和 1440/1100/900/760 宽度；默认目录 336px/三列，拖宽范围 216–520px，窄窗口抽屉与独立滚动可用。第九轮还验证隐藏搜索、全部隐藏后恢复、原顺序及焦点保持。它们是模拟组件证据，不替代当前原生补验。

历史生成绑定有 337 处空白告警，本任务不做全文件格式化；当前任务增量另行检查。原生链接保留既有重复 `-lobjc` 警告，前端已有大包体及 Select 废弃属性提示不作为加载性能结果。

## 实现与回归入口

| 职责 | 入口 |
| --- | --- |
| 已接入、目录、隐藏/恢复、编辑 | [供应商组件](../frontend/src/components/ai/AISettingsProvidersSection.tsx#L280-L432)、[布局偏好](../frontend/src/components/ai/useAIProviderLayout.ts#L3-L89)、[响应样式](../frontend/src/components/ai/AISettingsProvidersSection.css#L1-L125) |
| 草稿、复制、保存和切换 | [编辑状态](../frontend/src/components/AISettingsModal.tsx#L417-L475)、[保存与默认切换](../frontend/src/components/AISettingsModal.tsx#L634-L743)、[宿主离开保护](../frontend/src/App.tsx#L3818-L3889) |
| 模型搜索及启停 | [模型控件](../frontend/src/components/ai/AIProviderModelSelect.tsx#L35-L171)、[模型偏好后端](../internal/ai/service/provider_models.go#L10-L49) |
| 预设、协议、别名 | [预设定义](../frontend/src/components/ai/aiSettingsModalConfig.tsx#L66-L109)、[端点映射](../frontend/src/utils/aiProviderEndpoints.ts#L1-L65)、[CLI 目录](../internal/ai/provider/cli_model_catalog.go#L17-L97) |
| 检查与 CLI | [检查范围](../internal/ai/service/service.go#L592-L766)、[当前供应商回滚](../internal/ai/service/service.go#L1144-L1174)、[NVM 解析](../internal/ai/provider/cli_lookup.go#L14-L113)、[Cursor CLI](../internal/ai/provider/cursor_cli.go#L21-L103) |
| 行为验证 | [编辑异步回归](../frontend/src/components/AISettingsModal.async.test.tsx#L1-L50)、[目录挂载回归](../frontend/src/components/ai/AISettingsProvidersSection.mounted.test.tsx#L1-L45)、[供应商服务回归](../internal/ai/service/provider_management_test.go#L20-L50)、[NVM 进程回归](../internal/ai/provider/codebuddy_nvm_test.go#L14-L42) |

## 下一步实机验收

1. 解锁 Mac 后继续使用独立 r10 程序，勿并行打开其他版本修改同一份供应商配置。确认实际运行版本，不用安装版或预览页代验收。
2. 将 Claude 的非默认模型停用后再启用，确认恢复原值时不提示未保存；真实修改仍能取消离开并保留草稿。确认检查提示不会随恢复原值复活。
3. 在隐藏分组恢复 OpenAI，再检查普通 API 和已接入 CLI 的隐藏、搜索、恢复顺序及新增去重；不需保存或删除真实配置。
4. 检查顶部收展/密度、目录拖拽、最小窗口和深浅主题，保持编辑区及底部按钮可用。
5. Codex/Claude/Cursor 使用本机登录检查，Grok 使用模型列表检查，核对范围提示。未安装、未登录、超时等故障用离线注入，不破坏真实登录或 PATH。
6. 真实模型回复需另行明确进行无数据库上下文的单条最小消息验证，会消耗订阅额度；当前未执行。Devin 需先确认实际 CLI 来源与契约。

## 文档与交付

本轮采用主线程执行；Computer Use 负责原生操作，git-batch-commit-push 负责范围审查与提交推送，doc-memory-closeout 和 document-code-link-audit 同步本记录。历史只读代理报告不代替本轮证据，没有新增代理、任务或 Goal。

`document_impact=project-current`，同步范围仅本任务记录；CodeNote 来源、DB 治理、全局规则与记忆不变。问题只记入本任务，不提升为全局记忆。提交按 CLI 路径、供应商接口/模型、正式 UI、验证说明划分；测试随对应行为提交。构建产物及本机核验证据不发布。

```json documentation-sync-group-v1
{
  "schema": "documentation-sync-group-v1",
  "group_key": "gonavi-ai-provider-management",
  "group_owner": "czz-docs/ai-provider-management-task-card.md",
  "documents": ["czz-docs/ai-provider-management-task-card.md"],
  "dependencies": [
    "frontend/src/components/AISettingsModal.tsx",
    "frontend/src/components/ai/AISettingsProvidersSection.tsx",
    "frontend/src/components/ai/AISettingsProvidersSection.css",
    "frontend/src/components/ai/AIProviderModelSelect.tsx",
    "frontend/src/components/ai/useAIProviderLayout.ts",
    "frontend/src/utils/aiProviderManagement.ts",
    "frontend/src/App.tsx",
    "internal/ai/provider/cli_lookup.go",
    "internal/ai/provider/cursor_cli.go",
    "internal/ai/provider/cli_model_catalog.go",
    "internal/ai/service/service.go",
    "internal/ai/service/provider_models.go"
  ],
  "validators": [
    "frontend/src/utils/aiProviderManagement.test.ts",
    "frontend/src/components/AISettingsModal.async.test.tsx",
    "frontend/src/components/ai/AISettingsProvidersSection.mounted.test.tsx",
    "internal/ai/provider/cli_lookup_test.go",
    "internal/ai/provider/cursor_cli_test.go",
    "internal/ai/provider/codebuddy_nvm_test.go",
    "internal/ai/service/provider_management_test.go",
    "internal/ai/service/provider_models_test.go"
  ],
  "git_scope_prefixes": ["czz-docs/ai-provider-management-task-card.md"]
}
```
