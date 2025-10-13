package model

const (
	AdminUsername ConfKey = "admin_username"
	AdminPassword ConfKey = "admin_password"
	AdminSecure   ConfKey = "admin_secure"
	AdminToken    ConfKey = "admin_token"
	AdminLoginIP  ConfKey = "admin_login_ip"
	AdminLoginAt  ConfKey = "admin_login_at"

	ApiAuthToken ConfKey = "api_auth_token" // API 对接令牌
	ApiAppUri    ConfKey = "api_app_uri"    // API 对接地址（收银台地址）

	AtomUSDT ConfKey = "atom_usdt"
	AtomUSDC ConfKey = "atom_usdc"
	AtomTRX  ConfKey = "atom_trx"

	MonitorMinAmount  ConfKey = "monitor_min_amount" // 监控最小金额，低于此金额的入账不进行通知
	PaymentMinAmount  ConfKey = "payment_min_amount"
	PaymentMaxAmount  ConfKey = "payment_max_amount"
	PaymentTimeout    ConfKey = "payment_timeout"     // 订单支付超时时间，单位秒
	PaymentStaticPath ConfKey = "payment_static_path" // 收银台静态资源路径

	RpcEndpointTron     ConfKey = "rpc_endpoint_tron"     // TRON RPC节点
	RpcEndpointBsc      ConfKey = "rpc_endpoint_bsc"      // BSC RPC节点
	RpcEndpointSolana   ConfKey = "rpc_endpoint_solana"   // Solana RPC节点
	RpcEndpointXlayer   ConfKey = "rpc_endpoint_xlayer"   // Xlayer RPC节点
	RpcEndpointPolygon  ConfKey = "rpc_endpoint_polygon"  // Polygon RPC节点
	RpcEndpointArbitrum ConfKey = "rpc_endpoint_arbitrum" // Arbitrum RPC节点
	RpcEndpointEthereum ConfKey = "rpc_endpoint_ethereum" // Ethereum RPC节点
	RpcEndpointBase     ConfKey = "rpc_endpoint_base"     // Base RPC节点
	RpcEndpointAptos    ConfKey = "rpc_endpoint_aptos"    // APTOS RPC节点

	RateSyncInterval ConfKey = "rate_sync_interval" // 汇率同步间隔，单位秒

	NotifyMaxRetry     ConfKey = "notify_max_retry"      // 最大重试次数，订单回调失败
	BlockHeightMaxDiff ConfKey = "block_height_max_diff" // 区块高度最大差值，超过此值则以当前区块高度为准，重新开始扫描

	NotifierParams  ConfKey = "notifier_params"  // 通知参数 (token, chat_id, email
	NotifierChannel ConfKey = "notifier_channel" // 通知渠道 (telegram, wechat, email
)

const (
	DefaultUsername = "admin"
	DefaultPassword = "$2a$10$GHecK2k0JICojz4rg4G6nO6AHi6r3lgocRp.Ob4JB/O8JTQ13Obly" // 123456
)
