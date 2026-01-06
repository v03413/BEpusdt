<template>
  <div class="snow-page">
    <div class="container">
      <a-card :bordered="false">
        <a-row align="center">
          <a-col :span="2">
            <div>
              <a-avatar :size="100" @click="toast" trigger-type="mask">
                <img alt="avatar" :src="myImage" />
                <template #trigger-icon>
                  <IconEdit />
                </template>
              </a-avatar>
            </div>
          </a-col>
          <a-col :span="22">
            <a-space direction="vertical" size="large">
              <a-descriptions :data="data" :column="3" title="用户资料" :align="{ label: 'right' }" />
            </a-space>
          </a-col>
        </a-row>
      </a-card>
      <a-card class="margin-top" :bordered="false">
        <a-row align="center">
          <a-col :span="24">
            <a-tabs :type="type" :size="size">
              <a-tab-pane key="1" title="基本信息">
                <BasicInfo />
              </a-tab-pane>
              <a-tab-pane key="2" title="安全设置">
                <SecuritySettings />
              </a-tab-pane>
              <a-tab-pane key="3" title="实名认证">
                <Accreditation />
              </a-tab-pane>
            </a-tabs>
          </a-col>
          <a-col :span="24"> </a-col>
        </a-row>
      </a-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import myImage from "@/assets/img/my-image.jpg";
import BasicInfo from "@/views/personal/user-settings/components/basic-info.vue";
import SecuritySettings from "@/views/personal/user-settings/components/security-settings.vue";
import Accreditation from "@/views/personal/user-settings/components/accreditation.vue";
import useGlobalProperties from "@/hooks/useGlobalProperties";
const proxy = useGlobalProperties();

const type = ref("rounded");
const size = ref("medium");

const data = ref([
  {
    label: "用户ID：",
    value: "1000001"
  },
  {
    label: "用户名：",
    value: "admin"
  },
  {
    label: "用户昵称：",
    value: "兔子先森"
  },
  {
    label: "角色：",
    value: "系统管理员"
  },
  {
    label: "权限：",
    value: "全部"
  },
  {
    label: "实名认证：",
    value: "已认证"
  },
  {
    label: "邮箱：",
    value: "2547096351@qq.com"
  },
  {
    label: "手机号：",
    value: "15888888888"
  },
  {
    label: "注册时间：",
    value: "2024-06-20"
  }
]);

const toast = () => {
  proxy.$message.info("修改头像");
};
</script>

<style lang="scss" scoped>
.margin-top {
  margin-top: $padding;
}
</style>
