<template>
  <div class="snow-page">
    <div class="snow-inner">
      <s-layout-tools>
        <template #left>
          <a-space wrap>
            <a-input v-model="form.name" placeholder="请输入部门名称" allow-clear />
            <a-select placeholder="部门状态" v-model="form.status" style="width: 120px" allow-clear>
              <a-option v-for="item in openState" :key="item.value" :value="item.value">{{ item.name }}</a-option>
            </a-select>
            <a-button type="primary" @click="search">
              <template #icon><icon-search /></template>
              <span>查询</span>
            </a-button>
            <a-button @click="reset">
              <template #icon><icon-refresh /></template>
              <span>重置</span>
            </a-button>
          </a-space>
        </template>
        <template #right>
          <a-space wrap>
            <a-button type="primary" @click="onAdd">
              <template #icon><icon-plus /></template>
              <span>新增</span>
            </a-button>
            <a-button type="primary" status="danger">
              <template #icon><icon-delete /></template>
              <span>删除</span>
            </a-button>
          </a-space>
        </template>
      </s-layout-tools>
      <a-table
        ref="tableRef"
        :data="tableData"
        default-expand-all-rows
        :bordered="{ cell: true }"
        row-key="id"
        :loading="loading"
      >
        <template #columns>
          <a-table-column title="部门名称">
            <template #cell="{ record }">
              {{ record.name }}
            </template>
          </a-table-column>
          <a-table-column title="排序" data-index="sort" :width="100" align="center"></a-table-column>
          <a-table-column title="启用状态" align="center">
            <template #cell="{ record }">
              <a-space>
                <a-tag bordered size="small" color="arcoblue" v-if="record.status">启用</a-tag>
                <a-tag bordered size="small" color="red" v-else>禁用</a-tag>
              </a-space>
            </template>
          </a-table-column>
          <a-table-column title="描述" data-index="description" :ellipsis="true" :tooltip="true"></a-table-column>
          <a-table-column title="创建时间" data-index="createTime" :ellipsis="true" :tooltip="true"></a-table-column>
          <a-table-column title="操作" align="center" :fixed="'right'">
            <template #cell="{ record }">
              <a-space>
                <a-button size="mini" type="primary" @click="onUpdate(record)">
                  <template #icon><icon-edit /></template>
                  <span>修改</span>
                </a-button>
                <a-button size="mini" type="primary" status="success" @click="addDivision(record.id)">
                  <template #icon><icon-plus /></template>
                  <span>新增</span>
                </a-button>
                <a-popconfirm type="warning" content="确定删除该部门吗?">
                  <a-button size="mini" type="primary" status="danger" v-if="record.id != '100'">
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
          <a-form-item field="parentId" label="上级部门" validate-trigger="blur">
            <a-tree-select
              v-model="addFrom.parentId"
              :data="tableData"
              :field-names="{
                key: 'id',
                title: 'name',
                children: 'children'
              }"
              placeholder="选择上级部门"
            ></a-tree-select>
          </a-form-item>
          <a-row :gutter="24">
            <a-col :span="12">
              <a-form-item field="name" label="部门名称" validate-trigger="blur">
                <a-input v-model="addFrom.name" placeholder="请输入部门名称" allow-clear />
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item field="status" label="部门状态" validate-trigger="blur">
                <a-switch type="round" :checked-value="1" :unchecked-value="0" v-model="addFrom.status">
                  <template #checked> 启用 </template>
                  <template #unchecked> 禁用 </template>
                </a-switch>
              </a-form-item>
            </a-col>
          </a-row>
          <a-row :gutter="24">
            <a-col :span="12">
              <a-form-item field="leader" label="负责人" validate-trigger="blur">
                <a-input v-model="addFrom.leader" placeholder="请输入负责人" allow-clear />
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item field="phone" label="联系电话" validate-trigger="blur">
                <a-input v-model="addFrom.phone" placeholder="请输入联系电话" allow-clear />
              </a-form-item>
            </a-col>
          </a-row>
          <a-row :gutter="24">
            <a-col :span="12">
              <a-form-item field="email" label="邮箱" validate-trigger="blur">
                <a-input v-model="addFrom.email" placeholder="请输入邮箱" allow-clear />
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item field="sort" label="显示排序" validate-trigger="blur">
                <a-input-number
                  v-model="addFrom.sort"
                  placeholder="请输入排序值"
                  :min="1"
                  :max="999"
                  :step="1"
                  :precision="0"
                  allow-clear
                />
              </a-form-item>
            </a-col>
          </a-row>
          <a-form-item field="description" label="描述" validate-trigger="blur">
            <a-textarea v-model="addFrom.description" placeholder="请输入描述" allow-clear />
          </a-form-item>
        </a-form>
      </div>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { getDivisionAPI } from "@/api/modules/system/index";
import { deepClone } from "@/utils";

// 新增
const open = ref(false);
const rules = {
  parentId: [
    {
      required: true,
      message: "请选择上级部门"
    }
  ],
  name: [
    {
      required: true,
      message: "请输入部门名称"
    }
  ],
  status: [
    {
      required: true,
      message: "请选择部门状态"
    }
  ]
};
const addFrom = ref<any>({
  parentId: null,
  name: "",
  sort: 0,
  leader: "",
  phone: "",
  email: "",
  status: 1,
  description: ""
});
const formType = ref(0); // 0新增 1修改 2新增下级
const title = ref("");
const formRef = ref();
const onAdd = () => {
  title.value = "添加部门";
  formType.value = 0;
  open.value = true;
};
const handleOk = async () => {
  let state = await formRef.value.validate();
  if (state) return (open.value = true); // 校验不通过
  console.log("模拟提交", formType.value, addFrom.value);
  arcoMessage("success", "模拟提交成功");
  getDivision();
};
// 关闭对话框动画结束后触发
const afterClose = () => {
  formRef.value.resetFields();
  addFrom.value = {
    parentId: null,
    name: "",
    sort: 0,
    leader: "",
    phone: "",
    email: "",
    status: 1,
    description: ""
  };
};
const onUpdate = (row: any) => {
  title.value = "修改部门";
  formType.value = 1;
  addFrom.value = deepClone(row);
  row.parentId == 0 && (addFrom.value.parentId = null);
  open.value = true;
};
const addDivision = (id: any) => {
  title.value = "新增部门";
  formType.value = 2;
  addFrom.value.parentId = id == 0 ? null : id;
  open.value = true;
};

const openState = ref(dictFilter("status"));
const form = ref({
  name: "",
  status: ""
});
const search = () => {
  getDivision();
};
const reset = () => {
  form.value = {
    name: "",
    status: ""
  };
  getDivision();
};
const loading = ref(false);
const tableRef = ref();
const tableData = ref();
const getDivision = async () => {
  loading.value = true;
  let res = await getDivisionAPI();
  tableData.value = res.data;
  loading.value = false;
  setTimeout(() => {
    tableRef.value.expandAll();
  }, 0);
};

onMounted(() => {
  getDivision();
});
</script>

<style lang="scss" scoped></style>
