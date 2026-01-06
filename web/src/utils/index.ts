/**
 * 深拷贝
 * @param { string } data 需要深拷贝的数据
 * @returns 深拷贝的数据
 */
export function deepClone(data: any) {
  let stack = [];
  let cloned;
  if (Array.isArray(data)) {
    cloned = [];
  } else if (typeof data === "object" && data !== null) {
    cloned = {};
  } else {
    return data;
  }
  stack.push({
    original: data,
    copy: cloned
  });
  while (stack.length > 0) {
    let current: any = stack.pop();
    let original = current.original;
    let copy = current.copy;

    for (let key in original) {
      if (original.hasOwnProperty(key)) {
        let value = original[key];

        if (typeof value === "object" && value !== null) {
          copy[key] = Array.isArray(value) ? [] : {};

          stack.push({
            original: value,
            copy: copy[key]
          });
        } else {
          copy[key] = value;
        }
      }
    }
  }
  return cloned;
}

/**
 * 获取浏览器默认语言
 * @returns 语言类型
 */
export const webDefaultLanguage = () => navigator.language;

/**
 * 时间戳转 年月日时分秒
 * 若timestamp不传则取当前时间
 * 若type不传则取 年月日时分秒
 * @param { number } timestamp 时间戳
 * @return 返回年月日时分秒字符串
 */
export const getTimestamp = (timestamp: string | number | null, type: string) => {
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
 * 根据当前日期前推指定日期
 * @param { number } days 需要前推的天数
 * @return {Array[]} 是否为空对象 [前推天数的日期，当前日期]
 */
export const getDatesForwardDate = (days = 0) => {
  const today = new Date();
  const firstDay: any = new Date(today);
  firstDay.setDate(firstDay.getDate() - days); // 向前推指定天，得到前指定天数的第一天

  const lastDay: any = new Date(today);
  lastDay.setDate(lastDay.getDate() - 1); // 昨天是前指定天数的最后一天

  const firstDayFormatted = getTimestamp(firstDay, "YYYY-MM-DD");
  const lastDayFormatted = getTimestamp(lastDay, "YYYY-MM-DD");
  return [firstDayFormatted, lastDayFormatted];
};

/**
 * 给formData循环添加参数，过滤null、undefined、空字符串、NaN
 * 示例：let data = appendFormData(your-object);
 * @param { object } obj 参数对象
 * @return 返回formData对象
 */
export const appendFormData = (obj: any) => {
  let formData = new FormData();
  function deepAppendFormData(formData: any, data: any, parentKey = "") {
    if (Array.isArray(data) || (typeof data === "object" && data !== null)) {
      // 如果数据是数组或对象，序列化为 JSON 字符串
      formData.append(parentKey, JSON.stringify(data));
    } else if (data !== null && data !== undefined && !Number.isNaN(data) && data !== "") {
      // 如果数据是基本类型，直接添加
      formData.append(parentKey, data);
    }
  }
  for (let key in obj) {
    if (obj.hasOwnProperty(key)) {
      deepAppendFormData(formData, obj[key], key);
    }
  }
  deepAppendFormData(formData, obj);
  return formData;
};

/**
 * 判断是否为空对象
 * @param {object} obj 对象
 * @returns {boolean} 是否为空对象
 */
export const isEmptyObject = (obj: object) => {
  // 校验是否为对象且不为 null
  if (typeof obj !== "object" || obj === null) {
    return false;
  }
  return Object.keys(obj).length === 0 && obj.constructor === Object;
};

/**
 * 判断是否为安全环境https或localhost
 * @returns 是否为安全环境https或localhost
 */
export const isSecureEnvironment = () => {
  const { protocol, hostname } = window.location;

  // 检查协议是否为HTTPS
  const isHttps = protocol === "https:";

  // 检查主机名是否为localhost或本地IP（支持IPv4和IPv6）
  const isLocalhost =
    hostname === "localhost" || hostname === "127.0.0.1" || hostname === "0.0.0.0" || hostname === "[::1]" || hostname === "::1";

  return isHttps || isLocalhost;
};

/**
 * 获取不同路由模式所对应的 url + params
 * @returns {String} 返回路径
 */
export const getUrlWithParams = () => {
  type UrlParams = {
    hash: string;
    history: string;
  };

  if (typeof window === "undefined") return "";
  const url: UrlParams = {
    hash: window.location.hash.substring(1),
    history: window.location.pathname + window.location.search
  };
  const mode = (import.meta.env.VITE_ROUTER_MODE as keyof UrlParams) || "history";
  return url[mode] || "";
};

/**
 * 下划线转驼峰
 * @param v 下划线字符串
 * @returns 驼峰字符串
 */
export const getPascalCase = (v: string) => {
  // 处理空字符串情况
  if (!v) return v;

  // 1. 替换所有下划线序列及其后的首个字符
  // - 正则 /_+(.)/g 匹配一个或多个下划线后紧跟的任意字符
  // - 替换函数将匹配的字符转为大写（字母则大写，非字母则原样保留）
  let result = v.replace(/_+(.)/g, (_, char) => char.toUpperCase());

  // 2. 删除开头和结尾剩余的下划线
  result = result.replace(/^_+|_+$/g, "");

  // 3. 将首字母转为大写后返回
  return result.charAt(0).toUpperCase() + result.slice(1);
};
