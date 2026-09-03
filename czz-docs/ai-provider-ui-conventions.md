# AI 供应商设置的界面与实现约定

本文件汇总 2026-09-01 一轮实机反馈中沉淀下来的规则。每条都来自一次真实缺陷，不是预设偏好。改动供应商设置页时先读本文件；它不随上游 PR 提交（见 [upstream-pr-scope.md](upstream-pr-scope.md)）。

产物与轮次的操作口径在 [gonavi-verify-build-restart.md](gonavi-verify-build-restart.md)。

## A. 悬浮提示

- **不使用原生 `title`。** 系统气泡约一秒才出、指针一移开就消失，且无法配置。需要提示时用 antd `Tooltip`。
- **进出延迟统一走共享常量**，进场 300ms、离场 150ms，取值锁在 `frontend/src/components/common/tooltipTiming.ts`。antd 默认的 100ms/100ms 会让指针扫过一排按钮时逐个闪烁。
- **供应商设置页的提示不截鼠标。** 目录卡片、已隐藏行、区块 ⓘ 与模型下拉都用 `passThroughHintTooltip`：离场 0ms，浮层 `pointer-events: none`。指针离开触发源立刻关掉，被浮层挡住的相邻项可以马上换新提示。
- **含内部选择的标题 ⓘ 例外。** 「收起编辑」左侧的标题提示里有连接字段布局开关，必须能点到浮层内的按钮，因此走 `interactiveHintTooltip`：进场仍 300ms，离场 400ms，浮层 `pointer-events: auto`。取值仍锁在 `tooltipTiming.ts`，不要在组件里手写延迟。
- **说明 + 分段钮 + ? 是共用件。** 分段钮直接复用顶部「紧凑 / 正常」的 `gonavi-ai-provider-density` 样式，不在浮层里另写一套边框、加粗或 `!important` 覆盖；钮文两字（「并排」「纵向」）。问号用 `passThroughHintTooltip`，正文只两行——并排一句、纵向一句（`connection_layout.inline_hint` / `stacked_hint`），不再有「够宽时可同一行」这类前导句；ⓘ 正文不含布局说明，只留模型选择那一句。
- 仓库里仍有 `0.35 / 0.4 / 0.75` 三处手写延迟（`App.tsx`、`TableOverview.tsx`、`TitleBarQuickActions.tsx`），尚未收敛；新代码不要再新增手写值。

## B. 提示密度与信息层级

- **说明性文案不占整行。** 多行说明合成一个 `ⓘ` 图标挂到它所属的标题旁；配置页高度有限时，整行说明块会把表单本身挤出视野。
- **收进图标不等于从 DOM 移除。** 图标内保留一份视觉隐藏文本（`clip-path: inset(50%)`），读屏可读，静态渲染也仍能命中——否则只依赖 Tooltip 的内容对辅助技术不可见。
- 报错类信息（`role="alert"`、字段校验）**必须常驻可见**，不得收进悬浮提示。

## C. 滚动与定位

- **禁止用 `Element.scrollIntoView` 做页内定位。** 它会滚动每一层可滚动祖先；设置中心那几层写着 `overflow: hidden`，仍可被程序滚动，但没有滚动条能滚回来，结果是页头被永久顶掉。
- 正确做法是只改目标容器自己的 `scrollTop`，实现见 `revealFirstErrorIn`：按 `getBoundingClientRect` 差值算偏移，夹到 `>= 0`，优先用容器的 `scrollTo({behavior:'smooth'})`。
- **失败后要主动定位。** 校验消息在点击之后才挂载，异步校验器更慢，因此保存/测试的包装函数各查两次：下一帧一次、200ms 再一次。
- 这类包装会破坏「按回调引用查按钮」的测试写法，测试应改为按标签查找。

## D. 弹窗几何

- **可缩放弹窗的高度要能压过调用方的行内样式。** 调用方经 antd `styles.content` 传入的 `height` 是行内声明，普通样式表规则打不过它；改尺寸的规则必须带 `!important`，并用 `[data-has-resized-height='true']` 之类的属性把作用域限死在「用户真的拖过」之后。
- **每个缩放手柄只钉自己那根轴。** 拖南边不应该把宽度冻在当时的像素值。
- **内容会切换的浮层要预留固定高度。** 两个选项卡内容不等高、而浮层向上弹出时，高度一变就靠移动顶边吸收，表现为上下跳动。预留值由 JS 单点持有，经 CSS 变量下发，样式表不重复这个数字。

## E. 破坏性操作

- 删除入口可以做成悬停显形的角标，但**必须保留二次确认**（`Popconfirm`），并在 `pendingProviderId` 或 `loading` 期间禁用。
- 角标与相邻的编辑按钮要拉开距离，并在悬停时给出与相邻控件不同的反馈（本仓库用一次短抖动 + 危险色），降低误点。
- 抖动类动效要在 `prefers-reduced-motion: reduce` 下退化为静态强调。
- 绝对定位到卡片外的控件，父级滚动容器要留出对应内边距，否则会被 `overflow` 裁掉。

## F. 渲染性能

- 列表行抽成 `React.memo` 组件，避免一次开关切换重建整份列表。
- **memo 生效的前提是回调引用稳定。** 父级用 `useRef` 承接最新闭包，再用空依赖的 `useCallback` 暴露稳定包装；直接把内联箭头函数传下去会让 memo 永远无法 bail out。
- 按钮要有本地按压反馈（`:active` 缩放、`:hover` 背景），让点击在状态回流之前就被感知。

## G. 布局与断点

- **断点从可用宽度预算推导，不要写死。** 窄屏断点 = 目录最小宽 + 手柄宽 + 编辑区最小宽；写死的数值会让实测 633px 的工作区意外落入抽屉模式。
- 抽屉态要铺满工作区，否则遮罩会在旁边裸露成灰板；抽屉外壳不透明度不足时，要显式隐藏底层内容，避免文字透出。
- 工作区之上的都是固定占位，每削减一像素编辑区就多一像素。固定高度的列表区（如已接入列表）用 `min(px, vh)` 而不是纯像素，避免条目变多时挤压编辑区。
- 顶栏说明与添加框、已接入工具条（密度/搜索）对齐预览稿的水平分组，间距只做小幅回放（约 +4–6px），不要把第十一轮的压缩整段撤掉。
- **认证三字段默认单列。** URL 与 API Key 必须能看全，不要再用 730px 容器查询强制三列。并排是标题 ⓘ 里的可选布局：失焦时 URL 去 `https://` 并 `...` 头尾压缩、Key 只留头尾 4 位，所以一行只需容纳压缩后的文字——格式约 140px、URL 约 200px、Key 约 150px 起排（编辑区约 520px 即可三项同排），不够才换行；比压缩文字还窄时失焦输入框用省略号而不是裁切。点进输入框仍编辑全文。选择写入 `gonavi.ai.providers.layout.v1` 的 `connectionLayout`。
- **单选项不要下拉。** API 格式只有一种时仍保留输入框边框（antd `Input` `readOnly`，`gonavi-ai-provider-fixed-value`）与 URL、Key 对齐，文字用 `--provider-muted` 浅灰，不出现 Select 箭头、不可编辑、不进 Tab 序。
- **CLI 模型列表按格式缓存。** `AIGetCLIModelCatalog` 会真的拉起本机 CLI，进入编辑或重开设置页不重复调用：可用结果按 `apiFormat` 写进 `gonavi.ai.providers.modelCatalog.v1`（`cliModelCatalogCache.ts`），下次直接沿用；只有点 `n/m 已启用` 才强制重拉。`stale`、空列表、失败不入缓存，下次自动重试。
- **目录与已隐藏列表可拖拽排序。** 顺序存 `layout.v1` 的 `presetOrder`（只存 key，新增预设按默认顺序补在已知项之后）。用 `@dnd-kit` 指针传感器（6px 起拖，点击不受影响），`rectSortingStrategy` 覆盖多列网格；拖动时原卡片留在流里变成虚化虚线占位并提前滑到落点，指针下是浮起的副本（`DragOverlay`）。搜索过滤中不允许拖（子序列的落点无法映射到全序）。两组共用一份全序，隐藏列表内部拖动只交换隐藏项的槽位，可见项位置不动。
- **拖拽副本必须 portal 到 `<body>`。** 设置页弹窗带 `transform`，`position: fixed` 的副本若留在弹窗树内会以弹窗为参照、落后鼠标一个偏移量。`zIndex 1060`，主题变量用 `overlayStyle` 传过去；副本贴抓取点、微倾 + 阴影 + 半透明、`pointer-events: none`。可拖拽面悬浮 `cursor: grab`，拖动中由 `body.gonavi-ai-provider-dragging` 强制 `grabbing`；眼睛 / 恢复按钮保持 pointer。
- **模型启用反馈落在该行。** 「已停用 / 已启用 / 先换默认 / 已添加」都是该行开关旁 1.6 秒的小浮窗（`MODEL_ROW_FLASH_MS`），不在弹层底部放共享说明行；`role=status` 只保留为视觉隐藏的 live region。停用行整行置灰。
- **空的 CLI 折叠不要留。** 本机 CLI 没有认证字段时，不渲染「本机 CLI」`<details>`；ⓘ 挂到正在编辑的已接入芯片右上角。「已接入 CLI 无需重复添加」只在编辑已保存的 CLI 时出现在目录工具条原位。
- **模型启用管理不要双页签。** 默认模型用上方选择器；`n/m 已启用` 只打开启用管理，点选择器本体只出普通菜单、不带管理外壳。候选项少时弹层随内容变矮。
- **默认模型标题只一行。** 「选择模型；留空跟随 CLI 默认」、来源、范围都在标题左 ⓘ；右 ⓘ 仅 CLI 供应商出现，放 CLI 登录说明、命令、effort 未校验、能力读取失败与刷新说明。同一句说明不在两个 ⓘ 里重复。

## H. 命名

- **副本命名按重名判定，不无条件加后缀。** 用户已经改过名字且不与既有配置重名时，原样保留；只有撞名才追加后缀并递增序号。
- 判定只比较完整字符串，不做大小写或空白归一化。

## I. 单例与多实例

- 单例 CLI 预设复用同一份本机登录，复制第二份没有意义：这类配置**不渲染**「另存为」入口，而不是渲染后置灰。
- 普通 API 可多份接入，「另存为」收进保存按钮的下拉里，主按钮保持「保存」。

## J. 测试约定

- 组件里新增 antd 组件或 `@ant-design/icons` 图标，要同步补进各测试文件的 `vi.mock` 导出表，否则整个文件会以 `No "X" export is defined` 集体失败。
- 需要断言的浮层内容（`Tooltip` 的 `title`、`Popconfirm` 的 `onConfirm`）在 mock 里保留为可寻址元素，不要简化成 `<>{children}</>`。
- **测试随改动一起提，不单独剥离。** 上游有 521 个 `_test.go` 与 540 个 `.test.ts(x)`，且 PR 一提交 CI 就跑 `go test ./...`；改了函数签名却不同步既有测试，会直接编译失败。缩小 PR 的正确手段是按主题拆分，不是抽走测试。细节见 [upstream-pr-scope.md](upstream-pr-scope.md) §3.5。
- **不要把样式表当字符串断言。** 读 `.css` 文本去断言 `left: 31px` 这类像素字面量，会因任何一次格式化而误报；它还绕开了仓库 `testPolicy` 守卫（其正则只匹配 `.ts/.tsx`）。要验位置就渲染后断言类名或计算值。
- **生成文件提 PR 前先用 `git diff -w` 核对真实增量。** `frontend/wailsjs/go/models.ts` 曾出现 694 行改动，其中只有 34 行是内容，其余 660 行是空行风格差异——本机 wails 在空行写 `\t`，上游写空。提交前把「整行仅由空白构成」的行归一化即可消掉，不要动任何有内容的行。不做这一步，review 面会凭空放大二十倍。
- **DOM 行为抽成纯函数再测。** `react-test-renderer` 没有真实 DOM，滚动、测量类逻辑要抽出可注入容器的纯函数（如 `revealFirstErrorIn`），在单元层覆盖调用形状与边界；纯 CSS 数值不适合写断言，应明确标为待实机确认。

## K. 表单折叠行

- **不要把 `display: flex` 直接打在 `<summary>` 上。** WKWebView（Wails 桌面端）会藏掉原生三角，并且标题和 ⓘ 之间的空白点不着，看起来像“这栏不能收展”。
- 折叠行用内层满宽 flex 条（`.gonavi-ai-provider-disclosure`）承载标题列、caret 和 ⓘ；`<summary>` 保持块级并关掉 `::marker` / `::-webkit-details-marker`。
- 箭头复用已接入列表的 `gonavi-ai-provider-caret`（收起 `RightOutlined`，展开 `DownOutlined`），放在标题右侧的共享列里（`.gonavi-ai-provider-disclosure-lead`），不要 `margin-left: auto` 甩到行尾；两行标题列同宽，箭头竖向对齐。整栏包括空白仍是点击面。
- ⓘ 仍 `stopPropagation`，避免看提示时把栏收展开。
