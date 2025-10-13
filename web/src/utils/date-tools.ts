/**
 * 时间戳转 年月日时分秒
 * 若timestamp不传则取当前时间
 * 若type不传则取 年月日时分秒
 * @param { number } timestamp 时间戳
 * @return { string } 返回年月日时分秒字符串
 */
export const getTimestamp = (timestamp: string | number | null, type: string): string => {
  let date = null;
  if (timestamp) {
    date = new Date(timestamp);
  } else {
    date = new Date();
  }
  let Year = String(date.getFullYear());
  let Moth = String(date.getMonth() + 1).padStart(2, "0");
  let Day = String(date.getDate()).padStart(2, "0");
  let Hour = String(date.getHours()).padStart(2, "0");
  let Minute = String(date.getMinutes()).padStart(2, "0");
  let Seconds = String(date.getSeconds()).padStart(2, "0");
  if (type === "yyyy") {
    return `${Year}`;
  }
  if (type === "yyyy-MM") {
    return `${Year}-${Moth}`;
  }
  if (type === "yyyy-MM-dd") {
    return `${Year}-${Moth}-${Day}`;
  }
  return `${Year}-${Moth}-${Day} ${Hour}:${Minute}:${Seconds}`;
};

/**
 * 时间字符串转时间戳
 * @param { string } str 时间字符串，如：2025-06-28
 * @return 返回时间戳
 */
export const getTimes = (str: string) => {
  return new Date(str.substring(0, 19).replace(/-/g, "/")).getTime();
};

/**
 * 获取当前月份的第一天、当天和最后一天
 * @param {string} YearandMonty 例如: '2025-09'；参数可传可不传，不传获取当月，传则获取指定月
 * @return {array} [当月第一天,当前天,最后一天]
 **/
export const getFirstandLastDay = (YearandMonty?: string): Array<string> => {
  let date = YearandMonty ? new Date(YearandMonty) : new Date();
  let year = String(date.getFullYear());
  let month = String(date.getMonth() + 1).padStart(2, "0");
  let day = String(date.getDate()).padStart(2, "0");
  // new Date第三个参数默认1，也就是当月第一天,如果传0则是当月最后一天
  let lastDay = new Date(Number(year), Number(month), 0).getDate();
  let firstDate = `${year}-${month}-01`;
  let theDay = `${year}-${month}-${day}`;
  let lastDate = `${year}-${month}-${lastDay}`;
  return [firstDate, theDay, lastDate];
};

/**
 * 根据当前日期前推指定日期
 * @param { number } days 需要前推的天数
 * @param { string } type 返回的时间格式
 * @return {Array[]} 是否为空对象 [前推天数的日期，当前日期]
 */
export const getDatesForwardDate = (days = 0, type = "YYYY-MM-DD"): Array<string> => {
  const today = new Date();
  const firstDay: any = new Date(today);
  firstDay.setDate(firstDay.getDate() - days); // 向前推指定天，得到前指定天数的第一天

  const lastDay: any = new Date(today);
  lastDay.setDate(lastDay.getDate() - 1); // 昨天是前指定天数的最后一天

  const firstDayFormatted = getTimestamp(firstDay, type);
  const lastDayFormatted = getTimestamp(lastDay, type);
  return [firstDayFormatted, lastDayFormatted];
};

/**
 * 获取当前日期是第几周
 * @param dateTime 当前传入的日期值
 * @description 示例：getWeek(new Date()); 返回25，说明今年第25周
 * @returns 返回第几周数字值
 */
export const getWeek = (dateTime: Date): number => {
  let temptTime = new Date(dateTime.getTime());
  // 周几
  let weekday = temptTime.getDay() || 7;
  // 周1+5天=周六
  temptTime.setDate(temptTime.getDate() - weekday + 1 + 5);
  let firstDay = new Date(temptTime.getFullYear(), 0, 1);
  let dayOfWeek = firstDay.getDay();
  let spendDay = 1;
  if (dayOfWeek != 0) spendDay = 7 - dayOfWeek + 1;
  firstDay = new Date(temptTime.getFullYear(), 0, 1 + spendDay);
  let d = Math.ceil((temptTime.valueOf() - firstDay.valueOf()) / 86400000);
  let result = Math.ceil(d / 7);
  return result;
};

/**
 * 时间日期转换
 * @param date 当前时间，new Date() 格式
 * @param format 需要转换的时间格式字符串
 * @description format 字符串随意，如 `YYYY-mm、YYYY-mm-dd`
 * @description format 季度："YYYY-mm-dd HH:MM:SS QQQQ"
 * @description format 星期："YYYY-mm-dd HH:MM:SS WWW"
 * @description format 几周："YYYY-mm-dd HH:MM:SS ZZZ"
 * @description format 季度 + 星期 + 几周："YYYY-mm-dd HH:MM:SS WWW QQQQ ZZZ"
 * @returns 返回拼接后的时间字符串
 */
export const formatDate = (date: Date, format: string): string => {
  let we = date.getDay(); // 星期
  let z = getWeek(date); // 周
  let qut = Math.floor((date.getMonth() + 3) / 3).toString(); // 季度
  const opt: { [key: string]: string } = {
    "Y+": date.getFullYear().toString(), // 年
    "m+": (date.getMonth() + 1).toString(), // 月(月份从0开始，要+1)
    "d+": date.getDate().toString(), // 日
    "H+": date.getHours().toString(), // 时
    "M+": date.getMinutes().toString(), // 分
    "S+": date.getSeconds().toString(), // 秒
    "q+": qut // 季度
  };
  // 中文数字 (星期)
  const week: { [key: string]: string } = {
    "0": "日",
    "1": "一",
    "2": "二",
    "3": "三",
    "4": "四",
    "5": "五",
    "6": "六"
  };
  // 中文数字（季度）
  const quarter: { [key: string]: string } = {
    "1": "一",
    "2": "二",
    "3": "三",
    "4": "四"
  };
  if (/(W+)/.test(format))
    format = format.replace(
      RegExp.$1,
      RegExp.$1.length > 1 ? (RegExp.$1.length > 2 ? "星期" + week[we] : "周" + week[we]) : week[we]
    );
  if (/(Q+)/.test(format))
    format = format.replace(RegExp.$1, RegExp.$1.length == 4 ? "第" + quarter[qut] + "季度" : quarter[qut]);
  if (/(Z+)/.test(format)) format = format.replace(RegExp.$1, RegExp.$1.length == 3 ? "第" + z + "周" : z + "");
  for (let k in opt) {
    let r = new RegExp("(" + k + ")").exec(format);
    // 若输入的长度不为1，则前面补零
    if (r) format = format.replace(r[1], RegExp.$1.length == 1 ? opt[k] : opt[k].padStart(RegExp.$1.length, "0"));
  }
  return format;
};

/**
 * 将时间转换为 `几秒前`、`几分钟前`、`几小时前`、`几天前`
 * @param param 当前时间，new Date() 格式或者字符串时间格式
 * @param format 需要转换的时间格式字符串
 * @description param 10秒：  10 * 1000
 * @description param 1分：   60 * 1000
 * @description param 1小时： 60 * 60 * 1000
 * @description param 24小时：60 * 60 * 24 * 1000
 * @description param 3天：   60 * 60* 24 * 1000 * 3
 * @returns 返回转换后的字符串
 */
export const formatPast = (param: string | Date, format: string = "YYYY-mm-dd"): string => {
  // 传入格式处理、存储转换值
  let t: any, s: number;
  // 获取js 时间戳
  let time: number = new Date().getTime();
  // 是否是对象
  typeof param === "string" || "object" ? (t = new Date(param).getTime()) : (t = param);
  // 当前时间戳 - 传入时间戳
  time = Number.parseInt(`${time - t}`);
  if (time < 10000) {
    // 10秒内
    return "刚刚";
  } else if (time < 60000 && time >= 10000) {
    // 超过10秒少于1分钟内
    s = Math.floor(time / 1000);
    return `${s}秒前`;
  } else if (time < 3600000 && time >= 60000) {
    // 超过1分钟少于1小时
    s = Math.floor(time / 60000);
    return `${s}分钟前`;
  } else if (time < 86400000 && time >= 3600000) {
    // 超过1小时少于24小时
    s = Math.floor(time / 3600000);
    return `${s}小时前`;
  } else if (time < 259200000 && time >= 86400000) {
    // 超过1天少于3天内
    s = Math.floor(time / 86400000);
    return `${s}天前`;
  } else {
    // 超过3天
    let date = typeof param === "string" || "object" ? new Date(param) : param;
    return formatDate(date, format);
  }
};

/**
 * 时间问候语
 * @param param 当前时间，new Date() 格式
 * @description param 调用 `formatAxis(new Date())` 输出 `上午好`
 * @returns 返回根据时间处理后的问候语
 */
export function formatAxis(param: Date): string {
  let hour: number = new Date(param).getHours();
  if (hour < 6) return "凌晨好";
  else if (hour < 9) return "早上好";
  else if (hour < 12) return "上午好";
  else if (hour < 14) return "中午好";
  else if (hour < 17) return "下午好";
  else if (hour < 19) return "傍晚好";
  else if (hour < 22) return "晚上好";
  else return "夜里好";
}

/**
 * 秒转时分秒
 * @param {number} seconds 秒时间戳
 * @returns {string} 秒转时分秒后的字符串，如：00小时11分01秒
 */
export const formatTime = (seconds: number) => {
  let hours = Math.floor(seconds / 3600);
  let minutes = Math.floor((seconds - hours * 3600) / 60);
  let remainingSeconds = seconds - hours * 3600 - minutes * 60;
  return `${String(hours).padStart(2, "0")}小时${String(minutes).padStart(2, "0")}分${String(remainingSeconds).padStart(2, "0")}秒`;
};

/**
 * 判断开始时间和结束时间是否大于12个月
 * @param {string} startTime 开始时间，如：'2023-09-14'
 * @param {string} endTime 结束时间，如：'2023-12-10'
 * @return {boolean} true\false
 */
export const dateLimit = (startTime: string, endTime: string): boolean => {
  const start = new Date(startTime);
  const end = new Date(endTime);

  // 计算年差和月差
  const yearDiff = end.getFullYear() - start.getFullYear();
  const monthDiff = end.getMonth() - start.getMonth();

  // 总月差 = 年差*12 + 月差
  let totalMonths = yearDiff * 12 + monthDiff;

  // 如果日期部分 end 小于 start，总月差减1
  if (end.getDate() < start.getDate()) {
    totalMonths--;
  }

  return totalMonths > 12;
};
