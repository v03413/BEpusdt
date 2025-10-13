<template>
  <a-grid :colGap="12" :rowGap="12" :cols="{ xs: 1, sm: 1, md: 1, lg: 2, xl: 3, xxl: 4 }" :collapsed="collapsed">
    <slot name="form"></slot>
    <a-grid-item suffix>
      <a-space>
        <a-button type="primary" @click="onSearch">
          <template #icon>
            <icon-search />
          </template>
          查询
        </a-button>
        <a-button @click="onReset">
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
const emit = defineEmits(["search", "reset"]);
const collapsed = ref<boolean>(true);

const onSearch = () => emit("search");

const onReset = () => emit("reset");
</script>

<style lang="scss" scoped></style>
