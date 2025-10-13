/**
 * 小驼峰转小写下划线 userNameInfo
 * @param { string } v 需要转换的字符串，如：userNameInfo
 * @return 分割后的字符串 user_name_info
 */
export const toUnderline = (v: string) => {
  return v.replace(/[A-Z]/g, current => `_${current.toLowerCase()}`);
};

/**
 * 下划线命名转小驼峰 user_name_info
 * @param { string } v 需要转换的字符串，如：user_name_info
 * @return 转换后的字符串userNameInfo
 */
export const getCamelCase = (v: string) => {
  return v.replace(/_[a-z]/g, current => current.split("_")[1].toUpperCase());
};

/**
 * 生成随机Hex颜色
 * @return 返回色值，如：#1eb31
 */
export const getColorHexColor = () => `#${Math.floor(Math.random() * 0xfffff).toString(16)}`;

/**
 * 判断变量数据类型
 * @param { any } val 需要判断数据类型的变量
 * @return 数据类型
 */
export const getObjType = (val: any) => {
  const toString = Object.prototype.toString;
  const map: any = {
    "[object Boolean]": "boolean",
    "[object Number]": "number",
    "[object String]": "string",
    "[object Function]": "function",
    "[object Array]": "array",
    "[object Date]": "date",
    "[object RegExp]": "regExp",
    "[object Undefined]": "undefined",
    "[object Null]": "null",
    "[object Object]": "object",
    "[object Symbol]": "symbol"
  };
  // 如果是节点
  if (val instanceof Element) {
    return "element";
  }
  // 通过toString.call 判断是哪个类型
  // 判断的值为类型[object Boolean]， 通过对象取值返回
  return map[toString.call(val)];
};

/**
 * 基本数据类型数组去重
 * @param { array } v 需要去重的数组
 * @return 去重后的数组
 */
export const arrSet = (v: any) => [...new Set(v)];

/**
 * 数组对象去重
 * @param { array } v 需要去重的数组
 * @return 去重后的数组
 */
export const arrObjSet = (v: any) => {
  return [...new Set(v.map((el: any) => JSON.stringify(el)))].map((el: any) => JSON.parse(el));
};

/**
 * 数组对象根据指定key去重
 * @param { array } v 需要去重的数组
 * @param { string } k 根据指定key去重
 * @return 去重后的数组
 */
export const somethingSet = (v: any, k: any) => {
  let res = new Set();
  return v.filter((item: any) => !res.has(item[k]) && res.add(item[k]));
};

/**
 * 排序-根据条件对数组对象排序
 * v:数组对象，k:key，type:0升序|1降序
 * @param { array } v 需要排序的数组
 * @param { string } k 根据指定key排序
 * @param { number } type 0升序|1降序
 * @return 排序结果
 */
export const toSort = (v: any, k: string, type: number) => {
  v.sort(function (a: any, b: any) {
    if (type == 0) {
      return a[k] - b[k];
    } else {
      return b[k] - a[k];
    }
  });
  return v;
};

/**
 * 比对新对象和旧对象的差异值，返回新对象的差异值
 * @param { object } newObj 新对象
 * @param { object } oldObj 旧对象
 * @return { object } obj 新对象的差异值
 */
export const Difference = (newObj: any, oldObj: any) => {
  let diff: any = {};
  for (let k in newObj) {
    if (newObj[k] != oldObj[k]) {
      diff[k] = newObj[k];
    }
  }
  return diff;
};

/***
 * 删除数组中多个指定元素
 * 根据指定的key比对是否匹配，匹配则删除，最后返回删除后的数组
 * arrMoreDeletion(arr, [2,3,4], 'id')
 * @param {array} arr 原数组
 * @param {array} keys 需要删除的key数组
 * @param {string | number} key 根据指定key删除
 * @return {array} 删除后的数组
 */
export const arrMoreDeletion = (arr: any, keys: any, key: string | number) => {
  const keysSet = new Set(keys);
  return arr.filter((item: any) => !keysSet.has(item[key]));
};

/**
 * 无条件A/B数组交叉合并
 * @param {array} listA 数组A
 * @param {array} listB 数组B
 * @return {array} 数组A、B的交叉数组
 */
export const crossList = (listA = [], listB = []) => {
  const maxLength = Math.max(listA.length, listB.length);
  const result: any = [];

  for (let i = 0; i < maxLength; i++) {
    if (i < listA.length) result.push(listA[i]);
    if (i < listB.length) result.push(listB[i]);
  }

  return result;
};

/***
 * 生成指定长度的随机字符串
 * @param {number} num 需要生成多少位随机字符
 * @return {string} 生成的随机字符
 */
export const randomString = (num: number) => {
  let str = "",
    arr = [
      "0",
      "1",
      "2",
      "3",
      "4",
      "5",
      "6",
      "7",
      "8",
      "9",
      "a",
      "b",
      "c",
      "d",
      "e",
      "f",
      "g",
      "h",
      "i",
      "j",
      "k",
      "l",
      "m",
      "n",
      "o",
      "p",
      "q",
      "r",
      "s",
      "t",
      "u",
      "v",
      "w",
      "x",
      "y",
      "z",
      "A",
      "B",
      "C",
      "D",
      "E",
      "F",
      "G",
      "H",
      "I",
      "J",
      "K",
      "L",
      "M",
      "N",
      "O",
      "P",
      "Q",
      "R",
      "S",
      "T",
      "U",
      "V",
      "W",
      "X",
      "Y",
      "Z"
    ];
  let index = null;
  for (let i = 0; i < num; i++) {
    index = Math.round(Math.random() * (arr.length - 1));
    str += arr[index];
  }
  return str;
};

/***
 * 获取数值在数组中的近似值
 * @param {number[]} arr 数组,如：[23, 30, 35, 47, 16, 21]
 * @param {number} num 当前值，如37
 * @return {string} 当前值在数组中最接近的值
 */
export const closest = (arr: number[], num: number) => {
  let ret = arr[0];
  let distance = Math.abs(ret - num);
  for (let i = 1; i < arr.length; i++) {
    let newDistance = Math.abs(arr[i] - num);
    if (newDistance < distance) {
      distance = newDistance;
      ret = arr[i];
    }
  }
  return ret;
};

/***
 * 根据标准值计算当前值的涨幅和减幅
 * @param {number} current 当前值
 * @param {number} base 标准值
 * @description 返回{ percent: 10, type: 1, text: '涨幅10.00%' }格式
 * @description percent为涨幅或减幅，type: 1为涨幅，2为减幅，0为无变化
 * @return {object} 涨幅/减幅信息
 */
export const getPercent = (current: number, base: number) => {
  // 处理特殊情况
  if (base === 0) {
    return { error: "基数不能为零" };
  }

  if (!current && current !== 0) {
    return { error: "当前值无效" };
  }

  // 计算变化百分比
  const change = ((current - base) / base) * 100;

  // 根据变化值返回结果
  if (change > 0) {
    return {
      percent: change,
      type: 1, // 涨幅
      text: `涨幅${change.toFixed(2)}%`
    };
  } else if (change < 0) {
    return {
      percent: Math.abs(change),
      type: 2, // 减幅
      text: `减幅${Math.abs(change).toFixed(2)}%`
    };
  } else {
    return {
      percent: 0,
      type: 0, // 无变化
      text: "无变化"
    };
  }
};

/**
 * 根据涨幅和减幅计算实际值
 * @param {Number} num 要计算的值
 * @param {Object} extent 包含 type 和 percent 的对象
 * @returns {Number} 计算后的值，保留两位小数
 */
export const getValue = (num: number, extent: { type: number; percent: number }) => {
  if (!extent || (extent.type !== 1 && extent.type !== 2)) return num;

  const percent = extent.percent || 0;
  let result;

  if (extent.type === 1) {
    // 涨幅计算
    result = num * (1 + percent / 100);
  } else {
    // 减幅计算（修复了原代码中的item.value引用错误）
    result = num * (1 - percent / 100);
  }

  // 使用更精确的四舍五入方法
  return Math.round(result * 100) / 100;
};

/**
 * 判断是否为空对象
 * @param {any} obj 对象
 * @returns {boolean} 是否为空对象
 */
export const isEmptyObject = (obj: any) => {
  // 校验是否为对象且不为 null
  if (typeof obj !== "object" || obj === null) {
    return false;
  }
  return Object.keys(obj).length === 0 && obj.constructor === Object;
};

/**
 * 判断是否为安全环境https或localhost
 * @returns {boolean} 是否为安全环境https或localhost
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
 * 获取浏览器默认语言
 * @returns {string} 语言类型
 */
export const webDefaultLanguage = () => navigator.language;

/**
 * 深拷贝
 * @param { any } data 需要深拷贝的数据
 * @returns 深拷贝的数据
 */
export const deepClone = (data: any) => {
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
};

/**
 * 获取URL附带的参数
 * @param { string } url 浏览器链接
 * @returns 解析url?之后的参数转化为object的数据
 */
export const getQueryParams = (url = window.location.href) => {
  const searchParams = new URL(url).searchParams;
  const params: any = {};
  for (const [key, value] of searchParams.entries()) {
    params[key] = value;
  }
  return params;
};

/**
 * 截断字符串，如果字符串超过指定长度，则截断并添加省略号
 * @param { string } str 字符串
 * @param { number } maxLength 截断长度
 * @returns 截断后的字符串
 */
export const truncateString = (str: string, maxLength: number) => {
  if (str.length <= maxLength) return str;
  return str.slice(0, maxLength) + "...";
};

/**
 * 去掉字符串前后空格
 * @param {string} val 字符串
 * @returns {string} 返回处理后的字符串
 */
export const verifyAndSpace = (val: string): string => {
  // 匹配空格
  let v = val.replace(/(^\s*)|(\s*$)/g, "");
  // 返回结果
  return v;
};

/**
 * 去掉字符串的所有中文和空格
 * @param {string} val 当前值字符串
 * @returns {string} 纯非中文不包含空格的字符串
 */
export const verifyCnAndSpace = (val: string): string => {
  // 匹配中文与空格
  let v = val.replace(/[\u4e00-\u9fa5\s]+/g, "");
  // 匹配空格
  v = v.replace(/(^\s*)|(\s*$)/g, "");
  // 返回结果
  return v;
};

/**
 * 去掉前后空格和字符串中的英文
 * @param {string} val 当前值字符串
 * @returns {string} 纯非英文不包含前后空格的字符串
 */
export const verifyEnAndSpace = (val: string): string => {
  // 匹配英文与空格
  let v = val.replace(/[a-zA-Z]+/g, "");
  // 匹配空格
  v = v.replace(/(^\s*)|(\s*$)/g, "");
  // 返回结果
  return v;
};

/**
 * 数字转中文大写
 * @param {any} val 当前值字符串
 * @param {string} unit 默认：仟佰拾亿仟佰拾万仟佰拾元角分
 * @returns {string} 返回处理后的字符串
 */
export const verifyNumberCnUppercase = (val: any, unit = "仟佰拾亿仟佰拾万仟佰拾元角分", v = ""): string => {
  // 处理小数部分，确保能正确转换角分
  val += "00";
  // 返回某个指定的字符串值在字符串中首次出现的位置，没有出现，则该方法返回 -1
  let lookup = val.indexOf(".");
  // substring：不包含结束下标内容，substr：包含结束下标内容
  if (lookup >= 0) val = val.substring(0, lookup) + val.substr(lookup + 1, 2);
  // 根据内容 val 的长度，截取返回对应大写
  unit = unit.substr(unit.length - val.length);
  // 循环截取拼接大写
  for (let i = 0; i < val.length; i++) {
    v += "零壹贰叁肆伍陆柒捌玖".substr(val.substr(i, 1), 1) + unit.substr(i, 1);
  }
  // 正则处理
  v = v
    .replace(/零角零分$/, "整")
    .replace(/零[仟佰拾]/g, "零")
    .replace(/零{2,}/g, "零")
    .replace(/零([亿|万])/g, "$1")
    .replace(/零+元/, "元")
    .replace(/亿零{0,3}万/, "亿")
    .replace(/^元/, "零元");
  // 返回结果
  return v;
};
