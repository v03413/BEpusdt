<template>
  <div>
    <s-lang-provider>
      <component :is="layouts[layoutType]" />
    </s-lang-provider>
  </div>
</template>

<script setup lang="ts">
import { storeToRefs } from "pinia";
import { useThemeConfig } from "@/store/modules/theme-config";

const themeStore = useThemeConfig();
const { layoutType } = storeToRefs(themeStore);

// 引入组件-异步组件
const layouts: any = {
  layoutDefaults: defineAsyncComponent(() => import("@/layout/layout-defaults/index.vue")),
  layoutHead: defineAsyncComponent(() => import("@/layout/layout-head/index.vue")),
  layoutMixing: defineAsyncComponent(() => import("@/layout/layout-mixing/index.vue"))
};
</script>

<style lang="scss" scoped></style>
