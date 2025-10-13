<template>
  <div class="container">
    <s-svg-icon :name="'snow'" :size="200" class="svg-icon" />
    <a-button type="primary" @click="onBack">返回首页</a-button>
  </div>
</template>

<script setup lang="ts">
const router = useRouter();
const onBack = () => {
  router.replace("/home");
};
</script>

<style lang="scss" scoped>
.container {
  display: flex;
  flex-direction: column;
  row-gap: 20px;
  align-items: center;
  justify-content: center;
  height: 100vh;
  .svg-icon {
    animation: rotate 10s linear infinite;
  }

  @keyframes rotate {
    from {
      transform: rotate(0deg);
    }
    to {
      transform: rotate(360deg);
    }
  }
}
</style>
