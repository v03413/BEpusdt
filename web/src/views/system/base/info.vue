<template>
  <a-row align="center" :gutter="[0, 16]">
    <a-col :span="24">
      <a-card title="基本信息">
        <a-form :model="form" :rules="rules" :style="{ width: '600px' }" @submit="onSubmit">
          <a-form-item
            field="block_height_max_diff"
            label="区块最大差值"
            extra="区块高度最大差值，超过此值则以当前区块高度为准，重新开始扫描"
          >
            <a-input v-model="form.block_height_max_diff" placeholder="推荐 1000" />
          </a-form-item>

          <a-form-item
            field="block_offset_confirm"
            label="区块偏移确认"
            extra="启用后可进一步提高交易的安全性，但会增加交易确认的等待时间"
          >
            <a-select v-model="form.block_offset_confirm" placeholder="请选择">
              <a-option value="0">关闭</a-option>
              <a-option value="1">开启</a-option>
            </a-select>
          </a-form-item>

          <a-form-item
            field="notify_max_retry"
            label="回调最大重试"
            extra="支付回调失败时的最大重试次数，重试分钟间隔数：2 4 8 16 32 64 ..."
          >
            <a-input v-model="form.notify_max_retry" placeholder="推荐 10" />
          </a-form-item>

          <a-form-item
            field="payment_min_amount"
            label="单笔最小金额"
            extra="单笔支付允许的最小金额，单位为创建交易时传入的法币，用于一定风险控制"
          >
            <a-input v-model="form.payment_min_amount" placeholder="推荐 0.01" />
          </a-form-item>

          <a-form-item
            field="payment_max_amount"
            label="单笔最大金额"
            extra="单笔支付允许的最大金额，单位为创建交易时传入的法币，用于一定风险控制"
          >
            <a-input v-model="form.payment_max_amount" placeholder="建议 9999" />
          </a-form-item>

          <a-form-item
            field="payment_timeout"
            label="订单默认超时"
            extra="订单默认超时，单位为秒；超过此时间未支付的订单将被自动关闭"
          >
            <a-input v-model="form.payment_timeout" placeholder="推荐 1200" />
          </a-form-item>

          <a-form-item field="payment_match_mode" label="金额匹配模式">
            <template #extra>
              订单交易在金额确认时，使用不同算法的算法进行比对；详细区别请看
              <a href="https://github.com/v03413/BEpusdt/blob/main/docs/payment-match-mode/README.md" target="_blank">文档说明</a>
            </template>
            <a-select v-model="form.payment_match_mode" placeholder="请选择金额匹配模式">
              <a-option value="classic">传统模式</a-option>
              <a-option value="has_prefix">前缀匹配</a-option>
              <a-option value="round_off">修约匹配</a-option>
            </a-select>
          </a-form-item>

          <a-form-item>
            <a-button type="primary" html-type="submit">提交</a-button>
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
  block_offset_confirm: "0",
  notify_max_retry: "",
  payment_max_amount: "",
  payment_min_amount: "",
  payment_match_mode: "classic"
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
  block_offset_confirm: [
    {
      required: true,
      message: "区块偏移确认不能为空"
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
  payment_timeout: [
    {
      required: true,
      type: "number",
      min: 180,
      max: 3600,
      message: "订单默认超时必须在180到3600秒之间",
      positive: true
    }
  ],
  payment_match_mode: [
    {
      required: true,
      message: "金额匹配模式不能为空"
    }
  ]
};

const onSubmit = async ({ errors }: ArcoDesign.ArcoSubmit) => {
  if (errors) return;

  await setsConfAPI([
    { key: "block_height_max_diff", value: form.value.block_height_max_diff },
    { key: "block_offset_confirm", value: form.value.block_offset_confirm },
    { key: "notify_max_retry", value: form.value.notify_max_retry },
    { key: "payment_max_amount", value: form.value.payment_max_amount },
    { key: "payment_min_amount", value: form.value.payment_min_amount },
    { key: "payment_timeout", value: form.value.payment_timeout },
    { key: "payment_match_mode", value: form.value.payment_match_mode }
  ]);

  Message.success("保存成功");

  emit("refresh");
};

watch(
  () => data.value,
  () => {
    form.value.block_height_max_diff = data.value.block_height_max_diff;
    form.value.block_offset_confirm = data.value.block_offset_confirm || "0";
    form.value.notify_max_retry = data.value.notify_max_retry;
    form.value.payment_max_amount = data.value.payment_max_amount;
    form.value.payment_min_amount = data.value.payment_min_amount;
    form.value.payment_timeout = data.value.payment_timeout;
    form.value.payment_match_mode = data.value.payment_match_mode || "classic";
  }
);
</script>

<style lang="scss" scoped>
.row-title {
  font-size: $font-size-title-1;
}
</style>
