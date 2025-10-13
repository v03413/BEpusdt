<template>
  <div class="snow-page">
    <a-spin :loading="loading" tip="loading...">
      <a-card :bordered="false">
        <a-row align="center">
          <a-col :span="2">
            <div>
              <a-avatar :size="100" trigger-type="mask">
                <img alt="avatar" src="https://avatars.githubusercontent.com/u/49953737?v=4" />
                <template #trigger-icon>
                  <IconEdit />
                </template>
              </a-avatar>
            </div>
          </a-col>
          <a-col :span="22">
            <a-space direction="vertical" size="large">
              <a-descriptions :column="2" title="基本信息" :align="{ label: 'right' }">
                <a-descriptions-item label="管理员">
                  {{ Conf.admin_username }}
                </a-descriptions-item>
                <a-descriptions-item label="登录时间">
                  {{ Conf.last_login_at }}
                </a-descriptions-item>
                <a-descriptions-item label="登录IP">
                  {{ Conf.last_login_ip || "暂无" }}
                </a-descriptions-item>
              </a-descriptions>
            </a-space>
          </a-col>
        </a-row>
      </a-card>
      <a-card class="margin-top" :bordered="false">
        <a-row align="center">
          <a-col :span="24">
            <a-tabs :type="type" :size="size" :active-key="activeTabs" @change="onChangeTab">
              <a-tab-pane key="1" title="基本设置">
                <Info v-model="Conf" @refresh="refresh" />
              </a-tab-pane>
              <a-tab-pane key="2" title="交易通知">
                <Notifier v-model="Conf" @refresh="refresh" />
              </a-tab-pane>
              <a-tab-pane key="4" title="API设置">
                <Api v-model="Conf" @refresh="refresh" />
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
import Info from "@/views/system/base/info.vue";
import Security from "@/views/system/base/security.vue";
import Api from "@/views/system/base/api.vue";
import { getsConfAPI } from "@/api/modules/conf/index";
import Notifier from "./notifier.vue";

const route = useRoute();
const type = ref("rounded");
const size = ref("medium");
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
        "api_app_uri",
        "api_auth_token",
        "admin_username",
        "admin_secure",
        "block_height_max_diff",
        "last_login_at",
        "last_login_ip",
        "notify_max_retry",
        "payment_max_amount",
        "payment_min_amount",
        "payment_static_path",
        "payment_timeout",
        "monitor_min_amount",
        "notifier_params",
        "notifier_channel"
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
.margin-top {
  margin-top: $padding;
}
</style>
