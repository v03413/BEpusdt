import axios from "@/api";

export const listAPI = (data: any) => {
  return axios({
    url: "/api/order/list",
    method: "post",
    data
  });
};

export const detailAPI = (data: any) => {
  return axios({
    url: "/api/order/detail",
    method: "post",
    data
  });
};

export const paidAPI = (data: any) => {
  return axios({
    url: "/api/order/paid",
    method: "post",
    data
  });
};

export const delOrderApi = (data:any) => {
    return axios({
    url: "/api/order/del",
    method: "post",
    data
  });
}

export const batchDelOrderApi = (data:any) => {
    return axios({
    url: "/api/order/batch_del",
    method: "post",
    data
  });
}

export const manualNotifyAPI = (data: any) => {
  return axios({
    url: "/api/order/manual_notify",
    method: "post",
    data
  });
};
