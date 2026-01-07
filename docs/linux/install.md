# Linux 手动安装

**默认认为您已经掌握 Linux 的操作常识，否则后续没有查看必要。**

准备服务器，Debian12+，架构目前只对amd64做了测试，其他架构请自行测试；依次执行以下命令：

```bash

# 下载对应架构版本
wget -O ./bepusdt.zip https://github.com/v03413/bepusdt/releases/latest/download/bepusdt-linux-amd64.zip

# 解压
unzip ./bepusdt.zip

# 解压结构如下：
#root@debian:~# tree ./bepusdt
#./bepusdt
#├── bepusdt
#├── bepusdt.service

# 配置软件自启
mv ./bepusdt/bepusdt.service /etc/systemd/system
systemctl enable bepusdt.service

# 启动软件
systemctl start bepusdt.service

# 查看软件状态（看到 Active: active (running) 即成功启动）
systemctl status bepusdt.service
```