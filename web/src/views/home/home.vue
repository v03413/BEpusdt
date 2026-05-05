<template>
  <div class="snow-page">
    <div class="home-page">
      <div class="dashboard-toolbar">
        <!-- 法币选择器 -->
        <div class="fiat-selector">
          <span class="label">交易法币：</span>
          <div class="fiat-options">
            <div
              v-for="item in fiatOptions"
              :key="item.value"
              class="fiat-option"
              :class="{ active: fiat === item.value }"
              @click="handleFiatChange(item.value)"
            >
              <span class="currency-symbol">{{ item.symbol }}</span>
              <span class="currency-name">{{ item.label }}</span>
            </div>
          </div>
        </div>

        <div class="range-actions">
          <a-range-picker
            v-if="range === 'custom'"
            v-model="customDates"
            format="YYYY-MM-DD"
            style="width: 240px"
            @change="handleCustomDateChange"
          />
          <a-select v-model="range" style="width: 132px" @change="handleRangeChange">
            <a-option v-for="item in rangeOptions" :key="item.value" :value="item.value">
              {{ item.label }}
            </a-option>
          </a-select>
          <a-button type="primary" :loading="loading" @click="forceRefresh">
            <template #icon><icon-refresh /></template>
            强制刷新
          </a-button>
        </div>
      </div>

      <!-- 财务指标 -->
      <Finance :home-data="home" />
      <!-- 数据图 -->
      <DataBox :home-data="home" />
    </div>
  </div>
</template>

<script setup lang="ts">
import Finance from "@/views/home/components/finance.vue";
import DataBox from "@/views/home/components/data-box.vue";
import { getDashboardHomeAPI } from "@/api/modules/home/index";

const fiat = ref("CNY");
const range = ref("7d");
const customDates = ref<any[]>([]);
const home = ref<any>(null);
const loading = ref(false);
const timezone = Intl.DateTimeFormat().resolvedOptions().timeZone || "Asia/Shanghai";
let dashboardRetryTimer: ReturnType<typeof setTimeout> | null = null;

const fiatOptions = ref([
  { value: "CNY", label: "人民币", symbol: "¥" },
  { value: "USD", label: "美元", symbol: "$" },
  { value: "EUR", label: "欧元", symbol: "€" },
  { value: "GBP", label: "英镑", symbol: "£" },
  { value: "JPY", label: "日元", symbol: "¥" }
]);

const rangeOptions = [
  { value: "today", label: "今天" },
  { value: "7d", label: "最近 7 天" },
  { value: "30d", label: "最近 30 天" },
  { value: "custom", label: "自定义" }
];

const clearDashboardRetry = () => {
  if (dashboardRetryTimer) {
    clearTimeout(dashboardRetryTimer);
    dashboardRetryTimer = null;
  }
};

const getDashboardHome = async (force = false, retryCount = 0) => {
  if (range.value === "custom" && (!Array.isArray(customDates.value) || customDates.value.length !== 2)) return;

  if (retryCount === 0) {
    clearDashboardRetry();
  }

  loading.value = true;
  try {
    const params: any = {
      range: range.value,
      tz: timezone,
      fiat: fiat.value,
      force
    };
    if (range.value === "custom") {
      params.from = customDates.value[0];
      params.to = customDates.value[1];
    }

    const data = await getDashboardHomeAPI(params);
    if (!data?.data) {
      throw new Error("仪表盘数据为空");
    }
    home.value = data.data;
  } catch (error) {
    if (retryCount < 3) {
      dashboardRetryTimer = setTimeout(() => {
        getDashboardHome(force, retryCount + 1);
      }, (retryCount + 1) * 1000);
      return;
    }
    console.error("获取首页统计失败:", error);
  } finally {
    loading.value = false;
  }
};

// 处理法币切换
const handleFiatChange = (value: string) => {
  fiat.value = value;
  getDashboardHome();
};

const handleRangeChange = () => {
  if (range.value !== "custom") {
    getDashboardHome();
  }
};

const handleCustomDateChange = () => {
  getDashboardHome();
};

const forceRefresh = () => {
  getDashboardHome(true);
};

onMounted(() => {
  getDashboardHome();
});

onUnmounted(() => {
  clearDashboardRetry();
});
</script>

<style lang="scss" scoped>
.home-page {
  padding: $padding;
  background: $color-bg-1;
}

.dashboard-toolbar {
  margin-bottom: 16px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  flex-wrap: wrap;
}

.fiat-selector {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  background: $color-bg-2;
  border-radius: 8px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.05);
  border: 1px solid $color-border-2;

  .label {
    font-size: 13px;
    font-weight: 500;
    color: $color-text-2;
    white-space: nowrap;
  }

  .fiat-options {
    display: flex;
    gap: 6px;
  }

  .fiat-option {
    display: flex;
    align-items: center;
    gap: 4px;
    padding: 6px 12px;
    border-radius: 6px;
    cursor: pointer;
    transition: all 0.2s ease;
    border: 1px solid $color-border-2;
    background: $color-fill-1;
    min-width: 60px;
    justify-content: center;

    &:hover {
      border-color: $color-primary;
      background: rgb(var(--primary-1));
    }

    &.active {
      background: $color-primary;
      border-color: $color-primary;
      color: #fff;
      box-shadow: 0 2px 8px rgb(var(--primary-6) / 25%);

      .currency-symbol,
      .currency-name {
        color: #fff;
      }
    }

    .currency-symbol {
      font-size: 14px;
      font-weight: bold;
      color: $color-primary;
    }

    .currency-name {
      font-size: 11px;
      color: $color-text-2;
      font-weight: 500;
    }
  }
}

.range-actions {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 10px;
  margin-left: auto;
}

// 响应式设计
@media (max-width: 768px) {
  .dashboard-toolbar {
    align-items: stretch;
  }

  .fiat-selector {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
    padding: 10px 12px;

    .fiat-options {
      width: 100%;
      justify-content: space-between;
      flex-wrap: wrap;
    }

    .fiat-option {
      flex: 1;
      min-width: auto;
      margin-bottom: 4px;
    }
  }

  .range-actions {
    width: 100%;
    flex-wrap: wrap;
    justify-content: flex-start;
    margin-left: 0;
  }
}
</style>
