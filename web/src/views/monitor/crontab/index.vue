<template>
  <div class="snow-page">
    <div class="snow-inner">
      <a-space wrap>
        <a-input v-model="form.loginLocation" placeholder="请输入任务名称" allow-clear />
        <a-select placeholder="任务状态" v-model="form.status" style="width: 120px" allow-clear>
          <a-option v-for="item in openState" :key="item.value" :value="item.value">{{ item.name }}</a-option>
        </a-select>
        <a-range-picker v-model="form.createTime" show-time format="YYYY-MM-DD HH:mm" allow-clear />
        <a-button type="primary" @click="search">
          <template #icon><icon-search /></template>
          <span>查询</span>
        </a-button>
        <a-button @click="reset">
          <template #icon><icon-refresh /></template>
          <span>重置</span>
        </a-button>
      </a-space>

      <a-row>
        <a-space wrap>
          <a-button type="primary" @click="onAdd">
            <template #icon><icon-plus /></template>
            <span>新增</span>
          </a-button>
        </a-space>
      </a-row>

      <a-table
        row-key="id"
        :data="list"
        :bordered="{ cell: true }"
        :loading="loading"
        :scroll="{ x: '100%', y: '100%', minWidth: 1000 }"
        :pagination="pagination"
      >
        <template #columns>
          <a-table-column title="序号" :width="64">
            <template #cell="cell">{{ cell.rowIndex + 1 }}</template>
          </a-table-column>
          <a-table-column title="任务名称" data-index="name" ellipsis tooltip></a-table-column>
          <a-table-column title="执行服务" data-index="service" ellipsis tooltip></a-table-column>
          <a-table-column title="定时规则" data-index="taskType" ellipsis tooltip>
            <template #cell="{ record }">
              <span v-if="record.taskType == 0">CRON: {{ record.cron }}</span>
              <span v-else>间隔时间: {{ record.every }} 秒</span>
            </template>
          </a-table-column>
          <a-table-column title="状态" align="center" data-index="status" :width="80">
            <template #cell="{ record }">
              <a-tag bordered size="small" color="arcoblue" v-if="record.status === 1">启用</a-tag>
              <a-tag bordered size="small" color="red" v-else>禁用</a-tag>
            </template>
          </a-table-column>
          <a-table-column title="备注" data-index="remark" ellipsis tooltip></a-table-column>
          <a-table-column title="创建时间" data-index="createTime" :width="180"></a-table-column>
          <a-table-column title="操作" :width="250" align="center" :fixed="'right'">
            <template #cell="{ record }">
              <a-space>
                <a-button type="primary" status="success" size="mini" @click="onLogs(record)">
                  <template #icon><icon-file /></template>
                  <span>日志</span>
                </a-button>
                <a-button type="primary" size="mini" @click="onUpdate(record)">
                  <template #icon><icon-edit /></template>
                  <span>修改</span>
                </a-button>
                <a-popconfirm type="warning" content="确定删除吗?">
                  <a-button type="primary" status="danger" size="mini">
                    <template #icon><icon-delete /></template>
                    <span>删除</span>
                  </a-button>
                </a-popconfirm>
              </a-space>
            </template>
          </a-table-column>
        </template>
      </a-table>
    </div>

    <a-modal width="40%" v-model:visible="open" @close="afterClose" @ok="handleOk" @cancel="afterClose">
      <template #title> {{ title }} </template>
      <div>
        <a-form ref="formRef" auto-label-width :rules="rules" :model="addFrom">
          <a-form-item field="name" label="任务名称" validate-trigger="blur">
            <a-input v-model="addFrom.name" placeholder="请输入任务名称" allow-clear />
          </a-form-item>
          <a-form-item field="status" label="任务状态" validate-trigger="blur">
            <a-switch type="round" :checked-value="1" :unchecked-value="0" v-model="addFrom.status">
              <template #checked> 启用 </template>
              <template #unchecked> 禁用 </template>
            </a-switch>
          </a-form-item>
          <a-row :gutter="24">
            <a-col :span="12">
              <a-form-item field="taskType" label="任务类型" validate-trigger="blur">
                <a-select placeholder="请选择执行策略" v-model="addFrom.taskType" allow-clear @change="onTaskTypeChange">
                  <a-option v-for="item in taskTypeOption" :key="item.value" :value="item.value">{{ item.name }}</a-option>
                </a-select>
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item field="cron" label="Cron表达式" validate-trigger="blur" v-if="addFrom.taskType == 0">
                <a-input v-model="addFrom.cron" placeholder="******" allow-clear />
              </a-form-item>
              <a-form-item field="every" label="间隔时长" validate-trigger="blur" v-else>
                <a-input-number
                  v-model="addFrom.every"
                  placeholder="请输入"
                  :min="1"
                  :max="999999999"
                  :step="1"
                  :precision="0"
                  allow-clear
                >
                  <template #suffix> 秒 </template>
                </a-input-number>
              </a-form-item>
            </a-col>
          </a-row>

          <a-form-item field="service" label="执行服务" validate-trigger="blur">
            <a-input v-model="addFrom.service" placeholder="请输入执行服务service" allow-clear />
          </a-form-item>
          <a-form-item field="startDate" label="开始时间" validate-trigger="blur">
            <a-date-picker v-model="addFrom.startDate" show-time allow-clear format="YYYY-MM-DD HH:mm:ss" />
          </a-form-item>
          <a-form-item field="misfirePolicy" label="执行策略" validate-trigger="blur">
            <a-select placeholder="请选择执行策略" v-model="addFrom.misfirePolicy" allow-clear>
              <a-option v-for="item in misfirePolicyOption" :key="item.value" :value="item.value">{{ item.name }}</a-option>
            </a-select>
          </a-form-item>
          <a-form-item field="remark" label="备注" validate-trigger="blur">
            <a-textarea v-model="addFrom.remark" placeholder="请输入描述" allow-clear />
          </a-form-item>
        </a-form>
      </div>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { getCrontabAPI } from "@/api/modules/monitor/index";
import { deepClone } from "@/utils";
defineOptions({ name: "crontab" });
const router = useRouter();
const openState = ref(dictFilter("status"));
const misfirePolicyOption = ref([
  { name: "循环", value: 1 },
  { name: "执行一次", value: 2 }
]);
const taskTypeOption = ref([
  { name: "cron", value: 0 },
  { name: "时间间隔", value: 1 }
]);
const form = ref({
  loginLocation: "",
  userName: "",
  status: null,
  createTime: []
});
const search = () => {
  getCrontab();
};
const reset = () => {
  form.value = {
    loginLocation: "",
    userName: "",
    status: null,
    createTime: []
  };
  getCrontab();
};

// 新增
const open = ref(false);
const rules = {
  name: [{ required: true, message: "请输入任务名称" }],
  status: [{ required: true, message: "请选择任务状态" }],
  taskType: [{ required: true, message: "请选择任务类型" }],
  cron: [{ required: true, message: "请输入Cron表达式" }],
  every: [{ required: true, message: "请输入间隔时间(秒)" }],
  service: [{ required: true, message: "请输入执行服务service" }],
  startDate: [{ required: true, message: "请选择开始时间" }],
  misfirePolicy: [{ required: true, message: "请选择执行策略" }]
};
const addFrom = ref<any>({
  name: "",
  status: 1,
  taskType: 0,
  cron: "",
  every: null,
  service: "",
  startDate: "",
  misfirePolicy: 1,
  remark: ""
});
const title = ref("");
const formRef = ref();
const onAdd = () => {
  title.value = "添加任务";
  open.value = true;
};
const handleOk = async () => {
  let state = await formRef.value.validate();
  if (state) return (open.value = true); // 校验不通过
  arcoMessage("success", "模拟提交成功");
  getCrontab();
};
// 关闭对话框动画结束后触发
const afterClose = () => {
  formRef.value.resetFields();
  addFrom.value = {
    name: "",
    status: 1,
    taskType: 0,
    cron: "",
    every: null,
    service: "",
    startDate: "",
    misfirePolicy: 1,
    remark: ""
  };
};
const onTaskTypeChange = (value: number) => {
  if (value == 0) return (addFrom.value.every = null);
  addFrom.value.cron = "";
  addFrom.value.startDate = "";
};
const onUpdate = (row: any) => {
  title.value = "修改任务";
  addFrom.value = deepClone(row);
  open.value = true;
};

const onLogs = (row: any) => {
  router.push({
    path: "/monitor/crontab-logs",
    query: {
      id: row.id,
      name: row.name
    }
  });
};

// 获取列表
const loading = ref(false);
const pagination = ref({
  pageSize: 10,
  showPageSize: true
});
const list = ref([]);
const getCrontab = async () => {
  try {
    loading.value = true;
    let res = await getCrontabAPI();
    list.value = res.data;
  } finally {
    loading.value = false;
  }
};

getCrontab();
</script>

<style lang="scss" scoped>
.text-right-gap {
  margin-right: $margin;
}
</style>
