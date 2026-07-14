**❗️作者声明：本项目为研究学习区块链的开源项目，不提供任何形式的收费服务(谨防诈骗)
，不鼓励任何衍生金融属性的交易行为，不负责任何使用本项目进行的三方行为；使用过程中遇见问题请提`issue`
或群里交流，开源项目请自重！**

---

# BEpusdt (Better Easy Payment USDT)

<p align="center">
<img src="./static/checkout/official/assets/img/bepusdt.svg" width="15%" alt="tether">
</p>
<p align="center">
<a href="https://www.gnu.org/licenses/gpl-3.0.html">
    <img src="https://img.shields.io/github/license/v03413/bepusdt" alt="license GPLV3">
</a>
<a href="https://github.com/v03413/bepusdt">
  <img src="https://img.shields.io/github/v/release/v03413/bepusdt" alt="GitHub Release">
</a>
<a href="https://github.com/v03413/bepusdt">
  <img src="https://img.shields.io/github/downloads/v03413/bepusdt/total" alt="GitHub Release">
</a>
<a href="https://hub.docker.com/r/v03413/bepusdt">
    <img src="https://img.shields.io/docker/pulls/v03413/bepusdt?style=flat-square&logo=docker" alt="Docker Pulls">
</a>
<a href="https://github.com/gin-gonic/gin">
    <img src="https://img.shields.io/github/stars/v03413/bepusdt?style=flat-square&logo=github" alt="GitHub Stars">
</a>
</p>

## 🪧 介绍

缘起于`Epusdt`，但不仅限于此，加入了一些新特性，致力于成为一款更好用的个人`加密货币`收款网关。

> 原始 EPusdt 项目已被转手，现与本项目再没任何联系、两个项目完全独立运营，[具体说明](./docs/faq/epusdt.md)。

## 🎉 新特性

### 🌟 目前支持收款网络

🔥 主流网络：TRON Ethereum BSC Polygon<br>
⚡ 其他网络： X-Layer Solana Aptos Arbitrum-One Base [完整列表](./docs/trade-type.md)

- ✅ 完全兼容 `Epusdt` 插件无缝替换
- ️✅ 支持主流区块网络 不仅限于`USDT`
- ✅ 支持主流法定货币 汇率自动更新
- ✅ 轻依赖 单体二进制文件 部署便捷
- ✅ 底层区块扫描 安全确认 速度稳定
- ✅ 原生兼容`易支付`收单 接入便捷
- ✅ 地址独占模式 底层支持不定额收款
- ✅ 支持MQTT消息发布 交易信息即时广播
- ✅ 快速迭代期 超多实用特性功能等你发现

## 🚀 快捷启动

Docker 快速启动，执行完命令打开地址`http://服务器IP:8080`即可看到初始页面

```bash  
docker run -d --restart=unless-stopped -p 8080:8080 v03413/bepusdt:latest
```
## 📃 技术文档  

- 安装：[Docker](docs/docker/docker.md) [Linux](docs/linux/install.md) [1Panel](./docs/1panel/README.md) [宝塔](./docs/bt_panel/README.md)
- 开发：[API对接](docs/api/api.md) [订单回调](docs/notify/readme.md) [Python](https://github.com/luoyanglang/bepusdt-python-sdk) [PHP](https://github.com/v03413/bepusdt-php-sdk)
- 对接：[独角Next](docs/api/dujiao-next/dujiao-next.md) [彩虹易支付](https://github.com/v03413/Epay-BEpusdt) [whmcs](https://github.com/v03413/whmcs-gateway-epusdt) [其它](docs/api/other.md)
- 其它：[https 配置](./docs/ssl.md) [时钟同步](docs/linux/systemd-timesyncd.md) [收银台修改](docs/checkout/README.md)

## 🖼 功能截图

| 前台收银                                                   | 后台订单                                                 | 订单通知                                                          |
|--------------------------------------------------------|------------------------------------------------------|---------------------------------------------------------------|
| <img src=./docs/images/checkout.png alt=收银台 width=300> | <img src=./docs/images/admin.png alt=后台订单 width=300> | <img src=./docs/images/telegram.png alt=Telegram通知 width=300> |

## ❓ 常见问题

- [服务器配置性能选型推荐 ⚡️](./docs/faq/server.md)
- [服务器流量消耗说明解释](./docs/faq/bandwidth.md)
- [后台入口账密忘记重置教程](./docs/faq/login-reset.md)
- [Telegram 通知 Chat ID 获取教程](docs/faq/telegram-chat-id.md)
- [区块 RPC 节点稳定性说明指南‼️](./docs/faq/rpc-endpoint.md)

## 🏝️ 社区交流

- **Telegram 群组**：[https://t.me/BEpusdtChat](https://t.me/BEpusdtChat)
- **Telegram 频道**：[https://t.me/BEpusdtChannel](https://t.me/BEpusdtChannel)

## 🙏 致谢

感谢原项目提供的灵感和参考：

- [EPusdt](./docs/faq/epusdt.md)

感谢为本项目作出贡献的开发者：

<a href="https://github.com/v03413/BEpusdt/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=v03413/BEpusdt" alt="BEpusdt 贡献者清单">
</a>

## 🌟 Star 历史

[![Stargazers over time](https://starchart.cc/v03413/bepusdt.svg)](https://starchart.cc/v03413/bepusdt)
