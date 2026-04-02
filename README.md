# 博客系统 - 完整开发文档

## 📖 项目简介

这是一个基于 **Vue 3 + Beego v2** 的全栈博客系统，采用前后端分离架构。系统提供文章管理、分类标签、用户认证、评论互动等完整的博客功能，支持 Markdown 编辑器，拥有现代化的 UI 设计和响应式布局。

### ✨ 主要特性

- 📝 **文章管理** - 创建、编辑、删除文章，支持草稿和发布状态
- 🏷️ **分类标签** - 双轨制分类系统（系统标签 + 用户自定义合集）
- 💬 **评论互动** - 支持文章评论、回复功能
- 👤 **用户系统** - JWT 认证、邮箱验证、密码找回
- 🎨 **Markdown 编辑** - 集成 Vditor 和 WangEditor 双编辑器
- 📊 **数据统计** - 浏览量统计、用户数据概览
- 🔐 **权限管理** - 普通用户/管理员分级权限
- 📱 **响应式设计** - 完美适配 PC 和移动端

---

## 🛠️ 技术栈

### 前端技术栈

| 技术 | 版本 | 说明 |
|------|------|------|
| **框架** | Vue 3.4.0 | Composition API, 响应式系统 |
| **构建工具** | Vite 5.0.0 | 快速开发和热更新 |
| **UI 组件库** | Element Plus 2.5.0 | 桌面端组件库 |
| **状态管理** | Pinia 2.1.0 | Vue 3 官方推荐状态管理 |
| **路由** | Vue Router 4.2.0 | SPA 路由管理 |
| **HTTP 客户端** | Axios 1.6.0 | API 请求封装 |
| **编辑器** | Vditor 3.11.2<br>WangEditor 5.1.23 | Markdown 富文本编辑器 |
| **工具库** | @vueuse/core 10.7.0<br>Day.js 1.11.10<br>NProgress 0.2.0 | Vue 组合式函数<br>日期处理<br>进度条 |
| **样式** | Sass 1.69.0 | CSS 预处理器 |

### 后端技术栈

| 技术 | 版本 | 说明 |
|------|------|------|
| **Web 框架** | Beego v2.3.9 | Go 语言全栈框架 |
| **ORM** | Beego ORM | 内置 ORM 支持 |
| **数据库驱动** | MySQL Driver 1.9.3 | MySQL 数据库连接 |
| **JWT 认证** | golang-jwt/jwt v5.3.1 | Token 认证 |
| **加密** | bcrypt | 密码加密存储 |
| **邮件发送** | email 4.0.1 | SMTP 邮件服务 |

### 数据库

- **MySQL** - 主数据库
- **SQLite** - 可选的备用数据库

---

## 📁 项目结构

```
博客系统/
├── frontend/                    # 前端项目
│   ├── src/
│   │   ├── api/                # API 接口定义
│   │   ├── components/         # 公共组件
│   │   ├── router/             # 路由配置
│   │   ├── stores/             # Pinia 状态管理
│   │   ├── styles/             # 全局样式
│   │   ├── utils/              # 工具函数
│   │   └── views/              # 页面组件
│   ├── package.json            # 依赖配置
│   └── vite.config.js          # Vite 配置
│
├── backend/                     # 后端项目
│   ├── conf/
│   │   └── app.conf            # 配置文件
│   ├── controllers/            # 控制器层
│   │   ├── admin.go            # 管理后台控制器
│   │   ├── article.go          # 文章相关控制器
│   │   ├── auth.go             # 认证相关控制器
│   │   └── user_category.go    # 用户分类控制器
│   ├── models/                 # 数据模型
│   │   └── models.go           # 所有模型定义
│   ├── routers/                # 路由配置
│   ├── utils/                  # 工具函数
│   │   ├── captcha.go          # 验证码生成
│   │   ├── cors.go             # CORS 配置
│   │   ├── email.go            # 邮件发送
│   │   ├── jwt.go              # JWT 工具
│   │   └── middleware.go       # 中间件
│   ├── uploads/                # 上传文件存储
│   ├── sessions/               # Session 文件存储
│   ├── main.go                 # 入口文件
│   └── go.mod                  # Go 模块配置
│
└── README.md                   # 项目文档
```

---

## 🚀 快速开始

### 环境要求

- **Node.js**: >= 18.0.0
- **Go**: >= 1.25.0
- **MySQL**: >= 5.7 (推荐 8.0+)
- **包管理器**: npm >= 9.0.0

### 1. 数据库配置

#### 1.1 创建数据库

```sql
-- 登录 MySQL
mysql -u root -p

-- 创建数据库
CREATE DATABASE IF NOT EXISTS blog DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 创建用户并授权（可选）
CREATE USER 'blog_user'@'localhost' IDENTIFIED BY 'your_password';
GRANT ALL PRIVILEGES ON blog.* TO 'blog_user'@'localhost';
FLUSH PRIVILEGES;
```

#### 1.2 导入数据表结构

```bash
# 进入数据库目录
cd backend/database

# 导入初始化脚本
mysql -u root -p blog < init.sql

# 或者在 MySQL 命令行中执行
USE blog;
SOURCE init.sql;
```

#### 1.3 修改数据库配置

编辑 `backend/conf/app.conf`:

```ini
# 数据库配置 (MySQL)
dbdriver = "mysql"
dbconn = "root:123456@tcp(127.0.0.1:3306)/blog?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=true&loc=Local"
dbmaxidle = 10
dbmaxopen = 100
```

**连接字符串格式说明**:
```
用户名：密码@tcp(主机：端口)/数据库名？参数
```

请根据实际情况修改：
- `root`: 数据库用户名
- `123456`: 数据库密码
- `127.0.0.1:3306`: MySQL 地址和端口
- `blog`: 数据库名称

---

### 2. 后端启动

#### 2.1 安装 Bee 工具

```bash
# 安装 Bee 开发工具
go install github.com/beego/bee/v2@latest

# 验证安装
bee version
```

#### 2.2 安装依赖

```bash
cd backend
go mod download
```

#### 2.3 启动后端服务

```bash
# 方式一：使用 Bee 工具（推荐，支持热更新）
bee run

# 方式二：直接运行
go run main.go

# 方式三：编译后运行
go build -o blog.exe
./blog.exe
```

**启动成功标志**:
```
[INFO] [Main.Init] 正在初始化数据库连接 - driver: mysql
[INFO] [Main.Init] 数据库连接池配置 - maxIdle: 10, maxOpen: 100
[INFO] [Main.Init] 数据模型已注册
[INFO] [Main] ✅ 数据库连接正常
[INFO] [Main] 正在启动 Web 服务...
[INFO] [Main] 静态文件服务已配置：/uploads -> ./uploads
Server is running on :8080
```

**访问测试**: http://localhost:8080

---

### 3. 前端启动

#### 3.1 安装依赖

```bash
cd frontend
npm install
```

如果遇到网络问题，可以使用淘宝镜像：

```bash
npm config set registry https://registry.npmmirror.com
npm install
```

#### 3.2 启动开发服务器

```bash
# 启动开发服务器（支持热更新）
npm run dev

# 指定端口（默认 5173）
npm run dev -- --port 3000
```

**启动成功标志**:
```
  VITE v5.0.0  ready in 1234 ms

  ➜  Local:   http://localhost:5173/
  ➜  Network: use --host to expose
  ➜  press h to show help
```

#### 3.3 访问应用

打开浏览器访问：**http://localhost:5173**

---

## 📋 功能模块详解

### 1. 用户认证系统

#### 注册流程

1. 访问 `/register` 页面
2. 填写用户名、邮箱、密码
3. 输入图形验证码
4. 提交后自动登录并跳转到首页

#### 登录流程

1. 访问 `/login` 页面
2. 输入用户名和密码
3. 输入图形验证码
4. 登录成功后跳转到首页

#### JWT 认证机制

- **Token 有效期**: 7 天
- **存储位置**: localStorage
- **自动刷新**: 每次请求自动续期
- **安全策略**: 
  - 密码错误 5 次锁定账户 30 分钟
  - 敏感操作需要重新验证

---

### 2. 文章管理

#### 创作中心

**路径**: `/creator`

**功能**:
- 📝 创建新文章
- ✏️ 编辑已有文章
- 🗑️ 删除文章
- 📊 查看文章状态（草稿/已发布）
- 🔍 按状态筛选

#### 文章编辑器

**双编辑器支持**:

1. **Vditor** (Markdown 编辑器)
   - 实时预览
   - 语法高亮
   - 表格、流程图支持
   - 代码块插入

2. **WangEditor** (富文本编辑器)
   - 所见即所得
   - 图片上传
   - 视频嵌入
   - 格式化文本

**文章属性**:
- 标题（必填）
- 内容（必填）
- 摘要（自动生成或手动填写）
- 封面图（可选）
- 分类（系统标签或个人合集）
- 标签（多选）
- 状态（草稿/发布）

---

### 3. 分类系统

#### 双轨制分类

**系统标签 (Tag)**:
- 全局共享
- 管理员维护
- 用于文章归类

**个人合集 (User Category)**:
- 用户私有
- 自定义颜色标识
- 权限隔离
- 可删除（有文章关联时禁止删除）

#### 管理后台分类管理

**路径**: `/admin/categories`

**功能**:
- 创建/编辑/删除分类
- Slug 自动生成
- 关联检查（有文章使用时禁止删除）

---

### 4. 评论系统

**功能特性**:
- 嵌套回复（最多 3 层）
- @提及用户
- 表情支持
- Markdown 格式
- 实时通知

---

### 5. 管理后台

**路径**: `/admin`

**权限要求**: 管理员角色（role = 0）

**功能模块**:

| 模块 | 路径 | 功能 |
|------|------|------|
| **数据概览** | `/admin` | 用户数、文章数、评论数统计 |
| **文章管理** | `/admin/articles` | 审核、推荐、删除文章 |
| **分类管理** | `/admin/categories` | 管理系统标签 |
| **标签管理** | `/admin/tags` | 管理系统标签 |
| **用户管理** | `/admin/users` | 用户列表、角色分配、封禁 |

---

## 🔌 API 接口文档

### 认证接口

#### 用户注册
```http
POST /api/auth/register
Content-Type: application/json

{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "123456",
  "captcha_id": "xxx",
  "captcha_ans": "1234"
}
```

#### 用户登录
```http
POST /api/auth/login
Content-Type: application/json

{
  "username": "john_doe",
  "password": "123456",
  "captcha_id": "xxx",
  "captcha_ans": "1234"
}
```

**响应示例**:
```json
{
  "code": 0,
  "message": "登录成功",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "username": "john_doe",
      "email": "john@example.com",
      "nickname": "John",
      "avatar": "/uploads/avatars/xxx.jpg",
      "role": 1
    }
  }
}
```

---

### 文章接口

#### 获取文章列表
```http
GET /api/articles?page=1&page_size=10&category_id=1&tag_id=2&keyword=test
Authorization: Bearer {token}
```

#### 获取文章详情
```http
GET /api/articles/:id
Authorization: Bearer {token}
```

#### 创建文章
```http
POST /api/articles
Content-Type: application/json
Authorization: Bearer {token}

{
  "title": "我的第一篇博客",
  "content": "文章内容...",
  "summary": "文章摘要",
  "cover_image": "/uploads/covers/xxx.jpg",
  "category_id": 1,
  "user_category_id": 2,
  "tags": [1, 2, 3],
  "status": 1  // 0: 草稿，1: 发布
}
```

#### 更新文章
```http
PUT /api/articles/:id
Content-Type: application/json
Authorization: Bearer {token}

{
  "title": "更新后的标题",
  "content": "更新后的内容",
  ...
}
```

#### 删除文章
```http
DELETE /api/articles/:id
Authorization: Bearer {token}
```

---

### 用户分类接口

#### 获取我的分类
```http
GET /api/my-categories
Authorization: Bearer {token}
```

#### 创建分类
```http
POST /api/user-categories
Content-Type: application/json
Authorization: Bearer {token}

{
  "name": "技术笔记",
  "color": "#409EFF"
}
```

#### 更新分类
```http
PUT /api/user-categories/:id
Content-Type: application/json
Authorization: Bearer {token}

{
  "name": "新的分类名",
  "color": "#67C23A"
}
```

#### 删除分类
```http
DELETE /api/user-categories/:id
Authorization: Bearer {token}
```

---

## ⚙️ 配置说明

### 后端配置 (app.conf)

```ini
# 应用配置
appname = blog-system
httpport = "8080"
runmode = dev              # dev/prod
autorender = true
copyrequestbody = true

# 数据库配置
dbdriver = "mysql"
dbconn = "root:123456@tcp(127.0.0.1:3306)/blog?charset=utf8mb4"
dbmaxidle = 10
dbmaxopen = 100

# Session 配置
sessionon = true
sessionprovider = "file"
sessionproviderconfig = "./sessions"
sessiongcmaxlifetime = 3600

# CORS 配置（跨域）
cors_enabled = true
cors_allow_origins = "*"
cors_allow_methods = "GET,POST,PUT,DELETE,PATCH,OPTIONS"
cors_allow_headers = "Origin,X-Requested-With,Content-Type,Accept,Authorization"
cors_allow_credentials = "true"
cors_max_age = "86400"

# JWT 配置
jwt_secret_key = "your-secret-key-change-in-production"
jwt_expire_hours = 168     # 7 天

# 邮件配置（用于密码找回）
smtp_host = "smtp.example.com"
smtp_port = 587
smtp_username = "noreply@example.com"
smtp_password = "your-smtp-password"
```

---

### 前端配置

#### 环境变量 (.env)

```bash
# 开发环境
VITE_API_BASE_URL=http://localhost:8080/api
VITE_UPLOAD_BASE_URL=http://localhost:8080

# 生产环境
VITE_API_BASE_URL=https://api.yourblog.com/api
VITE_UPLOAD_BASE_URL=https://api.yourblog.com
```

#### API 请求封装 (utils/request.js)

```javascript
import axios from 'axios'

const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
  timeout: 10000
})

// 请求拦截器 - 添加 Token
request.interceptors.request.use(config => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// 响应拦截器 - 处理错误
request.interceptors.response.use(response => {
  const { code, message } = response.data
  if (code !== 0) {
    if (code === 401) {
      // Token 过期，跳转登录
      localStorage.removeItem('token')
      window.location.href = '/login'
    }
    return Promise.reject(new Error(message))
  }
  return response.data
}, error => {
  console.error('请求失败:', error)
  return Promise.reject(error)
})

export default request
```

---

## 🐛 常见问题

### 1. 后端启动失败

#### 错误：数据库连接失败

**原因**: 数据库配置错误或 MySQL 未启动

**解决方案**:
```bash
# 检查 MySQL 服务状态
# Windows
net start | findstr MySQL

# macOS/Linux
sudo systemctl status mysql

# 启动 MySQL
# Windows
net start MySQL80

# macOS
brew services start mysql

# Linux
sudo systemctl start mysql
```

#### 错误：端口被占用

**现象**: `bind: address already in use`

**解决方案**:
```bash
# 查看端口占用
# Windows
netstat -ano | findstr :8080

# macOS/Linux
lsof -i :8080

# 修改端口号
# 编辑 backend/conf/app.conf
httpport = "8081"
```

---

### 2. 前端启动失败

#### 错误：找不到模块

**现象**: `Cannot find module 'xxx'`

**解决方案**:
```bash
# 删除 node_modules 和 package-lock.json
rm -rf node_modules package-lock.json

# 重新安装
npm install
```

#### 错误：跨域请求失败

**现象**: `Access-Control-Allow-Origin` 错误

**解决方案**:
1. 确认后端 CORS 配置已启用
2. 检查前端 API 基础 URL 是否正确
3. 确保后端已启动且监听正确端口

---

### 3. 登录失败

#### 错误：验证码错误

**原因**: Session 配置问题或验证码过期

**解决方案**:
```bash
# 检查 sessions 目录是否存在且有写入权限
# Windows
icacls backend\sessions /grant Users:F

# macOS/Linux
chmod 755 backend/sessions
```

#### 错误：Token 无效

**原因**: JWT 密钥不一致或 Token 过期

**解决方案**:
```bash
# 清除本地 Token
localStorage.removeItem('token')

# 重新登录

# 检查后端 JWT 配置
# backend/utils/jwt.go
var secretKey = []byte("your-secret-key")
```

---

### 4. 文件上传失败

#### 错误：无法保存文件

**原因**: uploads 目录无写入权限

**解决方案**:
```bash
# Windows
icacls backend\uploads /grant Users:F

# macOS/Linux
chmod -R 755 backend/uploads
```

#### 错误：文件大小超限

**限制**: 单文件最大 5MB

**解决方案**:
- 压缩图片后再上传
- 或修改后端限制（article.go）:
```go
const maxFileSize = 10 << 20 // 改为 10MB
```

---

### 5. 数据库迁移问题

#### 新增字段后查询失败

**原因**: 数据库表结构与模型不匹配

**解决方案**:
```sql
-- 手动添加缺失字段
ALTER TABLE user ADD COLUMN IF NOT EXISTS email_verified INT DEFAULT 0;
ALTER TABLE user ADD COLUMN IF NOT EXISTS reset_token VARCHAR(255) NULL;
ALTER TABLE user ADD COLUMN IF NOT EXISTS reset_expires DATETIME NULL;

-- 或使用 Beego ORM 自动同步
-- backend/main.go 中添加:
orm.RunSyncdb("default", false, true)
```

---

## 📊 性能优化建议

### 前端优化

1. **按需加载组件**
```javascript
// 路由懒加载
const routes = [
  {
    path: '/articles',
    component: () => import('@/views/Articles.vue')
  }
]
```

2. **图片懒加载**
```vue
<img v-lazy="article.cover_image" alt="封面" />
```

3. **虚拟滚动长列表**
```bash
npm install vue-virtual-scroller
```

---

### 后端优化

1. **数据库索引优化**
```sql
-- 为常用查询字段添加索引
CREATE INDEX idx_article_created_at ON article(created_at);
CREATE INDEX idx_article_status ON article(status);
CREATE INDEX idx_comment_article_id ON comment(article_id);
```

2. **查询缓存**
```go
// 使用 Beego 缓存
cache.Set("article_list", articles, 300 * time.Second)
```

3. **连接池调优**
```ini
# backend/conf/app.conf
dbmaxidle = 20
dbmaxopen = 200
```

---

## 🚀 部署指南

### 生产环境部署

#### 1. 前端构建

```bash
cd frontend

# 设置生产环境变量
echo "VITE_API_BASE_URL=https://api.yourblog.com/api" > .env.production

# 构建
npm run build

# 构建产物在 dist/ 目录
```

#### 2. 后端编译

```bash
cd backend

# 交叉编译（Linux 服务器）
GOOS=linux GOARCH=amd64 go build -o blog

# 或直接编译当前系统
go build -o blog
```

#### 3. Nginx 配置示例

```nginx
server {
    listen 80;
    server_name yourblog.com;

    # 前端静态文件
    location / {
        root /var/www/blog/dist;
        try_files $uri $uri/ /index.html;
    }

    # 后端 API 代理
    location /api/ {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }

    # 上传文件
    location /uploads/ {
        alias /var/www/blog/backend/uploads/;
    }
}
```

#### 4. 使用 systemd 管理后端服务

```bash
# 创建服务文件
sudo vim /etc/systemd/system/blog.service
```

```ini
[Unit]
Description=Blog System Backend
After=network.target mysql.service

[Service]
Type=simple
User=www-data
WorkingDirectory=/var/www/blog/backend
ExecStart=/var/www/blog/backend/blog
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

```bash
# 启动服务
sudo systemctl daemon-reload
sudo systemctl start blog
sudo systemctl enable blog

# 查看状态
sudo systemctl status blog
```

---

## 📚 开发规范

### Git 提交规范

```bash
# 功能开发
git commit -m "feat: 添加文章收藏功能"

# Bug 修复
git commit -m "fix: 修复登录时验证码显示异常"

# 文档更新
git commit -m "docs: 更新 API 接口文档"

# 代码优化
git commit -m "refactor: 重构用户认证逻辑"

# 样式调整
git commit -m "style: 优化移动端导航栏样式"
```

### 分支管理

```bash
# 主分支
main          # 生产环境代码
develop       # 开发分支

# 功能分支
feature/*     # 新功能
bugfix/*      # Bug 修复
hotfix/*      # 紧急修复
```

---

## 🤝 贡献指南

欢迎提交 Issue 和 Pull Request！

### 提交 Issue

请提供以下信息：
1. 问题描述
2. 复现步骤
3. 预期行为
4. 实际行为
5. 环境信息（操作系统、浏览器、Node.js 版本等）
6. 截图或录屏（如适用）

### 提交 PR

1. Fork 本项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

---

## 📄 许可证

MIT License

---

## 🙏 致谢

感谢以下开源项目：

- [Vue.js](https://vuejs.org/)
- [Element Plus](https://element-plus.org/)
- [Beego](https://beego.me/)
- [Vditor](https://vditor.js.org/)
- [WangEditor](https://www.wangeditor.com/)

---

**最后更新日期**: 2026-04-01
