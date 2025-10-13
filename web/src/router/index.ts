import NProgress from "@/config/nprogress";
import pinia from "@/store/index";
import { createRouter, createWebHashHistory } from "vue-router";
import { staticRoutes, notFoundAndNoPower } from "@/router/route.ts";
import { currentlyRoute } from "@/router/route-output";
import { storeToRefs } from "pinia";
import { useUserInfoStore } from "@/store/modules/user-info";
import { useRouteConfigStore } from "@/store/modules/route-config";
import { useRoutingMethod } from "@/hooks/useRoutingMethod";
/**
 * 创建vue的路由示例
 * @method createRouter(options: RouterOptions): Router
 * @link 参考：https://next.router.vuejs.org/zh/api/#createrouter
 */
const router = createRouter({
  history: createWebHashHistory(),
  /**
   * 设置静态路由，其它的路由通过addRoute动态添加
   * 1、staticRoutes登录页、layout页、默认页面('/')
   * 2、notFoundAndNoPower 添加默认 401、500界面，防止提示 No match found for location with path 'xxx'
   * 2、后端控制路由中也需要添加 notFoundAndNoPower 401、500界面
   * 静态添加 notFoundAndNoPower 401、500界面将全屏显示
   * 如果要 notFoundAndNoPower 在layout容器展示，则需要移除静态添加并将其添加到缓存路由树
   */
  routes: [...staticRoutes, ...notFoundAndNoPower]
});

/**
 * 路由加载前需要判断用户是否登录
 * 1、去登录页，无token，放行
 * 2、没有token，直接重定向到登录页
 * 3、去登录页，有token，直接重定向到home页
 * 4、去非登录页，有token，用户信息是否存在，有则放行，否则重新获取路由信息、初始化路由
 * 注意：
 * 全局routeTree不能持久化缓存
 * 页面刷新会导致addRoute动态添加的路由失效，需要重新初始化路由
 */
router.beforeEach(async (to: any, _: any, next: any) => {
  NProgress.start(); // 开启进度条
  const store = useUserInfoStore(pinia);
  const routeStore = useRouteConfigStore(pinia);
  const { token } = storeToRefs(store);
  const { routeTree } = storeToRefs(routeStore);
  // console.log("去", to, "来自", from);
  // next()内部加了path等于跳转指定路由会再次触发router.beforeEach，内部无参数等于放行，不会触发router.beforeEach

  // 1、去登录页，无token，放行
  if (to.path === "/login" && !token.value) return next();

  // 2、没有token，直接重定向到登录页
  if (!token.value) return next("/login");

  // 3、去登录页，有token，直接重定向到home页
  if (to.path === "/login" && token.value) {
    // 项目内的跳转，处理跳转路由高亮
    currentlyRoute(to);
    return next("/home");
  }

  // 4、去非登录页，有token，用户信息是否存在，有则放行，否则重新获取路由信息、初始化路由
  // 判断路由是否获取，先获取账号信息和路由信息，添加路由后再跳转(页面刷新时触发)
  // 解决刷新页面404的问题
  if (!routeTree.value.length) {
    const routeStore = useRouteConfigStore(pinia);

    try {
      // 获取账号信息
      await store.setAccount();
      // 获取路由信息
      await routeStore.initSetRouter();

      // 判断是否是动态路由
      const { isDynamicRoute } = useRoutingMethod();
      if (isDynamicRoute(to.path)) {
        return next({ name: to.name, params: to.params });
      } else {
        return next({ path: to.path, query: to.query });
      }
    } catch (error) {
      console.error(error);

      return next(false);
    }
  }

  // 获取外链路由的处理函数
  // 所有的路由正常放行，只不过额外判断是否是外链，如果是，则打开新窗口跳转外链
  // 外链的页面依旧正常打开，只不过不会参与缓存与tabs显示，符合路由跳转的直觉
  const { openExternalLinks } = useRoutingMethod();
  // 处理外链跳转
  openExternalLinks(to);
  // 项目内的跳转，处理跳转路由高亮
  currentlyRoute(to);
  // 动态路由添加过走这里，直接放行
  next();
});

// 路由跳转错误
router.onError((error: any) => {
  NProgress.done();
  console.warn("路由错误", error.message);
});

// 路由加载后
router.afterEach(() => {
  NProgress.done();
});

export default router;
