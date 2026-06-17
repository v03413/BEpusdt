<template>
  <a-row align="center" :gutter="[0, 16]">
    <a-col :span="24">
      <a-card title="API设置">
        <a-alert type="info" style="margin-bottom: 16px">
          系统兼容彩虹易支付 <strong>submit.php</strong> 接口收单，对接时 PID 固定为 <strong>1000</strong>，KEY
          则是和对接令牌保持一致。
        </a-alert>
        <a-form :model="form" :rules="rules" :layout="layoutMode" class="base-setting-form" @submit="onSubmit">
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
            field="payment_template"
            label="收银台模板"
            extra="官方默认保持原有付款页；狼哥设计使用内置 LangGe design；自定义模板使用下方静态资源路径【修改重启生效】"
          >
            <a-select v-model="form.payment_template" placeholder="请选择收银台模板" :fallback-option="false">
              <a-option value="official">官方默认</a-option>
              <a-option value="wolf">狼哥设计</a-option>
              <a-option value="custom">自定义模板</a-option>
            </a-select>
          </a-form-item>

          <a-form-item
            v-if="form.payment_template === 'wolf'"
            field="payment_template_language"
            label="默认语言"
            extra="仅狼哥设计模板生效；用户手动切换后会优先使用用户选择"
          >
            <a-select v-model="form.payment_template_language" placeholder="请选择默认语言" :fallback-option="false">
              <a-option value="auto">跟随浏览器</a-option>
              <a-option value="zh">简体中文</a-option>
              <a-option value="zh-Hant">繁體中文</a-option>
              <a-option value="en">English</a-option>
              <a-option value="ru">Русский</a-option>
              <a-option value="vi">Tiếng Việt</a-option>
              <a-option value="tr">Türkçe</a-option>
              <a-option value="ja">日本語</a-option>
              <a-option value="ko">한국어</a-option>
            </a-select>
          </a-form-item>

          <a-form-item
            v-if="form.payment_template === 'custom'"
            field="payment_static_path"
            label="收银台静态资源"
            extra="收银台静态资源路径,可通过此功能自定义前端收银台样式;不懂请勿修改,否则可能导致收银台异常!【修改重启生效】"
          >
            <a-input v-model="form.payment_static_path" placeholder="/var/lib/bepusdt/payment/" allow-clear />
          </a-form-item>

          <a-form-item>
            <a-space>
              <a-button type="primary" html-type="submit">提交</a-button>
            </a-space>
          </a-form-item>
        </a-form>
      </a-card>
    </a-col>
  </a-row>
</template>

<script setup lang="ts">
import { useDevicesSize } from "@/hooks/useDevicesSize";

import { Message } from "@arco-design/web-vue";
import { setsConfAPI, resetApiAuthToken } from "@/api/modules/conf/index";

const emit = defineEmits(["refresh"]);
const data = defineModel() as any;
const { isMobile } = useDevicesSize();
const layoutMode = computed(() => (isMobile.value ? "vertical" : "horizontal"));

const form = ref({
  api_auth_token: "",
  api_app_uri: "",
  payment_template: "official",
  payment_template_language: "auto",
  payment_static_path: ""
});
const rules = {};
const paymentTemplateModes = ["official", "wolf", "custom"];
const paymentTemplateLanguages = ["auto", "zh", "zh-Hant", "en", "ru", "vi", "tr", "ja", "ko"];

const normalizePaymentTemplate = (value: string, staticPath: string) => {
  if (paymentTemplateModes.includes(value)) {
    return value;
  }

  return staticPath ? "custom" : "official";
};

const normalizePaymentTemplateLanguage = (value: string) => {
  if (paymentTemplateLanguages.includes(value)) {
    return value;
  }

  return "auto";
};

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
      key: "payment_template",
      value: form.value.payment_template
    },
    {
      key: "payment_template_language",
      value: form.value.payment_template_language
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
    form.value.payment_template = normalizePaymentTemplate(data.value.payment_template, form.value.payment_static_path);
    form.value.payment_template_language = normalizePaymentTemplateLanguage(data.value.payment_template_language);
  }
);
</script>

<style lang="scss" scoped>
.row-title {
  font-size: $font-size-title-1;
}

.token-input-group {
  width: 100%;
  min-width: 0;

  :deep(.arco-input-wrapper) {
    flex: 1;
    min-width: 0;
  }

  :deep(.arco-btn) {
    flex-shrink: 0;
  }
}
</style>
