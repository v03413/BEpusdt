# Docker 部署指南

## 环境准备

### 前置要求

1. **Docker 环境：** 确保服务器已安装 Docker（建议使用最新稳定版本）
2. **服务器推荐配置：**
    - **操作系统：** Debian 12+ 或其他主流 Linux 发行版
    - **性能要求：** 无特殊硬性要求，建议至少 1GB 内存
    - **网络要求：** 确保服务器网络通畅，能够稳定访问区块链网络
    - **云服务商：** 推荐使用知名云服务提供商（如 AWS、GCP、DigitalOcean 等），以确保服务稳定性和安全性

## 快速部署

### 基础部署命令

```bash
docker run -d --restart=unless-stopped -p 8080:8080 v03413/bepusdt:latest
```

### 参数说明

- **端口映射 `-p 8080:8080`：** 将容器内部的 8080 端口映射到宿主机的 8080 端口，支持通过宿主机 IP 地址访问服务
- **自动重启 `--restart=unless-stopped`：** 确保容器在异常退出或系统重启后自动恢复运行
- **镜像版本：**
    - `v03413/bepusdt:latest` - 稳定发行版（推荐生产环境使用）
    - `v03413/bepusdt:nightly` - 每日构建开发版（包含最新特性，适合测试环境）

## 数据持久化配置

### SQLite 数据库（默认）

默认情况下，系统使用 SQLite 数据库存储所有数据，数据保存在容器内部路径 `/var/lib/bepusdt`。

**重要提示：** 为避免容器删除导致数据丢失，强烈建议将数据目录挂载到宿主机：

```bash
docker run -d --restart=unless-stopped -p 8080:8080 -v [挂载路径]:/var/lib/bepusdt v03413/bepusdt:latest
```

**配置示例：**

将 `[挂载路径]` 替换为宿主机实际路径，例如：

```bash
docker run -d --restart=unless-stopped -p 8080:8080 -v /opt/bepusdt:/var/lib/bepusdt v03413/bepusdt:latest
```

### MySQL 数据库（可选）

系统支持使用 MySQL 作为数据存储方案，适用于需要更强大数据库功能的场景。

**前置准备：**

- MySQL 版本：5.7 或更高版本
- 需提前完成 MySQL 安装及安全配置
- 已创建专用数据库（如 `bepusdt`）

**部署命令：**

```bash
docker run -d \
  --restart unless-stopped \
  -p 8080:8080 \
  -e MYSQL_DSN=user:password@tcp(127.0.0.1:3306)/bepusdt?charset=utf8mb4&parseTime=True&loc=Local&timeout=3s&readTimeout=10s&writeTimeout=10s \
  v03413/bepusdt:latest
```

**配置说明：**

- 将 `user` 替换为 MySQL 用户名
- 将 `password` 替换为 MySQL 密码
- 将 `127.0.0.1:3306` 替换为 MySQL 服务器地址和端口
- 将 `bepusdt` 替换为实际数据库名称（如需要）

## 部署验证

### 查看容器日志

执行部署命令后，系统将返回容器 ID。使用以下命令查看启动日志：

```bash
docker logs -f [容器ID]
```

如果启动成功，日志将显示相应的成功信息。

### 访问管理界面

通过浏览器访问以下地址进入初始化页面：

```
http://[服务器IP]:8080/[安全入口]
```

首次访问时，系统将显示安装向导。请按照提示完成初始化配置，并妥善保存管理员账号信息。

![安装完成](./init.png)

## 运维管理

### 重置管理员密码

如忘记管理员密码，可使用以下命令重置：

```bash
docker exec -it [容器ID] bepusdt reset
```

### 重要提示

**强烈推荐配置 TronGrid API Key：** 配置后可显著提升 TRON 网络扫块的稳定性和可靠性。

详细配置步骤请参考：[TronGrid 配置指南](../tron-grid/readme.md)

---