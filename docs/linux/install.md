# Linux 手动安装指南

> ⚠️ **前置条件**：本指南假设您具备 Linux 基本操作能力，包括命令行使用、文件系统管理等基础知识。

## 系统需求

| 项目         | 要求               |
|------------|------------------|
| **操作系统**   | Debian 11 或更高版本  |
| **CPU 架构** | amd64（其他架构请自行验证） |

## 安装步骤

### 1. 获取安装包

从 [GitHub Releases](https://github.com/v03413/bepusdt/releases/latest/) 页面下载与系统架构对应的安装包。

```bash
wget -O ./linux-amd64-BEpusdt.tar.gz https://github.com/v03413/BEpusdt/releases/latest/download/linux-amd64-BEpusdt.tar.gz
```

### 2. 解压安装包

```bash
tar -zxvf ./linux-amd64-BEpusdt.tar.gz
```

解压后的目录结构：

```
./bepusdt
├── bepusdt              # 可执行程序文件
└── bepusdt.service     # systemd 服务配置文件
```

### 3. 系统集成与自启配置

```bash
# 复制可执行文件至系统路径
mv ./bepusdt/bepusdt /usr/local/bin/

# 设置可执行权限
chmod +x /usr/local/bin/bepusdt

# 复制服务配置文件
mv ./bepusdt/bepusdt.service /etc/systemd/system/

# 启用开机自启
systemctl enable bepusdt.service
```

### 4. 启动服务

```bash
systemctl start bepusdt.service
```

### 5. 验证服务状态

```bash
systemctl status bepusdt.service
```

✅ **成功指标**：状态显示 `Active: active (running)` 表示服务已正常启动

---

## 常用操作命令

| 操作       | 命令                                  |
|----------|-------------------------------------|
| **查看状态** | `systemctl status bepusdt.service`  |
| **查看日志** | `journalctl -u bepusdt.service -f`  |
| **重启服务** | `systemctl restart bepusdt.service` |
| **停止服务** | `systemctl stop bepusdt.service`    |