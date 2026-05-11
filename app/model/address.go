package model

import (
	"strings"

	"github.com/v03413/bepusdt/app/conf"
	"github.com/v03413/bepusdt/app/utils"
)

type networkAddressCodec func(string) (string, bool)

var networkAddressNormalizers = map[Network]networkAddressCodec{
	conf.Ton: utils.NormalizeTonAddress,
}

var networkAddressFormatters = map[Network]networkAddressCodec{
	conf.Ton: utils.FormatTonAddress,
}

// normalizeNetworkAddress 按网络将地址转换为内部匹配和索引使用的标准形式。
func normalizeNetworkAddress(network Network, address string) (string, bool) {
	address = strings.TrimSpace(address)
	if address == "" {
		return "", false
	}

	if normalize, ok := networkAddressNormalizers[network]; ok {
		return normalize(address)
	}

	return address, true
}

// NormalizeKnownAddress 使用已注册的网络规范化器尝试转换地址。
func NormalizeKnownAddress(address string) string {
	address = strings.TrimSpace(address)
	if address == "" {
		return ""
	}

	for _, normalize := range networkAddressNormalizers {
		if normalized, ok := normalize(address); ok {
			return normalized
		}
	}

	return address
}

// FormatNetworkAddress 按网络将内部地址转换为适合用户查看的展示格式。
func FormatNetworkAddress(network Network, address string) string {
	address = strings.TrimSpace(address)
	if address == "" {
		return ""
	}

	if format, ok := networkAddressFormatters[network]; ok {
		if formatted, ok := format(address); ok {
			return formatted
		}
	}

	return address
}

// NormalizeTradeAddress 按交易类型将地址转换为订单匹配使用的标准形式。
func NormalizeTradeAddress(t TradeType, address string) string {
	address = strings.TrimSpace(address)
	if c, ok := registry[t]; ok {
		if normalized, ok := normalizeNetworkAddress(c.Network, address); ok {
			address = normalized
		}
	}
	// 非大小写敏感的地址，统一转为小写存储
	if !AddrCaseSens(t) {
		return strings.ToLower(address)
	}

	return address
}

// FormatTradeAddress 按交易类型将内部地址转换为适合用户查看和付款的展示格式。
func FormatTradeAddress(t TradeType, address string) string {
	if c, ok := registry[t]; ok {
		return FormatNetworkAddress(c.Network, address)
	}

	return strings.TrimSpace(address)
}

// FormatOrderPaymentAddress 返回订单对外展示的付款地址。
func FormatOrderPaymentAddress(order Order) string {
	address := order.Address

	return FormatTradeAddress(order.TradeType, address)
}
