import axios from "@/api";

export const listAPI = (data: any) => {
  return axios({
    url: "/api/order/list",
    method: "post",
    data
  });
};

export const orderBinListAPI = (data: any) => {
  return axios({
    url: "/api/order/order_bin_list",
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

export const move2BinApi = (data:any) => {
    return axios({
    url: "/api/order/move_to_bin",
    method: "post",
    data
  });
}
export const removeFromBin = (data:any) => {
    return axios({
    url: "/api/order/remove_form_bin",
    method: "post",
    data
  });
}
export const emptyBin = () => {
    return axios({
    url: "/api/order/empty_bin",
    method: "post"
  });
}
export const manualNotifyAPI = (data: any) => {
  return axios({
    url: "/api/order/manual_notify",
    method: "post",
    data
  });
};
