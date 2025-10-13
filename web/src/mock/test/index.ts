import type { MockMethod } from "vite-plugin-mock";
import { resultSuccess, resultError } from "../_utils";

export default [
  {
    url: "/mock/test/success",
    method: "get",
    timeout: 300,
    response: () => {
      return resultSuccess("SnowAdmin开源免费，Welcome!");
    }
  },
  {
    url: "/mock/test/fail",
    method: "get",
    timeout: 300,
    response: () => {
      return resultError(null, "模拟mock数据出错", 500);
    }
  }
] as MockMethod[];
