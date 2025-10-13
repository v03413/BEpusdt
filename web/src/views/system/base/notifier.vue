<template>
  <div>
    <a-row align="center" :gutter="[0, 16]">
      <a-col :span="24">
        <a-card title="通知设置">
          <a-form :model="form" :rules="rules" :style="{ width: '600px' }" @submit="onSubmit">
            <a-form-item field="notifier_channel" label="通知渠道">
              <a-select v-model="form.notifier_channel" placeholder="请选择通知渠道" @change="onChannelChange">
                <a-option
                  v-for="channel in channelConfigs"
                  :key="channel.value"
                  :value="channel.value"
                  :disabled="channel.disabled"
                >
                  {{ channel.label }}
                </a-option>
              </a-select>
            </a-form-item>

            <template v-for="field in currentChannelFields" :key="field.key">
              <a-form-item :field="field.key" :label="field.label">
                <a-input
                  v-model="form.notifier_params[field.key]"
                  :placeholder="field.placeholder"
                  :type="field.type || 'text'"
                  allow-clear
                />
              </a-form-item>
            </template>

            <a-form-item>
              <a-space>
                <a-button type="primary" html-type="submit">保存配置</a-button>
                <a-button v-if="form.notifier_channel !== 'none'" type="outline" @click="onTest" :loading="testLoading">
                  推送测试
                </a-button>
              </a-space>
            </a-form-item>
          </a-form>
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script setup lang="ts">
import { Message } from "@arco-design/web-vue";
import { notifierAPI, notifierTestAPI } from "@/api/modules/conf/index";

const emit = defineEmits(["refresh"]);
const data = defineModel() as any;

interface FieldConfig {
  key: string;
  label: string;
  placeholder: string;
  type?: string;
  required: boolean;
  message?: string;
  validator?: string;
}

interface ChannelConfig {
  value: string;
  label: string;
  disabled: boolean;
  fields: FieldConfig[];
}

interface FormData {
  notifier_channel: string;
  notifier_params: Record<string, string>;
}

const channelConfigs: ChannelConfig[] = [
  {
    value: "none",
    label: "关闭通知",
    disabled: false,
    fields: []
  },
  {
    value: "telegram",
    label: "Telegram",
    disabled: false,
    fields: [
      {
        key: "bot_token",
        label: "Bot Token",
        placeholder: "请输入 Telegram Bot Token",
        required: true,
        message: "Bot Token 不能为空"
      },
      { key: "chat_id", label: "Chat ID", placeholder: "请输入 Telegram Chat ID", required: true, message: "Chat ID 不能为空" },
      { key: "topic_id", label: "Topic ID", placeholder: "请输入 Telegram Topic ID", required: false }
    ]
  },
  {
    value: "wechat",
    label: "企业微信（开发中）",
    disabled: true,
    fields: [
      {
        key: "webhook_url",
        label: "Webhook URL",
        placeholder: "请输入企业微信 Webhook URL",
        type: "url",
        required: true,
        message: "Webhook URL不能为空",
        validator: "url"
      }
    ]
  },
  {
    value: "email",
    label: "邮箱（开发中）",
    disabled: true,
    fields: [
      {
        key: "email",
        label: "邮箱地址",
        placeholder: "请输入邮箱地址",
        type: "email",
        required: true,
        message: "邮箱地址不能为空",
        validator: "email"
      },
      {
        key: "smtp_server",
        label: "SMTP服务器",
        placeholder: "请输入SMTP服务器地址",
        required: true,
        message: "SMTP服务器不能为空"
      }
    ]
  }
];

const form = ref<FormData>({
  notifier_channel: "telegram",
  notifier_params: {}
});

const testLoading = ref<boolean>(false);

const currentChannelFields = computed<FieldConfig[]>(
  () => channelConfigs.find(config => config.value === form.value.notifier_channel)?.fields || []
);

const currentChannelParamKeys = computed<string[]>(() => currentChannelFields.value.map(field => field.key));

const rules = computed(() => {
  const baseRules: Record<string, any[]> = {
    notifier_channel: [{ required: true, message: "请选择通知渠道" }]
  };

  currentChannelFields.value.forEach(field => {
    if (field.required) {
      const fieldPath = `notifier_params.${field.key}`;
      const fieldRules: any[] = [{ required: true, message: field.message }];

      if (field.validator === "email") {
        fieldRules.push({ type: "email", message: "请输入正确的邮箱格式" });
      } else if (field.validator === "url") {
        fieldRules.push({ type: "url", message: "请输入正确的URL格式" });
      }

      baseRules[fieldPath] = fieldRules;
    }
  });

  return baseRules;
});

const initParams = (): Record<string, string> => {
  const params: Record<string, string> = {};
  channelConfigs.forEach(config => {
    config.fields.forEach(field => {
      params[field.key] = "";
    });
  });
  return params;
};

const onChannelChange = (): void => {
  form.value.notifier_params = initParams();
};

const onSubmit = async ({ errors }: ArcoDesign.ArcoSubmit): Promise<void> => {
  if (errors) return;

  try {
    const filteredParams: Record<string, string> = {};
    currentChannelParamKeys.value.forEach(key => {
      const value = form.value.notifier_params[key];
      if (value !== undefined && value !== null) {
        filteredParams[key] = String(value);
      }
    });

    const response = await notifierAPI({
      channel: form.value.notifier_channel,
      params: filteredParams
    });

    if (response?.code === 200) {
      Message.success("配置成功！");
      emit("refresh");
    } else {
      Message.error(response?.msg || "配置保存失败");
    }
  } catch (error: any) {
    console.error("配置保存失败:", error);
    Message.error("配置保存失败，请稍后重试");
  }
};

const onTest = async (): Promise<void> => {
  try {
    testLoading.value = true;

    const filteredParams: Record<string, string> = {};
    currentChannelParamKeys.value.forEach(key => {
      const value = form.value.notifier_params[key];
      if (value !== undefined && value !== null) {
        filteredParams[key] = String(value);
      }
    });

    const response = await notifierTestAPI({
      channel: form.value.notifier_channel,
      params: filteredParams
    });

    if (response?.code === 200) {
      Message.success("推送测试成功！");
    } else {
      Message.error(response?.msg || "推送测试失败");
    }
  } catch (error: any) {
    console.error("推送测试失败:", error);
    Message.error("推送测试失败，请稍后重试");
  } finally {
    testLoading.value = false;
  }
};

watch(
  () => data.value,
  () => {
    if (data.value) {
      form.value.notifier_channel = String(data.value.notifier_channel || "telegram");

      if (data.value.notifier_params) {
        try {
          const params =
            typeof data.value.notifier_params === "string" ? JSON.parse(data.value.notifier_params) : data.value.notifier_params;

          const parsedParams: Record<string, string> = {};
          Object.keys(params).forEach(key => {
            parsedParams[key] = String(params[key] || "");
          });

          form.value.notifier_params = { ...initParams(), ...parsedParams };
        } catch (e) {
          console.error("解析 notifier_params 失败:", e);
          form.value.notifier_params = initParams();
        }
      } else {
        form.value.notifier_params = initParams();
      }
    }
  },
  { immediate: true }
);

onMounted(() => {
  form.value.notifier_params = initParams();
});
</script>

<style lang="scss" scoped>
.row-title {
  font-size: $font-size-title-1;
}
</style>
