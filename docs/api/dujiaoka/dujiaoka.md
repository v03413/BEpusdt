# 独角数卡对接教程

## 重要提示

由于新版 `BEpusdt` 支持多区块链 USDT、超时回调等新特性，独角数卡自带的 `Epusdt` 插件已无法兼容，**必须替换原有插件**。

> ⚠️ **重要警告**
>
> 独角数卡自带的 `Epusdt` 插件无法直接使用，强行使用必定出现问题！

## 插件替换说明

我已对独角数卡进行适配调整，并向原项目提交了 [PR #424](https://github.com/assimon/dujiaoka/pull/424)，期待官方尽快合并。

### 安装步骤

1. **下载适配文件**

   下载更新后的控制器文件：
   ```
   https://raw.githubusercontent.com/v03413/dujiaoka/refs/heads/master/app/Http/Controllers/Pay/EpusdtController.php
   ```

2. **替换原有文件**

   将下载的文件替换到独角数卡网站目录：
   ```
   app/Http/Controllers/Pay/EpusdtController.php
   ```

## 配置步骤

确保 `BEpusdt` 已成功安装并运行后，按照以下步骤配置：

### 1. 进入支付配置

独角数卡后台 → 支付配置 → 新增

![配置示例](1.png)

### 2. 参数填写说明

| 参数名称       | 填写内容                                                      | 说明                                                                        |
|------------|-----------------------------------------------------------|---------------------------------------------------------------------------|
| **商户ID**   | 与 `BEpusdt` 中 `AUTH_TOKEN` 参数保持一致                         | 用于身份验证                                                                    |
| **支付标识**   | 收款交易类型                                                    | 可选值：`tron.trx`、`usdt.trc20`、`usdt.erc20` 等<br/>[完整列表参考](../../trade-type.md) |
| **商户密钥**   | `https://your-domain.com/api/v1/order/create-transaction` | ⚠️ 请将 `your-domain.com` 替换为实际域名                                           |
| **支付处理路由** | `pay/epusdt`                                              | 固定值，无需修改                                                                  |

### 3. 完成配置

根据实际情况填写其他参数后，点击保存即可完成对接。

## 常见问题

- 如果遇到支付失败，请检查 `BEpusdt` 服务是否正常运行
- 确认商户密钥中的域名是否正确配置
- 验证 `AUTH_TOKEN` 是否与独角数卡中的商户ID一致

---