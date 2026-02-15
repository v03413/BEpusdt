# Dujiao-Next 对接教程

本文档介绍如何在 [Dujiao-Next](https://github.com/dujiao-next/dujiao-next)（独角next）中对接 BEpusdt。

## 前置条件

- BEpusdt 已成功安装并运行
- Dujiao-Next 已部署

## 配置步骤

### 1. 进入支付渠道配置

Dujiao-Next 后台 → 支付渠道 → 新增渠道

### 2. 基本参数填写

| 参数名称 | 填写内容 | 说明 |
|---------|---------|------|
| 渠道名称 | USDT-TRC20 支付 | 自定义名称，建议包含币种 |
| 渠道类型 | `epusdt` | 固定值，选择 BEpusdt |
| 支付方式 | `usdt-trc20` / `usdc-trc20` / `trx` | 选择支持的币种 |
| 交互方式 | `qr` 或 `redirect` | 二维码或跳转收银台 |

### 3. 支持的渠道类型

Dujiao-Next 支持以下 BEpusdt 币种：

| 渠道类型 | 说明 | 对应 trade_type |
|---------|------|----------------|
| `usdt-trc20` | USDT (TRC20) | `usdt.trc20` |
| `usdc-trc20` | USDC (TRC20) | `usdc.trc20` |
| `trx` | TRX | `tron.trx` |

> **注意**：选择渠道类型后，系统会自动设置对应的 `trade_type`，无需手动填写。

### 4. 渠道配置参数

在「渠道配置」中填写以下 JSON：

```json
{
  "gateway_url": "https://your-bepusdt-domain.com",
  "auth_token": "your-api-token",
  "fiat": "CNY",
  "notify_url": "https://your-api-domain.com/api/v1/payments/callback",
  "return_url": "https://your-user-domain.com/pay"
}
```

### 5. 配置参数说明

| 参数 | 必填 | 说明 | 示例 |
|------|------|------|------|
| `gateway_url` | ✅ | BEpusdt 网关地址 | `https://pay.example.com` |
| `auth_token` | ✅ | BEpusdt API Token | 从 BEpusdt 后台获取 |
| `trade_type` | ❌ | 交易类型（自动设置） | `usdt.trc20` / `usdc.trc20` / `tron.trx` |
| `fiat` | ❌ | 法币类型，默认 `CNY` | `CNY` / `USD` |
| `notify_url` | ✅ | 异步回调地址 | `https://api.example.com/api/v1/payments/callback` |
| `return_url` | ✅ | 支付成功跳转地址 | `https://shop.example.com/pay` |

> **重要提示**：
> - `trade_type` 会根据「支付方式」自动设置，无需在配置中填写
> - `notify_url` 路径必须是 `/api/v1/payments/callback`（不是 `/api/public/payment/callback`）
> - `return_url` 路径必须是 `/pay`（不是 `/order` 或 `/payment`）

## 获取 API Token

1. 登录 BEpusdt 后台
2. 系统管理 → 基本设置 → API 设置
3. 复制「对接令牌」

## 配置示例

### 示例 1：USDT-TRC20 支付

```
渠道名称: USDT-TRC20 支付
渠道类型: epusdt
支付方式: usdt-trc20
交互方式: qr

渠道配置:
{
  "gateway_url": "http://127.0.0.1:8088",
  "auth_token": "your-token-here",
  "fiat": "CNY",
  "notify_url": "https://api.example.com/api/v1/payments/callback",
  "return_url": "https://shop.example.com/pay"
}
```

### 示例 2：TRX 支付

```
渠道名称: TRX 支付
渠道类型: epusdt
支付方式: trx
交互方式: redirect

渠道配置:
{
  "gateway_url": "http://127.0.0.1:8088",
  "auth_token": "your-token-here",
  "fiat": "CNY",
  "notify_url": "https://api.example.com/api/v1/payments/callback",
  "return_url": "https://shop.example.com/pay"
}
```

## 回调机制说明

### 异步通知（notify_url）

BEpusdt 会在支付状态变化时，向 `notify_url` 发送 POST 请求：

```
POST /api/v1/payments/callback
Content-Type: application/x-www-form-urlencoded

order_id=xxx&trade_id=xxx&amount=xxx&actual_amount=xxx&token=xxx&...&signature=xxx
```

Dujiao-Next 会验证签名并更新订单状态。

### 同步跳转（return_url）

用户支付完成后，BEpusdt 会跳转到 `return_url`，并附带参数：

```
https://shop.example.com/pay?order_id=xxx&epusdt_return=1
```

Dujiao-Next 前端会检测 `epusdt_return=1` 参数，自动查询订单状态并显示支付结果。

## 常见问题

### 1. 支付创建失败

**可能原因**：
- `gateway_url` 配置错误或 BEpusdt 服务未运行
- `auth_token` 与 BEpusdt 配置不一致
- 网络连接问题（检查防火墙）

**解决方法**：
- 检查 BEpusdt 服务状态：`curl http://your-bepusdt/health`
- 验证 API Token 是否正确
- 查看 Dujiao-Next 后端日志

### 2. 回调未收到或订单状态未更新

**可能原因**：
- `notify_url` 配置错误（路径必须是 `/api/v1/payments/callback`）
- BEpusdt 无法访问回调地址（网络/防火墙问题）
- 签名验证失败

**解决方法**：
- 确认回调地址可被外网访问：`curl https://your-api/api/v1/payments/callback`
- 检查 BEpusdt 日志中的回调错误
- 检查 Dujiao-Next 后端日志中的签名验证信息

### 3. 支付完成后跳转错误

**可能原因**：
- `return_url` 路径配置错误（必须是 `/pay` 不是 `/payment` 或 `/order`）
- 前端路由配置问题

**解决方法**：
- 确认 `return_url` 配置为：`https://your-user-domain.com/pay`
- 清除浏览器缓存后重试

### 4. 签名验证失败

**可能原因**：
- `auth_token` 不一致（注意前后空格）
- 回调参数被中间件修改

**解决方法**：
- 确保 `auth_token` 完全一致
- 检查是否有反向代理修改了请求参数

### 5. 金额不匹配

**可能原因**：
- 汇率波动导致实际支付金额与订单金额不同
- BEpusdt 配置的汇率源问题

**解决方法**：
- BEpusdt 会返回 `actual_amount`（实际支付金额），Dujiao-Next 会记录此金额
- 可在订单详情中查看实际支付金额

## 技术支持

- Dujiao-Next: [GitHub Issues](https://github.com/dujiao-next/dujiao-next/issues)
- BEpusdt: [Telegram 群组](https://t.me/BEpusdtChat)

---

**贡献者**: 狼哥 ([@luoyanglang](https://github.com/luoyanglang) | Telegram: [@luoyanglang](https://t.me/luoyanglang))

> 本文档适用于 Dujiao-Next 最新版本和 BEpusdt 最新版本  
> 最后更新：2026-02-13
