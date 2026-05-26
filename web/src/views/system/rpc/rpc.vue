<template>
  <div class="snow-page">
    <a-spin :loading="loading" tip="正在加载..." class="full-height">
      <!-- 网络状态总览 -->
      <a-card :bordered="false" class="stats-card">
        <template #title>
          <div class="card-title">
            <div class="title-icon stats-icon">
              <icon-thunderbolt />
            </div>
            <span>网络状态总览</span>
            <a-tag size="small" color="arcoblue" class="live-tag">
              <icon-sync class="spin-icon" />
              实时
            </a-tag>
          </div>
        </template>

        <div class="network-stats-grid">
          <div v-for="net in allNetworks" :key="net.key" class="network-stat-card" :class="getNetworkClass(net.key)">
            <!-- 网络头部 -->
            <div class="net-card-header">
              <div class="net-badge" :style="{ background: net.color }">
                <component :is="net.icon" />
              </div>
              <div class="net-name">{{ net.label }}</div>
              <div class="net-status-dot" :class="getStatusDotClass(net.statKey)"></div>
            </div>

            <!-- 统计数据 -->
            <div class="net-card-body">
              <!-- 成功率 -->
              <div class="stat-row">
                <div class="stat-label">
                  <icon-check-circle class="stat-icon success-icon" />
                  成功率
                </div>
                <div class="stat-value" :class="getSuccRateClass(getNetStat(net.statKey)?.succ)">
                  {{ getNetStat(net.statKey)?.succ || "—" }}
                </div>
              </div>

              <!-- 最新高度 -->
              <div class="stat-row">
                <div class="stat-label">
                  <icon-layers class="stat-icon block-icon" />
                  最新高度
                </div>
                <div class="stat-value block-value">
                  {{ getNetStat(net.statKey)?.block ? "#" + Number(getNetStat(net.statKey)?.block).toLocaleString() : "—" }}
                </div>
              </div>

              <!-- 最后同步 -->
              <div class="stat-row">
                <div class="stat-label">
                  <icon-clock-circle class="stat-icon time-icon" />
                  最后同步
                </div>
                <div class="stat-value time-value" :title="formatTime(getNetStat(net.statKey)?.time)">
                  {{ formatTime(getNetStat(net.statKey)?.time) }}
                </div>
              </div>
            </div>

            <!-- 无数据占位 -->
            <div v-if="!getNetStat(net.statKey)" class="no-data-overlay">
              <icon-minus-circle />
              <span>暂无数据</span>
            </div>
          </div>
        </div>
      </a-card>

      <!-- 主要内容：配置 -->
      <a-card :bordered="false" class="main-card">
        <template #title>
          <div class="card-title">
            <div class="title-icon">
              <icon-settings />
            </div>
            <span>区块网络配置</span>
            <a
              href="https://github.com/v03413/BEpusdt/blob/main/docs/faq/rpc-endpoint.md"
              target="_blank"
              class="title-doc-link"
            >
              <icon-exclamation-circle />
              必读：节点配置说明
            </a>
          </div>
        </template>

        <template #extra>
          <a-space size="small" wrap>
            <a-button @click="handleReset" :loading="loading" size="small" class="action-btn">
              <template #icon>
                <icon-refresh />
              </template>
              重置
            </a-button>
            <a-button type="primary" @click="handleSave" :loading="saveLoading" size="small" class="action-btn save-btn">
              <template #icon>
                <icon-save />
              </template>
              保存配置
            </a-button>
          </a-space>
        </template>

        <div class="form-container">
          <a-form :model="formData" layout="vertical" ref="formRef">
            <!-- Tron 网络配置区域 -->
            <div class="tron-section">
              <div class="section-header">
                <div class="header-icon">
                  <icon-fire />
                </div>
                <span class="header-title">Tron 网络</span>
              </div>
              <a-row :gutter="16">
                <a-col :xs="24" :sm="24" :md="12">
                  <a-form-item
                    field="rpc_endpoint_tron"
                    label="Tron RPC"
                    :rules="[{ required: true, message: '请输入Tron RPC' }]"
                    class="network-form-item"
                  >
                    <a-input
                      v-model="formData.rpc_endpoint_tron"
                      placeholder="请输入 Tron RPC"
                      allow-clear
                      size="small"
                      class="network-input tron-input"
                    >
                      <template #prefix>
                        <div class="input-icon">
                          <icon-link />
                        </div>
                      </template>
                    </a-input>
                  </a-form-item>
                </a-col>
                <a-col :span="24">
                  <a-form-item field="rpc_endpoint_tron_grid_api_key" class="network-form-item">
                    <template #label>
                      <div class="tron-grid-label">
                        <span class="label-with-tip">
                          <span>Tron Grid Api Key</span>
                          <a-tooltip content="配置独立 Api Key 可提高扫块稳定性，多个可用半角符逗号隔开。" position="top">
                            <icon-question-circle class="tip-icon" />
                          </a-tooltip>
                          <span class="optional-tag recommend-tag">
                            <icon-exclamation-circle />
                            推荐配置
                          </span>
                        </span>
                      </div>
                    </template>
                    <a-textarea
                      v-model="formData.rpc_endpoint_tron_grid_api_key"
                      placeholder="请输入 Tron Grid Api Key (可选)，多个可用半角符逗号隔开"
                      allow-clear
                      size="small"
                      class="network-input tron-input tron-grid-api-key-input"
                      :auto-size="{ minRows: 1, maxRows: 6 }"
                    >
                      <template>
                        <div class="input-icon">
                          <icon-safe />
                        </div>
                      </template>
                    </a-textarea>
                  </a-form-item>
                </a-col>
              </a-row>
            </div>

            <!-- 其他网络配置 -->
            <div class="other-section">
              <div class="section-header">
                <div class="header-icon">
                  <icon-link />
                </div>
                <span class="header-title">其他网络</span>
              </div>

              <a-row :gutter="[16, 6]">
                <a-col
                  v-for="network in networks.filter(n => n.key !== 'rpc_endpoint_tron')"
                  :key="network.key"
                  :xs="24"
                  :sm="24"
                  :md="12"
                  :lg="8"
                >
                  <a-form-item
                    :field="network.key"
                    :label="network.label"
                    :rules="[{ required: true, message: `请输入${network.label}` }]"
                    class="network-form-item"
                  >
                    <a-input
                      v-model="formData[network.key]"
                      :placeholder="`请输入 ${network.label}`"
                      allow-clear
                      size="small"
                      class="network-input"
                    >
                      <template #prefix>
                        <div class="input-icon">
                          <component :is="network.icon" />
                        </div>
                      </template>
                    </a-input>
                  </a-form-item>
                </a-col>
              </a-row>
            </div>
          </a-form>
        </div>
      </a-card>
    </a-spin>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted } from "vue";
import { Message } from "@arco-design/web-vue";
import { setsConfAPI, getRpcConfAPI } from "@/api/modules/conf/index";
import {
  IconSettings,
  IconSave,
  IconRefresh,
  IconExclamationCircle,
  IconLink,
  IconFire,
  IconQuestionCircle,
  IconSafe,
  IconThunderbolt,
  IconSync,
  IconCheckCircle,
  IconLayers,
  IconClockCircle,
  IconMinusCircle
} from "@arco-design/web-vue/es/icon";

// 网络配置（用于表单）
const networks = [
  { key: "rpc_endpoint_ethereum", label: "Ethereum RPC", icon: IconLink },
  { key: "rpc_endpoint_bsc", label: "BSC RPC", icon: IconLink },
  { key: "rpc_endpoint_polygon", label: "Polygon RPC", icon: IconLink },
  { key: "rpc_endpoint_arbitrum", label: "Arbitrum RPC", icon: IconLink },
  { key: "rpc_endpoint_base", label: "Base RPC", icon: IconLink },
  { key: "rpc_endpoint_xlayer", label: "X Layer RPC", icon: IconLink },
  { key: "rpc_endpoint_tron", label: "Tron RPC", icon: IconLink },
  { key: "rpc_endpoint_solana", label: "Solana RPC", icon: IconLink },
  { key: "rpc_endpoint_aptos", label: "Aptos RPC", icon: IconLink },
  { key: "rpc_endpoint_plasma", label: "Plasma RPC", icon: IconLink }
];

// 所有网络（用于状态总览，statKey 对应 stats 中的 key）
const allNetworks = [
  { key: "rpc_endpoint_tron", statKey: "tron", label: "Tron", icon: IconFire, color: "#e8503a" },
  { key: "rpc_endpoint_ethereum", statKey: "ethereum", label: "Ethereum", icon: IconLink, color: "#627eea" },
  { key: "rpc_endpoint_bsc", statKey: "bsc", label: "BSC", icon: IconLink, color: "#f0b90b" },
  { key: "rpc_endpoint_polygon", statKey: "polygon", label: "Polygon", icon: IconLink, color: "#8247e5" },
  { key: "rpc_endpoint_arbitrum", statKey: "arbitrum", label: "Arbitrum", icon: IconLink, color: "#28a0f0" },
  { key: "rpc_endpoint_base", statKey: "base", label: "Base", icon: IconLink, color: "#0052ff" },
  { key: "rpc_endpoint_xlayer", statKey: "xlayer", label: "X Layer", icon: IconLink, color: "#000000" },
  { key: "rpc_endpoint_solana", statKey: "solana", label: "Solana", icon: IconLink, color: "#9945ff" },
  { key: "rpc_endpoint_aptos", statKey: "aptos", label: "Aptos", icon: IconLink, color: "#00c2cb" },
  { key: "rpc_endpoint_plasma", statKey: "plasma", label: "Plasma", icon: IconLink, color: "#ff6b35" }
];

interface StatInfo {
  block: string;
  succ: string;
  time: number;
}

const loading = ref<boolean>(false);
const saveLoading = ref<boolean>(false);
const formRef = ref();
const formData = reactive<Record<string, string>>({});
const originalData = ref<Record<string, string>>({});
const statsData = ref<Record<string, StatInfo>>({});

let refreshTimer: ReturnType<typeof setInterval> | null = null;

const getNetStat = (statKey: string): StatInfo | undefined => {
  return statsData.value[statKey];
};

const getNetworkClass = (key: string) => {
  return key === "rpc_endpoint_tron" ? "tron-net-card" : "";
};

const getStatusDotClass = (statKey: string) => {
  const stat = statsData.value[statKey];
  if (!stat) return "dot-unknown";
  const succ = parseFloat(stat.succ);
  if (succ >= 95) return "dot-good";
  if (succ >= 80) return "dot-warn";
  return "dot-bad";
};

const getSuccRateClass = (succ?: string) => {
  if (!succ) return "";
  const val = parseFloat(succ);
  if (val >= 95) return "succ-good";
  if (val >= 80) return "succ-warn";
  return "succ-bad";
};

const formatTime = (ts?: number): string => {
  if (!ts) return "—";
  const d = new Date(ts * 1000);
  const pad = (n: number) => String(n).padStart(2, "0");
  return `${d.getFullYear()}/${d.getMonth() + 1}/${d.getDate()} ${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`;
};

const getConf = async () => {
  try {
    loading.value = true;

    const response = await getRpcConfAPI();
    const rpc = response.data?.rpc || {};
    const stats = response.data?.stats || {};

    networks.forEach(network => {
      formData[network.key] = rpc[network.key] || "";
    });
    formData.rpc_endpoint_tron_grid_api_key = rpc.rpc_endpoint_tron_grid_api_key || "";
    originalData.value = { ...formData };
    statsData.value = stats;
  } catch (error) {
    Message.error("获取配置失败");
    console.error("获取配置失败:", error);
  } finally {
    loading.value = false;
  }
};

const handleRefreshStats = async () => {
  try {
    const response = await getRpcConfAPI();
    statsData.value = response.data?.stats || {};
  } catch {}
};

const handleSave = async () => {
  try {
    const errors = await formRef.value?.validate();
    if (errors) {
      Message.error("表单验证失败，请检查所有字段");
      return;
    }
  } catch (validationError) {
    console.error("表单验证失败:", validationError);
    Message.error("请填写所有必填项");
    return;
  }

  try {
    saveLoading.value = true;

    const saveData: Array<{ key: string; value: string }> = [];

    networks.forEach(network => {
      const value = formData[network.key]?.trim();
      if (value) {
        saveData.push({ key: network.key, value });
      }
    });

    if (saveData.length < networks.length) {
      Message.error("所有RPC节点都必须填写");
      return;
    }

    const tronApiKey = formData.rpc_endpoint_tron_grid_api_key?.trim() || "";
    saveData.push({ key: "rpc_endpoint_tron_grid_api_key", value: tronApiKey });

    await setsConfAPI(saveData);
    Message.success("配置保存成功");
    await getConf();
  } catch (error) {
    Message.error("保存配置失败");
    console.error("保存配置失败:", error);
  } finally {
    saveLoading.value = false;
  }
};

const handleReset = () => {
  Object.assign(formData, originalData.value);
  Message.info("已重置为原始配置");
};

onMounted(() => {
  getConf();
  // 每3秒自动刷新状态
  refreshTimer = setInterval(handleRefreshStats, 3000);
});

onUnmounted(() => {
  if (refreshTimer) clearInterval(refreshTimer);
});
</script>

<style lang="scss" scoped>
.full-height {
  min-height: 100%;
}

.title-doc-link {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  margin-left: 4px;
  padding: 2px 10px;
  border-radius: 20px;
  background: rgba(var(--danger-6), 0.1);
  border: 1px solid rgba(var(--danger-6), 0.3);
  color: $color-danger;
  font-size: 12px;
  font-weight: 600;
  text-decoration: none;
  transition: all 0.2s ease;

  &:hover {
    background: rgba(var(--danger-6), 0.18);
    border-color: rgba(var(--danger-6), 0.55);
    box-shadow: 0 2px 8px rgba(var(--danger-6), 0.2);
  }
}

/* ============ 状态总览卡片 ============ */
.stats-card {
  margin-bottom: 12px;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
  background: $color-bg-2;

  :deep(.arco-card-header) {
    border-bottom: 1px solid $color-border-2;
    padding: 12px 16px;
    background: $color-bg-3;
    border-radius: 8px 8px 0 0;
  }

  :deep(.arco-card-body) {
    padding: 14px 16px;
  }
}

.live-tag {
  margin-left: 8px;
  font-size: 11px;
  display: inline-flex;
  align-items: center;
  gap: 3px;
  padding: 1px 6px;
}

.spin-icon {
  font-size: 11px;
  animation: spin 2s linear infinite;
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

.network-stats-grid {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 10px;

  @media (max-width: 1200px) {
    grid-template-columns: repeat(4, 1fr);
  }
  @media (max-width: 960px) {
    grid-template-columns: repeat(3, 1fr);
  }
  @media (max-width: 680px) {
    grid-template-columns: repeat(2, 1fr);
  }
  @media (max-width: 420px) {
    grid-template-columns: 1fr;
  }
}

.network-stat-card {
  position: relative;
  background: $color-bg-3;
  border: 1px solid $color-border-2;
  border-radius: 8px;
  padding: 10px 12px;
  transition: all 0.2s ease;
  overflow: hidden;

  &::before {
    content: "";
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 2px;
    background: $color-fill-3;
    transition: background 0.2s ease;
  }

  &:hover {
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
    transform: translateY(-1px);
    border-color: rgb(var(--primary-3));
  }

  &.tron-net-card::before {
    background: linear-gradient(90deg, #e8503a, #ff7b5a);
  }
}

.net-card-header {
  display: flex;
  align-items: center;
  gap: 7px;
  margin-bottom: 10px;
}

.net-badge {
  width: 24px;
  height: 24px;
  border-radius: 5px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 12px;
  flex-shrink: 0;
}

.net-name {
  font-weight: 600;
  font-size: 13px;
  color: $color-text-1;
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.net-status-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  flex-shrink: 0;

  &.dot-good {
    background: #00b42a;
    box-shadow: 0 0 0 2px rgba(0, 180, 42, 0.2);
    animation: pulse-good 2s infinite;
  }

  &.dot-warn {
    background: #ff7d00;
    box-shadow: 0 0 0 2px rgba(255, 125, 0, 0.2);
  }

  &.dot-bad {
    background: #f53f3f;
    box-shadow: 0 0 0 2px rgba(245, 63, 63, 0.2);
  }

  &.dot-unknown {
    background: $color-fill-4;
  }
}

@keyframes pulse-good {
  0%,
  100% {
    box-shadow: 0 0 0 2px rgba(0, 180, 42, 0.2);
  }
  50% {
    box-shadow: 0 0 0 4px rgba(0, 180, 42, 0.1);
  }
}

.net-card-body {
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.stat-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 4px;
}

.stat-label {
  display: flex;
  align-items: center;
  gap: 4px;
  color: $color-text-3;
  font-size: 11px;
  white-space: nowrap;
  flex-shrink: 0;

  .stat-icon {
    font-size: 11px;
  }

  .success-icon {
    color: #00b42a;
  }
  .block-icon {
    color: $color-primary;
  }
  .time-icon {
    color: $color-text-3;
  }
}

.stat-value {
  font-size: 12px;
  font-weight: 600;
  color: $color-text-1;
  text-align: right;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  min-width: 0;

  &.succ-good {
    color: #00b42a;
  }
  &.succ-warn {
    color: #ff7d00;
  }
  &.succ-bad {
    color: #f53f3f;
  }

  &.block-value {
    font-family: monospace;
    font-size: 11px;
    color: $color-primary;
  }

  &.time-value {
    font-size: 11px;
    color: $color-text-2;
    font-weight: 500;
    cursor: help;
  }
}

.no-data-overlay {
  position: absolute;
  inset: 0;
  background: rgba($color-bg-3, 0.85);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 4px;
  border-radius: 8px;
  color: $color-text-4;
  font-size: 11px;
  padding-top: 20px;

  svg {
    font-size: 20px;
  }
}

/* ============ 主配置卡片 ============ */
.main-card {
  margin-top: 0;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
  background: $color-bg-2;

  :deep(.arco-card-header) {
    border-bottom: 1px solid $color-border-2;
    padding: 14px 18px;
    background: $color-bg-3;
    border-radius: 8px 8px 0 0;
  }

  :deep(.arco-card-body) {
    padding: 16px;
  }
}

.card-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  font-size: 15px;
  color: $color-text-1;

  .title-icon {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 28px;
    height: 28px;
    background: $color-primary;
    border-radius: 6px;
    color: #fff;
    font-size: 13px;

    &.stats-icon {
      background: linear-gradient(135deg, #165dff, #722ed1);
    }
  }
}

.action-btn {
  border-radius: 6px;
  font-weight: 500;
  transition: all 0.2s ease;
  padding: 5px 14px;
  height: 30px;

  &:hover {
    transform: translateY(-1px);
    box-shadow: 0 3px 8px rgba(0, 0, 0, 0.12);
  }
}

.save-btn {
  background: $color-primary;
  border: none;

  &:hover {
    background: rgb(var(--primary-5));
  }
}

.form-container {
  margin: 12px 0;
}

.tron-section {
  background: rgba(var(--danger-6), 0.06);
  border: 1px solid rgba(var(--danger-6), 0.18);
  border-radius: 6px;
  padding: 10px 12px;
  margin-bottom: 12px;
  position: relative;
  overflow: hidden;

  &::before {
    content: "";
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 2px;
    background: rgba(var(--danger-6), 0.72);
  }

  .section-header {
    display: flex;
    align-items: center;
    gap: 6px;
    margin-bottom: 8px;
    padding-bottom: 6px;
    border-bottom: 1px solid $color-border-2;

    .header-icon {
      display: flex;
      align-items: center;
      justify-content: center;
      width: 20px;
      height: 20px;
      background: $color-danger;
      border-radius: 4px;
      color: #fff;
      font-size: 11px;
      box-shadow: 0 2px 4px rgba(var(--danger-6), 0.3);
    }

    .header-title {
      font-weight: 600;
      font-size: 13px;
      color: $color-text-1;
    }
  }
}

.other-section {
  background: rgba(var(--success-6), 0.06);
  border: 1px solid rgba(var(--success-6), 0.18);
  border-radius: 6px;
  padding: 10px 12px;
  margin-bottom: 12px;
  position: relative;
  overflow: hidden;

  &::before {
    content: "";
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 2px;
    background: rgba(var(--success-6), 0.72);
  }

  .section-header {
    display: flex;
    align-items: center;
    gap: 6px;
    margin-bottom: 8px;
    padding-bottom: 6px;
    border-bottom: 1px solid $color-border-2;

    .header-icon {
      display: flex;
      align-items: center;
      justify-content: center;
      width: 20px;
      height: 20px;
      background: $color-success;
      border-radius: 4px;
      color: #fff;
      font-size: 11px;
      box-shadow: 0 2px 4px rgba(var(--success-6), 0.3);
    }

    .header-title {
      font-weight: 600;
      font-size: 13px;
      color: $color-text-1;
    }
  }

  .network-input {
    :deep(.arco-input-wrapper) {
      border-color: $color-border-2;
      background: $color-bg-2;

      &:hover {
        border-color: $color-success;
        box-shadow: 0 0 0 2px rgba(var(--success-6), 0.08);
      }

      &.arco-input-focus {
        border-color: $color-success;
        box-shadow: 0 0 0 2px rgba(var(--success-6), 0.1);
      }
    }
  }
}

.network-form-item {
  :deep(.arco-form-item-label-col) {
    margin-bottom: 4px;

    .arco-form-item-label {
      font-weight: 500;
      color: $color-text-1;
      font-size: 12px;
    }
  }
}

.network-input {
  border-radius: 6px;
  transition: all 0.2s ease;

  :deep(.arco-input-wrapper) {
    border: 1px solid $color-border-2;
    background: $color-bg-2;
    height: 32px;

    &:hover {
      border-color: $color-primary;
      box-shadow: 0 0 0 2px rgba(var(--primary-6), 0.08);
    }

    &.arco-input-focus {
      border-color: $color-primary;
      box-shadow: 0 0 0 2px rgba(var(--primary-6), 0.1);
    }
  }

  :deep(.arco-input) {
    font-size: 12px;
  }
}

.tron-input {
  :deep(.arco-input-wrapper) {
    border-color: $color-border-2;

    &:hover {
      border-color: $color-danger;
      box-shadow: 0 0 0 2px rgba(var(--danger-6), 0.08);
    }

    &.arco-input-focus {
      border-color: $color-danger;
      box-shadow: 0 0 0 2px rgba(var(--danger-6), 0.1);
    }
  }
}

.tron-grid-api-key-input {
  max-width: 100%;

  :deep(textarea) {
    max-height: 120px;
    overflow-y: auto;
    line-height: 20px;
  }
}

.input-icon {
  display: flex;
  align-items: center;
  color: $color-text-3;
  font-size: 13px;
}

.tron-grid-label {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  max-width: 100%;
  min-width: 0;
  flex-wrap: wrap;
}

.label-with-tip {
  display: flex;
  align-items: center;
  gap: 5px;

  .tip-icon {
    color: $color-text-3;
    cursor: help;
    font-size: 13px;

    &:hover {
      color: $color-primary;
    }
  }

  .optional-tag {
    color: $color-text-3;
    font-size: 11px;
    font-weight: normal;
  }

  .recommend-tag {
    display: inline-flex;
    align-items: center;
    gap: 3px;
    padding: 1px 7px;
    border-radius: 20px;
    background: rgba(var(--warning-6), 0.12);
    border: 1px solid rgba(var(--warning-6), 0.4);
    color: rgb(var(--warning-6));
    font-size: 11px;
    font-weight: 600;
    animation: recommend-pulse 2.5s ease-in-out infinite;
  }

  @keyframes recommend-pulse {
    0%, 100% { box-shadow: 0 0 0 0 rgba(var(--warning-6), 0); }
    50% { box-shadow: 0 0 0 3px rgba(var(--warning-6), 0.15); }
  }
}

.help-link {
  display: flex;
  align-items: center;
  gap: 3px;
  color: $color-danger;
  font-size: 11px;
  text-decoration: none;
  transition: all 0.2s ease;
  font-weight: 500;

  &:hover {
    color: rgb(var(--danger-5));
  }
}

:deep(.arco-card.arco-card-bordered) {
  border: 1px solid $color-border-2;
}
</style>
