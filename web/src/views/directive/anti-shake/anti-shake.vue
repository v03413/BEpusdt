<template>
  <div class="snow-page">
    <div class="snow-inner">
      <div class="title">自定义防抖指令，连续点击按钮，0.5s后执行一次</div>
      <div class="button-box"><a-button type="primary" style="width: 100px" v-antiShake="onClick">防抖</a-button></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Message } from "@arco-design/web-vue";
const onClick = () => {
  Message.success("0.5s后执行了");
};
</script>

<style lang="scss" scoped>
.title {
  font-size: $font-size-title-1;
  color: $color-text-1;
  text-align: center;
}
.button-box {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100px;
}
</style>
