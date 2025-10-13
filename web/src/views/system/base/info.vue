<template>
  <a-row align="center" :gutter="[0, 16]">
    <a-col :span="24">
      <a-card title="基本信息">
        <a-form :model="form" :rules="rules" :style="{ width: '600px' }" @submit="onSubmit">
          <a-form-item
            field="block_height_max_diff"
            label="区块最大差值"
            :label-col-props="{ span: 6 }"
            :wrapper-col-props="{ span: 18 }"
          >
            <a-row :gutter="16" align="center">
              <a-col :span="8">
                <a-input v-model="form.block_height_max_diff" placeholder="推荐 1000" />
              </a-col>
              <a-col :span="16">
                <div class="field-description">区块高度最大差值，超过此值则以当前区块高度为准，重新开始扫描</div>
              </a-col>
            </a-row>
          </a-form-item>

          <a-form-item
            field="notify_max_retry"
            label="回调最大重试"
            :label-col-props="{ span: 6 }"
            :wrapper-col-props="{ span: 18 }"
          >
            <a-row :gutter="16" align="center">
              <a-col :span="8">
                <a-input v-model="form.notify_max_retry" placeholder="推荐 10" />
              </a-col>
              <a-col :span="16">
                <div class="field-description">支付回调失败时的最大重试次数，重试分钟间隔数：2 4 8 16 32 64 ...</div>
              </a-col>
            </a-row>
          </a-form-item>

          <a-form-item
            field="monitor_min_amount"
            label="监控最小金额"
            :label-col-props="{ span: 6 }"
            :wrapper-col-props="{ span: 18 }"
          >
            <a-row :gutter="16" align="center">
              <a-col :span="8">
                <a-input v-model="form.monitor_min_amount" placeholder="推荐 0.01" />
              </a-col>
              <a-col :span="16">
                <div class="field-description">低于此金额的非订单转账交易不进行通知，可用于防范诱导式诈骗交易</div>
              </a-col>
            </a-row>
          </a-form-item>

          <a-form-item
            field="payment_min_amount"
            label="单笔支付金额"
            :label-col-props="{ span: 6 }"
            :wrapper-col-props="{ span: 18 }"
          >
            <a-row :gutter="16" align="center">
              <a-col :span="8">
                <a-input v-model="form.payment_min_amount" placeholder="推荐 0.01" />
              </a-col>
              <a-col :span="16">
                <div class="field-description">单笔支付允许的最小金额，单位为创建交易时传入的法币，用于一定风险控制</div>
              </a-col>
            </a-row>
          </a-form-item>

          <a-form-item
            field="payment_max_amount"
            label="单笔最大金额"
            :label-col-props="{ span: 6 }"
            :wrapper-col-props="{ span: 18 }"
          >
            <a-row :gutter="16" align="center">
              <a-col :span="8">
                <a-input v-model="form.payment_max_amount" placeholder="建议 9999" />
              </a-col>
              <a-col :span="16">
                <div class="field-description">单笔支付允许的最大金额，单位为创建交易时传入的法币，用于一定风险控制</div>
              </a-col>
            </a-row>
          </a-form-item>

          <a-form-item
            field="payment_timeout"
            label="订单默认超时"
            :label-col-props="{ span: 6 }"
            :wrapper-col-props="{ span: 18 }"
          >
            <a-row :gutter="16" align="center">
              <a-col :span="8">
                <a-input v-model="form.payment_timeout" placeholder="推荐 1200" />
              </a-col>
              <a-col :span="16">
                <div class="field-description">订单默认超时，单位为秒；超过此时间未支付的订单将被自动关闭</div>
              </a-col>
            </a-row>
          </a-form-item>

          <a-form-item :wrapper-col-props="{ span: 18, offset: 6 }">
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
import { Message } from "@arco-design/web-vue";
import { setsConfAPI } from "@/api/modules/conf/index";
const emit = defineEmits(["refresh"]);
const data = defineModel() as any;
const form = ref({
  payment_timeout: "",
  block_height_max_diff: "",
  notify_max_retry: "",
  payment_max_amount: "",
  payment_min_amount: "",
  monitor_min_amount: ""
});
const rules = {
  block_height_max_diff: [
    {
      required: true,
      type: "number",
      positive: true,
      message: "区块高度最大差值不能为空"
    }
  ],
  notify_max_retry: [
    {
      required: true,
      type: "number",
      positive: true,
      message: "回调最大重试次数不能为空"
    }
  ],
  payment_min_amount: [
    {
      required: true,
      type: "number",
      positive: true,
      message: "单笔支付最小金额不能为空"
    }
  ],
  payment_max_amount: [
    {
      required: true,
      type: "number",
      positive: true,
      message: "单笔支付最大金额不能为空"
    }
  ],
  monitor_min_amount: [
    {
      required: true,
      type: "number",
      positive: true,
      message: "监控最小金额不能为空"
    }
  ],
  payment_timeout: [
    {
      required: true,
      type: "number",
      min: 180,
      max: 3600,
      message: "订单默认超时必须在180到3600秒之间",
      positive: true
    }
  ]
};

const onSubmit = async ({ errors }: ArcoDesign.ArcoSubmit) => {
  if (errors) return;

  await setsConfAPI([
    { key: "block_height_max_diff", value: form.value.block_height_max_diff },
    { key: "notify_max_retry", value: form.value.notify_max_retry },
    { key: "payment_max_amount", value: form.value.payment_max_amount },
    { key: "payment_min_amount", value: form.value.payment_min_amount },
    { key: "monitor_min_amount", value: form.value.monitor_min_amount },
    { key: "payment_timeout", value: form.value.payment_timeout }
  ]);

  Message.success("保存成功");

  emit("refresh");
};

watch(
  () => data.value,
  () => {
    form.value.block_height_max_diff = data.value.block_height_max_diff;
    form.value.notify_max_retry = data.value.notify_max_retry;
    form.value.payment_max_amount = data.value.payment_max_amount;
    form.value.payment_min_amount = data.value.payment_min_amount;
    form.value.monitor_min_amount = data.value.monitor_min_amount;
    form.value.payment_timeout = data.value.payment_timeout;
  }
);
</script>

<style lang="scss" scoped>
.row-title {
  font-size: $font-size-title-1;
}

.field-description {
  color: #86909c;
  font-size: 12px;
  line-height: 1.4;
}
</style>
