# XiaoMiao Home 小喵回家后端项目

<div align="center">

![版本](https://img.shields.io/badge/版本-1.0.0-blue)
![Go版本](https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go)
![授权](https://img.shields.io/badge/授权-MIT-green)

</div>

## 📑 项目介绍

XiaoMiao Home 小喵回家后端项目，基于 [Kratos](https://go-kratos.dev/) 框架构建，提供喵咪领养、寻喵、问答等能力，并集成用户、角色、权限与通知等基础能力。

### 主要功能

| 模块       | 说明                                       |
| ---------- | ------------------------------------------ |
| **用户**   | 用户 CRUD、登录/登出、状态、密钥、菜单权限 |
| **角色**   | 角色 CRUD、状态管理                        |
| **用户组** | 用户组 CRUD、状态管理                      |
| **通知**   | 通知的创建、批量操作、列表、状态管理       |

### 技术栈

- **框架**: Kratos v2（HTTP + gRPC）
- **数据库**: MySQL（GORM）
- **缓存**: Redis
- **配置/注册**: Nacos
- **依赖注入**: Wire
- **API**: Protobuf + gRPC-Gateway，生成 OpenAPI 规范

## 🚀 快速开始

### 前置要求

- Go 1.23 或更高版本
- MySQL、Redis（本地或远程）
- Nacos（可选，用于配置与注册）
- Docker（可选，用于容器化部署）

### 本地开发

1. **安装代码生成与构建工具**

```bash
make init
```

2. **生成 API 与校验代码**

```bash
make api
```

3. **生成校验相关代码**

```bash
make validate
```

4. **生成配置与 Wire 等**

```bash
make config
make generate
```

5. **一键生成所有必要文件**

```bash
make all
```

6. **编译运行**

```bash
make build
# 运行前请将配置文件放到 /data/conf 或通过 -conf 指定目录
./bin/server -conf /data/conf
```

### 端口说明

- **8000**: HTTP API
- **9000**: gRPC 服务

## 🐳 Docker 部署

### 构建镜像

```bash
docker build -t xiaomiao-home-system:latest .
```

### 运行容器

```bash
docker run --rm -p 8000:8000 -p 9000:9000 \
  -v /path/to/your/configs:/data/conf \
  xiaomiao-home-system:latest
```

镜像内通过 `-conf /data/conf` 加载配置，请将实际配置文件挂载到 `/data/conf`。

## 📂 项目结构

```
.
├── api/              # API 定义（Proto）及生成的 Go/OpenAPI 代码
│   ├── role/         # 角色 API
│   └── user/         # 用户、用户组、通知 API
├── cmd/              # 应用入口与 Wire
│   └── xiaomiao-home-system/
├── internal/         # 内部实现
│   ├── biz/          # 业务逻辑
│   ├── conf/         # 配置结构（Proto）
│   ├── server/       # HTTP/gRPC 服务
│   ├── service/      # 对外服务实现
│   └── task/         # 定时/异步任务
├── third_party/      # 第三方 Proto（validate、openapi 等）
├── utils/            # 通用工具（加密、Nacos、存储等）
├── Makefile          # 构建与代码生成
├── openapi.yaml      # 生成的 OpenAPI 规范
└── Dockerfile        # 镜像构建
```

## 🔧 配置说明

应用通过 `-conf` 指定配置目录（默认示例为 `/data/conf`），目录内需包含 Kratos 所需的 Bootstrap 配置（如 `config.yaml` 等）。配置结构定义在 `internal/conf/conf.proto`，主要包括：

- **Server**: HTTP / gRPC 监听地址与超时
- **Data**: Database（MySQL）、Redis 连接与超时
- **Auth**: 服务密钥、API 密钥
- **Config**: 业务相关（如 dayu_url、api_key）
- **Static**: 静态资源目录与访问地址

请根据实际环境修改数据库连接串、Redis 地址、Nacos 等配置。

## 📝 开发规范

- 遵循 Go 官方编码规范
- 提交前建议运行测试与 `make build`
- API 变更后需执行 `make api` 与 `make validate`，必要时 `make all`
- 使用语义化版本管理

## 📄 许可证

MIT License
