import type { MockMethod } from "vite-plugin-mock";
import { resultSuccess } from "../_utils";
import { dictData, divisionData, roleData, accountData } from "../_data/system_data";

/**
 * post请求body,get请求query
 * /mock/system/getDict 获取字典数据
 * /mock/system/getDivision 获取部门数据
 * /mock/system/getRole 获取角色数据
 * /mock/system/getAccount 获取账号数据
 */
export default [
  {
    url: "/mock/system/getDict",
    method: "get",
    timeout: 300,
    response: () => {
      return resultSuccess(dictData);
    }
  },
  {
    url: "/mock/system/getDivision",
    method: "get",
    timeout: 300,
    response: () => {
      return resultSuccess(divisionData);
    }
  },
  {
    url: "/mock/system/getRole",
    method: "get",
    timeout: 300,
    response: () => {
      return resultSuccess(roleData);
    }
  },
  {
    url: "/mock/system/getAccount",
    method: "get",
    timeout: 300,
    response: () => {
      return resultSuccess(accountData);
    }
  }
] as MockMethod[];
