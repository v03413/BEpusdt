<template>
  <div class="snow-page">
    <div class="snow-inner">
      <s-layout-tools>
        <template #left>
          <a-space wrap>
            <a-input v-model="form.name" placeholder="请输入字典名称" allow-clear />
            <a-input v-model="form.code" placeholder="请输入字典编码" allow-clear />
            <a-select placeholder="启用状态" v-model="form.status" style="width: 120px" allow-clear>
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
        row-key="id"
        :data="dictList"
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
          <a-table-column title="字典名称" data-index="name"></a-table-column>
          <a-table-column title="字典编码" data-index="code"></a-table-column>
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
                <a-button type="primary" status="success" size="mini" @click="onDictData(record)">
                  <template #icon><icon-file /></template>
                  <span>字典值</span>
                </a-button>
                <a-button type="primary" size="mini" @click="onUpdate(record)">
                  <template #icon><icon-edit /></template>
                  <span>修改</span>
                </a-button>
                <a-popconfirm type="warning" content="确定删除该字典吗?">
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

    <a-modal v-model:visible="open" @close="afterClose" @ok="handleOk" @cancel="afterClose">
      <template #title> {{ title }} </template>
      <div>
        <a-form ref="formRef" auto-label-width :rules="rules" :model="addFrom">
          <a-form-item field="name" label="字典名称" validate-trigger="blur">
            <a-input v-model="addFrom.name" placeholder="请输入字典名称" allow-clear />
          </a-form-item>
          <a-form-item field="code" label="字典编码" validate-trigger="blur">
            <a-input v-model="addFrom.code" placeholder="请输入字典编码" allow-clear />
          </a-form-item>
          <a-form-item field="description" label="描述" validate-trigger="blur">
            <a-textarea v-model="addFrom.description" placeholder="请输入字典描述" allow-clear />
          </a-form-item>
          <a-form-item field="description" label="状态" validate-trigger="blur">
            <a-switch type="round" :checked-value="1" :unchecked-value="0" v-model="addFrom.status">
              <template #checked> 启用 </template>
              <template #unchecked> 禁用 </template>
            </a-switch>
          </a-form-item>
        </a-form>
      </div>
    </a-modal>

    <a-modal width="50%" v-model:visible="detailOpen" @ok="detailOk" ok-text="关闭" :hide-cancel="true">
      <template #title> 字典详情 </template>
      <div>
        <a-row>
          <a-space wrap>
            <a-button type="primary" @click="onAddDetail">
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
          :data="dictDetail.list"
          :bordered="{ cell: true }"
          :loading="detailLoading"
          :scroll="{ x: '100%', y: '100%' }"
          :pagination="pagination"
          :row-selection="{ type: 'checkbox', showCheckedAll: true }"
          :selected-keys="selectedKeysDetail"
          @select="selectDetail"
          @select-all="selectAllDetail"
        >
          <template #columns>
            <a-table-column title="序号" :width="64">
              <template #cell="cell">{{ cell.rowIndex + 1 }}</template>
            </a-table-column>
            <a-table-column title="字典名" data-index="name"></a-table-column>
            <a-table-column title="字典值" data-index="value"></a-table-column>
            <a-table-column title="状态" :width="100" align="center">
              <template #cell="{ record }">
                <a-tag bordered size="small" color="arcoblue" v-if="record.status === 1">启用</a-tag>
                <a-tag bordered size="small" color="red" v-else>禁用</a-tag>
              </template>
            </a-table-column>
            <a-table-column title="操作" align="center" :fixed="'right'">
              <template #cell="{ record }">
                <a-space>
                  <a-button type="primary" size="mini" @click="onDetailUpdate(record)">
                    <template #icon><icon-edit /></template>
                    <span>修改</span>
                  </a-button>
                  <a-popconfirm type="warning" content="确定删除该字典吗?">
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
    </a-modal>

    <a-modal v-model:visible="detailCaseOpen" @close="afterCloseDetail" @ok="handleOkDetail" @cancel="afterCloseDetail">
      <template #title> {{ detailTitle }} </template>
      <div>
        <a-form ref="detailFormRef" auto-label-width :rules="detaulRules" :model="deatilForm">
          <a-form-item field="name" label="字典名称" validate-trigger="blur">
            <a-input v-model="deatilForm.name" placeholder="请输入字典名称" allow-clear />
          </a-form-item>
          <a-form-item field="value" label="字典值" validate-trigger="blur">
            <a-input v-model="deatilForm.value" placeholder="请输入字典值" allow-clear />
          </a-form-item>
          <a-form-item field="description" label="状态" validate-trigger="blur">
            <a-switch type="round" :checked-value="1" :unchecked-value="0" v-model="deatilForm.status">
              <template #checked> 启用 </template>
              <template #unchecked> 禁用 </template>
            </a-switch>
          </a-form-item>
        </a-form>
      </div>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { getDictAPI } from "@/api/modules/system/index";
import { deepClone } from "@/utils";
const openState = ref(dictFilter("status"));
const form = ref({
  name: "",
  code: "",
  status: null
});
const search = () => {
  getDict();
};
const reset = () => {
  form.value = {
    name: "",
    code: "",
    status: null
  };
  getDict();
};

const open = ref<boolean>(false);
const title = ref<string>("");
const rules = {
  name: [
    {
      required: true,
      message: "请输入字典名称"
    }
  ],
  code: [
    {
      required: true,
      message: "请输入字典编码"
    }
  ],
  status: [
    {
      required: true,
      message: "请选择状态"
    }
  ]
};
const addFrom = ref({
  name: "",
  code: "",
  description: "",
  status: 1
});
const formRef = ref();
const onAdd = () => {
  open.value = true;
  title.value = "新增字典";
};
const handleOk = async () => {
  let state = await formRef.value.validate();
  if (state) return (open.value = true); // 校验不通过
  arcoMessage("success", "模拟提交成功");
  getDict();
};
// 关闭对话框动画结束后触发
const afterClose = () => {
  formRef.value.resetFields();
  addFrom.value = {
    name: "",
    code: "",
    description: "",
    status: 1
  };
};
const onUpdate = (record: any) => {
  title.value = "修改字典";
  addFrom.value = deepClone(record);
  open.value = true;
};

const loading = ref(false);
const pagination = ref({
  pageSize: 10,
  showPageSize: true
});
const selectedKeys = ref([]);
const select = (list: []) => {
  selectedKeys.value = list;
};
const selectAll = (state: boolean) => {
  selectedKeys.value = state ? (dictList.value.map((el: any) => el.id) as []) : [];
};
const dictList = ref();
const getDict = async () => {
  loading.value = true;
  let res = await getDictAPI();
  dictList.value = res.data || [];
  loading.value = false;
};

// 字典详情
const detailLoading = ref(false);
const detailOpen = ref(false);
const dictDetail = ref({
  list: []
});
const onDictData = (record: any) => {
  detailLoading.value = true;
  dictDetail.value = record;
  detailOpen.value = true;
  detailLoading.value = false;
};
const detailOk = () => {
  detailOpen.value = false;
};
const deatilForm = ref({
  name: "",
  value: "",
  status: 1
});
const detaulRules = ref({
  name: [
    {
      required: true,
      message: "请输入字典名称"
    }
  ],
  value: [
    {
      required: true,
      message: "请输入字典值"
    }
  ],
  status: [
    {
      required: true,
      message: "请选择状态"
    }
  ]
});
const detailFormRef = ref();
const detailTitle = ref("");
const detailCaseOpen = ref(false);
const onAddDetail = () => {
  detailTitle.value = "新增字典数据";
  detailCaseOpen.value = true;
};
const handleOkDetail = async () => {
  let state = await detailFormRef.value.validate();
  if (state) return (detailCaseOpen.value = true); // 校验不通过
  arcoMessage("success", "模拟提交成功");
};
const onDetailUpdate = (record: any) => {
  detailTitle.value = "修改字典数据";
  deatilForm.value = deepClone(record);
  detailCaseOpen.value = true;
};
// 关闭对话框动画结束后触发
const afterCloseDetail = () => {
  detailFormRef.value.resetFields();
  deatilForm.value = {
    name: "",
    value: "",
    status: 1
  };
};
const selectedKeysDetail = ref([]);
const selectDetail = (list: []) => {
  selectedKeysDetail.value = list;
};
const selectAllDetail = (state: boolean) => {
  selectedKeysDetail.value = state ? (dictDetail.value.list.map((el: any) => el.id) as []) : [];
};

getDict();
</script>

<style lang="scss" scoped></style>
