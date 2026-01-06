<template>
  <div class="snow-page">
    <div class="snow-inner">
      <a-form ref="formRef" auto-label-width :model="formData.form">
        <a-row :gutter="16">
          <a-col :xs="24" :sm="24" :md="12" :lg="12" :xl="6" :xxl="6">
            <a-form-item field="name" label="姓名">
              <a-input v-model="formData.form.name" placeholder="请输入姓名" allow-clear />
            </a-form-item>
          </a-col>
          <a-col :xs="24" :sm="24" :md="12" :lg="12" :xl="6" :xxl="6">
            <a-form-item field="phone" label="手机号">
              <a-input v-model="formData.form.phone" placeholder="请输入手机号" allow-clear />
            </a-form-item>
          </a-col>
          <a-col :xs="24" :sm="24" :md="12" :lg="12" :xl="6" :xxl="6">
            <a-form-item field="email" label="邮箱">
              <a-input v-model="formData.form.email" placeholder="请输入邮箱" allow-clear />
            </a-form-item>
          </a-col>
          <a-col :xs="24" :sm="24" :md="12" :lg="12" :xl="6" :xxl="6">
            <a-space class="search-btn">
              <a-button type="primary" @click="getCommonTableList">
                <template #icon>
                  <icon-search />
                </template>
                <template #default>查询</template>
              </a-button>
              <a-button @click="onReset">
                <template #icon>
                  <icon-refresh />
                </template>
                <template #default>重置</template>
              </a-button>
              <a-button type="text" @click="formData.search = !formData.search">
                <template #icon>
                  <icon-up v-if="formData.search" />
                  <icon-down v-else />
                </template>
                <template #default>{{ formData.search ? "收起" : "展开" }}</template>
              </a-button>
            </a-space>
          </a-col>
        </a-row>
        <a-row :gutter="16" v-if="formData.search">
          <a-col :xs="24" :sm="24" :md="12" :lg="12" :xl="6" :xxl="6">
            <a-form-item field="address" label="地址">
              <a-input v-model="formData.form.address" placeholder="请输入地址" allow-clear />
            </a-form-item>
          </a-col>
          <a-col :xs="24" :sm="24" :md="12" :lg="12" :xl="6" :xxl="6">
            <a-form-item field="status" label="用户状态">
              <a-select v-model="formData.form.status" placeholder="请选择用户状态" allow-clear>
                <a-option value="1">停用</a-option>
                <a-option value="2">启用</a-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>
      </a-form>
      <a-table
        row-key="key"
        size="small"
        :bordered="{
          cell: true
        }"
        :scroll="{ x: '100%', y: '100%', minWidth: 1000 }"
        :loading="loading"
        :columns="columns"
        :data="data"
        :row-selection="rowSelection"
        v-model:selectedKeys="selectedKeys"
        :pagination="pagination"
        @page-change="pageChange"
        @page-size-change="pageSizeChange"
      >
        <template #avatar="{ record }">
          <a-avatar
            auto-fix-font-size
            :size="30"
            :style="{
              backgroundColor: '#14a9f8'
            }"
          >
            {{ record.avatar }}
          </a-avatar>
        </template>
        <template #status="{ record }">
          <a-space>
            <a-tag size="small" color="green" v-if="record.status == 1">启用</a-tag>
            <a-tag size="small" color="red" v-else>停用</a-tag>
          </a-space>
        </template>
        <template #optional>
          <a-space>
            <a-button size="mini" type="primary">详情</a-button>
            <a-button size="mini">修改</a-button>
            <a-popconfirm content="确定删除这条数据吗?" type="warning">
              <a-button size="mini" type="primary" status="danger">删除</a-button>
            </a-popconfirm>
          </a-space>
        </template>
      </a-table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { getCommonTableListAPI } from "@/api/modules/table/index";
import { List, FormData, RowSelection, Pagination } from "./config";
const formData = reactive<FormData>({
  form: {
    name: "",
    phone: "",
    email: "",
    address: "",
    status: null
  },
  search: false
});
const selectedKeys = ref<string[]>([]);
const rowSelection = reactive<RowSelection>({
  type: "checkbox",
  showCheckedAll: true,
  onlyCurrent: false
});
const pagination = ref<Pagination>({ showPageSize: true, showTotal: true, current: 1, pageSize: 10, total: 10 });
const pageChange = (page: number) => {
  pagination.value.current = page;
};
const pageSizeChange = (pageSize: number) => {
  pagination.value.pageSize = pageSize;
};
const columns = [
  {
    title: "姓名",
    dataIndex: "name"
  },
  {
    title: "头像",
    dataIndex: "avatar",
    slotName: "avatar",
    align: "center"
  },
  {
    title: "手机号",
    dataIndex: "phone"
  },
  {
    title: "Email",
    dataIndex: "email"
  },
  {
    title: "居住地址",
    dataIndex: "address"
  },
  {
    title: "用户状态",
    dataIndex: "status",
    align: "center",
    slotName: "status"
  },
  {
    title: "创建时间",
    dataIndex: "createTime"
  },
  {
    title: "操作",
    slotName: "optional",
    align: "center",
    fixed: "right",
    width: "200"
  }
];
const formRef = ref();
const onReset = () => {
  formRef.value.resetFields();
  getCommonTableList();
};
const loading = ref<boolean>(false);
const data = reactive<List[]>([]);
const getCommonTableList = async () => {
  try {
    loading.value = true;
    let res = await getCommonTableListAPI();
    Object.assign(data, res.data.list);
    pagination.value.total = res.data.total;
  } finally {
    loading.value = false;
  }
};
getCommonTableList();
</script>

<style lang="scss" scoped>
.search-btn {
  margin-bottom: 20px;
}
</style>
