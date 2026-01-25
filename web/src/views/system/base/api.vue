<template>
  <div>
    <a-row align="center" :gutter="[0, 16]">
      <a-col :span="24">
        <a-card title="API设置">
          <a-form :model="form" :rules="rules" :style="{ width: '600px' }" @submit="onSubmit">
            <a-form-item field="api_auth_token" label="对接令牌" extra="API对接的身份验证令牌，请妥善保管">
              <a-input-group class="token-input-group">
                <a-input-password v-model="form.api_auth_token" placeholder="请输入 Auth Token" readonly />
                <a-button type="primary" @click="handleResetToken">重置</a-button>
              </a-input-group>
            </a-form-item>

            <a-form-item field="api_app_uri" label="应用URI" extra="API对接的应用URI,前端收银台地址">
              <a-input v-model="form.api_app_uri" placeholder="http(s)://your-host-uri" allow-clear />
            </a-form-item>

            <a-form-item
              field="payment_static_path"
              label="收银台静态资源"
              extra="收银台静态资源路径,可通过此功能自定义前端收银台样式;不懂请勿修改,否则可能导致收银台异常!【修改重启生效】"
            >
              <a-input v-model="form.payment_static_path" placeholder="/var/lib/bepusdt/payment/" allow-clear />
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
import { setsConfAPI, resetApiAuthToken } from "@/api/modules/conf/index";

const emit = defineEmits(["refresh"]);
const data = defineModel() as any;
const form = ref({ api_auth_token: "", api_app_uri: "", payment_static_path: "" });
const rules = {};

const handleResetToken = async () => {
  try {
    await resetApiAuthToken({});
    Message.success("令牌重置成功");
    emit("refresh");
  } catch {
    Message.error("令牌重置失败");
  }
};

const onSubmit = async ({ errors }: ArcoDesign.ArcoSubmit) => {
  if (errors) return;

  await setsConfAPI([
    {
      key: "api_app_uri",
      value: form.value.api_app_uri
    },
    {
      key: "payment_static_path",
      value: form.value.payment_static_path
    }
  ]);

  Message.success("保存成功");

  emit("refresh");
};

watch(
  () => data.value,
  () => {
    form.value.api_auth_token = data.value.api_auth_token;
    form.value.api_app_uri = data.value.api_app_uri;
    form.value.payment_static_path = data.value.payment_static_path;
  }
);
</script>

<style lang="scss" scoped>
.row-title {
  font-size: $font-size-title-1;
}

.token-input-group {
  width: 100%;

  :deep(.arco-input-wrapper) {
    flex: 1;
  }

  :deep(.arco-btn) {
    flex-shrink: 0;
  }
}
</style>
