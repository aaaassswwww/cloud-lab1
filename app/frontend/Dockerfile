# 多阶段构建 - 构建阶段
FROM golang:1.21.13-bullseye AS builder

# 设置工作目录为项目根目录
WORKDIR /workspace

# 复制整个项目
COPY . .

# 进入 frontend 服务目录
WORKDIR /workspace/app/frontend

# 编译服务
RUN bash build.sh

# 运行阶段 - 使用更小的基础镜像
FROM debian:bullseye-slim

# 安装必要的运行时依赖
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# 设置工作目录
WORKDIR /app

# 从构建阶段复制编译好的文件
COPY --from=builder /workspace/app/frontend/output /app

# 暴露服务端口
EXPOSE 8080

# 启动服务
CMD ["bash", "bootstrap.sh"]
