<template>
  <div class="rate-syntax">
    <div class="snow-page">
      <div class="snow-inner">
        <a-row :gutter="16" style="margin: 16px 0">
          <a-col :span="12">
            <a-space size="medium">
              <a-button type="primary" @click="showSyncModal">
                <template #icon>
                  <icon-settings />
                </template>
                同步配置
              </a-button>
              <a-button type="primary" @click="showAtomModal" :status="'danger'">
                <template #icon>
                  <icon-robot-add />
                </template>
                支付颗粒度
              </a-button>
            </a-space>
          </a-col>
        </a-row>
        <a-table
          row-key="key"
          :size="'medium'"
          :bordered="{ cell: true }"
          :scroll="{ x: 1400, y: 600 }"
          :loading="loading"
          :columns="columns"
          :data="data"
          v-model:selectedKeys="selectedKeys"
          :pagination="false"
        >
          <template #fiat="{ record }">
            <span class="fiat-display">
              {{ getFiatFlag(record.fiat) }} <strong>{{ record.fiat }}</strong>
            </span>
          </template>
          <template #crypto="{ record }">
            <a-tag :color="getCryptoColor(record.crypto)" :bordered="true">
              {{ record.crypto }}
            </a-tag>
          </template>
          <template #syntax="{ record }">
            <div class="syntax-display">
              <span class="syntax-value">{{ record.syntax || "无" }}</span>
              <span class="syntax-description">{{ getTableSyntaxDescription(record.syntax) }}</span>
            </div>
          </template>
          <template #optional="{ record }">
            <a-space>
              <a-button size="mini" type="primary" @click="onEdit(record)">编辑</a-button>
            </a-space>
          </template>
        </a-table>
      </div>
    </div>

    <!-- 编辑汇率语法模态框 -->
    <a-modal
      v-model:visible="editModalVisible"
      title="编辑汇率语法"
      @ok="handleEditSubmit"
      @cancel="handleEditCancel"
      :ok-loading="editLoading"
      width="480px"
      class="edit-modal"
    >
      <a-form ref="editFormRef" :model="editForm" layout="vertical">
        <a-row :gutter="12">
          <a-col :span="12">
            <a-form-item label="交易法币">
              <a-input v-model="editForm.fiat" readonly size="small">
                <template #prefix>{{ getFiatFlag(editForm.fiat) }}</template>
              </a-input>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="加密货币">
              <a-tag :color="getCryptoColor(editForm.crypto)" :bordered="true">
                {{ editForm.crypto }}
              </a-tag>
            </a-form-item>
          </a-col>
        </a-row>

        <a-form-item label="语法类型">
          <a-radio-group v-model="syntaxType" @change="handleSyntaxTypeChange">
            <a-radio value="">固定数值</a-radio>
            <a-radio value="+">固定增加</a-radio>
            <a-radio value="-">固定减少</a-radio>
            <a-radio value="~">百分比浮动</a-radio>
          </a-radio-group>
        </a-form-item>

        <a-form-item label="数值">
          <a-input-number
            v-model="syntaxValue"
            :placeholder="getSyntaxPlaceholder()"
            :min="syntaxType === '~' ? 0.000001 : 0"
            :max="syntaxType === '~' ? 10 : 999999"
            :step="syntaxType === '~' ? 0.000001 : 0.01"
            style="width: 100%"
          >
            <template #prefix v-if="syntaxType">
              <span class="syntax-prefix">{{ syntaxType }}</span>
            </template>
          </a-input-number>
        </a-form-item>

        <div v-if="getFormSyntaxDescription()" class="syntax-tip">
          <a-typography-text type="secondary">
            <icon-info-circle style="margin-right: 4px" />
            {{ getFormSyntaxDescription() }}
          </a-typography-text>
        </div>
      </a-form>
    </a-modal>

    <!-- 同步频率设置模态框 -->
    <a-modal
      v-model:visible="syncModalVisible"
      title="汇率同步配置"
      @ok="handleSyncSubmit"
      @cancel="handleSyncCancel"
      :ok-loading="syncLoading"
      width="480px"
      class="sync-modal"
    >
      <a-form ref="syncFormRef" :model="syncForm" layout="vertical">
        <a-form-item label="同步频率（分钟）">
          <a-input-number
            v-model="syncForm.minutes"
            :min="10"
            :max="1440"
            :precision="0"
            placeholder="请输入同步频率"
            style="width: 100%"
          />
        </a-form-item>

        <a-form-item label="API 接口">
          <a-select v-model="syncForm.apiUrl" placeholder="请选择 API 接口" style="width: 100%">
            <a-option v-for="option in apiUrlOptions" :key="option.value" :value="option.value" :label="option.label">
              {{ option.label }}
            </a-option>
          </a-select>
        </a-form-item>

        <a-form-item label="API Key">
          <a-input v-model="syncForm.apiKey" placeholder="请输入 API Key（可选）" allow-clear style="width: 100%" />
        </a-form-item>

        <a-form-item label="汇率保留天数">
          <a-input-number
            v-model="syncForm.historyDays"
            :min="1"
            :max="365"
            :precision="0"
            placeholder="请输入汇率保留天数"
            style="width: 100%"
          />
        </a-form-item>

        <div class="sync-tip">
          <a-typography-text type="secondary">
            <icon-info-circle style="margin-right: 4px" />
            同步频率：10-1440分钟，推荐60分钟<br />
            官方接口：免费但有速率限制，配置
            <a-link href="https://www.coingecko.com/" target="_blank" :hoverable="false">API Key</a-link>
            可解除限制<br />
            开源接口：作者提供的免费缓存接口（落后官方接口3分钟），无速率限制<br />
            <b class="sync-warning">官方接口特指 CoinGecko，是全球最大的独立加密货币数据聚合平台之一</b>
            <hr />
          </a-typography-text>
        </div>
      </a-form>
    </a-modal>

    <!-- 支付颗粒度设置模态框 -->
    <a-modal
      v-model:visible="atomModalVisible"
      title="设置支付颗粒度"
      @ok="handleAtomSubmit"
      @cancel="handleAtomCancel"
      :ok-loading="atomLoading"
      width="400px"
      class="atom-modal"
    >
      <a-form ref="atomFormRef" :model="atomForm" layout="vertical">
        <a-form-item label="USDT 颗粒度">
          <a-input-number
            v-model="atomForm.usdt"
            :min="0.000001"
            :max="100"
            :precision="undefined"
            :step="0.000001"
            placeholder="推荐0.01"
            style="width: 100%"
          />
        </a-form-item>

        <a-form-item label="USDC 颗粒度">
          <a-input-number
            v-model="atomForm.usdc"
            :min="0.000001"
            :max="100"
            :precision="undefined"
            :step="0.000001"
            placeholder="推荐0.01"
            style="width: 100%"
          />
        </a-form-item>

        <a-form-item label="TRX 颗粒度">
          <a-input-number
            v-model="atomForm.trx"
            :min="0.000001"
            :max="100"
            :precision="undefined"
            :step="0.000001"
            placeholder="推荐0.01"
            style="width: 100%"
          />
        </a-form-item>

        <a-form-item label="BNB 颗粒度">
          <a-input-number
            v-model="atomForm.bnb"
            :min="0.00000001"
            :max="100"
            :precision="undefined"
            :step="0.000001"
            placeholder="推荐0.00001"
            style="width: 100%"
          />
        </a-form-item>

        <a-form-item label="ETH 颗粒度">
          <a-input-number
            v-model="atomForm.eth"
            :min="0.00000001"
            :max="100"
            :precision="undefined"
            :step="0.000001"
            placeholder="推荐0.000001"
            style="width: 100%"
          />
        </a-form-item>

        <div class="atom-tip">
          <a-typography-text type="secondary">
            <icon-info-circle style="margin-right: 4px" />
            支付数额递增时的最小单位，支付数额的最终保留位数；除非你明确知道其功能含义，一般情况下不推荐修改。
          </a-typography-text>
        </div>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from "vue";
import { Message } from "@arco-design/web-vue";
import { IconInfoCircle } from "@arco-design/web-vue/es/icon";
import { getSyntaxListAPI, setSyntaxAPI } from "@/api/modules/rate/index";
import { getsConfAPI, setsConfAPI } from "@/api/modules/conf/index";
import { List, EditForm } from "./syntax";
import { getFiatFlag, getCryptoColor } from "@/views/rate/common";

const selectedKeys = ref<string[]>([]);
const loading = ref<boolean>(false);
const data = reactive<List[]>([]);
const editModalVisible = ref<boolean>(false);
const editLoading = ref<boolean>(false);
const editFormRef = ref();
const syntaxType = ref<string>("");
const syntaxValue = ref<number | undefined>(undefined);

const editForm = reactive<EditForm>({
  fiat: "",
  crypto: "",
  syntax: ""
});

const columns = [
  {
    title: "交易法币",
    dataIndex: "fiat",
    align: "center",
    width: 160,
    slotName: "fiat",
    filterable: {
      filters: [
        { text: "🇨🇳 CNY", value: "CNY" },
        { text: "🇺🇸 USD", value: "USD" },
        { text: "🇯🇵 JPY", value: "JPY" },
        { text: "🇪🇺 EUR", value: "EUR" },
        { text: "🇬🇧 GBP", value: "GBP" }
      ],
      filter: (fiat: any, record: any) => fiat.includes(record.fiat),
      multiple: true
    }
  },
  {
    title: "加密货币",
    dataIndex: "crypto",
    align: "center",
    width: 140,
    slotName: "crypto",
    filterable: {
      filters: [
        { text: "USDT", value: "USDT" },
        { text: "USDC", value: "USDC" },
        { text: "TRX", value: "TRX" },
        { text: "ETH", value: "ETH" },
        { text: "BNB", value: "BNB" }
      ],
      filter: (crypto: any, record: any) => crypto.includes(record.crypto),
      multiple: true
    }
  },
  {
    title: "汇率浮动",
    dataIndex: "syntax",
    slotName: "syntax",
    width: 480
  },
  {
    title: "操作",
    slotName: "optional",
    align: "center",
    fixed: "right",
    width: 120
  }
];

const parseSyntax = (syntax: string) => {
  if (!syntax) return { type: "", value: undefined };

  if (syntax.startsWith("~")) return { type: "~", value: parseFloat(syntax.substring(1)) };
  if (syntax.startsWith("+")) return { type: "+", value: parseFloat(syntax.substring(1)) };
  if (syntax.startsWith("-")) return { type: "-", value: parseFloat(syntax.substring(1)) };
  return { type: "", value: parseFloat(syntax) };
};

const generateSyntax = () => {
  if (syntaxValue.value === undefined || syntaxValue.value === null) return "";

  // 格式化数值,去除末尾的0
  const formatValue = (val: number) => {
    return parseFloat(val.toFixed(6)).toString();
  };

  // 百分比浮动
  if (syntaxType.value === "~") {
    return syntaxType.value + formatValue(syntaxValue.value);
  }

  // 其他类型
  return syntaxType.value + formatValue(syntaxValue.value);
};

const getSyntaxPlaceholder = () => {
  const placeholders = {
    "+": "如：0.3",
    "-": "如：0.2",
    "~": "如：1.020000 或 0.970000",
    "": "如：7.4"
  };
  return placeholders[syntaxType.value as keyof typeof placeholders];
};

const getTableSyntaxDescription = (syntax: string) => {
  if (!syntax) return "";

  const parsed = parseSyntax(syntax);
  if (parsed.value === undefined || parsed.value === null) return "";

  // 格式化数值,去除末尾的0
  const formatValue = (val: number) => {
    return parseFloat(val.toFixed(6)).toString();
  };

  switch (parsed.type) {
    case "+":
      return `订单汇率 = 基准汇率 + ${formatValue(parsed.value)}`;
    case "-":
      return `订单汇率 = 基准汇率 - ${formatValue(parsed.value)}`;
    case "~":
      return parsed.value != 1 ? `订单汇率 = 基准汇率 * ${formatValue(parsed.value)}` : `订单汇率 = 基准汇率`;
    default:
      return `订单汇率强制固定 ${formatValue(parsed.value)}`;
  }
};

const getFormSyntaxDescription = () => {
  if (!syntaxType.value || syntaxValue.value === undefined || syntaxValue.value === null) return "";

  // 格式化数值,去除末尾的0
  const formatValue = (val: number) => {
    return parseFloat(val.toFixed(6)).toString();
  };

  switch (syntaxType.value) {
    case "+":
      return `订单汇率 = 基准汇率 + ${formatValue(syntaxValue.value)}`;
    case "-":
      return `订单汇率 = 基准汇率 - ${formatValue(syntaxValue.value)}`;
    case "~":
      return syntaxValue.value != 1 ? `订单汇率 = 基准汇率 * ${formatValue(syntaxValue.value)}` : `订单汇率 = 基准汇率`;
    default:
      return `订单汇率强制固定 ${formatValue(syntaxValue.value)}`;
  }
};

const handleSyntaxTypeChange = () => {
  if (syntaxType.value === "~" && (syntaxValue.value === undefined || syntaxValue.value === null || syntaxValue.value === 0)) {
    syntaxValue.value = 1.0;
  } else if (syntaxType.value !== "~" && syntaxValue.value === 1) {
    syntaxValue.value = 0;
  }
};

const getCommonTableList = async () => {
  try {
    loading.value = true;
    const res = await getSyntaxListAPI();
    data.length = 0;
    data.push(...res.data);
  } finally {
    loading.value = false;
  }
};

const onEdit = (record: List) => {
  editForm.fiat = record.fiat;
  editForm.crypto = record.crypto;
  editForm.syntax = record.syntax;

  const parsed = parseSyntax(record.syntax);
  syntaxType.value = parsed.type;
  syntaxValue.value = parsed.value;
  editModalVisible.value = true;
};

const handleEditSubmit = async () => {
  try {
    if (!editForm.fiat || !editForm.crypto) {
      Message.error("交易对信息不完整");
      return;
    }

    if (syntaxValue.value === undefined || syntaxValue.value === null) {
      Message.error("请输入有效的数值");
      return;
    }

    if (syntaxType.value === "~") {
      if (syntaxValue.value <= 0) {
        Message.error("百分比浮动数值必须大于 0");
        return;
      }
    } else if (syntaxValue.value < 0) {
      Message.error("数值不能为负数");
      return;
    }

    editLoading.value = true;
    const syntax = generateSyntax();

    await setSyntaxAPI({
      fiat: editForm.fiat,
      crypto: editForm.crypto,
      syntax: syntax
    });

    Message.success("编辑成功");
    editModalVisible.value = false;
    await getCommonTableList();
  } catch (error) {
    console.error("编辑失败:", error);
    Message.error("编辑失败");
  } finally {
    editLoading.value = false;
  }
};

const handleEditCancel = () => {
  editModalVisible.value = false;
  editFormRef.value?.resetFields();
  syntaxType.value = "";
  syntaxValue.value = undefined;
};

// 同步频率相关状态
const syncModalVisible = ref<boolean>(false);
const syncLoading = ref<boolean>(false);
const syncFormRef = ref();

const syncForm = reactive({
  minutes: 60,
  apiUrl: "https://api.coingecko.com",
  apiKey: "",
  historyDays: 30
});

// API 接口选项
const apiUrlOptions = [
  {
    label: "官方接口 免费额度存在速率限制",
    value: "https://api.coingecko.com"
  },
  {
    label: "开源免费 作者自建 没有速率限制",
    value: "https://api-coingecko-com.bepusdt.online"
  }
];

// 显示同步频率模态框
const showSyncModal = async () => {
  try {
    const res = await getsConfAPI({
      keys: ["rate_sync_interval", "rate_sync_coingecko_api_url", "rate_sync_coingecko_api_key", "rate_sync_history_days"]
    });

    if (res.data) {
      if (res.data.rate_sync_interval) {
        const seconds = parseInt(res.data.rate_sync_interval);
        const minutes = Math.round(seconds / 60);
        syncForm.minutes = minutes;
      } else {
        syncForm.minutes = 60;
      }

      syncForm.apiUrl = res.data.rate_sync_coingecko_api_url || "https://api.coingecko.com";
      syncForm.apiKey = res.data.rate_sync_coingecko_api_key || "";
      syncForm.historyDays = res.data.rate_sync_history_days ? parseInt(res.data.rate_sync_history_days) : 30;
    } else {
      syncForm.minutes = 60;
      syncForm.apiUrl = "https://api.coingecko.com";
      syncForm.apiKey = "";
      syncForm.historyDays = 30;
    }
  } catch (error) {
    console.error("获取同步频率配置失败:", error);
    syncForm.minutes = 60;
    syncForm.apiUrl = "https://api.coingecko.com";
    syncForm.apiKey = "";
    syncForm.historyDays = 30;
    Message.warning("获取当前配置失败，使用默认值");
  }

  syncModalVisible.value = true;
};

const handleSyncSubmit = async () => {
  try {
    if (!syncForm.minutes || syncForm.minutes < 10 || syncForm.minutes > 1440) {
      Message.error("请输入有效的同步频率（10-1440分钟）");
      return;
    }

    if (!syncForm.apiUrl) {
      Message.error("请选择 API 接口");
      return;
    }

    if (!syncForm.historyDays || syncForm.historyDays < 1 || syncForm.historyDays > 365) {
      Message.error("请输入有效的汇率保留天数（1-365天）");
      return;
    }

    syncLoading.value = true;
    const seconds = syncForm.minutes * 60;

    await setsConfAPI([
      { key: "rate_sync_interval", value: seconds.toString() },
      { key: "rate_sync_coingecko_api_url", value: syncForm.apiUrl },
      { key: "rate_sync_coingecko_api_key", value: syncForm.apiKey },
      { key: "rate_sync_history_days", value: syncForm.historyDays.toString() }
    ]);

    Message.success("汇率同步配置已保存");
    syncModalVisible.value = false;
  } catch (error) {
    console.error("设置同步配置失败:", error);
    Message.error("设置失败");
  } finally {
    syncLoading.value = false;
  }
};

const handleSyncCancel = () => {
  syncModalVisible.value = false;
  syncFormRef.value?.resetFields();
  syncForm.minutes = 60;
  syncForm.apiUrl = "https://api.coingecko.com";
  syncForm.apiKey = "";
  syncForm.historyDays = 30;
};

// 支付颗粒度相关状态
const atomModalVisible = ref<boolean>(false);
const atomLoading = ref<boolean>(false);
const atomFormRef = ref();

const atomForm = reactive({
  usdt: 0.01,
  usdc: 0.01,
  trx: 0.01,
  eth: 0.000001,
  bnb: 0.00001
});

const showAtomModal = async () => {
  try {
    const res = await getsConfAPI({
      keys: ["atom_usdt", "atom_usdc", "atom_trx", "atom_eth", "atom_bnb"]
    });

    if (res.data) {
      console.log(res.data);
      atomForm.usdt = res.data.atom_usdt ? parseFloat(res.data.atom_usdt) : 0.01;
      atomForm.usdc = res.data.atom_usdc ? parseFloat(res.data.atom_usdc) : 0.01;
      atomForm.trx = res.data.atom_trx ? parseFloat(res.data.atom_trx) : 0.01;
      atomForm.eth = res.data.atom_eth ? parseFloat(res.data.atom_eth) : 0.000001;
      atomForm.bnb = res.data.atom_bnb ? parseFloat(res.data.atom_bnb) : 0.00001;
    }
  } catch (error) {
    console.error("获取支付颗粒度配置失败:", error);
    Message.warning("获取当前配置失败，使用默认值");
  }

  atomModalVisible.value = true;
};

const handleAtomSubmit = async () => {
  try {
    if (!atomForm.usdt || !atomForm.usdc || !atomForm.trx || !atomForm.eth || !atomForm.bnb) {
      Message.error("请填写所有颗粒度配置");
      return;
    }

    atomLoading.value = true;

    await setsConfAPI([
      { key: "atom_usdt", value: atomForm.usdt.toString() },
      { key: "atom_usdc", value: atomForm.usdc.toString() },
      { key: "atom_trx", value: atomForm.trx.toString() },
      { key: "atom_eth", value: atomForm.eth.toString() },
      { key: "atom_bnb", value: atomForm.bnb.toString() }
    ]);

    Message.success("支付颗粒度设置成功");
    atomModalVisible.value = false;
  } catch (error) {
    console.error("设置支付颗粒度失败:", error);
    Message.error("设置失败");
  } finally {
    atomLoading.value = false;
  }
};

const handleAtomCancel = () => {
  atomModalVisible.value = false;
  atomFormRef.value?.resetFields();
  atomForm.usdt = 0.01;
  atomForm.usdc = 0.01;
  atomForm.trx = 0.01;
  atomForm.eth = 0.000001;
  atomForm.bnb = 0.00001;
};

getCommonTableList();
</script>

<style lang="scss" scoped>
.fiat-display {
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.syntax-display {
  display: flex;
  align-items: center;
  gap: 12px;
  white-space: nowrap;

  .syntax-value {
    font-weight: 500;
    color: $color-text-1;
    min-width: 80px;
    text-align: left;
    flex-shrink: 0;
  }

  .syntax-description {
    font-size: 12px;
    color: $color-text-3;
    font-style: italic;
    flex-shrink: 0;
  }
}

.syntax-prefix {
  color: $color-primary;
  font-weight: bold;
}

.edit-modal {
  :deep(.arco-modal-body) {
    padding: 16px 24px;
  }

  .syntax-tip {
    padding: 8px 12px;
    background: $color-fill-1;
    border: 1px solid $color-border-2;
    border-radius: 4px;
    font-size: 12px;
    margin-top: 8px;
  }
}

.toolbar {
  margin-bottom: 16px;
  display: flex;
  justify-content: flex-start; // 改为左对齐
}

.sync-modal {
  :deep(.arco-modal-body) {
    padding: 16px 24px;
  }

  .sync-tip {
    padding: 6px 10px;
    background: $color-fill-1;
    border: 1px solid $color-border-2;
    border-radius: 4px;
    font-size: 11px;
    line-height: 1.4;
    margin-top: 8px;

    .sync-warning {
      color: $color-danger;
    }

    hr {
      margin: 6px 0 0 0;
      border: none;
      border-top: 1px solid $color-border-2;
    }
  }
}

.atom-modal {
  :deep(.arco-modal-body) {
    padding: 16px 24px;
  }

  .atom-tip {
    padding: 8px 12px;
    background: $color-fill-1;
    border: 1px solid $color-border-2;
    border-radius: 4px;
    font-size: 12px;
    margin-top: 8px;
  }
}
</style>
