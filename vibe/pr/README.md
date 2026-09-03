# 上游 PR 历史

节奏化索引。每条记录的**文件名即标题**：`序号-日期-方式-概要`。

- **序号**：对本 fork 向上游 `Syngnat/GoNavi` 发出的 PR 流水，从 `01` 起，不回收。
- **日期**：该条主事件的 Asia/Shanghai 日历 `yyMMdd`（合入日；未开则用准备日）。
- **方式**：`squash` / `followup` / `pending`（尚未开 PR）。
- **概要**：中文短题，与 GitHub title 互补，不替代 GitHub title。

当前权威在本目录。旧 [czz-docs/upstream-pr-scope.md](../../../czz-docs/upstream-pr-scope.md) 已退役为指针。

新建下一条时复制 [_template.md](_template.md)，不要在 task card 里另写一份范围表。

## 时间线

| 序号 | 标题 | 状态 | GitHub | 过程稿 |
| --- | --- | --- | --- | --- |
| 01 | [01-260901-squash-供应商目录与CLI发现](01-260901-squash-供应商目录与CLI发现.md) | MERGED | [#1131](https://github.com/Syngnat/GoNavi/pull/1131)（[#1130](https://github.com/Syngnat/GoNavi/pull/1130) CLOSED） | [供应商管理总卡](../../specs/260901/0000-ai-provider-management/task-card.md) |
| 02 | [02-260902-followup-设置页折叠与CLI续命](02-260902-followup-设置页折叠与CLI续命.md) | MERGED | [#1134](https://github.com/Syngnat/GoNavi/pull/1134) | 同上总卡第十四轮 |
| 03 | [03-260903-squash-编辑收缩展示](03-260903-squash-编辑收缩展示.md) | opened | [#1155](https://github.com/Syngnat/GoNavi/pull/1155) | [编辑收缩卡](../../specs/260903/0000-provider-editor-compact/task-card.md) |

## 对上游的永久排除

这些前缀**永远不进** `Syngnat/GoNavi` 的 PR：

- `vibe/`（本仓过程与知识，含本目录）
- `czz-docs/`（退役研究稿与归档桩）
- `.codemark/`、`build/evidence/`
- `.gitignore` 里为本机代理发现加的 `.agents/`
- `frontend/package.json.md5`
- 读 CSS 文本并断言像素字面量的用例（#1131 已剔，禁止回流）

## 压平规则

不要用 `czz-dev` 直接对 `upstream/dev` 开 PR。从最新 `upstream/dev` 另起分支，squash 或 cherry-pick **该条主题**的增量。

## 全局 PR 规则同步

写 GitHub 正文时对齐 CodeNote `AiRef/VibePractice/Vibe_Rules/github/rules.md` §3：Summary、关键文件/模块、用户可见行为、风险与兼容、核验、有意义的 UI 截图、相关时的 DB/安全/发布说明、文档/记忆更新。本仓另加：排除项清单、是否从最新 `dev` 另起、例程策略是否误带。
