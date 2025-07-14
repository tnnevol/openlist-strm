# 项目介绍

OpenList Stream 是一个现代化的流媒体管理系统，采用前后端分离的全栈架构设计。

## 项目架构

<div class="project-architecture">
  <div class="architecture-grid">
    <div class="architecture-item">
      <h3>前端应用 (apps/web-ele)</h3>
      <p>基于 Vue3 + Element Plus 的现代化前端应用</p>
      <ul>
        <li>Vue 3 + TypeScript</li>
        <li>Element Plus UI 组件库</li>
        <li>Vite 构建工具</li>
        <li>Pinia 状态管理</li>
        <li>Vue Router 路由管理</li>
      </ul>
    </div>
    <div class="architecture-item">
      <h3>后端服务 (backend-api)</h3>
      <p>基于 Gin 框架的 Go 语言后端 API 服务</p>
      <ul>
        <li>Gin Web 框架</li>
        <li>SQLite 数据库</li>
        <li>JWT 身份认证</li>
        <li>Swagger API 文档</li>
        <li>Zap 日志系统</li>
      </ul>
    </div>
    <div class="architecture-item">
      <h3>文档系统 (docs)</h3>
      <p>基于 VitePress 的项目文档站点</p>
      <ul>
        <li>VitePress 静态站点生成器</li>
        <li>Markdown 文档编写</li>
        <li>组件文档展示</li>
        <li>API 文档集成</li>
      </ul>
    </div>
  </div>
</div>

## 技术栈

### 前端技术栈

- **框架**: Vue 3 + TypeScript
- **UI 组件库**: Element Plus
- **构建工具**: Vite
- **状态管理**: Pinia
- **路由管理**: Vue Router
- **样式**: Tailwind CSS
- **代码规范**: ESLint + Prettier

### 后端技术栈

- **语言**: Go
- **Web 框架**: Gin
- **数据库**: SQLite
- **身份认证**: JWT
- **API 文档**: Swagger
- **日志系统**: Zap
- **测试框架**: Go 标准测试包

### 工程化工具

- **包管理**: pnpm
- **Monorepo**: Turbo
- **版本控制**: Git
- **CI/CD**: 待配置

## 项目特色

- 🚀 **现代化技术栈**: 采用最新的前端和后端技术
- 🏗️ **工程化架构**: Monorepo + Turbo 架构，规范且标准的大仓管理模式
- 🧪 **完善的测试体系**: 内置单元测试、集成测试、真实用户测试系统
- 📦 **模块化设计**: 前后端分离架构，支持独立部署和扩展
- 🎨 **现代化 UI**: 基于 Element Plus 组件库，提供美观、易用的用户界面
- 🔧 **开发工具链**: 集成现代化开发工具，提升开发效率

<style>
.project-architecture {
  margin: 2rem 0;
}

.architecture-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 2rem;
  margin-top: 2rem;
}

.architecture-item {
  padding: 1.5rem;
  border: 1px solid var(--vp-c-divider);
  border-radius: 8px;
  background: var(--vp-c-bg-soft);
}

.architecture-item h3 {
  margin: 0 0 1rem 0;
  color: var(--vp-c-brand);
  font-size: 1.2rem;
}

.architecture-item p {
  margin: 0 0 1rem 0;
  color: var(--vp-c-text-2);
}

.architecture-item ul {
  margin: 0;
  padding-left: 1.5rem;
}

.architecture-item li {
  margin: 0.5rem 0;
  color: var(--vp-c-text-1);
}

@media (max-width: 768px) {
  .architecture-grid {
    grid-template-columns: 1fr;
  }
}
</style>
