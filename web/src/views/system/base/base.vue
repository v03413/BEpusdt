<template>
  <div class="snow-page">
    <a-spin :loading="loading" class="container" tip="loading...">
      <a-card :bordered="false">
        <a-row align="center">
          <a-col :span="isMobile ? 24 : 2">
            <div :class="{ center: isMobile }">
              <a-avatar :size="100" trigger-type="mask">
                <img alt="avatar" src="https://avatars.githubusercontent.com/u/49953737?v=4" />
                <template #trigger-icon>
                  <IconEdit />
                </template>
              </a-avatar>
            </div>
          </a-col>
          <a-col :span="isMobile ? 24 : 22">
            <a-space direction="vertical" size="large" fill class="base-profile-space">
              <a-descriptions :column="descriptionsColumn(1, 2)" title="基本信息" :align="{ label: 'right' }">
                <a-descriptions-item label="管理员">
                  {{ Conf.admin_username }}
                </a-descriptions-item>
                <a-descriptions-item label="登录时间">
                  {{ Conf.admin_login_at }}
                </a-descriptions-item>
                <a-descriptions-item label="登录IP">
                  {{ Conf.admin_login_ip || "暂无" }}
                </a-descriptions-item>
              </a-descriptions>
            </a-space>
          </a-col>
        </a-row>
      </a-card>
      <a-card class="margin-top" :bordered="false">
        <a-row align="center">
          <a-col :span="24">
            <a-tabs :type="tabsType" :size="tabsSize" :active-key="activeTabs" @change="onChangeTab">
              <a-tab-pane key="1" title="基本设置">
                <Info v-model="Conf" @refresh="refresh" />
              </a-tab-pane>
              <a-tab-pane key="2" title="交易通知">
                <Notifier v-model="Conf" @refresh="refresh" />
              </a-tab-pane>
              <a-tab-pane key="4" title="API设置">
                <Api v-model="Conf" @refresh="refresh" />
              </a-tab-pane>
              <a-tab-pane key="5" title="MQTT设置">
                <Mqtt v-model="Conf" @refresh="refresh" />
              </a-tab-pane>
              <a-tab-pane key="3" title="安全设置">
                <Security v-model="Conf" @refresh="refresh" />
              </a-tab-pane>
            </a-tabs>
          </a-col>
        </a-row>
      </a-card>
    </a-spin>
  </div>
</template>

<script setup lang="ts">
import Info from "./components/info.vue";
import Security from "./components/security.vue";
import Api from "./components/api.vue";
import { getsConfAPI } from "@/api/modules/conf/index";
import Notifier from "./components/notifier.vue";
import Mqtt from "./components/mqtt.vue";
import { useDevicesSize } from "@/hooks/useDevicesSize";
import { useLayoutModel } from "@/hooks/useLayoutModel";

const route = useRoute();
const { isMobile } = useDevicesSize();
const { descriptionsColumn } = useLayoutModel();
const tabsType = computed(() => (isMobile.value ? "line" : "rounded"));
const tabsSize = computed(() => (isMobile.value ? "small" : "medium"));
const activeTabs = ref(route.query.type || "1");

const onChangeTab = (e: string) => {
  activeTabs.value = e;
};

const refresh = () => {
  getConf();
};

const loading = ref<boolean>(false);
const Conf = ref<any>({});
const getConf = async () => {
  try {
    loading.value = true;

    let data = await getsConfAPI({
      keys: [
        "payment_match_mode",
        "api_app_uri",
        "api_auth_token",
        "admin_username",
        "admin_secure",
        "block_height_max_diff",
        "block_offset_confirm",
        "admin_login_at",
        "admin_login_ip",
        "notify_max_retry",
        "payment_max_amount",
        "payment_min_amount",
        "payment_static_path",
        "payment_timeout",
        "notifier_params",
        "notifier_channel",
        "mqtt_host",
        "mqtt_port",
        "mqtt_user",
        "mqtt_pass",
        "mqtt_publish_qos",
        "mqtt_networks",
        "mqtt_topic_prefix"
      ]
    });
    Conf.value = data.data;
  } finally {
    loading.value = false;
  }
};

getConf();
</script>

<style lang="scss" scoped>
.container {
  display: flex;
  flex-direction: column;
}

.margin-top {
  margin-top: $padding;
}

.center {
  display: flex;
  justify-content: center;
  margin-bottom: 16px;
}

:deep(.base-setting-form) {
  width: 100%;
  max-width: 600px;
  min-width: 0;
}
</style>
