import { defineStore } from "pinia";
import router from "@/router/index";
import { staticRoutes } from "@/router/route";
import { getRoutersAPI } from "@/api/modules/system/index";
import { moduleReplacement, linearArray } from "@/router/route-output";
import { getUrlWithParams } from "@/utils/index";

/**
 * 路由列表
 * @method setTabsTitle 设置tabs标签页名称
 * @methods setRoutePaths 设置路由名称集合
 * @methods setTabs 添加tabs标签页
 * @methods setCurrentRoute 设置系统内的当前路由
 * @methods removeTabsList 删除tabs页的指定路由
 * @methods removeRouteName 删除缓存路由名，用于取消页面缓存，单个删除
 * @methods removeRoutePaths 删除缓存路由名，用于取消页面缓存，批量删除
 * @methods resetRoute 重置动态添加的路由
 * @methods initSetRouter 路由初始化
 */
export const routeConfigStore = () => {
  const routeTree = ref<any>([]); // 有访问权限的路由树
  const routeList = ref<any>([]); // 有访问权限的一维路由数组
  const cacheRoutes = ref<string[]>([]); // 所有可缓存路由的路由路径
  const tabsList = ref<Menu.MenuOptions[]>([]); // 标签页数据
  const currentRoute = ref<Menu.MenuOptions | object>({}); // 当前路由

  /**
   * 设置tabs名称
   * @param {string} title 路由名称
   */
  function setTabsTitle(title: string) {
    for (let i = 0; i < tabsList.value.length; i++) {
      if (tabsList.value[i].path == getUrlWithParams()) {
        tabsList.value[i].meta.title = title;
        break;
      }
    }
  }

  /**
   * 设置可缓存路由的路径 -> fullPath
   * 同路由不同参兼容，会开启多个标签页，fullPath可做标识
   * @param {string} path 路由fullPath
   */
  function setRoutePaths(path: string) {
    let state = cacheRoutes.value.some((item: string) => item === path);
    if (state) return;
    cacheRoutes.value.push(path);
  }

  /**
   * 添加tabs标签页
   * @param {object} route 当前tabs路由
   */
  function setTabs(route: Menu.MenuOptions) {
    // 顶层手动添加的全屏静态路由不参与tabs
    if (staticRoutes.some(item => item.name == route.name)) return;
    // 不在可访问的路由中则不参与tabs
    if (!routeList.value.some((item: any) => item.name == route.name)) return;
    // 当前路由在tags中是否存在
    let index = tabsList.value.findIndex((item: Menu.MenuOptions) => item.path === route.path);
    // 不存在，直接缓存
    if (index == -1) return tabsList.value.push(route);
  }
  /**
   * 设置系统内的当前路由数据
   * @param {object} data 当前路由
   */
  function setCurrentRoute(data: Menu.MenuOptions) {
    // 名称一样不代表参数相同，这不用做已存在匹配，直接存储当前路由
    currentRoute.value = data;
  }
  /**
   * 删除tabs页的指定路由
   * @param {string} path 路由fullPath
   */
  function removeTabsList(path: string) {
    const index = tabsList.value.findIndex((item: Menu.MenuOptions) => item.path === path);
    if (tabsList.value[index].meta.affix) return;
    if (index === -1) return;
    tabsList.value.splice(index, 1);
  }
  /**
   * 删除缓存路由，用于取消页面缓存，单个删除
   * @param {string} path 路由
   */
  function removeRouteName(path: string) {
    const index = cacheRoutes.value.findIndex((item: string) => item === path);
    if (index === -1) return;
    cacheRoutes.value.splice(index, 1);
  }
  /**
   * 删除缓存路由，用于取消页面缓存，批量删除
   * @param {Array} list 路由
   */
  function removeRoutePaths(list: Array<string>) {
    cacheRoutes.value = cacheRoutes.value.filter((item: string) => !list.includes(item));
  }
  /**
   * 重置动态添加的路由
   */
  async function resetRoute() {
    // 清除标签页数据
    tabsList.value = [];
    // 清除有访问权限的路由树
    routeTree.value = [];
    // 清除所有可缓存路由的路由路径
    cacheRoutes.value = [];
    // 清除当前路由
    currentRoute.value = {};
    // 清除动态添加的路由
    routeList.value.forEach((item: any) => {
      if (router.hasRoute(item.name)) router.removeRoute(item.name);
    });
    // 清除有访问权限的一维路由数组
    routeList.value = [];
  }

  /**
   * 路由初始化
   * 1、获取过滤角色权限后排过序的的路由树，后端处理
   * 2、获取路由树转换的一维路由
   * 3、将模块设置为真实模块
   * 4、动态添加路由
   * 5、存储路由树，用于生成菜单
   * 6、缓存一维路由
   */
  async function initSetRouter() {
    // 1、获取过滤角色权限后的树，后端做排序处理
    let { data } = await getRoutersAPI();

    // 验证数据有效性
    if (!data || !Array.isArray(data) || data.length === 0) {
      return;
    }

    // 2、获取路由树转换的一维路由
    let flatRoute = linearArray(data);
    // 3、将模块设置为真实模块
    let realTree = await moduleReplacement(flatRoute);
    // 4、动态添加路由
    realTree.forEach((route: any) => {
      if (route.meta.isFull) {
        router.addRoute(route);
      } else {
        router.addRoute("layout", route);
      }
    });
    // 5、存储路由树，用于生成菜单
    routeTree.value = data;
    // 6、缓存一维路由
    routeList.value = flatRoute;
  }

  return {
    routeTree,
    routeList,
    cacheRoutes,
    tabsList,
    currentRoute,
    setTabsTitle,
    setRoutePaths,
    setTabs,
    setCurrentRoute,
    removeTabsList,
    removeRouteName,
    removeRoutePaths,
    resetRoute,
    initSetRouter
  };
};

export const useRouteConfigStore = defineStore("route-config", routeConfigStore);
