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
                <icon-user />
                <span>商品名称</span>
              </div>
              <div class="detail-value">{{ detailData.name }}</div>
            </div>
          </a-col>
          <a-col :span="12">
            <div class="detail-item">
              <div class="detail-label">
                <icon-swap />
                <span>交易类型</span>
              </div>
              <div class="detail-value">
                <a-tag color="blue" class="trade-type-tag">{{ detailData.trade_type }}</a-tag>
              </div>
            </div>
          </a-col>
        </a-row>
        <a-row :gutter="24">
          <a-col :span="12">
            <div class="detail-item">
              <div class="detail-label">
                <icon-file />
                <span>商户订单</span>
              </div>
              <div class="detail-value">
                <a-typography-text copyable>{{ detailData.order_id }}</a-typography-text>
              </div>
            </div>
          </a-col>
          <a-col :span="12">
            <div class="detail-item">
              <div class="detail-label">
                <icon-tag />
                <span>交易编号</span>
              </div>
              <div class="detail-value">
                <a-typography-text copyable>{{ detailData.trade_id }}</a-typography-text>
              </div>
            </div>
          </a-col>
        </a-row>
        <a-row :gutter="24">
          <a-col :span="12">
            <div class="detail-item">
              <div class="detail-label">
                <icon-archive />
                <span>交易金额（汇率）</span>
              </div>
              <div class="detail-value">
                <span class="currency-symbol">{{ getCurrencySymbol(detailData.fiat) }}</span
                >{{ detailData.money }}
                <span class="rate-text">({{ detailData.rate }})</span>
              </div>
            </div>
          </a-col>
          <a-col :span="12">
            <div class="detail-item">
              <div class="detail-label">
                <icon-pushpin />
                <span>交易数额</span>
              </div>
              <div class="detail-value">
                {{ detailData.amount }}
                <a-tag size="mini" :color="getCryptoColor(detailData.crypto)" bordered style="margin-left: 4px">
                  {{ detailData.crypto }}
                </a-tag>
              </div>
            </div>
          </a-col>
        </a-row>
      </a-card>

      <!-- 地址信息卡片 -->
      <a-card class="detail-card" title="交易地址" :bordered="false">
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
      <a-card class="detail-card" title="回调信息" :bordered="false" v-if="detailData.status === 2 || detailData.status === 5">
        <a-row :gutter="24">
          <a-col :span="12">
            <div class="detail-item">
              <div class="detail-label">
                <icon-check />
                <span>回调状态</span>
              </div>
              <div class="detail-value">
                <a-tag v-if="detailData.notify_state === 1" color="green"> 成功 </a-tag>
                <a-tag v-else color="red"> 失败，等待第 {{ detailData.notify_num + 1 }} 次回调中 </a-tag>
              </div>
            </div>
          </a-col>
          <a-col :span="12" v-if="detailData.return_url">
            <div class="detail-item">
              <div class="detail-label">
                <icon-link />
                <span>商户网站</span>
              </div>
              <div class="detail-value">
                <a-link @click="openMerchantWebsite" :hoverable="false">
                  {{ getMerchantWebsite(detailData.return_url) }}
                </a-link>
              </div>
            </div>
          </a-col>
        </a-row>
      </a-card>

      <!-- 区块链信息卡片 -->
      <a-card class="detail-card" title="区块链数据" :bordered="false" v-if="detailData.status === 2 || detailData.status === 5">
        <a-row :gutter="24" v-if="detailData.ref_hash">
          <a-col :span="12" v-if="detailData.ref_block_num">
            <div class="detail-item">
              <div class="detail-label">
                <icon-layers />
                <span>区块索引</span>
              </div>
              <div class="detail-value">{{ detailData.ref_block_num }}</div>
            </div>
          </a-col>
          <a-col :span="12">
            <div class="detail-item">
              <div class="detail-label">
                <icon-safe />
                <span>链上详情</span>
              </div>
              <div class="detail-value hash-value">
                <a-link
                  v-if="detailData.status === 2 && detailData.tx_url"
                  @click="openTxUrl"
                  :hoverable="false"
                  class="tx-url-link"
                >
                  {{ detailData.tx_url }}
                </a-link>
                <a-tag v-else color="blue" size="small">
                  <template #icon><icon-clock-circle /></template>
                  等待交易确认
                </a-tag>
              </div>
            </div>
          </a-col>
        </a-row>
      </a-card>

      <!-- 时间信息卡片 -->
      <a-card class="detail-card" title="订单时间" :bordered="false">
        <a-row :gutter="24">
          <a-col :span="12" v-if="detailData.created_at">
            <div class="detail-item">
              <div class="detail-label">
                <icon-plus-circle />
                <span>创建订单</span>
              </div>
              <div class="detail-value">{{ formatDateTime(detailData.created_at) }}</div>
            </div>
          </a-col>
          <a-col :span="12">
            <div class="detail-item">
              <div class="detail-label">
                <icon-check-circle v-if="detailData.confirmed_at && (detailData.status === 2 || detailData.status === 5)" />
                <icon-schedule v-else-if="detailData.status === 3" />
                <icon-sync v-else />
                <span v-if="detailData.confirmed_at && (detailData.status === 2 || detailData.status === 5)">交易确认</span>
                <span v-else-if="detailData.status === 3">交易过期</span>
                <span v-else>最后更新</span>
              </div>
              <div class="detail-value">
                <span v-if="detailData.confirmed_at && (detailData.status === 2 || detailData.status === 5)">
                  {{ formatDateTime(detailData.confirmed_at) }}
                </span>
                <span v-else-if="detailData.status === 3">
                  {{ formatDateTime(detailData.expired_at) }}
                </span>
                <span v-else>
                  {{ formatDateTime(detailData.updated_at) }}
                </span>
              </div>
            </div>
          </a-col>
        </a-row>
      </a-card>
    </div>
  </a-modal>
</template>

<script setup lang="ts">
import { getCryptoColor } from "@/views/rate/common";

const props = defineProps({
  visible: Boolean,
  detailData: {
    type: Object,
    required: true
  }
});

const emits = defineEmits(["close"]);

const onClose = () => emits("close");

const openTxUrl = () => {
  if (props.detailData.tx_url) {
    window.open(props.detailData.tx_url, "_blank");
  }
};

const getMerchantWebsite = (returnUrl: string) => {
  if (!returnUrl) return "";
  try {
    const url = new URL(returnUrl);
    return `${url.protocol}//${url.host}/`;
  } catch {
    return returnUrl;
  }
};

const openMerchantWebsite = () => {
  if (props.detailData.return_url) {
    const merchantUrl = getMerchantWebsite(props.detailData.return_url);
    window.open(merchantUrl, "_blank");
  }
};

const statusMap: Record<number, { color: string; text: string }> = {
  1: { color: "blue", text: "等待支付" },
  2: { color: "green", text: "交易成功" },
  3: { color: "gray", text: "交易过期" },
  4: { color: "gold", text: "交易取消" },
  5: { color: "pinkpurple", text: "等待确认" },
  6: { color: "red", text: "确认失败" }
};

const getStatusColor = (status: number) => statusMap[status]?.color || "gray";
const getStatusText = (status: number) => statusMap[status]?.text || "未知状态";

const formatDateTime = (dateTimeStr: string) => {
  if (!dateTimeStr) return "";
  const date = new Date(dateTimeStr);
  return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, "0")}-${String(date.getDate()).padStart(2, "0")} ${String(date.getHours()).padStart(2, "0")}:${String(date.getMinutes()).padStart(2, "0")}:${String(date.getSeconds()).padStart(2, "0")}`;
};

const currencySymbolMap: Record<string, string> = {
  CNY: "¥",
  USD: "$",
  JPY: "¥",
  GBP: "£",
  EUR: "€"
};

const getCurrencySymbol = (fiat: string) => currencySymbolMap[fiat] || "";
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
}

.detail-value :deep(.arco-typography) {
  font-size: 14px;
  color: var(--color-text-1);
}

.detail-value :deep(.arco-typography-operation-copy) {
  color: #165dff;
  margin-left: 4px;
}

.detail-value :deep(.arco-typography-operation-copy:hover) {
  color: #0e42d2;
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
  color: var(--color-text-3);
  overflow: hidden;
}

.tx-url-link {
  display: inline-block;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 100%;
  vertical-align: bottom;
}

.hash-value :deep(.arco-link) {
  display: inline-block;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 100%;
  vertical-align: bottom;
}

.currency-symbol {
  font-weight: 700;
  font-size: 15px;
}

.rate-text {
  color: var(--color-text-3);
  font-size: 13px;
  margin-left: 4px;
}

.trade-type-tag {
  font-weight: 700;
  font-size: 14px;
}

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
