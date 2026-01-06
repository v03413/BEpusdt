import type { MockMethod } from "vite-plugin-mock";
import { deepClone, filterByDisable, buildTreeOptimized, treeSort, resultSuccess } from "../_utils";
import { systemMenu, permissionData } from "../_data/system_menu";

/**
 * 初始化
 * 后端接口需要做的事情：
 * 1、根据token判断角色
 * 2、过滤扁平路由，根据角色过滤权限和非禁用的节点
 * 3、将扁平路由转换为树结构
 * 4、路由树排序
 * 5、返回路由树
 *
 * 前端调用接口拿到路由数据后：
 * 1、将模块设置为真实模块
 * 2、存储路由树，用于生成菜单
 * 3、根据树生成一维路由数组
 * 4、动态添加路由
 * 5、缓存一维路由
 */

/**
 *  post请求body,get请求query
 *  /mock/menu/getRouters 获取路由数据
 *  /mock/menu/getMenuList 获取菜单列表数据-菜单管理
 *  /mock/menu/getUserPermission 根据角色获取权限数据
 *  /mock/menu/getPermission 获取权限数据
 */
export default [
  {
    url: "/mock/menu/getRouters",
    method: "get",
    timeout: 300,
    response: ({ headers }: any) => {
      let token = headers.authorization;
      // 根据token或登录账号判断是什么角色，这里模拟两个角色，admin、common
      let userRoles = token === "Admin-Token" ? ["admin"] : ["common"];
      const originMenu: any = deepClone(systemMenu);
      // 1. 过滤扁平路由，根据角色返回有权限且非禁用的节点
      const survivalTree = filterByDisable(originMenu, userRoles);
      // 2. 将扁平路由转换为树结构
      // 2. 给路由树排序
      // 3. 返回路由树
      return resultSuccess(treeSort(buildTreeOptimized(survivalTree)));
    }
  },
  {
    url: "/mock/menu/getMenuList",
    method: "get",
    timeout: 300,
    response: () => {
      // 菜单管理这里需要查看到所有菜单，无需判断角色权限、禁用菜单
      // 将扁平路由和权限菜单合并
      const originMenu: any = [...deepClone(systemMenu), ...deepClone(permissionData)];
      // 1. 将扁平路由转换为树结构
      // 2. 给路由树排序
      // 3. 返回路由树
      return resultSuccess(treeSort(buildTreeOptimized(originMenu)));
    }
  },
  {
    url: "/mock/menu/getUserPermission",
    method: "get",
    timeout: 300,
    response: ({ query }: any) => {
      let { role } = query;
      // 将扁平路由和权限菜单合并
      const originMenu: any = [...deepClone(systemMenu), ...deepClone(permissionData)];
      // 根据角色过滤id
      let idList = originMenu.filter((item: any) => item.meta.roles.includes(role)).map((item: any) => item.id);
      return resultSuccess(idList);
    }
  }
] as MockMethod[];
