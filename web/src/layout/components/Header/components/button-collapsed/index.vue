<template>
  <div class="menu_fold">
    <a-button size="mini" type="text" class="menu_fold_icon" id="system-collapsed" @click="onCollapsed">
      <template #icon>
        <icon-menu-fold :size="18" v-if="!collapsed" />
        <icon-menu-unfold :size="18" v-if="collapsed" />
      </template>
    </a-button>
  </div>

  <DrawerAside v-if="!isPc" />
</template>

<script setup lang="ts">
import DrawerAside from "@/layout/components/drawer-aside/index.vue";
import { storeToRefs } from "pinia";
import { useThemeConfig } from "@/store/modules/theme-config";
import { useDevicesSize } from "@/hooks/useDevicesSize";
const themeStore = useThemeConfig();
const { collapsed } = storeToRefs(themeStore);
const { isPc } = useDevicesSize();
// 刷新时，PC窗口展开菜单，移动端收起菜单
collapsed.value = isPc.value ? false : true;

// 折叠
const onCollapsed = () => {
  themeStore.setCollapsed(!collapsed.value);
};

// 监听屏幕尺寸变化，PC端自动展开菜单，移动端自动收起
watch(isPc, (newV: boolean) => {
  collapsed.value = newV ? false : true;
});
</script>

<style lang="scss" scoped>
.menu_fold {
  display: flex;
  align-items: center;
  justify-content: space-around;
  width: 32px;
  height: 32px;
  .menu_fold_icon {
    color: $color-text-1;
    border-radius: $radius-box-1;
  }
}
</style>
