<template>
  <div class="snow-page">
    <div class="snow-inner">
      <div class="my-page">
        <div class="title">自定义节流指令，连续点击按钮，每间隔1s执行一次</div>
        <div class="button-box">
          <a-button type="primary" style="width: 100px" v-throttle="onClick">1s节流</a-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Message } from "@arco-design/web-vue";
const onClick = () => {
  Message.success("节流执行了");
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
