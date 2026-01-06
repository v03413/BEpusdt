<template>
  <div class="snow-page">
    <div class="external-link-page">
      <div class="external-link-inner-page">
        <slot></slot>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts"></script>

<style lang="scss" scoped>
.external-link-page {
  box-sizing: border-box;
  height: 100%;
  .external-link-inner-page {
    height: 100%;
    overflow: hidden;
    background: $color-bg-1;
  }
}
</style>
