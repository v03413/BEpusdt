<template>
  <div class="wallet">
    <div class="snow-page">
      <div class="snow-inner">
        <a-form ref="formRef" auto-label-width :model="formData.form">
          <a-row :gutter="16">
            <a-col :xs="24" :sm="24" :md="12" :lg="12" :xl="6" :xxl="6">
              <a-input v-model="formData.form.name" placeholder="请输入名称" allow-clear />
            </a-col>
            <a-col :xs="24" :sm="24" :md="12" :lg="12" :xl="6" :xxl="6">
              <a-input v-model="formData.form.address" placeholder="请输入钱包地址" allow-clear />
            </a-col>
            <a-col :xs="24" :sm="24" :md="12" :lg="12" :xl="6" :xxl="6">
              <a-select v-model="formData.form.trade_type" placeholder="请选择交易类型" allow-clear allow-search>
                <a-option v-for="item in tradeTypeOptions" :key="item.value" :value="item.value">
                  {{ item.label }}
                </a-option>
              </a-select>
            </a-col>
            <a-col :xs="24" :sm="24" :md="12" :lg="12" :xl="6" :xxl="3">
              <a-space class="search-btn">
                <a-button type="primary" @click="getCommonTableList">
                  <template #icon><icon-search /></template>
                  查询
                </a-button>
                <a-button @click="onReset">
                  <template #icon><icon-refresh /></template>
                  重置
                </a-button>
              </a-space>
            </a-col>
            <a-col :xs="24" :sm="24" :md="12" :lg="12" :xl="6" :xxl="3">
              <a-button type="primary" status="success" @click="onAdd">
                <template #icon><icon-plus /></template>
                新增钱包
              </a-button>
            </a-col>
          </a-row>
        </a-form>

        <a-table
          row-key="key"
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
          <template #address="{ record }">
            <div class="address-cell">
              <a-typography-text copyable class="address-text">
                {{ record.address }}
              </a-typography-text>
            </div>
          </template>

          <template #status="{ record }">
            <a-tag size="small" :color="record.status === 1 ? 'green' : 'red'">
              {{ record.status === 1 ? "启用" : "停用" }}
            </a-tag>
          </template>

          <template #other_notify="{ record }">
            <a-tag size="small" :color="record.other_notify === 1 ? 'arcoblue' : 'gray'">
              {{ record.other_notify === 1 ? "开启" : "关闭" }}
            </a-tag>
          </template>

          <template #optional="{ record }">
            <a-space>
              <a-button size="mini" type="primary" @click="showDetail(record)">详情</a-button>
              <a-button size="mini" @click="onMod(record)">修改</a-button>
              <a-popconfirm content="确定删除这条数据吗?" type="warning" @ok="onDelete(record)">
                <a-button size="mini" type="primary" status="danger">删除</a-button>
              </a-popconfirm>
            </a-space>
          </template>
        </a-table>
      </div>
    </div>

    <!-- 新增钱包对话框 -->
    <a-modal width="40%" v-model:visible="open" @close="afterClose" @ok="addWallet" @cancel="afterClose">
      <template #title>{{ title }}</template>
      <a-form ref="formRef" auto-label-width :rules="rules" :model="addFrom">
        <a-form-item field="name" label="钱包名称" validate-trigger="blur">
          <a-input v-model="addFrom.name" placeholder="请输入钱包名称" allow-clear />
        </a-form-item>
        <a-form-item field="address" label="钱包地址" validate-trigger="blur">
          <a-input v-model="addFrom.address" placeholder="请输入钱包地址" allow-clear />
        </a-form-item>
        <a-form-item field="trade_type" label="交易类型" :rules="[{ required: true, message: '交易类型不能为空' }]">
          <a-select v-model="addFrom.trade_type" placeholder="请选择" allow-clear allow-search>
            <a-option v-for="item in tradeTypeOptions" :key="item.value" :value="item.value">
              {{ item.label }}
            </a-option>
          </a-select>
        </a-form-item>
        <a-form-item field="other_notify" label="其他通知">
          <a-select v-model="addFrom.other_notify" placeholder="请选择" allow-clear>
            <a-option :value="0">关闭</a-option>
            <a-option :value="1">启用</a-option>
          </a-select>
        </a-form-item>
        <a-form-item field="remark" label="备注信息" validate-trigger="blur">
          <a-textarea v-model="addFrom.remark" placeholder="请输入备注信息" allow-clear />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 修改钱包对话框 -->
    <a-modal width="40%" v-model:visible="modOpen" @close="afterModClose" @ok="modWallet" @cancel="afterModClose">
      <template #title>{{ modTitle }}</template>
      <a-form ref="modFormRef" auto-label-width :rules="rules" :model="modFrom">
        <a-form-item field="name" label="钱包名称" validate-trigger="blur">
          <a-input v-model="modFrom.name" placeholder="请输入钱包名称" allow-clear />
        </a-form-item>
        <a-form-item field="address" label="钱包地址" validate-trigger="blur">
          <a-input v-model="modFrom.address" placeholder="请输入钱包地址" allow-clear />
        </a-form-item>
        <a-form-item field="trade_type" label="交易类型" :rules="[{ required: true, message: '交易类型不能为空' }]">
          <a-select v-model="modFrom.trade_type" placeholder="请选择" allow-clear allow-search>
            <a-option v-for="item in tradeTypeOptions" :key="item.value" :value="item.value">
              {{ item.label }}
            </a-option>
          </a-select>
        </a-form-item>
        <a-form-item field="status" label="收款状态">
          <a-select v-model="modFrom.status" placeholder="请选择" allow-clear>
            <a-option :value="1">启用</a-option>
            <a-option :value="0">停用</a-option>
          </a-select>
        </a-form-item>
        <a-form-item field="other_notify" label="其他通知">
          <a-select v-model="modFrom.other_notify" placeholder="请选择" allow-clear>
            <a-option :value="0">关闭</a-option>
            <a-option :value="1">开启</a-option>
          </a-select>
        </a-form-item>
        <a-form-item field="remark" label="备注信息" validate-trigger="blur">
          <a-textarea v-model="modFrom.remark" placeholder="请输入备注信息" allow-clear />
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 详情对话框 -->
    <a-modal
      width="680px"
      v-model:visible="detailVisible"
      @close="closeDetail"
      @cancel="closeDetail"
      :footer="false"
      unmount-on-close
    >
      <template #title>
        <div class="detail-modal-title">
          <icon-star />
          <span>钱包详情信息</span>
        </div>
      </template>

      <div class="detail-content">
        <a-card class="detail-card" title="基础信息" :bordered="false">
          <template #extra>
            <a-tag size="medium" :color="detailData.status === 1 ? 'green' : 'red'" class="status-tag">
              <icon-check-circle v-if="detailData.status === 1" />
              <icon-close-circle v-else />
              {{ detailData.status === 1 ? "启用中" : "已停用" }}
            </a-tag>
          </template>

          <a-row :gutter="24">
            <a-col :span="12">
              <div class="detail-item">
                <div class="detail-label">
                  <icon-idcard />
                  <span>钱包ID</span>
                </div>
                <div class="detail-value">{{ detailData.id }}</div>
              </div>
            </a-col>
            <a-col :span="12">
              <div class="detail-item">
                <div class="detail-label">
                  <icon-user />
                  <span>钱包名称</span>
                </div>
                <div class="detail-value">{{ detailData.name }}</div>
              </div>
            </a-col>
          </a-row>

          <a-row :gutter="24">
            <a-col :span="24">
              <div class="detail-item">
                <div class="detail-label">
                  <icon-location />
                  <span>钱包地址</span>
                </div>
                <div class="detail-value address-value">
                  <a-typography-text copyable>{{ detailData.address }}</a-typography-text>
                </div>
              </div>
            </a-col>
          </a-row>

          <a-row :gutter="24">
            <a-col :span="12">
              <div class="detail-item">
                <div class="detail-label">
                  <icon-swap />
                  <span>交易类型</span>
                </div>
                <div class="detail-value">
                  <a-tag color="blue">{{ detailData.trade_type }}</a-tag>
                </div>
              </div>
            </a-col>
            <a-col :span="12">
              <div class="detail-item">
                <div class="detail-label">
                  <icon-notification />
                  <span>监控状态</span>
                </div>
                <div class="detail-value">
                  <a-tag :color="detailData.other_notify === 1 ? 'arcoblue' : 'gray'">
                    <icon-eye v-if="detailData.other_notify === 1" />
                    <icon-eye-invisible v-else />
                    {{ detailData.other_notify === 1 ? "已开启" : "已关闭" }}
                  </a-tag>
                </div>
              </div>
            </a-col>
          </a-row>
        </a-card>

        <a-card class="detail-card" title="备注信息" :bordered="false" v-if="detailData.remark">
          <div class="remark-content">
            <icon-message />
            <span>{{ detailData.remark }}</span>
          </div>
        </a-card>

        <a-card class="detail-card" title="时间信息" :bordered="false" v-if="detailData.created_at || detailData.updated_at">
          <a-row :gutter="24">
            <a-col :span="12" v-if="detailData.created_at">
              <div class="detail-item">
                <div class="detail-label">
                  <icon-plus-circle />
                  <span>创建时间</span>
                </div>
                <div class="detail-value">{{ detailData.created_at }}</div>
              </div>
            </a-col>
            <a-col :span="12" v-if="detailData.updated_at">
              <div class="detail-item">
                <div class="detail-label">
                  <icon-edit />
                  <span>更新时间</span>
                </div>
                <div class="detail-value">{{ detailData.updated_at }}</div>
              </div>
            </a-col>
          </a-row>
        </a-card>
      </div>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { getWalletListAPI, delWalletAPI, addWalletAPI, modWalletAPI } from "@/api/modules/wallet/index";
import { List, FormData, Pagination, AddForm, ModForm } from "./config";
import { Notification } from "@arco-design/web-vue";
import { useUserInfoStore } from "@/store/modules/user-info";
import { useWalletDetail } from "./detail";

const userStores = useUserInfoStore();
const { detailVisible, detailData, showDetail, closeDetail } = useWalletDetail();

const tradeTypeOptions = computed(() => Object.entries(userStores.trade_type).map(([value, label]) => ({ value, label })));

const formData = reactive<FormData>({
  form: { name: "", trade_type: "", address: "" },
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
  { title: "名称", align: "center", dataIndex: "name", width: 200 },
  { title: "交易类型", align: "center", dataIndex: "trade_type", width: 120 },
  { title: "钱包地址", align: "center", dataIndex: "address", slotName: "address", width: 300, ellipsis: true },
  { title: "收款状态", dataIndex: "status", align: "center", slotName: "status", width: 100 },
  { title: "其它通知", dataIndex: "other_notify", align: "center", slotName: "other_notify", width: 100 },
  { title: "操作", slotName: "optional", align: "center", fixed: "right", width: 200 }
];

const rules = {
  name: [{ required: true, message: "请输入钱包名称" }],
  address: [{ required: true, message: "请输入钱包地址" }],
  trade_type: [{ required: true, message: "请输入交易类型" }]
};

const formRef = ref();
const modFormRef = ref();
const title = ref("");
const modTitle = ref("");
const open = ref(false);
const modOpen = ref(false);

const addFrom = ref<AddForm>({
  name: "",
  address: "",
  trade_type: "",
  remark: "",
  other_notify: 0
});

const modFrom = ref<ModForm>({
  id: 0,
  name: "",
  address: "",
  trade_type: "",
  remark: "",
  other_notify: 0,
  status: 1
});

const pageChange = (page: number) => {
  pagination.value.current = page;
  getCommonTableList();
};

const pageSizeChange = (pageSize: number) => {
  pagination.value.pageSize = pageSize;
  getCommonTableList();
};

const onReset = () => {
  formData.form = { name: "", trade_type: "", address: "" };
  getCommonTableList();
};

const getCommonTableList = async () => {
  try {
    loading.value = true;
    const res = await getWalletListAPI({
      page: pagination.value.current,
      size: pagination.value.pageSize,
      sort: "desc",
      keyword: "",
      name: formData.form.name,
      trade_type: formData.form.trade_type,
      address: formData.form.address,
      status: 99
    });

    data.length = 0;
    data.push(...res.data);
    pagination.value.total = res.total;
  } finally {
    loading.value = false;
  }
};

const onDelete = async (record: List) => {
  try {
    await delWalletAPI({ id: record.id });
    getCommonTableList();
    Notification.success("删除成功");
  } catch (error) {
    Notification.error(error);
  }
};

const onAdd = () => {
  title.value = "新增钱包";
  open.value = true;
};

const onMod = (record: List) => {
  modTitle.value = "修改钱包";
  modFrom.value = {
    id: record.id,
    name: record.name,
    address: record.address,
    trade_type: record.trade_type || "",
    remark: record.remark || "",
    other_notify: record.other_notify || 0,
    status: record.status
  };
  modOpen.value = true;
};

const afterClose = () => {
  formRef.value.resetFields();
  addFrom.value = {
    name: "",
    address: "",
    trade_type: "",
    remark: "",
    other_notify: 0
  };
};

const afterModClose = () => {
  modFormRef.value?.resetFields();
  modFrom.value = {
    id: 0,
    name: "",
    address: "",
    trade_type: "",
    remark: "",
    other_notify: 0,
    status: 1
  };
};

const addWallet = async () => {
  const state = await formRef.value.validate();
  if (state) return;

  try {
    await addWalletAPI(addFrom.value);
    open.value = false;
    getCommonTableList();
    Notification.success("添加成功");
  } catch (error) {
    Notification.error(error);
  }
};

const modWallet = async () => {
  const state = await modFormRef.value.validate();
  if (state) return;

  try {
    await modWalletAPI(modFrom.value);
    modOpen.value = false;
    getCommonTableList();
    Notification.success("修改成功");
  } catch (error) {
    Notification.error(error);
  }
};

getCommonTableList();
</script>

<style lang="scss" scoped>
.wallet {
  height: 100%;

  .snow-page {
    height: 100%;

    .snow-inner {
      height: 100%;
      display: flex;
      flex-direction: column;

      .arco-table-container {
        flex: 1;
        overflow: hidden;
      }
    }
  }
}

.search-btn {
  margin-bottom: 20px;
}

.address-cell {
  max-width: 250px;

  .address-text {
    font-family: "Monaco", "Menlo", "Consolas", monospace;
    font-size: 12px;
    word-break: break-all;
    line-height: 1.4;

    :deep(.arco-typography-operation-copy) {
      color: #165dff;
      margin-left: 4px;

      &:hover {
        color: #0e42d2;
      }
    }
  }
}

.detail-modal-title {
  display: flex;
  align-items: center;
  gap: 8px;

  span {
    font-weight: 600;
    font-size: 16px;
  }
}

.detail-content {
  padding: 8px 0;

  .detail-card {
    margin-bottom: 16px;
    box-shadow: 0 1px 4px rgba(0, 0, 0, 0.08);

    &:last-child {
      margin-bottom: 0;
    }

    :deep(.arco-card-header) {
      border-bottom: 1px solid var(--color-border-2);
      padding: 16px 20px 12px;

      .arco-card-header-title {
        font-weight: 600;
        color: var(--color-text-1);
      }
    }

    :deep(.arco-card-body) {
      padding: 20px;
    }
  }

  .detail-item {
    margin-bottom: 20px;

    &:last-child {
      margin-bottom: 0;
    }

    .detail-label {
      display: flex;
      align-items: center;
      gap: 6px;
      margin-bottom: 8px;
      font-size: 13px;
      color: var(--color-text-3);
      font-weight: 500;

      .arco-icon {
        font-size: 14px;
        color: var(--color-text-4);
      }
    }

    .detail-value {
      font-size: 14px;
      color: var(--color-text-1);
      font-weight: 500;
      min-height: 22px;
      display: flex;
      align-items: center;

      &.address-value {
        word-break: break-all;
        font-family: "Monaco", "Menlo", monospace;
        font-size: 13px;

        :deep(.arco-typography) {
          font-family: inherit;
          font-size: inherit;
        }
      }
    }
  }

  .status-tag {
    display: flex;
    align-items: center;
    gap: 4px;
    font-weight: 500;

    .arco-icon {
      font-size: 12px;
    }
  }

  .remark-content {
    display: flex;
    align-items: flex-start;
    gap: 8px;
    padding: 12px 16px;
    background-color: var(--color-fill-2);
    border-radius: 6px;
    line-height: 1.6;

    .arco-icon {
      margin-top: 2px;
      color: var(--color-text-3);
      flex-shrink: 0;
    }

    span {
      color: var(--color-text-2);
    }
  }
}

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

  .detail-content {
    .detail-card :deep(.arco-card-body) {
      padding: 16px;
    }

    .detail-item .detail-value.address-value {
      font-size: 12px;
    }
  }
}
</style>
