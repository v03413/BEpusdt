# TON Center V3

> BEpusdt 的 TON 网络扫描使用 TON Center V3 API。
> 系统会通过 V3 接口查询原生 TON 转账和 Jetton 转账，并将解析后的交易继续交给统一的订单匹配、订单回调和 MQTT 发布流程。

## 默认端点

系统默认使用以下 TON Center V3 端点：

```text
https://toncenter.com/api/v3
```

一般情况下不建议修改该端点。只有在您使用自建 TON Center V3 服务，或使用兼容 TON Center V3 接口的第三方服务时，才需要修改。

## 为什么建议配置 Api Key

TON Center 的公开接口存在频率限制。未配置 Api Key 时，扫描地址、确认订单或持续 MQTT 监听都更容易受到限流影响。

建议配置 TON Center V3 Api Key，原因如下：

- 提高 TON 扫描稳定性。
- 降低因接口限流导致的订单确认延迟。
- 持续监听 TON 地址并发布 MQTT 消息时更可靠。
- 便于后续升级付费额度，而不需要修改系统代码。

## 获取 Ton Center Api Key

1. 在 Telegram 中点击访问 [@tonceter](https://t.me/toncenter) 机器人，点击 Start 开始。
2. Toncenter 机器人会回复您一条欢迎消息，点击消息底部的“管理 API 密钥 (Manage API Keys)”按钮，将会打开 toncenter bot 小程序。
3. 点击小程序底部的“创建 API 密钥 (Create API Key)”按钮，创建新的 Api Key。
4. 在 API 密钥详情 (API Key Detail) 页面，名称 (Name) / 描述 (Description) 随便填，网络 (Network) 选择“主网 (Mainnet)”，然后点击“创建 (Create)”按钮完成创建。
5. 复制生成的 Api Key，填入 BEpusdt 区块节点 TON Center V3 Api Key 输入框中。

> 说明：如果你有多个 TON Center Api Key，可以全部填入 BEpusdt，用半角逗号分隔。

## 配置 Api Key

登录 BEpusdt 后台，进入：

```text
系统管理 -> 区块节点 -> TON 网络
```

然后确认或填写：

```text
TON Center V3 Endpoint: https://toncenter.com/api/v3 （保持不变）
TON Center V3 Api Key: 你的 TON Center Api Key
```

保存后会立即生效。系统请求 TON Center V3 时会通过 HTTP Header 发送：

```http
X-API-Key: 你的 TON Center Api Key
```

如果需要配置多个 Api Key，请使用半角逗号分隔，例如：

```text
key-1,key-2,key-3
```

系统会在请求 TON Center V3 时从这些 Key 中选择一个使用，以降低单个 Key 触发频率限制的概率。

## MQTT 持续监听 TON

如果需要在没有待支付订单时也持续监听 TON 地址并发布 MQTT 消息，请进入：

```text
系统管理 -> 基本设置 -> MQTT 设置
```

在「区块链网络」中勾选 `Ton`，并确保 MQTT Host、Port、Topic 前缀等配置正确。

TON 网络的 MQTT Topic 格式为：

```text
bepusdt/transfer/ton
```

## 注意事项

- `TON Center V3 Endpoint` 不是传统 JSON-RPC 端点，不要填写 `/api/v2/jsonRPC`。
- 推荐使用 TonCenter 官方 `https://toncenter.com/api/v3`，除非你明确知道自定义端点兼容 TON Center V3。
- Api Key 不是钱包私钥，不涉及链上签名，但仍建议妥善保管。
- 如果出现接口限流或扫描延迟，优先检查 Api Key、端点可用性和 TON Center 套餐额度。
