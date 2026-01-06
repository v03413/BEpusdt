<template>
  <Transition :name="transitionPage" mode="out-in" appear>
    <!-- 向内传递插槽内容 -->
    <slot></slot>
  </Transition>
</template>

<script setup lang="ts">
import { storeToRefs } from "pinia";
import { useThemeConfig } from "@/store/modules/theme-config";

const themeStore = useThemeConfig();
const { transitionPage } = storeToRefs(themeStore);
</script>

<style lang="scss" scoped></style>
