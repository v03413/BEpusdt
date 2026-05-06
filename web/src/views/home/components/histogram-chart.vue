<template>
  <div class="histogram-scroll" :class="{ empty: isEmpty, scrollable: isScrollable }" ref="histogramScroll">
    <a-empty v-if="isEmpty" description="暂无订单数据" />
    <div v-show="!isEmpty" class="histogram-canvas" ref="sellHistogram" :style="{ width: `${chartWidth}px` }"></div>
  </div>
</template>

<script setup lang="ts">
import { default as VChart } from "@visactor/vchart";

const props = defineProps<{
  homeData: any;
}>();

const sellHistogram = ref();
const histogramScroll = ref();
const chartWidth = ref(0);
const isEmpty = ref(false);
const isScrollable = ref(false);
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

const getCssVarRgb = (name: string, fallback: string) => {
  const bodyValue = getComputedStyle(document.body).getPropertyValue(name).trim();
  const value = bodyValue || getComputedStyle(document.documentElement).getPropertyValue(name).trim();
  return value ? `rgb(${value})` : fallback;
};

const getThemeColors = () => {
  return [getCssVarRgb("--primary-6", "#165dff"), getCssVarRgb("--success-6", "#00b42a")];
};

const init = async (retryCount = 0) => {
  try {
    if (!props.homeData?.points || !sellHistogram.value || !histogramScroll.value) return;

    const points = props.homeData.points;
    if (!Array.isArray(points)) return;

    const gmvPaidByDate: Record<string, string> = {};
    const chartData = points.flatMap(point => {
      const date = String(point.date || "")
        .slice(5)
        .replace("-", "/");
      gmvPaidByDate[date] = Number(point.gmv_paid || 0).toFixed(2);
      return [
        {
          date,
          type: "订单总数",
          value: Number(point.orders_total || 0)
        },
        {
          date,
          type: "已支付订单",
          value: Number(point.orders_paid ?? point.orders_success ?? 0)
        }
      ];
    });

    if (chartData.length === 0) return;

    const hasOrderData = chartData.some(item => item.value > 0);
    if (!hasOrderData) {
      releaseChart();
      isEmpty.value = true;
      isScrollable.value = false;
      return;
    }

    isEmpty.value = false;
    await nextTick();

    const containerWidth = histogramScroll.value?.clientWidth || 0;
    const containerHeight = sellHistogram.value?.clientHeight || histogramScroll.value?.clientHeight || 0;
    if ((!containerWidth || !containerHeight) && retryCount < 8) {
      setTimeout(() => init(retryCount + 1), 100);
      return;
    }

    const renderWidth = containerWidth || 400;
    const renderHeight = containerHeight || 300;
    isScrollable.value = points.length > 7;
    chartWidth.value = isScrollable.value ? Math.max(renderWidth, points.length * 44) : renderWidth;
    await nextTick();

    releaseChart();

    const spec = {
      type: "bar",
      data: [{ id: "sellHistogramData", values: chartData }],
      xField: ["date", "type"],
      yField: "value",
      seriesField: "type",
      barWidth: 9,
      barGapInGroup: 2,
      legends: { visible: true, orient: "top", position: "middle" },
      color: {
        type: "ordinal",
        domain: ["订单总数", "已支付订单"],
        range: getThemeColors()
      },
      axes: [
        {
          orient: "bottom",
          label: {
            visible: true,
            formatMethod: (value: string) => (gmvPaidByDate[value] ? [value, gmvPaidByDate[value]] : value),
            style: {
              fontSize: 10,
              lineHeight: 13
            }
          },
          showAllGroupLayers: false
        },
        { orient: "left", label: { visible: true }, min: 0 }
      ],
      tooltip: {
        mark: {
          content: [
            {
              key: (datum: any) => datum["type"],
              value: (datum: any) => datum["value"]
            }
          ]
        }
      },
      animation: false,
      width: chartWidth.value || renderWidth,
      height: renderHeight
    };

    if (sellHistogram.value?.isConnected) {
      vchart = new VChart(spec as any, { dom: sellHistogram.value });
      vchart.renderSync();
    }
  } catch (error) {
    console.error("图表初始化错误:", error);
    releaseChart();
  }
};

watch(
  () => props.homeData,
  newData => {
    try {
      if (newData?.points) scheduleInit(120);
    } catch (error) {
      console.error("图表更新失败:", error);
    }
  },
  { immediate: true, deep: true }
);

onMounted(() => {
  try {
    if (typeof ResizeObserver !== "undefined" && histogramScroll.value) {
      resizeObserver = new ResizeObserver(() => scheduleInit(120));
      resizeObserver.observe(histogramScroll.value);
    }
    if (props.homeData?.points) scheduleInit(160);
  } catch (error) {
    console.error("图表初始化失败:", error);
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
.histogram-scroll {
  width: 100%;
  height: 100%;
  overflow-x: hidden;
  overflow-y: hidden;

  &.scrollable {
    overflow-x: auto;
  }

  &.empty {
    display: flex;
    align-items: center;
    justify-content: center;
  }
}

.histogram-canvas {
  height: 100%;
  min-width: 100%;
  min-height: 300px;
}
</style>
