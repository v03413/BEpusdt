package model

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/v03413/bepusdt/app/utils"
	"github.com/v03413/go-cache"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var confCache sync.Map

type Conf struct {
	K ConfKey `gorm:"column:k;type:varchar(32);not null;primaryKey" json:"key"`
	V string  `gorm:"column:v;type:varchar(255);not null" json:"val"`
}

func (c Conf) TableName() string {

	return "bep_conf"
}

func SetK(k ConfKey, v string) {
	if err = Db.Transaction(func(db *gorm.DB) error {
		if err2 := db.Where("`k` = ?", k).Delete(&Conf{}).Error; err2 != nil {

			return err2
		}
		if err2 := db.Create(&Conf{K: k, V: v}).Error; err2 != nil {

			return err2
		}

		defer RefreshC()

		return nil
	}); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, fmt.Sprintf("设置配置项 %s 错误：%s", k, err.Error()))
	}
}

func GetK(k ConfKey) string {
	var row Conf

	var tx = Db.Where("k = ?", k).Limit(1).Find(&row)
	if tx.Error == nil {

		return row.V
	}

	_, _ = fmt.Fprintln(os.Stderr, fmt.Sprintf("获取配置项 %s 错误：%s", k, tx.Error.Error()))

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

func CheckoutCounter(host, id string) string {
	uri := GetK(ApiAppUri)
	if uri == "" {
		uri = host
	}

	return fmt.Sprintf("%s/pay/checkout-counter/%s", uri, id)
}

func ConfInit() {
	var hash = utils.StrSha256(utils.Md5String(time.Now().String()))
	var secure = "/" + hash[:10]
	var token = strings.ToUpper(utils.Md5String(hash[18:28]))
	var username = hash[10:20]
	var password = hash[20:30]
	var encrypt, _ = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	var data = map[ConfKey]string{
		ApiAppUri:           "",
		ApiAuthToken:        token,
		AdminSecure:         secure,
		AdminUsername:       username,
		AdminPassword:       string(encrypt),
		RateSyncInterval:    "3600",
		AtomUSDT:            "0.01",
		AtomUSDC:            "0.01",
		AtomTRX:             "0.01",
		AtomBNB:             "0.00001",
		AtomETH:             "0.000001",
		MonitorMinAmount:    "0.01",
		PaymentMinAmount:    "0.01",
		PaymentMaxAmount:    "99999",
		RpcEndpointTron:     "grpc.trongrid.io:50051",
		RpcEndpointBsc:      "https://binance-smart-chain-public.nodies.app/",
		RpcEndpointSolana:   "https://solana-rpc.publicnode.com/",
		RpcEndpointXlayer:   "https://xlayerrpc.okx.com/",
		RpcEndpointPolygon:  "https://polygon-public.nodies.app/",
		RpcEndpointArbitrum: "https://arb1.arbitrum.io/rpc",
		RpcEndpointEthereum: "https://ethereum-public.nodies.app/",
		RpcEndpointBase:     "https://base-public.nodies.app/",
		RpcEndpointAptos:    "https://aptos-rest.publicnode.com/",
		RpcEndpointPlasma:   "https://rpc.plasma.to/",
		NotifyMaxRetry:      "10",
		BlockHeightMaxDiff:  "1000",
		PaymentTimeout:      "1200", // 20分钟
		PaymentStaticPath:   "",
		PaymentMatchMode:    string(Classic),
		SystemInstallLock:   "0",
	}

	var rows = make([]Conf, 0)
	for k, v := range data {
		rows = append(rows, Conf{K: k, V: v})
	}

	Db.Create(&rows)

	// 数据丢到缓存，前台首次访问时会展示这部分初始化信息；明文密码只这一次保存到缓存，不写入数据库
	cache.Set(string(SystemInstallLock), gin.H{
		"username": username,
		"password": password,
		"secure":   secure,
		"token":    token,
	}, -1)
}

func AuthToken() string {

	return GetK(ApiAuthToken)
}

func IsInstalled() bool {
	return GetC(SystemInstallLock) == "1"
}

func InstallLock() {
	SetK(SystemInstallLock, "1")
}

func GetInstallInfo() gin.H {
	if info, ok := cache.Get(string(SystemInstallLock)); ok {

		return info.(gin.H)
	}

	return gin.H{}
}
