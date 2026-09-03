# Task Card：供应商编辑收缩、布局提示与 CLI 提示范围

Tool: grok / DSH
Date: 2026-09-03
Status: `verified-r59`（2026-09-03 用户在 r59 核验包上确认通过）
Product: GoNavi AI 设置 · 模型供应商
Predecessor: [第十五轮认证字段默认纵向](../../260901/0000-ai-provider-management/task-card.md#第十五轮认证字段默认纵向) · [界面约定](../../../knowledge/ai-provider-ui-conventions.md)
上游记录: [03-260903-squash-编辑收缩展示](../../../knowledge/upstream-pr/03-260903-squash-编辑收缩展示.md)

> 当前路径：`vibe/specs/260903/0000-provider-editor-compact/task-card.md`。旧 `czz-docs/260903-provider-editor-compact-task-card.md` 只是归档桩。

## Task Documentation Sync Group

- Group key: `dsg:gonavi:260903-provider-editor-compact`
- Group owner: this `task-card.md`
- Predecessor: `gonavi-ai-provider-management`（落点已拆）

```json documentation-sync-group-v1
{
  "schema": "documentation-sync-group-v1",
  "group_key": "dsg:gonavi:260903-provider-editor-compact",
  "group_owner": "vibe/specs/260903/0000-provider-editor-compact/task-card.md",
  "documents": [
    "vibe/specs/260903/0000-provider-editor-compact/task-card.md",
    "vibe/knowledge/ai-provider-ui-conventions.md",
    "vibe/knowledge/upstream-pr/03-260903-squash-编辑收缩展示.md"
  ],
  "dependencies": [
    "frontend/src/components/ai/AISettingsProvidersSection.tsx",
    "frontend/src/components/ai/AIProviderModelSelect.tsx",
    "frontend/src/components/ai/AICompactValueInput.tsx",
    "frontend/src/components/ai/AIHintChoiceControl.tsx",
    "frontend/src/components/ai/AIProviderSortable.tsx",
    "frontend/src/components/ai/cliModelCatalogCache.ts",
    "frontend/src/components/ai/compactConnectionDisplay.ts",
    "frontend/src/components/ai/providerPresetOrder.ts",
    "frontend/src/components/common/tooltipTiming.ts"
  ],
  "validators": [
    "frontend/src/components/ai/AIProviderModelSelect.test.tsx",
    "frontend/src/components/ai/AISettingsProvidersSection.mounted.test.tsx",
    "frontend/src/components/ai/cliModelCatalogCache.test.ts",
    "frontend/src/components/ai/compactConnectionDisplay.test.ts",
    "frontend/src/components/ai/providerPresetOrder.test.ts"
  ],
  "git_scope_prefixes": [
    "vibe/specs/260903/0000-provider-editor-compact/",
    "vibe/knowledge/ai-provider-ui-conventions.md",
    "vibe/knowledge/upstream-pr/03-260903-squash-编辑收缩展示.md"
  ]
}
```

> 用户要求先核验说明，再写 task card 供本人核验。通过前不实施。

## 1. 核验对照（现状 vs 你的五条）

对照源码：`frontend/src/components/ai/AISettingsProvidersSection.tsx`、`AIProviderModelSelect.tsx`、[界面约定](../../../knowledge/ai-provider-ui-conventions.md) §A/B/D/G。

| # | 你的说明 | 现状（核验） | 判断 |
| --- | --- | --- | --- |
| 1 | 收缩时 URL 省略 `https://`；API 格式/Key 头尾 + 中间省略；够宽时格式+URL+Key 同一行 | 并排已存在（`connectionLayout=inline`），但字段是完整值、有最小宽 320/220，不够就换行，**没有**协议省略和头尾压缩 | 合理；是第十五轮「看全」之上的**收缩展示**，不是改存储值 |
| 2 | 布局开关两钮要有框/色和选中态；说明移到钮右侧问号；保留「上说明、下两钮、右问号」 | 现在 ⓘ 浮层里：多行说明 + 两个几乎无框的字按钮叠在一起，选中靠 `aria-pressed`，观感弱 | 合理；是现有 `interactiveHintTooltip` 的视觉升级，不是换交互模型 |
| 3 | 供应商目录选中后若 API 格式没有其它候选项，不要下拉 | 认证区始终渲染 `<Select>`，哪怕 `getProviderEndpointTypes(preset)` 只有 1 项 | 合理；单选项应改成只读文本 |
| 4 | `1/1 已启用` 旁加情景图标；底栏说明收到图标里；模型少时不要锁死高度；弹层只留「启用管理」，默认模型留在最上的选择器 | 弹层固定 `MODEL_MANAGEMENT_BODY_HEIGHT = 280`；头上有「选择默认 / 启用管理」双页签；底栏常驻 `models.scope` 长文 | 合理；与第十五轮「说明不占整行、ⓘ 收说明」一致。**「选择默认」页签可删**，因为表单第一行已经是默认模型选择器 |
| 5 | 「已接入 CLI 无需重复添加」只在：选 CLI 类供应商 **且** 正在编辑该类型已接入项时出现 | 目录工具条**无条件**渲染 `catalog_hint`（任何供应商编辑都看见） | 合理；**位置保持工具条原位**，平时隐藏，以后别的供应商可换别的文案 |
| 6 | 本机 CLI 折叠若没有内容：不要收缩，只留对齐的提示图标+短文；甚至可整块移除，把提示挂到「Codex 订阅」这类已接入标题右上角 | `usesLocalCLI` 时认证字段不进 `<details>`，折叠栏几乎只有标题「本机 CLI …」+ ⓘ + 空 caret | 合理；空壳折叠应取消。挂到已接入芯片右上角，或改成无箭头的一行短文+ⓘ |

## 2. 我对原话的归一化（请你改错）

1. **收缩展示只改显示，不改保存值。** 输入框/粘贴仍是完整 URL、完整 Key、完整 apiFormat。省略只发生在「并排/收缩」展示层；悬浮看全文。
2. **省略规则（已拍板）：** 格式与 Key 都用 `...`。URL 去掉开头 `https://` / `http://`，过长再头尾省略。短于阈值不省略。
3. **三字段一行的条件：** 仅 `connectionLayout=inline`，且编辑区宽度同时满足 URL≥320 与 Key≥220（沿用第十五轮下限）。不够则换行，**绝不**再压窄截断中间看不见。
4. **布局 ⓘ 结构：** 浮层上方仍是「默认纵向、够宽可选并排」这段说明；下方是「纵向 | 并排展示」两个分段按钮（有描边、有选中底）；右侧一个 `?`，再悬浮只解释这两个按钮干什么。`?` 走不截鼠标的短提示；整块 ⓘ 仍走可点击的 `interactiveHintTooltip`。
5. **默认模型弹层：** 外层选择器继续选默认模型。点 `n/m 已启用` 只打开「启用管理」（搜索、开关、添加）。不再要「选择默认」页签。`n/m 已启用` 右侧加 ⓘ，内容 = 现在底栏那两句（来源 + 仅影响本配置）。候选项少时弹层高度随内容，取消 280px 预留。
6. **CLI 工具条提示：** 位置仍在目录工具条原处。仅在编辑**已保存的 CLI 供应商**时显示「已接入 CLI 无需重复添加」；平时不显示。以后其它供应商可在同一位置换其它文案。卡片绿勾不受本条影响。
7. **空的「本机 CLI」折叠：** 认证字段对 CLI 本来就不在这个 `<details>` 里，折叠是空壳。取消 caret。待选：把 ⓘ 挂到当前已接入芯片（Codex 订阅 / Grok 订阅 / Cursor CLI）右上角，或改成表内一行短文+ⓘ、无折叠。

## 3. 已拍板 / 仍待选

已拍板：

- **B-1 单选 →** 格式和 Key 都用 `...`。
- **B-2 单选 →** 仍可编辑，只压缩展示。
- **B-3 单选 →** 提示留在目录工具条原位；平时不显示；以后别的供应商可换文案。

已落定（r54 实现、r55 确认）：

- **空 CLI 栏 → `F-1-a`**：删掉折叠，ⓘ 挂到正在编辑的已接入芯片右上角（`gonavi-ai-provider-chip-hint`）。
- **公共件 → 介于 `F-2-a` 与 `F-2-b`**：抽出 `AIHintChoiceControl`（分段钮 + ?）、`AICompactValueInput`（聚焦出全文的压缩输入）与 `compactConnectionDisplay.ts`（省略函数）；计数旁 ⓘ 沿用 `hintIcon`，未单独抽件。
- **来源说明只留左 ⓘ**（2026-09-03 拍板）：右 ⓘ 不再重复 `modelSourceKey`。
- **原空折叠里的 CLI 三条并入右 ⓘ**（2026-09-03 拍板）：CLI 命令、能力读取失败、effort 未校验，与登录说明、刷新说明同在 `n/n 已启用` 右侧 ⓘ。

## 4. Goal And Scope

- Goal: 并排收缩时看清格式/URL/Key 的头尾；布局开关可辨认；单选项不下拉；模型管理弹层变矮、只管家；CLI 提示不再误伤普通供应商。
- In scope: 供应商编辑页连接字段、布局 ⓘ、默认模型弹层、目录 CLI 提示；[界面约定](../../../knowledge/ai-provider-ui-conventions.md) 增补；CodeNote 侧给「分段选择 + 旁置问号」一条偏好指针（不复制像素）。
- Out of scope: 改保存协议、改 CLI 登录、改目录绿勾、改设置中心外壳、本轮不提交/不推送除非另授。
- Success evidence: 你按第 7 节清单点过；相关单测按影响面补上。

## 5. Decision（待你核验后才生效）

- Documentation level: `standard`（产品 UI 增量；规则只补约定，不升 Controlled）
- Execution: `main-only`
- High-risk / DB / 凭证: Key 只改展示层，完整值仍在输入/存储；悬浮全文不得进日志。
- Plan-mode: 本卡通过前不改代码。

## 6. 建议改动面

| 面 | 拟议 |
| --- | --- |
| 连接三字段 | `inline` 时显示压缩（`...`、去协议）；仍可编辑；宽度够则一行；`stacked` 仍完整单列 |
| 布局 ⓘ | 上说明、下描边分段钮、右 `?` |
| API 格式 | `options.length <= 1` 改为静态文本，不出现 Select 箭头 |
| 模型弹层 | 去掉 select 页签；高度自适应；底栏说明改到 `n/m 已启用` 旁 ⓘ |
| 目录提示 | 原位；仅编辑已接入 CLI 时显示 |
| 空 CLI 折叠 | 取消 caret；挂片或单行（待选） |
| 文档 | 本卡通过后写入 [界面约定](../../../knowledge/ai-provider-ui-conventions.md) 新节；[管理总卡](../../260901/0000-ai-provider-management/task-card.md) 加第十六轮指针；CodeNote `contextual-help` 只加「分段选择旁置问号」邻居句 |

## 7. 验收清单

2026-09-03 用户在 r59 上确认下列项通过。

1. 并排且够宽：格式、去协议 URL、压缩 Key 同一行；悬浮能看到完整值。
2. 并排但不够宽：换行，URL/Key 不被中间截断成看不懂。
3. 纵向：仍是完整字段，不做强制压缩。
4. 布局两钮有框/色，选中可辨；`?` 只解释按钮。
5. 目录里仅一种 API 格式的供应商：编辑区看不到下拉箭头。
6. 点 `n/m 已启用` 只有启用管理；默认模型仍用上方选择器。
7. 1～2 个模型时弹层明显矮于现在的 280px 空档。
8. 编辑已接入 Cursor CLI：可见「已接入 CLI 无需重复添加」。编辑 OpenAI / 新建普通 API：不可见。
9. CLI 编辑页不再出现空的「本机 CLI … >」折叠条；提示挂已接入芯片右上角。

## 8. Prior Task Overlap

- Relationship: `continuation` of 第十五轮（并排从强制改成可选）与第十三轮（ⓘ 收说明、模型双页签）。
- Traceability: `delta-only`；不回退「默认纵向、URL 必须能看全」。
- 本卡取代管理总卡里「模型下拉区分默认选择和启用管理」中**弹层双页签**那一句；外层默认模型选择器保留。

## 9. 明确没做（等你核验）

- 历史段。以 §10 融合说明为准。

## 10. 融合说明（当前理解）

原则：**提示进 ⓘ / ?，不进标题、不进列表行，不准把一行挤换行。** 保存值永远是全文；压缩只发生在并排展示。

### 10.1 收起编辑左侧 ⓘ

- 正文只留：「可搜索、自定义模型，并在下拉中管理启用状态。」
- 其下是分段钮：**并排展示 | 纵向**，样式对齐「紧凑 | 正常」，点选立刻换选中底，字段布局跟着变。
- **「够宽时格式、URL、密钥可放同一行。」不出现在 ⓘ 正文**，只进右侧 **?**。? 还要说明两钮含义（并排=同一行压缩展示，纵向=每项一行看全）。
- 不要写「默认纵向」。

### 10.2 认证三字段

- 纵向：完整 URL / Key，不强制压缩。
- 并排且够宽：格式 + 去协议 URL + 压缩 Key 同一行；仍可编辑，聚焦出全文。不够宽则换行，禁止中间截断。
- API 格式只有一项：不要 Select/假输入框，用固定示意字。

### 10.3 认证与连接 ⓘ

- 只讲：填写格式、接口地址和密钥。
- 禁止把模型候选说明挂在这里。

### 10.4 默认模型这一行

- 标题就是「默认模型」，**一行，禁止塞「留空跟随 CLI 默认」这类长文**。
- 左侧 ⓘ：CLI 时悬浮「选择模型；留空跟随 CLI 默认」+ 来源/范围。来源只在这里出现。
- 右侧：`n/n 已启用` + 另一 ⓘ，**仅 CLI 供应商渲染**（CLI 登录说明、CLI 命令、effort 未校验、能力读取失败、可点计数刷新且刷新完成前保留当前列表）。非 CLI 供应商右侧只有计数。
- 点 `n/n`：打开启用管理，并刷新候选（内存保留旧列表直到新结果）。

### 10.5 启用管理弹层

- 只要「启用管理」，不要「选择默认」页签。打开不抢搜索焦点。搜索栏保留。
- 列表不要「选择模型；留空跟随 CLI 默认」行。
- 每行单行、矮行：名称 | 非默认则「设为默认」 | 启用开关。
- 当前默认：标「默认」，开关在但不能关，点关提示至少留一个。
- 其他已启用项可点「设为默认」。

### 10.6 CLI / 目录

- 空的「本机 CLI」折叠删掉。
- 「已接入 CLI 无需重复添加」只在编辑已保存 CLI 时出现在目录工具条原位。
- CLI 芯片右上角可挂 ⓘ。

### 10.7 回复规则（CodeNote 全局，已进 owner）

- 选项用嵌套 `F-1` / `F-1-a`，标明单选/多选。
- 表格只给并列对照（核验），不把决策组摊平。
- 模块之间空行，Markdown 列表可见层级。

## 11. r55 核验与修正（2026-09-03）

按 §10 对照 r54 工作区改动，6 处不贴合、4 处遗留，已全部修正；文件行号为修正后实测。

| # | 不贴合 | 修正 |
| --- | --- | --- |
| 1 | §10.1「够宽时…同一行」出现在 ⓘ 正文（`AIHintChoiceControl` 的 `description`）及无障碍文本 | 句子并入 `connection_layout.choice_hint`，只在 `?` 里；删 `connection_layout.hint` 键（6 locale）；`description` 改为可选并不再传 |
| 2 | §10.4 默认模型标题旁塞「选择模型；留空跟随 CLI 默认」（`gonavi-ai-provider-model-aside`） | 删该 span 与 CSS；句子只在左 ⓘ |
| 3 | §10.4 左右 ⓘ 都含来源；原空折叠里的 CLI 命令/能力失败/effort 未校验三条无处显示 | 右 ⓘ 仅 CLI 渲染，去来源，并入三条 + 刷新说明 |
| 4 | §10.5 点选择器本体打开普通菜单时，顶部仍挂「启用管理 \| n/m \| ×」 | `renderManagement` 在 `mode !== 'manage'` 时直接返回 `menu`；外点关闭监听只在 manage 模式挂载 |
| 5 | §10.5 `is-default` 底色按 `badge` 非空判定，SQL 补全行也被高亮 | `ModelManagementRow` 增 `isDefault` prop |
| 6 | §10.2 并排 URL 的 `fullHint={undefined}`，悬浮看不到全文 | 传 `Form.useWatch('baseUrl')`；Key 仍不悬浮 |
| 7 | 模型目录回退不判 scope：切供应商时先显示上一个 CLI 的列表，拉取失败则永久残留；effect 开头少了 `setModelsLoading(false)` | `catalogResponse` 记 `cliScope`，只在同一编辑会话保留旧列表；失败时清异 scope；恢复 loading 复位 |
| 8 | `ModelSelectionManagement` 五个 effort prop 从未渲染；`.gonavi-ai-model-effort` 等孤儿 CSS；`QuestionCircleOutlined`、`hintChoiceInteractiveTooltip` 未用 | 全部删除 |
| 9 | `ai_settings.models.select` 键已无引用 | 删（6 locale） |
| 10 | 弹层高度注释仍说「预留 280 让两页签不跳动」 | 改写为「上限 280、短列表随内容」 |

另：`tooltipTiming.test.ts` 对 `interactiveHintTooltip` 的深比较缺 `getPopupContainer`，r54 起即失败，已补 `expect.any(Function)`。

r55 离线结果：`tsc --noEmit` 通过；全量 `vitest run` 5144 项中 5142 通过，两条失败为已知基线（`main.browserMock`、`testPolicy`），非本轮引入。

r55 核验包：`build/bin/GoNavi-provider-settings-260903-r55`（100547186 字节，SHA-256 `8633caaed2004d397e9195fb42d816f270f3e799a8965bf2b2d052d0b15d671d`），壳 `GoNavi-Provider-Verification-r55.app`。已被 r56 替换。

## 12. r56 实机反馈修正（2026-09-03）

r55 实机回看给出四条（三张截图 + 一条拖拽需求），r56 逐条落地：

| # | 反馈 | 修正 |
| --- | --- | --- |
| 1 | 图1：单选项 API 格式变成了裸文本，边框没了 | `gonavi-ai-provider-fixed-value` 改为 antd `Input readOnly tabIndex=-1`，边框与 URL/Key 对齐，文字 `--provider-muted` 浅灰，无箭头不可改 |
| 2 | 图2：CLI/订阅模型每次进编辑都实时拉取 | 新增 `cliModelCatalogCache.ts`：可用结果按 `apiFormat` 存 `gonavi.ai.providers.modelCatalog.v1`；进入编辑先读缓存不调 `AIGetCLIModelCatalog`；只有点绿色 `n/m 已启用` 置 `catalogRefreshForced` 才重拉并覆写缓存；`stale`/空/失败不入缓存。`models.refresh_hint` 六语改为「平时沿用上次结果」口径 |
| 3a | 图3：`?` 第一句「够宽时…」多余 | 删除；`choice_hint` 键去掉 |
| 3b | 图3：说明要自动换行、分并排 / 纵向两条 | 新键 `connection_layout.inline_hint` / `stacked_hint`，`choiceHint` 渲染为 `gonavi-ai-provider-hint-body` 内两行 |
| 3c | 图3：钮文两字、观感对齐顶部紧凑/正常 | `inline` 六语改「并排」/「Inline」等两字；删 CSS 顶部 16 行浮层内 `density` 覆盖（边框底色 `!important`、加粗、8px），直接复用顶部 `gonavi-ai-provider-density` 规则 |
| 4 | 目录与已隐藏供应商要能拖拽排序，多列时虚化 + 提前占位 | `presetOrder` 入 `layout.v1`；`providerPresetOrder.ts`（`applyPresetOrder` / `movePresetWithinGroup`）；`AIProviderSortable.tsx` 用 `@dnd-kit` `PointerSensor(distance 6)` + `rectSortingStrategy`（目录网格）/ `verticalListSortingStrategy`（已隐藏），原卡片 `is-drag-placeholder` 虚化虚线并随 transform 提前滑到落点，`DragOverlay` 浮起副本；搜索中禁拖 |

第 4 条按「供应商目录 + 已隐藏」理解（「调整宽度导致并排多个」指目录网格多列）。已接入芯片行（Codex 订阅 / Grok 订阅…）的顺序来自后端 `providers` 数组，本轮未做拖拽；若也要，需另加 `providerOrder` 偏好或后端排序接口。

新增测试：`providerPresetOrder.test.ts`（4）、`cliModelCatalogCache.test.ts`（3）、mounted 两条（缓存复用 + 计数强刷；拖拽落序 + 隐藏夹同序 + 搜索禁拖）。mounted 里 `hiddenRows()` 改为只数宿主 `div`，因为排序包装组件也带同名 `className`。

r56 离线结果：`tsc --noEmit` 通过；全量 `vitest run` 5153 项中 5151 通过（r55 为 5144/5142，增量 9 = 本轮新增用例数），两条失败仍为已知基线（`main.browserMock`、`testPolicy`）。

r56 核验包：`build/bin/GoNavi-provider-settings-260903-r56`（100547186 字节，SHA-256 `6ceed9b833f806ca3fe92da92f464e2deb1a7ce6c8de9c9d1beed69c8956953e`），已被 r57 替换。

## 13. r57 实机反馈修正（2026-09-03）

r56 实机回看三条：拖拽副本不贴鼠标、悬浮无拖拽光标、模型停用提示在底部。

| # | 反馈 | 根因 / 修正 |
| --- | --- | --- |
| 1 | 拖拽时浮起副本与鼠标差一个位置，不贴合 | 设置页是 `ResizableDraggableModal`，外壳带 `transform: translate(...)`；`DragOverlay` 的 `position: fixed` 因此以弹窗为参照，副本落后鼠标一个弹窗偏移量。改为 `createPortal(overlay, document.body)`，`zIndex 1060`（高于弹窗 1000、低于 tooltip 1070），主题变量经 `overlayStyle={rootStyle}` 带过去。副本样式按流行拖拽形态：贴抓取点、`rotate(1.5deg) scale(1.03)`、阴影 + 绿描边、`opacity .88` 让被压住的卡片仍可读、`pointer-events: none`，落下 180ms 回位动画 |
| 2 | 悬浮可拖拽区域时光标能否换成拖拽样式 | 可行：`AIProviderSortableItem` 未禁拖时加 `is-draggable`，卡片本体与已隐藏行名称区 `cursor: grab`（眼睛 / 恢复按钮保持 pointer）；拖动期间 `body.gonavi-ai-provider-dragging *` 强制 `grabbing`，副本经过文字或按钮时光标不再来回切换。搜索过滤中不加 `is-draggable`，光标回到 pointer |
| 3 | 图2「已停用」提示在弹层底部，不知道是哪一个；停用行要置灰 | 去掉底部 `gonavi-ai-model-management-foot`，改为 `role=status` 视觉隐藏 live region（读屏仍能听到）。可见反馈变成该行开关旁的小浮窗：`ModelManagementRow` 新增 `flash` prop，复用开关上的 Tooltip，`open` 强制、`placement="left"`、`gonavi-ai-model-flash` 白底绿字，`MODEL_ROW_FLASH_MS = 1600` 后自动消失；被拦的「先换默认」原因、新增模型「已添加」也走同一浮窗落在对应行。停用行整行 `color: --gn-fg-3`，名字 0.6 透明，徽标继承灰色 |

新增测试：`AIProviderModelSelect.test.tsx`「confirms a toggle with a short-lived popover on that row」（假计时器验证 1600ms 消失、底部 foot 不存在、被拦原因替换前一条闪现）。

r57 离线结果：`tsc --noEmit` 通过；全量 `vitest run` 5154 项中 5152 通过（r56 为 5153/5151，+1 = 本轮新增），两条失败仍为已知基线。

r57 核验包：`build/bin/GoNavi-provider-settings-260903-r57`（100547186 字节，SHA-256 `9a5d6000b7bb83bf9261e03e10c2e6521a90767b4160b09988f0f74147ede4c3`），已被 r58 替换。

## 14. r58 并排阈值（2026-09-03）

反馈：选了并排后 URL 常独占一行，希望宽度够时 URL 与 Key 同排，URL 可收缩或压缩。

原因：并排的换行阈值按「能显示典型完整 URL / Key」定的（URL `flex-basis 320px`、Key 220px、格式 148px，三项同排需约 712px），但失焦时字段本来就只显示压缩文本（URL 去 `https://` 头 18 尾 12，Key 头尾各 4 位），阈值过于保守。

修正 [AISettingsProvidersSection.css §connection-fields]：格式 140px、URL `2 1 200px`、Key `1 1 150px`，三项同排从约 520px 编辑区开始；比压缩文字还窄时失焦输入框 `text-overflow: ellipsis` 而不是裁切；聚焦仍编辑全文。仅 CSS，测试 58/58 通过；无新增用例。

r58 核验包：`build/bin/GoNavi-provider-settings-260903-r58`（100547186 字节，SHA-256 `698ac52eef253be6a2f6ddfc907e90abf901a8dde567c824044add59d930ce71`），壳 `GoNavi-Provider-Verification-r58.app`。已被 r59 替换。

## 15. r59 合入上游后重建（2026-09-03）

`czz-dev` 已合并 `upstream/dev`（`a70650b3`）。源码晚于 r58，重打 r59。`tsc` 通过；全量 vitest 5144/5144。

r59 核验包：`build/bin/GoNavi-provider-settings-260903-r59`（101861154 字节，SHA-256 `8e000301579a9f6f1bf505832414726a24831d2a552c6225a3ad17545e4f283e`），壳 `GoNavi-Provider-Verification-r59.app`，pid 62193 单实例。identifier 又换了：`connectionLayout` 回到默认纵向。

**实机结果（2026-09-03）**：用户确认 r59 实验通过。覆盖 §7 清单，以及 r57（贴鼠标 / 抓手 / 行内浮窗）与 r58 并排阈值。未覆盖、仍不在本卡：已接入芯片行拖拽、真实模型回复、Windows/Linux 实机、签名发布包。
