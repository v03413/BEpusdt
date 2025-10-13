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
          一般情况下不推荐修改RPC节点，除非您非常了解区块网络并确保节点的可用性和稳定性。
        </div>
      </a-alert>

      <!-- 主要内容 -->
      <a-card :bordered="false" class="main-card">
        <template #title>
          <div class="card-title">
            <div class="title-icon">
              <icon-settings />
            </div>
            <span>区块网络 RPC 节点配置</span>
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
            <a-row :gutter="[16, 8]">
              <a-col v-for="network in networks" :key="network.key" :span="8">
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
  IconFire
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
    const keys = networks.map(network => network.key);

    const response = await getsConfAPI({ keys });
    const data = response.data || {};

    networks.forEach(network => {
      formData[network.key] = data[network.key] || "";
    });

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
    networks.forEach(network => {
      const value = formData[network.key]?.trim();
      if (value) {
        saveData.push({
          key: network.key,
          value: value
        });
      }
    });

    if (saveData.length !== networks.length) {
      Message.error("所有RPC节点都必须填写");
      return;
    }

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
  margin-bottom: 16px;
  border-radius: 6px;
  box-shadow: 0 2px 6px rgba(255, 125, 0, 0.08);
  background: linear-gradient(135deg, #fff7e6 0%, #fff1d6 100%);
  border: 1px solid #ffb84d;

  :deep(.arco-alert-content) {
    font-size: 13px;
    line-height: 1.5;
  }

  :deep(.arco-alert) {
    padding: 12px 16px;
  }
}

.main-card {
  margin-top: 0;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);
  background: linear-gradient(135deg, #ffffff 0%, #fafbfc 100%);

  :deep(.arco-card-header) {
    border-bottom: 1px solid #f2f3f5;
    padding: 16px 20px;
    background: linear-gradient(135deg, #f8f9fa 0%, #ffffff 100%);
    border-radius: 8px 8px 0 0;
  }

  :deep(.arco-card-body) {
    padding: 20px;
  }
}

.card-title {
  display: flex;
  align-items: center;
  gap: 10px;
  font-weight: 600;
  font-size: 16px;
  color: #1d2129;

  .title-icon {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 32px;
    height: 32px;
    background: linear-gradient(135deg, #165dff 0%, #246fff 100%);
    border-radius: 6px;
    color: white;
    font-size: 14px;
  }
}

.action-btn {
  border-radius: 6px;
  font-weight: 500;
  transition: all 0.2s ease;
  padding: 6px 16px;
  height: 32px;

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
  margin: 16px 0;
}

.network-form-item {
  :deep(.arco-form-item-label-col) {
    margin-bottom: 6px;

    .arco-form-item-label {
      font-weight: 500;
      color: #1d2129;
      font-size: 13px;
    }
  }

  :deep(.arco-form-item) {
    margin-bottom: 16px;
  }
}

.network-input {
  border-radius: 6px;
  transition: all 0.2s ease;

  :deep(.arco-input-wrapper) {
    border: 1px solid #e5e6eb;
    background: #ffffff;
    height: 36px;

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
    font-size: 13px;
  }
}

.input-icon {
  display: flex;
  align-items: center;
  color: #86909c;
  font-size: 14px;
}

.info-divider {
  margin: 24px 0 16px 0;

  .divider-content {
    display: flex;
    align-items: center;
    gap: 6px;
    color: #1d2129;
    font-weight: 500;
    font-size: 14px;
  }

  :deep(.arco-divider-text) {
    background: linear-gradient(135deg, #f8f9fa 0%, #ffffff 100%);
    border: 1px solid #e5e6eb;
    border-radius: 16px;
    padding: 6px 12px;
    font-size: 13px;
  }
}

.info-section {
  background: linear-gradient(135deg, #f8f9fa 0%, #ffffff 100%);
  border-radius: 8px;
  padding: 16px;
  border: 1px solid #e5e6eb;
}

.info-paragraph {
  margin: 0;
}

.info-grid {
  display: grid;
  gap: 12px;
}

.info-item {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  padding: 10px;
  background: white;
  border-radius: 6px;
  border: 1px solid #f2f3f5;
  transition: all 0.2s ease;

  &:hover {
    transform: translateY(-1px);
    box-shadow: 0 3px 8px rgba(0, 0, 0, 0.06);
    border-color: #165dff;
  }

  .info-icon {
    flex-shrink: 0;
    width: 18px;
    height: 18px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #165dff;
    font-size: 14px;
    margin-top: 1px;
  }

  span {
    color: #4e5969;
    line-height: 1.5;
    font-size: 13px;
  }
}

// 响应式设计
@media (max-width: 768px) {
  .card-title {
    font-size: 15px;

    .title-icon {
      width: 28px;
      height: 28px;
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
      padding: 14px 16px;
    }

    :deep(.arco-card-body) {
      padding: 16px;
    }
  }
}

// 暗色主题适配
:deep(.arco-card.arco-card-bordered) {
  border: 1px solid #e5e6eb;
}
</style>
