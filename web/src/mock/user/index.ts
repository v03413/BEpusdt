import type { MockMethod } from "vite-plugin-mock";
import { resultSuccess, resultError, deepClone } from "../_utils";
import { accountData } from "../_data/system_data";
import { permissionData } from "../_data/system_menu";
import { roleData } from "../_data/system_data";
/**
 * post请求body,get请求query
 * /mock/login 用户登录
 * /mock/user/getUserInfo 获取用户信息
 */
export default [
  {
    url: "/mock/login",
    method: "post",
    timeout: 300,
    response: ({ body }: any) => {
      let { username, password } = body;
      // 根据登录账号返回token，这里模拟两个账号
      // 1、admin账号
      // 2、common账号
      if (username === "admin" && password === "123456") {
        return resultSuccess({ token: "Admin-Token" });
      }
      if (username === "common" && password === "123456") {
        return resultSuccess({ token: "Common-Token" });
      }
      return resultError(null, "账号或者密码错误", 500);
    }
  },
  {
    url: "/mock/user/getUserInfo",
    method: "get",
    timeout: 300,
    response: ({ headers, query }: any) => {
      let data: any = {
        user: {}, // 账户信息
        roles: [], // 角色
        permissions: [] // 权限
      };
      let account = deepClone(accountData);

      // 传入id则通过id匹配，否则通过token匹配
      if (query.id) {
        data.user = account.find((item: any) => item.id == query.id);
      } else {
        // 1、拿账户信息
        if (headers.authorization == "Admin-Token") {
          // 管理员账户
          data.user = account.find((item: any) => item.id == 1);
        } else {
          // 普通账户
          data.user = account.find((item: any) => item.id == 2);
        }
      }

      // 2、根据账户拿角色信息
      let roleInfo = roleData.filter((item: any) => data.user.roles.includes(item.code));
      data.user.roles = roleInfo;
      data.roles = roleInfo.map((item: any) => item.code);

      // 3、拿权限信息
      // 超级管理员
      if (data.roles.includes("admin")) {
        data.permissions = ["*:*:*"];
      } else {
        // 非admin账号则只拿对应角色的权限
        let list = permissionData.filter((permission: any) => {
          return data.roles.some((item: string) => permission.meta.roles.includes(item));
        });
        data.permissions = list.map(item => item.meta.permission);
      }

      return resultSuccess(data);
    }
  }
] as MockMethod[];
