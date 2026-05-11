<template>
  <div class="analysis-chart" :class="{ empty: isEmpty }">
    <a-empty v-if="isEmpty" description="暂无交易数据" />
    <s-echarts v-else :options="option" :theme="theme" :update-options="{ notMerge: true }" />
  </div>
</template>

<script setup lang="ts">
import { use } from "echarts/core";
import { CanvasRenderer } from "echarts/renderers";
import { PieChart } from "echarts/charts";
import { TooltipComponent, LegendComponent } from "echarts/components";
import { useChart } from "@/hooks/useChart";

use([CanvasRenderer, PieChart, TooltipComponent, LegendComponent]);

const props = defineProps<{
  homeData: any;
}>();

const colorMap: Record<string, string> = {
  USDT: "#1E90FF",
  USDC: "#32CD32",
  TRX: "#FF4500",
  BNB: "#F5A623",
  ETH: "#722ED1"
};

const chartData = computed(() => {
  const tokenMap = props.homeData?.token_map;
  if (!tokenMap || typeof tokenMap !== "object") return [];

  const totalAmount = Object.values(tokenMap).reduce((sum: number, value: any) => {
    const numValue = Number(value);
    return sum + (Number.isNaN(numValue) ? 0 : numValue);
  }, 0);
  if (totalAmount <= 0) return [];

  return Object.entries(tokenMap)
    .map(([token, amount]) => {
      const numAmount = Number(amount);
      if (Number.isNaN(numAmount) || numAmount <= 0) return null;

      return {
        name: token,
        value: Number(((numAmount / totalAmount) * 100).toFixed(2)),
        amount: Number(numAmount.toFixed(2))
      };
    })
    .filter(Boolean) as Array<{ name: string; value: number; amount: number }>;
});

const isEmpty = computed(() => chartData.value.length === 0);

const { option, theme } = useChart((_, palette = []) => {
  const colors = chartData.value.map(item => colorMap[item.name] || palette[chartData.value.indexOf(item) % palette.length] || "#666666");

  return {
    backgroundColor: "transparent",
    color: colors,
    tooltip: {
      trigger: "item",
      borderWidth: 0,
      formatter: (params: any) => {
        const amount = params.data?.amount ?? 0;
        return `<b>${params.name}</b><br/>${params.marker} 占比 <b style="margin-left: 16px">${params.value}%</b><br/>金额 <b style="margin-left: 16px">${amount}</b>`;
      }
    },
    legend: {
      orient: "vertical",
      left: "left",
      top: "middle",
      icon: "circle",
      itemHeight: 10,
      itemGap: 14,
      data: chartData.value.map(item => item.name)
    },
    series: [
      {
        type: "pie",
        center: ["62%", "50%"],
        padAngle: 1,
        radius: ["45%", "70%"],
        label: {
          color: "inherit",
          formatter: "{b}: {c}%"
        },
        data: chartData.value,
        emphasis: {
          scale: true,
          scaleSize: 10
        }
      }
    ]
  };
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
