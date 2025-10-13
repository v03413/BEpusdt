<template>
  <div style="height: 100%" ref="sellHistogram"></div>
</template>

<script setup lang="ts">
import { default as VChart } from "@visactor/vchart";

const props = defineProps<{
  homeData: any;
}>();

const sellHistogram = ref();
let vchart: any = null;

const init = () => {
  try {
    if (!props.homeData?.monthly || !sellHistogram.value) return;

    const monthlyData = props.homeData.monthly;
    if (typeof monthlyData !== "object" || monthlyData === null) return;

    const chartData = Object.entries(monthlyData).map(([month, sales]) => ({
      month,
      sales: Number(Number(sales).toFixed(2)) || 0
    }));

    if (chartData.length === 0) return;

    if (vchart) {
      vchart.release();
      vchart = null;
    }

    const spec = {
      type: "bar",
      data: [{ id: "sellHistogramData", values: chartData }],
      xField: "month",
      yField: "sales",
      barWidth: 10,
      barGapInGroup: 0,
      animation: false,
      width: sellHistogram.value?.clientWidth || 400,
      height: sellHistogram.value?.clientHeight || 300
    };

    if (sellHistogram.value?.isConnected) {
      vchart = new VChart(spec as any, { dom: sellHistogram.value });
      vchart.renderSync();
    }
  } catch (error) {
    console.error("图表初始化错误:", error);
    if (vchart) {
      vchart.release();
      vchart = null;
    }
  }
};

watch(
  () => props.homeData,
  newData => {
    try {
      if (newData?.monthly) init();
    } catch (error) {
      console.error("图表更新失败:", error);
    }
  },
  { immediate: true, deep: true }
);

onMounted(() => {
  try {
    if (props.homeData?.monthly) init();
  } catch (error) {
    console.error("图表初始化失败:", error);
  }
});

onUnmounted(() => {
  try {
    vchart?.release();
    vchart = null;
  } catch (error) {
    console.error("清理图表失败:", error);
  }
});
</script>

<style lang="scss" scoped></style>
