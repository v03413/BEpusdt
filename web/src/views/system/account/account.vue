<template>
  <div class="snow-fill">
    <div class="snow-fill-inner container">
      <div class="left-box">
        <a-input placeholder="请输入部门名称">
          <template #prefix>
            <icon-search />
          </template>
        </a-input>
        <div class="tree-box">
          <a-tree ref="treeRef" :field-names="fieldNames" :data="treeData" show-line @select="onSelectTree"> </a-tree>
        </div>
      </div>
      <div class="right-box">
        <a-space wrap>
          <a-input v-model="form.name" placeholder="请输入用户名称" allow-clear />
          <a-input v-model="form.phone" placeholder="请输入手机号码" allow-clear />
          <a-select placeholder="用户状态" v-model="form.status" style="width: 120px" allow-clear>
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
            <a-button type="primary" status="danger">
              <template #icon><icon-delete /></template>
              <span>删除</span>
            </a-button>
          </a-space>
        </a-row>

        <a-table
          row-key="id"
          :data="accountList"
          :bordered="{ cell: true }"
          :loading="loading"
          :scroll="{ x: '120%', y: '100%' }"
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
            <a-table-column title="用户名称" data-index="userName"></a-table-column>
            <a-table-column title="昵称" data-index="nickName"></a-table-column>
            <a-table-column title="性别" data-index="sex" align="center">
              <template #cell="{ record }">
                <a-tag bordered size="small" color="arcoblue" v-if="record.sex == 1">男</a-tag>
                <a-tag bordered size="small" color="red" v-else-if="record.sex == 0">女</a-tag>
                <a-tag bordered size="small" v-else>未知</a-tag>
              </template>
            </a-table-column>
            <a-table-column title="部门" data-index="deptName"></a-table-column>
            <a-table-column title="手机号" data-index="phone" :width="180"></a-table-column>
            <a-table-column title="状态" :width="100" align="center">
              <template #cell="{ record }">
                <a-tag bordered size="small" color="arcoblue" v-if="record.status === 1">启用</a-tag>
                <a-tag bordered size="small" color="red" v-else>禁用</a-tag>
              </template>
            </a-table-column>
            <a-table-column title="描述" data-index="description" :ellipsis="true" :tooltip="true"></a-table-column>
            <a-table-column title="创建时间" data-index="createTime" :width="180"></a-table-column>
            <a-table-column title="操作" :width="200" align="center" :fixed="'right'">
              <template #cell="{ record }">
                <a-space>
                  <a-button type="primary" size="mini" @click="onUpdate(record)">
                    <template #icon><icon-edit /></template>
                    <span>修改</span>
                  </a-button>
                  <a-popconfirm type="warning" content="确定删除该账号吗?">
                    <a-button type="primary" status="danger" size="mini" :disabled="record.admin">
                      <template #icon><icon-delete /></template>
                      <span>删除</span>
                    </a-button>
                  </a-popconfirm>
                  <a-tooltip content="用户详情">
                    <a-button type="primary" status="success" size="mini" @click="onDetail(record)">
                      <template #icon>
                        <icon-more />
                      </template>
                    </a-button>
                  </a-tooltip>
                </a-space>
              </template>
            </a-table-column>
          </template>
        </a-table>
      </div>
    </div>

    <a-modal width="40%" v-model:visible="open" @close="afterClose" @ok="handleOk" @cancel="afterClose">
      <template #title> {{ title }} </template>
      <div>
        <a-form ref="formRef" auto-label-width :rules="rules" :model="addFrom">
          <a-row>
            <a-col :span="12">
              <a-form-item field="userName" label="用户名称" validate-trigger="blur">
                <a-input v-model="addFrom.userName" placeholder="请输入用户名称" allow-clear />
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item field="nickName" label="昵称" validate-trigger="blur">
                <a-input v-model="addFrom.nickName" placeholder="请输入昵称" allow-clear />
              </a-form-item>
            </a-col>
          </a-row>
          <a-row>
            <a-col :span="12">
              <a-form-item field="phone" label="手机号码" validate-trigger="blur">
                <a-input v-model="addFrom.phone" placeholder="请输入手机号码" allow-clear />
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item field="email" label="邮箱" validate-trigger="blur">
                <a-input v-model="addFrom.email" placeholder="请输入邮箱" allow-clear />
              </a-form-item>
            </a-col>
          </a-row>
          <a-row>
            <a-col :span="24">
              <a-form-item field="sex" label="性别" validate-trigger="blur">
                <a-radio-group v-model="addFrom.sex" :options="sexOption">
                  <template #label="{ data }">
                    <div>{{ data.name }}</div>
                  </template>
                </a-radio-group>
              </a-form-item>
            </a-col>
          </a-row>
          <a-form-item field="deptId" label="所属部门" validate-trigger="blur">
            <a-tree-select
              v-model="addFrom.deptId"
              :data="treeData"
              :field-names="{
                key: 'id',
                title: 'name',
                children: 'children'
              }"
              placeholder="请选择所属部门"
            ></a-tree-select>
          </a-form-item>
          <a-form-item field="roles" label="角色" validate-trigger="blur">
            <a-select v-model="addFrom.roles" multiple placeholder="请选择角色">
              <a-option
                v-for="item in roleList"
                :key="item.code"
                :value="item.code"
                :label="item.name"
                :disabled="item.admin"
              ></a-option>
            </a-select>
          </a-form-item>
          <a-form-item field="status" label="状态" validate-trigger="blur">
            <a-switch type="round" :checked-value="1" :unchecked-value="0" v-model="addFrom.status">
              <template #checked> 启用 </template>
              <template #unchecked> 禁用 </template>
            </a-switch>
          </a-form-item>
          <a-form-item field="description" label="描述" validate-trigger="blur">
            <a-textarea v-model="addFrom.description" placeholder="请输入描述" allow-clear />
          </a-form-item>
        </a-form>
      </div>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { getDivisionAPI, getAccountAPI, getRoleAPI } from "@/api/modules/system/index";
import { deepClone } from "@/utils";

const router = useRouter();
const openState = ref(dictFilter("status"));
const sexOption = ref(dictFilter("gender"));
const form = ref({
  name: "",
  phone: "",
  status: null,
  createTime: []
});
const search = () => {
  getAccount();
};
const reset = () => {
  form.value = {
    name: "",
    phone: "",
    status: null,
    createTime: []
  };
  getAccount();
};
// 新增
const open = ref(false);
const rules = {
  userName: [
    {
      required: true,
      message: "请输入用户名称"
    }
  ],
  nickName: [
    {
      required: true,
      message: "请输入昵称"
    }
  ],
  sex: [
    {
      required: true,
      message: "请选择性别"
    }
  ],
  deptId: [
    {
      required: true,
      message: "请选择所属部门"
    }
  ],
  roles: [
    {
      required: true,
      message: "请选择角色"
    }
  ],
  status: [
    {
      required: true,
      message: "请选择状态"
    }
  ]
};
const addFrom = ref<any>({
  userName: "",
  nickName: "",
  phone: "",
  email: "",
  sex: 2,
  deptId: null,
  roles: [],
  status: 1,
  description: ""
});
const formType = ref(0); // 0新增 1修改
const title = ref("");
const formRef = ref();
const onAdd = () => {
  title.value = "新增账号";
  formType.value = 0;
  open.value = true;
};
const handleOk = async () => {
  let state = await formRef.value.validate();
  if (state) return (open.value = true); // 校验不通过
  arcoMessage("success", "模拟提交成功");
  getAccount();
};
// 关闭对话框动画结束后触发
const afterClose = () => {
  formRef.value.resetFields();
  addFrom.value = {
    userName: "",
    nickName: "",
    phone: "",
    email: "",
    sex: 2,
    deptId: null,
    roles: [],
    status: 1,
    description: ""
  };
};
const onUpdate = (row: any) => {
  title.value = "修改账号";
  formType.value = 1;
  addFrom.value = deepClone(row);
  open.value = true;
};

const onDetail = (row: any) => {
  console.log(row);

  router.push({
    path: "/system/userinfo",
    query: {
      id: row.id,
      userName: row.userName
    }
  });
};

const loading = ref(false);
const pagination = ref({
  pageSize: 10,
  showPageSize: true
});

// 账户
const accountList = ref();
const getAccount = async () => {
  loading.value = true;
  let res = await getAccountAPI();
  res.data.forEach((item: any) => item.admin && (item.disabled = true));
  accountList.value = res.data;
  loading.value = false;
};
const selectedKeys = ref([]);
const select = (list: []) => {
  selectedKeys.value = list;
};
const selectAll = (state: boolean) => {
  selectedKeys.value = state ? (accountList.value.map((el: any) => el.id) as []) : [];
};

// 部门树
const fieldNames = ref({
  key: "id",
  title: "name",
  children: "children"
});
const treeData = ref();
const treeRef = ref();
const getDivision = async () => {
  let res = await getDivisionAPI();
  treeData.value = res.data;
  setTimeout(() => {
    treeRef.value.expandAll();
  }, 0);
};
const onSelectTree = () => {
  getAccount();
};

// 角色列表
const roleList = ref<any>([]);
const getRole = async () => {
  let res = await getRoleAPI();
  roleList.value = res.data;
};

onMounted(() => {
  getDivision();
  getAccount();
  getRole();
});
</script>

<style lang="scss" scoped>
.container {
  display: flex;
  column-gap: $padding;
  .left-box {
    display: flex;
    flex-direction: column;
    width: 250px;
    height: 100%;
    .tree-box {
      flex: 1;
      margin-top: $padding;
      overflow: auto;
    }
  }
  .right-box {
    flex: 1;
  }
}
</style>
