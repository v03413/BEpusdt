<template>
  <VCharts v-if="renderChart" :option="options" :autoresize="autoResize" :style="{ width, height }" />
</template>

<script lang="ts" setup>
import VCharts from "vue-echarts";
defineProps({
  options: {
    type: Object,
    default() {
      return {};
    }
  },
  autoResize: {
    type: Boolean,
    default: true
  },
  width: {
    type: String,
    default: "100%"
  },
  height: {
    type: String,
    default: "100%"
  }
});
const renderChart = ref(false);
// wait container expand
nextTick(() => {
  renderChart.value = true;
});
</script>
