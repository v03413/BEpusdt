package model

import "github.com/shopspring/decimal"

type ConfKey string
type Fiat string
type CoinId string
type Crypto string
type TradeType string
type MatchMode string
type Network string
type Range struct {
	MinAmount decimal.Decimal
	MaxAmount decimal.Decimal
}
type TradeTypeConf struct {
	Alias        string  // 类型别名，主要用户前端展示
	NetworkName  string  // 网络名称，用于前端展示
	Network      Network // 所属区块链网络
	Crypto       Crypto  // 币种类型
	Native       bool    // 是否原生币
	Contract     string  // 合约地址，原生币为空
	Decimal      int32   // 小数位
	AmountRange  Range   // 合法数额范围；这里特指则扫块时[数额范围]，目前偷懒全部写死一个大概合理的范围，后面有问题再说...
	ExplorerFmt  string  // 区块浏览器交易链接格式化字符串，%s 位置替换为交易哈希
	EndpointKey  ConfKey // RPC 端点配置键
	AddrCaseSens bool    // 钱包地址是否大小写敏感，如果为 false 则会统一转为小写比较
}

const (
	AdminUsername ConfKey = "admin_username"
	AdminPassword ConfKey = "admin_password"
	AdminSecure   ConfKey = "admin_secure"
	AdminSecret   ConfKey = "admin_secret"
	AdminLoginIP  ConfKey = "admin_login_ip"
	AdminLoginAt  ConfKey = "admin_login_at"

	ApiAuthToken ConfKey = "api_auth_token" // API 对接令牌
	ApiAppUri    ConfKey = "api_app_uri"    // API 对接地址（收银台地址）

	AtomUSDT ConfKey = "atom_usdt"
	AtomUSDC ConfKey = "atom_usdc"
	AtomTRX  ConfKey = "atom_trx"
	AtomBNB  ConfKey = "atom_bnb"
	AtomETH  ConfKey = "atom_eth"

	MonitorMinAmount  ConfKey = "monitor_min_amount" // 监控最小金额，低于此金额的入账不进行通知
	PaymentMinAmount  ConfKey = "payment_min_amount"
	PaymentMaxAmount  ConfKey = "payment_max_amount"
	PaymentTimeout    ConfKey = "payment_timeout"     // 订单支付超时时间，单位秒
	PaymentStaticPath ConfKey = "payment_static_path" // 收银台静态资源路径
	PaymentMatchMode  ConfKey = "payment_match_mode"  // 订单金额匹配模式

	RpcEndpointPlasma         ConfKey = "rpc_endpoint_plasma"            // Plasma RPC节点
	RpcEndpointBsc            ConfKey = "rpc_endpoint_bsc"               // BSC RPC节点
	RpcEndpointSolana         ConfKey = "rpc_endpoint_solana"            // Solana RPC节点
	RpcEndpointXlayer         ConfKey = "rpc_endpoint_xlayer"            // Xlayer RPC节点
	RpcEndpointPolygon        ConfKey = "rpc_endpoint_polygon"           // Polygon RPC节点
	RpcEndpointArbitrum       ConfKey = "rpc_endpoint_arbitrum"          // Arbitrum RPC节点
	RpcEndpointEthereum       ConfKey = "rpc_endpoint_ethereum"          // Ethereum RPC节点
	RpcEndpointBase           ConfKey = "rpc_endpoint_base"              // Base RPC节点
	RpcEndpointAptos          ConfKey = "rpc_endpoint_aptos"             // APTOS RPC节点
	RpcEndpointTron           ConfKey = "rpc_endpoint_tron"              // TRON RPC节点
	RpcEndpointTronGridApiKey ConfKey = "rpc_endpoint_tron_grid_api_key" // TRON RPC节点 TronGrid Api Key

	RateSyncCoingeckoApiUrl ConfKey = "rate_sync_coingecko_api_url" // 汇率同步 Coingecko Api URL
	RateSyncCoingeckoApiKey ConfKey = "rate_sync_coingecko_api_key" // 汇率同步 Coingecko Api Key
	RateSyncInterval        ConfKey = "rate_sync_interval"          // 汇率同步间隔，单位秒
	NotifyMaxRetry          ConfKey = "notify_max_retry"            // 最大重试次数，订单回调失败
	BlockHeightMaxDiff      ConfKey = "block_height_max_diff"       // 区块高度最大差值，超过此值则以当前区块高度为准，重新开始扫描
	BlockOffsetConfirm      ConfKey = "block_offset_confirm"        // 区块偏移确认数，扫描时以当前区块高度减去此偏移量为准，避免重链导致的订单回调失败

	NotifierParams  ConfKey = "notifier_params"  // 通知参数 (token, chat_id, email
	NotifierChannel ConfKey = "notifier_channel" // 通知渠道 (telegram, wechat, email

	SystemInstallLock ConfKey = "system_install_lock" // 系统安装锁
)
const (
	CNY Fiat = "CNY"
	USD Fiat = "USD"
	JPY Fiat = "JPY"
	EUR Fiat = "EUR"
	GBP Fiat = "GBP"
)
const (
	USDT Crypto = "USDT"
	USDC Crypto = "USDC"
	TRX  Crypto = "TRX"
	BNB  Crypto = "BNB"
	ETH  Crypto = "ETH"
)
const (
	Classic   MatchMode = "classic"    // 经典模式，精确匹配
	HasPrefix MatchMode = "has_prefix" // 前缀匹配，允许多付
	RoundOff  MatchMode = "round_off"  // 数值修约，四舍五入，允许容错
)

// USD 交易类型常见扫描范围
var usdGeneralRange = Range{
	MinAmount: decimal.NewFromFloat(0.01),
	MaxAmount: decimal.NewFromFloat(1000000),
}

// registry 交易类型注册表【由init函数自动维护】
var networkTradesMap = make(map[Network][]TradeType)
var networkEndpointMap = make(map[Network]ConfKey)
var contractTradeMap = make(map[string]TradeType)
var contractDecimalMap = make(map[string]int32)
var tradeAmountRangeMap = make(map[TradeType]Range)
var explorerUrlMap = make(map[TradeType]string)
var cryptoAtomKeys = make(map[Crypto]ConfKey)
