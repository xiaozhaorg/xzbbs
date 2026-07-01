# XzBBS (Go Rewrite)

轻量级论坛系统，基于 Go + Vue 3 前后端分离架构重写。

## 特性

- 🚀 Go 后端，高性能、低资源占用
- 🎨 Vue 3 + TailwindCSS 前端，响应式设计
- 🔐 JWT 认证，无状态
- 📦 单二进制部署，内嵌前端资源
- 🗃️ 支持 MySQL / SQLite
- 📎 文件上传、图片管理
- 👮 版主系统：置顶、关闭、移动、封禁
- 🛡️ 限流、CORS、bcrypt 密码加密

## 快速开始

### 开发模式

```bash
# 后端
go run ./cmd/server

# 前端（另一个终端）
cd web && npm run dev
```

### 生产部署

```bash
# 构建
make build

# 运行
./bin/XzBBS
```

### Docker

```bash
docker build -t XzBBS .
docker run -p 8080:8080 XzBBS
```

## 配置

编辑 `config.yaml`：

```yaml
server:
  port: 8080
  mode: release

database:
  driver: mysql
  dsn: "user:password@tcp(127.0.0.1:3306)/XzBBS?charset=utf8mb4&parseTime=True&loc=Local"

jwt:
  secret: "your-secure-random-string"
  expire_hour: 72
```

## 默认账号

- 用户名: `admin`
- 密码: `admin123`

⚠️ 首次部署请立即修改管理员密码！

## API 文档

| Method | Endpoint                | Description |
| ------ | ----------------------- | ----------- |
| POST   | /api/auth/register      | 注册          |
| POST   | /api/auth/login         | 登录          |
| GET    | /api/auth/me            | 当前用户        |
| GET    | /api/forums             | 版块列表        |
| GET    | /api/forums/:id/threads | 版块帖子列表      |
| POST   | /api/threads            | 发帖          |
| GET    | /api/threads/:id        | 帖子详情        |
| POST   | /api/threads/:tid/posts | 回帖          |
| POST   | /api/attachments        | 上传附件        |
| PUT    | /api/mod/threads/top    | 置顶          |
| GET    | /api/admin/stats        | 管理统计        |

完整 API 见源码 `cmd/server/main.go` 路由定义。

## 项目结构

```
cmd/server/         入口
internal/
  config/           配置加载
  model/            数据模型 (GORM)
  repository/       数据访问层
  service/          业务逻辑层
  handler/          HTTP 处理器
  middleware/       中间件 (JWT, CORS, 限流)
  dto/              请求/响应结构
  pkg/              工具包
web/                Vue 3 前端
```

## License

MIT
