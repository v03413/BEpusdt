// 在线用户
export const onlineuser = [
  {
    tokenId: "2t6c278e-495h-821t-as26-31c98f8d2379",
    userName: "admin",
    deptName: "研发部门",
    ipaddr: "115.23.16.25",
    loginLocation: "上海市 上海市",
    browser: "Quark 7.9.2.771",
    os: "HarmonyOS NEXT 5.0.0.102",
    loginTime: "2025-04-12 10:30:59",
    status: 1 // 0离线 1在线
  },
  {
    tokenId: "37adftaf-172c-1527-dcf6-efa3f92adc5c",
    userName: "admin",
    deptName: "研发部门",
    ipaddr: "112.21.152.215",
    loginLocation: "浙江省 杭州市",
    browser: "Chrome 135.0.0.0",
    os: "Windows 11",
    loginTime: "2025-04-12 10:30:59",
    status: 1 // 0离线 1在线
  }
];

// 定时任务
export const crontab = [
  {
    id: 1,
    name: "cron任务，5秒执行一次", // 任务名称
    cron: "0/5 * * * * * ", // cron-执行规则
    every: null, // 间隔时间秒-执行规则
    service: "taskDemoService.test()", // 执行服务
    startDate: "2025-03-30 11:28:35", // 开始时间
    nextRunTime: "2025-03-13 11:20:50", // 下次执行时间
    taskType: 0, // 任务类型 0 cron 1 指定时间(秒)
    misfirePolicy: 1, // 执行策略 1 循环 2 执行一次
    remark: "这是定时任务", // 备注
    status: 1, // 0离线 1在线
    createBy: "admin",
    createTime: "2024-06-30 11:28:35",
    updateBy: null,
    updateTime: null
  },
  {
    id: 2,
    name: "每秒执行一次",
    cron: null,
    every: 1,
    service: "taskDemoService.test(1,2)",
    startDate: "2025-03-30 11:28:35",
    nextRunTime: "2025-03-13 11:20:50",
    taskType: 1,
    misfirePolicy: 1,
    remark: "这是定时任务",
    status: 1,
    createBy: "admin",
    createTime: "2024-06-30 11:28:35",
    updateBy: null,
    updateTime: null
  }
];

// 任务日志
export const crontabLogs = [
  {
    id: 1,
    taskId: 1,
    taskName: "cron任务，5秒执行一次",
    status: 1,
    remark: "",
    createTime: "2025-03-30 11:28:35",
    executionTime: "2025-03-30 11:28:35"
  },
  {
    id: 2,
    taskId: 2,
    taskName: "每秒执行一次",
    status: 0,
    remark:
      "You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use near ', `_isUTC` = false, `_pf` = '[object Object]', `_locale` = '[object Object]', `_' at line 1",
    createTime: "2025-03-30 11:28:35",
    executionTime: "2025-03-30 11:28:35"
  }
];
