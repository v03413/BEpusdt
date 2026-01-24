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
                  <a-popconfirm :content="batchDelConfirm" type="warning" @ok="onBatchDelete">
                    <a-button v-show="selectedKeys.length > 0" type="primary" status="danger">
                      <template #icon><icon-delete /></template>
                      删除
                    </a-button>
                  </a-popconfirm>
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
          :row-selection="orderSelection"
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

          <template #amount="{ record }">
            <span>
              {{ record.amount }}
              <a-tag size="mini" :color="getCryptoColor(record.crypto)" bordered style="margin-left: 4px">{{
                record.crypto
              }}</a-tag>
            </span>
          </template>

          <template #money="{ record }">
            <span>
              {{ record.money }}
              <a-tag size="mini" color="arcoblue" style="margin-left: 4px">{{ record.fiat }}</a-tag>
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

          <!-- 不常用的操作优先放置在详情页，尽量保持第一视角的干净整洁 -->
          <template #optional="{ record }">
            <a-space>
              <a-button size="mini" type="primary" @click="showDetail(record)">详情</a-button>
              <a-button
                size="mini"
                type="primary"
                status="warning"
                :disabled="record.status === 2"
                @click="showPaidModal(record)"
              >
                补单
              </a-button>
              <a-button
                v-if="record.status === 2 && record.notify_state === 0"
                size="mini"
                type="primary"
                status="danger"
                @click="handleManualNotify(record)"
              >
                手动回调
              </a-button>
            </a-space>
          </template>
        </a-table>
      </div>
    </div>

    <DetailModal :visible="detailVisible" :detailData="detailData" @close="closeDetail" />

    <!-- 补单弹窗 -->
    <a-modal
      v-model:visible="paidModalVisible"
      title="确认补单操作"
      @ok="confirmPaid"
      @cancel="closePaidModal"
      ok-text="确认补单"
      cancel-text="取消"
      width="500px"
      :mask-closable="false"
    >
      <div class="paid-modal-content">
        <a-alert type="warning" style="margin-bottom: 20px">
          <template #icon>
            <icon-exclamation-circle-fill />
          </template>
          <div style="font-weight: 500">注意</div>
          <div style="font-size: 13px; margin-top: 4px; color: #666">
            补单操作将强制标记订单为已支付，即使用户实际未付款、谨慎操作!
          </div>
        </a-alert>

        <a-form :model="paidForm" layout="vertical">
          <a-form-item field="ref_hash" label="交易哈希" :rules="[{ maxLength: 200, message: '哈希值不能超过200个字符' }]">
            <a-input v-model="paidForm.ref_hash" placeholder="请输入区块链交易哈希值(可选)" allow-clear>
              <template #prefix>
                <icon-link />
              </template>
            </a-input>
            <template #extra>
              <div style="font-size: 12px; color: #86909c; margin-top: 4px">如有实际交易,建议填写对应的区块链交易哈希值</div>
            </template>
          </a-form-item>
        </a-form>
      </div>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { listAPI, paidAPI, manualNotifyAPI, delOrderApi } from "@/api/modules/order/index";
import { List, FormData, Pagination } from "./config";
import { Notification, Modal } from "@arco-design/web-vue";
import { useUserInfoStore } from "@/store/modules/user-info";
import DetailModal from "./detail.vue";
import { useOrderDetail } from "./detail";
import { getCryptoColor } from "@/views/rate/common";

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
const orderSelection = reactive({
  type: "checkbox",
  showCheckedAll: true,
  onlyCurrent: false
});
const batchDelConfirm = computed(() => `确定删除这${selectedKeys.value.length}条数据吗？`);
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
  { title: "交易类型", align: "center", dataIndex: "trade_type", width: 120 },
  { title: "交易数额", align: "center", dataIndex: "amount", slotName: "amount", width: 150 },
  { title: "交易金额", align: "center", dataIndex: "money", slotName: "money", width: 150 },
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

const paidModalVisible = ref(false);
const paidForm = reactive({
  ref_hash: "",
  recordId: 0
});

const showPaidModal = (record: List) => {
  paidForm.recordId = record.id;
  paidForm.ref_hash = "";
  paidModalVisible.value = true;
};

const closePaidModal = () => {
  paidModalVisible.value = false;
  paidForm.ref_hash = "";
  paidForm.recordId = 0;
};

const confirmPaid = async () => {
  try {
    await paidAPI({
      id: paidForm.recordId,
      ref_hash: paidForm.ref_hash || "" // 确保空时传递空字符串
    });
    closePaidModal();
    getOrderList();
    Notification.success("补单成功");
  } catch (error) {
    Notification.error(error);
  }
};
const onBatchDelete = async () => {
  try {
    await delOrderApi({ ids: selectedKeys.value });
    pagination.value.current = 1;
    getOrderList();
    Notification.success("删除成功");
    selectedKeys.value = [];
  } catch (error) {
    Notification.error(error);
  }
};
const handleManualNotify = (record: List) => {
  Modal.confirm({
    title: "确认手动回调",
    content: `确定要手动触发订单 ${record.order_id} 的回调吗？系统将立即向商户发送回调通知。`,
    okText: "确认回调",
    cancelText: "取消",
    onOk: async () => {
      try {
        await manualNotifyAPI({ id: record.id });
        Notification.success("回调已触发，请稍后查看回调状态");
        getOrderList();
      } catch (error) {
        Notification.error(error);
      }
    }
  });
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

// 在 style 标签中添加
.paid-modal-content {
  padding: 4px 0;

  :deep(.arco-alert) {
    border-radius: 6px;
  }

  :deep(.arco-form-item-label-col) {
    font-weight: 500;
    color: #1d2129;
  }

  :deep(.arco-input-wrapper) {
    &:hover {
      border-color: #4080ff;
    }
  }
}

:deep(.arco-modal) {
  .arco-modal-header {
    border-bottom: 1px solid #e5e6eb;
    padding: 16px 20px;
  }

  .arco-modal-body {
    padding: 20px;
  }

  .arco-modal-footer {
    border-top: 1px solid #e5e6eb;
    padding: 12px 20px;
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
