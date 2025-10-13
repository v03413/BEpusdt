<template>
  <a-grid :colGap="12" :rowGap="12" :cols="{ xs: 1, sm: 1, md: 1, lg: 2, xl: 3, xxl: 4 }" :collapsed="collapsed">
    <a-grid-item>
      <div class="row-center-gap">
        <span class="text-nowrap">姓名</span>
        <a-input style="width: 100%" v-model="form.name" placeholder="请输入姓名" allow-clear />
      </div>
    </a-grid-item>
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
    <a-grid-item suffix>
      <a-space>
        <a-button type="primary">
          <template #icon>
            <icon-search />
          </template>
          查询
        </a-button>
        <a-button>
          <template #icon>
            <icon-refresh />
          </template>
          重置
        </a-button>
        <a-button type="text" @click="collapsed = !collapsed">
          <template #icon>
            <icon-down v-if="collapsed" />
            <icon-up v-else />
          </template>
          <span>{{ collapsed ? "展开" : "收起" }}</span>
        </a-button>
      </a-space>
    </a-grid-item>
  </a-grid>
</template>

<script setup lang="ts">
const form = ref<any>({
  name: "",
  phone: "",
  email: "",
  address: "",
  status: null
});

const collapsed = ref<boolean>(false);
</script>

<style lang="scss" scoped>
.row-center-gap {
  display: flex;
  column-gap: 12px;
  align-items: center;
  .text-nowrap {
    white-space: nowrap;
  }
}
</style>
