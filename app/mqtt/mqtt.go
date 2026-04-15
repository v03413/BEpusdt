package mqtt

import (
	"fmt"
	"os"
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/v03413/bepusdt/app/log"
	"github.com/v03413/bepusdt/app/model"
)

type activeConf struct {
	host string
	port string
	user string
	pass string
}

var (
	client  mqtt.Client
	mu      sync.RWMutex
	curConf activeConf
)

// Reload 按最新配置连接MQTT服务器
func Reload() error {
	host := model.GetC(model.MqttHost)
	port := model.GetC(model.MqttPort)
	if host == "" || port == "" {
		return nil // 未配置 MQTT，跳过
	}

	newConf := activeConf{
		host: host,
		port: port,
		user: model.GetC(model.MqttUser),
		pass: model.GetC(model.MqttPass),
	}

	mu.Lock()
	defer mu.Unlock()
	if newConf == curConf { // 配置未变更，无需重连
		return nil
	}

	// 断开旧连接（给 250ms 优雅退出）
	if client != nil {
		client.Disconnect(250)
		log.Info("🔄 MQTT 配置变动，正在连接...")
	}

	opts := mqtt.NewClientOptions()
	// 暂时只支持 tcp 模式
	opts.AddBroker(fmt.Sprintf("tcp://%s:%s", newConf.host, newConf.port))
	opts.SetUsername(newConf.user)
	opts.SetPassword(newConf.pass)
	opts.SetClientID(fmt.Sprintf("BEpusdt %d", os.Getpid()))
	opts.SetAutoReconnect(true)
	opts.SetOrderMatters(false)
	opts.SetOnConnectHandler(onConnectHandler)
	opts.SetConnectionLostHandler(onConnectionLost)

	newClient := mqtt.NewClient(opts)
	if token := newClient.Connect(); token.Wait() && token.Error() != nil {
		return fmt.Errorf("MQTT 连接失败: %s", token.Error())
	}

	client = newClient
	curConf = newConf

	return nil
}

func onConnectHandler(c mqtt.Client) {
	for topic, cb := range subscribeMap {
		c.Subscribe(topic, cb.Qos, cb.Handler).Wait()
	}

	log.Info("✅ MQTT 连接成功")
}

func onConnectionLost(_ mqtt.Client, err error) {
	log.Warn(fmt.Sprintf("❌ MQTT 连接断开: %s", err.Error()))
}
