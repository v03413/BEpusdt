<template>
  <div class="snow-page">
    <div class="snow-inner">
      <a-card title="路由测试页">
        <a-space direction="vertical" fill>
          <div>路由：{{ route.path }}</div>
          <div>参数：{{ JSON.stringify(route.query) }}</div>
          <a-input :style="{ width: '320px' }" placeholder="请输入内容测试页面缓存" allow-clear v-model="form" />
        </a-space>
      </a-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouteConfigStore } from "@/store/modules/route-config";
defineOptions({ name: "common-route" });
const routerStore = useRouteConfigStore();
const route = useRoute();
const form = ref("");

routerStore.setTabsTitle(`${route.query.name} - ${route.query.text}`);
</script>

<style lang="scss" scoped></style>
