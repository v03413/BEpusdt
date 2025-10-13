<template>
  <div class="snow-page">
    <div class="snow-inner">
      <s-layout-tools>
        <template #left>
          <a-space wrap>
            <a-input v-model="form.name" placeholder="请输入角色名称" allow-clear />
            <a-input v-model="form.code" placeholder="请输入角色标识" allow-clear />
            <a-select placeholder="角色状态" v-model="form.status" style="width: 120px" allow-clear>
              <a-option v-for="item in openState" :key="item.value" :value="item.value">{{ item.name }}</a-option>
            </a-select>
            <a-range-picker v-model="form.time" show-time format="YYYY-MM-DD HH:mm" allow-clear />
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
        row-key="id"
        :data="roleList"
        :bordered="{ cell: true }"
        :loading="loading"
        :scroll="{ x: '100%', y: '100%', minWidth: 1000 }"
        :pagination="pagination"
        :row-selection="{ type: 'checkbox', showCheckedAll: true }"
        :selected-keys="selectedKeys"
        @select="select"
        @select-all="selectAll"
      >
        <template #columns>
          <a-table-column title="序号" :width="64">
            <template #cell="cell">{{ cell.rowIndex + 1 }}</template>
          </a-table-column>
          <a-table-column title="角色名称" data-index="name"></a-table-column>
          <a-table-column title="角色标识" data-index="code"></a-table-column>
          <a-table-column title="排序" data-index="sort" :width="100" align="center"></a-table-column>
          <a-table-column title="状态" :width="100" align="center">
            <template #cell="{ record }">
              <a-tag bordered size="small" color="arcoblue" v-if="record.status === 1">启用</a-tag>
              <a-tag bordered size="small" color="red" v-else>禁用</a-tag>
            </template>
          </a-table-column>
          <a-table-column title="描述" data-index="description" :ellipsis="true" :tooltip="true"></a-table-column>
          <a-table-column title="创建时间" data-index="createTime" :width="180"></a-table-column>
          <a-table-column title="操作" :width="280" align="center" :fixed="'right'">
            <template #cell="{ record }">
              <a-space>
                <a-button type="primary" status="success" size="mini" :disabled="record.admin" @click="onPrivileges(record)">
                  <template #icon><icon-safe /></template>
                  <span>分配权限</span>
                </a-button>
                <a-button type="primary" size="mini" :disabled="record.admin" @click="onUpdate(record)">
                  <template #icon><icon-edit /></template>
                  <span>修改</span>
                </a-button>
                <a-popconfirm type="warning" content="确定删除该角色吗?">
                  <a-button type="primary" status="danger" size="mini" :disabled="record.admin">
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
          <a-form-item field="name" label="角色名称" validate-trigger="blur">
            <a-input v-model="addFrom.name" placeholder="请输入角色名称" allow-clear />
          </a-form-item>
          <a-form-item field="code" label="角色编码" validate-trigger="blur">
            <a-input v-model="addFrom.code" placeholder="请输入角色编码" allow-clear />
          </a-form-item>
          <a-form-item field="status" label="状态" validate-trigger="blur">
            <a-switch type="round" :checked-value="1" :unchecked-value="0" v-model="addFrom.status">
              <template #checked> 启用 </template>
              <template #unchecked> 禁用 </template>
            </a-switch>
          </a-form-item>
          <a-form-item field="sort" label="排序" validate-trigger="blur">
            <a-input-number
              v-model="addFrom.sort"
              :step="1"
              :precision="0"
              :min="1"
              :max="9999"
              :style="{ width: '150px' }"
              placeholder="请输入"
              mode="button"
              class="input-demo"
            />
          </a-form-item>
          <a-form-item field="description" label="描述" validate-trigger="blur">
            <a-textarea v-model="addFrom.description" placeholder="请输入描述" allow-clear />
          </a-form-item>
        </a-form>
      </div>
    </a-modal>

    <a-drawer :visible="drawerOpen" :width="500" @ok="drawerOk" @cancel="drawerCancel">
      <template #title> 分配权限 </template>
      <div>
        <a-card>
          <a-row :gutter="24" justify="center">
            <a-col :span="8">
              <span class="text-right-gap">展开全部</span>
              <a-switch type="round" v-model="treeSwitch.expandAll" @change="onExpandAll">
                <template #checked> 是 </template>
                <template #unchecked> 否 </template>
              </a-switch>
            </a-col>
            <a-col :span="8">
              <span class="text-right-gap">全选节点</span>
              <a-switch type="round" v-model="treeSwitch.selectAll" @change="onSelectAll">
                <template #checked> 是 </template>
                <template #unchecked> 否 </template>
              </a-switch>
            </a-col>
            <a-col :span="8">
              <a-tooltip
                content="权限树的父子节点独立，因为若节点关联，父节点会存在半选情况，半选节点的ID不会返回，会导致菜单无法渲染"
              >
                <span>父子关联 <icon-question-circle-fill /></span>
              </a-tooltip>
            </a-col>
          </a-row>
        </a-card>

        <a-tree
          ref="treeRef"
          :fieldNames="{
            key: 'id',
            title: 'i18n',
            children: 'children'
          }"
          :check-strictly="true"
          :checkable="true"
          :show-line="true"
          :unmount-on-close="true"
          v-model:checked-keys="permissionKeys"
          :data="permissionTree"
        />
      </div>
    </a-drawer>
  </div>
</template>

<script setup lang="ts">
import { getRoleAPI, getMenuListAPI, getUserPermissionAPI } from "@/api/modules/system/index";
import { deepClone } from "@/utils";
import useGlobalProperties from "@/hooks/useGlobalProperties";
const proxy = useGlobalProperties();
const openState = ref(dictFilter("status"));
const form = ref({
  name: "",
  code: "",
  time: [],
  status: null
});
const search = () => {
  getRole();
};
const reset = () => {
  form.value = {
    name: "",
    code: "",
    time: [],
    status: null
  };
  getRole();
};

// 新增
const open = ref(false);
const rules = {
  name: [{ required: true, message: "请输入角色名称" }],
  code: [{ required: true, message: "请输入角色编码" }],
  status: [{ required: true, message: "请选择状态" }]
};
const addFrom = ref<any>({
  name: "",
  code: "",
  status: 1,
  sort: 1,
  description: ""
});

const title = ref("");
const formRef = ref();
const onAdd = () => {
  title.value = "新增角色";
  open.value = true;
};
const handleOk = async () => {
  let state = await formRef.value.validate();
  if (state) return (open.value = true); // 校验不通过
  arcoMessage("success", "模拟提交成功");
  getRole();
};
// 关闭对话框动画结束后触发
const afterClose = () => {
  formRef.value.resetFields();
  addFrom.value = {
    name: "",
    code: "",
    status: 1,
    sort: 1,
    description: ""
  };
};
// 修改角色
const onUpdate = (row: any) => {
  title.value = "修改角色";
  addFrom.value = deepClone(row);
  open.value = true;
};

// 获取列表
const loading = ref(false);
const pagination = ref({
  pageSize: 10,
  showPageSize: true
});
const roleList = ref([]);
const getRole = async () => {
  try {
    loading.value = true;
    let res = await getRoleAPI();
    res.data.forEach((item: any) => item.admin && (item.disabled = true));
    roleList.value = res.data;
  } finally {
    loading.value = false;
  }
};
const selectedKeys = ref([]);
const select = (list: []) => {
  selectedKeys.value = list;
};
const selectAll = (state: boolean) => {
  selectedKeys.value = state ? (roleList.value.map((el: any) => el.id) as []) : [];
};

// 获取权限树
const treeRef = ref();
const treeSwitch = ref({
  expandAll: true, // 展开全部
  selectAll: false // 全选
});
// 展开全部
const onExpandAll = (state: boolean) => {
  treeRef.value.expandAll(state);
};
// 全选
const onSelectAll = (state: boolean) => {
  treeRef.value.checkAll(state);
};

// 重置节点操作开关
const treeSwitchReset = () => {
  treeSwitch.value = {
    expandAll: true, // 固定
    selectAll: false // 全选
  };
};
const permissionTree = ref([]);
const permissionKeys = ref([]);
const getMenuList = async () => {
  let { data } = await getMenuListAPI();
  translation(data);
  permissionTree.value = data;
};

// 分配权限
const drawerOpen = ref(false);
const onPrivileges = async (row: any) => {
  let res = await getUserPermissionAPI({ role: row.code });
  permissionKeys.value = res.data;
  drawerOpen.value = true;
  treeRef.value.expandAll(true);
};

const drawerOk = () => {
  drawerOpen.value = false;
  treeSwitchReset();
  arcoMessage("success", "模拟提交成功");
  getRole();
};
const drawerCancel = () => {
  drawerOpen.value = false;
  treeSwitchReset();
};

// 语言转化
const translation = (tree: any) => {
  tree.forEach((item: any) => {
    if (item.children) translation(item.children);
    if (item.meta.title) {
      item.i18n = proxy.$t(`menu.${item.meta.title}`);
    }
  });
};

getRole();
getMenuList();
</script>

<style lang="scss" scoped>
.text-right-gap {
  margin-right: $margin;
}
</style>
