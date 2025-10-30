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