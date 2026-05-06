<template>
  <div class="data-box">
    <div class="sell-histogram">
      <div>
        <span class="data-title">收款趋势</span>
        <span class="data-subtext">单位: 天</span>
      </div>

      <div class="chart-body">
        <HistogramChart :home-data="props.homeData" />
      </div>
    </div>
    <div class="daily-analysis">
      <div>
        <span class="data-title">交易占比</span>
      </div>

      <div class="chart-body">
        <AnalysisChart :home-data="props.homeData" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import HistogramChart from "@/views/home/components/histogram-chart.vue";
import AnalysisChart from "@/views/home/components/analysis-chart.vue";

const props = defineProps<{
  homeData: any;
}>();
</script>

<style lang="scss" scoped>
.data-box {
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
  margin-top: calc($padding * 2);

  .sell-histogram,
  .daily-analysis {
    display: flex;
    flex-direction: column;
    height: 400px;
    padding-bottom: $padding;
    min-width: 0;
  }

  .sell-histogram {
    width: 100%;
  }

  .daily-analysis {
    width: 100%;
  }

  .chart-body {
    flex: 1;
    min-height: 0;
    min-width: 0;
  }

  .data-title {
    font-size: $font-size-body-3;
    color: $color-text-1;
  }
  .data-subtext {
    margin-left: $margin-text;
    font-size: $font-size-body-2;
    color: $color-text-2;
  }
}

@media screen and (min-width: $xxl) {
  .data-box {
    .sell-histogram {
      width: calc(100% - 600px - $padding);
    }

    .daily-analysis {
      width: 600px;
      margin-left: $padding;
    }
  }
}
</style>
