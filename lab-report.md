# 3.1 MySQL 数据持久化

## 修改内容

**修改前**：
```yaml
mysql:
  volumes:
    - ./db/sql/ini:/docker-entrypoint-initdb.d
```

**修改后**：
```yaml
mysql:
  volumes:
    - ./db/sql/ini:/docker-entrypoint-initdb.d
    - mysql-data:/var/lib/mysql              # 添加持久化卷

volumes:
  mysql-data:                                 # 定义命名卷
```

## 数据持久化原理阐述

**关键原理**：通过卷映射将容器内的数据目录映射到 Docker 外部管理的卷中。

| 项目 | 改之前 | 改之后 |
|------|--------|--------|
| **数据存储位置** | 容器内部 `/var/lib/mysql` | 外部卷 `mysql-data` |
| **数据位置** | 容器内 | **容器外**（Docker 管理） |
| **容器删除** | 数据丢失 | 卷保留，数据保存 |
| **重建容器** | 空白数据库 | 自动挂载卷，数据恢复 |

### 改之前（无持久化）
```
容器内 /var/lib/mysql (数据库文件)
           ↓
容器被删除 → 所有数据丢失
```

### 改之后（有持久化）
```
容器内 /var/lib/mysql ← 映射 ← 外部卷 mysql-data (物理存储)
           ↓                           ↓
    容器被删除                    卷保留，数据保存
           ↓
新容器启动 → 自动挂载卷 → 数据恢复 
```



# 3.3 Docker 微服务容器化
## 构建微服务方法
```bash
docker build --no-cache -f app/cart/Dockerfile -t gomall-cart:latest .
docker build --no-cache -f app/checkout/Dockerfile -t gomall-checkout:latest .
docker build --no-cache -f app/email/Dockerfile -t gomall-email:latest .
docker build --no-cache -f app/frontend/Dockerfile -t gomall-frontend:latest .
docker build --no-cache -f app/order/Dockerfile -t gomall-order:latest .
docker build --no-cache -f app/payment/Dockerfile -t gomall-payment:latest .
docker build --no-cache -f app/product/Dockerfile -t gomall-product:latest .
docker build --no-cache -f app/user/Dockerfile -t gomall-user:latest .
```
## 启动微服务方法
```bash
docker run -d --name cart --network gomall -p 8883:8883 --env-file ./app/cart/.env gomall-cart:latest
docker run -d --name checkout --network gomall -p 8884:8884 --env-file ./app/checkout/.env gomall-checkout:latest
docker run -d --name email --network gomall -p 8885:8885 --env-file ./app/email/.env gomall-email:latest
docker run -d --name frontend --network gomall -p 8080:8080 --env-file ./app/frontend/.env gomall-frontend:latest
docker run -d --name order --network gomall -p 8886:8886 --env-file ./app/order/.env gomall-order:latest
docker run -d --name payment --network gomall -p 8887:8887 --env-file ./app/payment/.env gomall-payment:latest
docker run -d --name product --network gomall -p 8888:8888 --env-file ./app/product/.env gomall-product:latest
docker run -d --name user --network gomall -p 8889:8889 --env-file ./app/user/.env gomall-user:latest
```

## Dockerfile设计

下面对 `app/cart/Dockerfile` 的设计逻辑做一个清晰且简单的说明：

- 核心思路：多阶段构建（multi-stage build）
  - 构建阶段（builder）：
    - 基于官方 `golang:1.21.13-bullseye` 镜像，设置工作目录并复制项目源码。
    - 切换到服务目录 `/workspace/app/cart`，运行项目内的 `build.sh` 脚本完成编译，产出可执行文件和其它构建产物到 `output` 目录。
  - 运行阶段（runtime）：
    - 基于更小的运行时镜像 `debian:bullseye-slim`，只安装必要运行时依赖（例如 `ca-certificates`）。
    - 将构建阶段生成的 `output` 整体复制到运行镜像的 `/app`，避免在运行镜像中包含源码或构建工具。
    - 暴露服务端口 `8883`，并通过 `CMD ["bash", "bootstrap.sh"]` 启动服务（假定 `bootstrap.sh` 在 `output` 中并负责运行最终二进制或初始化工作）。