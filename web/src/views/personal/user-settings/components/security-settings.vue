<template>
  <div>
    <a-row align="center" :gutter="[0, 16]">
      <a-col :span="24">
        <a-card title="安全设置">
          <a-form :model="form" :rules="rules" :style="{ width: '600px' }" @submit="onSubmit">
            <a-form-item field="password" label="登录密码">
              <a-input-password v-model="form.password" placeholder="请输入登录密码" allow-clear />
            </a-form-item>
            <a-form-item field="phone" label="安全手机">
              <a-input v-model="form.phone" placeholder="请输入安全手机" allow-clear>
                <template #prepend> +86 </template>
              </a-input>
            </a-form-item>
            <a-form-item field="email" label="安全邮箱">
              <a-input v-model="form.email" placeholder="请输入安全邮箱" allow-clear />
            </a-form-item>
            <a-form-item field="question" label="密保问题">
              <a-input v-model="form.question" placeholder="请输入密保问题" allow-clear />
            </a-form-item>
            <a-form-item field="answer" label="密保答案">
              <a-input v-model="form.answer" placeholder="请输入密保答案" allow-clear />
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
import useGlobalProperties from "@/hooks/useGlobalProperties";
const proxy = useGlobalProperties();
const form = ref({
  password: "123456",
  phone: "15888888888",
  email: "",
  question: "",
  answer: ""
});
const rules = {
  password: [
    {
      required: true,
      message: "登录密码不能为空"
    }
  ],
  phone: [
    {
      required: true,
      message: "安全手机不能为空"
    }
  ]
};

const onSubmit = ({ errors }: ArcoDesign.ArcoSubmit) => {
  if (errors) return;
  proxy.$message.success("修改安全设置");
};
</script>

<style lang="scss" scoped>
.row-title {
  font-size: $font-size-title-1;
}
</style>
