<template>
  <div>
    <router-view />
  </div>
</template>

<script setup lang="ts">
import { useThemeMethods } from "@/hooks/useThemeMethods";

// 初始化主题
const onTheme = () => {
  let { initTheme } = useThemeMethods();
  initTheme();
};
onTheme();
</script>

<style lang="scss" scoped></style>
