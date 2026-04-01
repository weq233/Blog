# 🚀 博客系统 - 5 分钟快速启动指南

> 适合快速体验项目的开发者，详细文档请查看 [README.md](./README.md)

---

## ⚡ 三步启动（前提：已安装 Node.js、Go、MySQL）

### 第一步：准备数据库 (1 分钟)

```bash
# 1. 登录 MySQL


# 2. 创建数据库（数据库在backend\database里）


# 3. 导入表结构


# 4. 修改数据库密码（编辑配置文件）
# 打开 backend/conf/app.conf，找到第 9 行
# 将 "root:123456" 改为你的 MySQL 密码
```

---

### 第二步：启动后端 (1 分钟)

```bash
# 终端 1 - 进入后端目录
cd backend

# 安装依赖（首次运行需要）
go mod download

# 启动服务
bee run

# ✅ 看到 "✅ 数据库连接正常" 表示成功
# 后端地址：http://localhost:8080
```

**常见错误**:
- ❌ `command not found: bee` → 执行 `go install github.com/beego/bee/v2@latest`
- ❌ `数据库连接失败` → 检查 MySQL 是否启动，密码是否正确

---

### 第三步：启动前端 (1 分钟)

```bash
# 终端 2 - 进入前端目录
cd frontend

# 安装依赖（首次运行需要，约 1-2 分钟）
npm install

# 启动开发服务器
npm run dev

# ✅ 看到 "Local: http://localhost:5173/" 表示成功
# 前端地址：http://localhost:5173
```

**常见错误**:
- ❌ `Cannot find module` → 删除 `node_modules` 后重新 `npm install`
- ❌ `端口被占用` → 执行 `npm run dev -- --port 3000`

---

## 🎉 完成！

打开浏览器访问：**http://localhost:5173**

### 默认管理员账号

```
用户名：admin
密码：admin123
```

---

## 🔧 常用操作

### 注册新用户

1. 访问 http://localhost:5173/register
2. 填写用户名、邮箱、密码
3. 输入验证码提交

### 发布第一篇文章

1. 登录后点击右上角头像 → 创作中心
2. 点击"写文章"按钮
3. 输入标题和内容
4. 选择分类和标签
5. 点击"发布"

### 进入管理后台

1. 使用管理员账号登录
2. 访问 http://localhost:5173/admin
3. 管理文章、用户、分类

---

## 🛠️ 故障排查

### 后端无法启动

```bash
# 检查 MySQL 服务
# Windows
net start | findstr MySQL

# macOS
brew services list

# 启动 MySQL
# Windows: net start MySQL80
# macOS: brew services start mysql
```

### 前端无法启动

```bash
# 清理缓存
rm -rf node_modules package-lock.json
npm install

# 检查 Node.js 版本（需 >= 18）
node -v
```

### API 请求失败

```bash
# 确认后端已启动
curl http://localhost:8080

# 检查跨域配置
# 查看 backend/conf/app.conf 中 cors_enabled = true
```

---

## 📝 下一步

- 📖 查看完整文档：[README.md](./README.md)
- 🔌 查看 API 接口：[README.md#api-接口文档](./README.md#api-接口文档)
- ⚙️ 自定义配置：[README.md#配置说明](./README.md#配置说明)
- 🚀 生产部署：[README.md#部署指南](./README.md#部署指南)

---

**遇到问题？** 查看 [常见问题](./README.md#常见问题) 或提交 Issue

**最后更新**: 2026-04-01
