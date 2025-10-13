package model

import (
	"fmt"
	"sync"
	"time"

	"github.com/shopspring/decimal"
	"github.com/v03413/bepusdt/app/conf"
	"github.com/v03413/bepusdt/app/utils"
)

type ConfKey string

var confCache sync.Map

type Conf struct {
	K ConfKey `gorm:"column:k;type:varchar(32);primaryKey" json:"key"`
	V string  `gorm:"column:v;type:varchar(255)" json:"val"`
}

func (c Conf) TableName() string {

	return "bep_conf"
}

func SetK(k ConfKey, v string) {
	Db.Exec("REPLACE INTO bep_conf (k, v) VALUES (?, ?)", k, v)
}

func GetK(k ConfKey) string {
	var row Conf

	var tx = Db.Where("k = ?", k).Limit(1).Find(&row)
	if tx.Error == nil {

		return row.V
	}

	return ""
}

func GetVs(keys []ConfKey) map[ConfKey]string {
	var rows = make([]Conf, 0)
	Db.Where("k IN ?", keys).Find(&rows)

	var result = make(map[ConfKey]string)
	for _, row := range rows {
		result[row.K] = row.V
	}

	for _, k := range keys {
		if _, ok := result[k]; !ok {
			result[k] = ""
		}
	}

	return result
}

// GetC 从缓存获取配置，适用于高频读取，依赖 RefreshC 刷新缓存
func GetC(k ConfKey) string {
	value, ok := confCache.Load(k)
	if !ok {
		return ""
	}

	return value.(string)
}

func RefreshC() {
	var rows = make([]Conf, 0)
	Db.Find(&rows)

	for _, row := range rows {
		confCache.Store(row.K, row.V)
	}
}

func Payment() (decimal.Decimal, decimal.Decimal) {
	var vs = GetVs([]ConfKey{
		PaymentMinAmount,
		PaymentMaxAmount,
	})

	minAmount, _ := decimal.NewFromString(vs[PaymentMinAmount])
	maxAmount, _ := decimal.NewFromString(vs[PaymentMaxAmount])

	return minAmount, maxAmount
}

func CheckoutCounter(host, id string) string {
	uri := GetK(ApiAppUri)
	if uri == "" {
		uri = host
	}

	return fmt.Sprintf("%s/pay/checkout-counter/%s", uri, id)
}

func ConfInit() {
	var data = map[ConfKey]string{
		ApiAppUri:           "",
		ApiAuthToken:        utils.Md5String(time.Now().String()),
		AdminUsername:       DefaultUsername,
		AdminPassword:       DefaultPassword,
		AdminSecure:         utils.Md5String(time.Now().String())[:10],
		RateSyncInterval:    "3600",
		AtomUSDT:            "0.01",
		AtomUSDC:            "0.01",
		AtomTRX:             "0.01",
		MonitorMinAmount:    "0.01",
		PaymentMinAmount:    "0.01",
		PaymentMaxAmount:    "99999",
		RpcEndpointTron:     "18.141.79.38:50051",
		RpcEndpointBsc:      "https://binance-smart-chain-public.nodies.app/",
		RpcEndpointSolana:   "https://solana-rpc.publicnode.com/",
		RpcEndpointXlayer:   "https://xlayerrpc.okx.com/",
		RpcEndpointPolygon:  "https://polygon-public.nodies.app/",
		RpcEndpointArbitrum: "https://arb1.arbitrum.io/rpc",
		RpcEndpointEthereum: "https://ethereum-public.nodies.app/",
		RpcEndpointBase:     "https://base-public.nodies.app/",
		RpcEndpointAptos:    "https://aptos-rest.publicnode.com/",
		NotifyMaxRetry:      "10",
		BlockHeightMaxDiff:  "1000",
		PaymentTimeout:      "1200", // 20分钟
		PaymentStaticPath:   "",
	}

	var rows = make([]Conf, 0)
	for k, v := range data {
		rows = append(rows, Conf{K: k, V: v})
	}

	Db.Create(&rows)
}

func AuthToken() string {

	return GetK(ApiAuthToken)
}

func Endpoint(net string) string {
	switch net {
	case conf.Tron:
		return GetC(RpcEndpointTron)
	case conf.Bsc:
		return GetC(RpcEndpointBsc)
	case conf.Solana:
		return GetC(RpcEndpointSolana)
	case conf.Xlayer:
		return GetC(RpcEndpointXlayer)
	case conf.Polygon:
		return GetC(RpcEndpointPolygon)
	case conf.Arbitrum:
		return GetC(RpcEndpointArbitrum)
	case conf.Ethereum:
		return GetC(RpcEndpointEthereum)
	case conf.Base:
		return GetC(RpcEndpointBase)
	case conf.Aptos:
		return GetC(RpcEndpointAptos)
	}

	return ""
}
