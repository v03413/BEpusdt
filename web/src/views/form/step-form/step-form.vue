<template>
  <div class="snow-page">
    <a-spin :loading="loading" style="display: block">
      <div class="snow-inner container">
        <a-row justify="center">
          <a-col :xs="22" :sm="18" :md="16" :lg="16" :xl="12" :xxl="12">
            <a-steps :current="currentStep" line-less>
              <a-step description="创建渠道活动">基本信息</a-step>
              <a-step description="完善活动渠道信息">活动渠道</a-step>
              <a-step description="创建成功">完成创建</a-step>
            </a-steps>
          </a-col>
        </a-row>
        <a-row justify="center" class="margin-top">
          <a-col :xs="18" :sm="12" :md="12" :lg="12" :xl="8" :xxl="8">
            <a-form ref="formRef" auto-label-width :model="form" :rules="rules" @submit="handleSubmit">
              <div v-if="currentStep == 1">
                <a-form-item field="name" label="活动名称">
                  <a-input v-model="form.name" placeholder="请输入活动名称" />
                </a-form-item>
                <a-form-item field="type" label="渠道类型">
                  <a-select v-model="form.type" placeholder="请选择渠道类型" allow-clear>
                    <a-option value="1">互联网</a-option>
                    <a-option value="2">线下</a-option>
                  </a-select>
                </a-form-item>
                <a-form-item field="date" label="推广时间">
                  <a-range-picker v-model="form.date" style="width: 100%" />
                </a-form-item>
                <a-form-item field="address" label="推广地址">
                  <a-input v-model="form.address" placeholder="请输入推广地址" />
                  <template #extra>
                    <div>从哪些渠道了解到该项目，例如gitee：https://gitee.com/wang_fan_w/SnowAdmin</div>
                  </template>
                </a-form-item>
              </div>
              <div v-if="currentStep == 2">
                <a-form-item field="source" label="广告来源">
                  <a-input v-model="form.source" placeholder="请输入广告来源" />
                </a-form-item>
                <a-form-item field="medium" label="广告媒介">
                  <a-input v-model="form.medium" placeholder="请输入广告媒介" />
                </a-form-item>
                <a-form-item field="keyword" label="关键词">
                  <a-select v-model="form.keyword" placeholder="请选择关键词" allow-clear multiple>
                    <a-option value="1">gitee</a-option>
                    <a-option value="2">github</a-option>
                  </a-select>
                </a-form-item>
                <a-form-item field="tip" label="推送提醒">
                  <a-switch v-model="form.tip" />
                </a-form-item>
                <a-form-item field="content" label="广告内容">
                  <a-textarea v-model="form.content" placeholder="请输入广告内容" allow-clear />
                </a-form-item>
              </div>

              <div v-if="currentStep == 3">
                <a-result status="success" title="提交成功">
                  <template #subtitle> 提交表单成功 </template>
                  <template #extra>
                    <a-space>
                      <a-button type="primary">查看详情</a-button>
                      <a-button @click="currentStep = 1">再次创建</a-button>
                    </a-space>
                  </template>
                </a-result>
              </div>

              <a-form-item v-if="currentStep != 3">
                <a-space>
                  <a-button @click="onLastStep" v-if="currentStep != 1">上一步</a-button>
                  <a-button html-type="submit" type="primary">下一步</a-button>
                </a-space>
              </a-form-item>
            </a-form>
          </a-col>
        </a-row>
        <a-row v-if="currentStep == 3">
          <a-col :span="16" :offset="4">
            <a-typography style="padding: 24px; background: var(--color-fill-2)">
              <a-typography-paragraph>提示</a-typography-paragraph>
              <ul>
                <li>您可以在github或者gitee上完整预览该项目</li>
                <li>如果您觉得该项目对您有用，可以动动小手，点亮一下star噢</li>
              </ul>
            </a-typography>
          </a-col>
        </a-row>
      </div>
    </a-spin>
  </div>
</template>

<script setup lang="ts">
const loading = ref(false);
const currentStep = ref(1);
const form = ref({
  name: "",
  type: null,
  date: [],
  address: "https://gitee.com/wang_fan_w/SnowAdmin",
  source: "",
  medium: "",
  keyword: [],
  tip: true,
  content: ""
});
const rules = ref({
  name: [
    {
      required: true,
      message: "请输入活动名称"
    }
  ],
  type: [
    {
      required: true,
      message: "请选择渠道类型"
    }
  ],
  date: [
    {
      required: true,
      message: "请选择推广时间"
    }
  ],
  address: [
    {
      required: true,
      message: "请输入推广地址"
    }
  ],
  source: [
    {
      required: true,
      message: "请输入广告来源"
    }
  ],
  medium: [
    {
      required: true,
      message: "请输入广告媒介"
    }
  ],
  keyword: [
    {
      required: true,
      message: "请选择关键词"
    }
  ],
  tip: [
    {
      required: true,
      message: "请选择推送提醒"
    }
  ],
  content: [
    {
      required: true,
      message: "请输入广告内容"
    }
  ]
});
const formRef = ref();

const handleSubmit = ({ errors }: ArcoDesign.ArcoSubmit) => {
  if (errors) return;
  if (currentStep.value == 3) return;
  if (currentStep.value == 2) {
    loading.value = true;
    setTimeout(() => {
      loading.value = false;
      currentStep.value += 1;
    }, 500);
  } else {
    currentStep.value += 1;
  }
};

const onLastStep = () => {
  if (currentStep.value == 1) return;
  currentStep.value -= 1;
};
</script>

<style lang="scss" scoped>
.container {
  padding: 60px 0;
}
.margin-top {
  margin-top: 60px;
}
</style>
