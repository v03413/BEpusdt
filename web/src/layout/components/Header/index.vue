<template>
  <a-layout-header class="header">
    <HeaderLeft />
    <HeaderRight />
  </a-layout-header>
</template>
<script setup lang="ts">
import HeaderLeft from "@/layout/components/Header/components/header-left/index.vue";
import HeaderRight from "@/layout/components/Header/components/header-right/index.vue";
</script>

<style lang="scss" scoped>
.header {
  position: relative;
  box-sizing: border-box;
  display: flex;
  align-items: center;
  height: 60px;
  padding: 0 $padding;
  border-bottom: $border-1 solid $color-border-2;
}
</style>
