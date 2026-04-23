package mqtt

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/v03413/bepusdt/app/log"
)

// Publish 发布消息。若 MQTT 未连接则静默丢弃，返回 nil。
func Publish(topic string, qos byte, retained bool, payload interface{}) mqtt.Token {
	mu.RLock()
	c := client
	mu.RUnlock()

	if c == nil {
		return nil
	}

	token := c.Publish(topic, qos, retained, payload)
	if qos > 0 {
		go func() {
			token.Wait()
			if err := token.Error(); err != nil {
				log.Warn("MQTT publish error: " + err.Error())
				return
			}
		}()
	}

	return token
}
