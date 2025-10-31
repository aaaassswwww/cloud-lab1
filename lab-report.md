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
# 注意在 cloud-lab1 目录下执行
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