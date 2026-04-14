## 📚 API 对接文档

> 💡 **推荐阅读**：在开始对接前，建议先阅读 [《BEpusdt 执行原理说明》](./how-it-works.md)
> ，了解订单匹配机制、金额匹配模式与汇率浮动等核心原理，有助于更准确地理解本文档中的参数设计。

### 基础信息

**认证方式**：签名认证（详见[签名算法](#签名算法)）

**请求格式**：JSON

**响应格式**：JSON

---

## 接口列表

> **两种订单模式说明**
>
> BEpusdt 的订单创建接口支持两种工作模式，通过 `amount` 参数决定：
>
> | 模式 | 触发条件 | 地址分配 | 匹配规则 | 典型场景 |
> |------|---------|---------|---------|---------|
> | **普通模式** | `amount > 0` | 地址 + 金额同时锁定 | 到账金额须匹配订单金额（±匹配容差） | 商城购物、固定金额收款 |
> | **独占模式** | `amount = 0` 或不传 | 仅锁定地址，金额不限 | 任意金额到账均触发回调 | 充值场景、钱包监控、不定额收款 |
>
> 独占模式下，回调中的 `amount`（法币）和 `actual_amount`（加密货币）均以**实际到账金额**为准，创建订单时无需预设金额。

### 1. 创建交易

创建支付订单并获取收银台付款地址。

> **注意**：使用相同 `order_id` 创建订单时，系统会根据最新参数重建订单（订单金额、交易类型、收款地址、法币类型），同时超时时间重置。这意味着商户可以基于同一订单号实现独立收银台，灵活变更支付参数。

#### 请求地址

```http
POST /api/v1/order/create-transaction
```

#### 请求参数

| 参数名          | 类型     | 必填 | 说明                                                                                                                                                          |
|--------------|--------|----|-------------------------------------------------------------------------------------------------------------------------------------------------------------|
| order_id     | string | ✅  | 商户订单编号（唯一标识）                                                                                                                                                |
| notify_url   | string | ✅  | 支付结果异步回调地址                                                                                                                                                  |
| redirect_url | string | ✅  | 支付成功后商户跳转地址                                                                                                                                                 |
| signature    | string | ✅  | 签名字符串（详见[签名算法](#签名算法)）                                                                                                                                      |
| amount       | number | ❌  | 支付金额（法币金额）；留空或传 `0` 则进入**地址独占模式**，收到任意金额均触发回调（详见[执行原理 §3](../faq/how-it-works.md#3-地址独占模式)）                                                                 |
| trade_type   | string | ❌  | 交易类型，默认 `usdt.trc20`<br/>完整列表：[trade-type.md](../trade-type.md)                                                                                             |
| fiat         | string | ❌  | 法币类型，默认 `CNY`<br/>可选：`CNY`、`USD`、`EUR`、`GBP`、`JPY`                                                                                                          |
| address      | string | ❌  | 指定收款地址（留空则自动分配）                                                                                                                                             |
| name         | string | ❌  | 商品名称                                                                                                                                                        |
| timeout      | number | ❌  | 订单超时时间（秒），最低 120 秒<br/>留空则使用配置 `payment_timeout`，默认 600 秒                                                                                                   |
| rate         | string | ❌  | 强制指定汇率，支持多种写法：<br/>• `7.4` - 固定汇率 7.4<br/>• `~1.02` - 最新汇率上浮 2%<br/>• `~0.97` - 最新汇率下浮 3%<br/>• `+0.3` - 最新汇率加 0.3<br/>• `-0.2` - 最新汇率减 0.2<br/>留空则使用系统配置汇率 |

#### 请求示例

```json
{
  "order_id": "20250120001",
  "amount": 28.88,
  "fiat": "CNY",
  "trade_type": "usdt.trc20",
  "name": "测试商品",
  "address": "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t",
  "notify_url": "https://example.com/notify",
  "redirect_url": "https://example.com/success",
  "timeout": 1200,
  "signature": "1cd4b52df5587cfb1968b0c0c6e156cd"
}
```

#### 响应参数

| 参数名                  | 类型     | 说明           |
|----------------------|--------|--------------|
| status_code          | number | 状态码，200 表示成功 |
| message              | string | 响应消息         |
| data                 | object | 订单数据         |
| data.fiat            | string | 交易法币类型       |
| data.trade_id        | string | 系统交易 ID      |
| data.order_id        | string | 商户订单编号       |
| data.amount          | string | 请求支付金额（法币）   |
| data.actual_amount   | string | 实际支付金额（加密货币） |
| data.status          | string | 订单状态，1 表示待付款 |
| data.token           | string | 收款地址         |
| data.expiration_time | number | 订单有效期（秒）     |
| data.payment_url     | string | 收银台付款链接地址    |

> **计算公式**：`actual_amount` = `amount` ÷ 汇率

#### 响应示例

```json
{
  "status_code": 200,
  "message": "success",
  "data": {
    "fiat": "CNY",
    "trade_id": "b3d2477c-d945-41da-96b7-f925bbd1b415",
    "order_id": "20250120001",
    "amount": "28.88",
    "actual_amount": "4.25",
    "token": "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t",
    "expiration_time": 1200,
    "status": 1,
    "payment_url": "https://example.com/pay/checkout-counter/b3d2477c-d945-41da-96b7-f925bbd1b415"
  },
  "request_id": ""
}
```

---

<details>
<summary><strong>2. 取消交易</strong>　取消指定订单，系统将停止监控该订单并释放金额占用。</summary>

#### 请求地址

```http
POST /api/v1/order/cancel-transaction
```

#### 请求参数

| 参数名       | 类型     | 必填 | 说明      |
|-----------|--------|----|---------|
| trade_id  | string | ✅  | 系统交易 ID |
| signature | string | ✅  | 签名字符串   |

#### 请求示例

```json
{
  "trade_id": "b3d2477c-d945-41da-96b7-f925bbd1b415",
  "signature": "1cd4b52df5587cfb1968b0c0c6e156cd"
}
```

#### 响应示例

```json
{
  "status_code": 200,
  "message": "success",
  "data": {
    "trade_id": "b3d2477c-d945-41da-96b7-f925bbd1b415"
  },
  "request_id": ""
}
```

</details>

---

<details>
<summary><strong>3. 创建订单</strong>　创建收银台订单，让用户在付款页自由选择支付方式。</summary>

> **注意**：使用相同 `order_id` 创建订单时，系统会根据最新参数重建订单（订单金额、交易类型、收款地址、法币类型），同时超时时间重置。

#### 请求地址

```http
POST /api/v1/order/create-order
```

#### 请求参数

| 参数名          | 类型     | 必填 | 说明                                                                                                                                       |
|--------------|--------|----|------------------------------------------------------------------------------------------------------------------------------------------|
| order_id     | string | ✅  | 商户订单编号（唯一标识）                                                                                                                             |
| notify_url   | string | ✅  | 支付结果异步回调地址                                                                                                                               |
| redirect_url | string | ✅  | 支付成功后商户跳转地址                                                                                                                              |
| signature    | string | ✅  | 签名字符串（详见[签名算法](#签名算法)）                                                                                                                   |
| amount       | number | ❌  | 支付金额（法币金额）；留空或传 `0` 则进入**地址独占模式**，收到任意金额均触发回调（详见[执行原理 §3](../faq/how-it-works.md#3-地址独占模式)）                                              |
| currencies   | string | ❌  | 限定交易币种，留空则不限制付款币种。<br/>多个币种请使用半角逗号分隔，黑名单模式以短横线开头。<br/>例如：<br/>`USDT`（仅限 USDT）<br/>`USDT,USDC` （限 USDT/USDC）<br/>`-ETH,-BNB` （排除 ETH/BNB） |
| fiat         | string | ❌  | 法币类型，默认 `CNY`<br/>可选：`CNY`、`USD`、`EUR`、`GBP`、`JPY`                                                                                       |
| name         | string | ❌  | 商品名称                                                                                                                                     |
| timeout      | number | ❌  | 订单超时时间（秒），最低 180 秒<br/>留空则使用配置 `payment_timeout`，默认 600 秒                                                                                |

#### 请求示例

```json
{
  "order_id": "20250120001",
  "amount": 28.88,
  "fiat": "CNY",
  "currencies": "USDT,USDC",
  "name": "测试商品",
  "notify_url": "https://example.com/notify",
  "redirect_url": "https://example.com/success",
  "timeout": 1200,
  "signature": "1cd4b52df5587cfb1968b0c0c6e156cd"
}
```

#### 响应参数

| 参数名                  | 类型     | 说明           |
|----------------------|--------|--------------|
| status_code          | number | 状态码，200 表示成功 |
| message              | string | 响应消息         |
| data.fiat            | string | 交易法币类型       |
| data.trade_id        | string | 系统交易 ID      |
| data.order_id        | string | 商户订单编号       |
| data.amount          | string | 请求支付金额（法币）   |
| data.status          | string | 订单状态，1 表示待付款 |
| data.expiration_time | number | 订单有效期（秒）     |
| data.payment_url     | string | 收银台订单链接      |

#### 响应示例

```json
{
  "status_code": 200,
  "message": "success",
  "data": {
    "fiat": "CNY",
    "trade_id": "b3d2477c-d945-41da-96b7-f925bbd1b415",
    "order_id": "20250120001",
    "amount": "28.88",
    "expiration_time": 1200,
    "payment_url": "https://example.com/pay/cashier/b3d2477c-d945-41da-96b7-f925bbd1b415"
  },
  "request_id": ""
}
```

</details>

---

<details>
<summary><strong>4. 更新订单付款方式</strong>　为收银台订单选定支付链，获取具体收款地址与金额。</summary>

> **注意**：使用相同 `trade_id` 更新订单时，系统会根据最新交易类型更新订单付款方式。

#### 请求地址

```http
POST /api/v1/pay/update-order
```

#### 请求参数

| 参数名      | 类型     | 必填 | 说明                                  |
|----------|--------|----|-------------------------------------|
| trade_id | string | ✅  | 系统交易 ID                             |
| currency | string | ✅  | 加密货币币种，可选：`USDT`、`USDC`、`TRX` 等     |
| network  | string | ✅  | 加密货币所属网络，如：`tron`、`polygon`、`bsc` 等 |

#### 请求示例

```json
{
  "trade_id": "b3d2477c-d945-41da-96b7-f925bbd1b415",
  "currency": "USDT",
  "network": "tron"
}
```

#### 响应参数

| 参数名                  | 类型     | 说明           |
|----------------------|--------|--------------|
| status_code          | number | 状态码，200 表示成功 |
| data.fiat            | string | 交易法币类型       |
| data.trade_id        | string | 系统交易 ID      |
| data.order_id        | string | 商户订单编号       |
| data.amount          | string | 请求支付金额（法币）   |
| data.actual_amount   | string | 实际支付金额（加密货币） |
| data.status          | string | 订单状态，1 表示待付款 |
| data.expiration_time | number | 订单有效期（秒）     |
| data.payment_url     | string | 收银台付款链接      |

#### 响应示例

```json
{
  "status_code": 200,
  "message": "success",
  "data": {
    "fiat": "CNY",
    "trade_id": "b3d2477c-d945-41da-96b7-f925bbd1b415",
    "order_id": "20250120001",
    "amount": "28.88",
    "actual_amount": "4.25",
    "expiration_time": 1200,
    "status": 1,
    "payment_url": "https://example.com/pay/checkout-counter/b3d2477c-d945-41da-96b7-f925bbd1b415"
  },
  "request_id": ""
}
```

</details>

---

<details>
<summary><strong>5. 订单可用付款方式列表</strong>　返回指定订单当前所有可用的支付方式及实时报价。</summary>

#### 请求地址

```http
POST /api/v1/pay/methods
```

#### 请求参数

| 参数名      | 类型     | 必填 | 说明                                    |
|----------|--------|----|---------------------------------------|
| trade_id | string | ✅  | 系统交易 ID                               |
| currency | string | ❌  | 货币名称（如 USDT / USDC / TRX），不传则返回全部可用方式 |

#### 请求示例

```json
{
  "trade_id": "b3d2477c-d945-41da-96b7-f925bbd1b415"
}
```

#### methods 字段说明

| 字段                | 类型      | 说明           |
|-------------------|---------|--------------|
| amount            | string  | 应付金额（法币）     |
| actual_amount     | string  | 实际支付金额（加密货币） |
| fiat              | string  | 法币单位         |
| exchange_rate     | string  | 汇率           |
| currency          | string  | 货币名称         |
| network           | string  | 网络名称         |
| token_net_name    | string  | 代币协议标准名      |
| token_custom_name | string  | 用户自定义显示名称    |
| is_popular        | boolean | 是否为流行网络      |

#### 响应示例

```json
{
  "status_code": 200,
  "message": "success",
  "data": {
    "methods": [
      {
        "amount": "28.88",
        "actual_amount": "4.25",
        "fiat": "CNY",
        "exchange_rate": "6.79",
        "currency": "USDT",
        "network": "bsc",
        "token_net_name": "BEP20",
        "token_custom_name": "",
        "is_popular": false
      }
    ]
  },
  "request_id": ""
}
```

> TODO: `is_popular` 与 `token_custom_name` 当前为占位值，待后期完善。

</details>

---

<details>
<summary><strong>6. 支付回调通知</strong>　订单状态变更时系统主动推送至 <code>notify_url</code>。</summary>

#### 通知参数

| 参数名                  | 类型     | 说明                                                       |
|----------------------|--------|----------------------------------------------------------|
| trade_id             | string | 系统交易 ID                                                  |
| order_id             | string | 商户订单编号                                                   |
| amount               | number | 请求支付金额（法币）                                               |
| actual_amount        | number | 实际支付金额（加密货币）                                             |
| token                | string | 收款地址                                                     |
| block_transaction_id | string | 区块链交易哈希                                                  |
| signature            | string | 签名字符串                                                    |
| status               | number | 订单状态：<br/>• `1` - 等待支付<br/>• `2` - 支付成功<br/>• `3` - 支付超时 |

#### 通知示例

```json
{
  "trade_id": "b3d2477c-d945-41da-96b7-f925bbd1b415",
  "order_id": "20250120001",
  "amount": 28.88,
  "actual_amount": 4.25,
  "token": "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t",
  "block_transaction_id": "12ef6267b42e43959795cf31808d0cc72b3d0a48953ed19c61d4b6665a341d10",
  "signature": "1cd4b52df5587cfb1968b0c0c6e156cd",
  "status": 2
}
```

#### 响应要求

- **成功**：返回字符串 `success`（不区分大小写）
- **失败**：返回其他内容，系统将按退避策略重试通知

</details>

---

<details>
<summary><strong>签名算法</strong></summary>

### 签名流程

**第一步：参数排序与拼接**

1. 筛选所有 **非空** 且 **非 signature** 的参数
2. 按参数名 **ASCII 码** 从小到大排序（字典序）
3. 按 `key=value` 格式拼接，使用 `&` 连接

**第二步：加密生成签名**

1. 在拼接字符串末尾追加 API Token（无 `&` 符号）
2. 对完整字符串进行 MD5 加密
3. 将结果转为 **小写** 即为 `signature`

### 签名示例

假设请求参数如下：

```json
{
  "order_id": "20220201030210321",
  "amount": 42,
  "notify_url": "http://example.com/notify",
  "redirect_url": "http://example.com/redirect"
}
```

假设 API Token 为：`epusdt_password_xasddawqe`

**步骤 1：排序拼接**

```text
amount=42&notify_url=http://example.com/notify&order_id=20220201030210321&redirect_url=http://example.com/redirect
```

**步骤 2：追加 Token 并加密**

```text
MD5(amount=42&notify_url=http://example.com/notify&order_id=20220201030210321&redirect_url=http://example.com/redirectepusdt_password_xasddawqe)
```

**最终签名**：`1cd4b52df5587cfb1968b0c0c6e156cd`

### 代码参考

- **PHP 实现
  **：[点击查看](https://github.com/v03413/Epay-BEpusdt/blob/b7fa8fd608d71ce50e0f8eabb1717783c96761ac/bepusdt_plugin.php#L108:L127)
- **Python 实现**：[BEpusdt Python SDK](https://github.com/luoyanglang/bepusdt-python-sdk)

### 签名规则总结

| 规则              | 说明                            |
|-----------------|-------------------------------|
| 参数名区分大小写        | `Amount` 和 `amount` 是不同参数     |
| 空值不参与签名         | `null`、`""`、`undefined` 等均不参与 |
| signature 不参与签名 | 签名本身不参与签名计算                   |
| 必须按字典序排序        | 确保双方计算结果一致                    |
| MD5 结果必须小写      | 统一使用小写字母                      |

</details>

---

<details>
<summary><strong>常见问题</strong></summary>

### 1. 如何获取 API Token？

登录后台 -> 系统管理 -> 基本设置 -> API 设置 -> 对接令牌

### 2. 订单重建机制是什么？

当使用相同 `order_id` 创建订单时：

- ✅ 更新订单金额、交易类型、收款地址、法币类型
- ✅ **会** 重置超时时间

### 3. 支持哪些交易类型？

完整列表请查看：[`trade-type.md`](../trade-type.md)

常用类型：

- `usdt.trc20` - USDT (TRC20)
- `usdt.erc20` - USDT (ERC20)
- `tron.trx` - TRX
- `usdc.polygon` - USDC (Polygon)

### 4. 回调通知如何验证签名？

使用与请求相同的签名算法验证 `signature` 参数，确保通知来自可信源。

### 5. 如何测试对接？

建议步骤：

1. 搭建测试环境
2. 使用小额测试订单
3. 验证签名算法正确性
4. 测试回调通知处理逻辑
5. 确认超时和取消场景

</details>

---

## 参考资料

- [交易类型完整列表](../trade-type.md)
- [订单回调通知](../notify/readme.md)
- [易支付对接插件](https://github.com/v03413/Epay-BEpusdt)
- [Python SDK](https://github.com/luoyanglang/bepusdt-python-sdk)
- [原项目参考](https://github.com/assimon/epusdt/blob/master/wiki/API.md)
