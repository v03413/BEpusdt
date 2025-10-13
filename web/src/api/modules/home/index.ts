import axios from "@/api";

export const getDashboardHomeAPI = (data: any) => {
  return axios({
    url: "/api/dashboard/home",
    method: "post",
    data
  });
};
