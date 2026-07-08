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

          <a-form-item field="payment_checkout" label="前台收银模板">
            <template #extra>
              <span v-html="currentCheckoutInfo"></span>
            </template>
            <a-select
              v-model="form.payment_checkout"
              placeholder="请选择收银台模板"
              :fallback-option="false"
              :loading="checkoutListLoading"
            >
              <a-option v-for="option in checkoutList" :key="option.value" :value="option.value">
                {{ option.label }}
              </a-option>
            </a-select>
          </a-form-item>

          <a-form-item field="payment_support_url" label="前台收银客服" extra="收银台页面跳转的客服链接地址，留空则不启用">
            <a-input v-model="form.payment_support_url" placeholder="http(s)://your-support-url" allow-clear />
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
import { setsConfAPI, resetApiAuthToken, checkoutListAPI } from "@/api/modules/conf/index";

const emit = defineEmits(["refresh"]);
const data = defineModel() as any;
const { isMobile } = useDevicesSize();
const layoutMode = computed(() => (isMobile.value ? "vertical" : "horizontal"));

const form = ref({
  api_auth_token: "",
  api_app_uri: "",
  payment_checkout: "",
  payment_support_url: "",
});
const rules = {};
const checkoutList = ref<Array<{ label: string; value: string; author: string; desc: string; link: string }>>([]);
const checkoutListLoading = ref(false);

const normalizePaymentCheckout = (value?: string) => {
  const rawValue = value || "";
  const validValues = checkoutList.value.map(item => item.value);
  if (validValues.length === 0) {
    return rawValue;
  }
  if (validValues.includes(rawValue)) {
    return rawValue;
  }
  if (validValues.includes("official")) {
    return "official";
  }
  return validValues[0] || "";
};

const syncFormFromConfig = () => {
  if (!data.value) return;

  form.value.api_auth_token = data.value.api_auth_token || "";
  form.value.api_app_uri = data.value.api_app_uri || "";
  form.value.payment_checkout = normalizePaymentCheckout(data.value.payment_checkout || data.value.payment_template);
  form.value.payment_support_url = data.value.payment_support_url || "";
};

// 获取收银台模板列表
const fetchCheckoutList = async () => {
  try {
    checkoutListLoading.value = true;
    const res = await checkoutListAPI({});
    if (res.data && typeof res.data === "object") {
      checkoutList.value = Object.entries(res.data).map(([key, template]: [string, any]) => ({
        label: template.name,
        value: key,
        author: template.author,
        desc: template.desc,
        link: template.link
      }));
      syncFormFromConfig();
    }
  } catch (error) {
    Message.error("获取收银台模板列表失败");
  } finally {
    checkoutListLoading.value = false;
  }
};

// 获取当前选中模板的详细信息
const currentCheckoutInfo = computed(() => {
  const current = checkoutList.value.find(item => item.value === form.value.payment_checkout);
  if (!current) return "选择前台收银台模板";

  let info = `作者: ${current.author}，` + current.desc;

  if (current.link !== "") {
    info += ` <a href="${current.link}" target="_blank">#Link</a>`;
  }
  return info || "选择收银台模板样式";
});

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

  form.value.payment_checkout = normalizePaymentCheckout(form.value.payment_checkout);

  await setsConfAPI([
    {
      key: "api_app_uri",
      value: form.value.api_app_uri
    },
    {
      key: "payment_checkout",
      value: form.value.payment_checkout
    },
    {
      key: "payment_support_url",
      value: form.value.payment_support_url
    }
  ]);

  Message.success("保存成功");

  emit("refresh");
};

watch(
  () => data.value,
  syncFormFromConfig,
  { immediate: true }
);

// 组件挂载时获取模板列表
onMounted(() => {
  fetchCheckoutList();
});
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
