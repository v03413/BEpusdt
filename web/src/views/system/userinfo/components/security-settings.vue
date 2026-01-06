<template>
  <div>
    <a-row align="center" :gutter="[0, 16]">
      <a-col :span="24">
        <a-card title="安全设置">
          <a-form :model="form" :rules="rules" :style="{ width: '600px' }" @submit="onSubmit">
            <a-form-item field="password" label="登录密码">
              <a-input-password v-model="form.password" placeholder="请输入登录密码" allow-clear />
            </a-form-item>
            <a-form-item field="phone" label="手机号">
              <a-input v-model="form.phone" placeholder="请输入手机号" allow-clear>
                <template #prepend> +86 </template>
              </a-input>
            </a-form-item>
            <a-form-item field="email" label="邮箱">
              <a-input v-model="form.email" placeholder="请输入邮箱" allow-clear />
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

const emit = defineEmits(["refresh"]);
const proxy = useGlobalProperties();
const data = defineModel() as any;
const form = ref({
  id: "",
  password: "",
  phone: "",
  email: ""
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
      message: "手机号不能为空"
    }
  ]
};

const onSubmit = ({ errors }: ArcoDesign.ArcoSubmit) => {
  if (errors) return;
  proxy.$message.success("模拟修改成功");
  emit("refresh");
};

watch(
  () => data.value,
  () => {
    form.value.id = data.value.id;
    form.value.phone = data.value.phone;
    form.value.email = data.value.email;
  }
);
</script>

<style lang="scss" scoped>
.row-title {
  font-size: $font-size-title-1;
}
</style>
