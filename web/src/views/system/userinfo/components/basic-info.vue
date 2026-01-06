<template>
  <a-row align="center" :gutter="[0, 16]">
    <a-col :span="24">
      <a-card title="基本信息">
        <a-form :model="form" :rules="rules" :style="{ width: '600px' }" @submit="onSubmit">
          <a-form-item field="userName" label="用户名">
            <a-input v-model="form.userName" placeholder="请输入用户名" allow-clear />
          </a-form-item>
          <a-form-item field="nickName" label="用户昵称">
            <a-input v-model="form.nickName" placeholder="请输入用户昵称" allow-clear />
          </a-form-item>
          <a-form-item field="sex" label="性别" validate-trigger="blur">
            <a-radio-group v-model="form.sex" :options="sexOption">
              <template #label="{ data }">
                <div>{{ data.name }}</div>
              </template>
            </a-radio-group>
          </a-form-item>
          <a-form-item field="description" label="描述">
            <a-textarea placeholder="请输入描述" v-model="form.description" allow-clear />
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
import useGlobalProperties from "@/hooks/useGlobalProperties";

const emit = defineEmits(["refresh"]);
const proxy = useGlobalProperties();
const data = defineModel() as any;
const sexOption = ref(dictFilter("gender"));
const form = ref({
  id: "",
  userName: "",
  nickName: "",
  sex: null,
  description: ""
});
const rules = {
  userName: [
    {
      required: true,
      message: "用户名不能为空"
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
    form.value.userName = data.value.userName;
    form.value.nickName = data.value.nickName;
    form.value.sex = data.value.sex;
    form.value.description = data.value.description;
  }
);
</script>

<style lang="scss" scoped>
.row-title {
  font-size: $font-size-title-1;
}
</style>
