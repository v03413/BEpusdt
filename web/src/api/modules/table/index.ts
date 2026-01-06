import axios from "@/api";

// 获取列表数据
export const getCommonTableListAPI = () => {
  return axios({
    url: "/mock/common-table/list",
    method: "get"
  });
};

// 获取列表数据
export const getCustomTableListAPI = () => {
  return axios({
    url: "/mock/custom-table/list",
    method: "get"
  });
};
