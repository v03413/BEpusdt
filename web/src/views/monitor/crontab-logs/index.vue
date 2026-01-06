<template>
  <div class="snow-page">
    <div class="snow-inner">
      <a-space wrap>
        <a-input v-model="form.name" placeholder="请输入任务名称" allow-clear />
        <a-select placeholder="执行状态" v-model="form.status" style="width: 120px" allow-clear>
          <a-option v-for="item in openState" :key="item.value" :value="item.value">{{ item.name }}</a-option>
        </a-select>
        <a-range-picker v-model="form.datetimeRange" show-time format="YYYY-MM-DD HH:mm" allow-clear />

        <a-button type="primary" @click="search">
          <template #icon><icon-search /></template>
          <span>查询</span>
        </a-button>
        <a-button @click="reset">
          <template #icon><icon-refresh /></template>
          <span>重置</span>
        </a-button>
      </a-space>

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
          <a-table-column title="任务名称" data-index="taskName" ellipsis tooltip></a-table-column>
          <a-table-column title="执行状态" data-index="status" align="center" :width="100">
            <template #cell="{ record }">
              <a-tag bordered size="small" color="arcoblue" v-if="record.status === 1">成功</a-tag>
              <a-tag bordered size="small" color="red" v-else>失败</a-tag>
            </template>
          </a-table-column>
          <a-table-column title="备注" data-index="remark" ellipsis tooltip>
            <template #cell="{ record }">
              {{ record.remark || "无" }}
            </template>
          </a-table-column>
          <a-table-column title="预设时间" data-index="createTime"></a-table-column>
          <a-table-column title="执行时间" data-index="executionTime"></a-table-column>
          <a-table-column title="操作" :width="100" align="center" :fixed="'right'">
            <template #cell="{ record }">
              <a-space>
                <a-button type="primary" status="success" size="mini" @click="onDetail(record)">
                  <template #icon><icon-file /></template>
                  <span>详情</span>
                </a-button>
              </a-space>
            </template>
          </a-table-column>
        </template>
      </a-table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { getCrontabLogsAPI } from "@/api/modules/monitor/index";
import { useRouteConfigStore } from "@/store/modules/route-config";
defineOptions({ name: "crontab-logs" });
let route = useRoute();
const routerStore = useRouteConfigStore();
const openState = ref(dictFilter("taskStatus"));

const form = ref({
  name: "",
  status: null,
  datetimeRange: []
});
const search = () => {
  getCrontabLogs();
};
const reset = () => {
  form.value = {
    name: "",
    status: null,
    datetimeRange: []
  };
  getCrontabLogs();
};

const onDetail = (record: any) => {
  console.log("详情", record);
};

// 获取列表
const loading = ref(false);
const pagination = ref({
  pageSize: 10,
  showPageSize: true
});
const list = ref([]);
const getCrontabLogs = async () => {
  try {
    loading.value = true;
    let res = await getCrontabLogsAPI(form.value);
    list.value = res.data;
  } finally {
    loading.value = false;
  }
};

const init = () => {
  if (route.query.id) {
    form.value.name = route.query.name;
  }
  getCrontabLogs();
};

init();

routerStore.setTabsTitle(`任务日志${route.query.id ? " - " + route.query.id : ""}`);
</script>

<style lang="scss" scoped>
.text-right-gap {
  margin-right: $margin;
}
</style>
