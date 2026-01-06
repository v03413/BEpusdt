<template>
  <div style="height: 100%" ref="monthlyAnalysis"></div>
</template>

<script setup lang="ts">
import { default as VChart } from "@visactor/vchart";

const props = defineProps<{
  homeData: any;
}>();

const monthlyAnalysis = ref();
let vchart: any = null;

const init = () => {
  try {
    if (!props.homeData?.token_map || !monthlyAnalysis.value) return;

    const tokenMap = props.homeData.token_map;
    if (typeof tokenMap !== "object" || tokenMap === null) return;

    const totalAmount = Object.values(tokenMap).reduce((sum: number, value: any) => {
      const numValue = Number(value);
      return sum + (isNaN(numValue) ? 0 : numValue);
    }, 0);

    if (totalAmount === 0) return;

    const chartData = Object.entries(tokenMap)
      .map(([token, amount]) => {
        const numAmount = Number(amount);
        if (isNaN(numAmount) || numAmount <= 0) return null;
        const percentage = (numAmount / totalAmount) * 100;
        return {
          type: token,
          value: Number(percentage.toFixed(2)),
          amount: Number(numAmount.toFixed(2))
        };
      })
      .filter(item => item !== null);

    if (chartData.length === 0) return;

    if (vchart) {
      vchart.release();
      vchart = null;
    }

    const colorMap = {
      USDT: "#1E90FF",
      USDC: "#32CD32",
      TRX: "#FF4500"
    };

    const spec = {
      type: "pie",
      data: [{ id: "monthlyAnalysisData", values: chartData }],
      outerRadius: 0.8,
      innerRadius: 0.5,
      padAngle: 0.6,
      valueField: "value",
      categoryField: "type",
      color: {
        type: "ordinal",
        domain: chartData.map(item => item.type),
        range: chartData.map(item => colorMap[item.type as keyof typeof colorMap] || "#666666")
      },
      pie: {
        style: { cornerRadius: 0 },
        state: {
          hover: { outerRadius: 0.85, stroke: "#fff", lineWidth: 1 },
          selected: { outerRadius: 0.85, stroke: "#fff", lineWidth: 1 }
        }
      },
      legends: { visible: true, orient: "left" },
      label: { visible: true, content: "{type}: {value}%" },
      tooltip: {
        mark: {
          content: [
            {
              key: (datum: any) => datum["type"],
              value: (datum: any) => `${datum["value"]}% (${datum["amount"]})`
            }
          ]
        }
      },
      animation: false,
      width: monthlyAnalysis.value?.clientWidth || 400,
      height: monthlyAnalysis.value?.clientHeight || 300
    };

    if (monthlyAnalysis.value?.isConnected) {
      vchart = new VChart(spec as any, { dom: monthlyAnalysis.value });
      vchart.renderSync();
    }
  } catch (error) {
    console.error("分析图表初始化错误:", error);
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
      if (newData?.token_map) init();
    } catch (error) {
      console.error("分析图表更新失败:", error);
    }
  },
  { immediate: true, deep: true }
);

onMounted(() => {
  try {
    if (props.homeData?.token_map) init();
  } catch (error) {
    console.error("分析图表初始化失败:", error);
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
