<template>
  <div class="header_crumb">
    <ButtonCollapsed />
    <Breadcrumb />
  </div>
</template>

<script setup lang="ts">
import Breadcrumb from "@/layout/components/Header/components/Breadcrumb/index.vue";
import ButtonCollapsed from "@/layout/components/Header/components/button-collapsed/index.vue";
</script>

<style lang="scss" scoped>
.header_crumb {
  display: flex;
  align-items: center;
  width: 100%;
}
</style>
