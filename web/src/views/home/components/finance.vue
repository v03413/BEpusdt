<template>
  <div class="shortcut-box">
    <div class="box-title">
      <div>数据汇总</div>
    </div>
    <a-divider :margin="16" />
    <a-grid class="finance-card" :cols="{ xs: 1, sm: 2, lg: 4, xl: 4 }" :col-gap="16" :row-gap="28">
      <a-grid-item v-for="(item, index) in financeData" :key="index">
        <a-card hoverable class="finance-a-card" :class="'animated-fade-up-' + index">
          <div class="finance-nav">
            <div class="tag-dot" :style="{ border: `3px solid ${item.color}` }"></div>
            <span class="finance-nav-title">{{ item.title }}</span>
          </div>
          <a-statistic
            :value-style="{
              fontSize: '13px',
              marginLeft: '16px',
              marginTop: '12px'
            }"
            :value="item.value"
            :value-from="0"
            :start="true"
            animation
            :show-group-separator="false"
            :precision="item.precision"
          />
        </a-card>
      </a-grid-item>
    </a-grid>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{
  homeData: any;
}>();

const financeData = ref([
  {
    id: 1,
    title: "累计订单",
    value: props.homeData?.total_count || 0,
    color: "#165DFF",
    precision: 0
  },
  {
    id: 2,
    title: "今日订单",
    value: props.homeData?.today_count || 0,
    color: "#6c73ff",
    precision: 0
  },
  {
    id: 3,
    title: "今日收款",
    value: props.homeData?.today_money || 0,
    color: "#39cbab",
    precision: 2
  },

  {
    id: 4,
    title: "总计收款",
    value: props.homeData?.total_money || 0,
    color: "#ff8625",
    precision: 2
  }
]);

watch(
  () => props.homeData,
  newData => {
    if (newData) {
      financeData.value[0].value = newData.total_count || 0;
      financeData.value[1].value = newData.today_count || 0;
      financeData.value[2].value = newData.today_money || 0;
      financeData.value[3].value = newData.total_money || 0;
    }
  },
  { immediate: true }
);
</script>

<style lang="scss" scoped>
.shortcut-box {
  .card-box {
    margin-bottom: $padding;
    .shortcut-card-label {
      width: 100px;
      margin-left: 20px;
      font-size: $font-size-body-3;
      color: $color-text-2;
    }
  }
  .card-middling {
    width: 200px;
  }
  .row-center {
    display: flex;
    align-items: center;
    justify-content: center;
  }
}
.box-title {
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: $font-size-body-3;
  color: $color-text-1;
}
.margin-left-text {
  margin-left: $margin-text;
}
</style>
