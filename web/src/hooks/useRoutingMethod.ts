import pinia from "@/store/index";
import { storeToRefs } from "pinia";
import { useRouteConfigStore } from "@/store/modules/route-config";
import { findCategoryById, findPathOfParentNode } from "@/utils/tree-tools";

/**
 * 路由处理hooks，内置多种路由处理场景
 * @returns 路由方法
 */
export const useRoutingMethod = () => {
  /**
   * 从一维路由中查找路由
   * @param {string} path 路由的path
   * @returns 查找到的路由，undefined则表示未找到
   */
  const findLinearArray = (path: string) => {
    const routerStore = useRouteConfigStore(pinia);
    const { routeTree } = storeToRefs(routerStore);
    return findCategoryById(routeTree.value, "path", path);
  };

  /**
   * 根据当前路由找到所有直属父级路由
   * @param {string} path 路由的path
   * @returns 查找到的所有父级路由，未找到则null
   */
  const getAllParentRoute = (path: string) => {
    const routerStore = useRouteConfigStore(pinia);
    const { routeTree } = storeToRefs(routerStore);
    return findPathOfParentNode(routeTree.value, "path", path);
  };

  /**
   * 从一维路由中判断路由是否存在
   * @param {string} key 路由的name
   * @returns 路由是否存在，true存在 false不存在
   */
  const hasRoute = (key: string) => {
    const routerStore = useRouteConfigStore(pinia);
    const { routeList } = storeToRefs(routerStore);
    return routeList.value.some((item: Menu.MenuOptions) => item.name == key);
  };

  /**
   * 从tabs路由中查找路由
   * @param {string} key 路由的name
   * @returns 查找到的路由，undefined则表示未找到
   */
  const findTagsList = (key: string) => {
    const routerStore = useRouteConfigStore(pinia);
    const { tabsList } = storeToRefs(routerStore);
    return tabsList.value.find((item: Menu.MenuOptions) => item.name == key);
  };

  /**
   * 处理外链跳转，打开一个新窗口并根据url跳转
   * @param {any} route 路由
   */
  const openExternalLinks = (route: any) => {
    // 处理外链跳转
    if (route.meta.link && !route.meta.iframe) {
      window.open(route.meta.link as string, "_blank");
    }
  };

  /**
   * 检测是否是动态匹配路由，如果是动态匹配路由，则path必然带有"/:"字样，例如：/user/:id
   * @param {string} path 路由path
   * @returns 是否是动态匹配路由
   */
  const isDynamicRoute = (path: string) => {
    return path.includes("/:");
  };

  return {
    findLinearArray,
    getAllParentRoute,
    findTagsList,
    openExternalLinks,
    isDynamicRoute,
    hasRoute
  };
};
