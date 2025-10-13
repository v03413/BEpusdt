<template>
  <div class="rate-list">
    <div class="snow-page">
      <div class="snow-inner">
        <a-form ref="formRef" auto-label-width :model="formData.form">
          <a-row :gutter="16">
            <a-col :xs="24" :sm="24" :md="12" :lg="12" :xl="6" :xxl="6">
              <a-select v-model="formData.form.fiat" placeholder="请选择交易法币" allow-clear allow-search>
                <a-option v-for="item in tradeFiatOptions" :key="item.value" :value="item.value">
                  {{ item.value }}
                </a-option>
              </a-select>
            </a-col>
            <a-col :xs="24" :sm="24" :md="12" :lg="12" :xl="6" :xxl="6">
              <a-select v-model="formData.form.crypto" placeholder="请选择加密货币" allow-clear allow-search>
                <a-option v-for="item in tradeCryptoOptions" :key="item.value" :value="item.value">
                  {{ item.value }}
                </a-option>
              </a-select>
            </a-col>
            <a-col :xs="24" :sm="24" :md="12" :lg="12" :xl="6" :xxl="6">
              <a-range-picker
                v-model="formData.form.datetime"
                :placeholder="['开始时间', '结束时间']"
                show-time
                format="YYYY-MM-DD HH:mm"
                style="width: 100%"
                allow-clear
              />
            </a-col>
            <a-col :xs="24" :sm="24" :md="12" :lg="12" :xl="6" :xxl="6">
              <a-space class="search-btn">
                <a-button type="primary" @click="getCommonTableList">
                  <template #icon><icon-search /></template>
                  查询
                </a-button>
                <a-button @click="onReset">
                  <template #icon><icon-refresh /></template>
                  重置
                </a-button>
                <a-button type="primary" status="success" @click="onSync">
                  <template #icon><icon-loop /></template>
                  立刻同步
                </a-button>
              </a-space>
            </a-col>
          </a-row>
        </a-form>
        <a-divider :margin="10" />
        <a-table
          row-key="key"
          size="medium"
          :bordered="{ cell: true }"
          :scroll="{ x: 1400, y: 600 }"
          :loading="loading"
          :columns="columns"
          :data="data"
          v-model:selectedKeys="selectedKeys"
          :pagination="pagination"
          @page-change="pageChange"
          @page-size-change="pageSizeChange"
        >
          <template #fiat="{ record }">
            <span class="fiat-display">
              {{ getFiatFlag(record.fiat) }} <strong>{{ record.fiat }}</strong>
            </span>
          </template>
          <template #crypto="{ record }">
            <a-tag :color="getCryptoColor(record.crypto)" bordered>
              {{ record.crypto }}
            </a-tag>
          </template>
        </a-table>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { getRateListAPI, onSyncAPI } from "@/api/modules/rate/index";
import { List, FormData, Pagination } from "./list";
import { useUserInfoStore } from "@/store/modules/user-info";
import { getFiatFlag, getCryptoColor } from "@/views/rate/common";
import { Message } from "@arco-design/web-vue";

const userStores = useUserInfoStore();
const tradeCryptoOptions = computed(() =>
  (userStores.trade_crypto ?? []).map((crypto: string) => ({ value: crypto, label: crypto }))
);
const tradeFiatOptions = computed(() => (userStores.trade_fiat ?? []).map((fiat: string) => ({ value: fiat, label: fiat })));

const formData = reactive<FormData>({
  form: { fiat: "", crypto: "", datetime: [] },
  search: false
});

const selectedKeys = ref<string[]>([]);
const pagination = ref<Pagination>({
  showPageSize: true,
  showTotal: true,
  current: 1,
  pageSize: 10,
  total: 10
});

const loading = ref(false);
const data = reactive<List[]>([]);

const columns = [
  { title: "ID", dataIndex: "id", align: "center", width: "100" },
  { title: "交易法币", align: "center", dataIndex: "fiat", width: "180", slotName: "fiat" },
  { title: "加密货币", align: "center", dataIndex: "crypto", width: "180", slotName: "crypto" },
  { title: "订单汇率", align: "center", dataIndex: "rate" },
  { title: "基准汇率", dataIndex: "raw_rate", align: "center", slotName: "raw_rate" },
  { title: "同步时间", dataIndex: "created_at", align: "center", slotName: "created_at" }
];

const getCommonTableList = async () => {
  try {
    loading.value = true;
    const res = await getRateListAPI({
      page: pagination.value.current,
      size: pagination.value.pageSize,
      sort: "desc",
      keyword: "",
      fiat: formData.form.fiat,
      crypto: formData.form.crypto,
      datetime: formData.form.datetime,
      status: 99
    });

    data.length = 0;
    data.push(...res.data);
    pagination.value.total = res.total;
  } finally {
    loading.value = false;
  }
};

const onSync = async () => {
  try {
    await onSyncAPI({});
    Message.success("同步成功！");
    getCommonTableList();
  } catch (error) {
    console.error("同步失败:", error);
    Message.error("同步失败，请重试");
  }
};

const onReset = () => {
  formData.form = { fiat: "", crypto: "", datetime: [] };
  getCommonTableList();
};

const pageChange = (page: number) => {
  pagination.value.current = page;
  getCommonTableList();
};

const pageSizeChange = (pageSize: number) => {
  pagination.value.pageSize = pageSize;
  getCommonTableList();
};

getCommonTableList();
</script>

<style lang="scss" scoped>
.fiat-display {
  display: inline-flex;
  align-items: center;
  gap: 4px;
}
</style>
