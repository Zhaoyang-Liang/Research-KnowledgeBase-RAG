# 00_AI接手上下文 — 给任何接手此知识库的 AI

> 本目录存放 OpenClaw workspace 的关键上下文文件。
> 目标：**任何 AI（DeepSeek / GPT / Claude / Qwen 等）** 打开此目录即可理解用户身份、偏好和工作规则，无需管理 OpenClaw 或 QClaw。

## 文件说明

| 文件 | 用途 | 必读 |
|------|------|:--:|
| `USER.md` | 用户身份、研究领域、工作偏好、安全原则 | ⭐ |
| `IDENTITY.md` | 本 Agent 的公开定位和职责描述 | ⭐ |
| `SOUL.md` | 核心人设、沟通风格、工作原则 | ⭐ |
| `MEMORY.md` | 长期记忆：关键决策、偏好、教训 | ⭐ |
| `AGENTS.md` | 工作区规则、文件约定、工具使用 | ◯ |
| `TOOLS.md` | 本地工具笔记（摄像头、SSH、TTS等） | △ |
| `HEARTBEAT.md` | 心跳任务提醒（当前为空） | △ |

> ⭐ = 首次接手必读 | ◯ = 写文件/执行操作前查阅 | △ = 按需

## 接手阅读顺序

```
1. USER.md      → 理解用户是谁、研究什么
2. IDENTITY.md  → 理解自己是什么角色
3. SOUL.md      → 理解沟通风格和边界
4. MEMORY.md    → 了解历史决策和偏好
5. AGENTS.md    → 了解工作区规则
6. AI_无缝接手说明.md  → 全局 AI 接手流程与阅读顺序（本目录内）
7. 总流程入口.md        → 知识库唯一全局入口（根目录）
```

## ⚠️ 安全提示

- 本目录不含 API Key（安全规则禁止写入这些文件）
- USER.md 中可能含本地路径，不要外传
- MEMORY.md 是个人长期记忆，仅在 main session 加载

## 来源

原始文件位于 OpenClaw workspace：
```
~/.qclaw/workspace-agent-a0d46bbd/
```
本目录为副本，方便非 OpenClaw 环境下的 AI 直接读取。
