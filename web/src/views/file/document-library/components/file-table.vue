<template>
  <div class="container-page">
    <a-row :gutter="16" justify="space-between" class="bottom-gap">
      <a-col :xs="24" :sm="24" :md="12" :lg="12" :xl="6" :xxl="6" style="margin-bottom: 16px">
        <a-space>
          <a-button type="primary">
            <template #icon>
              <icon-plus />
            </template>
            <template #default>新增</template>
          </a-button>
          <a-button type="primary" status="danger">
            <template #icon>
              <icon-delete />
            </template>
            <template #default>删除</template>
          </a-button>
          <a-button>
            <template #icon>
              <icon-download />
            </template>
            <template #default>导出</template>
          </a-button>
        </a-space>
      </a-col>

      <a-col :xs="24" :sm="24" :md="12" :lg="12" :xl="8" :xxl="6">
        <a-space>
          <a-input v-model="form.name" placeholder="请输入搜索关键词" />
          <a-button type="primary" status="success" @click="getList">
            <template #icon>
              <icon-search />
            </template>
            <template #default>搜索</template>
          </a-button>
        </a-space>
      </a-col>
    </a-row>

    <a-table
      row-key="key"
      size="small"
      :bordered="{
        cell: true
      }"
      :scroll="{ x: '100%', y: '100%', minWidth: 1000 }"
      :loading="loading"
      :data="list"
      :row-selection="rowSelection"
      v-model:selectedKeys="selectedKeys"
      :pagination="pagination"
      @page-change="pageChange"
      @page-size-change="pageSizeChange"
    >
      <template #columns>
        <a-table-column title="序号" :width="68">
          <template #cell="cell">{{ cell.rowIndex + 1 }}</template>
        </a-table-column>
        <a-table-column title="渠道来源" data-index="source">
          <template #cell="{ record }">
            <a-space>
              <s-svg-icon :name="record.sourceSvg" />
              <span>{{ record.source }}</span>
            </a-space>
          </template>
        </a-table-column>
        <a-table-column title="源类型" data-index="sourceType"></a-table-column>
        <a-table-column title="组流量共享" data-index="flow">
          <template #cell="{ record }">
            <a-badge status="processing" :text="record.flow" v-if="record.percent >= 60" />
            <a-badge status="warning" :text="record.flow" v-else-if="record.percent < 60 && record.percent >= 50" />
            <a-badge status="danger" :text="record.flow" v-else />
          </template>
        </a-table-column>
        <a-table-column title="数据比例" data-index="percent" :width="200">
          <template #cell="{ record }">
            <a-progress v-if="record.percent >= 60" :percent="record.percent / 100" />
            <a-progress
              v-else-if="record.percent < 60 && record.percent >= 50"
              status="warning"
              :percent="record.percent / 100"
            />
            <a-progress v-else status="danger" :percent="record.percent / 100" />
          </template>
        </a-table-column>
        <a-table-column title="状态" data-index="status" align="center">
          <template #cell="{ record }">
            <a-space>
              <a-tag size="small" color="green" v-if="record.status == 1">在线</a-tag>
              <a-tag size="small" color="red" v-else>离线</a-tag>
            </a-space>
          </template>
        </a-table-column>
        <a-table-column title="操作" :width="200" fixed="right">
          <template #cell="cell">
            <a-space>
              <a-button size="mini" type="primary">详情</a-button>
              <a-button size="mini">修改</a-button>
              <a-popconfirm content="确定删除这条数据吗?" type="warning">
                <a-button size="mini" type="primary" status="danger" @click="onDelete(cell)">删除</a-button>
              </a-popconfirm>
            </a-space>
          </template>
        </a-table-column>
      </template>
    </a-table>
  </div>
</template>

<script setup lang="ts">
import { getDocumentLibraryTableAPI } from "@/api/modules/file/index";
import { SelectedKeys, Form, Pagination, List } from "./config";

const form = ref<Form>({ name: "" });
const selectedKeys = ref<string[]>([]);
const rowSelection = reactive<SelectedKeys>({
  type: "checkbox",
  showCheckedAll: true,
  onlyCurrent: false
});
const pagination = ref<Pagination>({ showPageSize: true, showTotal: true, current: 1, pageSize: 10, total: 0 });
const pageChange = (page: number) => {
  pagination.value.current = page;
};
const pageSizeChange = (pageSize: number) => {
  pagination.value.pageSize = pageSize;
};
const onDelete = (cell: any) => {
  console.log("删除", cell);
};

const loading = ref<boolean>(false);
const list = ref<List[]>([]);
const getList = async () => {
  try {
    loading.value = true;
    let { data } = await getDocumentLibraryTableAPI();
    pagination.value.total = data.total;
    list.value = shuffleArray(data.list);
    loading.value = false;
  } catch {
    loading.value = false;
  }
};
getList();

const shuffleArray = (arr: List[]) => {
  // 循环遍历数组，从最后一个元素开始向前,
  // 递减生成，确保数组的每个值都被照顾到
  for (let i = arr.length - 1; i > 0; i--) {
    // 生成一个随机索引 j，范围是从 0 到 i
    // 生成的最大值就是arr数组的最大索引
    const j = Math.floor(Math.random() * (i + 1));
    // 交换 arr[i] 和 arr[j] 的值 (当前遍历的值和新生成的索引值)
    [arr[i], arr[j]] = [arr[j], arr[i]];
    // 随机修改值
    arr[i].flow = String(Math.floor(Math.random() * 5000));
    arr[i].percent = String(Math.floor(Math.random() * 100));
  }
  return arr;
};

defineExpose({
  getList
});
</script>

<style lang="scss" scoped>
.container-page {
  box-sizing: border-box;
  padding: $padding;
}
.bottom-gap {
  margin-bottom: $padding;
}
:deep(.arco-form-item) {
  margin-bottom: 0;
}
</style>
