import axios from "@/api";

// 获取在线用户
export const getOnlineuserAPI = () => {
  return axios({
    url: "/mock/monitor/getOnlineuser",
    method: "get"
  });
};

// 获取定时任务
export const getCrontabAPI = () => {
  return axios({
    url: "/mock/monitor/getCrontab",
    method: "get"
  });
};

// 获取定时任务日志
export const getCrontabLogsAPI = (params: any) => {
  return axios({
    url: "/mock/monitor/getCrontabLogs",
    method: "get",
    params
  });
};
