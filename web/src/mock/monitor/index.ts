import type { MockMethod } from "vite-plugin-mock";
import { resultSuccess } from "../_utils";
import { onlineuser, crontab, crontabLogs } from "../_data/monitor_data";
/**
 * post请求body,get请求query
 * /mock/monitor/getOnlineuser 获取在线用户
 * /mock/monitor/getCrontab 定时任务
 */
export default [
  {
    url: "/mock/monitor/getOnlineuser",
    method: "get",
    timeout: 300,
    response: () => {
      return resultSuccess(onlineuser);
    }
  },
  {
    url: "/mock/monitor/getCrontab",
    method: "get",
    timeout: 300,
    response: () => {
      return resultSuccess(crontab);
    }
  },
  {
    url: "/mock/monitor/getCrontabLogs",
    method: "get",
    timeout: 300,
    response: ({ query }: any) => {
      console.log(query);
      return resultSuccess(crontabLogs);
    }
  }
] as MockMethod[];
