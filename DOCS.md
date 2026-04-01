# 📚 博客系统文档导航

欢迎使用博客系统！本文档中心提供了完整的使用指南、API 参考和开发教程。

---

## 🗂️ 文档目录

### 🚀 [快速启动指南](./QUICKSTART.md) ⭐ **新手必读**

**适合人群**: 首次接触项目的开发者  
**阅读时间**: 5 分钟

如果你想在 **5 分钟内快速启动项目**，请直接阅读此文档！

**内容概览**:
- ✅ 三步启动流程（数据库 → 后端 → 前端）
- ✅ 默认管理员账号
- ✅ 常用操作演示
- ✅ 快速故障排查

[👉 立即开始](./QUICKSTART.md)

---

### 📖 [完整项目文档](./README.md) 📘 **详细教程**

**适合人群**: 所有开发者  
**阅读时间**: 30 分钟

如果你想**深入了解项目的架构、功能和配置**，请仔细阅读此文档！

**内容概览**:
- 📋 项目简介和技术栈
- 🏗️ 项目结构详解
- 🔧 环境配置和依赖安装
- 🚀 详细的启动步骤
- 📝 功能模块详解（文章、分类、评论等）
- 🔌 API 接口文档
- ⚙️ 配置说明（前后端配置）
- 🐛 常见问题大全
- 📊 性能优化建议
- 🚀 生产环境部署指南
- 🤝 贡献指南

[👉 阅读全文](./README.md)

---

### 📡 [API 接口速查手册](./API.md) 🔌 **开发必备**

**适合人群**: 前后端开发者、测试人员  
**阅读时间**: 10 分钟（查阅式）

如果你需要**调用 API 或进行接口测试**，这是你的最佳参考！

**内容概览**:
- 🔐 认证接口（注册、登录、验证码）
- 📝 文章接口（CRUD、上传封面）
- 🏷️ 分类和标签接口
- 💬 评论接口
- 👍 互动接口（点赞、收藏）
- 👤 用户管理接口
- 📊 统计接口
- 🔍 搜索接口
- 📋 响应格式规范
- 🌐 Postman 集合示例

[👉 查看 API](./API.md)

---

## 🎯 推荐阅读路径

### 路径一：快速体验（15 分钟）

```
1. 阅读 [QUICKSTART.md](./QUICKSTART.md) 
   ↓
2. 按照步骤启动项目
   ↓
3. 访问 http://localhost:5173 体验功能
```

---

### 路径二：深入学习（1-2 小时）

```
1. 阅读 [QUICKSTART.md](./QUICKSTART.md) 启动项目
   ↓
2. 阅读 [README.md](./README.md) 了解架构
   ↓
3. 查看 [API.md](./API.md) 理解接口设计
   ↓
4. 阅读源代码深入理解实现
```

---

### 路径三：二次开发（按需查阅）

```
1. 确定开发需求
   ↓
2. 查阅 [API.md](./API.md) 了解现有接口
   ↓
3. 查阅 [README.md](./README.md) 相关章节
   ↓
4. 修改代码并测试
```

---

## 📂 项目文件结构

```
博客系统/
│
├── DOCS.md                 # 📚 本文档（文档导航）
├── README.md               # 📖 完整项目文档
├── QUICKSTART.md           # 🚀 快速启动指南
├── API.md                  # 📡 API 接口速查
│
├── frontend/               # 前端项目
│   ├── src/
│   │   ├── api/            # API 调用封装
│   │   ├── components/     # 公共组件
│   │   ├── router/         # 路由配置
│   │   ├── stores/         # Pinia 状态管理
│   │   ├── utils/          # 工具函数
│   │   └── views/          # 页面组件
│   ├── package.json
│   └── vite.config.js
│
└── backend/                # 后端项目
    ├── conf/
    │   └── app.conf        # 配置文件
    ├── controllers/        # 控制器层
    ├── models/             # 数据模型
    ├── routers/            # 路由配置
    ├── utils/              # 工具函数
    └── main.go             # 入口文件
```

---

## 🔍 快速查找

### 我想...

#### ...启动项目
👉 查看 [QUICKSTART.md](./QUICKSTART.md)

#### ...了解技术栈
👉 查看 [README.md - 技术栈章节](./README.md#技术栈)

#### ...调用某个 API
👉 查看 [API.md](./API.md) 对应接口

#### ...配置数据库
👉 查看 [README.md - 数据库配置](./README.md#1-数据库配置)

#### ...解决启动错误
👉 查看 [README.md - 常见问题](./README.md#常见问题)

#### ...部署到服务器
👉 查看 [README.md - 部署指南](./README.md#部署指南)

#### ...提交代码贡献
👉 查看 [README.md - 贡献指南](./README.md#贡献指南)

---

## 🆘 获取帮助

### 遇到问题？

1. **首先查看** [常见问题](./README.md#常见问题)
2. **搜索 Issue**: https://github.com/yourname/blog-system/issues
3. **提交 Issue**: 提供详细描述和复现步骤
4. **联系支持**: support@yourblog.com

### 文档问题？

如果发现文档有错误或缺失，欢迎：
- 📝 提交 Issue 指出问题
- 🔧 提交 Pull Request 修复文档

---

## 📝 文档更新记录

| 日期 | 文档 | 更新内容 |
|------|------|----------|
| 2026-04-01 | DOCS.md | 创建文档导航中心 |
| 2026-04-01 | README.md | 创建完整项目文档 |
| 2026-04-01 | QUICKSTART.md | 创建快速启动指南 |
| 2026-04-01 | API.md | 创建 API 接口速查手册 |

---

## 🎓 学习资源

### Vue 3 相关
- [Vue 3 官方文档](https://cn.vuejs.org/)
- [Composition API 入门](https://cn.vuejs.org/guide/extras/composition-api-faq.html)
- [Element Plus 组件库](https://element-plus.org/zh-CN/)

### Go & Beego 相关
- [Beego v2 官方文档](https://beego.me/)
- [Go 语言圣经](https://books.studygolang.com/gopl-zh/)
- [Beego ORM 使用指南](https://beego.me/docs/mvc/model/overview.md)

### 其他工具
- [Vditor 编辑器文档](https://vditor.js.org/)
- [WangEditor 文档](https://www.wangedor.com/)
- [Axios 中文文档](https://axios-http.cn/)

---

## 📄 许可证

MIT License

---

## 👥 反馈与建议

我们非常重视你的反馈！

- 💡 **功能建议**: 提交 Feature Request Issue
- 🐛 **Bug 报告**: 提交 Bug Report Issue  
- 📖 **文档改进**: 提交文档纠错 PR
- 💬 **讨论交流**: GitHub Discussions

---

**最后更新**: 2026-04-01

**文档版本**: v1.0.0

---

<div align="center">

### 🎉 祝你使用愉快！

如有任何问题，请随时查阅文档或寻求帮助。

[返回顶部](#博客系统文档导航)

</div>
