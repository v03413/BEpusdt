<template>
  <a-layout class="layout">
    <div :class="asideDark ? 'aside dark' : 'aside'" v-if="isPc">
      <Logo />
      <a-layout-sider :collapsed="collapsed" breakpoint="xl" class="layout_side" :width="220">
        <a-scrollbar style="height: 100%; overflow: auto" outer-class="scrollbar"><Menu :route-tree="routeList" /></a-scrollbar>
      </a-layout-sider>
    </div>
    <a-layout class="layout-right">
      <a-layout-header class="header">
        <div class="header-left">
          <ButtonCollapsed />
          <Breadcrumb v-if="!isPc" />
        </div>

        <div class="layout-head-menu" v-if="isPc">
          <a-menu
            v-if="drawing"
            mode="horizontal"
            :selected-keys="[selectedMenu]"
            @menu-item-click="onMenuItem"
            :popup-max-height="600"
          >
            <template v-for="item in routeTree" :key="item.path">
              <a-menu-item v-if="!item.meta.hide" :key="item.path" :popup-max-height="600">
                <template #icon v-if="item.meta.svgIcon || item.meta.icon">
                  <MenuItemIcon :svg-icon="item.meta.svgIcon" :icon="item.meta.icon" />
                </template>
                <span>{{ $t(`menu.${item.meta.title}`) }}</span>
              </a-menu-item>
            </template>
          </a-menu>
        </div>
        <HeaderRight />
      </a-layout-header>
      <Main />
      <Footer v-if="isFooter" />
    </a-layout>
  </a-layout>
</template>

<script setup lang="ts">
import Logo from "@/layout/components/Logo/index.vue";
import Main from "@/layout/components/Main/index.vue";
import Footer from "@/layout/components/Footer/index.vue";
import Menu from "@/layout/components/Menu/index.vue";
import HeaderRight from "@/layout/components/Header/components/header-right/index.vue";
import MenuItemIcon from "@/layout/components/Menu/menu-item-icon.vue";
import ButtonCollapsed from "@/layout/components/Header/components/button-collapsed/index.vue";
import Breadcrumb from "@/layout/components/Header/components/Breadcrumb/index.vue";
import { useRouteConfigStore } from "@/store/modules/route-config";
import { useRoutingMethod } from "@/hooks/useRoutingMethod";
import { storeToRefs } from "pinia";
import { useThemeConfig } from "@/store/modules/theme-config";
import { useDevicesSize } from "@/hooks/useDevicesSize";
defineOptions({ name: "LayoutMixing" });
const route = useRoute();
const router = useRouter();
const routerStore = useRouteConfigStore();
const themeStore = useThemeConfig();
const { isFooter, collapsed, asideDark, language } = storeToRefs(themeStore);
const { routeTree } = storeToRefs(routerStore);
const { isPc } = useDevicesSize();

const drawing = ref<boolean>(true);
watch(language, () => {
  drawing.value = false;
  nextTick(() => (drawing.value = true));
});

// 横向菜单点击事件
// 将一级菜单下的children给左侧菜单
// 如果没有children则直接自身菜单给左侧菜单
const routeList = ref<any>([]);
const onMenuItem = (path: string) => {
  const { findLinearArray } = useRoutingMethod();
  const find = findLinearArray(path);
  // 路由存在则存入并跳转，不存在则跳404
  if (find) {
    // 给左侧树赋值
    setAside(find);
    // 若有重定向，则跳转到重定向的路由
    // 如果有子路由则重定向到自己的第一个菜单
    // 如果没有子路由则说明当前父级是一个菜单，直接跳转
    let path = "";
    if (find.redirect) {
      path = find.redirect;
    } else if (find.children && find.children.length > 0) {
      path = find.children[0].path;
    } else {
      path = find.path;
    }
    router.push(path);
  } else {
    router.push("/404");
  }
};

// 给左侧树赋值
const setAsideMenu = (find: Menu.MenuOptions) => {
  // 将父级的chindren给左侧树
  if (find.children && find.children.length > 0) {
    routeList.value = find.children;
  } else {
    // 如果没有则直接将父级给左侧树，做一级兜底
    routeList.value = [find];
  }
};

const setAside = debounce(setAsideMenu, 150);

// 由于刷新后，路由信息丢失，所以需要重新获取
// 混合布局的横向菜单为顶层路由下的一级菜单
const selectedMenu = computed(() => {
  const { getAllParentRoute } = useRoutingMethod();
  // 动态路由参数会在path拼接，导致匹配失败
  // 这取matched做路由匹配
  const find = getAllParentRoute(route.matched.at(-1).path);
  if (find) {
    setAside(find[0]);
    return find[0].path;
  }
  return "";
});
</script>

<style lang="scss" scoped>
.layout {
  height: 100vh;
}
.dark {
  background: #232324;
}
.layout_side {
  height: calc(100% - 60px);
  .scrollbar {
    height: 100%;
  }
}

// 修改左侧滚动条宽度
:deep(.arco-scrollbar-thumb-direction-vertical .arco-scrollbar-thumb-bar) {
  width: 4px;
  margin-left: 8px;
}

// 去掉右侧阴影并替换为边线
:deep(.arco-layout-sider-light) {
  border-right: $border-1 solid $color-border-2;
  box-shadow: unset;
}

// 解决折叠菜单的icon不居中问题
:deep(.arco-menu-vertical.arco-menu-collapsed) {
  // 消除icon的自带padding值，并且让元素居中
  .arco-menu-has-icon {
    justify-content: center;
    padding: 0;
  }

  // 消除icon的自带margin-right值，并且设置icon的padding值以保留icon空隙
  .arco-menu-icon {
    padding: 10px 0;
    margin-right: 0;
  }

  // 消除title占位
  .arco-menu-title {
    display: none;
  }
}

// 去掉sider背景
.arco-layout-sider {
  background: unset;
}
.layout-right {
  display: grid;
  grid-template-rows: auto 1fr auto;
  height: 100%;
  .header {
    position: relative;
    box-sizing: border-box;
    display: flex;
    align-items: center;
    justify-content: space-between;
    height: 60px;
    padding: 0 $padding;
    overflow: hidden;
    border-bottom: $border-1 solid $color-border-2;
    .header-left {
      display: flex;
      align-items: center;
    }
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
