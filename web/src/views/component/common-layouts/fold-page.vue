<template>
  <div class="container">
    <s-fold-page>
      <template #sider>
        <div>左侧</div>
      </template>
      <template #content>
        <div>右侧</div>
      </template>
    </s-fold-page>
  </div>
</template>

<script setup lang="ts"></script>

<style lang="scss" scoped>
.container {
  height: 100%;
  border: 1px solid $color-border-2;
}
</style>
