<template>
  <div>
    <a-row align="center" :gutter="[0, 16]">
      <a-col :span="24">
        <a-card title="MQTT 发布设置">
          <a-alert type="info" style="margin-bottom: 16px">
            系统可将扫描到的交易信息发布到 MQTT 服务器，其它系统通过订阅实时获取数据，
            <a href="https://github.com/v03413/BEpusdt/blob/main/docs/api/mqtt.md" target="_blank">查看文档</a>
          </a-alert>
          <a-alert type="warning" style="margin-bottom: 16px">
            <template #icon><icon-exclamation-circle-fill /></template>
            目前 MQTT 通信协议为 <strong>MQTT over TCP</strong>，该功能只有在 <strong>Host</strong> 和
            <strong>Port</strong> 均配置后才会正常启用
          </a-alert>
          <a-form :model="form" :rules="rules" :style="{ width: '600px' }" @submit="onSubmit">
            <a-form-item field="mqtt_host" label="MQTT Host" extra="MQTT 服务器地址">
              <a-input v-model="form.mqtt_host" placeholder="例如：127.0.0.1" allow-clear />
            </a-form-item>

            <a-form-item field="mqtt_port" label="MQTT Port" extra="MQTT 服务器端口">
              <a-input v-model="form.mqtt_port" placeholder="例如：1883" allow-clear />
            </a-form-item>

            <a-form-item field="mqtt_user" label="用户名" extra="MQTT 连接用户名，无则留空">
              <a-input v-model="form.mqtt_user" placeholder="请输入用户名" allow-clear />
            </a-form-item>

            <a-form-item field="mqtt_pass" label="密码" extra="MQTT 连接密码，无则留空">
              <a-input-password v-model="form.mqtt_pass" placeholder="请输入密码" allow-clear />
            </a-form-item>

            <a-form-item
              field="mqtt_topic_prefix"
              label="消息路径前缀"
              extra="消息发布的 Topic 路径前缀，只允许字母、数字、下划线和斜杠，默认为 bepusdt"
            >
              <a-input v-model="form.mqtt_topic_prefix" placeholder="例如：bepusdt" allow-clear />
            </a-form-item>

            <a-form-item field="mqtt_publish_qos" label="Publish QoS" extra="消息发布服务质量等级">
              <a-radio-group v-model="form.mqtt_publish_qos">
                <a-radio value="0">0 - 最多一次</a-radio>
                <a-radio value="1">1 - 至少一次</a-radio>
                <a-radio value="2">2 - 恰好一次</a-radio>
              </a-radio-group>
            </a-form-item>

            <a-form-item field="mqtt_networks" label="区块链网络" extra="选择需要持续监听的区块链网络，多选">
              <a-checkbox-group v-model="networksSelected">
                <a-checkbox value="tron">Tron</a-checkbox>
                <a-checkbox value="bsc">Bsc</a-checkbox>
                <a-checkbox value="polygon">Polygon</a-checkbox>
                <a-checkbox value="ethereum">Ethereum</a-checkbox>
                <a-checkbox value="aptos">Aptos</a-checkbox>
                <a-checkbox value="solana">Solana</a-checkbox>
                <a-checkbox value="xlayer">XLayer</a-checkbox>
                <a-checkbox value="plasma">Plasma</a-checkbox>
                <a-checkbox value="arbitrum">Arbitrum</a-checkbox>
                <a-checkbox value="base">Base</a-checkbox>
              </a-checkbox-group>
            </a-form-item>

            <a-form-item>
              <a-button type="primary" html-type="submit">提交</a-button>
            </a-form-item>
          </a-form>
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script setup lang="ts">
import { Message } from "@arco-design/web-vue";
import { setsConfAPI } from "@/api/modules/conf/index";

const emit = defineEmits(["refresh"]);
const data = defineModel() as any;
const rules = {
  mqtt_topic_prefix: [
    {
      validator: (value: string, callback: (error?: string) => void) => {
        if (!value) return callback();
        if (!/^[a-zA-Z0-9_/]+$/.test(value)) {
          return callback("Topic 前缀只允许包含字母、数字、下划线和斜杠");
        }
        if (value.startsWith("/") || value.endsWith("/")) {
          return callback("Topic 前缀不能以斜杠开头或结尾");
        }
        if (value.includes("//")) {
          return callback("Topic 前缀不能包含连续斜杠");
        }
        callback();
      }
    }
  ]
};

const form = ref({
  mqtt_host: "",
  mqtt_port: "",
  mqtt_user: "",
  mqtt_pass: "",
  mqtt_publish_qos: "0",
  mqtt_networks: "",
  mqtt_topic_prefix: "bepusdt"
});

const networksSelected = computed({
  get() {
    return form.value.mqtt_networks ? form.value.mqtt_networks.split(",").filter(Boolean) : [];
  },
  set(val: string[]) {
    form.value.mqtt_networks = val.join(",");
  }
});

const onSubmit = async ({ errors }: ArcoDesign.ArcoSubmit) => {
  if (errors) return;

  await setsConfAPI([
    { key: "mqtt_host", value: form.value.mqtt_host },
    { key: "mqtt_port", value: form.value.mqtt_port },
    { key: "mqtt_user", value: form.value.mqtt_user },
    { key: "mqtt_pass", value: form.value.mqtt_pass },
    { key: "mqtt_publish_qos", value: form.value.mqtt_publish_qos },
    { key: "mqtt_networks", value: form.value.mqtt_networks },
    { key: "mqtt_topic_prefix", value: form.value.mqtt_topic_prefix || "bepusdt" }
  ]);

  Message.success("保存成功");
  emit("refresh");
};

watch(
  () => data.value,
  () => {
    form.value.mqtt_host = data.value.mqtt_host ?? "";
    form.value.mqtt_port = data.value.mqtt_port ?? "";
    form.value.mqtt_user = data.value.mqtt_user ?? "";
    form.value.mqtt_pass = data.value.mqtt_pass ?? "";
    form.value.mqtt_publish_qos = data.value.mqtt_publish_qos ?? "0";
    form.value.mqtt_networks = data.value.mqtt_networks ?? "";
    form.value.mqtt_topic_prefix = data.value.mqtt_topic_prefix ?? "bepusdt";
  }
);
</script>
