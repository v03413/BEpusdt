import axios from "@/api";

export const getRateListAPI = (data: any) => {
  return axios({
    url: "/api/rate/list",
    method: "post",
    data
  });
};

export const getSyntaxListAPI = () => {
  return axios({
    url: "/api/rate/syntax",
    method: "post"
  });
};

export const setSyntaxAPI = (data: any) => {
  return axios({
    url: "/api/rate/set_syntax",
    method: "post",
    data
  });
};

export const onSyncAPI = (data: any) => {
  return axios({
    url: "/api/rate/sync",
    method: "post",
    data
  });
};
