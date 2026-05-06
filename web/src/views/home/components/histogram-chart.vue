<template>
  <div class="histogram-scroll" :class="{ empty: isEmpty, scrollable: isScrollable }" ref="histogramScroll">
    <a-empty v-if="isEmpty" description="暂无订单数据" />
    <div v-else class="histogram-canvas" :style="{ width: chartWidth }">
      <s-echarts :options="option" :theme="theme" :update-options="{ notMerge: true }" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { use } from "echarts/core";
import { BarChart } from "echarts/charts";
import { CanvasRenderer } from "echarts/renderers";
import { TooltipComponent, GridComponent, LegendComponent } from "echarts/components";
import { useChart } from "@/hooks/useChart";

use([TooltipComponent, GridComponent, LegendComponent, BarChart, CanvasRenderer]);

const props = defineProps<{
  homeData: any;
}>();

const histogramScroll = ref<HTMLElement>();
const containerWidth = ref(0);
let resizeObserver: ResizeObserver | null = null;

const points = computed(() => {
  const source = props.homeData?.points;
  return Array.isArray(source) ? source : [];
});

const dates = computed(() =>
  points.value.map(point =>
    String(point.date || "")
      .slice(5)
      .replace("-", "/")
  )
);

const ordersTotalData = computed(() => points.value.map(point => Number(point.orders_total || 0)));
const ordersPaidData = computed(() => points.value.map(point => Number(point.orders_paid ?? point.orders_success ?? 0)));

const gmvPaidByDate = computed(() => {
  return points.value.reduce<Record<string, string>>((result, point, index) => {
    result[dates.value[index]] = Number(point.gmv_paid || 0).toFixed(2);
    return result;
  }, {});
});

const isEmpty = computed(() => !ordersTotalData.value.some(Boolean) && !ordersPaidData.value.some(Boolean));
const isScrollable = computed(() => points.value.length > 7);
const chartWidth = computed(() => {
  if (!isScrollable.value) return "100%";
  return `${Math.max(containerWidth.value || 0, points.value.length * 58)}px`;
});

const { option, theme } = useChart((isDark, palette = []) => {
  return {
    backgroundColor: "transparent",
    color: [palette[0] || "#4080FF", palette[3] || "#4CD263"],
    tooltip: {
      trigger: "axis",
      borderWidth: 0,
      axisPointer: {
        type: "shadow",
        shadowStyle: {
          color: isDark ? "rgba(247, 248, 250,0.03)" : "rgba(0, 0, 0,0.03)"
        }
      },
      formatter: (params: any) => {
        const items = Array.isArray(params) ? params : [params];
        const date = items[0]?.name || "";
        const rows = items
          .map(item => `${item.marker} ${item.seriesName} <b style="margin-left: 16px">${item.value}</b>`)
          .join("<br/>");
        const gmv = gmvPaidByDate.value[date] || "0.00";
        return `<b>${date}</b><br/>${rows}<br/>已收款 <b style="margin-left: 16px">${gmv}</b>`;
      }
    },
    legend: {
      top: 0,
      left: "center",
      icon: "circle"
    },
    grid: {
      left: "3%",
      right: "4%",
      top: 42,
      bottom: "3%",
      containLabel: true
    },
    xAxis: [
      {
        type: "category",
        data: dates.value,
        axisLabel: {
          color: "#86909c",
          fontSize: 10,
          formatter: (value: string) => (gmvPaidByDate.value[value] ? `${value}\n${gmvPaidByDate.value[value]}` : value),
          interval: 0,
          lineHeight: 11
        },
        axisLine: {
          lineStyle: {
            color: isDark ? "#373738" : "#ebedf0"
          }
        }
      }
    ],
    yAxis: [
      {
        type: "value",
        min: 0,
        axisLabel: {
          color: "#86909c"
        },
        splitLine: {
          lineStyle: {
            color: isDark ? "#373738" : "#ebedf0"
          }
        }
      }
    ],
    series: [
      {
        name: "订单总数",
        type: "bar",
        barWidth: 9,
        barGap: "25%",
        data: ordersTotalData.value
      },
      {
        name: "已支付订单",
        type: "bar",
        barWidth: 9,
        data: ordersPaidData.value
      }
    ]
  };
});

onMounted(() => {
  if (typeof ResizeObserver !== "undefined" && histogramScroll.value) {
    resizeObserver = new ResizeObserver(entries => {
      containerWidth.value = entries[0]?.contentRect.width || histogramScroll.value?.clientWidth || 0;
    });
    resizeObserver.observe(histogramScroll.value);
  }
  containerWidth.value = histogramScroll.value?.clientWidth || 0;
});

onUnmounted(() => {
  resizeObserver?.disconnect();
  resizeObserver = null;
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
