<template>
  <div class="snow-page">
    <div class="snow-inner">
      <a-row class="grid-demo" :gutter="24">
        <a-col :span="12">
          <a-card title="图标选择器" :style="{ width: '100%' }">
            <s-select-icon type="arco" v-model:="iconName" />
            <div class="target-title">当前选择的图标：</div>
            <component v-if="iconName" :is="iconName" :size="50"></component>
            <a-empty v-else />
          </a-card>
        </a-col>
        <a-col :span="12">
          <a-card title="SVG选择器" :style="{ width: '100%' }">
            <s-select-icon type="svg" v-model:="svgName" />
            <div class="target-title">当前选择的图标：</div>
            <s-svg-icon v-if="svgName" :name="svgName" :size="50" />
            <a-empty v-else />
          </a-card>
        </a-col>
      </a-row>
    </div>
  </div>
</template>

<script setup lang="ts">
const iconName = ref<string>("icon-menu");
const svgName = ref<string>("home");
</script>

<style lang="scss" scoped>
.target-title {
  margin: $margin 0;
}
</style>
