**❗️作者声明：本项目为研究学习区块链的开源项目，不提供任何形式的收费服务(谨防诈骗)
，不鼓励任何衍生金融属性的交易行为，不负责任何使用本项目进行的三方行为；使用过程中遇见问题请提`issue`
或群里交流，开源项目请自重！**

---

# BEpusdt (Better Easy Payment USDT)

<p align="center">
<img src="./static/payment/assets/img/tether.svg" width="15%" alt="tether">
</p>
<p align="center">
<a href="https://www.gnu.org/licenses/gpl-3.0.html"><img src="https://img.shields.io/badge/license-GPLV3-blue" alt="license GPLV3"></a>
<a href="https://golang.org"><img src="https://img.shields.io/badge/Golang-1.24-red" alt="Go version 1.24"></a>
<a href="https://github.com/gin-gonic/gin">
    <img src="https://img.shields.io/badge/Gin-v1.10-blue" alt="Gin Web Framework">
</a>
<a href="https://github.com/v03413/bepusdt">
  <img src="https://img.shields.io/github/v/release/v03413/bepusdt" alt="GitHub Release">
</a>
<a href="https://hub.docker.com/r/v03413/bepusdt">
    <img src="https://img.shields.io/docker/pulls/v03413/bepusdt?style=flat-square&logo=docker" alt="Docker Pulls">
</a>
</p>

## 🪧 介绍

缘起于`Epusdt`，但不仅限于此，加入了一些新特性，致力于成为一款更好用的个人`加密货币`收款网关。

## 🎉 新特性

### 🌟 目前支持收款网络

🔥 主流网络：TRON Ethereum BSC Polygon<br>
⚡ 其他网络： X-Layer Solana Aptos Arbitrum-One Base [完整列表](./docs/trade-type.md)

- ✅ 完全兼容 `Epusdt` 插件无缝替换
- ️✅ 支持主流区块网络 不仅限于`USDT`
- ✅ 支持主流法定货币 汇率自动更新
- ✅ 轻依赖 单体二进制文件 部署便捷
- ✅ 支持非订单交易监控 余额变动通知
- ✅ 支持自定义支付精度与递增颗粒度
- ✅ 底层区块扫描 安全确认 速度稳定
- ✅ 支持波场能量代理 回收实时监控
- ✅ 原生兼容`易支付`收单 接入便捷
- ✅ 完整独立WEB后台 方便配置管理
- ✅ 收银台支持中英适配 助力出海需求
- ✅ 快速迭代期 超多实用特性功能等你发现

## 🚀 快捷启动

Docker 快速启动，执行完命令打开地址`http://服务器IP:8080`即可看到初始页面

```bash  
docker run -d --restart=unless-stopped -p 8080:8080 v03413/bepusdt:latest
```
## 📃 技术文档  

- 安装：[Docker](docs/docker/docker.md) [Linux](docs/linux/install.md) [1Panel](./docs/1panel/README.md) [宝塔](./docs/bt_panel/README.md)
- 开发：[API对接](docs/api/api.md) [订单回调](docs/notify/readme.md) [Python](https://github.com/luoyanglang/bepusdt-python-sdk) [PHP](https://github.com/v03413/bepusdt-php-sdk)
- 对接：[独角数卡](docs/dujiaoka/dujiaoka.md) [异次元](./docs/acg-faka.md) [萌次元](./docs/mcy-shop.md) [彩虹易支付](https://github.com/v03413/Epay-BEpusdt) [WHMCS](https://github.com/v03413/whmcs-gateway-epusdt)
- 其它：[https 配置](./docs/ssl.md) [时钟同步](docs/linux/systemd-timesyncd.md) [收银台修改](docs/payment-template/README.md)

## 🖼 功能截图

| 前台收银                                            | 后台订单                                             | Telegram 通知                                            |
|-------------------------------------------------|--------------------------------------------------|--------------------------------------------------------|
| <img src=./docs/images/1.png alt=收银台 width=300> | <img src=./docs/images/2.png alt=后台订单 width=300> | <img src=./docs/images/3.png alt=Telegram通知 width=300> |

## ❓ 常见问题

- [服务器配置性能选型推荐 ⚡️](./docs/faq/server.md)
- [服务器流量消耗说明解释](./docs/faq/bandwidth.md)
- [后台入口账密忘记重置教程](./docs/faq/login-reset.md)
- [Telegram 通知 Chat ID 获取教程](docs/faq/telegram-chat-id.md)
- [推荐配置提高 Tron 扫块稳定性‼️](./docs/tron-grid/readme.md)

## ⚠️ 重要提示

- **订单交易强依赖时间**：请确保服务器时间准确，否则可能导致订单异常
- **网络环境要求**：请确保服务器网络环境稳定，否则可能影响功能正常运行

## 🏝️ 社区交流

- **Telegram 群组**：[https://t.me/BEpusdtChat](https://t.me/BEpusdtChat)
- **Telegram 频道**：[https://t.me/BEpusdtChannel](https://t.me/BEpusdtChannel)

## 🙏 致谢

- [Epusdt](https://github.com/assimon/epusdt)

## 🌟 Star 历史

[![Stargazers over time](https://starchart.cc/v03413/bepusdt.svg)](https://starchart.cc/v03413/bepusdt)
