<template>
  <div class="page-404">
    <div>
      <s-svg-icon name="暂无权限" :size="500" />
    </div>
    <div class="prompt">
      <div class="title">401</div>
      <div class="text">抱歉，暂无访问权限~</div>
      <a-button type="primary" @click="onBack">立即返回</a-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from "vue-router";
const router = useRouter();
const onBack = () => {
  if (window.history.state.back !== null) {
    router.replace(window.history.state.back);
  } else {
    router.replace("/login");
  }
};
</script>

<style lang="scss" scoped>
.page-404 {
  box-sizing: border-box;
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100vh;
  padding: $padding;
  overflow: hidden;
  .prompt {
    row-gap: $padding;
    width: 250px;
  }
  .title {
    font-size: 80px;
    color: $color-text-1;
  }
  .text {
    margin-bottom: $padding;
    font-size: $font-size-body-3;
    color: $color-text-2;
  }
}
</style>
