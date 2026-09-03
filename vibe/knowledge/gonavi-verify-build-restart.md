# 核验产物的构建与重启通路

本文件是本仓库「改一轮 → 出一个独立核验程序 → 停旧起新」这条通路的规则正文。它不随上游 PR 提交（见 [上游 PR 史](upstream-pr/README.md)）。

适用于任何需要用户在真实 macOS 桌面确认观感或行为的改动。编译通过与离线测试通过**都不构成**实机通过。

## 1. 何时走这条通路

| 情形 | 是否重建 |
| --- | --- |
| 改了 `frontend/src`、`internal/`、`shared/i18n` 中任一文件 | 是，出新一轮 |
| 只改了文档、测试夹具或注释 | 否 |
| 用户只说「重启」且工作区源码 mtime 全部早于当前产物 | 否，直接停旧起新 |

判断是否需要重建，用实测时间戳而不是记忆：

```bash
stat -f "%Sm %N" -t "%Y-%m-%d %H:%M:%S" build/bin/GoNavi-provider-settings-<yymmdd>-r<N>
git status --porcelain=v1 | awk '{print $2}' | grep -E '^(frontend/src|internal|shared)' \
  | xargs stat -f "%Sm %N" -t "%Y-%m-%d %H:%M:%S" | sort | tail -3
```

产物晚于全部源码即无需重建。

## 2. 出一轮的完整步骤

轮次号 `rN` 单调递增，不复用、不回退。以下命令按顺序执行，任一步失败即停止，不得跳过后续核验直接启动。

### 2.1 落库前核验

```bash
cd frontend && npx tsc --noEmit          # 必须退出码 0
cd frontend && npx vitest run            # 见 §4 的已知失败基线
go build ./internal/... ./cmd/...        # 仅后端改动时至少跑这条
```

### 2.2 构建

```bash
cd frontend && npm run build             # Vite 生产包
cd <repo-root> && wails build -s -skipbindings -nosyncgomod
```

`-skipbindings` 与 `-nosyncgomod` 是本仓库既定口径：绑定与 `go.mod` 由人工维护，构建期不自动改写。

### 2.3 打包独立核验程序

从上一轮的 `.app` 复制模板，然后**四个键都要改**——这是踩过的坑，漏改任何一个都会让人误判运行版本：

```bash
cd build/bin
cp GoNavi.app/Contents/MacOS/GoNavi GoNavi-provider-settings-<yymmdd>-r<N>
rm -rf GoNavi-Provider-Verification-r<N>.app
cp -R GoNavi-Provider-Verification-r<N-1>.app GoNavi-Provider-Verification-r<N>.app
rm -f GoNavi-Provider-Verification-r<N>.app/Contents/MacOS/*
cp GoNavi-provider-settings-<yymmdd>-r<N> GoNavi-Provider-Verification-r<N>.app/Contents/MacOS/

P=GoNavi-Provider-Verification-r<N>.app/Contents/Info.plist
plutil -replace CFBundleExecutable  -string GoNavi-provider-settings-<yymmdd>-r<N> $P
plutil -replace CFBundleIdentifier  -string com.czz.gonavi.provider-verification.r<N> $P
plutil -replace CFBundleDisplayName -string "GoNavi Provider Verification r<N>" $P
plutil -replace CFBundleName        -string "GoNavi Provider Verification r<N>" $P

shasum -a 256 GoNavi-provider-settings-<yymmdd>-r<N>
```

### 2.4 停旧起新

必须保证单实例：多个实例会争抢同一份供应商配置，观察到的状态不可信。只杀上一轮不够。

纯重启（不重建）从仓库根跑 Skill 里的脚本：

```bash
.agents/skills/gonavi-verify-build-restart/restart.sh        # 最高 rN
.agents/skills/gonavi-verify-build-restart/restart.sh r45    # 指定轮次
```

脚本会杀掉所有 `GoNavi-provider-settings-` 进程再打开目标 `.app`。回读必须只剩一行，且路径是本轮的 `.app`。打包完成后的起新也走这条，不要手抄 kill/open。

### 2.5 记录

同轮更新当前 task card 的验收产物段（供应商主线：[管理总卡](../specs/260901/0000-ai-provider-management/task-card.md)；编辑收缩：[收缩卡](../specs/260903/0000-provider-editor-compact/task-card.md)）：产物路径、字节数、SHA-256、核验包名。旧轮产物保留，不删。

## 3. 不变量

- **运行版本以进程路径为准，不看菜单名。** r14 与 r15 的 `CFBundleDisplayName` 滞留在 r13，菜单显示的版本号是错的；自 r16 起才随轮次同步。
- **每轮换 identifier 意味着 WebView 本地存储清空。** 布局偏好、`catalogWidth`、模型启停等本地状态都会回到默认。这在验证「默认行为」时是干净起点，在验证「延续状态」时是陷阱。
- **应用内「检查更新」会整包替换核验程序。** r14 的 `.app` 曾被官方 0.9.4 发布包顶掉（identifier 变 `com.wails.GoNavi`、可执行文件改名、多出 `_CodeSignature`）。独立二进制 `GoNavi-provider-settings-*` 不受影响，可据此重新打包。核验期间不要点更新。
- **构建目录、截图与本机证据不入库。** `build/bin/` 已被 `.gitignore` 忽略；`build/evidence/` 未忽略但按任务卡口径同样不提交。
- **编译与离线测试不等于实机通过。** 交付说明里必须把「已验」与「未验」分开写，未验项要指名由用户在哪一版、点哪里确认。

## 4. 已知失败基线

2026-09-03 合入 `upstream/dev`（`a70650b3`）后，本机全量 `npx vitest run` 为 **5144/5144**。原先两条长期失败已被上游消掉，不得再当「允许红」用。

合入前的历史基线（仅作对照，不要拿来开脱新失败）：

| 用例 | 当时性质 |
| --- | --- |
| `src/main.browserMock.test.ts > localizes browser mock provider test messages` | 合并上游前即已复现 |
| `src/testPolicy.test.ts > adds no new test that reads source files from disk` | 上游自带守卫误伤，后已在上游修掉 |

判断是否引入新失败，看**通过数增量**是否等于本轮新增用例数，而不是只看失败条数。

## 5. 与 Skill 的关系

这条通路已升级为 `project` 作用域的 Skill：`gonavi-verify-build-restart`。

- **唯一物理来源**在 CodeNote 的 `AiRef/VibePractice/Skills/projects/gonavi/gonavi-verify-build-restart/`，由 `init_skill.py --scope project --project-key gonavi` 生成。
- **发现链接**有两处，都是指向该包的软链而非副本：治理默认面 `.agents/skills/`，本宿主消费面 `.claude/skills/`（Claude Code 不扫 `.agents/skills`）。两个目录均已在 `.gitignore` 中忽略，软链指向用户 Home，不得进版本库。
- 前置的仓库登记也已完成：`vibe/knowledge/project-index.json` 的 `gonavi` 条目与 `workspace-config/workspace.local.json` 的 `repository_bindings.gonavi`。未登记时 `init_skill.py` 会报 `Project root is not a verified canonical project`。

本文件与该 Skill 是同一份知识的两个投影：Skill 供 Agent 按步执行，本文件供人阅读与评审。**改动时两边同改**，不要只改一边；也不要在其他发现目录下另建可独立编辑的副本，那会违反 Skill 治理的「单一物理来源」。
