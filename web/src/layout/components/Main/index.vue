<template>
  <a-watermark :content="watermark" v-bind="watermarkConfig">
    <a-layout-content class="layout-main-content">
      <Tabs v-if="isTabs" />
      <router-view v-slot="{ Component }">
        <keep-alive :include="cacheRoutes">
          <component :is="Component" v-if="refreshPage" />
        </keep-alive>
      </router-view>
    </a-layout-content>
  </a-watermark>
</template>

<script setup lang="ts">
import Tabs from "@/layout/components/Tabs/index.vue";
import { storeToRefs } from "pinia";
import { useThemeConfig } from "@/store/modules/theme-config";
import { useRouteConfigStore } from "@/store/modules/route-config";
const themeStore = useThemeConfig();
let { refreshPage, isTabs, watermark, watermarkStyle, watermarkRotate, watermarkGap } = storeToRefs(themeStore);
const routerStore = useRouteConfigStore();
const { cacheRoutes } = storeToRefs(routerStore);

// 水印配置
const watermarkConfig = computed(() => {
  return {
    font: watermarkStyle.value,
    rotate: watermarkRotate.value,
    gap: watermarkGap.value
  };
});

watch(watermarkConfig, newv => {
  console.log(newv);
});
</script>

<style lang="scss" scoped>
.layout-main-content {
  display: flex;
  flex-direction: column;
  height: 100%;
}

// 修改左侧滚动条宽度-主要针对main窗口内的滚动条
:deep(.arco-scrollbar-thumb-direction-vertical .arco-scrollbar-thumb-bar) {
  width: 4px;
  margin-left: 8px;
}
</style>
