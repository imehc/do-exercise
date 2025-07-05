## 首次运行说明

本项目前端依赖于根据后端 `openapi.yaml` 自动生成的 API 接口代码，首次运行前请确保 `src/do-exercise-api` 目录已存在。

如未生成，请执行以下命令生成前端 API 代码：

```bash
pnpm gen:apis
```

如果你在 Windows 平台，或无法执行 `pnpm gen:apis`，请执行：

```bash
pnpm gen:apis-other
```

> **注意**：如果启动或开发过程中提示找不到 `do-exercise-api` 相关模块，请先确认 `src/do-exercise-api` 目录是否存在，若不存在请按上述命令生成。
