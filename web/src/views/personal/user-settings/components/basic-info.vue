<template>
  <a-row align="center" :gutter="[0, 16]">
    <a-col :span="24">
      <a-card title="基本信息">
        <a-form :model="form" :rules="rules" :style="{ width: '600px' }" @submit="onSubmit">
          <a-form-item field="userId" label="用户ID">
            <a-input v-model="form.userId" disabled />
          </a-form-item>
          <a-form-item field="username" label="用户名">
            <a-input v-model="form.username" placeholder="请输入用户名" allow-clear />
          </a-form-item>
          <a-form-item field="nickname" label="用户昵称">
            <a-input v-model="form.nickname" placeholder="请输入用户昵称" allow-clear />
          </a-form-item>
          <a-form-item field="role" label="角色">
            <a-select :style="{ width: '320px' }" v-model="form.role" placeholder="请选择角色" allow-clear>
              <a-option>系统管理员</a-option>
              <a-option>普通角色</a-option>
            </a-select>
          </a-form-item>
          <a-form-item field="address" label="所在区域">
            <a-cascader
              v-model="form.address"
              :options="options"
              expand-trigger="hover"
              :style="{ width: '320px' }"
              placeholder="请选择所在区域"
              allow-clear
            />
          </a-form-item>
          <a-form-item field="introduce" label="简介">
            <a-textarea placeholder="请输入简介" v-model="form.introduce" allow-clear />
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
const proxy = useGlobalProperties();
const form = ref({
  userId: "1000001",
  username: "admin",
  nickname: "兔子先森",
  role: "系统管理员",
  address: "guanggu",
  introduce: ""
});
const rules = {
  userId: [
    {
      required: true,
      message: "用户ID不能为空"
    }
  ],
  username: [
    {
      required: true,
      message: "用户名不能为空"
    }
  ]
};

const options = ref([
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
        value: "jingzhou",
        label: "荆州"
      }
    ]
  },
  {
    value: "hunan",
    label: "湖南",
    children: [
      {
        value: "changsha",
        label: "长沙"
      }
    ]
  }
]);
const onSubmit = ({ errors }: ArcoDesign.ArcoSubmit) => {
  if (errors) return;
  proxy.$message.success("修改基本信息");
};
</script>

<style lang="scss" scoped>
.row-title {
  font-size: $font-size-title-1;
}
</style>
