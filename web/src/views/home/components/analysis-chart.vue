<template>
  <div class="analysis-chart" :class="{ empty: isEmpty }" ref="dailyAnalysis">
    <a-empty v-if="isEmpty" description="暂无交易数据" />
  </div>
</template>

<script setup lang="ts">
import { default as VChart } from "@visactor/vchart";

const props = defineProps<{
  homeData: any;
}>();

const dailyAnalysis = ref();
const isEmpty = ref(false);
let vchart: any = null;
let resizeObserver: ResizeObserver | null = null;
let initTimer: ReturnType<typeof setTimeout> | null = null;

const releaseChart = () => {
  if (vchart) {
    vchart.release();
    vchart = null;
  }
};

const scheduleInit = (delay = 0) => {
  if (initTimer) {
    clearTimeout(initTimer);
  }
  initTimer = setTimeout(() => {
    initTimer = null;
    init();
  }, delay);
};

const init = async (retryCount = 0) => {
  try {
    if (!props.homeData?.token_map || !dailyAnalysis.value) return;

    const tokenMap = props.homeData.token_map;
    if (typeof tokenMap !== "object" || tokenMap === null) return;

    const totalAmount = Object.values(tokenMap).reduce((sum: number, value: any) => {
      const numValue = Number(value);
      return sum + (isNaN(numValue) ? 0 : numValue);
    }, 0);

    if (totalAmount === 0) {
      releaseChart();
      isEmpty.value = true;
      return;
    }

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

    if (chartData.length === 0) {
      releaseChart();
      isEmpty.value = true;
      return;
    }

    isEmpty.value = false;
    await nextTick();

    const containerWidth = dailyAnalysis.value?.clientWidth || 0;
    const containerHeight = dailyAnalysis.value?.clientHeight || 0;
    if ((!containerWidth || !containerHeight) && retryCount < 8) {
      setTimeout(() => init(retryCount + 1), 100);
      return;
    }

    releaseChart();

    const colorMap = {
      USDT: "#1E90FF",
      USDC: "#32CD32",
      TRX: "#FF4500",
      BNB: "#F5A623",
      ETH: "#722ED1"
    };

    const spec = {
      type: "pie",
      data: [{ id: "dailyAnalysisData", values: chartData }],
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
      width: containerWidth || 400,
      height: containerHeight || 300
    };

    if (dailyAnalysis.value?.isConnected) {
      vchart = new VChart(spec as any, { dom: dailyAnalysis.value });
      vchart.renderSync();
    }
  } catch (error) {
    console.error("分析图表初始化错误:", error);
    releaseChart();
  }
};

watch(
  () => props.homeData,
  newData => {
    try {
      if (newData?.token_map) scheduleInit(120);
    } catch (error) {
      console.error("分析图表更新失败:", error);
    }
  },
  { immediate: true, deep: true }
);

onMounted(() => {
  try {
    if (typeof ResizeObserver !== "undefined" && dailyAnalysis.value) {
      resizeObserver = new ResizeObserver(() => scheduleInit(120));
      resizeObserver.observe(dailyAnalysis.value);
    }
    if (props.homeData?.token_map) scheduleInit(160);
  } catch (error) {
    console.error("分析图表初始化失败:", error);
  }
});

onUnmounted(() => {
  try {
    if (initTimer) {
      clearTimeout(initTimer);
      initTimer = null;
    }
    resizeObserver?.disconnect();
    resizeObserver = null;
    releaseChart();
  } catch (error) {
    console.error("清理图表失败:", error);
  }
});
</script>

<style lang="scss" scoped>
.analysis-chart {
  width: 100%;
  height: 100%;
  min-height: 300px;

  &.empty {
    display: flex;
    align-items: center;
    justify-content: center;
  }
}
</style>
