package mqtt

import mqtt "github.com/eclipse/paho.mqtt.golang"

// Publish 发布消息。若 MQTT 未连接则静默丢弃，返回 nil。
func Publish(topic string, qos byte, retained bool, payload interface{}) mqtt.Token {
	mu.RLock()
	c := client
	mu.RUnlock()

	if c == nil {
		return nil
	}

	return c.Publish(topic, qos, retained, payload)
}
