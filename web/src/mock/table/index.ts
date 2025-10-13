import type { MockMethod } from "vite-plugin-mock";
import { resultSuccess } from "../_utils";

export default [
  {
    url: "/mock/common-table/list",
    method: "get",
    timeout: 500,
    response: () => {
      let list = [
        {
          key: "1",
          name: "陈思源",
          avatar: "陈",
          phone: "13812345678",
          email: "zhangsan@example.com",
          address: "北京市朝阳区",
          status: 1,
          createTime: "2024-05-27 09:00:00"
        },
        {
          key: "2",
          name: "李婉娴",
          avatar: "李",
          phone: "13987654321",
          email: "lisi@example.com",
          address: "上海市浦东新区",
          status: 0,
          createTime: "2024-05-26 15:30:00"
        },
        {
          key: "3",
          name: "王雨菲",
          avatar: "王",
          phone: "13666666666",
          email: "wangwu@example.com",
          address: "广州市天河区",
          status: 1,
          createTime: "2024-05-25 12:45:00"
        },
        {
          key: "4",
          name: "张晨曦",
          avatar: "张",
          phone: "13788888888",
          email: "zhaoliu@example.com",
          address: "深圳市福田区",
          status: 0,
          createTime: "2024-05-24 11:20:00"
        },
        {
          key: "5",
          name: "赵梦琪",
          avatar: "赵",
          phone: "13599999999",
          email: "qianqi@example.com",
          address: "成都市锦江区",
          status: 1,
          createTime: "2024-05-23 14:10:00"
        },
        {
          key: "6",
          name: "刘昊然",
          avatar: "刘",
          phone: "13377777777",
          email: "sunba@example.com",
          address: "杭州市西湖区",
          status: 0,
          createTime: "2024-05-22 10:05:00"
        },
        {
          key: "7",
          name: "孙梦洁",
          avatar: "孙",
          phone: "13266666666",
          email: "zhoujiu@example.com",
          address: "南京市鼓楼区",
          status: 1,
          createTime: "2024-05-21 08:45:00"
        },
        {
          key: "8",
          name: "黄俊杰",
          avatar: "黄",
          phone: "13155555555",
          email: "wushi@example.com",
          address: "重庆市渝中区",
          status: 0,
          createTime: "2024-05-20 16:30:00"
        },
        {
          key: "9",
          name: "周雨萱",
          avatar: "周",
          phone: "13044444444",
          email: "zhengshiyi@example.com",
          address: "武汉市江汉区",
          status: 1,
          createTime: "2024-05-19 09:20:00"
        },
        {
          key: "10",
          name: "韩雪儿",
          avatar: "韩",
          phone: "13933333333",
          email: "kongshier@example.com",
          address: "西安市雁塔区",
          status: 0,
          createTime: "2024-05-18 13:55:00"
        }
      ];
      let data = {
        list,
        total: list.length
      };
      return resultSuccess(data);
    }
  },
  {
    url: "/mock/custom-table/list",
    method: "get",
    timeout: 500,
    response: () => {
      let list = [
        {
          key: "1",
          serial: "57",
          name: "WEIJEIWGSDG",
          content: "图文",
          searchType: "阅读量",
          contentSize: 3700,
          status: 1,
          createTime: "2024-05-27 09:00:00"
        },
        {
          key: "2",
          serial: "58",
          name: "DHDWEH",
          content: "图文",
          searchType: "阅读量",
          contentSize: 3700,
          status: 0,
          createTime: "2024-05-26 15:30:00"
        },
        {
          key: "3",
          serial: "59",
          name: "SHEHWAE",
          content: "视频",
          searchType: "播放量",
          contentSize: 2500,
          status: 1,
          createTime: "2024-05-25 12:45:00"
        },
        {
          key: "4",
          serial: "60",
          name: "WEHJUK",
          content: "视频",
          searchType: "播放量",
          contentSize: 2500,
          status: 0,
          createTime: "2024-05-24 11:20:00"
        },
        {
          key: "5",
          serial: "61",
          name: "EYTUJDGSD",
          content: "音频",
          searchType: "播放量",
          contentSize: 2500,
          status: 1,
          createTime: "2024-05-23 14:10:00"
        },
        {
          key: "6",
          serial: "62",
          name: "RTIFSDF",
          content: "图文",
          searchType: "阅读量",
          contentSize: 3700,
          status: 0,
          createTime: "2024-05-22 10:05:00"
        },
        {
          key: "7",
          serial: "63",
          name: "OIKSHHFD",
          content: "图文",
          searchType: "阅读量",
          contentSize: 3700,
          status: 1,
          createTime: "2024-05-21 08:45:00"
        },
        {
          key: "8",
          serial: "64",
          name: "QEUHRKO",
          content: "视频",
          searchType: "播放量",
          contentSize: 2500,
          status: 0,
          createTime: "2024-05-20 16:30:00"
        },
        {
          key: "9",
          serial: "65",
          name: "JSGSDRI",
          content: "视频",
          searchType: "播放量",
          contentSize: 2300,
          status: 1,
          createTime: "2024-05-19 09:20:00"
        },
        {
          key: "10",
          serial: "66",
          name: "AOUSTHREGG",
          content: "音频",
          searchType: "播放量",
          contentSize: 2300,
          status: 0,
          createTime: "2024-05-18 13:55:00"
        }
      ];
      let data = {
        list,
        total: list.length
      };
      return resultSuccess(data);
    }
  }
] as MockMethod[];
