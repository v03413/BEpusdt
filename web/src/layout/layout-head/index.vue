<template>
  <div class="layout-head">
    <div class="layout-head-top">
      <a-layout-header class="header">
        <div class="header-logo" v-if="!isMobile">
          <Logo />
        </div>
        <div class="layout-head-menu" v-if="!isMobile">
          <a-menu
            v-if="drawing"
            mode="horizontal"
            :selected-keys="[selectedMenu]"
            @menu-item-click="onMenuItem"
            :popup-max-height="600"
          >
            <template v-for="item in routeTree" :key="item.path">
              <a-sub-menu v-if="menuShow(item)" :key="item.path" :popup-max-height="600">
                <template #icon v-if="item.meta.svgIcon || item.meta.icon">
                  <MenuItemIcon :svg-icon="item.meta.svgIcon" :icon="item.meta.icon" />
                </template>
                <template #title>{{ $t(`menu.${item.meta.title}`) }}</template>
                <MenuItem :route-tree="item.children" />
              </a-sub-menu>
              <a-menu-item v-else-if="aMenuShow(item)" :key="item?.path">
                <template #icon v-if="item.meta.svgIcon || item.meta.icon">
                  <MenuItemIcon :svg-icon="item.meta.svgIcon" :icon="item.meta.icon" />
                </template>
                <span>{{ $t(`menu.${item.meta.title}`) }}</span>
              </a-menu-item>
            </template>
          </a-menu>
        </div>
        <ButtonCollapsed v-else />

        <HeaderRight />
      </a-layout-header>
      <Main />
      <Footer v-if="isFooter" />
    </div>
  </div>
</template>

<script setup lang="ts">
import Logo from "@/layout/components/Logo/index.vue";
import HeaderRight from "@/layout/components/Header/components/header-right/index.vue";
import Main from "@/layout/components/Main/index.vue";
import Footer from "@/layout/components/Footer/index.vue";
import MenuItem from "@/layout/components/Menu/menu-item.vue";
import MenuItemIcon from "@/layout/components/Menu/menu-item-icon.vue";
import ButtonCollapsed from "@/layout/components/Header/components/button-collapsed/index.vue";
import { useRoutingMethod } from "@/hooks/useRoutingMethod";
import { storeToRefs } from "pinia";
import { useRouteConfigStore } from "@/store/modules/route-config";
import { useThemeConfig } from "@/store/modules/theme-config";
import { useMenuMethod } from "@/hooks/useMenuMethod";
import { useDevicesSize } from "@/hooks/useDevicesSize";
defineOptions({ name: "LayoutHead" });
const route = useRoute();
const router = useRouter();
const routerStore = useRouteConfigStore();
const themeStore = useThemeConfig();
const { routeTree } = storeToRefs(routerStore);
const { isFooter, language } = storeToRefs(themeStore);
const { isMobile } = useDevicesSize();
const { menuShow, aMenuShow } = useMenuMethod();

const drawing = ref<boolean>(true);
watch(language, () => {
  drawing.value = false;
  nextTick(() => (drawing.value = true));
});

const onMenuItem = (path: string) => router.push(path);

const selectedMenu = computed(() => {
  const { getAllParentRoute } = useRoutingMethod();
  const find = getAllParentRoute(route.matched.at(-1).path);
  if (!find) return "";
  let path = "";
  for (let i = find.length - 1; i >= 0; i--) {
    if (find[i].meta && !find[i].meta.hide) {
      path = find[i].path;
      break;
    }
  }
  return path;
});
</script>

<style lang="scss" scoped>
.layout-head {
  height: 100vh;
  &-top {
    position: relative;
    display: grid;
    grid-template-rows: auto 1fr auto;
    height: 100%;
  }
}
.header {
  position: relative;
  box-sizing: border-box;
  display: flex;
  align-items: center;
  width: 100%;
  height: 60px;
  padding: 0 $padding;
  overflow: hidden;
  border-bottom: $border-1 solid $color-border-2;
  .header-logo {
    max-width: 180px;
  }
  .layout-head-menu {
    display: flex;
    flex: 1;
    overflow: hidden;
  }
}
:deep(.arco-menu-pop) {
  white-space: nowrap;
}

// 横向菜单样式修改
:deep(.arco-menu-horizontal) {
  flex: 1;
  overflow: hidden;
  .arco-menu-inner {
    padding-left: 0; // 横向排列，禁用左padding
    .arco-menu-overflow-wrap {
      white-space: nowrap; // 禁用换行，否则会导致菜单换行闪烁
    }
  }
}
</style>
