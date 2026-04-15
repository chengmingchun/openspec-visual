# OpenSpec 规格驱动 AI 编码系统规范

## 概述

OpenSpec 是一个用于 AI 结对编程的规格驱动开发工具，旨在解决 AI 编码中的"对齐"问题。它通过将需求、设计决策和实现步骤持久化为 Markdown 文件，为 AI 提供清晰的"行为契约"，从而减少沟通成本，提高开发效率。

### 核心问题
- **记忆丢失**：AI 的记忆仅存在于当前对话，对话关闭后所有上下文消失
- **对齐困难**：AI 对需求的理解可能与开发者不一致，导致代码偏离预期
- **上下文重建**：每次新对话都需要重新解释项目背景、技术栈和架构约束
- **中断恢复**：开发被打断后难以继续，需要手动拼凑上下文

## 系统架构

### 目录结构
```
openspec/
├── specs/          # 主规格：系统当前行为的"源真相"（Source of Truth）
│   ├── auth/       # 认证模块规格
│   ├── payments/   # 支付模块规格
│   └── ...         # 其他模块规格
├── changes/        # 活跃变更：正在进行的修改
│   ├── add-dark-mode/      # 深色模式功能变更
│   ├── fix-login-bug/      # 登录Bug修复变更
│   └── ...                 # 其他变更
├── archive/        # 归档变更：已完成的历史记录
│   └── 2026-02-27_add-github-oauth/
└── config.yaml     # 项目配置
```

### 核心概念

#### 1. Specs（主规格）
- **定义**：系统当前行为的权威描述，回答"系统现在是怎么运作的"
- **特点**：
  - 反映系统的最新真实状态
  - 随变更归档而更新
  - 只描述外部可观察的行为，不描述内部实现

#### 2. Changes（变更）
- **定义**：正在进行的修改，每个功能或Bug修复独立一个文件夹
- **特点**：
  - 互不干扰，支持并行开发
  - 包含完整的工件链（proposal → specs → design → tasks）
  - 归档后规格变化合并到主规格

#### 3. Delta Specs（增量规格）
- **定义**：只描述"这次改了什么"，而不是重新描述整个系统
- **类型**：
  - ADDED：新增的行为要求
  - MODIFIED：修改的行为要求
  - REMOVED：删除的行为要求
- **优势**：
  - 变更审查一目了然
  - 并行开发不冲突
  - 存量项目友好

## 工件体系

每个变更包含4个核心工件，按依赖关系生成：

### 1. proposal.md - 回答"为什么要做"
**目的**：定义变更的动机、范围和预期收益
**内容**：
- 问题陈述
- 解决方案概述
- 范围界定（做什么和不做什么）
- 成功标准
- 预期收益

### 2. specs/ - 回答"系统行为会怎么改变"
**目的**：用Delta Specs描述新增、修改、删除的行为
**格式**：
```markdown
## ADDED Requirements

### Requirement: Theme Switching
系统 MUST 提供深色/浅色主题切换功能。
系统 SHOULD 支持跟随操作系统主题设置。

#### Scenario: Manual Theme Switch
Given 用户当前使用浅色主题
When 用户点击主题切换开关
Then 界面切换为深色主题
And 选择结果持久化到 localStorage

## MODIFIED Requirements

### Requirement: Page Background (MODIFIED)
- 原：系统 MUST 使用固定白色背景（#FFFFFF）
- 新：系统 MUST 根据当前主题设置显示对应的背景色

## REMOVED Requirements

### Requirement: Fixed Color Scheme (REMOVED)
- 原：系统 MUST 使用预设的固定配色方案
- 原因：被新的主题系统取代
```

**RFC 2119 关键字**：
- **MUST**：必须实现，不实现就是Bug
- **SHOULD**：强烈建议，除非有充分理由可以不做
- **MAY**：可选的增强功能

**场景格式**：Given/When/Then
- Given：前置条件
- When：触发动作
- Then：预期结果

### 3. design.md - 回答"技术上怎么实现"
**目的**：描述架构决策、组件设计和技术选型的理由
**内容**：
- 技术方案
- 组件设计
- 数据流设计
- 技术选型理由
- 性能考虑

### 4. tasks.md - 回答"具体要干哪几件事"
**目的**：带复选框的实现清单，AI按此逐条执行
**格式**：
```markdown
## Implementation Tasks

- [x] 1. Install passport-github2 dependency
- [ ] 2. Create OAuth strategy configuration
- [ ] 3. Add callback route handler
- [ ] 4. Implement account linking logic
- [ ] 5. Update user model with OAuth fields
- [ ] 6. Add GitHub login button to UI
- [ ] 7. Write tests for OAuth flow
- [ ] 8. Update documentation
```

## 命令系统

### Core Profile（默认可用）

#### `/opsx:propose` - 需求清晰时一步到位
**功能**：一口气生成全套工件（proposal → specs → design → tasks）
**场景**：需求明确，直接开干
**示例**：
```
/opsx:propose 用户可以通过 GitHub OAuth 登录，
登录后自动创建账号，支持关联已有邮箱账号
```

#### `/opsx:explore` - 不确定时先探索
**功能**：调研分析，不创建任何文件（零副作用）
**场景**：需求模糊、技术选型、瓶颈分析
**示例**：
```
/opsx:explore 我们应该用 WebSocket 还是 SSE 来实现实时通知？
请分析当前的架构，评估两种方案
```

#### `/opsx:apply` - 按清单执行实现
**功能**：按tasks.md逐条执行任务
**特点**：
- 支持断点续传
- 智能跳过已完成任务
- 感知实际代码状态
**示例**：
```
/opsx:apply add-github-oauth
```

#### `/opsx:archive` - 收尾归档
**功能**：
1. 合并Delta Specs到主规格
2. 移动变更到归档目录
3. 关闭变更生命周期
**示例**：
```
/opsx:archive add-github-oauth
```

### Expanded Profile（需启用）

#### `/opsx:new` - 只建骨架
**功能**：创建变更目录和元数据，不生成工件内容
**场景**：手动控制节奏

#### `/opsx:continue` - 逐步生成
**功能**：每次执行生成一个工件，按依赖链推进
**场景**：需求还在打磨，每步都想审查

#### `/opsx:ff` - 快进生成
**功能**：补全所有剩余工件
**场景**：确认方向后加速

#### `/opsx:verify` - 质量检查
**功能**：从三个维度验证实现：
1. **完整性**：所有任务是否完成？需求场景是否覆盖？
2. **正确性**：实现是否匹配规格意图？边界条件是否处理？
3. **一致性**：代码结构是否反映设计决策？命名和模式是否统一？

**报告级别**：
- CRITICAL：必须修复
- WARNING：建议修复
- SUGGESTION：优化建议

#### `/opsx:sync` - 只同步规格
**功能**：合并Delta Specs到主规格，不变更归档
**场景**：并行变更需要引用

#### `/opsx:bulk-archive` - 批量归档
**功能**：一次性归档多个变更，自动检测并解决规格冲突
**场景**：并行开发后的统一收尾

#### `/opsx:onboard` - 交互式教程
**功能**：用真实代码库走一遍完整流程
**场景**：新手上手

### CLI 工具命令

```bash
openspec list                    # 查看所有活跃变更
openspec view                    # 交互式仪表盘
openspec show <change-name>      # 查看变更详情
openspec status <change-name>    # 查看工件完成进度
openspec validate --all --strict # 检查所有变更和规格格式
openspec archive <change-name>   # 从终端归档
openspec config profile          # 切换Profile
```

## 配置系统

### config.yaml
```yaml
schema: spec-driven

context: |
  技术栈：TypeScript、React 18、Node.js、PostgreSQL
  API 风格：RESTful，文档在 docs/api.md
  测试框架：Vitest + React Testing Library
  代码规范：参考 .eslintrc.js

rules:
  proposal:
    - 必须包含回滚方案
    - 标注影响的模块范围
  specs:
    - 使用 Given/When/Then 格式描述测试场景
```

**context字段**：注入到所有工件的生成过程中，避免重复交代技术栈
**rules字段**：针对特定工件类型的额外要求

### Schema 系统
- **默认schema**：`spec-driven`（proposal → specs → design → tasks）
- **自定义schema**：支持创建工作流变体
- **创建命令**：`openspec schema fork spec-driven <new-name>`

## 工作流程与底层运作模式

### 核心运作模式（The Core Operation Model）

OpenSpec 的核心运作模式可以抽象为一个“需求降维与契约沉淀”的漏斗模型。这不只是一个流水账式的工具，而是一个通过建立“不可变契约”来解决 AI 幻觉和偏离问题的协作范式：

1. **输入阶段（Input & Intent）**：用户通过输入一句话或一段散乱的想法，触发引擎进入。引擎此时使用大模型对意图进行提纯，剥离无效信息，生成正式的 `proposal.md`，这就好比传统开发中的“需求评审”。
2. **沉淀阶段（Architecture & Specification）**：AI 读取项目全局的 `config.yaml`（项目技术栈与基调），以此为标尺推导出技术实现方案 `design.md` 和用 Given/When/Then 精确定量边界条件的 `specs/`。这一步将“口语化的意想”转变为“确定性的工程边界”。
3. **推演与工件化阶段（Task Breakdown）**：得到契约后，将工程任务进一步降维打散分解为带复选框的 `tasks.md`，实现从不可见的代码重构工作，具象为可被原子化追踪的 Task Item。
4. **落地与归档阶段（Execution & Archiving）**：最终的落地环节通过 `/opsx:apply` 逐条执行。若在执行期间遭遇代码环境变化，`tasks.md` 允许动态插队和复选，AI 永远依赖该表单作为自己的“短期记忆”，完成功能后直接 `/opsx:archive` 进行封卷存入历史记录，实现历史可追溯。

---

### 典型工作场景流转

#### 场景一：需求清晰，直接开干
**路径**：`propose → 审查 → apply → verify → archive`
**适用**：需求明确，追求效率

### 场景二：需求模糊，边探索边明确
**路径**：`explore → new → continue → ... → apply → archive`
**适用**：需求不明确，需要逐步澄清

### 场景三：做到一半被打断，继续
**路径**：`（新对话）→ apply <change-name>`
**特点**：AI直接读取tasks.md，从断点继续

### 场景四：apply跑完，发现需求没做完
**处理**：直接编辑tasks.md，重新执行apply

### 场景五：纯技术调研
**路径**：`explore → （多轮对话）→ （确定后propose或new）`
**特点**：零副作用，只探索不执行

### 场景六：并行开发多个功能
**特点**：每个变更独立目录，互不干扰
**管理**：`openspec list`查看进度，`bulk-archive`统一归档

### 场景七：存量项目引入OpenSpec
**策略**：渐进式推进，用到哪补到哪
**步骤**：
1. `openspec init`初始化
2. `explore`分析现有代码
3. 逐步补充主规格
4. 新功能按正常流程走

## 最佳实践

### 1. 归档前先verify
养成习惯：apply完成后先verify再archive，在问题成本最低的阶段修复。

### 2. 一个变更，一个职责
变更名称应简洁明了：`add-user-avatar`、`fix-login-timeout`、`refactor-payment-module`
避免：`misc-improvements`、`feature-1`、`wip`

### 3. 审查工件，别急着apply
propose或ff生成工件后，花几分钟审查：
- tasks.md：任务拆分是否合理？颗粒度是否够细？
- specs/：边界条件是否覆盖？有没有遗漏？

### 4. 用config.yaml减少重复
将技术栈、代码规范、API风格等背景信息写入config.yaml的context字段。

### 5. apply前开新对话
在跑apply之前，开新对话窗口（清空历史上下文），避免对话噪音影响代码质量。

### 6. 选择高推理能力模型
- 高推理环节（propose、ff、continue）：Claude Opus、GPT-4
- 纯执行环节（apply）：要求相对较低

### 7. 更新还是新建的判断标准
| 情况 | 选择 |
|------|------|
| 意图没变，只是执行方案要调整 | 更新现有变更 |
| 范围在缩小（先做MVP） | 更新现有变更 |
| 做着做着发现要做的完全不是一回事 | 新建变更 |
| 范围膨胀到可以拆成两个独立功能 | 新建变更 |
| 原来的变更已经可以独立交付 | 归档旧的，新建新的 |

## 设计原则

### 1. 灵活，而非死板
- 没有"规划阶段不许写代码"的锁定
- 写到一半发现specs不对？回去改就是了
- 不是瀑布流程，不存在"阶段门禁"

### 2. 迭代，而非瀑布
- 不要求一次把所有事情想清楚
- 先写个大概，边做边完善
- 需求在实现过程中逐渐清晰

### 3. 简单，而非复杂
- 就是几个Markdown文件
- 没有数据库、没有服务端、没有Dashboard
- `openspec init`之后就能用

### 4. 存量优先
- 为"改存量系统"设计
- 不是只能在空白项目上玩的"理想流程"
- 渐进式推进，不需要一次性补齐所有规格

## 渐进式严格

### Lite Spec（日常开发）
- 简短的行为描述
- 清晰的范围界定
- 基本的验收条件
- 够AI理解要干什么就行

### Full Spec（高风险变更）
- 涉及跨团队协作、API变更、数据迁移
- 完整的Given/When/Then场景
- 边界条件分析
- 错误处理路径
- 返工成本高的场景值得多花时间

**判断标准**：如果这个变更搞砸了的返工成本很高，就多花点时间写规格；如果改错了5分钟就能修好，那写个大概就够了。

## 安装与初始化

### 前置要求
- Node.js 20.19.0+

### 安装
```bash
npm install -g @fission-ai/openspec@latest
```

### 初始化
```bash
cd your-project
openspec init
```

### 启用扩展命令
```bash
openspec config profile
openspec update
```
选择 **Expanded Profile** 后解锁完整命令集。

## Profile选择指南

| 工作方式 | 推荐Profile | 典型路径 |
|----------|-------------|----------|
| 需求通常清晰，追求效率 | Core（默认） | `propose → apply → archive` |
| 想逐步审查每个工件 | Expanded | `new → continue → apply → verify → archive` |
| 需求常常不明确 | Expanded | `explore → new → continue → apply` |
| 经常并行多个功能 | Expanded | 多个独立变更 + `bulk-archive` |

**建议**：从Core开始用，当需要更精细控制时再切到Expanded。

## 命令速查表

| 命令 | 说明 | 场景 |
|------|------|------|
| `/opsx:propose` | 一步生成完整变更 | 需求清晰 |
| `/opsx:explore` | 探索调研，不产生文件 | 需求模糊、技术选型 |
| `/opsx:apply` | 按任务清单写代码 | 实现阶段 |
| `/opsx:archive` | 归档，合并规格 | 功能完成收尾 |
| `/opsx:new` | 创建变更骨架 | 手动逐步推进 |
| `/opsx:continue` | 生成下一个工件 | 逐步审查 |
| `/opsx:ff` | 快进生成所有工件 | 确认方向后加速 |
| `/opsx:verify` | 三维度验证实现 | 归档前质量检查 |
| `/opsx:sync` | 只同步规格不归档 | 并行变更需引用 |
| `/opsx:bulk-archive` | 批量归档 | 多功能统一收尾 |
| `/opsx:onboard` | 交互式教程 | 新手上手 |
| `openspec list` | 查看活跃变更 | 日常管理 |
| `openspec status` | 查看工件完成度 | 了解当前进度 |
| `openspec view` | 交互式仪表盘 | 浏览变更和规格 |
| `openspec validate` | 验证格式 | 检查规格质量 |

## 参考资料

- [OpenSpec GitHub](https://github.com/Fission-AI/OpenSpec)
- [Getting Started](https://github.com/Fission-AI/OpenSpec/blob/main/docs/getting-started.md)
- [Workflows](https://github.com/Fission-AI/OpenSpec/blob/main/docs/workflows.md)
- [Commands Reference](https://github.com/Fission-AI/OpenSpec/blob/main/docs/commands.md)
- [CLI Reference](https://github.com/Fission-AI/OpenSpec/blob/main/docs/cli.md)
- [Concepts](https://github.com/Fission-AI/OpenSpec/blob/main/docs/concepts.md)
- [Customization](https://github.com/Fission-AI/OpenSpec/blob/main/docs/customization.md)

---

**版本**：1.0.0  
**创建日期**：2026-04-16  
**基于**：OpenSpec 完全使用指南（https://www.notemi.cn/openspec-complete-user-guide--driving-ai-encoding-with-specifications.html）  
**状态**：正式规范