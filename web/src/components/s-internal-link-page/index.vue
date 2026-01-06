<template>
  <div class="snow-fill-pure">
    <div class="link-page-box">
      <div v-if="!isLoading" class="loading-page">
        <div class="dc-loader"></div>
      </div>
      <iframe v-show="isLoading" :src="link" class="iframe-size" id="internalLinkPage"> </iframe>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRoute } from "vue-router";

const route = useRoute();

const link = ref<string>(route.meta.link as string);

const isLoading = ref(false); // iframe是否加载完成
onMounted(() => {
  let iframe = document.getElementById("internalLinkPage") as HTMLElement;
  if (iframe) {
    // 加载完成
    iframe.onload = () => {
      isLoading.value = true;
    };
  }
});
</script>

<style lang="scss" scoped>
.link-page-box {
  display: flex;
  width: 100%;
  height: 100%;
  .iframe-size {
    flex: 1;
    border: 0;
  }
}
</style>
