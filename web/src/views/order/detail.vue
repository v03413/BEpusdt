<template>
  <a-modal
    width="680px"
    :visible="visible"
    @close="onClose"
    @cancel="onClose"
    @update:visible="onClose"
    :footer="false"
    unmount-on-close
  >
    <template #title>
      <div class="detail-modal-title">
        <icon-star />
        <span>订单详情</span>
      </div>
    </template>
    <div class="detail-content">
      <!-- 基础信息卡片 -->
      <a-card class="detail-card" title="基础信息" :bordered="false">
        <template #extra>
          <a-tag size="medium" :color="getStatusColor(detailData.status)" class="status-tag">
            <icon-check-circle v-if="detailData.status === 2" />
            <icon-clock-circle v-else-if="detailData.status === 1 || detailData.status === 5" />
            <icon-close-circle v-else />
            {{ getStatusText(detailData.status) }}
          </a-tag>
        </template>
        <a-row :gutter="24">
          <a-col :span="12">
            <div class="detail-item">
              <div class="detail-label">
                <icon-idcard />
                <span>订单ID</span>
              </div>
              <div class="detail-value">{{ detailData.id }}</div>
            </div>
          </a-col>
          <a-col :span="12">
            <div class="detail-item">
              <div class="detail-label">
                <icon-file />
                <span>商户订单号</span>
              </div>
              <div class="detail-value">{{ detailData.order_id }}</div>
            </div>
          </a-col>
        </a-row>
        <a-row :gutter="24">
          <a-col :span="12">
            <div class="detail-item">
              <div class="detail-label">
                <icon-tag />
                <span>本地交易ID</span>
              </div>
              <div class="detail-value">{{ detailData.trade_id }}</div>
            </div>
          </a-col>
          <a-col :span="12">
            <div class="detail-item">
              <div class="detail-label">
                <icon-user />
                <span>商品名称</span>
              </div>
              <div class="detail-value">{{ detailData.name }}</div>
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
                <icon-code />
                <span>API类型</span>
              </div>
              <div class="detail-value">
                <a-tag color="arcoblue">{{ detailData.api_type }}</a-tag>
              </div>
            </div>
          </a-col>
        </a-row>
      </a-card>

      <!-- 金额信息卡片 -->
      <a-card class="detail-card" title="金额信息" :bordered="false">
        <a-row :gutter="24">
          <a-col :span="12">
            <div class="detail-item">
              <div class="detail-label">
                <icon-computer />
                <span>订单汇率</span>
              </div>
              <div class="detail-value money-value">{{ detailData.rate }}</div>
            </div>
          </a-col>
          <a-col :span="12">
            <div class="detail-item">
              <div class="detail-label">
                <icon-pushpin />
                <span>实际付款</span>
              </div>
              <div class="detail-value money-value">{{ detailData.amount }}</div>
            </div>
          </a-col>
        </a-row>
        <a-row :gutter="24">
          <a-col :span="12">
            <div class="detail-item">
              <div class="detail-label">
                <icon-archive />
                <span>交易金额</span>
              </div>
              <div class="detail-value money-value highlight">{{ getCurrencySymbol(detailData.fiat) }}{{ detailData.money }}</div>
            </div>
          </a-col>
        </a-row>
      </a-card>

      <!-- 地址信息卡片 -->
      <a-card class="detail-card" title="地址信息" :bordered="false">
        <a-row :gutter="24">
          <a-col :span="24">
            <div class="detail-item">
              <div class="detail-label">
                <icon-location />
                <span>收款地址</span>
              </div>
              <div class="detail-value address-value">
                <a-typography-text copyable>{{ detailData.address }}</a-typography-text>
              </div>
            </div>
          </a-col>
        </a-row>
        <a-row :gutter="24" v-if="detailData.from_address">
          <a-col :span="24">
            <div class="detail-item">
              <div class="detail-label">
                <icon-send />
                <span>支付地址</span>
              </div>
              <div class="detail-value address-value">
                <a-typography-text copyable>{{ detailData.from_address }}</a-typography-text>
              </div>
            </div>
          </a-col>
        </a-row>
      </a-card>

      <!-- 回调信息卡片 -->
      <a-card class="detail-card" title="回调信息" :bordered="false">
        <a-row :gutter="24">
          <a-col :span="12">
            <div class="detail-item">
              <div class="detail-label">
                <icon-notification />
                <span>回调次数</span>
              </div>
              <div class="detail-value">{{ detailData.notify_num }}</div>
            </div>
          </a-col>
          <a-col :span="12">
            <div class="detail-item">
              <div class="detail-label">
                <icon-check />
                <span>回调状态</span>
              </div>
              <div class="detail-value">
                <a-tag :color="detailData.notify_state === 1 ? 'green' : 'red'">
                  {{ detailData.notify_state === 1 ? "成功" : "失败" }}
                </a-tag>
              </div>
            </div>
          </a-col>
        </a-row>
        <a-row :gutter="24" v-if="detailData.return_url">
          <a-col :span="24">
            <div class="detail-item">
              <div class="detail-label">
                <icon-link />
                <span>同步地址</span>
              </div>
              <div class="detail-value url-value">{{ detailData.return_url }}</div>
            </div>
          </a-col>
        </a-row>
        <a-row :gutter="24" v-if="detailData.notify_url">
          <a-col :span="24">
            <div class="detail-item">
              <div class="detail-label">
                <icon-sync />
                <span>异步地址</span>
              </div>
              <div class="detail-value url-value">{{ detailData.notify_url }}</div>
            </div>
          </a-col>
        </a-row>
      </a-card>

      <!-- 区块链信息卡片 -->
      <a-card class="detail-card" title="区块链信息" :bordered="false" v-if="detailData.ref_hash || detailData.ref_block_num">
        <a-row :gutter="24" v-if="detailData.ref_hash">
          <a-col :span="24">
            <div class="detail-item">
              <div class="detail-label">
                <icon-safe />
                <span>交易哈希</span>
              </div>
              <div class="detail-value hash-value">
                <a-typography-text copyable>{{ detailData.ref_hash }}</a-typography-text>
              </div>
            </div>
          </a-col>
        </a-row>
        <a-row :gutter="24" v-if="detailData.ref_block_num">
          <a-col :span="12">
            <div class="detail-item">
              <div class="detail-label">
                <icon-layers />
                <span>区块索引</span>
              </div>
              <div class="detail-value">{{ detailData.ref_block_num }}</div>
            </div>
          </a-col>
        </a-row>
        <a-row :gutter="24">
          <a-col :span="24">
            <div class="detail-item">
              <a-button
                size="small"
                type="primary"
                :disabled="!(detailData.status === 2 || detailData.status === 5)"
                @click="openTxUrl"
              >
                <template #icon>
                  <icon-export />
                </template>
                链上交易详情
              </a-button>
            </div>
          </a-col>
        </a-row>
      </a-card>

      <!-- 时间信息卡片 -->
      <a-card class="detail-card" title="时间信息" :bordered="false">
        <a-row :gutter="24">
          <a-col :span="12" v-if="detailData.created_at">
            <div class="detail-item">
              <div class="detail-label">
                <icon-plus-circle />
                <span>创建时间</span>
              </div>
              <div class="detail-value">{{ formatDateTime(detailData.created_at) }}</div>
            </div>
          </a-col>
          <a-col :span="12" v-if="detailData.updated_at">
            <div class="detail-item">
              <div class="detail-label">
                <icon-edit />
                <span>更新时间</span>
              </div>
              <div class="detail-value">{{ formatDateTime(detailData.updated_at) }}</div>
            </div>
          </a-col>
        </a-row>
        <a-row :gutter="24">
          <a-col :span="12">
            <div class="detail-item">
              <div class="detail-label">
                <icon-schedule />
                <span>失效时间</span>
              </div>
              <div class="detail-value">{{ formatDateTime(detailData.expired_at) }}</div>
            </div>
          </a-col>
          <a-col :span="12" v-if="detailData.confirmed_at">
            <div class="detail-item">
              <div class="detail-label">
                <icon-check-circle />
                <span>确认时间</span>
              </div>
              <div class="detail-value">{{ formatDateTime(detailData.confirmed_at) }}</div>
            </div>
          </a-col>
        </a-row>
      </a-card>
    </div>
  </a-modal>
</template>

<script setup lang="ts">
const props = defineProps({
  visible: Boolean,
  detailData: {
    type: Object,
    required: true
  }
});
const emits = defineEmits(["close"]);
const onClose = () => {
  emits("close");
};

// 打开链上交易详情
const openTxUrl = () => {
  if (props.detailData.tx_url) {
    window.open(props.detailData.tx_url, "_blank");
  }
};

// 获取状态颜色
const getStatusColor = (status: number) => {
  const statusMap: Record<number, string> = {
    1: "blue", // 等待支付
    2: "green", // 交易成功
    3: "gray", // 交易过期
    4: "gold", // 交易取消
    5: "pinkpurple", // 等待确认
    6: "red" // 确认失败
  };
  return statusMap[status] || "gray";
};

// 获取状态文本
const getStatusText = (status: number) => {
  const statusMap: Record<number, string> = {
    1: "等待支付",
    2: "交易成功",
    3: "交易过期",
    4: "交易取消",
    5: "等待确认",
    6: "确认失败"
  };
  return statusMap[status] || "未知状态";
};

// 格式化时间
const formatDateTime = (dateTimeStr: string) => {
  if (!dateTimeStr) return "";

  const date = new Date(dateTimeStr);
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, "0");
  const day = String(date.getDate()).padStart(2, "0");
  const hours = String(date.getHours()).padStart(2, "0");
  const minutes = String(date.getMinutes()).padStart(2, "0");
  const seconds = String(date.getSeconds()).padStart(2, "0");

  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
};

// 获取货币符号
const getCurrencySymbol = (fiat: string) => {
  const currencyMap: Record<string, string> = {
    CNY: "¥",
    USD: "$",
    JPY: "¥",
    GBP: "£",
    EUR: "€"
  };
  return currencyMap[fiat] || "";
};
</script>

<style scoped>
.detail-modal-title {
  display: flex;
  align-items: center;
  gap: 8px;
}

.detail-content {
  padding: 8px 0;
  max-height: 70vh;
  overflow-y: auto;
}

.detail-card {
  margin-bottom: 16px;
}

.detail-card:last-child {
  margin-bottom: 0;
}

.status-tag {
  display: flex;
  align-items: center;
  gap: 4px;
  font-weight: 500;
}

.detail-item {
  margin-bottom: 20px;
  display: flex;
  flex-direction: column;
}

.detail-item:last-child {
  margin-bottom: 0;
}

.detail-label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 500;
  color: var(--color-text-2);
  margin-bottom: 6px;
  font-size: 14px;
}

.detail-value {
  font-size: 14px;
  color: var(--color-text-1);
  line-height: 1.5;
  padding-left: 26px;
  position: relative;
}

.address-value {
  font-family: "Monaco", "Menlo", "Consolas", monospace;
  font-size: 13px;
  word-break: break-all;
}

.address-value :deep(.arco-typography-operation-copy) {
  color: #165dff;
  margin-left: 4px;
}

.address-value :deep(.arco-typography-operation-copy:hover) {
  color: #0e42d2;
}

.hash-value {
  font-family: "Monaco", "Menlo", "Consolas", monospace;
  font-size: 12px;
  word-break: break-all;
  color: var(--color-text-3);
}

.hash-value :deep(.arco-typography-operation-copy) {
  color: #165dff;
  margin-left: 4px;
}

.money-value {
  font-weight: 600;
  color: var(--color-text-1);
}

.money-value.highlight {
  font-size: 16px;
  color: #f53f3f;
  font-weight: 700;
}

.url-value {
  font-size: 12px;
  color: var(--color-text-3);
  word-break: break-all;
  line-height: 1.4;
}

/* 响应式设计 */
@media (max-width: 768px) {
  :deep(.arco-modal) {
    width: 95vw !important;
    margin: 10px;
  }

  .detail-content {
    max-height: 80vh;
  }
}
</style>
