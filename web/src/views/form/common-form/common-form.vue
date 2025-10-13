<template>
  <div class="snow-page">
    <div class="snow-inner">
      <a-row :gutter="16">
        <a-col :xs="24" :sm="24" :md="12" :lg="12" :xl="12" :xxl="12">
          <a-form ref="formRef" auto-label-width :size="form.size" :model="form" @submit="handleSubmit">
            <a-form-item field="size" label="表单大小">
              <a-radio-group v-model="form.size" type="button">
                <a-radio value="mini">迷你</a-radio>
                <a-radio value="small">偏小</a-radio>
                <a-radio value="medium">中等</a-radio>
                <a-radio value="large">偏大</a-radio>
              </a-radio-group>
            </a-form-item>
            <a-form-item
              field="name"
              label="用户名称"
              :rules="[
                { required: true, message: '用户名称不能为空' },
                { minLength: 5, message: '必须大于 5 个字符' }
              ]"
              :validate-trigger="['change', 'input']"
            >
              <a-input :style="{ width: '100%' }" v-model="form.name" placeholder="请输入用户名称" allow-clear />
            </a-form-item>
            <a-form-item
              field="age"
              label="年龄"
              :rules="[
                { required: true, message: '年龄为必填项' },
                { type: 'number', max: 200, message: '年龄上限为 200 岁' }
              ]"
            >
              <a-input-number v-model="form.age" placeholder="请输入您的年龄" allow-clear />
            </a-form-item>
            <a-form-item field="section" label="选项" :rules="[{ required: true, message: '选项不能为空' }]">
              <a-select v-model="form.section" placeholder="请选择" allow-clear>
                <a-option value="section one">Section One</a-option>
                <a-option value="section two">Section Two</a-option>
                <a-option value="section three">Section Three</a-option>
              </a-select>
            </a-form-item>
            <a-form-item field="province" label="地区" :rules="[{ required: true, message: '地区不能为空' }]">
              <a-cascader v-model="form.province" :options="options" placeholder="请选择地区" allow-clear />
            </a-form-item>
            <a-form-item field="options" label="爱好" :rules="[{ type: 'array', minLength: 2, message: '请选择' }]">
              <a-checkbox-group v-model="form.options">
                <a-checkbox value="option one">听歌</a-checkbox>
                <a-checkbox value="option two">看电影</a-checkbox>
                <a-checkbox value="option three">旅游</a-checkbox>
                <a-checkbox value="option four">电竞</a-checkbox>
              </a-checkbox-group>
            </a-form-item>
            <a-form-item field="date" label="日期">
              <a-date-picker v-model="form.date" placeholder="请选择日期" />
            </a-form-item>
            <a-form-item field="time" label="时间">
              <a-time-picker v-model="form.time" placeholder="请选择时间" />
            </a-form-item>
            <a-form-item field="radio" label="性别" :rules="[{ required: true, message: '性别不能为空' }]">
              <a-radio-group v-model="form.radio">
                <a-radio value="man">男</a-radio>
                <a-radio value="woman">女</a-radio>
              </a-radio-group>
            </a-form-item>
            <a-form-item field="slider" label="进度" :rules="[{ type: 'number', min: 5, message: '最小为5' }]">
              <a-slider v-model="form.slider" :max="10" />
            </a-form-item>
            <a-form-item field="score" label="评分">
              <a-rate v-model="form.score" allow-clear />
            </a-form-item>
            <a-form-item field="switch" label="状态" :rules="[{ type: 'boolean', true: true, message: '请打开开关' }]">
              <a-switch v-model="form.switch" />
            </a-form-item>
            <a-form-item field="multiSelect" label="多选">
              <a-select v-model="form.multiSelect" placeholder="请选择" multiple>
                <a-option value="section one">多选一</a-option>
                <a-option value="section two">多选二</a-option>
                <a-option value="section three">多选三</a-option>
              </a-select>
            </a-form-item>
            <a-form-item field="treeSelect" label="分组">
              <a-tree-select :data="treeData" v-model="form.treeSelect" placeholder="请选择分组" />
            </a-form-item>
            <a-form-item>
              <a-space>
                <a-button html-type="submit">提交</a-button>
                <a-button @click="onReset">重置</a-button>
              </a-space>
            </a-form-item>
          </a-form>
        </a-col>
        <a-col :xs="24" :sm="24" :md="12" :lg="12" :xl="12" :xxl="12">
          <s-code-view :code-json="codeJson" />
        </a-col>
      </a-row>
    </div>
  </div>
</template>

<script setup lang="ts">
const handleSubmit = ({ values, errors }: any) => {
  console.log("values:", values, "\nerrors:", errors);
};
// 将序列化的字符串美观输出，\t 换行
const codeJson = computed(() => JSON.stringify(form, null, "\t"));
const form = reactive({
  size: "medium",
  name: "",
  age: undefined,
  section: "",
  province: "guanggu",
  options: [],
  date: "",
  time: "",
  radio: "man",
  slider: 5,
  score: 5,
  switch: false,
  multiSelect: ["section one"],
  treeSelect: "node1"
});
const options = [
  {
    value: "hubei",
    label: "湖北",
    children: [
      {
        value: "wuhan",
        label: "武汉",
        children: [
          {
            value: "guanggu",
            label: "光谷"
          }
        ]
      },
      {
        value: "xiaogan",
        label: "孝感"
      },
      {
        value: "huanggang",
        label: "黄冈"
      },
      {
        value: "xianning",
        label: "咸宁"
      }
    ]
  },
  {
    value: "hunan",
    label: "湖南",
    children: [
      {
        value: "changsha",
        label: "长沙",
        children: [
          {
            value: "yuelu",
            label: "岳麓"
          }
        ]
      }
    ]
  }
];
const treeData = [
  {
    key: "node1",
    title: "分组01",
    children: [
      {
        key: "node1-01",
        title: "分组01-01"
      }
    ]
  },
  {
    key: "node2",
    title: "分组02",
    children: [
      {
        key: "node2-01",
        title: "分组02-01"
      },
      {
        key: "node2-01",
        title: "分组02-02"
      }
    ]
  }
];

const formRef = ref();
const onReset = () => {
  formRef.value.resetFields();
};
</script>

<style lang="scss" scoped></style>
