/**
 * 处理输入内容，只允许输入数字和2位小数，其它的替换为空
 * 用于输入框限制输入数字和2位小数的场景
 * @param {string} val 当前值字符串
 * @returns {string} 返回处理后的字符串
 */
export const verifyNumberIntegerAndFloat = (val: string): string => {
  // 匹配空格
  let v = val.replace(/(^\s*)|(\s*$)/g, "");
  // 只能是数字和小数点，不能是其他输入
  v = v.replace(/[^\d.]/g, "");
  // 以0开始只能输入一个
  v = v.replace(/^0{2}$/g, "0");
  // 保证第一位只能是数字，不能是点
  v = v.replace(/^\./g, "");
  // 小数只能出现1位
  v = v.replace(".", "$#$").replace(/\./g, "").replace("$#$", ".");
  // 小数点后面保留2位
  v = v.replace(/^(\-)*(\d+)\.(\d\d).*$/, "$1$2.$3");
  // 返回结果
  return v;
};

/**
 * 正整数验证
 * 用于输入框限制正整数的场景
 * @param {string} val 当前值字符串
 * @returns {string} 返回正整数字符串
 */
export const verifiyNumberInteger = (val: string): string => {
  // 匹配空格
  let v = val.replace(/(^\s*)|(\s*$)/g, "");
  // 去掉 '.' , 防止贴贴的时候出现问题 如 0.1.12.12
  v = v.replace(/[\.]*/g, "");
  // 去掉以 0 开始后面的数, 防止贴贴的时候出现问题 如 00121323
  v = v.replace(/(^0[\d]*)$/g, "0");
  // 首位是0,只能出现一次
  v = v.replace(/^0\d$/g, "0");
  // 只匹配数字
  v = v.replace(/[^\d]/g, "");
  // 返回结果
  return v;
};

/**
 * 手机号码格式验证
 * @param {string} val 当前值字符串
 * @returns {boolean} 返回 true: 手机号码正确
 */
export const verifyPhone = (val: string): boolean => {
  // false: 手机号码不正确
  if (!/^((12[0-9])|(13[0-9])|(14[5|7])|(15([0-3]|[5-9]))|(18[0|1,5-9]))\d{8}$/.test(val)) return false;
  // true: 手机号码正确
  else return true;
};

/**
 * 国内电话号码格式验证
 * @param {string} val 当前值字符串
 * @returns {boolean} 返回 true: 国内电话号码正确
 */
export const verifyTelPhone = (val: string): boolean => {
  // false: 国内电话号码不正确
  if (!/\d{3}-\d{8}|\d{4}-\d{7}/.test(val)) return false;
  // true: 国内电话号码正确
  else return true;
};

/**
 * 登录账号验证 (字母开头，允许5-16字节，允许字母数字下划线)
 * @param {string} val 当前值字符串
 * @returns {boolean} 返回 true: 登录账号正确
 */
export const verifyAccount = (val: string): boolean => {
  // false: 登录账号不正确
  if (!/^[a-zA-Z][a-zA-Z0-9_]{4,15}$/.test(val)) return false;
  // true: 登录账号正确
  else return true;
};

/**
 * 密码 (以字母开头，长度在6~16之间，只能包含字母、数字和下划线)
 * @param {string} val 当前值字符串
 * @returns {boolean} 返回 true: 密码正确
 */
export const verifyPassword = (val: string): boolean => {
  // false: 密码不正确
  if (!/^[a-zA-Z]\w{5,15}$/.test(val)) return false;
  // true: 密码正确
  else return true;
};

/**
 * 强密码 (字母+数字+特殊字符，长度在6-16之间)
 * @param {string} val 当前值字符串
 * @returns {boolean} 返回 true: 强密码正确
 */
export const verifyPasswordPowerful = (val: string): boolean => {
  // false: 强密码不正确
  if (
    !/^(?![a-zA-z]+$)(?!\d+$)(?![!@#$%^&\.*]+$)(?![a-zA-z\d]+$)(?![a-zA-z!@#$%^&\.*]+$)(?![\d!@#$%^&\.*]+$)[a-zA-Z\d!@#$%^&\.*]{6,16}$/.test(
      val
    )
  )
    return false;
  // true: 强密码正确
  else return true;
};

/**
 * 密码强度
 * @param {string} val 当前值字符串
 * @description 弱：纯数字，纯字母，纯特殊字符
 * @description 中：字母+数字，字母+特殊字符，数字+特殊字符
 * @description 强：字母+数字+特殊字符
 * @returns {string} 密码强度字符串：弱、中、强
 */
export const verifyPasswordStrength = (val: string): string => {
  let v = "";
  // 弱：纯数字，纯字母，纯特殊字符
  if (/^(?:\d+|[a-zA-Z]+|[!@#$%^&\.*]+){6,16}$/.test(val)) v = "弱";
  // 中：字母+数字，字母+特殊字符，数字+特殊字符
  if (/^(?![a-zA-z]+$)(?!\d+$)(?![!@#$%^&\.*]+$)[a-zA-Z\d!@#$%^&\.*]{6,16}$/.test(val)) v = "中";
  // 强：字母+数字+特殊字符
  if (
    /^(?![a-zA-z]+$)(?!\d+$)(?![!@#$%^&\.*]+$)(?![a-zA-z\d]+$)(?![a-zA-z!@#$%^&\.*]+$)(?![\d!@#$%^&\.*]+$)[a-zA-Z\d!@#$%^&\.*]{6,16}$/.test(
      val
    )
  )
    v = "强";
  // 返回结果
  return v;
};

/**
 * 校验IP地址是否正确
 * @param {string} val 当前值字符串
 * @returns {boolean} 返回 true: IP地址正确
 */
export const verifyIPAddress = (val: string): boolean => {
  // false: IP地址不正确
  if (
    !/^(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])$/.test(
      val
    )
  )
    return false;
  // true: IP地址正确
  else return true;
};

/**
 * 校验邮箱是否正确
 * @param {string} val 当前值字符串
 * @returns {boolean} 返回 true: 邮箱正确
 */
export const verifyEmail = (val: string): boolean => {
  // false: 邮箱不正确
  if (
    !/^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/.test(
      val
    )
  )
    return false;
  // true: 邮箱正确
  else return true;
};

/**
 * 身份证验证
 * @param {string} val 当前值字符串
 * @returns {boolean} 返回 true: 身份证正确
 */
export const verifyIdCard = (val: string): boolean => {
  // false: 身份证不正确
  if (!/^[1-9]\d{5}(18|19|20)\d{2}((0[1-9])|(1[0-2]))(([0-2][1-9])|10|20|30|31)\d{3}[0-9Xx]$/.test(val)) return false;
  // true: 身份证正确
  else return true;
};

/**
 * 中文姓名验证
 * @param {string} val 当前值字符串
 * @returns {boolean} 返回 true: 姓名正确
 */
export const verifyFullName = (val: string): boolean => {
  // false: 姓名不正确
  if (!/^[\u4e00-\u9fa5]{1,6}(·[\u4e00-\u9fa5]{1,6}){0,2}$/.test(val)) return false;
  // true: 姓名正确
  else return true;
};

/**
 * 邮政编码验证
 * @param {string} val 当前值字符串
 * @returns {boolean} 返回 true: 邮政编码正确
 */
export const verifyPostalCode = (val: string): boolean => {
  // false: 邮政编码不正确
  if (!/^[1-9][0-9]{5}$/.test(val)) return false;
  // true: 邮政编码正确
  else return true;
};

/**
 * url正确性验证
 * @param {string} val 当前值字符串
 * @returns {boolean} 返回 true: url 正确
 */
export const verifyUrl = (val: string): boolean => {
  // false: url不正确
  if (
    !/^(?:(?:(?:https?|ftp):)?\/\/)(?:\S+(?::\S*)?@)?(?:(?!(?:10|127)(?:\.\d{1,3}){3})(?!(?:169\.254|192\.168)(?:\.\d{1,3}){2})(?!172\.(?:1[6-9]|2\d|3[0-1])(?:\.\d{1,3}){2})(?:[1-9]\d?|1\d\d|2[01]\d|22[0-3])(?:\.(?:1?\d{1,2}|2[0-4]\d|25[0-5])){2}(?:\.(?:[1-9]\d?|1\d\d|2[0-4]\d|25[0-4]))|(?:(?:[a-z\u00a1-\uffff0-9]-*)*[a-z\u00a1-\uffff0-9]+)(?:\.(?:[a-z\u00a1-\uffff0-9]-*)*[a-z\u00a1-\uffff0-9]+)*(?:\.(?:[a-z\u00a1-\uffff]{2,})).?)(?::\d{2,5})?(?:[/?#]\S*)?$/i.test(
      val
    )
  )
    return false;
  // true: url正确
  else return true;
};

/**
 * 车牌号验证
 * @param {string} val 当前值字符串
 * @returns {boolean} 返回 true：车牌号正确
 */
export const verifyCarNum = (val: string): boolean => {
  // false: 车牌号不正确
  if (
    !/^(([京津沪渝冀豫云辽黑湘皖鲁新苏浙赣鄂桂甘晋蒙陕吉闽贵粤青藏川宁琼使领][A-Z](([0-9]{5}[DF])|([DF]([A-HJ-NP-Z0-9])[0-9]{4})))|([京津沪渝冀豫云辽黑湘皖鲁新苏浙赣鄂桂甘晋蒙陕吉闽贵粤青藏川宁琼使领][A-Z][A-HJ-NP-Z0-9]{4}[A-HJ-NP-Z0-9挂学警港澳使领]))$/.test(
      val
    )
  )
    return false;
  // true：车牌号正确
  else return true;
};
