# XzBBS Go 版 — 宝塔面板 Go 项目管理器部署指南

## 📋 目录

1. [环境准备](#一环境准备)
2. [安装宝塔 Go 项目管理器](#二安装宝塔-go-项目管理器)
3. [上传项目代码](#三上传项目代码)
4. [创建 MySQL 数据库](#四创建-mysql-数据库可选)
5. [修改配置文件](#五修改配置文件)
6. [添加 Go 项目](#六添加-go-项目)
7. [构建前端](#七构建前端)
8. [添加网站（Nginx 反向代理）](#八添加网站nginx-反向代理)
9. [配置 HTTPS（推荐）](#九配置-https推荐)
10. [验证部署](#十验证部署)
11. [更新升级](#十一更新升级)
12. [常见问题 FAQ](#十二常见问题-faq)

---

## 一、环境准备

### 服务器要求

| 项目     | 最低配置  | 推荐配置     |
| -------- | --------- | ------------ |
| 操作系统 | CentOS 7+ | Ubuntu 22.04 |
| 内存     | 512MB     | 1GB ~ 2GB    |
| 硬盘     | 2GB       | 10GB+        |
| 宝塔版本 | 7.9+      | 8.x          |

### 宝塔面板已安装

确保你已经安装了宝塔面板，如果没有安装：

```bash
# CentOS
yum install -y wget && wget -O install.sh https://download.bt.cn/install/install_6.0.sh && bash install.sh ed8484bec

# Ubuntu/Debian
wget -O install.sh https://download.bt.cn/install/install-6.0.sh && bash install.sh ed8484bec
```

---

## 二、安装宝塔 Go 项目管理器

### 2.1 安装 Go 管理器

1. 登录宝塔面板
2. 左侧菜单 → **软件商店**
3. 搜索 **"Go"** 或 **"Go项目管理器"**
4. 找到 **Go 项目管理器** → 点击 **安装**
5. 等待安装完成

### 2.2 安装 Go SDK

1. 安装完成后，点击 Go 项目管理器右侧的 **设置**
2. 进入 **版本管理** 或 **SDK版本管理**
3. 安装 **Go 1.22.x** 或更高版本（与项目 go.mod 中的版本一致）
4. 等待下载安装完成

> 💡 项目要求 Go 1.22+，建议安装 1.22.5 或更高版本。

---

## 三、上传项目代码

### 方式一：Git 克隆（推荐）

宝塔面板 → **终端**：

```bash
cd /www/wwwroot
git clone https://github.com/你的用户名/XzBBS.git XzBBS
cd XzBBS
ls -la
```

### 方式二：宝塔文件管理器上传

1. 本地将项目文件夹打包为 `XzBBS.zip`
2. 宝塔 → **文件** → 进入 `/www/wwwroot/`
3. 点击 **上传**，上传 zip 包
4. 右键 zip 包 → **解压**
5. 解压后确认文件夹名为 `XzBBS`

### 方式三：FTP 上传

使用 FileZilla 等工具将整个项目上传到 `/www/wwwroot/XzBBS/`

### 确认目录结构

上传完成后，目录结构应该是这样的：

```
/www/wwwroot/XzBBS/
├── go.mod
├── go.sum
├── config.yaml
├── cmd/
│   └── server/
│       └── main.go    ← 入口文件
├── internal/
├── web/
│   ├── package.json
│   └── src/
└── ...
```

---

## 四、创建 MySQL 数据库（可选）

> 如果只是测试，可以用 SQLite，跳过这一节。生产环境推荐 MySQL。

宝塔面板 → **数据库** → **添加数据库**：

| 字段       | 示例值          | 说明             |
| ---------- | --------------- | ---------------- |
| 数据库名   | XzBBS        | 自己起名         |
| 用户名     | XzBBS_user   | 自己起名         |
| 密码       | 强密码          | 点击生成随机密码 |
| 访问权限   | 本地服务器      | 生产环境建议本地 |
| 字符集     | utf8mb4         | 建议 utf8mb4     |

**记录好数据库名、用户名、密码**，后面配置要用。

---

## 五、修改配置文件

宝塔 → **文件** → 进入 `/www/wwwroot/XzBBS/` → 右键 `config.yaml` → **编辑**

### 关键配置说明

```yaml
server:
  port: 8080              # 后端端口（Go项目管理器里也要填这个）
  mode: release           # ⚠️ 生产环境必须改成 release
  trusted_proxies:        # 可信代理（Nginx 的 IP）
    - "127.0.0.1"

database:
  driver: mysql           # mysql 或 sqlite，生产推荐 mysql
  # MySQL 连接串（把下面的 用户名/密码/数据库名 换成你自己的）
  dsn: "数据库用户名:数据库密码@tcp(127.0.0.1:3306)/数据库名?charset=utf8mb4&parseTime=True&loc=Local"
  # 用 SQLite 测试的话改成下面这样：
  # driver: sqlite
  # dsn: "./XzBBS.db"

redis:
  enabled: false          # 没有 Redis 保持 false
  addr: "127.0.0.1:6379"
  password: ""
  db: 0

jwt:
  secret: "改成你自己的随机长字符串"  # ⚠️ 务必修改！安全起见
  expire_hour: 72         # 登录有效期（小时）

upload:
  path: "./uploads"
  max_size: 10            # 单文件最大 MB
  allow_ext:
    - ".jpg"
    - ".jpeg"
    - ".png"
    - ".gif"
    - ".webp"
    - ".zip"
    - ".rar"
    - ".pdf"

site:
  name: "我的论坛"         # 网站名称
  brief: "一个轻量级论坛"   # 网站简介
  page_size: 20            # 每页条数

email:
  enabled: false          # 是否启用邮件
  host: "smtp.qq.com"
  port: 465
  username: "your@email.com"
  password: "授权码"
  from: "我的论坛 <your@email.com>"
```

> 💡 **生成随机 JWT Secret**：
> 宝塔终端执行：`openssl rand -hex 32`

保存配置文件。

---

## 六、添加 Go 项目

### 6.1 打开 Go 项目管理器

宝塔面板 → **网站** → 顶部选择 **Go项目** 标签

### 6.2 添加项目

点击 **添加项目** 按钮，填写以下信息：

| 字段           | 填写值                                      | 说明                                         |
| -------------- | ------------------------------------------- | -------------------------------------------- |
| 项目名称       | `XzBBS`                                  | 随便起个名字，方便识别                       |
| 项目目录       | `/www/wwwroot/XzBBS`                     | 项目根目录                                   |
| 启动文件       | `cmd/server/main.go`                        | 入口文件路径（相对项目目录）                 |
| 输出文件名     | `XzBBS`                                  | 编译后的二进制文件名                         |
| 运行端口       | `8080`                                      | 和 config.yaml 里的 port 一致                |
| Go 版本        | 选择你安装的 Go 1.22.x 版本                 |                                              |
| 环境变量       | `GOPROXY=https://goproxy.cn,direct`         | 国内加速，必须加！                           |
| 编译参数       | `-ldflags="-s -w"`                          | 可选，减小二进制体积                         |
| 运行用户       | `www`                                       | 推荐用 www 用户                              |
| 开机自启       | ✅ 勾选                                      | 服务器重启后自动启动                         |

### 6.3 环境变量说明

在 **环境变量** 中添加：

```
GOPROXY=https://goproxy.cn,direct
CONFIG_PATH=/www/wwwroot/XzBBS/config.yaml
GIN_MODE=release
```

> ⚠️ **重要**：`GOPROXY` 一定要设置为国内镜像，否则依赖下载会失败！

### 6.4 启动项目

填写完成后点击 **提交** 或 **创建**。

项目创建后，会自动开始下载依赖并编译。等待一会儿，查看状态：

- 状态显示 **运行中** 表示启动成功
- 如果失败，点击 **日志** 查看错误信息

### 6.5 常用操作

在 Go 项目管理器中可以：

| 操作 | 说明 |
|------|------|
| **启动** | 启动项目 |
| **停止** | 停止项目 |
| **重启** | 重启项目（修改配置后需要重启） |
| **日志** | 查看运行日志（排错用） |
| **配置** | 修改项目配置 |
| **删除** | 删除项目 |

### 验证后端是否正常

宝塔终端执行：

```bash
curl http://127.0.0.1:8080/health
```

正常返回：
```json
{"status":"ok"}
```

---

## 七、构建前端

### 7.1 安装 Node.js 版本管理器

如果还没装 Node.js：

宝塔 → **软件商店** → 搜索 **"Node.js版本管理器"** → 安装

安装后打开，安装 **Node.js 18.x** 或更高版本。

### 7.2 构建前端

宝塔 → **终端**：

```bash
cd /www/wwwroot/XzBBS/web

# 安装依赖
npm install

# 构建生产版本
npm run build
```

构建成功后，会生成 `dist/` 目录：

```
web/dist/
├── index.html
└── assets/
    ├── index.xxx.js
    ├── index.xxx.css
    └── ...
```

> 💡 如果内存不足构建失败，可以在本地构建好后把 `dist/` 目录上传。

---

## 八、添加网站（Nginx 反向代理）

### 8.1 添加站点

宝塔面板 → **网站** → **添加站点**：

| 字段       | 填写值                           | 说明                     |
| ---------- | -------------------------------- | ------------------------ |
| 域名       | `bbs.example.com`                | 你的域名                 |
| 根目录     | `/www/wwwroot/XzBBS/web/dist` | 指向前端构建的 dist 目录 |
| FTP        | 不创建                           |                          |
| 数据库     | 不创建                           | 已单独创建               |
| PHP版本    | 纯静态                           | 项目不需要 PHP           |

点击 **提交**。

### 8.2 修改 Nginx 配置

站点创建后 → 点击右侧 **设置** → **配置文件**

将原有内容替换为以下（**注意替换 `bbs.example.com` 为你的域名**）：

```nginx
server
{
    listen 80;
    server_name bbs.example.com;
    
    root /www/wwwroot/XzBBS/web/dist;
    index index.html;

    # 前端 SPA：所有未知路径回退到 index.html
    location / {
        try_files $uri $uri/ /index.html;
    }

    # API 反向代理到 Go 后端
    location /api {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # 超时设置
        proxy_connect_timeout 60s;
        proxy_read_timeout 60s;
        proxy_send_timeout 60s;
    }

    # 上传文件由 Go 处理
    location /uploads {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    # 健康检查
    location /health {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
    }

    # 前端静态资源缓存（30天）
    location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2|ttf|eot)$ {
        expires 30d;
        add_header Cache-Control "public, immutable";
        access_log off;
    }

    # Gzip 压缩
    gzip on;
    gzip_min_length 1k;
    gzip_comp_level 5;
    gzip_types text/plain text/css application/json application/javascript text/xml application/xml application/rss+xml text/javascript image/svg+xml;
    gzip_vary on;

    # 上传大小限制
    client_max_body_size 20M;

    # 禁止访问隐藏文件
    location ~ /\. {
        deny all;
    }

    access_log  /www/wwwlogs/bbs.example.com.log;
    error_log   /www/wwwlogs/bbs.example.com.error.log;
}
```

点击 **保存**。

---

## 九、配置 HTTPS（推荐）

### 9.1 申请免费 SSL 证书

宝塔 → 网站设置 → **SSL** → **Let's Encrypt**：

1. 勾选你的域名
2. 填写邮箱
3. 点击 **申请**

申请成功后，打开右上角的 **强制 HTTPS**。

### 9.2 验证

浏览器访问 `https://你的域名`，地址栏应该显示小锁图标。

---

## 十、验证部署

### 10.1 访问网站

浏览器打开你的域名，应该能看到论坛首页。

### 10.2 测试注册登录

1. 点击 **注册**，创建一个新用户
2. 用新账号登录
3. 测试发帖、回复功能

### 10.3 默认管理员账号

| 用户名   | 密码       |
| -------- | ---------- |
| admin    | admin123   |

> ⚠️ **重要**：登录后请立即修改管理员密码！

### 10.4 健康检查

```bash
curl https://bbs.example.com/health
# 或
curl http://你的IP/health
```

---

## 十一、更新升级

### 方式一：通过宝塔 Go 项目管理器

1. **拉取最新代码**
   ```bash
   cd /www/wwwroot/XzBBS
   git pull
   ```

2. **重新构建前端**
   ```bash
   cd /www/wwwroot/XzBBS/web
   npm install
   npm run build
   ```

3. **重启 Go 项目**
   - 宝塔 → **网站** → **Go项目** → 找到 XzBBS → 点击 **重启**

### 方式二：手动更新

1. 上传新的代码文件覆盖
2. Go 项目管理器中点击 **重启**（会自动重新编译）
3. 前端重新 `npm run build`

---

## 十二、常见问题 FAQ

### Q1: Go 项目启动失败，日志显示依赖下载失败

**原因**：没有设置 GOPROXY 国内镜像

**解决**：
1. Go 项目管理器 → 找到项目 → 点击 **配置**
2. 在环境变量中添加：`GOPROXY=https://goproxy.cn,direct`
3. 保存后重启项目

### Q2: 访问网站显示 502 Bad Gateway

**原因**：Go 后端没有正常运行

**排查**：
1. 宝塔 → **网站** → **Go项目** → 查看项目状态是否为 **运行中**
2. 点击 **日志** 查看错误信息
3. 终端执行 `curl http://127.0.0.1:8080/health` 测试

### Q3: 前端页面空白

**原因**：前端构建失败或路径配置错误

**排查**：
- 确认 `web/dist/index.html` 文件存在
- 浏览器按 F12 打开开发者工具，看 Console 和 Network 报错
- 确认 Nginx 的 root 路径正确
- 确认 Nginx 配置了 `try_files $uri $uri/ /index.html;`

### Q4: 上传文件报错 413 Request Entity Too Large

**原因**：Nginx 上传大小限制

**解决**：
在 Nginx 配置的 server 块中添加：
```nginx
client_max_body_size 20M;
```
然后保存并重载 Nginx。

### Q5: 数据库连接失败

**排查**：
1. 宝塔 → **数据库** → 确认数据库状态是 **运行中**
2. 检查 `config.yaml` 中的 dsn 连接串是否正确
3. 确认数据库名、用户名、密码无误
4. Go 项目管理器 → 重启项目 → 看日志

### Q6: 修改配置后不生效

**原因**：Go 服务需要重启才能加载新配置

**解决**：
宝塔 → **网站** → **Go项目** → 找到 XzBBS → 点击 **重启**

### Q7: 宝塔 Go 项目管理器在哪里？

**回答**：
宝塔面板 → 左侧菜单 **网站** → 顶部有一排标签（PHP项目、Java项目、Node项目、Go项目...）→ 点击 **Go项目**

如果没有看到 Go项目 标签，说明还没安装 Go 项目管理器，去 **软件商店** 搜索安装。

### Q8: 启动文件路径怎么填？

**回答**：
填写入口文件相对于项目目录的路径。本项目的入口是 `cmd/server/main.go`，所以：
- 项目目录：`/www/wwwroot/XzBBS`
- 启动文件：`cmd/server/main.go`

### Q9: 内存不足编译失败

**解决**：
- 方案一：增加服务器内存（推荐 1GB+）
- 方案二：添加 swap 分区：
  ```bash
  dd if=/dev/zero of=/swapfile bs=1M count=1024
  mkswap /swapfile
  swapon /swapfile
  echo '/swapfile none swap defaults 0 0' >> /etc/fstab
  ```
- 方案三：本地 Windows 编译好 Linux 二进制再上传：
  ```powershell
  # Windows PowerShell 中执行
  set GOOS=linux
  set GOARCH=amd64
  set CGO_ENABLED=0
  go build -ldflags="-s -w" -o XzBBS ./cmd/server
  ```
  然后把生成的 `XzBBS` 二进制文件上传到服务器项目目录，在 Go 项目管理器中直接用这个二进制启动。

---

## 📁 部署完成后的目录结构

```
/www/wwwroot/XzBBS/
├── config.yaml           # 配置文件
├── XzBBS              # 编译后的 Go 二进制（Go管理器生成）
├── uploads/              # 上传文件目录
│   └── avatars/
├── cmd/                  # 源码
├── internal/             # 源码
├── web/
│   └── dist/             # 前端构建产物（Nginx root）
│       ├── index.html
│       └── assets/
└── ...
```

---

## 🔒 安全检查清单

部署完成后，建议逐项检查：

- [ ] 修改 JWT Secret 为随机字符串
- [ ] 修改默认管理员密码（admin / admin123）
- [ ] 开启 HTTPS 并强制跳转
- [ ] 宝塔面板修改默认端口（8888）
- [ ] 宝塔面板设置强密码
- [ ] 数据库不对外开放（仅本地访问）
- [ ] 定期备份数据库和 uploads 目录
- [ ] 关闭不必要的端口（8080 不对外开放）
- [ ] 服务器防火墙已配置

---

## 🎯 部署流程总结

```
1. 软件商店安装 Go 项目管理器
2. Go 项目管理器安装 Go 1.22+ SDK
3. 上传项目代码到 /www/wwwroot/XzBBS
4. 修改 config.yaml 配置
5. Go 项目管理器 → 添加项目 → 启动
6. 构建前端（npm run build）
7. 网站 → 添加站点 → 配置 Nginx 反代
8. 申请 SSL 证书 → 强制 HTTPS
9. 完成！🎉
```

---

祝部署顺利！🎉
