<template>
  <div class="tools">
    <a-space>
      <slot name="left"> </slot>
    </a-space>
    <a-space>
      <slot name="right"> </slot>
    </a-space>
  </div>
</template>

<script setup lang="ts"></script>

<style lang="scss" scoped>
.tools {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  justify-content: space-between;
}
</style>
