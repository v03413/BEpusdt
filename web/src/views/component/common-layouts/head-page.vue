<template>
  <a-space direction="vertical" fill size="medium">
    <a-card title="头部工具栏">
      <s-layout-tools>
        <template #left>
          <a-button type="primary">
            <template #icon>
              <icon-plus />
            </template>
            新增
          </a-button>
          <a-button type="primary" status="danger">
            <template #icon>
              <icon-delete />
            </template>
            删除
          </a-button>
        </template>
        <template #right>
          <a-input-search :style="{ width: '200px' }" placeholder="请输入关键字" search-button />
          <a-button type="primary" status="success">
            <template #icon>
              <icon-download />
            </template>
          </a-button>
          <a-button type="primary">
            <template #icon>
              <icon-refresh />
            </template>
          </a-button>
        </template>
      </s-layout-tools>
    </a-card>
    <a-card title="搜索栏布局">
      <s-layout-search @search="onSearch" @reset="onReset">
        <template #form>
          <a-grid-item>
            <div class="row-center-gap">
              <span class="text-nowrap">手机号</span>
              <a-input v-model="form.phone" placeholder="请输入手机号" allow-clear />
            </div>
          </a-grid-item>
          <a-grid-item>
            <div class="row-center-gap">
              <span class="text-nowrap">邮箱</span>
              <a-input v-model="form.email" placeholder="请输入邮箱" allow-clear />
            </div>
          </a-grid-item>
          <a-grid-item>
            <div class="row-center-gap">
              <span class="text-nowrap">状态</span>
              <a-input style="width: 100%" v-model="form.name" placeholder="请输入状态" allow-clear />
            </div>
          </a-grid-item>
          <a-grid-item>
            <div class="row-center-gap">
              <span class="text-nowrap">类型</span>
              <a-input style="width: 100%" v-model="form.name" placeholder="请输入类型" allow-clear />
            </div>
          </a-grid-item>
          <a-grid-item>
            <div class="row-center-gap">
              <span class="text-nowrap">机构</span>
              <a-input style="width: 100%" v-model="form.name" placeholder="请输入机构" allow-clear />
            </div>
          </a-grid-item>
        </template>
      </s-layout-search>
    </a-card>
  </a-space>
</template>

<script setup lang="ts">
const form = ref<any>({
  name: "",
  age: null,
  select: null,
  adderss: "",
  hobby: null,
  date: "",
  time: "",
  gender: null,
  score: null,
  rate: null,
  state: null,
  multiple: null,
  treeSelect: null
});
const onSearch = () => {
  arcoMessage("success", "搜索");
};
const onReset = () => {
  arcoMessage("success", "重置");
};
</script>

<style lang="scss" scoped>
.tools-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.row-center-gap {
  display: flex;
  column-gap: 12px;
  align-items: center;
  .text-nowrap {
    white-space: nowrap;
  }
}
</style>
