<template>
  <div class="snow-page">
    <a-spin :loading="loading" style="display: block">
      <div class="snow-inner container">
        <a-row justify="center">
          <a-col :xs="22" :sm="18" :md="16" :lg="16" :xl="12" :xxl="12">
            <a-steps :current="currentStep" line-less>
              <a-step description="创建订单">基本信息</a-step>
              <a-step description="创建成功">完成创建</a-step>
            </a-steps>
          </a-col>
        </a-row>
        <a-row justify="center" class="margin-top">
          <a-col :xs="18" :sm="12" :md="12" :lg="12" :xl="8" :xxl="8">
            <a-form ref="formRef" auto-label-width :layout="formLayout" :model="form" :rules="rules" @submit="handleSubmit">
              <div v-if="currentStep == 1">
                <a-form-item
                  field="name"
                  label="收款项目"
                  :validate-trigger="['change', 'input']"
                >
                  <a-input :style="{ width: '100%' }" v-model="form.name" placeholder="请输入收款项目" allow-clear />
                </a-form-item>
                <a-form-item
                  field="order_id"
                  label="订单号"
                  :validate-trigger="['change', 'input']"
                >
                  <a-input :style="{ width: '100%' }" v-model="form.order_id" placeholder="请输入订单号" allow-clear />
                </a-form-item>
                <a-form-item
                  field="amount"
                  label="订单金额"
                  :validate-trigger="['change', 'input']"
                >
                  <a-input-number :style="{ width: '100%' }" v-model="form.amount" placeholder="请输入订单金额" allow-clear />
                </a-form-item>
                <a-form-item field="trade_fiat" label="法币币种" :rules="[{ required: true, message: '法币币种不能为空' }]">
                  <a-select v-model="form.trade_fiat" placeholder="请选择" allow-clear>
                    <a-option v-for="(_, key) in userInfoStore.trade_fiat" :key="key" :value="key">{{ key }}</a-option>
                  </a-select>
                </a-form-item>
                <a-form-item field="trade_crypto" label="限定加密货币（留空不限制）">
                  <a-select v-model="form.trade_crypto" placeholder="请选择" multiple>
                    <a-option v-for="(_, key) in userInfoStore.trade_crypto" :key="key" :value="key">{{ key }}</a-option>
                  </a-select>
                </a-form-item>
                <a-form-item field="timeout" label="订单有效期（小时）">
                  <a-slider v-model="form.timeout" :max="3" />
                </a-form-item>
              </div>

              <div v-if="currentStep == 2">
                <a-result :status="resultStatus" :title="resultTitle">
                  <template #subtitle> {{ resultSubtitle }} </template>
                  <template #extra>
                    <a-space direction="vertical" size="large">
                      <div v-if="resultStatus === 'success'">订单链接：<a :href="paymentUrl" class="link" target="_blank">{{ paymentUrl }}</a></div>
                      <a-space>
                        <a-button type="primary" v-if="resultStatus === 'success'" @click="copyLink">复制订单链接</a-button>
                        <a-button @click="resetForm">再次创建</a-button>
                      </a-space>
                    </a-space>
                  </template>
                </a-result>
              </div>
              <a-form-item v-if="currentStep != 2">
                <a-space>
                  <a-button @click="onLastStep" v-if="currentStep != 1">上一步</a-button>
                  <a-button html-type="submit" type="primary">下一步</a-button>
                </a-space>
              </a-form-item>
            </a-form>
          </a-col>
        </a-row>
        <a-row v-if="currentStep == 2">
          <a-col :span="16" :offset="4">
            <a-typography style="padding: 24px; background: var(--color-fill-2)">
              <a-typography-paragraph>提示</a-typography-paragraph>
              <ul>
                <li>您可以将该链接发给客户让其支付</li>
                <li>该订单将在您设置的有效期内有效</li>
              </ul>
            </a-typography>
          </a-col>
        </a-row>
      </div>
    </a-spin>
  </div>
</template>
<script setup lang="ts">
import { ref, reactive, onMounted } from "vue";
import { useUserInfoStore } from "@/store/modules/user-info";
import { createOrderApi } from "@/api/modules/order";
import { Message } from "@arco-design/web-vue";

import dayjs from "dayjs";

const copyLink = async () => {
  try {
    if (navigator.clipboard && window.isSecureContext) {
      await navigator.clipboard.writeText(paymentUrl.value);
      Message.success("复制成功");
    } else {
      // Fallback for non-secure contexts or older browsers
      const textArea = document.createElement("textarea");
      textArea.value = paymentUrl.value;
      textArea.style.position = "fixed";
      textArea.style.left = "-9999px";
      document.body.appendChild(textArea);
      textArea.focus();
      textArea.select();
      try {
        document.execCommand('copy');
        Message.success("复制成功");
      } catch (err) {
        console.error('Fallback: Oops, unable to copy', err);
        Message.error("复制失败");
      }
      document.body.removeChild(textArea);
    }
  } catch (e) {
    console.error(e);
    Message.error("复制失败");
  }
};

const loading = ref(false);
const currentStep = ref(1);
const userInfoStore = useUserInfoStore();

const formRef = ref();
const formLayout = ref("vertical");
const rules = {
  name: [
    { required: true, message: '收款项目不能为空' },
    { minLength: 3, message: '必须大于 3 个字符' }
  ],
  order_id: [{ required: true, message: '订单号不能为空' }],
  amount: [
    { required: true, message: '订单金额不能为空' },
    { type: 'number', min: 1, max: 99999999, message: '订单金额必须在 1 到 99999999 之间' }
  ],
  trade_fiat: [{ required: true, message: '法币币种不能为空' }],
  timeout: [{ type: 'number', min: 1, message: '最小为1小时' }]
};

const form = reactive({
  name: "",
  order_id: "",
  amount: 0,
  trade_fiat: undefined,
  trade_crypto: [],
  timeout: 1
});

const generateOrderId = () => {
  return "PAY" + dayjs().format("YYYYMMDDHHmmss");
};

const initForm = async () => {
  form.name = "";
  form.order_id = generateOrderId();

  form.amount = 0;
  form.trade_fiat = undefined;
  form.trade_crypto = [];
  form.timeout = 1;
};

onMounted(() => {
  initForm();
});

const resultStatus = ref("success");
const resultTitle = ref("创建成功");
const resultSubtitle = ref("订单创建成功");
const paymentUrl = ref("");

const handleSubmit = async ({ errors, values }: ArcoDesign.ArcoSubmit) => {
  if (errors) return;
  if (currentStep.value == 2) return;
  
  loading.value = true;
  try {
    const payload = {
      ...values,
      fiat: form.trade_fiat,
      currencies: Array.isArray(form.trade_crypto) ? form.trade_crypto.join(",") : "",
      timeout: form.timeout * 3600
    };

    const res = await createOrderApi(payload);
    if (res.code === 200) {
      resultStatus.value = "success";
      resultTitle.value = "创建成功";
      resultSubtitle.value = "订单创建成功";
      paymentUrl.value = res.data.payment_url;
      currentStep.value += 1;
    } else {
      resultStatus.value = "error";
      resultTitle.value = "订单创建失败";
      resultSubtitle.value = res.data.message || "未知错误";
      paymentUrl.value = "";
      currentStep.value += 1;
    }
  } catch (err: any) {
    console.error(err);
    resultStatus.value = "error";
    resultTitle.value = "订单创建失败";
    resultSubtitle.value = err.message || "请求失败";
     paymentUrl.value = "";
    currentStep.value += 1;
  } finally {
    loading.value = false;
  }
};

const onLastStep = () => {
  if (currentStep.value == 1) return;
  currentStep.value -= 1;
};

const resetForm = () => {
  currentStep.value = 1;
  initForm();
}
</script>

<style lang="scss" scoped>
.container {
  padding: 60px 0;
}
.margin-top {
  margin-top: 60px;
}

.link {
  color: #409eff;
  text-decoration: none;
}

.link:hover {
  color: #66b1ff;
}
</style>