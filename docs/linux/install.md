# Linux 手动安装

> ⚠️ **前置说明**：本教程假定您已掌握 Linux 基本操作，包括命令行使用、文件管理等常识。

## 系统要求

| 项目        | 要求               |
|-----------|------------------|
| **操作系统**  | Debian 12 或更高版本  |
| **处理器架构** | amd64（其他架构请自行测试） |

## 安装步骤

执行以下命令依次完成安装：

### 1. 下载安装包

```bash
# 下载对应架构版本
wget -O ./bepusdt.zip https://github.com/v03413/bepusdt/releases/latest/download/bepusdt-linux-amd64.zip
```

### 2. 解压文件

```bash
unzip ./bepusdt.zip
```

解压后的目录结构如下：

```
./bepusdt
├── bepusdt              # 可执行文件
└── bepusdt.service     # systemd 服务配置文件
```

### 3. 配置系统自启

```bash
# 将 service 文件移动到 systemd 目录
mv ./bepusdt/bepusdt.service /etc/systemd/system

# 启用服务自启
systemctl enable bepusdt.service
```

### 4. 启动服务

```bash
# 启动 BEpusdt 服务
systemctl start bepusdt.service
```

### 5. 验证启动状态

```bash
# 查看服务状态
systemctl status bepusdt.service
```

✅ **成功标志**：看到 `Active: active (running)` 表示服务已成功启动

---

## 常用命令

| 操作       | 命令                                  |
|----------|-------------------------------------|
| **查看状态** | `systemctl status bepusdt.service`  |
| **查看日志** | `journalctl -u bepusdt.service -f`  |
| **重启服务** | `systemctl restart bepusdt.service` |
| **停止服务** | `systemctl stop bepusdt.service`    |