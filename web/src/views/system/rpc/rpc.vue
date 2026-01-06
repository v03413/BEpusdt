<template>
  <div class="snow-page">
    <a-spin :loading="loading" tip="正在加载..." class="full-height">
      <!-- 警告横幅 -->
      <a-alert type="warning" show-icon class="warning-banner">
        <template #icon>
          <icon-exclamation-circle />
        </template>
        <div>
          <strong>重要提醒：</strong>
          一般情况下不推荐修改RPC节点,除非您非常了解区块网络并确保节点的可用性和稳定性。
        </div>
      </a-alert>

      <!-- 主要内容 -->
      <a-card :bordered="false" class="main-card">
        <template #title>
          <div class="card-title">
            <div class="title-icon">
              <icon-settings />
            </div>
            <span>区块网络配置</span>
          </div>
        </template>

        <template #extra>
          <a-space size="small">
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
                <a-col :span="12">
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
                <a-col :span="12">
                  <a-form-item field="rpc_endpoint_tron_grid_api_key" class="network-form-item">
                    <template #label>
                      <div class="label-with-tip">
                        <span>Tron Grid Api Key</span>
                        <a-tooltip content="使用个人独立的Api Key，可以提高扫块成功率" position="top">
                          <icon-question-circle class="tip-icon" />
                        </a-tooltip>
                        <span class="optional-tag">(可选)</span>
                      </div>
                    </template>
                    <a-input
                      v-model="formData.rpc_endpoint_tron_grid_api_key"
                      placeholder="请输入 Tron Grid Api Key (可选)"
                      allow-clear
                      size="small"
                      class="network-input tron-input"
                    >
                      <template #prefix>
                        <div class="input-icon">
                          <icon-safe />
                        </div>
                      </template>
                      <template #suffix>
                        <a href="https://github.com/v03413/BEpusdt/" target="_blank" class="help-link">
                          <icon-question-circle />
                          获取方法
                        </a>
                      </template>
                    </a-input>
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
                <a-col v-for="network in networks.filter(n => n.key !== 'rpc_endpoint_tron')" :key="network.key" :span="8">
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

        <!-- 配置说明 -->
        <a-divider orientation="left" class="info-divider">
          <div class="divider-content">
            <icon-info-circle />
            <span>配置说明</span>
          </div>
        </a-divider>

        <div class="info-section">
          <div class="info-grid">
            <div v-for="(info, index) in infoList" :key="index" class="info-item">
              <div class="info-icon">
                <component :is="info.icon" />
              </div>
              <span>{{ info.text }}</span>
            </div>
          </div>
        </div>
      </a-card>
    </a-spin>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from "vue";
import { Message } from "@arco-design/web-vue";
import { getsConfAPI, setsConfAPI } from "@/api/modules/conf/index";
import {
  IconSettings,
  IconSave,
  IconRefresh,
  IconExclamationCircle,
  IconLink,
  IconInfoCircle,
  IconCheckCircle,
  IconStar,
  IconThunderbolt,
  IconFire,
  IconQuestionCircle,
  IconSafe
} from "@arco-design/web-vue/es/icon";

// 网络配置
const networks = [
  { key: "rpc_endpoint_ethereum", label: "Ethereum RPC", icon: IconLink },
  { key: "rpc_endpoint_bsc", label: "BSC RPC", icon: IconLink },
  { key: "rpc_endpoint_polygon", label: "Polygon RPC", icon: IconLink },
  { key: "rpc_endpoint_arbitrum", label: "Arbitrum RPC", icon: IconLink },
  { key: "rpc_endpoint_base", label: "Base RPC", icon: IconLink },
  { key: "rpc_endpoint_xlayer", label: "X Layer RPC", icon: IconLink },
  { key: "rpc_endpoint_tron", label: "Tron RPC", icon: IconLink },
  { key: "rpc_endpoint_solana", label: "Solana RPC", icon: IconLink },
  { key: "rpc_endpoint_aptos", label: "Aptos RPC", icon: IconLink }
];

const infoList = [
  { icon: IconCheckCircle, text: "RPC节点是与区块链网络通信的关键接口，请确保所配置的节点稳定可靠" },
  { icon: IconStar, text: "建议使用官方推荐的RPC节点或知名的第三方服务商" },
  { icon: IconThunderbolt, text: "配置前请先测试节点的连通性和响应速度" },
  { icon: IconFire, text: "修改配置后系统将立即生效，请谨慎操作" }
];

const loading = ref<boolean>(false);
const saveLoading = ref<boolean>(false);
const formRef = ref();
const formData = reactive<Record<string, string>>({});
const originalData = ref<Record<string, string>>({});

const getConf = async () => {
  try {
    loading.value = true;
    const keys = [...networks.map(network => network.key), "rpc_endpoint_tron_grid_api_key"];

    const response = await getsConfAPI({ keys });
    const data = response.data || {};

    networks.forEach(network => {
      formData[network.key] = data[network.key] || "";
    });
    formData.rpc_endpoint_tron_grid_api_key = data.rpc_endpoint_tron_grid_api_key || "";

    originalData.value = { ...formData };
  } catch (error) {
    Message.error("获取配置失败");
    console.error("获取配置失败:", error);
  } finally {
    loading.value = false;
  }
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

    // 构建保存数据数组
    const saveData: Array<{ key: string; value: string }> = [];

    // 添加所有网络的 RPC 配置
    networks.forEach(network => {
      const value = formData[network.key]?.trim();
      if (value) {
        saveData.push({
          key: network.key,
          value: value
        });
      }
    });

    // 验证所有必填的 RPC 节点是否都已填写
    if (saveData.length < networks.length) {
      Message.error("所有RPC节点都必须填写");
      return;
    }

    // 添加 Tron Grid Api Key (可选，但即使为空也要保存)
    const tronApiKey = formData.rpc_endpoint_tron_grid_api_key?.trim() || "";
    saveData.push({
      key: "rpc_endpoint_tron_grid_api_key",
      value: tronApiKey
    });

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

// 重置配置
const handleReset = () => {
  Object.assign(formData, originalData.value);
  Message.info("已重置为原始配置");
};

onMounted(() => {
  getConf();
});
</script>

<style lang="scss" scoped>
.full-height {
  min-height: 100%;
}

.warning-banner {
  margin-bottom: 12px;
  border-radius: 6px;
  box-shadow: 0 2px 6px rgba(255, 125, 0, 0.08);
  background: linear-gradient(135deg, #fff7e6 0%, #fff1d6 100%);
  border: 1px solid #ffb84d;

  :deep(.arco-alert-content) {
    font-size: 13px;
    line-height: 1.4;
  }

  :deep(.arco-alert) {
    padding: 10px 14px;
  }
}

.main-card {
  margin-top: 0;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
  background: linear-gradient(135deg, #ffffff 0%, #fafbfc 100%);

  :deep(.arco-card-header) {
    border-bottom: 1px solid #f2f3f5;
    padding: 14px 18px;
    background: linear-gradient(135deg, #f8f9fa 0%, #ffffff 100%);
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
  color: #1d2129;

  .title-icon {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 28px;
    height: 28px;
    background: linear-gradient(135deg, #165dff 0%, #246fff 100%);
    border-radius: 6px;
    color: white;
    font-size: 13px;
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
  background: linear-gradient(135deg, #165dff 0%, #246fff 100%);
  border: none;

  &:hover {
    background: linear-gradient(135deg, #0e42d2 0%, #1a5dff 100%);
  }
}

.form-container {
  margin: 12px 0;
}

// Tron 配置区域样式 - 使用 Tron 官方红色系
.tron-section {
  background: linear-gradient(135deg, #fff5f5 0%, #fffafa 100%);
  border: 1px solid #ffccc7;
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
    background: linear-gradient(90deg, #ef1e23 0%, #ff4d4f 100%);
  }

  .section-header {
    display: flex;
    align-items: center;
    gap: 6px;
    margin-bottom: 8px;
    padding-bottom: 6px;
    border-bottom: 1px solid #ffccc7;

    .header-icon {
      display: flex;
      align-items: center;
      justify-content: center;
      width: 20px;
      height: 20px;
      background: linear-gradient(135deg, #ef1e23 0%, #ff4d4f 100%);
      border-radius: 4px;
      color: white;
      font-size: 11px;
      box-shadow: 0 2px 4px rgba(239, 30, 35, 0.3);
    }

    .header-title {
      font-weight: 600;
      font-size: 13px;
      color: #1d2129;
    }
  }
}

// 其他网络配置区域样式 - 使用柔和的浅绿色系
.other-section {
  background: linear-gradient(135deg, #f6ffed 0%, #fcfff9 100%);
  border: 1px solid #d9f7be;
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
    background: linear-gradient(90deg, #52c41a 0%, #73d13d 100%);
  }

  .section-header {
    display: flex;
    align-items: center;
    gap: 6px;
    margin-bottom: 8px;
    padding-bottom: 6px;
    border-bottom: 1px solid #d9f7be;

    .header-icon {
      display: flex;
      align-items: center;
      justify-content: center;
      width: 20px;
      height: 20px;
      background: linear-gradient(135deg, #52c41a 0%, #73d13d 100%);
      border-radius: 4px;
      color: white;
      font-size: 11px;
      box-shadow: 0 2px 4px rgba(82, 196, 26, 0.3);
    }

    .header-title {
      font-weight: 600;
      font-size: 13px;
      color: #1d2129;
    }
  }

  .network-input {
    :deep(.arco-input-wrapper) {
      border-color: #d9f7be;
      background: #ffffff;

      &:hover {
        border-color: #52c41a;
        box-shadow: 0 0 0 2px rgba(82, 196, 26, 0.08);
      }

      &.arco-input-focus {
        border-color: #52c41a;
        box-shadow: 0 0 0 2px rgba(82, 196, 26, 0.1);
      }
    }
  }
}

.network-form-item {
  :deep(.arco-form-item-label-col) {
    margin-bottom: 4px;

    .arco-form-item-label {
      font-weight: 500;
      color: #1d2129;
      font-size: 12px;
    }
  }

  :deep(.arco-form-item) {
    margin-bottom: 0;
  }
}

.network-input {
  border-radius: 6px;
  transition: all 0.2s ease;

  :deep(.arco-input-wrapper) {
    border: 1px solid #e5e6eb;
    background: #ffffff;
    height: 32px;

    &:hover {
      border-color: #165dff;
      box-shadow: 0 0 0 2px rgba(22, 93, 255, 0.08);
    }

    &.arco-input-focus {
      border-color: #165dff;
      box-shadow: 0 0 0 2px rgba(22, 93, 255, 0.1);
    }
  }

  :deep(.arco-input) {
    font-size: 12px;
  }
}

// Tron 输入框特殊样式
.tron-input {
  :deep(.arco-input-wrapper) {
    border-color: #ffccc7;

    &:hover {
      border-color: #ef1e23;
      box-shadow: 0 0 0 2px rgba(239, 30, 35, 0.08);
    }

    &.arco-input-focus {
      border-color: #ef1e23;
      box-shadow: 0 0 0 2px rgba(239, 30, 35, 0.1);
    }
  }
}

.input-icon {
  display: flex;
  align-items: center;
  color: #86909c;
  font-size: 13px;
}

.label-with-tip {
  display: flex;
  align-items: center;
  gap: 5px;

  .tip-icon {
    color: #86909c;
    cursor: help;
    font-size: 13px;

    &:hover {
      color: #165dff;
    }
  }

  .optional-tag {
    color: #86909c;
    font-size: 11px;
    font-weight: normal;
  }
}

.help-link {
  display: flex;
  align-items: center;
  gap: 3px;
  color: #ef1e23;
  font-size: 11px;
  text-decoration: none;
  transition: all 0.2s ease;
  font-weight: 500;

  &:hover {
    color: #d11a1f;
  }
}

.info-divider {
  margin: 16px 0 12px 0;

  .divider-content {
    display: flex;
    align-items: center;
    gap: 5px;
    color: #1d2129;
    font-weight: 500;
    font-size: 13px;
  }

  :deep(.arco-divider-text) {
    background: linear-gradient(135deg, #f8f9fa 0%, #ffffff 100%);
    border: 1px solid #e5e6eb;
    border-radius: 12px;
    padding: 4px 10px;
    font-size: 12px;
  }
}

.info-section {
  background: linear-gradient(135deg, #f8f9fa 0%, #ffffff 100%);
  border-radius: 6px;
  padding: 12px;
  border: 1px solid #e5e6eb;
}

.info-grid {
  display: grid;
  gap: 8px;
}

.info-item {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  padding: 7px;
  background: white;
  border-radius: 5px;
  border: 1px solid #f2f3f5;
  transition: all 0.2s ease;

  &:hover {
    transform: translateY(-1px);
    box-shadow: 0 2px 6px rgba(0, 0, 0, 0.06);
    border-color: #165dff;
  }

  .info-icon {
    flex-shrink: 0;
    width: 16px;
    height: 16px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #165dff;
    font-size: 13px;
    margin-top: 1px;
  }

  span {
    color: #4e5969;
    line-height: 1.4;
    font-size: 12px;
  }
}

// 响应式设计
@media (max-width: 768px) {
  .card-title {
    font-size: 14px;

    .title-icon {
      width: 26px;
      height: 26px;
      font-size: 12px;
    }
  }

  :deep(.arco-col) {
    span: 24 !important;
  }

  .info-grid {
    grid-template-columns: 1fr;
  }

  .main-card {
    :deep(.arco-card-header) {
      padding: 12px 14px;
    }

    :deep(.arco-card-body) {
      padding: 14px;
    }
  }

  .tron-section,
  .other-section {
    padding: 8px 10px;
  }
}

// 暗色主题适配
:deep(.arco-card.arco-card-bordered) {
  border: 1px solid #e5e6eb;
}
</style>
