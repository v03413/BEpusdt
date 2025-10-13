<template>
  <a-menu
    :breakpoint="layoutType != 'layoutHead' ? 'xl' : undefined"
    :mode="'vertical'"
    :theme="asideDark ? 'dark' : 'light'"
    :collapsed="collapsed"
    :auto-scroll-into-view="true"
    :auto-open-selected="true"
    :accordion="isAccordion"
    :selected-keys="[selectedMenu]"
    @menu-item-click="onMenuItem"
  >
    <MenuItem :route-tree="props.routeTree" />
  </a-menu>
</template>

<script setup lang="ts">
import MenuItem from "@/layout/components/Menu/menu-item.vue";
import { useRoutingMethod } from "@/hooks/useRoutingMethod";
import { storeToRefs } from "pinia";
import { useThemeConfig } from "@/store/modules/theme-config";
const route = useRoute();
const router = useRouter();
const themeStore = useThemeConfig();
const { collapsed, isAccordion, layoutType, asideDark } = storeToRefs(themeStore);

interface Props {
  routeTree: Menu.MenuOptions[];
}
// props的数据类型
// type类型参考：https://cn.vuejs.org/guide/typescript/composition-api.html#typing-component-props
const props = withDefaults(defineProps<Props>(), {
  routeTree: () => []
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

<style lang="scss" scoped></style>
