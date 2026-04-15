# 📡 MQTT 对接开发文档

> 系统可将扫描到的交易信息广播到 MQTT 服务器，其它系统通过订阅 MQTT 即可实时获取数据。

---

## 概述

BEpusdt 内置了 MQTT 发布能力，每当区块扫描器捕获到一笔**金额合法**的链上转账时，系统会立即将该交易信息以 JSON 格式发布到对应的
MQTT Topic。

外部系统（如监控大屏、自动充值系统、风控系统等）只需连接同一 MQTT Broker 并订阅对应 Topic，即可**实时、被动地**获取链上交易数据，无需轮询
API。

```
区块链网络
    │
    ▼
BEpusdt（扫块）
    │  检测到合法转账
    ▼
MQTT Broker（广播）
    │
    ├──► 订阅方 A（充值系统）
    ├──► 订阅方 B（数据分析）
    └──► 订阅方 C（监控告警）
```

---

## 一、后台配置

在 BEpusdt 后台 → **系统管理 → 基本设置 → MQTT 设置** 中填写以下参数：

| 配置项       | 说明                                         | 示例                  |
|-----------|--------------------------------------------|---------------------|
| MQTT Host | Broker 服务器地址（仅支持 TCP 模式）                   | `mqtt.example.com`  |
| MQTT Port | Broker 端口                                  | `1883`              |
| MQTT 用户名  | Broker 认证用户名，无认证可留空                        | `admin`             |
| MQTT 密码   | Broker 认证密码，无认证可留空                         | `password`          |
| 发布 QoS    | 消息服务质量：`0` 最多一次 / `1` 至少一次 / `2` 恰好一次      | `1`                 |
| 持续监控网络    | 需要**持续扫块**的网络列表，多个用英文逗号分隔；配置后即使无待支付订单也持续推送 | `tron,bsc,ethereum` |

> ⚠️ **注意**：系统当前仅支持 `tcp://` 连接协议，不支持 WebSocket（`ws://`）或 TLS（`ssl://`）。
>
> 客户端 ID 由系统自动生成，格式为 `BEpusdt {进程PID}`，每次启动可能不同，请勿依赖客户端 ID 做业务判断。

---

## 二、Topic 规则

系统按**区块链网络**维度发布消息，Topic 格式如下：

```
bepusdt/transfer/{network}
```

### 网络标识对照表

| 网络       | Topic                       | 支持币种                          |
|----------|-----------------------------|-------------------------------|
| Tron     | `bepusdt/transfer/tron`     | USDT (TRC20)、USDC (TRC20)、TRX |
| Ethereum | `bepusdt/transfer/ethereum` | USDT (ERC20)、USDC (ERC20)、ETH |
| BSC      | `bepusdt/transfer/bsc`      | USDT (BEP20)、USDC (BEP20)、BNB |
| Polygon  | `bepusdt/transfer/polygon`  | USDT、USDC                     |
| Arbitrum | `bepusdt/transfer/arbitrum` | USDT、USDC                     |
| Base     | `bepusdt/transfer/base`     | USDC                          |
| Solana   | `bepusdt/transfer/solana`   | USDT、USDC                     |
| Aptos    | `bepusdt/transfer/aptos`    | USDT、USDC                     |
| X Layer  | `bepusdt/transfer/xlayer`   | USDT、USDC                     |
| Plasma   | `bepusdt/transfer/plasma`   | USDT                          |

### 订阅示例

```
# 仅订阅 Tron 网络
bepusdt/transfer/tron

# 同时订阅 BSC 和 Polygon（分两次 Subscribe）
bepusdt/transfer/bsc
bepusdt/transfer/polygon

# 使用通配符订阅所有网络（MQTT v3.1.1）
bepusdt/transfer/#
```

---

## 三、消息格式

消息体为 **UTF-8 编码的 JSON 字符串**，结构如下：

### 字段说明

| 字段名            | 类型     | 说明                                           |
|----------------|--------|----------------------------------------------|
| `network`      | string | 区块链网络标识，如 `tron`、`bsc`、`ethereum`            |
| `tx_hash`      | string | 链上交易哈希（全局唯一，可作为幂等键）                          |
| `amount`       | string | 实际到账的**加密货币数量**（非法币金额）                       |
| `from_address` | string | 发款方钱包地址                                      |
| `recv_address` | string | 收款方钱包地址                                      |
| `timestamp`    | number | 交易发生时间，Unix 时间戳（秒）                           |
| `trade_type`   | string | 交易类型，完整列表见 [trade-type.md](../trade-type.md) |
| `block_num`    | number | 交易所在区块高度                                     |

### 消息示例

```json
{
  "network": "tron",
  "tx_hash": "12ef6267b42e43959795cf31808d0cc72b3d0a48953ed19c61d4b6665a341d10",
  "amount": "4.25",
  "from_address": "TJRabPrwbZy45sbavfcjinPJC18iYKbPa5",
  "recv_address": "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t",
  "timestamp": 1744732800,
  "trade_type": "usdt.trc20",
  "block_num": 68123456
}
```

---

## 四、注意事项

### 消息去重

同一笔交易**可能被发布多次**（例如系统重启、扫块重叠时），建议以 `tx_hash` 作为幂等键，在业务层去重处理。

### 消息触发时机

MQTT 消息在**扫块阶段**触发，此时交易已出现在区块链上，但系统层面尚未完成订单匹配或回调流程。如果需要确认订单是否支付成功，应以
**订单回调通知**（`notify_url`）为准；MQTT 消息可作为辅助的实时感知手段。

### 金额过滤

系统仅广播金额在合法范围内的交易（由后台「监控最小金额」配置决定）。超出范围的链上转账不会触发 MQTT 推送。

### 网络断线重连

BEpusdt 的 MQTT 客户端已启用自动重连，断线后会自动恢复并重新订阅所有 Topic。订阅方同样建议使用支持自动重连的 MQTT 客户端库。

### QoS 建议

| 场景                | 推荐 QoS |
|-------------------|--------|
| 数据监控、日志记录（允许少量丢失） | `0`    |
| 充值到账、财务类业务（不能漏消息） | `1`    |
| 对消息重复极度敏感的场景      | `2`    |

---

## 五、常见问题

**Q：MQTT 消息和订单回调（notify_url）有什么区别？**

> MQTT 是对**所有捕获到的合法链上转账**广播，不区分是否匹配到订单，实时性更高；订单回调仅针对成功匹配的支付订单，包含完整订单信息。两者可结合使用。

**Q：如何只监听特定收款地址的转账？**

> MQTT 消息包含 `recv_address` 字段，订阅方自行在消费逻辑中按地址过滤即可，Topic 层面暂不支持按地址过滤。

**Q：没有活跃订单时为什么收不到消息？**

> 默认情况下，BEpusdt 仅在有待支付订单时才对对应网络扫块。需在后台**「持续监控网络」**中填入目标网络标识（如 `tron,bsc`
> ），才能保证无订单时也持续推送。

**Q：支持 WebSocket 或 TLS 连接吗？**

> 当前版本仅支持 TCP 模式（`tcp://`），暂不支持 `ws://` / `wss://` / `ssl://`。

**Q：`amount` 字段的单位是什么？**

> `amount` 为**加密货币数量**（如 USDT 数量），并非法币金额。例如 `"amount": "4.25"` 表示实际到账 4.25 USDT。

---

## 参考资料

- [交易类型完整列表](../trade-type.md)
- [API 对接文档](./api.md)
- [订单回调通知](../notify/readme.md)
- [BEpusdt GitHub](https://github.com/v03413/BEpusdt)
