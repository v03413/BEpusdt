package model

import (
	"fmt"

	"github.com/shopspring/decimal"
	"github.com/v03413/bepusdt/app/conf"
)

type Range struct {
	MinAmount decimal.Decimal // 最小扫描金额
	MaxAmount decimal.Decimal // 最大扫描金额
}

type TradeTypeConf struct {
	Alias       string  // 类型别名，主要用户前端展示
	Network     Network // 所属区块链网络
	Crypto      Crypto  // 币种类型
	Native      bool    // 是否原生币
	Contract    string  // 合约地址，原生币为空
	Decimal     int32   // 小数位
	AmountRange Range   // 合法数额范围；这个范围有两个一个创建订单时[法币范围]，后台可配置，另一个则是扫块时[数额范围]，目前偷懒全部写死一个大概合理的范围，后面有问题再说...
	EndpointKey ConfKey // RPC 端点配置键
}

// USD 交易类型常见扫描范围
var usdGeneralRange = Range{
	MinAmount: decimal.NewFromFloat(0.01),
	MaxAmount: decimal.NewFromFloat(1000000),
}

// TradeType 交易类型，当下类型开始增多，以前的写法乱七八糟满天飞，现在这里统一管理、尽量收缩配置项
var registry = map[TradeType]TradeTypeConf{
	UsdtPlasma: {
		Alias:       "USDT・Plasma",
		Network:     conf.Plasma,
		Crypto:      USDT,
		Contract:    "0xb8ce59fc3717ada4c02eadf9682a9e934f625ebb",
		Decimal:     -6,
		AmountRange: usdGeneralRange,
		EndpointKey: RpcEndpointPlasma,
	},
	UsdtTrc20: {
		Alias:       "USDT・TRC20",
		Network:     conf.Tron,
		Crypto:      USDT,
		Contract:    "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t",
		Decimal:     conf.UsdtTronDecimals,
		AmountRange: usdGeneralRange,
		EndpointKey: RpcEndpointTron,
	},
	UsdtErc20: {
		Alias:       "USDT・ERC20",
		Network:     conf.Ethereum,
		Crypto:      USDT,
		Contract:    conf.UsdtErc20,
		Decimal:     conf.UsdtEthDecimals,
		AmountRange: usdGeneralRange,
		EndpointKey: RpcEndpointEthereum,
	},
	UsdtBep20: {
		Alias:       "USDT・BEP20",
		Network:     conf.Bsc,
		Crypto:      USDT,
		Contract:    conf.UsdtBep20,
		Decimal:     conf.UsdtBscDecimals,
		AmountRange: usdGeneralRange,
		EndpointKey: RpcEndpointBsc,
	},
	UsdtAptos: {
		Alias:       "USDT・Aptos",
		Network:     conf.Aptos,
		Crypto:      USDT,
		Contract:    conf.UsdtAptos,
		Decimal:     conf.UsdtAptosDecimals,
		AmountRange: usdGeneralRange,
		EndpointKey: RpcEndpointAptos,
	},
	UsdtXlayer: {
		Alias:       "USDT・X Layer",
		Network:     conf.Xlayer,
		Crypto:      USDT,
		Contract:    conf.UsdtXlayer,
		Decimal:     conf.UsdtXlayerDecimals,
		AmountRange: usdGeneralRange,
		EndpointKey: RpcEndpointXlayer,
	},
	UsdtSolana: {
		Alias:       "USDT・Solana",
		Network:     conf.Solana,
		Crypto:      USDT,
		Contract:    conf.UsdtSolana,
		Decimal:     conf.UsdtSolanaDecimals,
		AmountRange: usdGeneralRange,
		EndpointKey: RpcEndpointSolana,
	},
	UsdtPolygon: {
		Alias:       "USDT・Polygon",
		Network:     conf.Polygon,
		Crypto:      USDT,
		Contract:    conf.UsdtPolygon,
		Decimal:     conf.UsdtPolygonDecimals,
		AmountRange: usdGeneralRange,
		EndpointKey: RpcEndpointPolygon,
	},
	UsdtArbitrum: {
		Alias:       "USDT・Arbitrum",
		Network:     conf.Arbitrum,
		Crypto:      USDT,
		Contract:    conf.UsdtArbitrum,
		Decimal:     conf.UsdtArbitrumDecimals,
		AmountRange: usdGeneralRange,
		EndpointKey: RpcEndpointArbitrum,
	},
	UsdcErc20: {
		Alias:       "USDC・ERC20",
		Network:     conf.Ethereum,
		Crypto:      USDC,
		Contract:    conf.UsdcErc20,
		Decimal:     conf.UsdcEthDecimals,
		AmountRange: usdGeneralRange,
		EndpointKey: RpcEndpointEthereum,
	},
	UsdcBep20: {
		Alias:       "USDC・BEP20",
		Network:     conf.Bsc,
		Crypto:      USDC,
		Contract:    conf.UsdcBep20,
		Decimal:     conf.UsdcBscDecimals,
		AmountRange: usdGeneralRange,
		EndpointKey: RpcEndpointBsc,
	},
	UsdcXlayer: {
		Alias:       "USDC・X Layer",
		Network:     conf.Xlayer,
		Crypto:      USDC,
		Contract:    conf.UsdcXlayer,
		Decimal:     conf.UsdcXlayerDecimals,
		AmountRange: usdGeneralRange,
		EndpointKey: RpcEndpointXlayer,
	},
	UsdcPolygon: {
		Alias:       "USDC・Polygon",
		Network:     conf.Polygon,
		Crypto:      USDC,
		Contract:    conf.UsdcPolygon,
		Decimal:     conf.UsdcPolygonDecimals,
		AmountRange: usdGeneralRange,
		EndpointKey: RpcEndpointPolygon,
	},
	UsdcArbitrum: {
		Alias:       "USDC・Arbitrum",
		Network:     conf.Arbitrum,
		Crypto:      USDC,
		Contract:    conf.UsdcArbitrum,
		Decimal:     conf.UsdcArbitrumDecimals,
		AmountRange: usdGeneralRange,
		EndpointKey: RpcEndpointArbitrum,
	},
	UsdcBase: {
		Alias:       "USDC・Base",
		Network:     conf.Base,
		Crypto:      USDC,
		Contract:    conf.UsdcBase,
		Decimal:     conf.UsdcBaseDecimals,
		AmountRange: usdGeneralRange,
		EndpointKey: RpcEndpointBase,
	},
	UsdcTrc20: {
		Alias:       "USDC・TRC20",
		Network:     conf.Tron,
		Crypto:      USDC,
		Contract:    "TEkxiTehnzSmSe2XqrBj4w32RUN966rdz8",
		Decimal:     conf.UsdcTronDecimals,
		AmountRange: usdGeneralRange,
		EndpointKey: RpcEndpointTron,
	},
	UsdcSolana: {
		Alias:       "USDC・Solana",
		Network:     conf.Solana,
		Crypto:      USDC,
		Contract:    conf.UsdcSolana,
		Decimal:     conf.UsdcSolanaDecimals,
		AmountRange: usdGeneralRange,
		EndpointKey: RpcEndpointSolana,
	},
	UsdcAptos: {
		Alias:       "USDC・Aptos",
		Network:     conf.Aptos,
		Crypto:      USDC,
		Contract:    conf.UsdcAptos,
		Decimal:     conf.UsdcAptosDecimals,
		AmountRange: usdGeneralRange,
		EndpointKey: RpcEndpointAptos,
	},
	TronTrx: {
		Alias:   "Tron・Trx",
		Network: conf.Tron,
		Crypto:  TRX,
		Native:  true,
		Decimal: -6,
		AmountRange: Range{
			MinAmount: decimal.NewFromFloat(0.1),
			MaxAmount: decimal.NewFromFloat(1000000),
		},
		EndpointKey: RpcEndpointTron,
	},
	EthereumEth: {
		Alias:   "Ethereum・Eth",
		Network: conf.Ethereum,
		Crypto:  ETH,
		Native:  true,
		Decimal: conf.EthereumEthDecimals,
		AmountRange: Range{
			MinAmount: decimal.NewFromFloat(0.000001),
			MaxAmount: decimal.NewFromFloat(1000000),
		},
		EndpointKey: RpcEndpointEthereum,
	},
	BscBnb: {
		Alias:   "Bsc・Bnb",
		Network: conf.Bsc,
		Crypto:  BNB,
		Native:  true,
		Decimal: conf.BscBnbDecimals,
		AmountRange: Range{
			MinAmount: decimal.NewFromFloat(0.00001),
			MaxAmount: decimal.NewFromFloat(1000000),
		},
		EndpointKey: RpcEndpointBsc,
	},
}

var networkTradesMap = make(map[Network][]TradeType)
var networkEndpointMap = make(map[Network]ConfKey)
var contractTradeMap = make(map[string]TradeType)
var contractDecimalMap = make(map[string]int32)
var tradeAmountRangeMap = make(map[TradeType]Range)

func init() {
	for t, c := range registry {
		networkTradesMap[c.Network] = append(networkTradesMap[c.Network], t)
		networkEndpointMap[c.Network] = c.EndpointKey
		if c.Contract != "" {
			contractDecimalMap[c.Contract] = c.Decimal
			contractTradeMap[c.Contract] = t
		}
		if c.AmountRange != (Range{}) {
			tradeAmountRangeMap[t] = c.AmountRange
		}
	}
}

func IsSupportedTradeType(t TradeType) bool {
	_, ok := registry[t]

	return ok
}

func GetCrypto(t TradeType) (Crypto, error) {
	if c, ok := registry[t]; ok {
		return c.Crypto, nil
	}
	return "", fmt.Errorf("unsupported trade type: %s", t)
}

func GetAllAlias() map[string]string {
	var alias = make(map[string]string)

	for t, c := range registry {
		alias[string(t)] = c.Alias
	}

	return alias
}

func GetNetworkTrades(n Network) []TradeType {
	list, ok := networkTradesMap[n]
	if !ok {
		return []TradeType{}
	}

	return list
}

func GetContractTrade(addr string) (TradeType, bool) {
	t, ok := contractTradeMap[addr]

	return t, ok
}

func GetContractDecimal(addr string) int32 {
	if d, ok := contractDecimalMap[addr]; ok {

		return d
	}

	return -6
}

func IsAmountValid(t TradeType, d decimal.Decimal) bool {
	r, ok := tradeAmountRangeMap[t]
	if !ok {

		return false
	}

	if r.MaxAmount.Cmp(d) < 0 {

		return false
	}

	if r.MinAmount.Cmp(d) > 0 {

		return false
	}

	return true
}

func Endpoint(net Network) string {
	if endpointKey, ok := networkEndpointMap[net]; ok {
		return GetC(endpointKey)
	}
	return ""
}
