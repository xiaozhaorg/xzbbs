# XzBBS Go 版 — 宝塔面板部署指南

## 一、环境要求

| 项目   | 要求                                     |
| ---- | -------------------------------------- |
| 系统   | CentOS 7+ / Ubuntu 20.04+ / Debian 11+ |
| 宝塔面板 | 7.x 或 8.x                              |
| 内存   | ≥ 256MB（推荐 512MB+）                     |
| 存储   | ≥ 500MB 可用空间                           |

---

## 二、宝塔面板安装步骤

### 2.1 安装宝塔面板（如未安装）

```bash
# CentOS
yum install -y wget && wget -O install.sh https://download.bt.cn/install/install_6.0.sh && bash install.sh ed8484bec

# Ubuntu
wget -O install.sh https://download.bt.cn/install/install_6.0.sh && bash install.sh ed8484bec
```

安装完成后登录宝塔面板。

### 2.2 安装依赖

登录宝塔面板 → **软件商店** → 搜索安装：

| 软件    | 版本要求  | 用途                   |
| ----- | ----- | -------------------- |
| Nginx | 1.22+ | 反向代理 / 前端静态服务        |
| MySQL | 8.0+  | 数据库（如果用 SQLite 则不需要） |
| Redis | 6.x+  | 可选缓存                 |

> 如果用 SQLite 开发测试，可以跳过 MySQL。

### 2.3 安装 Go 环境

**方式一：宝塔软件商店（推荐）**

- 宝塔 → **软件商店** → 搜索 **Go** → 安装（会自动配置环境变量）

**方式二：手动安装**

```bash
# 下载 Go 1.22
cd /root
wget https://go.dev/dl/go1.22.5.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.22.5.linux-amd64.tar.gz

# 配置环境变量
echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
echo 'export GOPATH=/root/go' >> /etc/profile
echo 'export PATH=$PATH:$GOPATH/bin' >> /etc/profile
source /etc/profile

# 验证
go version
```

---

## 三、上传项目代码

### 方式一：宝塔文件管理器

1. 宝塔面板 → **文件**
2. 进入 `/www/wwwroot/` 目录
3. 点击 **上传**，将 `XzBBS` 文件夹上传到服务器

### 方式二：Git（推荐）

```bash
cd /www/wwwroot
git clone <你的仓库地址> XzBBS
cd XzBBS
```

### 方式三：FTP / SCP

将本地项目压缩为 zip 后上传，再解压：

```bash
cd /www/wwwroot
unzip XzBBS.zip
mv XzBBS-master XzBBS   # 根据实际目录名调整
```

---

## 四、配置后端

### 4.1 编辑配置文件

```bash
cd /www/wwwroot/XzBBS
vi config.yaml
```

根据你的服务器修改以下关键配置：

```yaml
server:
  port: 8080
  mode: release                    # 生产环境用 release

database:
  driver: mysql                    # 用 mysql 或 sqlite
  dsn: "用户名:密码@tcp(127.0.0.1:3306)/XzBBS?charset=utf8mb4&parseTime=True&loc=Local"
  # 如果用 SQLite: dsn: "./XzBBS.db"

redis:
  enabled: false                   # 没有 Redis 保持 false
  addr: "127.0.0.1:6379"

jwt:
  secret: "改为你自己的随机字符串"    # ⚠️ 务必修改！
  expire_hour: 72

upload:
  path: "./uploads"
  max_size: 10
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
  name: "你的论坛名称"
  brief: "论坛简介"
  page_size: 20
```

### 4.2 安装 Go 依赖

```bash
cd /www/wwwroot/XzBBS
go mod tidy
```

### 4.3 编译

```bash
# 直接编译
go build -o XzBBS ./cmd/server

# 或使用 Makefile
make build
```

编译完成后会在当前目录生成 `XzBBS` 可执行文件。

### 4.4 设置目录权限

```bash
cd /www/wwwroot/XzBBS

# 确保 uploads 目录可写
chown -R www:www uploads
chmod -R 755 uploads

# 确保二进制可执行
chmod +x XzBBS
```

---

## 五、配置 Systemd 守护进程（开机自启 + 后台运行）

创建服务文件：

```bash
vi /etc/systemd/system/XzBBS.service
```

写入以下内容（注意修改路径和用户名）：

```ini
[Unit]
Description=XzBBS Go Server
After=network.target mysql.service
Wants=mysql.service

[Service]
Type=simple
User=root
WorkingDirectory=/www/wwwroot/XzBBS
ExecStart=/www/wwwroot/XzBBS/XzBBS
Restart=always
RestartSec=5
Environment=CONFIG_PATH=/www/wwwroot/XzBBS/config.yaml

# 资源限制（可选）
LimitNOFILE=65535

[Install]
WantedBy=multi-user.target
```

启动服务：

```bash
systemctl daemon-reload
systemctl enable XzBBS      # 开机自启
systemctl start XzBBS       # 启动服务

# 查看状态
systemctl status XzBBS

# 查看日志
journalctl -u XzBBS -f
```

### 常用操作

```bash
systemctl restart XzBBS     # 重启
systemctl stop XzBBS        # 停止
journalctl -u XzBBS -n 100  # 最近100行日志
```

---

## 六、配置 Nginx 反向代理

### 方式一：API 分离 + 前端 SPA

宝塔面板 → **网站** → **添加站点** → 填写域名 → **确认**

然后进入站点设置 → **网站目录** → 运行目录改为 `/www/wwwroot/XzBBS/web/dist`（前端构建后）

进入 **配置文件**，替换为以下内容：

```nginx
server {
    listen 80;
    server_name 你的域名.com;
    root /www/wwwroot/XzBBS/web/dist;

    # 前端 SPA — 所有路由回退到 index.html
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
    }

    # 上传文件静态访问
    location /uploads {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
    }

    # 前端静态资源缓存
    location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2|ttf|eot)$ {
        expires 30d;
        add_header Cache-Control "public, immutable";
    }
}
```

### 方式二：纯 API + 独立前端

如果前端部署在 CDN 或单独的静态服务，Nginx 只需转发 API：

```nginx
server {
    listen 80;
    server_name api.你的域名.com;

    location /api {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }

    location /uploads {
        proxy_pass http://127.0.0.1:8080;
    }
}
```

---

## 七、配置 HTTPS（推荐）

宝塔面板 → 你的站点 → **SSL** → 申请 Let's Encrypt 免费证书 → 启用

启用后，Nginx 会自动添加 443 监听。如果没有，手动添加：

```nginx
listen 443 ssl http2;
ssl_certificate /www/server/panel/vhost/cert/你的域名.com/fullchain.pem;
ssl_certificate_key /www/server/panel/vhost/cert/你的域名.com/privkey.pem;

# 强制 HTTP → HTTPS
server {
    listen 80;
    server_name 你的域名.com;
    return 301 https://$host$request_uri;
}
```

---

## 八、构建前端

### 8.1 在服务器上构建

```bash
cd /www/wwwroot/XzBBS/web

# 安装 Node.js（宝塔软件商店搜索安装 Node.js）
# 或使用 nvm
nvm install 20
nvm use 20

# 安装依赖
npm install

# 构建
npm run build
```

构建完成后，`web/dist/` 目录会包含编译后的前端文件。

### 8.2 也可以本地构建后上传

```bash
# 本地
cd web && npm install && npm run build

# 将 web/dist/ 上传到服务器的 /www/wwwroot/XzBBS/web/dist/
```

---

## 九、初始化数据库

### 方式一：MySQL（推荐生产环境）

宝塔面板 → **数据库** → **添加数据库**：

- 数据库名：`XzBBS`
- 用户名：自己设置
- 密码：自己设置

修改 `config.yaml` 中的数据库连接：

```yaml
database:
  driver: mysql
  dsn: "你设置的账号:密码@tcp(127.0.0.1:3306)/XzBBS?charset=utf8mb4&parseTime=True&loc=Local"
```

然后重启服务，启动时会自动创建表结构和初始数据：

```bash
systemctl restart XzBBS
```

查看日志确认：

```bash
journalctl -u XzBBS -f
# 应该看到 "🚀 XzBBS server starting on :8080"
```

### 方式二：SQLite（测试用）

```yaml
database:
  driver: sqlite
  dsn: "./XzBBS.db"
```

无需创建数据库，启动时自动生成文件。

---

## 十、宝塔防火墙配置

宝塔面板 → **安全** → 确保以下端口放行：

| 端口   | 用途            |
| ---- | ------------- |
| 80   | HTTP          |
| 443  | HTTPS         |
| 8080 | 内部代理（不需要对外开放） |

如果云服务器使用安全组（阿里云/腾讯云等），也需要在云控制台放行 80 和 443。

---

## 十一、性能优化建议

### 11.1 开启 Gzip

在宝塔站点的 Nginx 配置中添加：

```nginx
gzip on;
gzip_min_length 1k;
gzip_comp_level 5;
gzip_types text/plain text/css application/json application/javascript text/xml application/xml application/rss+xml text/javascript image/svg+xml;
gzip_vary on;
```

### 11.2 开启 Redis 缓存

```bash
# 宝塔安装 Redis 后，修改 config.yaml
redis:
  enabled: true
  addr: "127.0.0.1:6379"
```

### 11.3 Nginx 静态资源优化

```nginx
# 前端构建后的静态资源由 Nginx 直接服务（不走 Go）
location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff2?|ttf|eot)$ {
    root /www/wwwroot/XzBBS/web/dist;
    expires 30d;
    add_header Cache-Control "public, immutable";
    access_log off;
}
```

### 11.4 宝塔 PHP 清理

这个项目不需要 PHP，如果不需要可以停掉 PHP 服务节省内存。

---

## 十二、更新部署

项目更新后重新部署流程：

```bash
cd /www/wwwroot/XzBBS

# 1. 拉取最新代码（如果用 Git）
git pull

# 2. 更新依赖
go mod tidy

# 3. 重新编译
go build -o XzBBS ./cmd/server

# 4. 前端重新构建
cd web && npm install && npm run build && cd ..

# 5. 重启服务
systemctl restart XzBBS
```

---

## 十三、常见问题

### Q1: 编译时报错 `no Go installation found`

- 宝塔安装 Go 或手动安装，确认 `go version` 能输出版本号

### Q2: 启动后无法访问

- 检查 `systemctl status XzBBS` 是否有报错
- 检查 8080 端口是否被占用：`ss -tlnp | grep 8080`
- 检查 Nginx 配置是否正确：宝塔 → 网站 → 你的站点 → 配置文件

### Q3: 上传文件 413 错误

- 宝塔 → Nginx → 配置文件 → 添加 `client_max_body_size 20M;`

### Q4: 数据库连接失败

- 检查 MySQL 是否运行：`systemctl status mysqld`
- 检查账号密码是否正确
- 检查数据库是否已创建
- MySQL 8.0 注意密码插件兼容性

### Q5: 前端白屏

- 确认 `npm run build` 成功
- 确认 `web/dist/` 目录存在 `index.html`
- 检查 Nginx root 路径是否正确

### Q6: 502 Bad Gateway

- 后端是否正常运行：`systemctl status XzBBS`
- 8080 端口是否监听：`curl http://127.0.0.1:8080/health`
- Nginx proxy_pass 地址是否正确

---

## 十四、目录结构说明

```
/www/wwwroot/XzBBS/
├── config.yaml           ← 配置文件
├── XzBBS              ← 编译后的二进制
├── uploads/              ← 上传文件（avatar, attach）
│   └── avatars/
├── cmd/server/main.go    ← 源码
├── internal/             ← 后端源码
├── web/
│   └── dist/             ← 前端构建产物
│       ├── index.html
│       └── assets/
```

---

## 十五、安全清单

- [ ] 修改 JWT secret（必须）
- [ ] 修改默认管理员密码（必须：admin / admin123）
- [ ] 关闭宝塔面板外网访问（或修改端口）
- [ ] 定期备份数据库和 uploads 目录
- [ ] 开启 HTTPS
- [ ] 关闭不需要的服务（PHP 等）
- [ ] 配置 fail2ban 防暴力破解
