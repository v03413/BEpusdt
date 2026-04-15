package task

import (
	"context"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
	"github.com/v03413/bepusdt/app/model"
	appMqtt "github.com/v03413/bepusdt/app/mqtt"
	"github.com/v03413/go-cache"
)

type mqttTransfer struct {
	Network     string          `json:"network"`
	TxHash      string          `json:"tx_hash"`
	Amount      decimal.Decimal `json:"amount"`
	FromAddress string          `json:"from_address"`
	RecvAddress string          `json:"recv_address"`
	Timestamp   int64           `json:"timestamp"`
	TradeType   model.TradeType `json:"trade_type"`
	BlockNum    int             `json:"block_num"`
}

func init() {
	Register(Task{Duration: time.Second * 5, Callback: mqttWatcher})
}

func mqttWatcher(_ context.Context) {
	if err := appMqtt.Reload(); err != nil {
		log.Printf("❌ MQTT 连接失败: %s\n", err.Error())
	}

	networks := strings.Split(model.GetC(model.MqttNetworks), ",")
	for _, n := range networks {
		cache.Set("mqtt_subscribed_"+n, true, time.Second*8)
	}
}

func mqttPublish(t transfer) {
	if !mqttSubscribed(t.Network) {
		return
	}

	var qos = cast.ToUint8(model.GetC(model.MqttPublishQos))
	var topic = model.GetC(model.MqttTopicPrefix) + "/transfer/" + t.Network
	var data = mqttTransfer{
		Network:     t.Network,
		TxHash:      t.TxHash,
		Amount:      t.Amount,
		FromAddress: t.FromAddress,
		RecvAddress: t.RecvAddress,
		Timestamp:   t.Timestamp.Unix(),
		TradeType:   t.TradeType,
		BlockNum:    t.BlockNum,
	}
	go func() {
		payload, _ := json.Marshal(data)

		appMqtt.Publish(topic, qos, false, payload)
	}()
}

func mqttSubscribed(n string) bool {
	_, found := cache.Get("mqtt_subscribed_" + n)
	return found
}
