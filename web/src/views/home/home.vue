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
const home = ref<any>(null);

const fiatOptions = ref([
  { value: "CNY", label: "人民币", symbol: "¥" },
  { value: "USD", label: "美元", symbol: "$" },
  { value: "EUR", label: "欧元", symbol: "€" },
  { value: "GBP", label: "英镑", symbol: "£" },
  { value: "JPY", label: "日元", symbol: "¥" }
]);

const getDashboardHome = async () => {
  const data = await getDashboardHomeAPI({ fiat: fiat.value });
  home.value = data.data;
};

// 处理法币切换
const handleFiatChange = (value: string) => {
  fiat.value = value;
  getDashboardHome();
};

onMounted(() => {
  getDashboardHome();
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
}
</style>
