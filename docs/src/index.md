---
# https://vitepress.dev/reference/default-theme-home-page
layout: home
sidebar: false

hero:
  name: OpenList Stream
  text: 现代化流媒体管理系统
  tagline: 基于 Vue3 + Element Plus + Gin 的全栈解决方案
  image:
    src: https://unpkg.com/@vbenjs/static-source@0.1.7/source/logo-v1.webp
    alt: OpenList Stream
  actions:
    - theme: brand
      text: 快速开始 ->
      link: /guide/project/standard
    - theme: alt
      text: 在 GitHub 查看
      link: https://github.com/tnnevol/openlist-strm
    - theme: alt
      text: 后端 API 文档
      link: http://localhost:8890/swagger/index.html

features:
  - icon: 🚀
    title: 现代化技术栈
    details: 前端基于 Vue3、Element Plus、TypeScript，后端基于 Gin 框架，提供完整的全栈解决方案。
    link: /guide/project/standard
    linkText: 开发规范
  - icon: 🏗️
    title: 工程化架构
    details: 采用 Monorepo + Turbo 架构，规范且标准的大仓管理模式，提供企业级开发规范。
    link: /guide/project/cli
    linkText: CLI 工具
  - icon: 🧪
    title: 完善的测试体系
    details: 内置单元测试、集成测试、真实用户测试系统，确保代码质量和系统稳定性。
    link: /guide/project/test
    linkText: 测试文档
  - icon: 📦
    title: 模块化设计
    details: 前后端分离架构，支持独立部署和扩展，提供灵活的模块化解决方案。
    link: /guide/project/dir
    linkText: 目录说明
  - icon: 🎨
    title: 现代化 UI
    details: 基于 Element Plus 组件库，提供美观、易用的用户界面和丰富的交互体验。
    link: /components/introduction
    linkText: 组件文档
  - icon: 🔧
    title: 开发工具链
    details: 集成 Vite、Tailwind CSS、ESLint 等现代化开发工具，提升开发效率。
    link: /guide/project/vite
    linkText: 构建配置
---

<!-- 已彻底移除贡献者、团队成员相关内容 -->
