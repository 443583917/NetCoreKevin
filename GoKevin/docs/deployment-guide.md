# GoKevin 运行和部署指南

**文档版本**: v1.0
**创建日期**: 2026/06/07
**最后更新**: 2026/06/07

---

## 目录

- [1. 环境要求](#1-环境要求)
- [2. 本地开发环境搭建](#2-本地开发环境搭建)
- [3. 后端运行](#3-后端运行)
- [4. 前端运行](#4-前端运行)
- [5. Docker 部署](#5-docker-部署)
- [6. 生产环境部署](#6-生产环境部署)
- [7. 配置说明](#7-配置说明)
- [8. 常见问题](#8-常见问题)

---

## 1. 环境要求

### 1.1 后端环境要求

| 软件 | 版本要求 | 说明 |
|------|----------|------|
| Go | 1.22+ | 主要开发语言 |
| MySQL | 8.0+ | 主数据库 |
| Redis | 7.0+ | 缓存和会话存储 |
| RabbitMQ | 3.x | 消息队列（可选） |
| Qdrant | 1.7+ | 向量数据库（可选） |

### 1.2 前端环境要求

| 软件 | 版本要求 | 说明 |
|------|----------|------|
| Node.js | 18+ | JavaScript 运行时 |
| npm | 8+ | 包管理器 |

### 1.3 开发工具推荐

| 工具 | 说明 |
|------|------|
| GoLand / VS Code | Go 开发 IDE |
| WebStorm / VS Code | 前端开发 IDE |
| Navicat / DBeaver | 数据库管理工具 |
| Postman / Apifox | API 测试工具 |
| Docker Desktop | 容器化部署 |

---

## 2. 本地开发环境搭建

### 2.1 安装 Go

#### Windows

1. 下载 Go 安装包：https://golang.org/dl/
2. 运行安装包，按提示完成安装
3. 配置环境变量：

```bash
# 添加到 PATH
C:\Go\bin

# 设置 GOPATH
set GOPATH=%USERPROFILE%\go
set PATH=%PATH%;%GOPATH%\bin
```

#### macOS

```bash
# 使用 Homebrew
brew install go

# 或下载安装包
# https://golang.org/dl/
```

#### Linux

```bash
# 下载并解压
wget https://golang.org/dl/go1.22.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz

# 配置环境变量
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
echo 'export GOPATH=$HOME/go' >> ~/.bashrc
echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.bashrc
source ~/.bashrc
```

### 2.2 安装 Node.js

#### Windows

1. 下载 Node.js 安装包：https://nodejs.org/
2. 运行安装包，按提示完成安装
3. 验证安装：

```bash
node --version
npm --version
```

#### macOS

```bash
# 使用 Homebrew
brew install node
```

#### Linux

```bash
# 使用 NodeSource
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt-get install -y nodejs
```

### 2.3 安装 MySQL

#### Windows

1. 下载 MySQL 安装包：https://dev.mysql.com/downloads/installer/
2. 运行安装包，按提示完成安装
3. 配置 root 密码
4. 创建数据库：

```sql
CREATE DATABASE kevin_app CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

#### macOS

```bash
# 使用 Homebrew
brew install mysql
brew services start mysql

# 创建数据库
mysql -u root -p
CREATE DATABASE kevin_app CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

#### Linux

```bash
# Ubuntu/Debian
sudo apt update
sudo apt install mysql-server
sudo systemctl start mysql
sudo systemctl enable mysql

# 创建数据库
sudo mysql
CREATE DATABASE kevin_app CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
ALTER USER 'root'@'localhost' IDENTIFIED WITH mysql_native_password BY 'your_password';
FLUSH PRIVILEGES;
```

### 2.4 安装 Redis

#### Windows

1. 下载 Redis：https://github.com/microsoftarchive/redis/releases
2. 解压并运行 redis-server.exe

#### macOS

```bash
brew install redis
brew services start redis
```

#### Linux

```bash
sudo apt update
sudo apt install redis-server
sudo systemctl start redis
sudo systemctl enable redis
```

### 2.5 安装 RabbitMQ（可选）

#### 使用 Docker（推荐）

```bash
docker run -d --name rabbitmq \
  -p 5672:5672 \
  -p 15672:15672 \
  -e RABBITMQ_DEFAULT_USER=guest \
  -e RABBITMQ_DEFAULT_PASS=guest \
  rabbitmq:3-management
```

#### 本地安装

```bash
# macOS
brew install rabbitmq
brew services start rabbitmq

# Linux
sudo apt install rabbitmq-server
sudo systemctl start rabbitmq-server
```

---

## 3. 后端运行

### 3.1 克隆项目

```bash
cd /path/to/your/workspace
git clone https://github.com/your-org/GoKevin.git
cd GoKevin/backend
```

### 3.2 安装依赖

```bash
go mod download
```

### 3.3 配置数据库

编辑 `configs/config.yaml`：

```yaml
server:
  port: 9901
  mode: debug
  readTimeout: 30
  writeTimeout: 30

database:
  host: 127.0.0.1
  port: 3306
  user: root
  password: your_password
  dbname: kevin_app
  maxIdleConns: 10
  maxOpenConns: 100

redis:
  host: 127.0.0.1
  port: 6379
  password: ""
  db: 0

jwt:
  secret: your-secret-key-here
  expireHour: 24

ai:
  defaultProvider: openai
  providers:
    - name: openai
      apiKey: ${OPENAI_API_KEY}
      baseUrl: https://api.openai.com/v1
      model: gpt-4

rabbitmq:
  host: 127.0.0.1
  port: 5672
  user: guest
  password: guest
  vhost: /

log:
  level: info
  filename: logs/app.log
  maxSize: 100
  maxBackups: 7
  maxAge: 30
```

### 3.4 运行数据库迁移

```bash
# 使用 Makefile
make migrate

# 或直接运行
go run cmd/server/main.go
```

### 3.5 启动后端服务

```bash
# 使用 Makefile
make run

# 或直接运行
go run cmd/server/main.go

# 或编译后运行
make build
./bin/go-kevin
```

### 3.6 验证后端服务

```bash
# 访问 API
curl http://localhost:9901/api/v1/auth/login \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{"userName":"admin","password":"123456"}'

# 访问 Swagger 文档
# 浏览器打开 http://localhost:9901/swagger/index.html
```

---

## 4. 前端运行

### 4.1 进入前端目录

```bash
cd GoKevin/frontend
```

### 4.2 安装依赖

```bash
npm install
```

### 4.3 配置环境变量

编辑 `.env.development`：

```env
# API 基础URL
VUE_APP_BASE_API=http://localhost:9901

# 应用标题
VUE_APP_TITLE=GoKevin

# 其他配置
VUE_APP_VERSION=1.0.0
```

### 4.4 启动前端服务

```bash
# 开发模式
npm run serve

# 或使用 yarn
yarn serve
```

### 4.5 访问前端应用

浏览器打开：http://localhost:8080

### 4.6 构建生产版本

```bash
# 构建
npm run build

# 构建产物在 dist/ 目录
```

---

## 5. Docker 部署

### 5.1 使用 Docker Compose 部署（推荐）

#### 5.1.1 进入后端目录

```bash
cd GoKevin/backend
```

#### 5.1.2 配置环境变量

创建 `.env` 文件：

```env
# 数据库配置
MYSQL_ROOT_PASSWORD=admin123
MYSQL_DATABASE=kevin_app

# Redis 配置
REDIS_PASSWORD=123456

# RabbitMQ 配置
RABBITMQ_DEFAULT_USER=guest
RABBITMQ_DEFAULT_PASS=guest

# 应用配置
CONFIG_PATH=/app/configs/config.prod.yaml
```

#### 5.1.3 启动所有服务

```bash
# 使用 Makefile
make docker-run

# 或直接使用 docker-compose
docker-compose up -d
```

#### 5.1.4 查看服务状态

```bash
docker-compose ps
```

#### 5.1.5 查看日志

```bash
# 查看所有服务日志
docker-compose logs -f

# 查看特定服务日志
docker-compose logs -f app
docker-compose logs -f mysql
docker-compose logs -f redis
```

#### 5.1.6 停止服务

```bash
# 使用 Makefile
make docker-stop

# 或直接使用 docker-compose
docker-compose down
```

### 5.2 单独构建 Docker 镜像

#### 5.2.1 构建后端镜像

```bash
cd GoKevin/backend

# 使用 Makefile
make docker-build

# 或直接构建
docker build -t go-kevin-backend .
```

#### 5.2.2 运行后端容器

```bash
docker run -d \
  --name go-kevin-backend \
  -p 9901:9901 \
  -v $(pwd)/configs:/app/configs \
  -v $(pwd)/logs:/app/logs \
  go-kevin-backend
```

### 5.3 Docker Compose 配置说明

```yaml
# docker-compose.yaml
version: '3.8'

services:
  # 应用服务
  app:
    build: .
    ports:
      - "9901:9901"
    environment:
      - CONFIG_PATH=/app/configs/config.prod.yaml
    depends_on:
      mysql:
        condition: service_healthy
      redis:
        condition: service_started
    volumes:
      - ./configs:/app/configs
      - ./logs:/app/logs
    restart: unless-stopped

  # MySQL 数据库
  mysql:
    image: mysql:8.0
    ports:
      - "3306:3306"
    environment:
      MYSQL_ROOT_PASSWORD: admin123
      MYSQL_DATABASE: kevin_app
    volumes:
      - mysql_data:/var/lib/mysql
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost"]
      interval: 10s
      timeout: 5s
      retries: 3

  # Redis 缓存
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    command: redis-server --requirepass 123456
    volumes:
      - redis_data:/data

  # RabbitMQ 消息队列
  rabbitmq:
    image: rabbitmq:3-management
    ports:
      - "5672:5672"
      - "15672:15672"
    environment:
      RABBITMQ_DEFAULT_USER: guest
      RABBITMQ_DEFAULT_PASS: guest
    volumes:
      - rabbitmq_data:/var/lib/rabbitmq

volumes:
  mysql_data:
  redis_data:
  rabbitmq_data:
```

---

## 6. 生产环境部署

### 6.1 部署架构

```
                    ┌─────────────┐
                    │   Nginx     │
                    │ (反向代理)   │
                    └──────┬──────┘
                           │
            ┌──────────────┼──────────────┐
            │              │              │
     ┌──────▼──────┐ ┌────▼────┐ ┌───────▼───────┐
     │   前端静态   │ │  后端   │ │   WebSocket   │
     │   文件服务   │ │  API    │ │    服务       │
     └─────────────┘ └────┬────┘ └───────────────┘
                           │
            ┌──────────────┼──────────────┐
            │              │              │
     ┌──────▼──────┐ ┌────▼────┐ ┌───────▼───────┐
     │   MySQL     │ │  Redis  │ │   RabbitMQ    │
     │   数据库    │ │  缓存   │ │   消息队列    │
     └─────────────┘ └─────────┘ └───────────────┘
```

### 6.2 服务器环境准备

#### 6.2.1 系统要求

| 资源 | 最低配置 | 推荐配置 |
|------|----------|----------|
| CPU | 2核 | 4核+ |
| 内存 | 4GB | 8GB+ |
| 硬盘 | 50GB | 100GB+ |
| 操作系统 | Ubuntu 20.04/CentOS 7 | Ubuntu 22.04 |

#### 6.2.2 安装 Docker

```bash
# Ubuntu
sudo apt update
sudo apt install docker.io
sudo systemctl start docker
sudo systemctl enable docker

# 安装 Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
```

#### 6.2.3 配置防火墙

```bash
# 开放端口
sudo ufw allow 80/tcp    # HTTP
sudo ufw allow 443/tcp   # HTTPS
sudo ufw allow 9901/tcp  # 后端 API（可选，通过 Nginx 代理）
```

### 6.3 部署步骤

#### 6.3.1 克隆代码

```bash
cd /opt
sudo git clone https://github.com/your-org/GoKevin.git
cd GoKevin
```

#### 6.3.2 配置生产环境

```bash
# 编辑后端配置
cd backend
cp configs/config.yaml configs/config.prod.yaml
vim configs/config.prod.yaml
```

配置生产环境数据库密码：

```yaml
database:
  host: mysql
  port: 3306
  user: root
  password: ${MYSQL_ROOT_PASSWORD}
  dbname: kevin_app
```

#### 6.3.3 构建并启动服务

```bash
# 构建并启动
docker-compose up -d

# 查看服务状态
docker-compose ps
```

#### 6.3.4 配置 Nginx

安装 Nginx：

```bash
sudo apt install nginx
```

创建 Nginx 配置文件：

```nginx
# /etc/nginx/sites-available/gokevin
server {
    listen 80;
    server_name your-domain.com;

    # 前端静态文件
    location / {
        root /opt/GoKevin/frontend/dist;
        try_files $uri $uri/ /index.html;
    }

    # 后端 API 代理
    location /api/ {
        proxy_pass http://localhost:9901;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # WebSocket 代理
    location /ws/ {
        proxy_pass http://localhost:9901;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
    }

    # 静态资源缓存
    location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg)$ {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }
}
```

启用配置：

```bash
sudo ln -s /etc/nginx/sites-available/gokevin /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx
```

#### 6.3.5 配置 HTTPS（可选）

```bash
# 安装 Certbot
sudo apt install certbot python3-certbot-nginx

# 获取 SSL 证书
sudo certbot --nginx -d your-domain.com

# 自动续期
sudo certbot renew --dry-run
```

### 6.4 数据备份

#### 6.4.1 数据库备份

```bash
# 创建备份脚本
cat > /opt/GoKevin/scripts/backup.sh << 'EOF'
#!/bin/bash
BACKUP_DIR="/opt/backups/mysql"
DATE=$(date +%Y%m%d_%H%M%S)
mkdir -p $BACKUP_DIR

# 备份数据库
docker exec gokevin-mysql-1 mysqldump -u root -p${MYSQL_ROOT_PASSWORD} kevin_app > $BACKUP_DIR/kevin_app_$DATE.sql

# 保留最近 7 天的备份
find $BACKUP_DIR -name "*.sql" -mtime +7 -delete
EOF

chmod +x /opt/GoKevin/scripts/backup.sh

# 添加定时任务
crontab -e
# 添加：0 2 * * * /opt/GoKevin/scripts/backup.sh
```

#### 6.4.2 文件备份

```bash
# 备份上传文件
tar -czf /opt/backups/files/files_$(date +%Y%m%d).tar.gz /opt/GoKevin/backend/uploads
```

### 6.5 监控和日志

#### 6.5.1 查看应用日志

```bash
# 查看后端日志
docker-compose logs -f app

# 查看系统日志
tail -f /var/log/nginx/access.log
tail -f /var/log/nginx/error.log
```

#### 6.5.2 健康检查

```bash
# 检查后端服务
curl http://localhost:9901/api/v1/health

# 检查数据库连接
docker exec gokevin-mysql-1 mysqladmin ping

# 检查 Redis 连接
docker exec gokevin-redis-1 redis-cli ping
```

---

## 7. 配置说明

### 7.1 后端配置文件

配置文件位于 `backend/configs/config.yaml`：

```yaml
# 服务器配置
server:
  port: 9901                    # 监听端口
  mode: debug                   # 运行模式：debug/release
  readTimeout: 30               # 读取超时（秒）
  writeTimeout: 30              # 写入超时（秒）

# 数据库配置
database:
  host: 127.0.0.1               # 数据库地址
  port: 3306                    # 数据库端口
  user: root                    # 数据库用户名
  password: your_password       # 数据库密码
  dbname: kevin_app             # 数据库名称
  maxIdleConns: 10              # 最大空闲连接数
  maxOpenConns: 100             # 最大打开连接数

# Redis 配置
redis:
  host: 127.0.0.1               # Redis 地址
  port: 6379                    # Redis 端口
  password: ""                  # Redis 密码
  db: 0                         # Redis 数据库编号

# JWT 配置
jwt:
  secret: your-secret-key-here  # JWT 密钥
  expireHour: 24                # Token 过期时间（小时）

# AI 配置
ai:
  defaultProvider: openai       # 默认 AI 提供者
  providers:
    - name: openai              # 提供者名称
      apiKey: ${OPENAI_API_KEY} # API 密钥
      baseUrl: https://api.openai.com/v1  # API 基础 URL
      model: gpt-4              # 默认模型

# RabbitMQ 配置
rabbitmq:
  host: 127.0.0.1               # RabbitMQ 地址
  port: 5672                    # RabbitMQ 端口
  user: guest                   # RabbitMQ 用户名
  password: guest               # RabbitMQ 密码
  vhost: /                      # RabbitMQ 虚拟主机

# 日志配置
log:
  level: info                   # 日志级别：debug/info/warn/error
  filename: logs/app.log        # 日志文件路径
  maxSize: 100                  # 单个日志文件最大大小（MB）
  maxBackups: 7                 # 保留旧日志文件最大数量
  maxAge: 30                    # 保留旧日志文件最大天数
```

### 7.2 前端配置文件

#### .env.development

```env
# 开发环境配置
VUE_APP_BASE_API=http://localhost:9901
VUE_APP_TITLE=GoKevin (Dev)
VUE_APP_VERSION=1.0.0
```

#### .env.production

```env
# 生产环境配置
VUE_APP_BASE_API=https://api.your-domain.com
VUE_APP_TITLE=GoKevin
VUE_APP_VERSION=1.0.0
```

### 7.3 环境变量

可以通过环境变量覆盖配置文件中的值：

```bash
# 数据库配置
export DB_HOST=localhost
export DB_PORT=3306
export DB_USER=root
export DB_PASSWORD=your_password
export DB_NAME=kevin_app

# Redis 配置
export REDIS_HOST=localhost
export REDIS_PORT=6379
export REDIS_PASSWORD=

# JWT 配置
export JWT_SECRET=your-secret-key
export JWT_EXPIRE_HOUR=24

# AI 配置
export OPENAI_API_KEY=sk-xxx

# 应用配置
export CONFIG_PATH=/path/to/config.yaml
```

---

## 8. 常见问题

### 8.1 后端启动失败

#### 问题：数据库连接失败

```
错误信息：failed to connect database
```

**解决方案：**

1. 检查 MySQL 服务是否启动
2. 检查数据库配置是否正确
3. 检查数据库用户权限
4. 检查防火墙设置

```bash
# 检查 MySQL 状态
systemctl status mysql

# 测试数据库连接
mysql -h 127.0.0.1 -P 3306 -u root -p
```

#### 问题：端口被占用

```
错误信息：bind: address already in use
```

**解决方案：**

1. 查找占用端口的进程
2. 停止该进程或修改配置端口

```bash
# 查找占用端口的进程
lsof -i :9901

# 停止进程
kill -9 <PID>
```

### 8.2 前端启动失败

#### 问题：npm install 失败

```
错误信息：npm ERR! code ERESOLVE
```

**解决方案：**

1. 清除 npm 缓存
2. 使用 `--legacy-peer-deps` 参数

```bash
# 清除缓存
npm cache clean --force

# 重新安装
npm install --legacy-peer-deps
```

#### 问题：接口请求失败

```
错误信息：Network Error
```

**解决方案：**

1. 检查后端服务是否启动
2. 检查 API 地址配置
3. 检查 CORS 配置

```bash
# 检查后端服务
curl http://localhost:9901/api/v1/health
```

### 8.3 Docker 部署问题

#### 问题：容器启动失败

```bash
# 查看容器日志
docker-compose logs app

# 检查容器状态
docker-compose ps
```

**解决方案：**

1. 检查配置文件是否正确
2. 检查依赖服务是否启动
3. 检查端口是否被占用

#### 问题：数据库初始化失败

```bash
# 进入数据库容器
docker exec -it gokevin-mysql-1 mysql -u root -p

# 手动创建数据库
CREATE DATABASE kevin_app CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 8.4 性能优化

#### 8.4.1 数据库优化

```sql
-- 添加索引
ALTER TABLE t_user ADD INDEX idx_tenant_id (tenant_id);
ALTER TABLE t_ai_apps ADD INDEX idx_tenant_id (tenant_id);
ALTER TABLE t_ai_chat_sessions ADD INDEX idx_user_id (user_id);

-- 优化查询
EXPLAIN SELECT * FROM t_user WHERE tenant_id = 1000;
```

#### 8.4.2 Redis 优化

```bash
# 检查 Redis 性能
redis-cli info stats

# 监控 Redis 命令
redis-cli monitor
```

#### 8.4.3 应用优化

```yaml
# 调整数据库连接池
database:
  maxIdleConns: 20
  maxOpenConns: 200

# 调整 Redis 连接池
redis:
  poolSize: 100
```

### 8.5 安全配置

#### 8.5.1 修改默认密码

```yaml
# 数据库密码
database:
  password: StrongPassword123!

# Redis 密码
redis:
  password: StrongRedisPassword!

# JWT 密钥
jwt:
  secret: VeryLongAndRandomSecretKeyHere
```

#### 8.5.2 配置 CORS

```go
// 限制允许的源
func CORS() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("Access-Control-Allow-Origin", "https://your-domain.com")
        c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
        // ...
    }
}
```

#### 8.5.3 配置 HTTPS

```nginx
server {
    listen 443 ssl;
    server_name your-domain.com;

    ssl_certificate /etc/letsencrypt/live/your-domain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/your-domain.com/privkey.pem;

    # 强制 HTTPS
    if ($scheme != "https") {
        return 301 https://$host$request_uri;
    }

    # ... 其他配置
}
```

---

## 附录

### A. Makefile 命令参考

```bash
# 构建
make build

# 运行
make run

# 测试
make test

# 单元测试
make test-unit

# 集成测试
make test-integration

# 测试覆盖率
make test-coverage

# 代码检查
make lint

# 清理
make clean

# Docker 构建
make docker-build

# Docker 运行
make docker-run

# Docker 停止
make docker-stop

# 数据库迁移
make migrate

# 生成文档
make docs
```

### B. 常用 Docker 命令

```bash
# 查看容器
docker ps

# 查看日志
docker logs <container_id>

# 进入容器
docker exec -it <container_id> /bin/sh

# 停止容器
docker stop <container_id>

# 删除容器
docker rm <container_id>

# 查看镜像
docker images

# 删除镜像
docker rmi <image_id>
```

### C. 常用 Git 命令

```bash
# 克隆仓库
git clone https://github.com/your-org/GoKevin.git

# 创建分支
git checkout -b feature/your-feature

# 提交更改
git add .
git commit -m "feat: add your feature"

# 推送分支
git push origin feature/your-feature

# 合并分支
git checkout master
git merge feature/your-feature

# 更新代码
git pull origin master
```

---

**文档版本**: v1.0
**最后更新**: 2026/06/07
**维护者**: NetCoreKevin 开发团队
