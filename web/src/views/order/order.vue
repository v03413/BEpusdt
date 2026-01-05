<template>
  <div class="order">
    <div class="snow-page">
      <div class="snow-inner">
        <a-form ref="formRef" :model="formData.form" auto-label-width>
          <a-row :gutter="16">
            <a-col :xs="24" :sm="12" :md="12" :lg="8" :xl="6">
              <a-form-item field="order_id" label="商户订单">
                <a-input v-model="formData.form.order_id" placeholder="请输入商户订单" allow-clear />
              </a-form-item>
            </a-col>
            <a-col :xs="24" :sm="12" :md="12" :lg="8" :xl="6">
              <a-form-item field="trade_type" label="交易类型">
                <a-select v-model="formData.form.trade_type" placeholder="请选择交易类型" allow-clear allow-search>
                  <a-option v-for="item in tradeTypeOptions" :key="item.value" :value="item.value">
                    {{ item.label }}
                  </a-option>
                </a-select>
              </a-form-item>
            </a-col>
            <a-col :xs="24" :sm="12" :md="12" :lg="8" :xl="6">
              <a-form-item field="status" label="订单状态">
                <a-select v-model="formData.form.status" placeholder="请选择订单状态" allow-clear>
                  <a-option v-for="item in statusOptions" :key="item.value" :value="item.value">
                    {{ item.label }}
                  </a-option>
                </a-select>
              </a-form-item>
            </a-col>

            <a-col :xs="24" :sm="12" :md="12" :lg="8" :xl="6" class="btn-col">
              <a-form-item label=" " style="margin-bottom: 0">
                <a-space>
                  <a-button type="primary" @click="getOrderList">
                    <template #icon><icon-search /></template>
                    查询
                  </a-button>
                  <a-button @click="onReset">
                    <template #icon><icon-refresh /></template>
                    重置
                  </a-button>
                  <a-button type="text" @click="formData.search = !formData.search">
                    {{ formData.search ? "收起" : "展开" }}
                    <icon-down :class="{ 'rotate-icon': formData.search }" />
                  </a-button>
                </a-space>
              </a-form-item>
            </a-col>

            <template v-if="formData.search">
              <a-col :xs="24" :sm="12" :md="12" :lg="8" :xl="6">
                <a-form-item field="trade_id" label="交易ID">
                  <a-input v-model="formData.form.trade_id" placeholder="请输入交易ID" allow-clear />
                </a-form-item>
              </a-col>
              <a-col :xs="24" :sm="12" :md="12" :lg="8" :xl="6">
                <a-form-item field="address" label="钱包地址">
                  <a-input v-model="formData.form.address" placeholder="请输入钱包地址" allow-clear />
                </a-form-item>
              </a-col>
              <a-col :xs="24" :sm="12" :md="12" :lg="8" :xl="6">
                <a-form-item field="createTime" label="创建时间">
                  <a-range-picker v-model="formData.form.createTime" show-time format="YYYY-MM-DD HH:mm:ss" style="width: 100%" />
                </a-form-item>
              </a-col>
            </template>
          </a-row>
        </a-form>

        <a-table
          row-key="id"
          size="small"
          :bordered="{ cell: true }"
          :scroll="{ x: 'max-content', y: '60vh' }"
          :loading="loading"
          :columns="columns"
          :data="data"
          v-model:selectedKeys="selectedKeys"
          :pagination="pagination"
          @page-change="pageChange"
          @page-size-change="pageSizeChange"
        >
          <template #wallet="{ record }">
            <a-tooltip :content="record.address" position="top">
              <span class="wallet-name">
                {{ record.wallet?.name || `⁉ ${record.address?.slice(-8) || "-"}` }}
              </span>
            </a-tooltip>
          </template>

          <template #money="{ record }">
            <span>
              {{ record.money }}
              <a-tag size="mini" color="arcoblue">{{ record.fiat }}</a-tag>
            </span>
          </template>

          <template #status="{ record }">
            <a-tag size="small" :color="getStatusColor(record.status)">
              {{ getStatusText(record.status) }}
            </a-tag>
          </template>

          <template #notify_state="{ record }">
            <a-tag size="small" :color="record.status === 2 ? (record.notify_state === 1 ? 'blue' : 'red') : 'gray'">
              {{ record.status === 2 ? (record.notify_state === 1 ? "成功" : "失败") : "-" }}
            </a-tag>
          </template>

          <template #optional="{ record }">
            <a-space>
              <a-button size="mini" type="primary" @click="showDetail(record)">详情</a-button>
              <a-popconfirm
                content="即使用户没付款,也确认强制补单吗?"
                type="warning"
                @ok="onPaid(record)"
                :disabled="record.status === 2"
              >
                <a-button size="mini" type="primary" status="warning" :disabled="record.status === 2"> 补单 </a-button>
              </a-popconfirm>
            </a-space>
          </template>
        </a-table>
      </div>
    </div>

    <DetailModal :visible="detailVisible" :detailData="detailData" @close="closeDetail" />
  </div>
</template>

<script setup lang="ts">
import { listAPI, paidAPI } from "@/api/modules/order/index";
import { List, FormData, Pagination } from "./config";
import { Notification } from "@arco-design/web-vue";
import { useUserInfoStore } from "@/store/modules/user-info";
import DetailModal from "./detail.vue";
import { useOrderDetail } from "./detail";

const userStores = useUserInfoStore();
const { detailVisible, detailData, showDetail, closeDetail } = useOrderDetail();

const tradeTypeOptions = computed(() => Object.entries(userStores.trade_type).map(([value, label]) => ({ value, label })));

const statusOptions = [
  { value: 1, label: "等待支付" },
  { value: 2, label: "支付成功" },
  { value: 3, label: "交易过期" },
  { value: 4, label: "交易取消" },
  { value: 5, label: "等待确认" },
  { value: 6, label: "确认失败" }
];

const formData = reactive<FormData>({
  form: {
    order_id: "",
    trade_id: "",
    trade_type: "",
    address: "",
    status: undefined,
    createTime: []
  },
  search: false
});

const selectedKeys = ref<string[]>([]);
const loading = ref(false);
const data = reactive<List[]>([]);
const pagination = ref<Pagination>({
  showPageSize: true,
  showTotal: true,
  current: 1,
  pageSize: 10,
  total: 10
});

const columns = [
  { title: "ID", align: "center", dataIndex: "id", width: 80 },
  { title: "商户订单", align: "center", dataIndex: "order_id", width: 200, ellipsis: true, tooltip: true },
  { title: "交易类型", align: "center", dataIndex: "trade_type", width: 100 },
  { title: "订单汇率", align: "center", dataIndex: "rate", width: 100 },
  { title: "实际付款", align: "center", dataIndex: "amount", width: 120 },
  { title: "订单金额", align: "center", dataIndex: "money", slotName: "money", width: 150 },
  { title: "收款钱包", align: "center", dataIndex: "wallet.name", slotName: "wallet", width: 150, ellipsis: true },
  { title: "交易状态", dataIndex: "status", align: "center", slotName: "status", width: 100 },
  { title: "回调", dataIndex: "notify_state", align: "center", slotName: "notify_state", width: 80 },
  { title: "创建时间", dataIndex: "created_at", align: "center", width: 160 },
  { title: "操作", slotName: "optional", align: "center", fixed: "right", width: 150 }
];

const statusMap: Record<number, { color: string; text: string }> = {
  1: { color: "blue", text: "等待支付" },
  2: { color: "green", text: "交易成功" },
  3: { color: "gray", text: "交易过期" },
  4: { color: "gold", text: "交易取消" },
  5: { color: "pinkpurple", text: "等待确认" },
  6: { color: "red", text: "确认失败" }
};

const getStatusColor = (status: number): string => statusMap[status]?.color || "gray";
const getStatusText = (status: number): string => statusMap[status]?.text || "未知";

const pageChange = (page: number) => {
  pagination.value.current = page;
  getOrderList();
};

const pageSizeChange = (pageSize: number) => {
  pagination.value.pageSize = pageSize;
  getOrderList();
};

const onReset = () => {
  formData.form = {
    order_id: "",
    trade_id: "",
    trade_type: "",
    address: "",
    status: undefined,
    createTime: []
  };
  getOrderList();
};

const getOrderList = async () => {
  try {
    loading.value = true;

    const params: any = {
      page: pagination.value.current,
      size: pagination.value.pageSize,
      sort: "desc",
      keyword: "",
      order_id: formData.form.order_id,
      trade_id: formData.form.trade_id,
      address: formData.form.address,
      trade_type: formData.form.trade_type
    };

    // 添加状态筛选
    if (formData.form.status !== undefined) {
      params.status = formData.form.status;
    }

    // 添加时间范围筛选
    if (formData.form.createTime && formData.form.createTime.length === 2) {
      params.start_at = formData.form.createTime[0];
      params.end_at = formData.form.createTime[1];
    }

    const res = await listAPI(params);

    data.length = 0;
    data.push(...res.data);
    pagination.value.total = res.total;
  } finally {
    loading.value = false;
  }
};

const onPaid = async (record: List) => {
  try {
    await paidAPI({ id: record.id });
    getOrderList();
    Notification.success("补单成功");
  } catch (error) {
    Notification.error(error);
  }
};

getOrderList();
</script>

<style lang="scss" scoped>
.order {
  height: 100%;

  .snow-page {
    height: 100%;

    .snow-inner {
      height: 100%;
      display: flex;
      flex-direction: column;

      .arco-form {
        margin-bottom: 16px;
      }

      .arco-table-container {
        flex: 1;
        overflow: hidden;
      }
    }
  }
}

.btn-col {
  .arco-form-item {
    margin-bottom: 0;
  }
}

.rotate-icon {
  transform: rotate(180deg);
  transition: transform 0.3s;
}

.wallet-name {
  cursor: help;
  color: #165dff;

  &:hover {
    text-decoration: underline;
  }
}

// 响应式处理
@media (max-width: 1200px) {
  :deep(.arco-table-th),
  :deep(.arco-table-td) {
    padding: 8px 6px !important;
    font-size: 12px;
  }
}

@media (max-width: 768px) {
  :deep(.arco-modal) {
    width: 95vw !important;
    margin: 10px;
  }

  :deep(.arco-table-th),
  :deep(.arco-table-td) {
    padding: 6px 4px !important;
    font-size: 11px;
  }
}
</style>
