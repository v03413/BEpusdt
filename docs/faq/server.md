# 服务器选项

## 配置要求

**内存**：单网络扫块内存占用通常不超过 100MB，开启多网络扫块时内存会线性增加。推荐：1GB 内存。

**磁盘**：空间需求较小，但强烈建议使用 SSD 或 NVMe 高速磁盘。使用 HDD 机械硬盘可能会触发 Slow SQL 告警。

**网络**：

- 必须保证网络畅通稳定，避免使用网络受限地区的服务器
- 扫块会持续下载区块数据，日流量消耗可达数 GB 以上

**时间同步**：务必配置 NTP 服务自动同步时间，系统多处逻辑依赖时间准确性。

**地区**：优先选择新加坡或美国机房。

## 推荐配置

最低配置：1 核 CPU / 1GB 内存 / 10GB 磁盘

推荐服务商：AWS、Google Cloud、Azure、DigitalOcean 等大厂商，可提供更好的稳定性和网络质量。

## 开发者推荐

[![DigitalOcean Referral Badge](https://web-platforms.sfo2.cdn.digitaloceanspaces.com/WWW/Badge%203.svg)](https://www.digitalocean.com/?refcode=7204f1398099&utm_campaign=Referral_Invite&utm_medium=Referral_Program&utm_source=badge)

本项目基于 [DigitalOcean](https://m.do.co/c/7204f1398099) 环境开发测试，兼容性最佳；月费用低于 $10，适合个人和小型项目。
