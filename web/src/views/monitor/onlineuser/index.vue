<template>
  <div class="snow-page">
    <div class="snow-inner">
      <a-space wrap>
        <a-input v-model="form.loginLocation" placeholder="请输入登录地址" allow-clear />
        <a-input v-model="form.userName" placeholder="请输入账户名称" allow-clear />
        <a-range-picker v-model="form.loginTime" show-time format="YYYY-MM-DD HH:mm" allow-clear />
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
          <a-table-column title="会话编号" data-index="tokenId" ellipsis tooltip></a-table-column>
          <a-table-column title="登录账户" data-index="userName" ellipsis tooltip></a-table-column>
          <a-table-column title="部门名称" data-index="deptName" ellipsis tooltip></a-table-column>
          <a-table-column title="IP地址" data-index="ipaddr" ellipsis tooltip></a-table-column>
          <a-table-column title="登录地址" data-index="loginLocation" ellipsis tooltip></a-table-column>
          <a-table-column title="状态" data-index="status" align="center" :width="80">
            <template #cell="{ record }">
              <a-space>
                <a-badge status="success" text="在线" v-if="record.status == 1" />
                <a-badge status="normal" text="离线" v-else />
              </a-space>
            </template>
          </a-table-column>
          <a-table-column title="浏览器" data-index="browser" ellipsis tooltip></a-table-column>
          <a-table-column title="操作系统" data-index="os" ellipsis tooltip></a-table-column>
          <a-table-column title="登录时间" data-index="loginTime" ellipsis tooltip></a-table-column>
          <a-table-column title="操作" :width="100" align="center" :fixed="'right'">
            <template #cell="{ record }">
              <a-space>
                <a-popconfirm type="warning" content="确定强制退出该账号吗?" @ok="onLogout(record)">
                  <a-button type="primary" status="danger" size="mini">
                    <template #icon><icon-export /></template>
                    <span>强退</span>
                  </a-button>
                </a-popconfirm>
              </a-space>
            </template>
          </a-table-column>
        </template>
      </a-table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { getOnlineuserAPI } from "@/api/modules/monitor/index";
defineOptions({ name: "onlineuser" });
const form = ref({
  loginLocation: "",
  userName: "",
  loginTime: []
});
const search = () => {
  getOnlineuser();
};
const reset = () => {
  form.value = {
    loginLocation: "",
    userName: "",
    loginTime: []
  };
  getOnlineuser();
};

const onLogout = (row: any) => {
  console.log("退出", row);
  arcoMessage("success", "模拟退出成功");
  getOnlineuser();
};

// 获取列表
const loading = ref(false);
const pagination = ref({
  pageSize: 10,
  showPageSize: true
});
const list = ref([]);
const getOnlineuser = async () => {
  try {
    loading.value = true;
    let res = await getOnlineuserAPI();
    list.value = res.data;
  } finally {
    loading.value = false;
  }
};

getOnlineuser();
</script>

<style lang="scss" scoped>
.text-right-gap {
  margin-right: $margin;
}
</style>
