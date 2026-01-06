<template>
  <a-layout class="layout" id="s-fold-page">
    <a-layout-sider ref="siderContent" class="layout-sider" :collapsed-width="12" :width="width" :collapsed="collapsed">
      <div class="packup-btn" @click="collapsed = !collapsed">
        <icon-left v-if="!collapsed" />
        <icon-right v-else />
      </div>
      <div
        v-show="!collapsed"
        :style="{
          padding: siderPadding
        }"
      >
        <slot name="sider"> </slot>
      </div>
    </a-layout-sider>
    <a-layout-content>
      <div class="content">
        <slot name="content"> </slot>
      </div>
    </a-layout-content>
  </a-layout>
</template>

<script setup lang="ts">
const width = defineModel("width", {
  type: Number,
  default: 200
});
const collapsed = ref<boolean>(false);
const siderPadding = computed(() => (collapsed.value ? "0px" : "12px"));
</script>

<style lang="scss" scoped>
.layout {
  box-sizing: border-box;
  height: 100%;
  overflow: hidden;
  .layout-sider {
    position: relative;
    border-right: 1px solid $color-border-2;
    .packup-btn {
      position: absolute;
      top: 50%;
      right: calc(-24px / 2);
      display: flex;
      align-items: center;
      justify-content: center;
      width: 24px;
      height: 24px;
      color: $color-text-2;
      background: $color-bg-white;
      border-radius: 50%;
      box-shadow:
        inset 0 0 0 1px $color-border-2,
        0 0 10px 0 rgb(0 0 0 / 5%);
      box-shadow:
        inset 0 0 0 1px $color-border-2,
        0 0 10px 0 rgb(0 0 0 / 5%);
      box-shadow:
        inset 0 0 0 1px $color-border-2,
        0 0 10px 0 rgb(0 0 0 / 5%);
      transform: translateY(-50%);
      transition: all 0.1s;
      &:hover {
        transform: translateY(-50%) scale(1.1);
      }
    }
  }
  .content {
    box-sizing: border-box;
    height: 100%;
    padding: 12px;
  }
}
:deep(.arco-layout-sider-light) {
  box-shadow: unset;
}
</style>
