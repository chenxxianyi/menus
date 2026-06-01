# 今天吃什么 — 菜谱推荐系统

一个帮助用户解决"今天吃什么"难题的全栈应用，支持智能菜谱推荐、一周菜单规划、情侣点餐、购物清单生成等功能。

## 技术栈

### 后端

| 技术 | 版本 | 用途 |
|------|------|------|
| Go | 1.22 | 编程语言 |
| Gin | 1.10 | Web 框架 |
| GORM | 1.25 | ORM |
| MySQL | 8.0 | 数据库 |
| Redis | 7.x | 缓存 |
| JWT | - | 身份认证 |
| Zap | 1.27 | 日志 |

### 前端

| 技术 | 版本 | 用途 |
|------|------|------|
| Vue | 3.5 | 核心框架 |
| TypeScript | 6.0 | 类型系统 |
| Vite | 8.0 | 构建工具 |
| Pinia | 3.0 | 状态管理 |
| Vue Router | 4.6 | 路由 |
| Vant | 4.9 | 移动端 UI 组件库 |
| Tailwind CSS | 3.4 | 样式框架 |
| Axios | 1.16 | HTTP 请求 |

## 项目结构

```
menus/
├── backend/                    # Go 后端
│   ├── cmd/server/main.go      # 入口
│   ├── config/                 # 配置文件
│   │   ├── config.go
│   │   └── config.yaml
│   ├── internal/
│   │   ├── api/                # HTTP 处理器
│   │   │   ├── v1/             # 用户端接口
│   │   │   └── admin/          # 管理后台接口
│   │   ├── middleware/         # 中间件（认证、CORS、日志）
│   │   ├── model/              # 数据模型
│   │   ├── pkg/                # 公共工具（响应、错误码）
│   │   ├── repository/         # 数据访问层
│   │   ├── router/             # 路由注册
│   │   └── service/            # 业务逻辑层
│   └── migrations/             # 数据库迁移 SQL
│
├── web-user/                   # Vue 3 用户端
│   └── src/
│       ├── api/                # API 接口封装
│       ├── assets/             # 静态资源
│       ├── components/         # 公共组件
│       ├── router/             # 路由配置
│       ├── stores/             # Pinia 状态管理
│       ├── types/              # TypeScript 类型
│       └── views/              # 页面组件
│
├── design/                     # 设计稿
├── 前端开发任务文档.md
└── 后端开发任务文档.md
```

## 功能特性

### 核心功能

- **首页** — 今日推荐、热门菜谱、快捷入口
- **菜谱浏览** — 分类筛选、关键词搜索、菜谱详情（食材、步骤、营养信息）
- **智能推荐** — 根据用餐人数、口味偏好、饮食目标、已有食材推荐菜单
- **一周菜单** — 自动生成 7 天食谱规划
- **购物清单** — 根据菜单自动生成食材采购清单，支持勾选已购买
- **收藏系统** — 收藏喜欢的菜谱，个人中心查看收藏数

### 情侣点餐

- **邀请绑定** — 生成 6 位邀请码，对方输入即可配对
- **双向点餐** — 两人各自选择想吃的菜品，支持选餐次、日期、备注
- **点餐清单** — 按日期查看双方点餐，支持确认/取消/删除
- **合意菜单** — 合并两人已确认的点餐，自动生成食材购物清单
- **情侣昵称** — 自定义情侣关系昵称

### 用户系统

- 注册 / 登录（JWT 认证）
- 个人中心（头像、昵称、收藏统计）
- 口味偏好设置
- 意见反馈

### 管理后台

- 数据概览 Dashboard
- 菜谱管理（CRUD）
- 分类管理、食材管理
- 轮播图管理
- 用户管理、反馈管理

## 快速开始

### 环境要求

- Go >= 1.22
- Node.js >= 18
- MySQL >= 8.0
- Redis >= 7.x

### 1. 数据库

创建数据库：

```sql
CREATE DATABASE menu_recommend CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 2. 后端

```bash
cd backend

# 修改配置（数据库连接等）
# vim config/config.yaml

# 安装依赖
go mod tidy

# 启动（首次运行会自动建表）
go run ./cmd/server/main.go
```

后端默认运行在 `http://localhost:8000`。

### 3. 前端

```bash
cd web-user

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

前端默认运行在 `http://localhost:5199`。

## API 接口

### 公开接口

| 方法 | 路径 | 描述 |
|------|------|------|
| POST | `/api/auth/register` | 用户注册 |
| POST | `/api/auth/login` | 用户登录 |
| GET | `/api/home` | 首页数据 |
| GET | `/api/recipes` | 菜谱列表 |
| GET | `/api/recipes/:id` | 菜谱详情 |
| GET | `/api/categories` | 分类列表 |
| GET | `/api/banners` | 轮播图 |

### 用户接口（需登录）

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | `/api/user/info` | 用户信息 |
| PUT | `/api/user/profile` | 更新资料 |
| GET/PUT | `/api/user/preferences` | 偏好设置 |
| POST/DELETE | `/api/recipes/:id/favorite` | 收藏/取消收藏 |
| GET | `/api/user/favorites` | 收藏列表 |
| GET | `/api/user/favorites/count` | 收藏数量 |
| GET/POST/PUT/DELETE | `/api/shopping-list` | 购物清单 CRUD |
| POST | `/api/recommend/menu` | 智能菜单推荐 |
| POST | `/api/recommend/by-ingredients` | 按食材推荐 |
| POST | `/api/recommend/week-menu` | 一周菜单 |

### 情侣点餐接口（需登录）

| 方法 | 路径 | 描述 |
|------|------|------|
| GET | `/api/couple/invite-code` | 获取/生成邀请码 |
| POST | `/api/couple/bind` | 绑定情侣 |
| GET | `/api/couple/info` | 情侣信息 |
| POST | `/api/couple/unbind` | 解除绑定 |
| PUT | `/api/couple/name` | 设置情侣昵称 |
| POST | `/api/couple/orders` | 创建点餐 |
| GET | `/api/couple/orders` | 点餐列表 |
| PUT | `/api/couple/orders/:id` | 更新点餐状态 |
| DELETE | `/api/couple/orders/:id` | 删除点餐 |
| POST | `/api/couple/orders/generate-shopping-list` | 生成食材清单 |

### 管理后台接口

| 方法 | 路径 | 描述 |
|------|------|------|
| POST | `/api/admin/login` | 管理员登录 |
| GET | `/api/admin/dashboard` | 数据概览 |
| CRUD | `/api/admin/recipes` | 菜谱管理 |
| CRUD | `/api/admin/categories` | 分类管理 |
| CRUD | `/api/admin/ingredients` | 食材管理 |
| CRUD | `/api/admin/banners` | 轮播图管理 |
| GET | `/api/admin/users` | 用户列表 |
| GET/PUT | `/api/admin/feedback` | 反馈管理 |

## 配置说明

后端配置文件位于 `backend/config/config.yaml`：

```yaml
server:
  host: "0.0.0.0"
  port: 8000

database:
  host: localhost
  port: 3306
  user: root
  password: your_password
  dbname: menu_recommend

redis:
  host: localhost
  port: 6379

jwt:
  secret_key: your-secret-key    # 生产环境务必修改
  expire_hours: 24

cors:
  origins:
    - "http://localhost:5199"    # 前端地址
```

## 部署

### Docker Compose

```bash
cd backend
docker-compose up -d
```

### 生产环境

```bash
# 后端
cd backend
go build -o server ./cmd/server/main.go
./server

# 前端
cd web-user
npm run build
# 将 dist/ 目录部署到 Nginx
```
