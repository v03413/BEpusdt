<template>
  <div class="snow-page">
    <div class="snow-inner">
      <a-alert>
        当前页面
        <a-tag color="green"> 超级管理员 </a-tag>
        和
        <a-tag color="green"> 普通角色 </a-tag>
        都可查看
      </a-alert>

      <br />
      <a-card hoverable title="根据角色判断页面是否显示">
        <a-typography-paragraph> <a-typography-text bold> 当前页面权限 </a-typography-text></a-typography-paragraph>
        <s-code-view :code-json="pageRule" style="width: 100%" />
        <br />
        <a-typography-paragraph> <a-typography-text bold> 当前角色权限 </a-typography-text></a-typography-paragraph>
        <s-code-view :code-json="roles" style="width: 100%" />
      </a-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { storeToRefs } from "pinia";
import { useUserInfoStore } from "@/store/modules/user-info";
const route = useRoute();
let userInfoStore = useUserInfoStore();
let { account } = storeToRefs(userInfoStore);
const roles = computed(() => JSON.stringify(account.value.roles, null));
const pageRule = computed(() => JSON.stringify(route.meta.roles, null));
</script>

<style lang="scss" scoped></style>
