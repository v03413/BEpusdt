<template>
  <a-layout-footer class="footer">
    <div class="footer_title">
      <div class="left-placeholder"></div>
      <span @click="onBEpusdt" class="primary-copyright">© 2025 BEpusdt. All rights reserved.</span>
      <span class="secondary-info">Powered by SnowAdmin</span>
    </div>
  </a-layout-footer>
</template>

<script setup lang="ts">
const onBEpusdt = () => {
  window.open("https://github.com/v03413/BEpusdt", "_blank");
};
</script>

<style lang="scss" scoped>
.footer {
  box-sizing: border-box;
  display: flex;
  align-items: center;
  height: 30px;
  border-top: $border-1 solid $color-border-2;
  padding: 0 16px;
}

.footer_title {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;

  .left-placeholder {
    flex: 1;
  }

  .primary-copyright {
    color: $color-text-2;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s ease;
    flex: 1;
    text-align: center;

    &:hover {
      color: $color-primary;
      transform: translateY(-1px);
    }
  }

  .secondary-info {
    color: $color-text-4;
    font-size: 12px;
    opacity: 0.6;
    font-weight: 300;
    flex: 1;
    text-align: right;
  }
}
</style>
