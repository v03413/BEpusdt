<template>
  <div class="tabs">
    <a-tabs
      :editable="true"
      :hide-content="true"
      :active-key="currentRoute.path"
      size="medium"
      type="line"
      @tab-click="onTabs"
      @delete="onDelete"
    >
      <a-tab-pane v-for="item of tabsList" :key="item.path" :title="$t(`menu.${item.meta.title}`)" :closable="!item.meta.affix" />
    </a-tabs>
    <div class="tabs_setting">
      <a-space>
        <a-tooltip :content="$t(`system.refresh`)" position="bottom" mini>
          <span ref="refreshRef" id="system-tabs-refresh" :class="rotateOpen && 'refresh'">
            <icon-refresh :size="18" @click="refresh" />
          </span>
        </a-tooltip>
        <a-dropdown trigger="hover" :popup-max-height="false">
          <div class="setting" id="system-tabs-setting"><icon-apps :size="18" /></div>
          <template #content>
            <a-doption @click="closeCurrent">
              <template #icon><icon-close /></template>
              <template #default>{{ $t(`system.close-current`) }}</template>
            </a-doption>
            <a-doption @click="closeSides('left')">
              <template #icon><icon-left /></template>
              <template #default>{{ $t(`system.close-left-side`) }}</template>
            </a-doption>
            <a-doption @click="closeSides('right')">
              <template #icon><icon-right /></template>
              <template #default>{{ $t(`system.close-right-side`) }}</template>
            </a-doption>
            <a-doption @click="closeOther('other')">
              <template #icon><icon-close-circle /></template>
              <template #default>{{ $t(`system.close-other`) }}</template>
            </a-doption>
            <a-doption @click="closeOther('all')">
              <template #icon><icon-folder-delete /></template>
              <template #default>{{ $t(`system.close-all`) }}</template>
            </a-doption>
          </template>
        </a-dropdown>
      </a-space>
    </div>
  </div>
</template>

<script setup lang="ts">
import { storeToRefs } from "pinia";
import { useRouteConfigStore } from "@/store/modules/route-config";
import { useThemeConfig } from "@/store/modules/theme-config";
const router = useRouter();
const routerStore = useRouteConfigStore();
const { tabsList, currentRoute } = storeToRefs(routerStore);

// 点击标签页，如果标签页存在，则跳转
const onTabs = (key: string) => {
  router.push(key);
};

// 删除当前标签页并跳转到最后一个标签页
const onDelete = (path: string) => {
  routerStore.removeTabsList(path);
  routerStore.removeRouteName(path);
  if (tabsList.value.length == 0) return;
  if (currentRoute.value.path != path) return;
  router.push(tabsList.value.at(-1).path);
};

// 刷新当前页
const rotateOpen = ref(false);
const refresh = () => {
  rotateOpen.value = true;
  setTimeout(() => {
    rotateOpen.value = false;
  }, 500);
  const themeStore = useThemeConfig();
  themeStore.setRefreshPage(false);
  currentRoute.value.meta.keepAlive && routerStore.removeRouteName(currentRoute.value.path);
  nextTick(() => {
    themeStore.setRefreshPage(true);
    currentRoute.value.meta.keepAlive && routerStore.setRoutePaths(currentRoute.value.path);
  });
};

// 关闭当前
const closeCurrent = () => {
  onDelete(currentRoute.value.path);
};

// 关闭右侧&关闭左侧
const closeSides = (type: string) => {
  // 获得当前index
  let currentIndex = tabsList.value.findIndex((item: Menu.MenuOptions) => item.path === currentRoute.value.path);
  // 过滤出两侧可关闭的 affix: false 表示可关闭
  let rightList = tabsList.value.filter((item: Menu.MenuOptions, index: number) => {
    if (type == "right") {
      if (index > currentIndex && !item.meta.affix) return item;
    } else {
      if (index < currentIndex && !item.meta.affix) return item;
    }
  });
  // 返回可关闭名称
  let rightPaths = rightList.map((item: Menu.MenuOptions) => item.path);
  // 删除右侧
  tabsList.value = tabsList.value.filter((item: Menu.MenuOptions) => !rightPaths.includes(item.path));
  // 删除缓存
  routerStore.removeRoutePaths(rightPaths);
};

// 关闭其它&关闭全部
const closeOther = (type: string) => {
  // 过滤出可关闭项 affix: false 表示可关闭
  let list = tabsList.value.filter((item: Menu.MenuOptions) => {
    if (type == "other") {
      if (item.path != currentRoute.value.path && !item.meta.affix) {
        return item;
      }
    } else {
      if (!item.meta.affix) {
        return item;
      }
    }
  });
  // 返回可关闭名称
  let rightNames = list.map((item: Menu.MenuOptions) => item.path);
  // 删除可关闭项
  tabsList.value = tabsList.value.filter((item: Menu.MenuOptions) => !rightNames.includes(item.path));
  // 删除缓存
  routerStore.removeRoutePaths(rightNames);
  // 关闭全部，若当前被关闭则跳转最后一个
  if (tabsList.value.length != 0 && !currentRoute.value.meta.affix && type == "all") {
    router.push(tabsList.value.at(-1).path);
  }
};
</script>

<style lang="scss" scoped>
.tabs {
  box-sizing: border-box;
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 40px;
  border-bottom: $border-1 solid $color-border-2;
  .tabs_setting {
    margin: 0 0 0 $margin;
    .setting {
      margin-right: $margin;
      color: $color-text-2;
    }
    .refresh {
      transform: rotate(360deg);
      transition: transform 0.5s;
    }
  }
}
:deep(.arco-tabs-nav-tab) {
  // 移入展示关闭icon
  .arco-tabs-tab-closable {
    svg {
      width: 0;
      transition: all 0.2s;
    }
    &:hover {
      svg {
        width: 1em;
      }
    }
  }

  // 消除tab移入的背景色
  &:hover .arco-tabs-tab-title::before {
    background: unset;
  }
}

// 消除tabs底部边线
:deep(.arco-tabs-nav) {
  &::before {
    background: unset;
  }
}
</style>
