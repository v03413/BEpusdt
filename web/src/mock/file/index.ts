import type { MockMethod } from "vite-plugin-mock";
import { resultSuccess } from "../_utils";

export default [
  {
    url: "/mock/document-library/tree",
    method: "get",
    timeout: 500,
    response: () => {
      let list = [
        {
          title: "总公司",
          key: "0-0",
          children: [
            {
              title: "行政部门",
              key: "0-0-0",
              children: [
                {
                  title: "人力资源部",
                  key: "0-0-0-0",
                  children: [{ title: "员工出差申请表", key: "0-0-0-0-0" }]
                },
                {
                  title: "财务部",
                  key: "0-0-0-1",
                  children: [{ title: "行政费用预算分析", key: "0-0-0-0-1" }]
                },
                {
                  title: "法务部",
                  key: "0-0-0-2",
                  children: [{ title: "公司章程", key: "0-0-0-0-2" }]
                },
                { title: "行政工作计划", key: "0-0-0-3" },

                { title: "会议室预订统计", key: "0-0-0-4" },
                { title: "月度工作总结", key: "0-0-0-5" }
              ]
            },
            {
              title: "业务部门",
              key: "0-0-1",
              children: [
                { title: "客户反馈集锦", key: "0-0-1-0" },
                { title: "市场分析报告-季度", key: "0-0-1-1" }
              ]
            },
            {
              title: "技术部门",
              key: "0-0-2",
              children: [
                { title: "代码审查记录", key: "0-0-2-0" },
                { title: "数据库优化方案", key: "0-0-2-1" },
                { title: "技术支持工单", key: "0-0-2-2" },
                { title: "安全漏洞扫描报告", key: "0-0-2-3" },
                { title: "系统性能调优总结", key: "0-0-2-4" },
                { title: "软件版本发布计划", key: "0-0-2-5" },
                { title: "技术培训材料", key: "0-0-2-6" },
                { title: "技术文档更新计划", key: "0-0-2-7" }
              ]
            },
            {
              title: "后勤部门",
              key: "0-0-3",
              children: [
                { title: "设施维护报告", key: "0-0-3-0" },
                { title: "物流运输安排", key: "0-0-3-1" },
                { title: "仓库库存清单", key: "0-0-3-2" },
                { title: "办公室布局设计", key: "0-0-3-3" },
                { title: "设备采购申请表", key: "0-0-3-4" },
                { title: "食品安全检查记录", key: "0-0-3-5" }
              ]
            }
          ]
        }
      ];
      return resultSuccess(list);
    }
  },
  {
    url: "/mock/document-library/table",
    method: "get",
    timeout: 500,
    response: () => {
      let list = [
        {
          key: "1",
          source: "bilibili",
          sourceSvg: "bilibili",
          sourceType: "分享",
          flow: "5015",
          percent: "75",
          status: 1,
          createTime: "2024-05-27 09:00:00"
        },
        {
          key: "2",
          source: "微信",
          sourceSvg: "wechat",
          sourceType: "分享",
          flow: "3217",
          percent: "60",
          status: 0,
          createTime: "2024-05-26 15:30:00"
        },
        {
          key: "3",
          source: "QQ音乐",
          sourceSvg: "QQ音乐",
          sourceType: "搜索",
          flow: "1343",
          percent: "55",
          status: 1,
          createTime: "2024-05-25 12:45:00"
        },
        {
          key: "4",
          source: "百度",
          sourceSvg: "百度",
          sourceType: "搜索",
          flow: "1235",
          percent: "20",
          status: 0,
          createTime: "2024-05-24 11:20:00"
        },
        {
          key: "5",
          source: "优酷",
          sourceSvg: "优酷",
          sourceType: "首页推荐",
          flow: "3456",
          percent: "20",
          status: 1,
          createTime: "2024-05-23 14:10:00"
        },
        {
          key: "6",
          source: "网易云音乐",
          sourceSvg: "网易云音乐",
          sourceType: "每日一推",
          flow: "5873",
          percent: "30",
          status: 0,
          createTime: "2024-05-22 10:05:00"
        },
        {
          key: "7",
          source: "抖音",
          sourceSvg: "抖音",
          sourceType: "自然",
          flow: "5465",
          percent: "86",
          status: 1,
          createTime: "2024-05-21 08:45:00"
        },
        {
          key: "8",
          source: "西瓜视频",
          sourceSvg: "西瓜视频",
          sourceType: "搜索",
          flow: "4642",
          percent: "45",
          status: 0,
          createTime: "2024-05-20 16:30:00"
        },
        {
          key: "9",
          source: "微博",
          sourceSvg: "微博",
          sourceType: "搜索",
          flow: "2375",
          percent: "60",
          status: 1,
          createTime: "2024-05-19 09:20:00"
        },
        {
          key: "10",
          source: "豆瓣",
          sourceSvg: "豆瓣",
          sourceType: "搜索",
          flow: "2465",
          percent: "40",
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
